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
	schoolsData,sucess:= GetRequest(getSchoolListApi,nil,nil)
	if sucess {
		bytes, err := fastjson.ParseBytes(schoolsData)
		if err != nil {
			fmt.Println("解析schoolsData出错")
			return ""
		}
		errCode := bytes.GetInt("errCode")
		if errCode != 0 {
			fmt.Println("状态码出错")
			return ""
		}
		data := bytes.GetArray("data")
		for _,v := range data {
			name := string(v.GetStringBytes("name"))
			if schoolName == name {
				id := string(v.GetStringBytes("id"))
				return id
			}
		}
	}else {
		fmt.Println("网络请求失败")
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
		fmt.Println("网络请求失败")
		return schoolInfo
	}
	err := json.Unmarshal(request, &schoolInfo)
	if err != nil {
		fmt.Println("json序列化失败")
		return schoolInfo
	}
	return schoolInfo
}

//获取到域名
func GetCpdailyApis(schoolName string) map[string]string {
	apis := make(map[string]string)

	id := GetSchoolId(schoolName)
	schoolInfo := GetSchoolInfo(id)
	if schoolInfo.IsEmpty() {
		return nil
	}
	idsUrl := schoolInfo.Data[0].IdsURL
	ampUrl := schoolInfo.Data[0].AmpURL
	if strings.Contains(ampUrl,"campusphere") || strings.Contains(ampUrl,"cpdaily") {
		host := GetNetLocol(ampUrl)
		resUrl := GetScheme(ampUrl)+"://"+host
		apis["login-url"] = idsUrl + "/login?service="+GetScheme(resUrl)+`%3A%2F%2F`+host+`%2Fportal%2Flogin`
		apis["host"] = host
	}

	ampUrl2 := schoolInfo.Data[0].AmpURL2
	if strings.Contains(ampUrl2,"campusphere") || strings.Contains(ampUrl2,"cpdaily") {
		host := GetNetLocol(ampUrl2)
		resUrl := GetScheme(ampUrl2)+"://"+host
		apis["login-url"] = idsUrl + "/login?service="+GetScheme(resUrl)+`%3A%2F%2F`+host+`%2Fportal%2Flogin`
		apis["host"] = host
	}
	return apis
}

//获取cookie
func GetCookie(user *User,apis map[string]string,loginApi string) string {
	params := make(map[string]string)
	params["login_url"] = apis["login-url"]
	params["needcaptcha_url"] = ""
	params["captcha_url"] = ""
	params["username"] = user.UserName
	params["password"] = user.PassWord
	//cookies := make(map[string]string)
	sucess,bytes := SendPostForm(loginApi, params)
	fmt.Println("ck返回结果:",string(bytes))
	if !sucess {
		fmt.Println("Post请求失败")
		return ""
	}

	value, err := fastjson.ParseBytes(bytes)
	if err != nil {
		fmt.Println("cookie的json解析失败")
		return ""
	}

	ck := string(value.GetStringBytes("cookies"))
	if ck == "None" || ck == ""{
		fmt.Println(loginApi,"登录失败,切换子墨登录接口重试")
		sucess,bytes := SendPostForm(zimoApi, params)
		if !sucess {
			fmt.Println("Post请求失败")
			return ""
		}

		value, err := fastjson.ParseBytes(bytes)
		if err != nil {
			fmt.Println("cookie的json解析失败")
			return ""
		}

		ck = string(value.GetStringBytes("cookies"))
		if ck == "None" || ck == "" {
			msg := string(value.GetStringBytes("msg"))
			if strings.Contains(msg,"用户名或者密码") {
				fmt.Println("登录信息",msg)
				return "1"
			}
			fmt.Println("ck等于None,出错")
			return ""
		}
	}
	return ck
}

