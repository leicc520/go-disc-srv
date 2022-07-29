package client

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/leicc520/go-orm"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/leicc520/go-orm/cache"
	"github.com/leicc520/go-orm/log"
)

type MicroClient interface {
	Register(name, srv, proto, version string)
	UnRegister(name, srv string)
	Discover(name string) ([]string, error)
	Reload() error
}

type MicroRegSrv struct {
	RegSrv string
	JwtKey string
}

const (
	DCSRV = "DCSRV"
	DCJWT = "DCJWT"
	DCENV = "DCENV"
)

var gMicroRegSrv *MicroRegSrv = nil
//不设置的话环境变量获取注册地址
func NewMicroRegSrv(srv string) *MicroRegSrv {
	if gMicroRegSrv == nil { //只能生成一个全局即可
		if len(srv) == 0 {
			srv = os.Getenv(DCSRV)
		}
		//非http请求的地址的情况
		if !strings.HasPrefix(srv,"http") {
			srv = "http://"+srv
		}
		token := os.Getenv(DCJWT)
		token = fmt.Sprintf("%x", md5.Sum([]byte(token)))
		gMicroRegSrv = &MicroRegSrv{RegSrv: srv, JwtKey: token}
	}
	log.Write(log.INFO, "register server:{"+gMicroRegSrv.RegSrv+"} token:{"+gMicroRegSrv.JwtKey+"}")
	return gMicroRegSrv
}

//检测微服务端状态 --最多尝试三次
func (a *MicroRegSrv) Health(nTry int, protoSt, srv string) bool {
	status := false
	for i := 0; i < nTry; i++ {
		status = a.zHealth(protoSt, srv)
		if status {//状态检测到的情况
			break
		}
	}
	return status
}

//检测微服务端状态
func (a *MicroRegSrv) zHealth(protoSt, srv string) bool {
	if protoSt == "grpc" {//GRPC服务的心跳

	} else {//处理HTTP的心跳

	}
	return true
}

//发起一个网络请求
func (a *MicroRegSrv) _request(url string, body []byte, method string) (result []byte) {
	var sp *http.Response = nil
	defer func() {//补货异常的处理逻辑
		if sp != nil && sp.Body != nil {
			sp.Body.Close()
		}
		if r := recover(); r != nil {
			log.Write(log.ERROR, "request url ", url, "error", r)
			result = nil
		}
	}()
	log.Write(log.INFO, url, string(body), a.JwtKey)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Write(log.ERROR, url, err, string(body))
		return nil
	}
	req.Header.Set("X-TOKEN", a.JwtKey)
	req.Header.Set("Content-Type", "application/json")
	client  := &http.Client{Timeout: time.Second*3}
	sp, err = client.Do(req)
	if err != nil || sp == nil || sp.StatusCode != http.StatusOK {
		log.Write(log.ERROR, url, err, string(body))
		return nil
	}
	if result, err = ioutil.ReadAll(sp.Body); err != nil {
		log.Write(log.ERROR, url, err, string(body))
		return nil
	} else {
		return result
	}
}

//提交申请接口注册处理逻辑 返回注册的服务地址
func (a *MicroRegSrv) Register(name, srv, protoSt, version string) string {
	req := map[string]interface{}{"name":name,
		"srv":srv, "proto":protoSt, "version":version}
	body, _ := json.Marshal(req)
	data := a._request(a.RegSrv+"/micsrv/register", body, "POST")
	if data == nil || len(data) < 3 {//请求返回异常直接panic
		panic(errors.New("Register microsrv{"+name+"} error"))
	}
	srvAddr := struct {
		Code int64
		Srv string
	}{}
	if err := json.Unmarshal(data, &srvAddr); err != nil || srvAddr.Code != 0 {
		log.Write(log.FATAL, "Register",  err, srvAddr)
		panic(errors.New("Register microsrv{"+name+"} error"))
	}
	log.Write(-1, "register server{"+name+"} success")
	return srvAddr.Srv
}

