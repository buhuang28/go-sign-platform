// 代码由简易GoVCL IDE自动生成。
// 不要更改此文件名
// 在这里写你的事件。

package main

import (
	"encoding/json"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/win"
	"strconv"
	"strings"
)

//::private::
type TForm1Fields struct {
	subItemHit win.TLVHitTestInfo
}

func (f *TForm1) OnFormCreate(sender vcl.IObject) {
	for i := 0; i < 24; i++ {
		f.MorningHour.Items().Add(strconv.FormatInt(int64(i), 10))
		f.NoonHour.Items().Add(strconv.FormatInt(int64(i), 10))
		f.EveningHour.Items().Add(strconv.FormatInt(int64(i), 10))
	}
	for i := 0; i < 60; i++ {
		f.MorningMin.Items().Add(strconv.FormatInt(int64(i), 10))
		f.NoonMin.Items().Add(strconv.FormatInt(int64(i), 10))
		f.EveningMin.Items().Add(strconv.FormatInt(int64(i), 10))
	}
	for i := 20; i < 100; i++ {
		f.SignStep.Items().Add(strconv.FormatInt(int64(i), 10))
	}

	exits := CheckFileIsExits("data.json")
	if exits {
		settingData := ReadSetting()
		f.SchoolEdit.SetText(settingData.SchoolName)
		i := 0
		for k, v := range settingData.QuestionAndAnswer {
			item := f.TaskListView.Items().Add()
			item.SetCaption(strconv.FormatInt(int64(i), 10))
			i++
			item.SubItems().Add(k)
			item.SubItems().Add(v)
		}
		f.PortEdit.SetText(settingData.RunPort)
		f.CallBackEdit.SetText(settingData.CallBackApi)
		f.MorningHour.SetText(settingData.MorningHour)
		f.MorningMin.SetText(settingData.MorningMin)
		f.NoonHour.SetText(settingData.NoonHour)
		f.NoonMin.SetText(settingData.NoonMin)
		f.EveningHour.SetText(settingData.EveningHour)
		f.EveningMin.SetText(settingData.EveningMin)
		f.SignStep.SetText(settingData.SignStep)

		for _, v := range settingData.LoginApiList {
			if v != "" {
				item := f.LoginApiListView.Items().Add()
				item.SetCaption("")
				item.SubItems().Add(v)
			}
		}

	}
	f.ScreenCenter()
	f.SetCaption(f.TForm.Caption() + "带图片版")
}

func (f *TForm1) OnGetTaskButtonClick(sender vcl.IObject) {
	f.TaskListView.Clear()
	if strings.TrimSpace(f.SchoolEdit.Text()) == "" || len(strings.TrimSpace(f.SchoolEdit.Text())) < 4 {
		vcl.ShowMessage("学校不可为空")
		return
	}
	if f.UserEdit.Text() == "" || f.PassWordEdit.Text() == "" {
		vcl.ShowMessage("学号密码不可为空")
		return
	}
	apis := GetCpdailyApis(strings.TrimSpace(f.SchoolEdit.Text()))
	if apis == nil {
		vcl.ShowMessage("该学校未加入今日校园，或者学校名字错误")
		return
	}
	var user User
	user.UserName = strings.TrimSpace(f.UserEdit.Text())
	user.PassWord = strings.TrimSpace(f.PassWordEdit.Text())

	tempLoginApi := ""

	for i := 0; i < int(f.LoginApiListView.Items().Count()); i++ {
		api := f.LoginApiListView.Items().Item(int32(i)).SubItems().Strings(0)
		api = strings.TrimSpace(api)
		if api != "" {
			tempLoginApi = api
			break
		}
	}
	if tempLoginApi == "" {
		vcl.ShowMessage("没有获取登录cookie的API")
		return
	}
	cookie := GetCookie(&user, apis, tempLoginApi)
	header2 := make(map[string]string)
	header2["Content-Type"] = "application/json"
	header2["Cookie"] = cookie

	api := GetSignInfoApi(&apis)
	realCookie, _ := SchoolDayGetLocationCookie(api, header2)
	cookie = realCookie
	if cookie == "" {
		vcl.ShowMessage("获取不到任务问卷，可能是账号被限制，请切换账号或者自己排查问题")
		return
	}
	//任务总标题 ---- 签到问题 ---- 回答答案
	QAndA := GetSignTaskQA(cookie, &apis, &user)
	if QAndA == nil {
		vcl.ShowMessage("获取不到任务问卷，可能是账号被限制，请切换账号或者自己排查问题")
		return
	}
	i := 1
	Q := make(map[string]bool)
	A := make(map[string]bool)

	for _, v := range QAndA {
		for k, v2 := range v {
			Q[k] = true
			for _, v3 := range v2 {
				A[v3] = true
			}
		}
	}

	for k, _ := range Q {
		item := f.TaskListView.Items().Add()
		item.SetCaption(strconv.FormatInt(int64(i), 10))
		i++
		item.SubItems().Add(k)
		item.SubItems().Add("选择回答")
	}

	for k, _ := range A {
		f.ComboBox1.Items().Add(k)
	}
	f.UserEdit.SetText("")
	f.PassWordEdit.SetText("")
}

