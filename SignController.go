package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/ying32/govcl/vcl"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	addCache      = make(map[string]bool)
	Lock          sync.Mutex //日志锁
	ListViewIndex int64      = 1
)

func WebStart(addr string) {
	var logfile *os.File
	exits := CheckFileIsExits("./gin.log")
	if !exits {
		logfile, _ = os.Create("./gin.log")
	} else {
		//如果存在文件则 追加log
		logfile, _ = os.OpenFile("./gin.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	}
	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = io.MultiWriter(logfile)
	router := gin.Default()
	router.Use(Use)
	router.GET("/index", Index)
	router.POST("/addUser", AddUser)
	router.POST("/getTask", GetTaskService)
	router.POST("/addPic", AddPic)
	router.Run(addr)
}

func Index(context *gin.Context) {
	context.Header("Content-Type", "text/html; charset=utf-8")
	context.String(200, htmlfile)
}

//func AddUser(context *gin.Context) {
//	defer func() {
//		err := recover()
//		if err != nil {
//			logger.Println(err)
//		}
//	}()
//	data, err := ioutil.ReadAll(context.Request.Body)
//	if data == nil || err != nil {
//		context.JSON(200, gin.H{
//			"code":    -1,
//			"message": "数据丢失",
//		})
//		return
//	}
//	var u User
//	err = json.Unmarshal(data, &u)
//	if err != nil {
//		context.JSON(200, Result{Message: "信息序列化失败"})
//		return
//	}
//	Lock.Lock()
//	isCache := addCache[u.UserName]
//	Lock.Unlock()
//	if isCache {
//		context.JSON(200, Result{Message: "请勿重复提交表单"})
//		return
//	}
//	Lock.Lock()
//	addCache[u.UserName] = true
//	Lock.Unlock()
//
//	defer func() {
//		addCache[u.UserName] = false
//	}()
//	u.AbnormalReason = "在家"
//	//进行用户校验
//	result := CheckUserData(&u)
//	if result.Code < 0 {
//		context.JSON(200, result)
//		return
//	}
//	//var result Result
//	//校验MD5和时间戳
//	result = MD5Sign(&u, true)
//	if result.Code < 0 {
//		context.JSON(200, result)
//		return
//	}
//	marshal, err := json.Marshal(u)
//	apis := GetCpdailyApis(schoolName)
//	if apis == nil {
//		result.Message = "找不到该学校"
//		context.JSON(-1, result)
//		return
//	}
//	cookie := GetCookie(&u, apis, loginApiList[0])
//	if cookie == "" || cookie == "1" {
//		time.Sleep(time.Second * 2)
//		cookie = GetCookie(&u, apis, backupApi)
//		if cookie == "" || cookie == "1" {
//			result.Message = "无法登录，可能学号或者密码错误"
//			context.JSON(-2, result)
//			return
//		}
//	}
//	sucess := WriteContent("./user/"+u.UserName+".json", string(marshal))
//	if sucess && err == nil {
//		result.Message = "用户添加成功"
//	} else {
//		result.Message = "请重试一次"
//	}
//	context.JSON(200, result)
//	return
//}

//这个好像是带图片的
func AddUser(context *gin.Context) {
	data, err := ioutil.ReadAll(context.Request.Body)
	if data == nil || err != nil {
		context.JSON(200, gin.H{
			"code":    -1,
			"message": "数据丢失",
		})
		return
	}
	var u User
	err = json.Unmarshal(data, &u)
	if err != nil {
		context.JSON(200, Result{Message: "信息序列化失败"})
		return
	}
	Lock.Lock()
	isCache := addCache[u.UserName]
	Lock.Unlock()
	if isCache {
		context.JSON(200, Result{Message: "请勿重复提交表单"})
		return
	}
	Lock.Lock()
	addCache[u.UserName] = true
	Lock.Unlock()

	defer func() {
		addCache[u.UserName] = false
	}()
	u.AbnormalReason = "在家"
	//进行用户校验
	result := CheckUserData(&u)
	if result.Code < 0 {
		context.JSON(200, result)
		return
	}
	//var result Result
	//校验MD5和时间戳
	result = MD5Sign(&u, true)
	if result.Code < 0 {
		context.JSON(200, result)
		return
	}
	marshal, err := json.Marshal(u)
	apis := GetCpdailyApis(schoolName)
	if apis == nil {
		result.Message = "找不到该学校"
		context.JSON(-1, result)
		return
	}
	cookie := GetCookie(&u, apis, loginApiList[0])
	if cookie == "" || cookie == "1" {
		time.Sleep(time.Second * 2)
		cookie = GetCookie(&u, apis, loginApiList[0])
		if cookie == "" || cookie == "1" {
			result.Message = "无法登录，可能学号或者密码错误"
			context.JSON(-2, result)
			return
		}
	}
	sucess := WriteContent("./user/"+u.UserName+".json", string(marshal))
	if sucess && err == nil {
		result.Message = "用户添加成功"
	} else {
		result.Message = "请重试一次"
	}
	context.JSON(200, result)
	return
}

func GetTaskService(context *gin.Context) {
	data, err := ioutil.ReadAll(context.Request.Body)
	if data == nil || err != nil {
		context.JSON(200, gin.H{
			"code":    -1,
			"message": "数据丢失",
		})
		return
	}
	var u User
	err = json.Unmarshal(data, &u)

	//用户校验
	result := CheckUserData(&u)
	if result.Code < 0 {
		context.JSON(200, result)
		return
	}
	//校验MD5和时间戳
	result = MD5Sign(&u, false)
	if result.Code < 0 {
		context.JSON(result.Code, result)
		return
	}

	apis := GetCpdailyApis(schoolName)
	if apis == nil {
		context.JSON(-1, "找不到该学校")
		return
	}
	cookie := GetCookie(&u, apis, loginApiList[0])

	if cookie == "" {
		context.JSON(-2, "无法登录，可能用户名或者密码错误")
		return
	}
	signTaskQA := GetSignTaskQA(cookie, &apis, &u)
	if signTaskQA == nil {
		context.JSON(-3, "找不到任务")
		return
	}
	context.JSON(200, signTaskQA)
	return
}

func AddPic(context *gin.Context) {
	file, _ := context.FormFile("file")
	filename := strconv.FormatInt(int64(time.Now().UnixNano()), 10) + ".jpg"
	var r Result
	if err := context.SaveUploadedFile(file, "./img/"+filename); err != nil {
		//logger.Println(err.Error())
		//自己完成信息提示
		r.Code = -1
		r.Message = "上传图片失败"
		context.JSON(200, "")
	} else {
		r.Code = 200
		r.Message = "sucess"
		r.Data = filename
		context.JSON(200, filename)
	}
	return
}

func Use(context *gin.Context) {
	start := time.Now().UnixNano() / 1e6
	method := context.Request.Method
	ip := context.ClientIP()
	path := context.Request.URL.String()
	end := time.Now().UnixNano() / 1e6
	useTime := strconv.FormatInt(end-start, 10)
	vcl.ThreadSync(func() {
		index := strconv.FormatInt(ListViewIndex, 10)
		ListViewIndex++
		item := Form1.HttpListView.Items().Add()
		item.SetCaption(index)
		item.SubItems().Add(ip)
		item.SubItems().Add(path)
		item.SubItems().Add(useTime)
		item.SubItems().Add("")
		item.SubItems().Add(method)
	})
}
