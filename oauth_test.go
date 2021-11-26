package main

import (
	"fmt"
	"jrxy/sign/entity"
	"jrxy/sign/tool"
	"strconv"
	"testing"
	"time"
)

var (
	APIADDR = "http://127.0.0.1:8800/"
	ADDCODE = APIADDR + "AddOauthCode"
	ACTIVITYOAUTHCODE = APIADDR + "ActiviteOuathCode"
	CHECKOAUTHCODE = APIADDR + "CheckOauthCode"
	AUTHOAUTHCODE = APIADDR + "AuthOauthCode"
	CREATEORDER = APIADDR + "CreateOrder"
	SELECTORDER = APIADDR + "SelectOrder"
)

func TestOauth(t *testing.T)  {

}

func TestAddOauthCode(t *testing.T)  {
	data := make(map[string]interface{})
	data["count"] = 50
	request, bytes := tool.PostRequest(ADDCODE, data)
	fmt.Println(request)
	fmt.Println(string(bytes))
}

func TestActiviteOuathCode(t *testing.T)  {
	acTime := time.Now().Unix()
	data := make(map[string]interface{})
	data["invitor"] = "ZBIFLH"
	data["activity_time"] = acTime
	data["sign"] = tool.GetSign(strconv.FormatInt(acTime,10),"ZBIFLH")
	request, bytes := tool.PostRequest(ACTIVITYOAUTHCODE, data)
	fmt.Println(request)
	fmt.Println(string(bytes))
}

func TestCheckOauthCode(t *testing.T) {
	var oauthData entity.OauthBody
	code := "D1E721AF1BE1EBBA916A3A8FD187989C"
	oauthData.OauthCode = code
	oauthData.OauthKey = tool.GetSign(code)
	oauthData.OauthTime = time.Now().Unix()
	oauthData.OauthSign = tool.GetSign(strconv.FormatInt(oauthData.OauthTime,10),oauthData.OauthKey,oauthData.OauthCode)
	request, bytes := tool.PostRequest(CHECKOAUTHCODE, oauthData)
	fmt.Println(request)
	fmt.Println(string(bytes))
}

func TestAuthOauthCode(t *testing.T)  {
	var oauthData entity.OauthBody
	code := "D1E721AF1BE1EBBA916A3A8FD187989C"
	oauthData.OauthCode = code
	oauthData.OauthKey = tool.GetSign(code+"123")
	oauthData.OauthTime = time.Now().Unix()
	oauthData.OauthSign = tool.GetSign(strconv.FormatInt(oauthData.OauthTime,10),oauthData.OauthKey,oauthData.OauthCode)
	request, bytes := tool.PostRequest(AUTHOAUTHCODE, oauthData)
	fmt.Println(request)
	fmt.Println(string(bytes))
}

func TestProLongOuathCodeTime(t *testing.T)  {

}

var (
	CREATEORDERAPI = "http://127.0.0.1:8800/CreateOrder"
)

//创建订单
func TestCreateOrder(t *testing.T) {
	var data entity.CreateOrderBody
	data.Code = "huaji"
	data.Price = 1
	data.CreateTime = time.Now().Unix()
	data.RandKey = "lalalla"
	getSign := tool.GetSign(data.Code, data.RandKey, strconv.FormatInt(data.Price, 10), strconv.FormatInt(data.CreateTime, 10))
	data.Sign = getSign
	request, bytes := tool.PostRequest(CREATEORDERAPI,  data)
	fmt.Println(request)
	fmt.Println(string(bytes))
}

//查询订单
func TestSelectOrder(t *testing.T) {





}

//回调测试
func TestPayCallBackOrder(t *testing.T)  {

}