func (f *TForm1) OnApplyButtonClick(sender vcl.IObject) {
	if strings.TrimSpace(f.SchoolEdit.Text()) == "" || len(strings.TrimSpace(f.SchoolEdit.Text())) < 4 {
		vcl.ShowMessage("学校不可为空")
		return
	}
	schoolName = strings.TrimSpace(f.SchoolEdit.Text())
	newQA := make(map[string]string)
	for i := 0; i < int(Form1.TaskListView.Items().Count()); i++ {
		Q := Form1.TaskListView.Items().Item(int32(i)).SubItems().Strings(0)
		if Q == "" {
			vcl.ShowMessage("问题不可为空")
			return
		}
		A := Form1.TaskListView.Items().Item(int32(i)).SubItems().Strings(1)
		if A == "" {
			vcl.ShowMessage("签到答案不可为空")
			return
		}
		newQA[Q] = A
	}
	if len(newQA) == 0 {
		vcl.ShowMessage("签到问题不可为空")
		return
	}
	questions = newQA

	cbApi := strings.TrimSpace(f.CallBackEdit.Text())
	if cbApi != "" {
		runes := []rune(cbApi)
		protocol := ""
		if strings.Contains(cbApi, "https://") && len(runes) > 10 {
			protocol = string(runes[0:8])
		} else if strings.Contains(cbApi, "http://") && len(runes) > 10 {
			protocol = string(runes[0:7])
		} else {
			vcl.ShowMessage("非法回调地址")
			return
		}
		if protocol != "https://" && "http://" != protocol {
			vcl.ShowMessage("非法回调地址")
			return
		}
		callBackApi = cbApi
	}

	pt := strings.TrimSpace(f.PortEdit.Text())
	if pt == "" {
		vcl.ShowMessage("端口号不可为空")
		return
	}

	parseInt, err := strconv.ParseInt(pt, 10, 64)
	if err != nil {
		vcl.ShowMessage("端口号错误")
		return
	}
	port = ":" + strconv.FormatInt(parseInt, 10)

	if f.SignStep.Text() == "" {
		vcl.ShowMessage("签到间隔最低为30秒")
		return
	}

	var data SettingData

	if f.MorningMin.Text() != "" && f.MorningHour.Text() != "" {
		MorningSignTime = "0 " + f.MorningMin.Text() + " " + f.MorningHour.Text() + " * * ?"
		data.MorningHour = f.MorningHour.Text()
		data.MorningMin = f.MorningMin.Text()
	}
	if f.NoonMin.Text() != "" && f.NoonHour.Text() != "" {
		NoonSignTime = "0 " + f.NoonMin.Text() + " " + f.NoonHour.Text() + " * * ?"
		data.NoonHour = f.NoonHour.Text()
		data.NoonMin = f.NoonMin.Text()
	}

	if f.EveningMin.Text() != "" && f.EveningHour.Text() != "" {
		EveningSignTime = "0 " + f.EveningMin.Text() + " " + f.EveningHour.Text() + " * * ?"
		data.EveningHour = f.EveningHour.Text()
		data.EveningMin = f.EveningMin.Text()
	}

	data.SignStep = f.SignStep.Text()
	SignStepTime, _ = strconv.ParseInt(f.SignStep.Text(), 10, 64)

	var apiList []string

	for i := 0; i < int(f.LoginApiListView.Items().Count()); i++ {
		api := f.LoginApiListView.Items().Item(int32(i)).SubItems().Strings(0)
		api = strings.TrimSpace(api)
		if api != "" {
			apiList = append(apiList, api)
		} else {
			f.LoginApiListView.Items().Delete(int32(i))
			i--
			continue
		}
	}

	loginApiList = apiList

	data.LoginApiList = apiList
	data.SchoolName = strings.TrimSpace(f.SchoolEdit.Text())
	data.QuestionAndAnswer = newQA
	data.CallBackApi = cbApi
	data.RunPort = strings.TrimSpace(f.PortEdit.Text())
	marshal, _ := json.Marshal(data)
	WriteContent("data.json", string(marshal))
	go CreateCron()
	go WebStart(port)
	vcl.ShowMessage("数据应用成功")
}

