package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetRequest(u string, headerData map[string]string, urlParam map[string]string) ([]byte, bool) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Println(err)
		}
	}()
	if u == "" {
		return nil, false
	}
	request, _ := http.NewRequest("GET", u, nil)
	if headerData != nil {
		for k, v := range headerData {
			request.Header.Set(k, v)
		}
	}

	//加入get参数
	q := request.URL.Query()
	if urlParam != nil {
		for k, v := range urlParam {
			q.Add(k, v)
		}
	}
	request.URL.RawQuery = q.Encode()

	timeout := time.Duration(6 * time.Second)
	//urli := url.URL{}
	//urlproxy, _ := urli.Parse("http://127.0.0.1:8080")
	//fmt.Println(urlproxy)
	client := http.Client{
		//Transport: &http.Transport{
		//	Proxy: http.ProxyURL(urlproxy),
		//},
		Timeout: timeout,
	}
	resp, err := client.Do(request)
	if err != nil {
		return []byte{}, false
	}

	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return nil, false
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return data, true
}

func GetRqeustCookie(u string, headerData map[string]string, urlParam map[string]string) (data []byte, cookie string) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Println(err)
		}
	}()
	if u == "" {
		return nil, ""
	}
	request, _ := http.NewRequest("GET", u, nil)
	if headerData != nil {
		for k, v := range headerData {
			request.Header.Set(k, v)
		}
	}

	//加入get参数
	q := request.URL.Query()
	if urlParam != nil {
		for k, v := range urlParam {
			q.Add(k, v)
		}
	}
	request.URL.RawQuery = q.Encode()

	timeout := time.Duration(6 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, ""
	}

	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return nil, ""
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	ck := ""
	cookies := resp.Cookies()
	if len(cookies) == 0 {
		ck = headerData["Cookie"]
	} else {
		for _, v := range cookies {
			ck = ck + ";" + v.Name + "=" + v.Value
		}
	}

	ck = strings.Trim(ck, ";")
	return data, ck
}

func GetScheme(data string) string {
	split := strings.Split(data, "://")
	return split[0]
}

var (
	hostReg = regexp.MustCompile(`\w{4,5}\:\/\/.*?\/`)
)

func GetNetLocol(data string) string {
	split := strings.Split(data, "//")
	i := strings.Split(split[1], "/")
	return i[0]
}

func GetRegData(data string) string {
	allString := hostReg.FindAllString(data, -1)
	return allString[0]
}

func PostRequest(url, cookie string, header map[string]string, data interface{}) (bool, []byte) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Println(err)
		}
	}()

	bytesData := []byte(`{}`)

	if data != nil {
		bytesData, _ = json.Marshal(data)
	}

	timeout := time.Duration(6 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("POST", url, bytes.NewReader(bytesData))

	if err != nil {
		fmt.Println(err)
		return false, []byte{}
	}

	//if cookie != "" && cookie != "1" {
	//	request.Header.Add("cookie", cookie)
	//	cookie = strings.Trim(cookie,";")
	//	split := strings.Split(cookie, ";")
	//	for _, v := range split {
	//		v = strings.TrimSpace(v)
	//		if v == "" {
	//			continue
	//		}
	//		i := strings.Split(v, "=")
	//		if len(i) > 1 {
	//			ck := &http.Cookie{Name: i[0], Value: i[1], HttpOnly: true}
	//			request.AddCookie(ck)
	//		}
	//	}
	//}

	//request.Header.Add("Content-Type", "application/json")
	//request.Header.Add("Cookie", "iPlanetDirectoryPro=9j6D5f5dk9PiYbaHGSgcMi;JSESSIONID=6EuiY1QW3vQDXmVDbLTOF27wtX1AmzQBbEVfxZIdr2TRQqYHASbm!-1406639286;route=922000a575f0ac992fd468a3d924eed4;CASTGC=TGT-68306-PJMItgoxbPfrgTHbXh3Va2fPIjkbPOpHZfkm5iRtte1RSIpkOG1621927089431-w0vA-cas;HWWAFSESID=6c56610041ec03906f2;HWWAFSESTIME=1621927083406;HWWAFSESID=70af373fec19952ef9;HWWAFSESTIME=1621927088620; MOD_AUTH_CAS=ST-iap:1018615964750600:ST:05f782ca-e8a9-42c0-83bf-a336d272cc37:20210525153409")

	if cookie != "" {
		request.Header.Add("Cookie", cookie)
	}

	if len(header) != 0 && cookie != "1" {
		for k, v := range header {
			request.Header.Add(k, v)
		}
	} else {
		request.Header.Add("Content-Type", "application/json")
	}

	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		logger.Println("请求失败")
		return false, []byte{}
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println("数据读取失败")
		return false, []byte{}
	}
	return true, respBytes
}

