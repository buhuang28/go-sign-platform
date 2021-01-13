package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetRequest(url string,headerData map[string]string,urlParam map[string]string) ([]byte,bool) {
	if url == "" {
		return nil,false
	}
	request,_ := http.NewRequest("GET", url, nil)
	if headerData != nil{
		for k,v := range headerData {
			request.Header.Set(k, v)
		}
	}
	//加入get参数
	q := request.URL.Query()
	if urlParam != nil {
		for k,v := range urlParam{
			q.Add(k,v)
		}
	}
	request.URL.RawQuery = q.Encode()

	timeout := time.Duration(6 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(request)
	if err != nil {
		return []byte{},false
	}

	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return nil,false
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return data,true
}

func GetScheme(data string) string {
	split := strings.Split(data, "://")
	return split[0]
}

func GetNetLocol(data string) string {
	split := strings.Split(data, "//")
	i := strings.Split(split[1], "/")
	return i[0]
}

func PostRequest(url,cookie string,header map[string]string,data interface{}) (bool,[]byte)  {
	bytesData,_ := json.Marshal(data)
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return false,[]byte{}
	}

	if header != nil {
		for k,v := range header {
			request.Header.Set(k,v)
		}
	}else {
		request.Header.Set("Content-Type", "application/json")
	}

	if cookie != "" {
		split := strings.Split(cookie, ";")
		for _,v := range split {
			i := strings.Split(v, "=")
			ck := &http.Cookie{Name: i[0],Value: i[1],HttpOnly: true}
			request.AddCookie(ck)
		}
	}

	timeout := time.Duration(6 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(request)

	if err != nil {
		logger.Println("请求失败")
		return false,[]byte{}
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println("数据读取失败")
		return false,[]byte{}
	}
	return true,respBytes
}

func SendPostForm(api string,params map[string]string) (bool,[]byte) {
	data := make(url.Values)
	if params != nil {
		for k,v := range params {
			data[k] = []string{v}
		}
	}

	res, err := http.PostForm(api, data)
	if err != nil {
		logger.Println(err.Error())
		return false,[]byte{}
	}
	respBytes,err := ioutil.ReadAll(res.Body)
	if err != nil || respBytes == nil {
		logger.Println(err)
		return false,[]byte{}
	}

	defer res.Body.Close()
	return true,respBytes
}

func Encrypt(origData, key,iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func GetSignInfoApi(apis *map[string]string) string {
	return "https://"+(*apis)["host"]+"/wec-counselor-sign-apps/stu/sign/getStuSignInfosInOneDay"
}

func GetSignTaskDetailApi(apis *map[string]string) string {
	return "https://"+(*apis)["host"]+"/wec-counselor-sign-apps/stu/sign/detailSignInstance"
}

func GetSubmitSignApi(apis *map[string]string) string {
	return "https://"+(*apis)["host"]+"/wec-counselor-sign-apps/stu/sign/submitSign"
}

func MD5Sign(user *User,full bool) Result {
	var result Result
	timeStamp, _ := strconv.ParseInt(user.Time, 10, 64)
	nowTime := time.Now().Unix()
	if nowTime - timeStamp > 100 {
		result.Code = -1
		result.Message = "invaid sign"
		return result
	}

	sign := ""
	fieldNames := GetFieldName(*user)
	tgName := GetTagName(*user)
	sort.Strings(fieldNames)
	sort.Strings(tgName)
	checkSign := ""
	immutable := reflect.ValueOf(*user)
	for k,v := range fieldNames {
		if v == "Sign" {
			checkSign = immutable.FieldByName(v).String()
			continue
		}
		if v == "AbnormalReason" {
			continue
		}
		value := ""
		if v != "FileList" {
			value = immutable.FieldByName(v).String()
		}else {
			continue
			//list := user.FileList
			//for _,v2 := range list {
			//	value += v2+","
			//}
			//value = strings.Trim(value,",")
		}
		if value == "" {
			if full {
				result.Code = -1
				result.Message = "invaid "+v
				return result
			}else {
				continue
			}
		}
		sign +=tgName[k]+value
	}
	sign = sessionKey+sign+sessionKey
	sign = Md5(sign)
	//sign 是自己算的，checkSign是获取结构体的
	if checkSign != sign {
		result.Code = -2
		result.Message = "invaid sign"
		return result
	}

	if sign != "" {
		result.Code = 200
	}else {
		result.Code = -3
		result.Message = "invaid sign"
		return result
	}
	return result
}

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

func GetTagName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Tag.Get("json"))
	}
	return result
}

