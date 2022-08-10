package app

import (
	"bytes"
	"encoding/base64"
	"net/http"
	
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/leicc520/go-gin-http"
)

//初始化验证码
func init()  {
	core.NewInitCap()
}


// @Summary 生成验证码
// @Description 中台的登录管理平台-生成验证码
// @Tags 通用接口
// @Success 200 {json} HttpView
// @Router /captcha [get]
func doCaptcha(c *gin.Context) {
	sumId, _ := c.Cookie(core.CapCookie)
	idStr := core.Gcaptcha.CheckCaptchaSum(sumId)
	if idStr == "" || !captcha.Reload(idStr) {
		idStr = captcha.NewLen(4)
		sumId = core.Gcaptcha.CaptchaSum(idStr)
		c.SetCookie(core.CapCookie, sumId, 0, "/", "", false, false)
	}
	lang, ext := "zh", ".png" //默认生成图片验证码
	if core.Gcaptcha.Serve(c.Writer, c.Request, idStr, ext, lang, false, 110, 38) == captcha.ErrNotFound {
		http.NotFound(c.Writer, c.Request)
	}
}

// @Summary 验证码跨域
// @Description 中台的登录管理平台-生成验证码
// @Tags 通用接口
// @Param sumid formData string true "验证的对话密钥"
// @Success 200 {json} HttpView
// @Router /captcha/json [post]
func doCaptchaJson(c *gin.Context) {
	args := struct {
		SumId string `form:"sumid"`
	}{}
	if err := c.ShouldBind(&args); err != nil {
		core.PanicValidateHttpError(1001, err)
	}
	idStr := core.Gcaptcha.CheckCaptchaSum(args.SumId)
	if idStr == "" || !captcha.Reload(idStr) { //cookie不存在的情况
		idStr = captcha.NewLen(4)
		args.SumId = core.Gcaptcha.CaptchaSum(idStr)
	}
	var buf bytes.Buffer
	width, height := 110, 38
	if err := captcha.WriteImage(&buf, idStr, width, height); err != nil {
		core.PanicHttpError(6001, "获取验证码图片失败.")
	}
	image := base64.StdEncoding.EncodeToString(buf.Bytes())
	datas := map[string]string{core.CapCookie: args.SumId, "image": image}
	core.NewHttpView(c).JsonDisplay(datas)
}

// @Summary 生成验证码
// @Description 中台的登录管理平台-生成验证码
// @Tags 通用接口
// @Param sumid formData string true "验证码sessionid"
// @Param vcode formData string true "验证码"
// @Success 200 {json} HttpView
// @Router /captcha/check [post]
func doCheckCaptcha(c *gin.Context) {
	args := struct {
		SumId string `form:"sumid" json:"sumid" binding:"required"`
		VCode string `form:"vcode" json:"vcode" binding:"required"`
	}{}
	if err := c.ShouldBind(&args); err != nil {
		core.PanicValidateHttpError(1001, err)
	}
	if len(args.SumId) < 6 { //默认取传过来的参数
		args.SumId, _ = c.Cookie(core.CapCookie)
	}
	if !core.Gcaptcha.CheckSum(args.SumId, args.VCode) {
		core.PanicHttpError(1010)
	}
	tkStr := core.Gcaptcha.GenerateHash(c)
	core.NewHttpView(c).JsonDisplay(gin.H{"xtoken": tkStr})
}