func SendPostForm(api string, params map[string]string) (bool, []byte) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Println(err)
		}
	}()
	data := make(url.Values)
	if params != nil {
		for k, v := range params {
			data[k] = []string{v}
		}
	}

	timeout := time.Duration(6 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	res, err := client.PostForm(api, data)
	if err != nil {
		logger.Println(err.Error())
		return false, []byte{}
	}

	respBytes, err := ioutil.ReadAll(res.Body)
	if err != nil || respBytes == nil {
		logger.Println(err)
		return false, []byte{}
	}

	defer res.Body.Close()
	return true, respBytes
}

func Encrypt(origData, key, iv []byte) ([]byte, error) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Println(err)
		}
	}()
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
	defer func() {
		err := recover()
		if err != nil {
			logger.Println(err)
		}
	}()
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func GetSignInfoApi(apis *map[string]string) string {
	return (*apis)["host"] + "wec-counselor-sign-apps/stu/sign/getStuSignInfosInOneDay"
}

func GetSignTaskDetailApi(apis *map[string]string) string {
	return (*apis)["host"] + "wec-counselor-sign-apps/stu/sign/detailSignInstance"
}

func GetSubmitSignApi(apis *map[string]string) string {
	return (*apis)["host"] + "wec-counselor-sign-apps/stu/sign/submitSign"
}

func MD5Sign(user *User, full bool) Result {
	var result Result
	timeStamp, _ := strconv.ParseInt(user.Time, 10, 64)
	nowTime := time.Now().Unix()
	if nowTime-timeStamp > 100 {
		result.Code = -1
		result.Message = "invaid time"
		return result
	}

	sign := ""
	fieldNames := GetFieldName(*user)
	tgName := GetTagName(*user)
	sort.Strings(fieldNames)
	sort.Strings(tgName)
	checkSign := ""
	immutable := reflect.ValueOf(*user)
	for k, v := range fieldNames {
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
		} else {
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
				result.Message = "invaid " + v
				return result
			} else {
				continue
			}
		}
		sign += tgName[k] + value
	}
	sign = sessionKey + sign + sessionKey
	sign = Md5(sign)
	//sign 是自己算的，checkSign是获取结构体的
	if checkSign != sign {
		result.Code = -2
		result.Message = "invaid sign"
		return result
	}

	if sign != "" {
		result.Code = 200
	} else {
		result.Code = -3
		result.Message = "invaid sign"
		return result
	}
	return result
}

//func MD5Sign(user *User,full bool) Result {
//	var result Result
//	timeStamp, _ := strconv.ParseInt(user.Time, 10, 64)
//	nowTime := time.Now().Unix()
//	if nowTime - timeStamp > 100 {
//		result.Code = -1
//		result.Message = "invaid sign"
//		return result
//	}
//
//	sign := ""
//	fieldNames := GetFieldName(*user)
//	tgName := GetTagName(*user)
//	sort.Strings(fieldNames)
//	sort.Strings(tgName)
//	checkSign := ""
//	immutable := reflect.ValueOf(*user)
//	for k,v := range fieldNames {
//		if v == "Sign" {
//			checkSign = immutable.FieldByName(v).String()
//			continue
//		}
//		if v == "AbnormalReason" {
//			continue
//		}
//		value := ""
//		if v != "FileList" {
//			value = immutable.FieldByName(v).String()
//		}else {
//			list := user.FileList
//			for _,v2 := range list {
//				value += v2+","
//			}
//			value = strings.Trim(value,",")
//		}
//		if value == "" {
//			if full {
//				result.Code = -1
//				result.Message = "invaid "+v
//				return result
//			}else {
//				continue
//			}
//		}
//		sign +=tgName[k]+value
//	}
//	sign = sessionKey+sign+sessionKey
//	sign = Md5(sign)
//	//sign 是自己算的，checkSign是获取结构体的
//	if checkSign != sign {
//		result.Code = -2
//		result.Message = "invaid sign"
//		return result
//	}
//
//	if sign != "" {
//		result.Code = 200
//	}else {
//		result.Code = -3
//		result.Message = "invaid sign"
//		return result
//	}
//	return result
//}

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
	if !strings.Contains(longitude, ".") {
		result.Message = "错误的经度"
		result.Code = -6
		return result
	}
	split := strings.Split(longitude, ".")
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
	if !strings.Contains(latitude, ".") {
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

	if !strings.Contains(user.MorningTime, ":") || !strings.Contains(user.NoonTime, ":") || !strings.Contains(user.EveningTime, ":") {
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

	minute, err := strconv.ParseInt(split[1], 10, 64)
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

	minute, err = strconv.ParseInt(split[1], 10, 64)
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

	minute, err = strconv.ParseInt(split[1], 10, 64)
	if err != nil || minute > 59 || minute < 0 {
		result.Code = -9
		result.Message = "非正常时间数据"
		return result
	}

	result.Code = 200
	result.Message = "Check Sucess"
	return result
}

func WriteContent(fileName, content string) bool {
	fd, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	buf := []byte(content)
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
	json.Unmarshal(f, &u)
	return &u
}

func ReadFile2(fileName string) []byte {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Println("read fail", err)
		return nil
	}
	return f
}

