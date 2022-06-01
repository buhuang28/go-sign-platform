package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robfig/cron"
	uuid "github.com/satori/go.uuid"
	"github.com/valyala/fastjson"
	"net/http"
	"strings"
	"time"
)

//获取学校ID
func GetSchoolId(schoolName string) string {
	schoolsData, sucess := GetRequest(getSchoolListApi, nil, nil)
	if sucess {
		bytes, err := fastjson.ParseBytes(schoolsData)
		if err != nil {
			logger.Println(schoolName, "解析schoolsData出错")
			return ""
		}
		errCode := bytes.GetInt("errCode")
		if errCode != 0 {
			logger.Println("GetSchoolId状态码出错")
			return ""
		}
		data := bytes.GetArray("data")
		for _, v := range data {
			name := string(v.GetStringBytes("name"))
			if schoolName == name {
				id := string(v.GetStringBytes("id"))
				return id
			}
		}
	} else {
		logger.Println("GetSchoolId网络请求失败")
	}
	return ""
}

//获取学校信息
func GetSchoolInfo(id string) SchoolInfo {
	var schoolInfo SchoolInfo
	urlParams := make(map[string]string)
	urlParams["ids"] = id

	request, sucess := GetRequest(getSchoolInfoApi, nil, urlParams)
	if !sucess {
		logger.Println("GetSchoolInfo网络请求失败")
		return schoolInfo
	}
	err := json.Unmarshal(request, &schoolInfo)
	if err != nil {
		logger.Println(request, "json序列化失败")
		return schoolInfo
	}
	return schoolInfo
}

//var (
//	hostReg = regexp.MustCompile(`\w{4,5}\:\/\/.*?\/`)
//)

//获取到域名
//func GetCpdailyApis(schoolName string) map[string]string {
//	defer func() {
//		err := recover()
//		if err != nil {
//			fmt.Println(err)
//			logger.Println(err)
//		}
//	}()
//
//	apis := make(map[string]string)
//	id := GetSchoolId(schoolName)
//	if id == "" {
//		return nil
//	}
//	schoolInfo := GetSchoolInfo(id)
//	if schoolInfo.IsEmpty() {
//		return nil
//	}
//	//idsUrl := schoolInfo.Data[0].IdsURL
//	ampUrl := schoolInfo.Data[0].AmpURL
//	ampUrl2 := schoolInfo.Data[0].AmpURL2
//
//	loginUrl := ""
//	if strings.Contains(ampUrl, "campusphere") {
//		loginUrl = ampUrl
//	} else if strings.Contains(ampUrl2, "campusphere") {
//		loginUrl = ampUrl2
//	}
//	loginUrl, ck := GetLocation2(loginUrl, nil)
//	apis["login-url"] = loginUrl
//	host := GetRegData(loginUrl)
//	apis["login-host"] = host
//	apis["cookie"] = ck
//	fmt.Println(ck)
//	//
//	//if strings.Contains(ampUrl, "campusphere") || strings.Contains(ampUrl, "cpdaily") {
//	//	host := GetRegData(ampUrl)
//	//	apis["host"] = host
//	//
//	//	sucess, location := GetLocation(ampUrl, nil)
//	//	if !sucess {
//	//		time.Sleep(time.Second)
//	//		sucess, location = GetLocation(ampUrl, nil)
//	//	}
//	//	if sucess && location != "" {
//	//		ampUrl = location
//	//	}
//	//	apis["login-url"] = ampUrl
//	//
//	//	loginHost := GetRegData(ampUrl)
//	//	apis["login-host"] = loginHost
//	//	//loginHost := GetRegData
//	//	//resUrl := GetScheme(ampUrl) + "://" + host
//	//	//apis["login-url"] = idsUrl + "/login?service=" + GetScheme(resUrl) + `%3A%2F%2F` + host + `%2Fportal%2Flogin`
//	//	//apis["host"] = host
//	//}
//	//
//	//if strings.Contains(ampUrl2, "campusphere") || strings.Contains(ampUrl2, "cpdaily") {
//	//	//host := GetNetLocol(ampUrl2)
//	//	host := GetRegData(ampUrl2)
//	//	apis["host"] = host
//	//	host = GetNetLocol(host)
//	//	resUrl := GetScheme(ampUrl2) + "://" + host
//	//	apis["login-url"] = idsUrl + "/login?service=" + GetScheme(resUrl) + `%3A%2F%2F` + host + `%2Fportal%2Flogin`
//	//	apis["login-host"] = GetRegData(apis["login-url"])
//	//}
//	return apis
//}

