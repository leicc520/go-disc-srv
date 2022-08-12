如果使用sqlite3数据的话，go需要开启cgo(环境变量:CGO_ENABLED=1)
如果是window环境的话，需要按照gcc/g++以及标准库，可以安装mingw，然后环境变量然后将安装之后的bin目录加入环境变量PATH当中...

服务/配置管理
对应的前端代码:git@github.com:leicc520/go-disc-web.git


# go-配置加载服务发现

go-服务发现以及配置加载接口

静态资源打包处理逻辑
https://github.com/gobuffalo/packr

安装工具箱:
go get -u github.com/gobuffalo/packr/v2@latest

代码中导入"github.com/gobuffalo/packr/v2"
先先写代码装箱BOX
之后运行，生成装箱静态代码编异成二进制代码
packr2
然后再做go build

#!/bin/sh
git pull
go mod tidy

CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build main.go
systemctl reload go-disc.service

本地启动的话需要配置三个环境变量


golang 后台业务管理服务，微服务管理以及服务配置加载管理
默认监听端口 0.0.0.0:7000

环境变量的设置
GOPROXY=https://goproxy.cn,direct;DCENV=dev;DCSRV=127.0.0.1:7000;DCJWT=xxxxxx

由于这个服务比较少变动，这里从管理后台迁移出来, 部署的时候目录发布到/data/web/star_micsrv
同时帮忙创建一个 /data/web/star_micsrv/cachedir cache目录写日志

需要开通内网访问即可
正式绑定域名 disc.xxx-xx.com
测试绑定域名 dev-disc.xxx-xx.com

upstream backend_micsrv {
    ip_hash;
    server 192.168.138.1:7000 max_fails=3 fail_timeout=10s weight=100;
    keepalive 1024;
}

server {
    listen       80;
    server_name  disc.xxx-xx.com;
    root   /home/webroot/xxx-xx.com;
    location / {
        index	index.html;
    }
    location ^~ /micsrv/ {
        proxy_pass        http://backend_micsrv;
        proxy_set_header  Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $remote_addr;
    }
    location ~ \.(png|jpg|jpeg|gif|html|css|js)$ {
        expires max;
        access_log off;
    }
}