//获取签到任务并且签到
func GetScoolSignTasksAndSign(cookie string,apis *map[string]string,user *User) {

	//得到签到路径的api
	api := GetSignInfoApi(apis)

	PostRequest(api, cookie, RequestHeader, nil)
	//这个是为了构造空的json body
	var n name
	sucess, bytes := PostRequest(api, cookie, RequestHeader, n)
	if !sucess {
		fmt.Println("Post网络请求失败")
		return
	}
	fmt.Println(string(bytes))
	value, err := fastjson.ParseBytes(bytes)
	if err != nil {
		fmt.Println("签到任务反序列化失败")
		return
	}
	datas := value.Get("datas")
	fmt.Println("全部任务:",datas.String())
	unSignedTasks := datas.GetArray("unSignedTasks")
	if unSignedTasks == nil || len(unSignedTasks) < 1 {
		fmt.Println("没有需要签到的任务")
		return
	}

	for _,v := range unSignedTasks {
		fmt.Println("签到任务:",v.String())
		params := make(map[string]string)
		params["signInstanceWid"] = string(v.GetStringBytes("signInstanceWid"))
		params["signWid"] = string(v.GetStringBytes("signWid"))
		task := GetDetailTask(cookie, &params, apis)
		if task.IsEmpty() {
			fmt.Println("空任务详情")
			//continue
		}
		form := FuckForm(task, user)
		//form := FuckForm2(task, user,cookie,apis)
		SubmitForm(cookie,user,form,apis)
	}
}

//获取任务详情
func GetDetailTask(cookie string,params,apis *map[string]string) TaskDeatil {
	api := GetSignTaskDetailApi(apis)

	sucess, bytes := PostRequest(api, cookie, RequestHeader, params)

	var task TaskDeatil

	if !sucess {
		return task
	}
	//fmt.Println("任务:",string(bytes))
	err := json.Unmarshal(bytes, &task)
	if err != nil {
		fmt.Println("TaskDeatil反序列化失败")
		return task
	}
	return task
}