//获取cookie和host
func GetCookieAndHost(userName, passWord, schoolName, loginApi string) (bool, string, string) {
	data := make(map[string]string)
	data["user_name"] = userName
	data["pass_word"] = passWord
	data["school_name"] = schoolName
	success, bytes := PostRequest(loginApi, "", nil, data)
	if !success {
		return false, "", ""
	}
	parseBytes, err := fastjson.ParseBytes(bytes)
	if err != nil {
		return false, "", ""
	}
	if string(parseBytes.GetStringBytes("cookie")) == "" || string(parseBytes.GetStringBytes("host")) == "" {
		return false, "", ""
	}
	return true, string(parseBytes.GetStringBytes("cookie")), string(parseBytes.GetStringBytes("host"))
}

//获取cookie
func GetCookie(user *User, apis map[string]string, loginApi string) (ck string) {
	ck = ""
	defer func() {
		err := recover()
		if err != nil {
			logger.Println(user.UserName, "出错", err)
			ck = ""
		}
	}()
	loginUrl := apis["login-url"]
	header := make(map[string]string)
	header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36"
	_, realUrl := SchoolDayGetLocationCookie(loginUrl, header)
	params := make(map[string]string)
	params["login_url"] = realUrl
	params["needcaptcha_url"] = ""
	params["captcha_url"] = ""
	params["username"] = user.UserName
	params["password"] = user.PassWord
	params["cookie"] = apis["cookie"]
	//cookies := make(map[string]string)
	sucess, bytes := SendPostForm(loginApi, params)
	logger.Println(user.UserName, "在", loginApi, "ck返回结果:", string(bytes))
	if !sucess {
		logger.Println("GetCookie Post请求失败:", loginApi)
		return ck
	}

	value, err := fastjson.ParseBytes(bytes)
	if err != nil {
		logger.Println("cookie的json解析失败")
		return ck
	}

	ck = string(value.GetStringBytes("cookies"))
	if ck == "None" || ck == "" {
		logger.Println(user.UserName, "在", loginApi, "登录失败,切换备用登录接口重试")
		sucess, bytes := SendPostForm(backupApi, params)
		if !sucess {
			logger.Println("Post请求失败")
			return ck
		}

		value, err := fastjson.ParseBytes(bytes)
		if err != nil {
			logger.Println(user.UserName, " cookie的json解析失败")
			return ck
		}

		ck = string(value.GetStringBytes("cookies"))
		if ck == "None" || ck == "" {
			msg := string(value.GetStringBytes("msg"))
			if strings.Contains(msg, "用户名或者密码") {
				logger.Println(user.UserName, "登录信息", msg)
				return "1"
			}
			logger.Println("ck等于None,出错")
			return ck
		}
	}
	return ck
}

