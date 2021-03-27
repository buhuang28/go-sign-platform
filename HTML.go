package main

const htmlfile = `<!DOCTYPE html>
<html>
<head>
   <meta charset="UTF-8">
   <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/element-ui@2.14.1/lib/theme-chalk/index.min.css">
   <script src="https://cdn.bootcss.com/blueimp-md5/2.10.0/js/md5.js"></script>
   <script src="https://cdn.jsdelivr.net/npm/vue@2.6.12/dist/vue.min.js"></script>
   <script src="https://cdn.jsdelivr.net/npm/element-ui@2.14.1/lib/index.min.js"></script>
</head>
<body>
<div id="app">
   <el-form ref="form" :model="form" label-width="100px">
       <el-form-item label="用户名:">
           <el-input v-model="form.user_name"></el-input>
       </el-form-item>
       <el-form-item label="密码:">
           <el-input v-model="form.pass_word" show-password></el-input>
       </el-form-item>
       <el-form-item label="签到地址:">
           <el-input v-model="form.address"></el-input>
       </el-form-item>
       <el-form-item label="签到经度:">
           <el-input v-model="form.longitude"></el-input>
       </el-form-item>
       <el-form-item label="签到纬度:">
           <el-input v-model="form.latitude"></el-input>   <el-button type="primary" @click="getLoc">获取签到经纬度</el-button>
       </el-form-item>
       <el-form-item label="早上签到时间:">
           <el-time-select
                   v-model="form.morning_time"
                   :picker-options="{
                       start: '7:05',
                       step: '00:01',
                       end: '08:30'
                       }"
                   placeholder="早上签到时间">
           </el-time-select>
       </el-form-item>
       <el-form-item label="中午签到时间:">
           <el-time-select
                   v-model="form.noon_time"
                   :picker-options= "{
						start: '12:05',
                       step: '00:01',
                       end: '13:30'
                       }"
                   placeholder="中午签到时间">
           </el-time-select>
       </el-form-item>
       <el-form-item label="晚上签到时间:">
           <el-time-select
                   v-model="form.evening_time"
                   :picker-options="{
						start: '18:05',
                       step: '00:01',
                       end: '19:30'
                       }"
                   placeholder="晚上签到时间">
           </el-time-select>
       <el-form-item>
           <el-button type="primary" @click="onSubmit">提交签到信息</el-button>
       </el-form-item>
   </el-form>
   <div width="100em">
       <a href="https://jq.qq.com/?_wv=1027&k=HxwFRLvn">
           不慌的自动签到平台(免费+开源,如果收费肯定是被骗了)交流群：696129128
       </a>
   </div>
</div>
</body>
<!-- import Vue before Element -->
<script src="https://cdn.jsdelivr.net/npm/axios@0.21.1/dist/axios.min.js"></script>
<script>window["\x64\x6f\x63\x75\x6d\x65\x6e\x74"]["\x74\x69\x74\x6c\x65"] = '\u4e0d\u614c\u7684\u81ea\u52a8\u7b7e\u5230\x28\u7b7e\u5230\u7b97\u6cd5\u6765\u81ea\u5b50\u58a8\x29\u4ea4\u6d41\u7fa4\uff1a\x36\x39\x36\x31\x32\x39\x31\x32\x38';</script>
<script>
   Vue.prototype.$http = axios
   new Vue({
       el: '#app',
       data() {
           return {
               form: {
                   user_name: '',
                   pass_word: '',
                   longitude: '',
                   latitude: '',
                   address: '',
                   morning_time:'',
                   noon_time:'',
                   evening_time:'',
                   sign:''
               }
           }
       },
       created () {
       },
       methods: {
           onSubmit() {
              this['$refs']['form']['validate'](async _0xa63ae2=>{var _0x11b61e={'iVlPI':'auto_sign','mnydf':function(_0x1d7c3d,_0x5e8e31){return _0x1d7c3d<_0x5e8e31;},'Kxjrj':function(_0x407382,_0x3fe657){return _0x407382==_0x3fe657;},'TWvDg':'sign','PjRyU':function(_0x4ca838,_0x3c029d){return _0x4ca838+_0x3c029d;},'VzxuB':function(_0x1a11b0,_0x345fef){return _0x1a11b0+_0x345fef;},'RtdRB':function(_0x3b515c,_0x23f36a){return _0x3b515c(_0x23f36a);},'jLmEU':'/addUser','fSfFh':function(_0x4f4027,_0x1090e0){return _0x4f4027!=_0x1090e0;}};if(!_0xa63ae2)return;let _0x2325a5=Date['parse'](new Date())['toString']();_0x2325a5=_0x2325a5['substr'](0x0,0xa);this['form']['time']=_0x2325a5;var _0x3f9351=new Array();let _0x6f7ebd=0x0;for(var _0x5bb6aa in this['form']){_0x3f9351[_0x6f7ebd]=_0x5bb6aa;_0x6f7ebd++;}_0x3f9351['sort']();const _0x2dd442=_0x11b61e['iVlPI'];let _0x5baa42=new String('');for(j=0x0;_0x11b61e['mnydf'](j,_0x3f9351['length']);j++){if(_0x11b61e['Kxjrj'](_0x3f9351[j],_0x11b61e['TWvDg'])){continue;}_0x5baa42+=_0x11b61e['PjRyU'](_0x3f9351[j],this['form'][_0x3f9351[j]]);}_0x5baa42=_0x11b61e['VzxuB'](_0x11b61e['VzxuB'](_0x11b61e['iVlPI'],_0x5baa42),_0x2dd442);_0x5baa42=_0x11b61e['RtdRB'](md5,_0x5baa42);this['form']['sign']=_0x5baa42;const {data:res}=await this['$http']['post'](_0x11b61e['jLmEU'],this['form']);if(_0x11b61e['fSfFh'](res['Status'],0xc8)){_0x11b61e['RtdRB'](alert,res['message']);}else{this['$message']['success']('成功');}});
           },
           getLoc(){
               location.href="http://api.map.baidu.com/lbsapi/getpoint/index.html";
           },
       }
   })
</script>
</html>`

