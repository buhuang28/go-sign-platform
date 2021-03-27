package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/valyala/fastjson"
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

//获取到域名
func GetCpdailyApis(schoolName string) map[string]string {
	defer func() {
		err := recover()
		if err != nil {
			logger.Println(err)
		}
	}()

	apis := make(map[string]string)
	id := GetSchoolId(schoolName)
	if id == "" {
		return nil
	}
	schoolInfo := GetSchoolInfo(id)
	if schoolInfo.IsEmpty() {
		return nil
	}
	idsUrl := schoolInfo.Data[0].IdsURL
	ampUrl := schoolInfo.Data[0].AmpURL
	if strings.Contains(ampUrl, "campusphere") || strings.Contains(ampUrl, "cpdaily") {
		host := GetNetLocol(ampUrl)
		resUrl := GetScheme(ampUrl) + "://" + host
		apis["login-url"] = idsUrl + "/login?service=" + GetScheme(resUrl) + `%3A%2F%2F` + host + `%2Fportal%2Flogin`
		apis["host"] = host
	}

	ampUrl2 := schoolInfo.Data[0].AmpURL2
	if strings.Contains(ampUrl2, "campusphere") || strings.Contains(ampUrl2, "cpdaily") {
		host := GetNetLocol(ampUrl2)
		resUrl := GetScheme(ampUrl2) + "://" + host
		apis["login-url"] = idsUrl + "/login?service=" + GetScheme(resUrl) + `%3A%2F%2F` + host + `%2Fportal%2Flogin`
		apis["host"] = host
	}
	return apis
}

//获取cookie
func GetCookie(user *User, apis map[string]string, loginApi string) string {

	defer func() {
		err := recover()
		if err != nil {
			logger.Println(user.UserName,"出错",err)
		}
	}()


	params := make(map[string]string)
	params["login_url"] = apis["login-url"]
	params["needcaptcha_url"] = ""
	params["captcha_url"] = ""
	params["username"] = user.UserName
	params["password"] = user.PassWord
	//cookies := make(map[string]string)
	sucess, bytes := SendPostForm(loginApi, params)
	logger.Println(user.UserName, "在", loginApi, "ck返回结果:", string(bytes))
	if !sucess {
		logger.Println("GetCookie Post请求失败:",loginApi)
		return ""
	}

	value, err := fastjson.ParseBytes(bytes)
	if err != nil {
		logger.Println("cookie的json解析失败")
		return ""
	}

	ck := string(value.GetStringBytes("cookies"))
	if ck == "None" || ck == "" {
		logger.Println(user.UserName, "在", loginApi, "登录失败,切换子墨登录接口重试")
		sucess, bytes := SendPostForm(zimoApi, params)
		if !sucess {
			logger.Println("Post请求失败")
			return ""
		}

		value, err := fastjson.ParseBytes(bytes)
		if err != nil {
			logger.Println(user.UserName, " cookie的json解析失败")
			return ""
		}

		ck = string(value.GetStringBytes("cookies"))
		if ck == "None" || ck == "" {
			msg := string(value.GetStringBytes("msg"))
			if strings.Contains(msg, "用户名或者密码") {
				logger.Println(user.UserName, "登录信息", msg)
				return "1"
			}
			logger.Println("ck等于None,出错")
			return ""
		}
	}
	return ck
}

//获取签到任务并且签到
func GetScoolSignTasksAndSign(cookie string, apis *map[string]string, user *User) (bool, string) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Println(err)
		}
	}()
	//得到签到路径的api
	api := GetSignInfoApi(apis)

	PostRequest(api, cookie, RequestHeader, nil)
	//这个是为了构造空的json body : {}
	sucess, bytes := PostRequest(api, cookie, RequestHeader, nil)
	if !sucess {
		logger.Println("GetScoolSignTasksAndSign Post网络请求失败")
		return false, "网络波动，导致签到失败"
	}
	value, err := fastjson.ParseBytes(bytes)
	if err != nil {
		logger.Println(user.UserName, "签到任务反序列化失败")
		return false, "无法获取到签到任务"
	}
	datas := value.Get("datas")
	fmt.Println(datas.String())

	unSignedTasks := datas.GetArray("unSignedTasks")
	if unSignedTasks == nil || len(unSignedTasks) < 1 {
		logger.Println(user.UserName, "没有需要签到的任务")
		return false, "没有需要签到的任务"
	}

	for _, v := range unSignedTasks {
		params := make(map[string]string)
		params["signInstanceWid"] = string(v.GetStringBytes("signInstanceWid"))
		params["signWid"] = string(v.GetStringBytes("signWid"))
		task := GetDetailTask(cookie, &params, apis)
		if task.IsEmpty() {
			logger.Println(user.UserName, "空任务详情")
			continue
		}
		//form := FuckForm(task, user)
		form := FuckForm(task, user)
		SubmitForm(cookie, user, form, apis)
	}
	return false, "没有未签到的任务"
}