//填写表单
func FuckForm(task TaskDeatil,user *User) map[string]interface{} {

	form := make(map[string]interface{})
	if task.Datas.IsNeedExtra == 1 {
		extraFields := task.Datas.ExtraField
		var extraFieldItemValues []map[string]interface{}
		for _,v := range extraFields {
			//检测问题是否对得上
			if questions[v.Title] == "" {
				fmt.Println("问题对不上:",v.Title)
				return nil
			}
			for _,v2 := range v.ExtraFieldItems {
				extraFieldItemValue := make(map[string]interface{})
				if v2.Content == questions[v.Title] {
					extraFieldItemValue["extraFieldItemValue"] = questions[v.Title]
					extraFieldItemValue["extraFieldItemWid"] = v2.Wid
					extraFieldItemValues = append(extraFieldItemValues,extraFieldItemValue)
				}

				if v2.IsOtherItems == 1 {
					fmt.Println("有额外任务")
					fmt.Println(task)
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
func FuckForm2(task TaskDeatil,user *User,cookie string,apis *map[string]string) map[string]interface{} {

	form := make(map[string]interface{})
	if task.Datas.IsNeedExtra == 1 {
		extraFields := task.Datas.ExtraField
		var extraFieldItemValues []map[string]interface{}
		for _,v := range extraFields {
			//检测问题是否对得上
			if questions[v.Title] == "" {
				fmt.Println("问题对不上:",v.Title)
				return nil
			}
			for _,v2 := range v.ExtraFieldItems {
				extraFieldItemValue := make(map[string]interface{})
				if v2.Content == questions[v.Title] {
					extraFieldItemValue["extraFieldItemValue"] = questions[v.Title]
					extraFieldItemValue["extraFieldItemWid"] = v2.Wid
					extraFieldItemValues = append(extraFieldItemValues,extraFieldItemValue)
				}

				if v2.IsOtherItems == 1 {
					fmt.Println("有额外任务")
					fmt.Println(task)
					continue
				}
			}
		}
		form["extraFieldItems"] = extraFieldItemValues
	}

	if task.Datas.IsPhoto == 1 {
		list := user.FileList
		//上传图片到今日校园的oss
		picMax := len(user.FileList)
		randInt := RandInt64(int64(picMax))
		fileName := UploadPicture(apis,list[randInt],cookie)
		pic := GetPic(fileName, cookie, apis)
		if pic == "" {
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
	return form
}

func SubmitForm(cookie string,user *User,form map[string]interface{},apis *map[string]string)  {
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
	jsonString := strings.ReplaceAll(string(marshal),":",": ")
	jsonString = strings.ReplaceAll(jsonString,",",", ")
	encrypt,_ :=  Encrypt([]byte(jsonString),KEY,IV)
	encoded := base64.StdEncoding.EncodeToString(encrypt)

	header["Cpdaily-Extension"] = encoded
	header["Content-Type"] = "application/json; charset=utf-8"
	header["Accept-Encoding"] = "gzip"
	header["Connection"] = "Keep-Alive"

	sucess, bytes := PostRequest(GetSubmitSignApi(apis), cookie, header, form)
	if !sucess {
		fmt.Println("提交任务失败")
		return
	}
	parseBytes, err := fastjson.ParseBytes(bytes)
	if err != nil {
		fmt.Println("返回json序列化失败")
		fmt.Println(parseBytes.String())
		return
	}
	message := string(parseBytes.GetStringBytes("message"))
	fmt.Println("签到信息:",message)
	if message == "SUCCESS" {
		fmt.Println("签到成功")
	}
}

func Sign(u *User,isFailProcess bool,loginApi string) {
	fmt.Println(u.UserName,"开始签到")
	apis := GetCpdailyApis(schoolName)
	cookie := GetCookie(u, apis,loginApi)
	if cookie == "1" {
		//密码错误的
		fmt.Println(u.UserName,"账号密码错误")
		return
	}

	if cookie == "" {
		if isFailProcess {
			fmt.Println(u.UserName,"登录失败,加入续命队列")
			failSlice = append(failSlice,u)
		}else {
			fmt.Println(u.UserName,"登录教务失败")
		}
		return
	}

	//GetSignTaskQA(cookie,&apis,user)
	GetScoolSignTasksAndSign(cookie,&apis,u)
	fmt.Println(u.UserName,"签到完毕")
}

func GetSignTaskQA(cookie string,apis *map[string]string,user *User) map[string]map[string][]string {
	//得到签到路径的api
	api := GetSignInfoApi(apis)

	PostRequest(api, cookie, RequestHeader, nil)
	//这个是为了构造空的json body
	var n name
	sucess, bytes := PostRequest(api, cookie, RequestHeader, n)
	if !sucess {
		fmt.Println("Post网络请求失败")
		return nil
	}
	fmt.Println(string(bytes))
	value, err := fastjson.ParseBytes(bytes)
	if err != nil {
		fmt.Println("签到任务反序列化失败")
		return nil
	}
	datas := value.Get("datas")
	fmt.Println("全部任务:",datas.String())
	unSignedTasks := datas.GetArray("unSignedTasks")
	signedTasks := datas.GetArray("signedTasks")

	//问卷 -- 问题 -- []答案
	tasks := make(map[string]map[string][]string)

	for _,v := range unSignedTasks {
		if tasks[string(v.GetStringBytes("taskName"))] != nil {
			continue
		}

		params := make(map[string]string)
		params["signInstanceWid"] = string(v.GetStringBytes("signInstanceWid"))
		params["signWid"] = string(v.GetStringBytes("signWid"))
		task := GetDetailTask(cookie, &params, apis)

		taskTitleAndAnswer := make(map[string][]string)
		for _,v := range task.Datas.ExtraField {
			var answer []string
			for _,v := range v.ExtraFieldItems {
				answer = append(answer,v.Content)
			}
			taskTitleAndAnswer[v.Title] = answer
		}
		tasks[task.Datas.TaskName] = taskTitleAndAnswer
	}

	for _,v := range signedTasks {
		if tasks[string(v.GetStringBytes("taskName"))] != nil {
			continue
		}

		params := make(map[string]string)
		params["signInstanceWid"] = string(v.GetStringBytes("signInstanceWid"))
		params["signWid"] = string(v.GetStringBytes("signWid"))
		task := GetDetailTask(cookie, &params, apis)

		taskTitleAndAnswer := make(map[string][]string)
		for _,v := range task.Datas.ExtraField {
			var answer []string
			for _,v := range v.ExtraFieldItems {
				answer = append(answer,v.Content)
			}
			taskTitleAndAnswer[v.Title] = answer
		}
		tasks[task.Datas.TaskName] = taskTitleAndAnswer
	}
	//fmt.Println(tasks)
	return tasks
	//marshal, _ := json.Marshal(tasks)
	//fmt.Println("json后的任务数据:",string(marshal))
}

func SignAllUser()  {
	dir := ReadDir("./user")
	var users []*User
	for _,v := range dir {
		user := ReadFile("./user/"+v.Name())
		users = append(users,user)
	}
	//每组用户数量
	step := len(users) / len(loginApiList)
	if step == 0 {
		step = 1
	}
	for k,v := range loginApiList {
		lgApi := v
		var signUserSlice []*User
		if (k+1)*step > len(users) {
			continue
		}
		if k != len(loginApiList)-1 {
			signUserSlice = users[k*step:(k+1)*step]
		}else {
			signUserSlice = users[k*step:]
		}
		go func() {
			for _,v2 := range signUserSlice {
				Sign(v2,false,lgApi)
				//fmt.Println(v2,lgApi)
				if SignStepTime != 0 {
					time.Sleep(time.Duration(SignStepTime) * time.Second)
				}else {
					time.Sleep(time.Second * 30)
				}
			}
		}()
	}
}

func SignFallUser(users []*User)  {
	step := len(users) / len(loginApiList)
	for k,v := range loginApiList {
		lgApi := v
		var signUserSlice []*User
		if k != len(loginApiList)-1 {
			signUserSlice = users[k*step:(k+1)*step]
		}else {
			signUserSlice = users[k*step:]
		}
		go func() {
			for _,v2 := range signUserSlice {
				Sign(v2,false,lgApi)
				time.Sleep(time.Second * 25)
			}
		}()
	}
}

//上传图片到今日校园的OSS
func UploadPicture(apis *map[string]string,imgName,cookie string) string {
	url := "https://"+(*apis)["host"]+"/wec-counselor-sign-apps/stu/oss/getUploadPolicy"
	params := make(map[string]int)
	params["fileType"] = 1

	sucess,res := PostRequest(url,cookie,nil,params)
	if !sucess || res == nil {
		return ""
	}
	fmt.Println(string(res))
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
	PostMultipartImage(&data,&header,imgName,policyHost)
	return fileName
}

//获取图片
func GetPic(fileName,cookie string,apis *map[string]string) string {
	url := "https://"+(*apis)["host"]+"/wec-counselor-sign-apps/stu/sign/previewAttachment"
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

func CreateCron()  {
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

	if MorningSignTime != "" {
		bc.AddFunc(MorningSignTime, SignAllUser)
	}

	if NoonSignTime != "" {
		bc.AddFunc(NoonSignTime, SignAllUser)
	}

	if EveningSignTime != "" {
		bc.AddFunc(EveningSignTime, SignAllUser)
	}

	bc.AddFunc(morningEndSepc, func() {
		failSlice = failSlice[0:0]
	})

	bc.AddFunc(noonEndSpec, func() {
		failSlice = failSlice[0:0]
	})

	bc.AddFunc(eveningEndSpec, func() {
		failSlice = failSlice[0:0]
	})
}