//提交申请注销微服务的处理逻辑
func (a *MicroRegSrv) UnRegister(protoSt, name, srv string)  {
	req := map[string]string{"name":name, "proto":protoSt, "srv":srv}
	body, _ := json.Marshal(req)
	data := a._request(a.RegSrv+"/micsrv/unregister", body, "POST")
	log.Write(log.INFO, "unregister server{"+name+"-->"+srv+"}-{"+string(data)+"} success")
}

//提交请求申请微服务发现逻辑
func (a *MicroRegSrv) Discover(protoSt, name string) ([]string, error) {
	url := a.RegSrv+"/micsrv/discover/"+protoSt+"/"+name
	data := a._request(url, nil, "GET")
	if data == nil {//服务异常的情况
		return nil, errors.New("发现服务异常,无法获得数据.")
	}
	log.Write(log.INFO, url, string(data))
	srvs := struct {
		Code int64 	`json:"code"`
		Msg  string `json:"msg"`
		Srvs []string `json:"srvs"`
	}{}
	if err := json.Unmarshal(data, &srvs); err != nil || srvs.Code != 0 {
		return nil, err
	}
	return srvs.Srvs, nil
}

//加载配置数据资料信息
func (a *MicroRegSrv) _config(name string) string {
	url := a.RegSrv+"/micsrv/config/"+name
	data := a._request(url, nil, "GET")
	if data == nil || len(data) < 1 {//服务异常的情况
		return ""
	}
	log.Write(log.INFO, url, string(data))
	item := struct {
		Code int64 	`json:"code"`
		Msg  string `json:"msg"`
		Yaml string `json:"yaml"`
	}{}
	if err := json.Unmarshal(data, &item); err != nil || item.Code != 0 {
		return ""
	}
	return item.Yaml
}

//获取微服务配置管理 配置写文件缓存
func (a *MicroRegSrv) Config(name string) string {
	cache := cache.NewFileCache("./cachedir", 1)
	if yaml := a._config(name); len(yaml) > 0 {
		cache.Set("config@"+name, yaml, 0)
		return yaml
	}
	item := cache.Get("config@"+name)
	if item != nil {//数据不为空的情况
		if yaml, ok := item.(string); ok && len(yaml) > 0 {
			return yaml
		}
	}
	log.Write(log.ERROR, "load Config {"+name+"} failed")
	panic("load Config {"+name+"} failed")
}

//提交请求申请微服务发现逻辑
func (a *MicroRegSrv) Reload() error {
	data := a._request(a.RegSrv+"/micsrv/reload", nil, "GET")
	if data == nil {//服务异常的情况
		return errors.New("重启服务异常,无法获得数据.")
	}
	log.Write(log.INFO, a.RegSrv+"/micsrv/reload")
	log.Write(log.INFO, string(data))
	srvs := struct {
		Code int64 	`json:"code"`
		Msg  string `json:"msg"`
	}{}
	if err := json.Unmarshal(data, &srvs); err != nil || srvs.Code != 0 {
		return errors.New(srvs.Msg)
	}
	return nil
}

//申请获取微服务注册的地址信息
func MicSrvGrpcServer(regsrv, protoSt, name string) string {
	rs := NewMicroRegSrv(regsrv)
	cache := orm.GetMCache()
	//通过注册服务 获取数据资料信息 且记录到内存当中 失败的时候取
	srvs, err, ckey := []string{}, errors.New(""), "grpc@"+name
	if srvs, err = rs.Discover(protoSt, name); err != nil || len(srvs) < 1 {
		log.Write(log.ERROR, "grpc 服务地址获取异常{", name, "},通过cache检索")
		data := cache.Get(ckey)
		if data != nil {//数据不为空的情况
			srvs, _ = data.([]string)
		}
	} else {//数据获取成功的情况
		cache.Set(ckey, srvs, 0)
	}
	nidx := len(srvs) - 1
	if nidx > 0 {//大于2条记录做负载均衡
		nidx = int(time.Now().Unix()) % len(srvs)
	}
	if nidx >= 0 {
		return srvs[nidx]
	}
	return ""
}