//获取签到任务并且签到
func GetScoolSignTasksAndSign(cookie, host string, user *User) (bool, string) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Println(err)
		}
	}()
	//得到签到路径的api
	api := GetSignInfoApi(host)
	//这个是为了构造空的json body : {}
	header := make(map[string]string)
	//header["User-Agent"] = "Mozilla/5.0 (Linux; U; Android 8.1.0; zh-cn; BLA-AL00 Build/HUAWEIBLA-AL00) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.132 MQQBrowser/8.9 Mobile Safari/537.36"
	header["Content-Type"] = "application/json"
	header["Accept-Encoding"] = "identity"
	header["Content-Length"] = "2"
	header["User-Agent"] = "Mozilla/5.0 (Linux; U; Android 8.1.0; zh-cn; BLA-AL00 Build/HUAWEIBLA-AL00) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.132 MQQBrowser/8.9 Mobile Safari/537.36"
	//PostRequest(api, cookie, RequestHeader, nil)
	header2 := make(map[string]string)
	header2["Content-Type"] = "application/json"
	header2["Cookie"] = cookie
	realCookie, _ := SchoolDayGetLocationCookie(api, header2)
	header["Cookie"] = realCookie
	cookie = realCookie

	sucess, bytes := PostRequest(api, cookie, header, nil)
	if !sucess {
		logger.Println("GetScoolSignTasksAndSign Post网络请求失败")
		return false, "网络波动，导致签到失败"
	}
	value, err := fastjson.ParseBytes(bytes)
	if err != nil {
		logger.Println(user.UserName, "签到任务反序列化失败:", string(bytes))
		return false, "无法获取到签到任务"
	}
	datas := value.Get("datas")

	unSignedTasks := datas.GetArray("unSignedTasks")
	if unSignedTasks == nil || len(unSignedTasks) < 1 {
		logger.Println(user.UserName, "没有需要签到的任务")
		return false, "没有需要签到的任务"
	}

	for _, v := range unSignedTasks {
		params := make(map[string]string)
		params["signInstanceWid"] = string(v.GetStringBytes("signInstanceWid"))
		params["signWid"] = string(v.GetStringBytes("signWid"))
		task := GetDetailTask(realCookie, &params, host)
		//if task.IsEmpty() {
		//	logger.Println(user.UserName, "空任务详情")
		//	continue
		//}
		//form := FuckForm(task, user)
		form := FuckForm(host, cookie, task, user)
		//form := FuckForm2(task, user,cookie,apis)
		SubmitForm(realCookie, user, form, host)
	}
	return false, "没有未签到的任务"
}

//获取任务详情
func GetDetailTask(cookie string, params *map[string]string, host string) TaskDeatil {
	api := GetSignTaskDetailApi(host)

	sucess, bytes := PostRequest(api, cookie, RequestHeader, params)

	var task TaskDeatil

	if !sucess {
		return task
	}
	err := json.Unmarshal(bytes, &task)
	if err != nil {
		logger.Println(bytes, "GetDetailTask反序列化失败")
		return task
	}
	return task
}

//填写表单 --不需要图片的
func FuckForm(host, cookie string, task TaskDeatil, user *User) map[string]interface{} {
	form := make(map[string]interface{})
	if task.Datas.IsNeedExtra == 1 {
		extraFields := task.Datas.ExtraField
		var extraFieldItemValues []map[string]interface{}
		for _, v := range extraFields {
			//检测问题是否对得上
			if questions[v.Title] == "" {
				logger.Println("问题对不上:", v.Title)
				return nil
			}
			for _, v2 := range v.ExtraFieldItems {
				extraFieldItemValue := make(map[string]interface{})
				if v2.Content == questions[v.Title] || v2.Value == questions[v.Title] {
					extraFieldItemValue["extraFieldItemValue"] = questions[v.Title]
					extraFieldItemValue["extraFieldItemWid"] = v2.Wid
					extraFieldItemValues = append(extraFieldItemValues, extraFieldItemValue)
				}

				if v2.IsOtherItems == 1 {
					logger.Println("有额外任务")
					logger.Println(task)
					continue
				}
			}
		}
		form["extraFieldItems"] = extraFieldItemValues
	}

	//if task.Datas.IsPhoto == 1 {
	//	form["signPhotoUrl"] = ""
	//}

	if task.Datas.IsPhoto == 1 {
		list := user.FileList
		//上传图片到今日校园的oss
		picMax := len(user.FileList)
		randInt := RandInt64(int64(picMax))
		imgName := "./img/" + list[randInt]
		fileName := UploadPicture(host, imgName, cookie)
		pic := GetPic(fileName, cookie, host)
		if pic == "" {
			fmt.Println("获取不到图片:", fileName)
			return nil
		}
		//从今日校园oss获取图片
		form["signPhotoUrl"] = pic
	}

	form["signInstanceWid"] = task.Datas.SignInstanceWid
	form["longitude"] = user.Longitude
	form["latitude"] = user.Latitude
	form["isMalposition"] = task.Datas.IsMalposition
	form["abnormalReason"] = user.AbnormalReason
	form["position"] = user.Address
	form["uaIsCpadaily"] = true
	form["signVersion"] = "1.0.0"
	return form
}

