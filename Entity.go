package main

import "reflect"

type SchoolInfo struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	Data    []struct {
		ID                         string `json:"id"`
		Name                       string `json:"name"`
		TenantCode                 string `json:"tenantCode"`
		Img                        string `json:"img"`
		Distance                   string `json:"distance"`
		ShortName                  string `json:"shortName"`
		JoinType                   string `json:"joinType"`
		CasLoginURL                string `json:"casLoginUrl"`
		IsEnter                    int    `json:"isEnter"`
		IdsURL                     string `json:"idsUrl"`
		AmpURL                     string `json:"ampUrl"`
		AmpURL2                    string `json:"ampUrl2"`
		PriorityURL                string `json:"priorityUrl"`
		AppID                      string `json:"appId"`
		AppSecret                  string `json:"appSecret"`
		MsgURL                     string `json:"msgUrl"`
		MsgAccessToken             string `json:"msgAccessToken"`
		MsgAppID                   string `json:"msgAppId"`
		MsgAppIDIos                string `json:"msgAppIdIos"`
		ZgAppKey                   string `json:"zgAppKey"`
		YktBalanceURL              string `json:"yktBalanceUrl"`
		YktTransferURL             string `json:"yktTransferUrl"`
		YktQrCodeURL               string `json:"yktQrCodeUrl"`
		XykURL                     string `json:"xykUrl"`
		UserShowCollege            string `json:"userShowCollege"`
		ScheduleOpenURL            string `json:"scheduleOpenUrl"`
		ScheduleDataURL            string `json:"scheduleDataUrl"`
		IsIdsProxy                 string `json:"isIdsProxy"`
		TenantNameImg              string `json:"tenantNameImg"`
		IsNeedAlias                string `json:"isNeedAlias"`
		ModifyPassURL              string `json:"modifyPassUrl"`
		ModifyPassSuccessURL       string `json:"modifyPassSuccessUrl"`
		ModifyPassDescr            string `json:"modifyPassDescr"`
		TaskURL                    string `json:"taskUrl"`
		TaskAppID                  string `json:"taskAppId"`
		CircleShowType             string `json:"circleShowType"`
		IsShowHotList              string `json:"isShowHotList"`
		AppStyleVersionID          string `json:"appStyleVersionId"`
		AppStyleResURL             string `json:"appStyleResUrl"`
		LikeBtnSpace               string `json:"likeBtnSpace"`
		IRobotURL                  string `json:"iRobotUrl"`
		ServicePagePlace           string `json:"servicePagePlace"`
		ScheduleAllDataURL         string `json:"scheduleAllDataUrl"`
		ScheduleUpdateDataURL      string `json:"scheduleUpdateDataUrl"`
		ShopURL                    string `json:"shopUrl"`
		HomePageDisplayItem        string `json:"homePageDisplayItem"`
		TaoBannerID                string `json:"taoBannerId"`
		CanIdsLogin                string `json:"canIdsLogin"`
		AppCacheDisable            string `json:"appCacheDisable"`
		HomePageDisplayItemTeacher string `json:"homePageDisplayItemTeacher"`
		LossPwdDesc                string `json:"lossPwdDesc"`
		IsAmpProxy                 string `json:"isAmpProxy"`
		ProvinceID                 string `json:"provinceId"`
		YbSwitch                   string `json:"ybSwitch"`
		Amp3URL                    string `json:"amp3Url"`
		IsOpenFission              string `json:"isOpenFission"`
		IsOpenOauth                string `json:"isOpenOauth"`
		AmpRobotURL                string `json:"ampRobotUrl"`
		MediaVersion               string `json:"mediaVersion"`
		BadHTTPSBlock              string `json:"badHttpsBlock"`
		FaqForumID                 string `json:"faqForumId"`
		CampusReqProxy             string `json:"campusReqProxy"`
		AppStoreURL                string `json:"appStoreUrl"`
		StudentVersion             string `json:"studentVersion"`
		CircleCanSeeOffCampus      string `json:"circleCanSeeOffCampus"`
		ContactDisplayItem         string `json:"contactDisplayItem"`
		ContactDisplayItemTeacher  string `json:"contactDisplayItemTeacher"`
		HomeFirstShow              string `json:"homeFirstShow"`
		AllowSendMsg               int    `json:"allowSendMsg"`
		TeacherVersion             string `json:"teacherVersion"`
		YibanBuild                 int    `json:"yibanBuild"`
		FreshPostRange             string `json:"freshPostRange"`
		YibanAuthType              string `json:"yibanAuthType"`
		CanInteractive             int    `json:"canInteractive"`
		SecondHandSwitch           string `json:"secondHandSwitch"`
		YwtStatus                  string `json:"ywtStatus"`
		YwtPrefixURL               string `json:"ywtPrefixUrl"`
		YwtServiceURL              string `json:"ywtServiceUrl"`
		CollegeTown                string `json:"collegeTown"`
		HasOpenMessageFresh        string `json:"hasOpenMessageFresh"`
	} `json:"data"`
}

