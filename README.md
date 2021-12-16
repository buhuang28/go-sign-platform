# go-sign-platform-vcl

用Golang编写的某日校园签到平台。
仅限学习
![Image text](https://github.com/buhuang28/go-sign-platform/blob/main/QQ%E5%9B%BE%E7%89%8720210112161446.png)
应该是适配的很多学校的签到了，如果签到需要上传图片，请把FuckForm函数替换成SchoolTool.go中的FuckForm2。另外HTML需要加上前端上传图片的代码。

需要配合 https://github.com/ZimoLoveShuang/wisedu-unified-login-api 获取cookie去获取签到任务、提交签到信息

GUI使用govcl绘制：https://github.com/ying32/govcl

后端框架使用了Gin

签到主要算法来自子墨：https://github.com/ZimoLoveShuang/auto-sign

运行go_build.bat即可编译
如果签到需要图片请查看源码后自行修改一下，有注释。

或者把SchoolTool.go中的UploadPicture函数注释掉。