func (f *TForm1) OnTaskListViewClick(sender vcl.IObject) {
	p := f.TaskListView.ScreenToClient(vcl.Mouse.CursorPos())
	f.subItemHit.Pt.X = p.X
	f.subItemHit.Pt.Y = p.Y
	win.ListView_SubItemHitTest(f.TaskListView.Handle(), &f.subItemHit)
	if f.subItemHit.IItem != -1 {
		var r types.TRect

		if f.TaskListView.RowSelect() {
			r = f.TaskListView.Selected().DisplayRect(types.DrBounds)
		} else {
			win.ListView_GetItemRect(f.TaskListView.Handle(), f.subItemHit.IItem, &r, 0)
		}

		colWidht := f.TaskListView.Column(f.subItemHit.ISubItem).Width()

		var left, i int32
		// 差2个像素???????????????????????????????????
		left += 232
		r.Top += 2
		for i = 0; i < f.subItemHit.ISubItem; i++ {
			left += f.TaskListView.Column(i).Width() //ListView_GetColumnWidth(f.ListView1.Handle(), i)
		}
		switch f.subItemHit.ISubItem {
		case 0:
		case 1:
		case 2:
			f.ComboBox1.SetText(f.TaskListView.Items().Item(f.subItemHit.IItem).SubItems().Strings(f.subItemHit.ISubItem - 1))
			f.ComboBox1.SetBounds(left, f.TaskListView.Top()+r.Top, colWidht, r.Bottom-r.Top)
			f.ComboBox1.Show()
			f.ComboBox1.SetFocus()
		}
	}
}

func (f *TForm1) OnComboBox1Exit(sender vcl.IObject) {
	f.ComboBox1.Hide()
	if f.subItemHit.IItem != -1 {
		if f.subItemHit.ISubItem == 2 {
			//f.ListView4.Items().Item(f.subItemHit.IItem).SubItems().SetStrings(f.subItemHit.ISubItem-1, f.ComboBox1.Text())
			f.TaskListView.Items().Item(f.subItemHit.IItem).SubItems().SetStrings(f.subItemHit.ISubItem-1, f.ComboBox1.Text())
		}
		//f.TaskListView.Items().Item(f.subItemHit.IItem).SubItems().SetStrings(f.subItemHit.ISubItem-1, f.ComboBox1.Text())
	}
}

func (f *TForm1) OnAddLoginApiButtonClick(sender vcl.IObject) {
	item := f.LoginApiListView.Items().Add()
	item.SetCaption("")
	item.SubItems().Add("http://127.0.0.1:8090/api")
}

func (f *TForm1) OnLoginApiListViewClick(sender vcl.IObject) {
	p := f.LoginApiListView.ScreenToClient(vcl.Mouse.CursorPos())
	f.subItemHit.Pt.X = p.X
	f.subItemHit.Pt.Y = p.Y
	win.ListView_SubItemHitTest(f.LoginApiListView.Handle(), &f.subItemHit)
	if f.subItemHit.IItem != -1 {

		var r types.TRect

		if f.LoginApiListView.RowSelect() {
			r = f.LoginApiListView.Selected().DisplayRect(types.DrBounds)
		} else {
			win.ListView_GetItemRect(f.LoginApiListView.Handle(), f.subItemHit.IItem, &r, 0)
		}

		//colWidht := ListView_GetColumnWidth(f.ListView2.Handle(), f.subItemHit.iSubItem)
		colWidht := f.LoginApiListView.Column(f.subItemHit.ISubItem).Width()

		//var itemPoint types.TPoint
		//ListView_GetItemPosition(f.ListView2.Handle(), f.subItemHit.iItem, &itemPoint)

		var left, i int32
		// 差2个像素???????????????????????????????????
		left += 326
		//r.Top +=
		for i = 0; i < f.subItemHit.ISubItem; i++ {
			left += f.LoginApiListView.Column(i).Width() //ListView_GetColumnWidth(f.ListView2.Handle(), i)
		}

		switch f.subItemHit.ISubItem {
		case 0: // 不处理
		case 1:
			// edit
			f.LoginApiEdit.SetText(f.LoginApiListView.Items().Item(f.subItemHit.IItem).SubItems().Strings(f.subItemHit.ISubItem - 1))
			f.LoginApiEdit.SetBounds(left, f.LoginApiListView.Top()+r.Top, colWidht, r.Bottom-r.Top)
			f.LoginApiEdit.Show()
			f.LoginApiEdit.SetFocus()
		}
	}
}

func (f *TForm1) OnLoginApiEditExit(sender vcl.IObject) {

	f.LoginApiEdit.Hide()
	if f.subItemHit.IItem != -1 {
		f.LoginApiListView.Items().Item(f.subItemHit.IItem).SubItems().SetStrings(f.subItemHit.ISubItem-1, f.LoginApiEdit.Text())
	}
}
