package main

var (
	getSchoolListApi = "https://mobile.campushoy.com/v6/config/guest/tenant/list"
	getSchoolInfoApi = "https://mobile.campushoy.com/v6/config/guest/tenant/info"
	zimoApi          = "http://www.zimo.wiki:8080/wisedu-unified-login-api-v1.0/api/login"
	LoginApi2        = "http://152.136.185.60:8001/api/login"
	LoginApi         = "http://47.114.146.19:8001/api/login"
	LoginApi3        = "http://106.52.129.73:8001/api/login"
	LoginApi4        = "http://81.70.164.119:8001/api/login"
	LoginApi5        = "http://127.0.0.1:8001/api/login"
	loginApiList     = []string{LoginApi, LoginApi2, LoginApi3, LoginApi4}
	//定义问题
	//schoolName = "深圳信息职业技术学院"
	//questions = map[string]string{"您的体温是否在37.3°及以上？":"否",
	//	"您是否存在咳嗽、流鼻涕等感冒症状？":"否",
	//	"您是否存在咳嗽，流鼻涕等感冒症状？":"否"}
	schoolName      = ""
	questions       = make(map[string]string)
	callBackApi     = ""
	port            = ""
	MorningSignTime = ""
	NoonSignTime    = ""
	EveningSignTime = ""
	SignStepTime    int64
)

var (
	//定义一个全局的请求头
	RequestHeader = map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36",
		"content-type":    "application/json",
		"Accept-Encoding": "gzip,deflate",
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
