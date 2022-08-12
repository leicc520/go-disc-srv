package test

import (
	"encoding/json"
	"fmt"
	"github.com/leicc520/go-gin-http"
	"github.com/leicc520/go-gin-http/micro"
	"github.com/leicc520/go-orm"
	"testing"
)

const BaseUrl = "http://127.0.0.1:7000"

func TestCaptcha(t *testing.T) {
	urlStr := BaseUrl + "/captcha/check"
	bodyStr, _ := json.Marshal(orm.SqlMap{"sumid":"870969-T3uRnzOWhaExutiZp3Sv",
		"vcode":"5431"})
	
	fmt.Println(string(bodyStr))
	
	core.NewHttpRequest().SetContentType("json").Request(urlStr, bodyStr, "POST")
}

func TestSignIn(t *testing.T) {
	urlStr := BaseUrl + "/signin/check"
	bodyStr, _ := json.Marshal(orm.SqlMap{"account":"admin",
		"loginpw":"simlife@520", "xtoken":"c145f-1660122586-7cbd3"})
	
	acl := core.NewAcl(0, nil)
	fmt.Println(acl.Crypt("simlife@520"))
	
	//1f9abbabf9926d579a3c5d1140421be8
	fmt.Println(string(bodyStr))
	core.NewHttpRequest().SetContentType("json").Request(urlStr, bodyStr, "POST")
}

func TestConfig(t *testing.T) {
	regsrv := micro.NewRegSrvClient("")
	config := regsrv.Config("go.disc.srv")
	fmt.Println(config)
}

func TestDiscover(t *testing.T) {
	regsrv := micro.NewRegSrvClient("")
	srv, err := regsrv.Discover("http", "go.demov5.srv")
	fmt.Println(srv, err)
}
