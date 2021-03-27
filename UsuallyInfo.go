package main

var (
	getSchoolListApi = "https://mobile.campushoy.com/v6/config/guest/tenant/list"
	getSchoolInfoApi = "https://mobile.campushoy.com/v6/config/guest/tenant/info"
	zimoApi          = "http://www.zimo.wiki:8080/wisedu-unified-login-api-v1.0/api/login"

	loginApiList []string
	schoolName      = ""
	questions       = make(map[string]string)
	callBackApi     = ""
	port            = ""
	MorningSignTime = ""
	NoonSignTime    = ""
	EveningSignTime = ""
	SignStepTime    int64
	ImgPath = ".\\img\\"
)

var (
	//定义一个全局的请求头
	RequestHeader = map[string]string{
		"Accept":       "application/json, text/plain, */*",
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36",
		"content-type": "application/json",
		//"Accept-Encoding": "gzip,deflate",
		"Accept-Language": "zh-CN,en-US;q=0.8",
		"Content-Type":    "application/json;charset=UTF-8"}
)

var (
	IV         = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	KEY        = []byte("b3L26XNL")
	sessionKey = "auto_sign"
)

type name struct {
	ha string
}