//填写表单(带图的)
//func FuckForm2(task TaskDeatil, user *User, cookie string, apis *map[string]string) map[string]interface{} {
//	form := make(map[string]interface{})
//	if task.Datas.IsNeedExtra == 1 {
//		extraFields := task.Datas.ExtraField
//		var extraFieldItemValues []map[string]interface{}
//		for _, v := range extraFields {
//			if v.Title == "" {
//				continue
//			}
//			//检测问题是否对得上
//			if questions[v.Title] == "" {
//				logger.Println("问题对不上:", v.Title)
//				return nil
//			}
//			for _, v2 := range v.ExtraFieldItems {
//				extraFieldItemValue := make(map[string]interface{})
//				if v2.Content == questions[v.Title] {
//					extraFieldItemValue["extraFieldItemValue"] = questions[v.Title]
//					extraFieldItemValue["extraFieldItemWid"] = v2.Wid
//					extraFieldItemValues = append(extraFieldItemValues, extraFieldItemValue)
//				}
//
//				if v2.IsOtherItems == 1 {
//					logger.Println("有额外任务")
//					logger.Println(task)
//					continue
//				}
//			}
//		}
//		form["extraFieldItems"] = extraFieldItemValues
//	}
//
//	if task.Datas.TaskType == "1" {
//		list := user.FileList
//		//上传图片到今日校园的oss
//		picMax := len(user.FileList) - 1
//		randInt := RandInt64(int64(picMax))
//		imgName := "./img/" + list[randInt]
//		fileName := UploadPicture(apis,imgName, cookie)
//		fmt.Println(imgName)
//		pic := GetPic(fileName, cookie, apis)
//		fmt.Println(pic)
//		if pic == "" {
//			return nil
//		}
//		//从今日校园oss获取图片
//		form["signPhotoUrl"] = pic
//	}
//
//	form["signInstanceWid"] = task.Datas.SignInstanceWid
//	form["longitude"] = user.Longitude
//	form["latitude"] = user.Latitude
//	form["isMalposition"] = task.Datas.IsMalposition
//	form["abnormalReason"] = user.AbnormalReason
//	form["position"] = user.Address
//	form["uaIsCpadaily"] = true
//	return form
//}