func ReadSetting() *SettingData {
	f, err := ioutil.ReadFile("data.json")
	if err != nil {
		logger.Println("read fail", err)
		return nil
	}
	var data SettingData
	json.Unmarshal(f, &data)
	return &data
}

func ReadDir(path string) []os.FileInfo {
	FileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Println("读取 img 文件夹出错")
		return nil
	}
	if FileInfo == nil || len(FileInfo) == 0 {
		return nil
	}
	return FileInfo
}

func createLog() {
	logFileNmae := `./log/` + time.Now().Format("20060102") + ".log"
	logFileAllPath := logFileNmae
	_, err := os.Stat(logFileAllPath)
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
	if err != nil {
		f, _ = os.Create(logFileAllPath)
	} else {
		//如果存在文件则 追加log
		f, _ = os.OpenFile(logFileAllPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
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
func PostMultipartImage(data, header *map[string]string, imgName, url string) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Println(err)
		}
	}()
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
	form, err = CreateFormFile2(Multipar, "file", "blob", "image/jpg")
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
		for k, v := range *header {
			req.Header.Set(k, v)
		}
	} else {
		req.Header.Set("Content-Type", "application/json")
	}

	// Don不要忘记设置内容类型,这将包含边界。
	req.Header.Set("Content-Type", Multipar.FormDataContentType())
	logger.Println(Multipar.FormDataContentType())
	//提交请求
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logger.Println("提交错误", err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(res.Body)
	logger.Println("提交返回数据", string(respBytes))
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

func CreateFormFile2(w *multipart.Writer, fieldname, filename, contentType string) (io.Writer, error) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Println(err)
		}
	}()
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

//重定向
func GetLocation(url string, headerData map[string]string) (bool, string) {
	req, _ := http.NewRequest("GET", url, nil)

	if headerData != nil {
		for k, v := range headerData {
			req.Header.Set(k, v)
		}
	}
	timeout := time.Duration(6 * time.Second)

	client := http.Client{
		Timeout: timeout,
	}
	//client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
	//	return errors.New("Redirect")
	//}
	response, _ := client.Do(req)
	defer func() {
		response.Body.Close()
	}()
	if response == nil {
		return false, ""
	}
	location, _ := response.Location()
	if location == nil || location.String() == "" {
		header := response.Header
		location := header.Get("location")
		if location != "" {
			return true, location
		}
		return false, ""
	}
	return true, location.String()
}

func GetLocation2(url string, headerData map[string]string) (string, string) {
	req, _ := http.NewRequest("GET", url, nil)
	if headerData != nil {
		for k, v := range headerData {
			req.Header.Set(k, v)
		}
	}
	timeout := time.Duration(6 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	response, _ := client.Do(req)
	defer func() {
		response.Body.Close()
	}()
	if response == nil {
		return "", ""
	}
	location := response.Request.URL.String()
	ck := ""
	cookies := response.Cookies()
	for _, v := range cookies {
		ck = ck + ";" + v.Name + "=" + v.Value
	}
	ck = strings.Trim(ck, ";")
	return location, ck
}
