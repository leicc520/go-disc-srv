package discCTRL

import (
	"github.com/leicc520/go-orm"
	"os"
	"sync"
	"time"

	"git.cht-group.net/leicc/go-orm"
	"git.cht-group.net/leicc/go-orm/log"
	"git.cht-group.net/leicc/go-disc-srv/app/models/sys"
)

type micsrvNodeSt struct {
	Id      int64  `json:"id"`
	Status  int8   `json:"status"`
	Name    string `json:"name" binding:"required"`
	Srv     string `json:"srv"  binding:"required"`
	Proto   string `json:"proto" binding:"required"`
	Version string `json:"version" binding:"required"`
}

type MicSrvMapSt map[string]map[string]int64

//定义存储服务发现存取的设计
type MicSrvPoolSt struct {
	mOnce sync.Once
	mPool MicSrvMapSt
	mLock sync.RWMutex
}

var gGrpcPools *MicSrvPoolSt = nil
var gHttpPools *MicSrvPoolSt = nil

func init() {
	gGrpcPools = &MicSrvPoolSt{mPool: make(MicSrvMapSt), mLock: sync.RWMutex{}, mOnce: sync.Once{}}
	gHttpPools = &MicSrvPoolSt{mPool: make(MicSrvMapSt), mLock: sync.RWMutex{}, mOnce: sync.Once{}}
}

//执行心跳检测 处理逻辑
func (s *MicSrvPoolSt) checkLoop(proto string) {
	sorm := sys.NewSysMsrv()
	smap := orm.SqlMap{"status": 2, "stime": time.Now().Unix()}
	log.Write(log.INFO, "start check {"+proto+"} server")
	for { //开启一个协程循环执行检测任务
		for sname, srvs := range s.mPool {
			for srv, oldid := range srvs {
				//状态不一致的情况删除节点 更新db 重复三次都是异常
				if !regSrv.Health(3, proto, srv) {
					s.Del(sname, srv)
					if oldid > 0 { //记录ID大于0的情况
						smap["stime"] = time.Now().Unix()
						sorm.Save(oldid, smap)
					}
					log.Write(log.INFO, "check server health {"+sname+"} -->"+srv+" status:ERROR")
				} else {
					log.Write(log.INFO, "check server health {"+sname+"} -->"+srv+" status:OK")
				}
			}
		}
		time.Sleep(time.Second * 60)
	}
}

//加载数据库的最新数据完成初始化
func (s *MicSrvPoolSt) Load(proto string) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	sorm := sys.NewSysMsrv()
	list := sorm.GetList(0, -1, func(st *orm.QuerySt) string {
		st.Where("proto", proto)
		st.OrderBy("status", orm.ASC).OrderBy("stime", orm.DESC)
		return st.GetWheres()
	}, "id,name,srv,status,proto,version")
	node := micsrvNodeSt{}
	for _, msrv := range list {
		if err := msrv.ToStruct(&node); err != nil || node.Id < 0 {
			log.Write(log.ERROR, err)
			continue
		}
		//状态异常的情况 且检测不到心跳的情况
		if node.Status != 1 && !regSrv.Health(1, node.Proto, node.Srv) {
			if os.Getenv("DCLOC") != "loc" {//本地执行的情况
				sorm.Delete(node.Id) //移除记录
			}
			continue
		}
		if node.Status != 1 { //更新重置状态
			sorm.Save(node.Id, orm.SqlMap{"status":1, "stime":time.Now().Unix()})
		}
		if _, ok := s.mPool[node.Name]; !ok {
			s.mPool[node.Name] = make(map[string]int64)
		}
		s.mPool[node.Name][node.Srv] = node.Id
	}
	s.mOnce.Do(func() { //只要启动一个执行检测示例即可
		time.AfterFunc(time.Second, func() {
			go s.checkLoop(proto)
		})
	})
}

//添加一个记录到内存当中
func (s *MicSrvPoolSt) Put(name, srv string, id int64) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	if _, ok := s.mPool[name]; !ok {
		s.mPool[name] = make(map[string]int64)
	}
	s.mPool[name][srv] = id
}

//获取服务列表
func (s *MicSrvPoolSt) Get(name string) []string {
	s.mLock.RLock()
	defer s.mLock.RUnlock()
	if _, ok := s.mPool[name]; !ok {
		return nil
	}
	srv := make([]string, 0)
	for osrv, oldid := range s.mPool[name] {
		if oldid > 0 { //附加到节点
			srv = append(srv, osrv)
		}
	}
	return srv
}

//删除指定服务器的一个节点
func (s *MicSrvPoolSt) Del(name, srv string) {
	s.mLock.Lock()
	defer s.mLock.Unlock()
	if _, ok := s.mPool[name]; ok {
		if _, ok := s.mPool[name][srv]; ok {
			delete(s.mPool[name], srv)
		}
	}
}