func SubmitForm(cookie string, user *User, form map[string]interface{}, host string) (bool, string) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Println(user.UserName, "错误了:", err)
		}
	}()

	deviceId := uuid.NewV4().String()
	extension := make(map[string]string)
	extension["lon"] = user.Longitude
	extension["model"] = "Mi 10"
	extension["appVersion"] = "8.2.14"
	extension["systemVersion"] = "10.0"
	extension["userId"] = user.UserName
	extension["systemName"] = "android"
	extension["lat"] = user.Latitude
	extension["deviceId"] = deviceId

	header := make(map[string]string)
	header["User-Agent"] = "Mozilla/5.0 (Linux; Android 10.0.0; Mi 10 Build/KTU84P) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/33.0.0.0 Safari/537.36 okhttp/3.12.4"
	header["CpdailyStandAlone"] = "0"
	header["extension"] = "1"

	//Cpdaily-Extension需要将extension的内容转为json然后使用des cbc PKCS5Padding 加密
	marshal, _ := json.Marshal(extension)
	jsonString := strings.ReplaceAll(string(marshal), ":", ": ")
	jsonString = strings.ReplaceAll(jsonString, ",", ", ")
	encrypt, _ := Encrypt([]byte(jsonString), KEY, IV)
	encoded := base64.StdEncoding.EncodeToString(encrypt)

	header["Cpdaily-Extension"] = encoded
	header["Content-Type"] = "application/json; charset=utf-8"
	//header["Accept-Encoding"] = "gzip"
	header["Connection"] = "Keep-Alive"

	submitData := make(map[string]string)
	var (
		appVersion    = "9.0.12"
		systemName    = "android"
		model         = "SEA-AL10"
		calVersion    = "fitstv"
		systemVersion = "11"
		version       = "first_v2"
	)
	submitData["appVersion"] = appVersion
	submitData["systemName"] = systemName
	m, _ := json.Marshal(form)
	aesEncrypt, _ := AESEncrypt(m)
	submitData["bodyString"] = aesEncrypt
	submitData["sign"] = Md5(string(m))
	submitData["model"] = model
	submitData["lon"] = user.Longitude
	submitData["calVersion"] = calVersion
	submitData["systemVersion"] = systemVersion
	submitData["deviceId"] = deviceId
	submitData["userId"] = user.UserName
	submitData["version"] = version
	submitData["lat"] = user.Latitude

	sucess, bytes := PostRequest(GetSubmitSignApi(host), cookie, header, submitData)
	if !sucess {
		logger.Println(user.UserName, "提交任务失败")
		return false, "提交任务失败"
	}
	parseBytes, err := fastjson.ParseBytes(bytes)
	if err != nil {
		logger.Println(user.UserName, "返回json序列化失败:", string(bytes))
		return false, "签到提交表单返回数据不正确"
	}
	message := string(parseBytes.GetStringBytes("message"))
	if message == "SUCCESS" {
		logger.Println(user.UserName, "签到成功")
		return true, "签到成功"
	}
	return false, "签到失败"
}

func Sign(u *User, isFailProcess bool, thisApi string) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Println(err)
		}
	}()
	logger.Println(u.UserName, "开始签到", thisApi)
	var cbData CallBackData
	cbData.UserName = u.UserName

	success, cookie, host := GetCookieAndHost(u.UserName, u.PassWord, schoolName, thisApi)
	if !success || cookie == "" {
		logger.Println(u.UserName, "可能账号密码错误")
		cbData.Status = -1
		if isFailProcess {
			cbData.SignResult = "登录失败,加入续命队列,一段时间后会重新尝试签到"
			logger.Println(u.UserName, "登录失败,加入续命队列")
			failSlice = append(failSlice, u)
		} else {
			cbData.SignResult = "登录教务失败,可能密码错误或者同一时间登录人数太多，系统拥挤"
			logger.Println(u.UserName, "在", thisApi, "登录教务失败,可能密码错误或者同一时间登录人数太多，系统拥挤")
		}
		if callBackApi != "" {
			PostRequest(callBackApi, "", nil, cbData)
		}
		return
	}
	sucess, signResult := GetScoolSignTasksAndSign(cookie, host, u)
	cbData.SignResult = signResult
	if sucess {
		cbData.Status = 0
		logger.Println(u.UserName, signResult)
	} else {
		cbData.Status = -1
		logger.Println(u.UserName, signResult)
	}
	if callBackApi != "" {
		PostRequest(callBackApi, "", nil, cbData)
	}
}