//获取任务详情
func GetDetailTask(cookie string, params, apis *map[string]string) TaskDeatil {
	api := GetSignTaskDetailApi(apis)

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

//填写表单
func FuckForm(task TaskDeatil, user *User) map[string]interface{} {

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
				if v2.Content == questions[v.Title] {
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

	if task.Datas.IsPhoto == 1 {
		form["signPhotoUrl"] = ""
	}

	form["signInstanceWid"] = task.Datas.SignInstanceWid
	form["longitude"] = user.Longitude
	form["latitude"] = user.Latitude
	form["isMalposition"] = task.Datas.IsMalposition
	form["abnormalReason"] = user.AbnormalReason
	form["position"] = user.Address
	form["uaIsCpadaily"] = true
	return form
}

//填写表单(带图的)
//func FuckForm2(task TaskDeatil, user *User, cookie string, apis *map[string]string) map[string]interface{} {
//	form := make(map[string]interface{})
//	if task.Datas.IsNeedExtra == 1 {
//		extraFields := task.Datas.ExtraField
//		var extraFieldItemValues []map[string]interface{}
//		for _, v := range extraFields {
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
//	if task.Datas.IsPhoto == 1 {
//		list := user.FileList
//		//上传图片到今日校园的oss
//		picMax := len(user.FileList)
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

func SubmitForm(cookie string, user *User, form map[string]interface{}, apis *map[string]string) (bool, string) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Println(user.UserName,"错误了:",err)
		}
	}()

	extension := make(map[string]string)
	extension["lon"] = user.Longitude
	extension["model"] = "Mi 10"
	extension["appVersion"] = "8.2.14"
	extension["systemVersion"] = "10.0"
	extension["userId"] = user.UserName
	extension["systemName"] = "android"
	extension["lat"] = user.Latitude
	extension["deviceId"] = uuid.NewV4().String()

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

	sucess, bytes := PostRequest(GetSubmitSignApi(apis), cookie, header, form)
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
	apis := GetCpdailyApis(schoolName)
	if apis == nil {
		return
	}
	cookie := GetCookie(u, apis, thisApi)
	if cookie == "1" {
		//密码错误的
		logger.Println(u.UserName, "账号密码错误")
		if callBackApi != "" {
			cbData.Status = -1
			cbData.SignResult = "账号密码错误"
			PostRequest(callBackApi, "", nil, cbData)
		}
		return
	}

	if cookie == "" {
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

	sucess, signResult := GetScoolSignTasksAndSign(cookie, &apis, u)
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

func GetSignTaskQA(cookie string, apis *map[string]string, user *User) map[string]map[string][]string {
	//得到签到路径的api
	api := GetSignInfoApi(apis)

	PostRequest(api, cookie, RequestHeader, nil)
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
		task := GetDetailTask(cookie, &params, apis)

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
		task := GetDetailTask(cookie, &params, apis)

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
func UploadPicture(apis *map[string]string, imgName, cookie string) string {
	url := "https://" + (*apis)["host"] + "/wec-counselor-sign-apps/stu/oss/getUploadPolicy"
	params := make(map[string]int)
	params["fileType"] = 1

	sucess, res := PostRequest(url, cookie, nil, params)
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
func GetPic(fileName, cookie string, apis *map[string]string) string {
	url := "https://" + (*apis)["host"] + "/wec-counselor-sign-apps/stu/sign/previewAttachment"
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

	bc.AddFunc(morningEndSepc, func() {
		failSlice = failSlice[0:0]
	})

	bc.AddFunc(noonEndSpec, func() {
		failSlice = failSlice[0:0]
	})

	bc.AddFunc(eveningEndSpec, func() {
		failSlice = failSlice[0:0]
	})
	bc.Start()
}

func CheckTaskStart(date, bgtime string) bool {
	split := strings.Split(date, " ")
	formatTimeStr := split[0] + " " + bgtime + ":00"
	formatTime, _ := time.Parse("2006-01-02 15:04:05", formatTimeStr)
	return time.Now().Unix() > formatTime.Unix()
}
