package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestAllSign(t *testing.T) {
	dir := ReadDir("./user")
	for _, v := range dir {
		user := ReadFile("./user/" + v.Name())
		Sign(user, false, LoginApi)
		time.Sleep(time.Second * 30)
	}
}

func TestWriteContent(t *testing.T) {
	WriteContent("test.txt", "12345677")
}

func TestMap2(t *testing.T) {
	m := make(map[string]bool)
	logger.Println(m["123"])
}

func TestLock(t *testing.T) {
	m := make(map[string]string)
	for i := 0; i < 10000; i++ {
		go logger.Println(i)
		//i0 := i
		//go func() {
		//Lock.Lock()
		//logger.Println(i0)
		//Lock.Unlock()
		//Lock.Lock()
		//is := strconv.FormatInt(int64(i),10)
		//m[is] = is
		//Lock.Unlock()
		//}()
	}
	//for  j := 0;j<10000;j++ {
	//	go func() {
	//		Lock.Lock()
	//		is := strconv.FormatInt(int64(j),10)
	//		_ = m[is]
	//		Lock.Unlock()
	//		//logger.Println(value)
	//	}()
	//}
	time.Sleep(time.Second * 2)
	logger.Println(m)
}

func lc2(m *map[string]string, key string) {
	Lock.Lock()
	v := (*m)[key]
	Lock.Unlock()
	logger.Println(v)
}

func TestLogin(t *testing.T) {
	var u User
	u.UserName = "213"
	u.PassWord = "213213.."
	apis := GetCpdailyApis(schoolName)
	cookie := GetCookie(&u, apis, LoginApi)
	logger.Println(cookie)
}

func TestSignAllUser(t *testing.T) {
	dir := ReadDir("./user")
	for _, v := range dir {
		user := ReadFile("./user/" + v.Name())
		Sign(user, true, LoginApi)
	}

}

func TestShangchuan(t *testing.T) {
	//准备一个您将提交的表单该网址。

	var url, filename string
	filename = "1.jpg"
	url = "http://127.0.0.1:8081/test"
	var bytedata bytes.Buffer
	//转换成对应的格式
	Multipar := multipart.NewWriter(&bytedata)
	//添加您的镜像文件
	fileData, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fileData.Close()

	//这里添加图片数据
	form, err := Multipar.CreateFormFile2("file", "blob","image/jpg")
	if err != nil {
		return
	}
	if _, err = io.Copy(form, fileData); err != nil {
		return
	}

	//添加其他字段
	if form, err = Multipar.CreateFormField("key"); err != nil {
		return
	}
	if _, err = form.Write([]byte("KEY")); err != nil {
		return
	}
	//不要忘记关闭multipart writer。
	//如果你不关闭它,你的请求将丢失终止边界。
	Multipar.Close()

	//现在你有一个表单,你可以提交它给你的处理程序。
	req, err := http.NewRequest("POST", url, &bytedata)
	if err != nil {
		return
	}
	// Don不要忘记设置内容类型,这将包含边界。
	req.Header.Set("Content-Type", Multipar.FormDataContentType())

	//提交请求
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	//检查响应
	//if res.StatusCode != http.StatusOK {
	//	err = fmt.Errorf("bad status:％s", res.Status)
	//}
	logger.Println(res.Body)
	return
}