func GetSignTaskQA(cookie, host string, user *User) map[string]map[string][]string {
	//得到签到路径的api
	api := GetSignInfoApi(host)

	//PostRequest(api, cookie, RequestHeader, nil)
	//这个是为了构造空的json body
	var n name
	sucess, bytes := PostRequest(api, cookie, RequestHeader, n)
	if !sucess {
		logger.Println("Post网络请求失败")
		return nil
	}
	logger.Println(string(bytes))
	value, err := fastjson.ParseBytes(bytes)
	if err != nil {
		logger.Println("签到任务反序列化失败")
		return nil
	}
	datas := value.Get("datas")
	unSignedTasks := datas.GetArray("unSignedTasks")
	signedTasks := datas.GetArray("signedTasks")

	//问卷 -- 问题 -- []答案
	tasks := make(map[string]map[string][]string)
	for _, v := range unSignedTasks {
		if tasks[string(v.GetStringBytes("taskName"))] != nil {
			continue
		}

		params := make(map[string]string)
		params["signInstanceWid"] = string(v.GetStringBytes("signInstanceWid"))
		params["signWid"] = string(v.GetStringBytes("signWid"))
		task := GetDetailTask(cookie, &params, host)

		taskTitleAndAnswer := make(map[string][]string)
		for _, v := range task.Datas.ExtraField {
			var answer []string
			for _, v := range v.ExtraFieldItems {
				answer = append(answer, v.Content)
			}
			taskTitleAndAnswer[v.Title] = answer
		}
		tasks[task.Datas.TaskName] = taskTitleAndAnswer
	}

	for _, v := range signedTasks {
		if tasks[string(v.GetStringBytes("taskName"))] != nil {
			continue
		}
		params := make(map[string]string)
		params["signInstanceWid"] = string(v.GetStringBytes("signInstanceWid"))
		params["signWid"] = string(v.GetStringBytes("signWid"))
		task := GetDetailTask(cookie, &params, host)

		taskTitleAndAnswer := make(map[string][]string)
		for _, v := range task.Datas.ExtraField {
			var answer []string
			for _, v := range v.ExtraFieldItems {
				answer = append(answer, v.Content)
			}
			taskTitleAndAnswer[v.Title] = answer
		}
		tasks[task.Datas.TaskName] = taskTitleAndAnswer
	}
	//logger.Println(tasks)
	return tasks
	//marshal, _ := json.Marshal(tasks)
	//logger.Println("json后的任务数据:",string(marshal))
}

func SignAllUser() {
	dir := ReadDir("./user")
	var users []*User
	for _, v := range dir {
		user := ReadFile("./user/" + v.Name())
		users = append(users, user)
	}
	//每组用户数量
	step := len(users) / len(loginApiList)
	if step == 0 {
		step = 1
	}
	for k, v := range loginApiList {
		lgApi := v
		var signUserSlice []*User
		if (k+1)*step > len(users) {
			continue
		}
		if k != len(loginApiList)-1 {
			signUserSlice = users[k*step : (k+1)*step]
		} else {
			signUserSlice = users[k*step:]
		}
		go func() {
			for _, v2 := range signUserSlice {
				Sign(v2, false, lgApi)
				if SignStepTime != 0 {
					time.Sleep(time.Duration(SignStepTime) * time.Second)
				} else {
					time.Sleep(time.Second * 30)
				}
			}
		}()
	}
}

func SignFallUser(users []*User) {
	step := len(users) / len(loginApiList)
	for k, v := range loginApiList {
		lgApi := v
		var signUserSlice []*User
		if k != len(loginApiList)-1 {
			signUserSlice = users[k*step : (k+1)*step]
		} else {
			signUserSlice = users[k*step:]
		}
		go func() {
			for _, v2 := range signUserSlice {
				Sign(v2, false, lgApi)
				time.Sleep(time.Second * 25)
			}
		}()
	}
}