func (schoolInfo *SchoolInfo) IsEmpty() bool {
	return reflect.DeepEqual(schoolInfo, SchoolInfo{})
}

type User struct {
	UserName       string   `json:"user_name"`
	PassWord       string   `json:"pass_word"`
	Longitude      string   `json:"longitude"` //当前位置经度
	Latitude       string   `json:"latitude"`  //纬度
	AbnormalReason string   `json:"abnormal_reason"`
	Address        string   `json:"address"`
	Time           string   `json:"time"` //提交时的时间戳
	Sign           string   `json:"sign"`
	MorningTime    string   `json:"morning_time"`
	NoonTime       string   `json:"noon_time"`
	EveningTime    string   `json:"evening_time"`
	FileList       []string `json:"file_list"`
	//SignTime string `json:"sign_time"`
	//SchoolName string `json:"school_name"`
	//QA map[string]string `json:"qa"`
}

type TaskDeatil struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Datas   struct {
		SignInstanceWid string `json:"signInstanceWid"`
		SignMode        int    `json:"signMode"`
		SignRate        string `json:"signRate"`
		SignCondition   int    `json:"signCondition"`
		TaskType        string `json:"taskType"`
		TaskName        string `json:"taskName"`
		TaskDesc        string `json:"taskDesc"`
		QrCodeRcvdUsers []struct {
			TargetWid      string `json:"targetWid"`
			TargetType     string `json:"targetType"`
			TargetName     string `json:"targetName"`
			TargetGrade    string `json:"targetGrade"`
			TargetDegree   string `json:"targetDegree"`
			TargetUserType string `json:"targetUserType"`
		} `json:"qrCodeRcvdUsers"`
		SenderUserName      string      `json:"senderUserName"`
		CurrentTime         string      `json:"currentTime"`
		SingleTaskBeginTime interface{} `json:"singleTaskBeginTime"`
		SingleTaskEndTime   interface{} `json:"singleTaskEndTime"`
		RateSignDate        string      `json:"rateSignDate"`
		RateTaskBeginTime   string      `json:"rateTaskBeginTime"`
		RateTaskEndTime     string      `json:"rateTaskEndTime"`
		SignStatus          string      `json:"signStatus"`
		SignTime            interface{} `json:"signTime"`
		SignPhotoURL        interface{} `json:"signPhotoUrl"`
		SignType            interface{} `json:"signType"`
		ChangeTime          interface{} `json:"changeTime"`
		ChangeActorName     string      `json:"changeActorName"`
		SignPlaceSelected   []struct {
			Wid            interface{} `json:"wid"`
			PlaceWid       interface{} `json:"placeWid"`
			Address        string      `json:"address"`
			Longitude      string      `json:"longitude"`
			Latitude       string      `json:"latitude"`
			Radius         int         `json:"radius"`
			CreatorUserWid interface{} `json:"creatorUserWid"`
			CreatorUserID  interface{} `json:"creatorUserId"`
			CreatorName    interface{} `json:"creatorName"`
			CurrentStatus  interface{} `json:"currentStatus"`
			IsShare        interface{} `json:"isShare"`
		} `json:"signPlaceSelected"`
		IsPhoto       int           `json:"isPhoto"`
		Photograph    []interface{} `json:"photograph"`
		DownloadURL   string        `json:"downloadUrl"`
		LeaveAppURL   string        `json:"leaveAppUrl"`
		CatQrURL      string        `json:"catQrUrl"`
		SignAddress   interface{}   `json:"signAddress"`
		Longitude     string        `json:"longitude"`
		Latitude      string        `json:"latitude"`
		IsMalposition int           `json:"isMalposition"`
		SignedStuInfo struct {
			UserWid        string      `json:"userWid"`
			UserID         string      `json:"userId"`
			UserName       string      `json:"userName"`
			Sex            string      `json:"sex"`
			Nation         string      `json:"nation"`
			Mobile         interface{} `json:"mobile"`
			Grade          string      `json:"grade"`
			Dept           string      `json:"dept"`
			Major          string      `json:"major"`
			Cls            string      `json:"cls"`
			SchoolStatus   interface{} `json:"schoolStatus"`
			Malposition    interface{} `json:"malposition"`
			StuDormitoryVo struct {
				Area     string `json:"area"`
				Building string `json:"building"`
				Unit     string `json:"unit"`
				Room     string `json:"room"`
				Sex      string `json:"sex"`
			} `json:"stuDormitoryVo"`
			ExtraFieldItemVos []struct {
				FieldIndex            int         `json:"fieldIndex"`
				ExtraTitle            string      `json:"extraTitle"`
				ExtraDesc             string      `json:"extraDesc"`
				ExtraFieldItemWid     string      `json:"extraFieldItemWid"`
				ExtraFieldItem        interface{} `json:"extraFieldItem"`
				IsExtraFieldOtherItem string      `json:"isExtraFieldOtherItem"`
				IsAbnormal            string      `json:"isAbnormal"`
			} `json:"extraFieldItemVos"`
		} `json:"signedStuInfo"`
		IsNeedExtra   int         `json:"isNeedExtra"`
		IsAllowUpdate bool        `json:"isAllowUpdate"`
		UpdateLimit   interface{} `json:"updateLimit"`
		LeftNum       interface{} `json:"leftNum"`
		ExtraField    []struct {
			Wid             int    `json:"wid"`
			Title           string `json:"title"`
			Description     string `json:"description"`
			HasOtherItems   int    `json:"hasOtherItems"`
			ExtraFieldItems []struct {
				Content      string      `json:"content"`
				Wid          int         `json:"wid"`
				IsOtherItems int         `json:"isOtherItems"`
				Value        interface{} `json:"value"`
				IsSelected   interface{} `json:"isSelected"`
				IsAbnormal   bool        `json:"isAbnormal"`
			} `json:"extraFieldItems"`
		} `json:"extraField"`
		ExtraFieldItemVos []interface{} `json:"extraFieldItemVos"`
	} `json:"datas"`
}

func (task *TaskDeatil) IsEmpty() bool {
	return reflect.DeepEqual(task, TaskDeatil{})
}

type Result struct {
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

type SettingData struct {
	SchoolName        string            `json:"school_name"`
	QuestionAndAnswer map[string]string `json:"question_and_answer"`
	RunPort           string            `json:"run_port"`
	CallBackApi       string            `json:"call_back_api"`
	MorningHour       string            `json:"morning_hour"`
	MorningMin        string            `json:"morning_min"`
	NoonHour          string            `json:"noon_hour"`
	NoonMin           string            `json:"noon_min"`
	EveningHour       string            `json:"evening_hour"`
	EveningMin        string            `json:"evening_min"`
	SignStep          string            `json:"sign_step"`
	LoginApiList      []string          `json:"login_api_list"`
}

type CallBackData struct {
	Status     int64  `json:"status"`
	UserName   string `json:"user_name"`
	SignResult string `json:"result"`
}

type U struct {
}