func CheckUserData(user *User) Result {
	var result Result
	if len(user.UserName) < 5 {
		result.Code = -4
		result.Message = "用户名出错"
		return result
	}

	if len(user.UserName) > 15 {
		result.Code = -4
		result.Message = "用户名出错了"
		return result
	}

	if len(user.PassWord) > 20 {
		result.Code = -5
		result.Message = "密码出错了"
		return result
	}

	longitude := user.Longitude
	if !strings.Contains(longitude,".") {
		result.Message = "错误的经度"
		result.Code = -6
		return result
	}
	split := strings.Split(longitude,".")
	prefix, _ := strconv.ParseInt(split[0], 10, 64)
	if prefix > 180 {
		result.Message = "错误经度"
		result.Code = -6
		return result
	}

	if len(split[1]) > 6 {
		result.Message = "错误的经度..."
		result.Code = -6
		return result
	}

	latitude := user.Latitude
	if !strings.Contains(latitude,"."){
		result.Message = "错误的纬度"
		result.Code = -7
		return result
	}
	split = strings.Split(latitude, ".")
	prefix, _ = strconv.ParseInt(split[0], 10, 64)
	if prefix > 90 {
		result.Message = "错误纬度"
		result.Code = -7
		return result
	}

	if len(split[1]) > 6 {
		result.Message = "错误的纬度..."
		result.Code = -7
		return result
	}

	if len(user.Address) < 5 {
		result.Code = -8
		result.Message = "地址错误"
		return result
	}

	if !strings.Contains(user.MorningTime,":") || !strings.Contains(user.NoonTime,":") || !strings.Contains(user.EveningTime,":")  {
		result.Code = -9
		result.Message = "无效时间"
		return result
	}

	split = strings.Split(user.MorningTime, ":")
	hour, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil || hour > 23 || hour < 0 {
		result.Code = -9
		result.Message = "非正常时间"
		return result
	}

	minute,err := strconv.ParseInt(split[1],10,64)
	if err != nil || minute > 59 || minute < 0 {
		result.Code = -9
		result.Message = "非正常时间数据"
		return result
	}

	split = strings.Split(user.NoonTime, ":")
	hour, err = strconv.ParseInt(split[0], 10, 64)
	if err != nil || hour > 23 || hour < 0 {
		result.Code = -9
		result.Message = "非正常时间"
		return result
	}

	minute,err = strconv.ParseInt(split[1],10,64)
	if err != nil || minute > 59 || minute < 0 {
		result.Code = -9
		result.Message = "非正常时间数据"
		return result
	}

	split = strings.Split(user.EveningTime, ":")
	hour, err = strconv.ParseInt(split[0], 10, 64)
	if err != nil || hour > 23 || hour < 0 {
		result.Code = -9
		result.Message = "非正常时间"
		return result
	}

	minute,err = strconv.ParseInt(split[1],10,64)
	if err != nil || minute > 59 || minute < 0 {
		result.Code = -9
		result.Message = "非正常时间数据"
		return result
	}

	result.Code = 200
	result.Message = "Check Sucess"
	return result
}

func WriteContent(fileName,content string) bool {
	fd,_:=os.OpenFile(fileName,os.O_WRONLY|os.O_TRUNC|os.O_CREATE,0644)
	buf:=[]byte(content)
	_, err := fd.Write(buf)
	fd.Close()
	if err == nil {
		return true
	}
	return false
}

func ReadFile(fileName string) *User {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Println("read fail", err)
		return nil
	}
	var u User
	json.Unmarshal(f,&u)
	return &u
}

func ReadSetting() *SettingData {
	f, err := ioutil.ReadFile("data.json")
	if err != nil {
		logger.Println("read fail", err)
		return nil
	}
	var data SettingData
	json.Unmarshal(f,&data)
	return &data
}

func ReadDir(path string) []os.FileInfo {
	FileInfo,err := ioutil.ReadDir(path )
	if err != nil{
		logger.Println("读取 img 文件夹出错")
		return nil
	}
	if FileInfo == nil || len(FileInfo) == 0 {
		return nil
	}
	return FileInfo
}

func createLog() {
	logFileNmae := `./log/`+time.Now().Format("20060102")+".log"
	logFileAllPath := logFileNmae
	_,err :=os.Stat(logFileAllPath)
	exits := CheckFileIsExits(`log`)
	if !exits {
		_ = os.Mkdir("./log", os.ModePerm)
	}

	exits = CheckFileIsExits(`user`)
	if !exits {
		_ = os.Mkdir("./user", os.ModePerm)
	}

	exits = CheckFileIsExits(`img`)
	if !exits {
		_ = os.Mkdir("./img", os.ModePerm)
	}

	var f *os.File
	if  err != nil{
		f, _= os.Create(logFileAllPath)
	}else{
		//如果存在文件则 追加log
		f ,_= os.OpenFile(logFileAllPath,os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	}
	logger = log.New(f, "", log.LstdFlags)
}

func CheckFileIsExits(fileName string) bool {
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//上传图片
func PostMultipartImage(data,header *map[string]string,imgName,url string)  {
	var bytedata bytes.Buffer
	//转换成对应的格式
	Multipar := multipart.NewWriter(&bytedata)
	//添加您的镜像文件
	fileData, err := os.Open(imgName)
	if err != nil {
		return
	}
	defer fileData.Close()

	form, _ := Multipar.CreateFormField("key")
	form.Write([]byte((*data)["key"]))

	form, _ = Multipar.CreateFormField("policy")
	form.Write([]byte((*data)["policy"]))

	form, _ = Multipar.CreateFormField("OSSAccessKeyId")
	form.Write([]byte((*data)["OSSAccessKeyId"]))

	form, _ = Multipar.CreateFormField("success_action_status")
	form.Write([]byte((*data)["success_action_status"]))

	form, _ = Multipar.CreateFormField("signature")
	form.Write([]byte((*data)["signature"]))

	//这里添加图片数据
	form, err = Multipar.CreateFormFile2("file", "blob","image/jpg")
	if err != nil {
		return
	}
	if _, err = io.Copy(form, fileData); err != nil {
		return
	}
	Multipar.Close()

	//现在你有一个表单,你可以提交它给你的处理程序。
	req, err := http.NewRequest("POST", url, &bytedata)
	if err != nil {
		return
	}
	if header != nil {
		for k,v := range *header {
			req.Header.Set(k,v)
		}
	}else {
		req.Header.Set("Content-Type", "application/json")
	}

	// Don不要忘记设置内容类型,这将包含边界。
	req.Header.Set("Content-Type", Multipar.FormDataContentType())
	logger.Println(Multipar.FormDataContentType())
	//提交请求
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logger.Println("提交错误",err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(res.Body)
	logger.Println("提交返回数据",string(respBytes))
	return
}

func RandInt64(max int64) int64 {
	rand.Seed(int64(time.Now().Nanosecond()))
	if max == 0 {
		return 0
	}
	b := rand.Int63n(max)
	return b
}