//上传图片到今日校园的OSS
func UploadPicture(host, imgName, cookie string) string {
	url := host + "wec-counselor-sign-apps/stu/oss/getUploadPolicy"
	params := make(map[string]int)
	params["fileType"] = 1
	h := make(map[string]string)
	h["content-type"] = "application/json"

	sucess, res := PostRequest(url, cookie, h, params)
	if !sucess || res == nil {
		return ""
	}
	logger.Println(string(res))
	value, err := fastjson.ParseBytes(res)
	if err != nil || value == nil {
		return ""
	}
	datas := value.Get("datas")
	fileName := string(datas.GetStringBytes("fileName"))
	policy := string(datas.GetStringBytes("policy"))
	accessKeyId := string(datas.GetStringBytes("accessid"))
	signature := string(datas.GetStringBytes("signature"))
	policyHost := string(datas.GetStringBytes("host"))

	header := make(map[string]string)
	header["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0"

	data := make(map[string]string)
	data["key"] = fileName
	data["policy"] = policy
	data["OSSAccessKeyId"] = accessKeyId
	data["success_action_status"] = "200"
	data["signature"] = signature
	PostMultipartImage(&data, &header, imgName, policyHost)
	return fileName
}

//获取图片
func GetPic(fileName, cookie, host string) string {
	url := host + "wec-counselor-sign-apps/stu/sign/previewAttachment"
	data := make(map[string]string)
	data["ossKey"] = fileName
	sucess, bytes := PostRequest(url, cookie, nil, data)
	if !sucess || bytes == nil {
		return ""
	}
	value, err := fastjson.ParseBytes(bytes)
	if err != nil || value == nil {
		return ""
	}
	photoUrl := string(value.GetStringBytes("datas"))
	return photoUrl
}

func CreateCron() {
	bc.Stop()
	bc = cron.New()

	if MorningSignTime != "" {
		bc.AddFunc(MorningSignTime, SignAllUser)
	}

	if NoonSignTime != "" {
		bc.AddFunc(NoonSignTime, SignAllUser)
	}

	if EveningSignTime != "" {
		bc.AddFunc(EveningSignTime, SignAllUser)
	}

	bc.AddFunc(morningProcessFailSepc, func() {
		if len(failSlice) > 0 {
			SignFallUser(failSlice)
		}
	})

	bc.AddFunc(noonProcessFailSepc, func() {
		if len(failSlice) > 0 {
			SignFallUser(failSlice)
		}
	})

	bc.AddFunc(eveningProcessFailSepc, func() {
		if len(failSlice) > 0 {
			SignFallUser(failSlice)
		}
	})

	//bc.AddFunc(morningEndSepc, func() {
	//	failSlice = failSlice[0:0]
	//})
	//
	//bc.AddFunc(noonEndSpec, func() {
	//	failSlice = failSlice[0:0]
	//})
	//
	//bc.AddFunc(eveningEndSpec, func() {
	//	failSlice = failSlice[0:0]
	//})
	bc.Start()
}

func CheckTaskStart(date, bgtime string) bool {
	split := strings.Split(date, " ")
	formatTimeStr := split[0] + " " + bgtime + ":00"
	formatTime, _ := time.Parse("2006-01-02 15:04:05", formatTimeStr)
	return time.Now().Unix() > formatTime.Unix()
}

func SchoolDayGetLocationCookie(u string, header map[string]string) (string, string) {
	req, _ := http.NewRequest("GET", u, nil)
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	timeout := time.Duration(6 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}
	response, _ := client.Do(req)
	defer func() {
		response.Body.Close()
	}()
	if response == nil {
		return "", ""
	}
	realLc := ""
	realCk := ""
	localtion, _ := response.Location()
	if localtion != nil && localtion.String() != "" {
		ck := ""
		cookies := response.Cookies()
		for _, v := range cookies {
			ck = ck + ";" + v.Name + "=" + v.Value
		}
		ck = strings.Trim(ck, ";")
		if ck != "" {
			header["Cookie"] = header["Cookie"] + ";" + ck
		}
		realCk, realLc = SchoolDayGetLocationCookie(localtion.String(), header)
	} else {
		realLc = response.Request.URL.String()
		realCk = header["Cookie"]
	}
	return realCk, realLc
}