//const htmlfile =  `<!DOCTYPE html>
//<html>
//<head>
//    <meta charset="UTF-8">
//    <link rel="stylesheet" href="https://unpkg.com/element-ui/lib/theme-chalk/index.css">
//    <script src="https://cdn.bootcss.com/blueimp-md5/2.10.0/js/md5.js"></script>
//    <script src="//unpkg.com/vue/dist/vue.js"></script>
//    <script src="//unpkg.com/element-ui@2.14.1/lib/index.js"></script>
//</head>
//<body>
//<div id="app">
//    <el-form ref="form" :model="form" label-width="100px">
//        <el-form-item label="用户名:">
//            <el-input v-model="form.user_name"></el-input>
//        </el-form-item>
//        <el-form-item label="密码:">
//            <el-input v-model="form.pass_word" show-password></el-input>
//        </el-form-item>
//        <el-form-item label="签到地址:">
//            <el-input v-model="form.address"></el-input>
//        </el-form-item>
//        <el-form-item label="签到经度:">
//            <el-input v-model="form.longitude"></el-input>
//        </el-form-item>
//        <el-form-item label="签到纬度:">
//            <el-input v-model="form.latitude"></el-input>   <el-button type="primary" @click="getLoc">获取签到经纬度</el-button>
//        </el-form-item>
//        <el-form-item label="早上签到时间:">
//            <el-time-select
//                    v-model="form.morning_time"
//                    :picker-options="{
//                        start: '7:05',
//                        step: '00:01',
//                        end: '08:30'
//                        }"
//                    placeholder="早上签到时间">
//            </el-time-select>
//        </el-form-item>
//        <el-form-item label="中午签到时间:">
//            <el-time-select
//                    v-model="form.noon_time"
//                    :picker-options= "{
//						start: '12:05',
//                        step: '00:01',
//                        end: '13:30'
//                        }"
//                    placeholder="中午签到时间">
//            </el-time-select>
//        </el-form-item>
//        <el-form-item label="晚上签到时间:">
//            <el-time-select
//                    v-model="form.evening_time"
//                    :picker-options="{
//						start: '18:05',
//                        step: '00:01',
//                        end: '19:30'
//                        }"
//                    placeholder="晚上签到时间">
//            </el-time-select>
//        </el-form-item>
//        <el-upload
//                class="upload-demo"
//                action="/addPic"
//                :on-preview="handlePreview"
//                :on-remove="handleRemove"
//                :on-success = "imgSucess"
//                list-type="picture">
//            <el-button size="small" type="primary">点击上传图片</el-button>
//            <div slot="tip" class="el-upload__tip">上传签到用的图片，可多张</div>
//        </el-upload>
//        <el-form-item>
//            <el-button type="primary" @click="onSubmit">提交签到信息</el-button>
//        </el-form-item>
//    </el-form>
//    <div width="100em">
//        <a href="https://jq.qq.com/?_wv=1027&k=HxwFRLvn">
//            不慌的自动签到平台(免费+开源，如果收费肯定是被骗了).(签到算法来自子墨).交流群：696129128
//        </a>
//    </div>
//</div>
//</body>
//<script src="https://unpkg.com/axios/dist/axios.min.js"></script>
//<script>window["\x64\x6f\x63\x75\x6d\x65\x6e\x74"]["\x74\x69\x74\x6c\x65"] = '\u4e0d\u614c\u7684\u81ea\u52a8\u7b7e\u5230\x28\u7b7e\u5230\u7b97\u6cd5\u6765\u81ea\u5b50\u58a8\x29\u4ea4\u6d41\u7fa4\uff1a\x36\x39\x36\x31\x32\x39\x31\x32\x38';</script>
//<script>
//    Vue.prototype.$http = axios
//    new Vue({
//        el: '#app',
//        data() {
//            return {
//                form: {
//                    file_list:[],
//                    user_name: '',
//                    pass_word: '',
//                    longitude: '',
//                    latitude: '',
//                    address: '',
//                    morning_time:'',
//                    noon_time:'',
//                    evening_time:'',
//                    sign:''
//                }
//            }
//        },
//        created () {
//        },
//        methods: {
//            onSubmit() {
//                this.$refs.form.validate(async valid => {
//                    if (!valid) return
//                    let time = Date.parse(new Date()).toString();//获取到毫秒的时间戳，精确到毫秒
//                    time = time.substr(0,10);
//                    this.form.time = time;
//                    var ar = [];
//                    let i = 0;
//                    for(var index in this.form){
//                        ar[i] = index;
//                        i++;
//                    }
//                    ar.sort();
//                    const sk = "auto_sign";
//                    let s = String("");
//                    for(j = 0; j < ar.length; j++) {
//                        if (ar[j] === "sign") {
//                            continue;
//                        }
//                        s += ar[j]+this.form[ar[j]];
//                    }
//                    s = "auto_sign"+s+sk;
//                    s = md5(s);
//                    this.form.sign = s;
//                    const { data: res } = await this.$http.post('/addUser', this.form)
//                    if (res.Status !== 200) {
//                        alert(res.message);
//                    } else {
//                        this.$message.success('成功')
//                    }
//                })
//            },
//            getLoc(){
//                location.href="http://api.map.baidu.com/lbsapi/getpoint/index.html";
//            },
//            handleRemove(file, file_list) {
//                console.log(file, file_list);
//            },
//            handlePreview(file) {
//                console.log(file);
//            },
//            imgSucess(data){
//                if (data !== "" && data !== null && data !== undefined) {
//                    this.form.file_list.push(data);
//                    this.$message.success("添加图片成功");
//                }else {
//                    this.$message.success("添加图片失败");
//                }
//            }
//        }
//    })
//</script>
//</html>`