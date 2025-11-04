import {fetchEventSource} from "../sse/index.js";
import { store } from '@/store/index'
import Print from '../utils/printPlus2.js'
import {parseSub, convertLatexSyntax} from "@/utils/util.js"
import {mapActions, mapGetters} from 'vuex'
import {i18n} from "@/lang"

var originalFetch = window.fetch;

import {md} from './marksown-it'
import $ from './jquery.min.js'
import { file } from "jszip";


export default {
    data() {
        return {
            isTestChat:false,
            defaultUrl: '/img/smart/logo.png',
            inputVal: '',
            eventSource: null,
            ctrlAbort: null,
            sseParams: {},
            sseResponse:{},
            echo: true,
            conversationId: '', //会话id
            chatList: [],
            reminderList: [],
            queryFilePath: '',
            stopBtShow: false,
            origin:window.location.origin,
            reconnectCount:0,
            isEnd: true,
            sseApi: "/user/api/v1/assistant/stream",
            rag_sseApi:"/user/api/v1/rag/chat",
            token:store.getters['user/token'],
            lastIndex:0,
            query:'',
            isStoped : false,
            access_token:'',
            runResponse: "",
            fileList: [],  // 文件列表
        };
    },
    created() {
        const vuex = JSON.parse(localStorage.getItem("access_cert"));
        if(vuex){
            this.access_token = vuex.user.token;
        }

    },
    mounted(){
        //this.addVisibilitychangeEvent()
    },
    beforeDestroy() {
        this.setStoreSessionStatus(-1)
        this.stopEventSource()
        this._print && this._print.stop()
    },
    computed: {
        ...mapGetters('app', ['sessionStatus']),
    },
    methods: {
        ...mapActions("app", ["setStoreSessionStatus"]),
        newFetch(url, options){
                // 可以调用原始的 fetch 函数
                if(this.isStoped){
                    return 
                }
                return originalFetch(url, options).then(response => {
                    // 可以在这里修改响应或者添加额外的处理
                    console.log('原始 fetch 响应:', response);
                    let query = this.query

                    if(response.status != 200){
                        let me = this
                        try{
                            const stream = response.body
                           
                            const reader = stream.getReader();
                            const decoder = new TextDecoder('utf-8');

                            function readStream() {
                                reader.read().then(({ done, value }) => {
                                if (done) {
                                    console.log('Stream complete');
                                    reader.releaseLock()
                                    return;
                                }
                        
                                // Decode and process each chunk of data.
                                const decodedValue = decoder.decode(value, { stream: true });
                                                               
                                if(decodedValue){
                                    let msg = JSON.parse(decodedValue).msg
                                    me.setStoreSessionStatus(-1)
                                    var fillData = {
                                        "query": query,
                                        qa_type:0,
                                        finish:1,
                                        response:msg,  //非代码文本使用自定义转换规则，不使用markdown,(markdown渲染会导致卡顿或样式丢失)
                                        oriResponse:"",
                                        searchList:[]  //过滤包含yunyingshang文件的出处
                                    }
                                    this.runResponse = msg
                                }                                
                                readStream();
                        
                                // Continue reading the stream.
                                }).catch(err => {
                                    console.error('Reading stream failed1:', err);
                                });
                            }
                        
                        readStream();
                        me.isStoped = true
                        }catch(e){
                            console.error('Reading stream failed:', e);
                        }

                    }

                    return response;
                }).catch((err)=>{
                    this.$message.warning("连接失败，请稍后重试")
                    this.isEnd = true
                    this.setStoreSessionStatus(-1)
                    this.runDisabled = false
                });
        },
        ...mapActions("app", ["setStoreSessionStatus"]),
        queryCopy(text){
            this.setPrompt(text)
        },
        /*过滤掉markdown中自定义的行号*/
        getContentInBraces(shtml) {
            let temp = document.createElement('div')
            temp.setAttribute('id','temp')
            temp.innerHTML = shtml
            document.body.appendChild(temp)
            $(temp).find('.line-num').remove()
            return temp.innerText
        },
        addCopyClick(){
            let copys =document.getElementsByClassName('copy-btn')
            for(var i=0;i<copys.length;i++){
                copys[i].addEventListener('click', (e)=> {
                    let innerText = e.target.parentNode.nextElementSibling.innerText
                    this.$copy(innerText)
                    e.target.innerText ='复制成功'
                    //document.getElementById('temp') && document.getElementById('temp').remove() //复制完成把临时生成的元素删除
                    setTimeout(()=>{
                        e.target.innerText ='复制'
                    },1500)
                })
            }
        },
        // 填充开场白
        setProloguePrompt(val) {
            this.$refs['editable'].setPrompt(val)
            this.preSend()
        },
        //获取上传的文件
        getFileIdList() {
            let list = this.$refs['editable'].getFileIdList()  //[{fileId:"b67e79615"url:"b4d0a8.jpg"}]
            let fileIds = []
            this.queryFilePath = ''
            if (list.length) {
                fileIds = list.map(n => {
                    return n.fileId
                })
                this.queryFilePath = list[0].url
            }
            return fileIds.join(',')
        },
        mouseEnter(n) {
            n.hover = true
        },
        mouseLeave(n) {
            n.hover = false
        },
        setSessionStatus(status) {
            this.setStoreSessionStatus(status)
        },
        setSseParams(data) {
            this.sseParams = data
        },
        doragSend(){
            this.stopBtShow = true
            this.isStoped= false
            let _history = this.$refs['session-com'].getList()
            this.sendEventStream(this.inputVal,'', _history.length)
        },
        sendEventStream(prompt, msgStr, lastIndex){
            if (this.sessionStatus === 0) {
                this.$message.warning('上个问题没有回答完！')
                return
            }

            this.sseResponse = {}
            this.setStoreSessionStatus(0)
            this.clearInput()
            let params = {
                query: prompt, 
                pending: true, 
                responseLoading: true, 
                requestFileUrls:[],
                fileList:this.fileList,
                pendingResponse:''
            }
            this.$refs['session-com'].pushHistory(params)
            let endStr = ''
            this._print = new Print({
                onPrintEnd: () => {
                    // this.setStoreSessionStatus(-1)
                }
            })

            let history_list = this.$refs['session-com'].getSessionData();
            const history = history_list['history'].length > 1 ? history_list['history'][history_list['history'].length - 2]['history'] : [];
            this.ctrlAbort = new AbortController();
            const userInfo = this.$store.state.user.userInfo || {};
            this.eventSource = new fetchEventSource(this.origin + this.rag_sseApi, {
                method: 'POST',
                headers: {
                    "Content-Type": 'application/json',
                    'Authorization': 'Bearer ' + this.token,
                    "x-user-id": userInfo.uid,
                    "x-org-id": userInfo.orgId
                },
                signal: this.ctrlAbort.signal,
                body: JSON.stringify({...this.sseParams,'history':history}),
                openWhenHidden: true, //页面退至后台保持连接
                onopen: async(e) => {
                    //console.log("已建立SSE连接~",new Date().getTime());
                    if (e.status !== 200) {
                        try {
                            const errorData = await e.json();
                            let commonData = {
                                ...this.sseParams,
                                "query": prompt,
                            }
                            let fillData = {
                                ...commonData,
                                "response": errorData.msg                                
                            }
                            this.$refs['session-com'].replaceLastData(lastIndex, fillData)
                        } catch (e) {
                            const text = await e.text();
                            this.$message.error(text || '未知错误');
                        }

                        this.stopEventSource();
                        this.setStoreSessionStatus(-1);
                        return;
                    }
                },
                onmessage: (e) => {
                    if (e && e.data) {
                        let data;
                        try {
                            data = JSON.parse(e.data);
                            console.log('===>',new Date().getTime(), data);
                        } catch (error) {
                            return; // 如果解析失败，直接返回，不处理这条消息
                        }
                        
                        this.sseResponse = data;
                        //待替换的数据，需要前端组装
                        let commonData = {
                            ...this.sseResponse,
                            ...this.sseParams,
                            "query": prompt,
                            "fileName":'',
                            "fileSize":'',
                            "response": '',
                            "filepath": '',
                            "requestFileUrls":'',
                            "searchList": this.sseResponse.data && this.sseResponse.data.searchList ? this.sseResponse.data.searchList: [],
                            "gen_file_url_list": [],
                            "thinkText":'思考中',
                            "isOpen":true,
                            "citations":[]
                        }
                        if(data.code === 0 || data.code === 1){
                            //finish 0：进行中  1：关闭   2:敏感词关闭
                            let _sentence = data.data.output;
                                this._print.print(
                                    {
                                        response:_sentence,
                                        finish:data.finish
                                    },
                                    commonData,
                                    (worldObj,search_list) => {
                                        this.setStoreSessionStatus(0)
                                        endStr += worldObj.world  
                                        endStr = convertLatexSyntax(endStr)
                                        endStr = parseSub(endStr, lastIndex)     
                                        let fillData = {
                                            ...commonData,
                                            "response":md.render(endStr),
                                            oriResponse:endStr,
                                            finish:worldObj.finish,
                                            searchList:(search_list && search_list.length) ? search_list.map(n => ({
                                                  ...n,
                                                  snippet: md.render(n.snippet)
                                                }))
                                            : []
                                        }
                                        this.$refs['session-com'].replaceLastData(lastIndex, fillData)
                                        if(worldObj.isEnd && worldObj.finish === 1){
                                            this.setStoreSessionStatus(-1)
                                        }
                                    })
                            // this.$nextTick(()=>{
                            //     this.$refs['session-com'].scrollBottom()
                            // })
                        }else if(data.code === 7 || data.code === -1){
                            this.setStoreSessionStatus(-1)
                            let fillData = {
                                ...commonData,
                                "response": data.message                            
                            }
                            this.$refs['session-com'].replaceLastData(lastIndex, fillData)
                        }
                    }
                },
                onclose: () => {
                    console.log('===> eventSource onClose')
                    this.setStoreSessionStatus(-1)//关闭后改变状态
                    this.sseOnCloseCallBack()
                },
                onerror: (e) => {
                    console.log("服务连接异常请重试！");
                    if (e.readyState === EventSource.CLOSED) {
                        console.log("connection is closed");
                    } else {
                        console.log("Error occured", e);
                    }
                    this.stopEventSource()//前端主动关闭连接
                    this.setStoreSessionStatus(-1)//关闭后改变状态
                }
            });
                      
        },
        doSend(params) {
            this.stopBtShow = true
            this.isStoped= false
            let _history = this.$refs['session-com'].getList()
            this.sendEventSource(this.inputVal, '', _history.length)
        },
        sendEventSource(prompt, msgStr, lastIndex) {
            console.log('####  sendEventSource',new Date().getTime())
            const userInfo = this.$store.state.user.userInfo || {}
            if (this.sessionStatus === 0) {
                this.$message.warning('上个问题没有回答完！')
                return
            }

            this.sseResponse = {}
            //发送问题后不允许继续提问
            this.setStoreSessionStatus(0)

            this.clearInput()

            let params = {
                query: prompt, 
                pending: true, 
                responseLoading: true, 
                requestFileUrls: this.queryFilePath?[this.queryFilePath]:[],
                fileList:this.fileList,
                pendingResponse:''
            }
            //正式环境传模型参数
            this.$refs['session-com'].pushHistory(params)
            
            let endStr = ''
            this._print = new Print({
                onPrintEnd: () => {
                }
            })

            let data = null;
            let headers = null;
            //判断是是不是openurl对话
            if(this.type === 'agentChat'){
                this.sseApi = "/user/api/v1/assistant/stream";
                const trial = this.isTestChat ? true : false
                data = {
                    ...this.sseParams,
                    prompt,
                    trial
                };
                headers = {
                    "Content-Type": 'application/json',
                    'Authorization': 'Bearer ' + this.token,
                    "x-user-id": userInfo.uid,
                    "x-org-id": userInfo.orgId
                }
            }else{
                this.sseApi = `/service/url/openurl/v1/agent/${this.sseParams.assistantId}/stream`;
                data = {
                   conversationId:this.sseParams.conversationId, 
                   prompt
                }
                headers = {
                    'X-Client-ID':this.getHeaderConfig().headers['X-Client-ID']
                };
            }

            this.ctrlAbort = new AbortController();
            this.eventSource = new fetchEventSource(this.origin + this.sseApi, {
                method: 'POST',
                headers,
                signal: this.ctrlAbort.signal,
                body: JSON.stringify(data),
                openWhenHidden: true, //页面退至后台保持连接
                ...(this.type === 'webChat' && { isOpenUrl: true }),
                onopen: async(e) => {
                    console.log("已建立SSE连接~",new Date().getTime());
                    if (e.status !== 200) {
                        try {
                            const errorData = await e.json();
                            let commonData = {
                                ...this.sseParams,
                                "query": prompt,
                            }
                            let fillData = {
                                ...commonData,
                                "response": errorData.msg                                
                            }
                            this.$refs['session-com'].replaceLastData(lastIndex, fillData)
                        } catch (e) {
                            const text = await e.text();
                            this.$message.error(text || '未知错误');
                        }

                        this.stopEventSource();
                        this.setStoreSessionStatus(-1);
                        return;
                    }
                },
                onmessage: (e) => {
                    if (e && e.data) {
                        let data = JSON.parse(e.data)
                        console.log('===>',new Date().getTime(),data)
                        this.sseResponse = data
                        //待替换的数据，需要前端组装
                        let commonData = {
                            ...data,
                            ...this.sseParams,
                            "query": prompt,
                            "fileList":this.fileList,
                            "response": '',
                            "filepath": data.file_url || '',
                            "requestFileUrls": this.queryFilePath?[this.queryFilePath] : data.requestFileUrls,
                            "searchList": data.search_list || [],
                            "gen_file_url_list":data.gen_file_url_list || [],
                            "thinkText":i18n.t('agent.thinking'),
                            'toolText':'使用工具中...',
                            "isOpen":true,
                            "showScrollBtn":null,
                            "citations":[]
                        }

                        if(data.code === 0){
                            //finish 0：进行中  1：关闭   2:敏感词关闭
                            let _sentence = data.response
                                this._print.print(
                                    {
                                        response:_sentence,
                                        finish:data.finish
                                    },
                                    commonData,
                                    (worldObj,search_list) => {
                                        this.setStoreSessionStatus(0)
                                        endStr += worldObj.world
                                        endStr = convertLatexSyntax(endStr)
                                        endStr = parseSub(endStr, lastIndex)
                                        const finalResponse = String(endStr)
                                        let fillData = {
                                            ...commonData,
                                            response: [0,1,2,3,4,6,20,21,10].includes(commonData.qa_type)?md.render(finalResponse):finalResponse.replaceAll('\n-','<br/>•').replaceAll('\n','<br/>'),
                                            finish:worldObj.finish,
                                            oriResponse:endStr,
                                            searchList:(search_list && search_list.length) ? search_list.map(n => ({
                                                  ...n,
                                                  snippet: md.render(n.snippet)
                                                }))
                                            : []
                                        }
                                        this.$refs['session-com'].replaceLastData(lastIndex, fillData)
                                        if(worldObj.finish !== 0){
                                            if(worldObj.finish === 4){
                                                let fillData = {
                                                    ...commonData,
                                                    "response":i18n.t('yuanjing.sensitiveTips')
                                                }
                                                this.$refs['session-com'].replaceLastData(lastIndex, fillData)
                                            }
                                            this.setStoreSessionStatus(-1)
                                        }
                                        if(worldObj.isEnd && worldObj.finish === 1){
                                          this.setStoreSessionStatus(-1)
                                       }
                                    })

                            this.$nextTick(()=>{
                                this.$refs['session-com'].scrollBottom()
                            })

                        }else if(data.code === 7 || data.code === -1 || data.code === 1){
                            this.setStoreSessionStatus(-1)
                            let fillData = {
                                ...commonData,
                                "response": data.message                               
                            }
                            this.$refs['session-com'].replaceLastData(lastIndex, fillData)
                        }
                    }
                },
                onclose: () => {
                    console.log('===> eventSource onClose')
                    this.setStoreSessionStatus(-1)//关闭后改变状态
                    this.sseOnCloseCallBack()
                },
                onerror: (e) => {
                    console.log("服务连接异常请重试！");
                    if (e.readyState === EventSource.CLOSED) {
                        console.log("connection is closed");
                    } else {
                        console.log("Error occured", e);
                    }
                    this.stopEventSource()//前端主动关闭连接
                    this.setStoreSessionStatus(-1)//关闭后改变状态
                }
            });
        },
        preStop() {
            //获取已经拿到的全部回答,一次性回显出来
            this.sseOnCloseCallBack(true)
        },
        sseOnCloseCallBack(isStoped){
            this.stopEventSource()
            //图文问答不使用打字机
           /* if(this.sseResponse.qa_type === 6){
                return
            }*/
            //主动停止
            if(isStoped) {
                this.stopAndEcho()
            }else{
                //收到onclose,且使用的是文生代码
                if(this.sseResponse.qa_type === 4){
                    this.stopAndEcho()
                }else{
                    //接口405等
                    let history_list = []
                    let lastIndex = history_list.length - 1
                    let lastRQ = history_list[lastIndex]
                    let endStr =this._print.getAllworld()
                    endStr = convertLatexSyntax(endStr)
                    // 替换标签
                    endStr = parseSub(endStr)
                    // 如果返回有结果，则在结束时不展示“本次回答已终止”
                    this.runResponse = md.render(endStr)
                    this.runDisabled = false
                    this.setStoreSessionStatus(-1)
                }
            }
        },
        stopAndEcho(){
            //暂存已经收到的所有response
            let endResponse = this._print.getAllworld()

            this._print && this._print.stop()

            setTimeout(()=>{
                this.setStoreSessionStatus(-1)

                let history_list = []
                let lastIndex = history_list.length - 1
                let lastRQ = history_list[lastIndex]
                if(endResponse){
                    endResponse = convertLatexSyntax(endResponse)
                    // 替换标签
                    endResponse = parseSub(endResponse)
                    this.runResponse = md.render(endResponse)
                    this.runDisabled = false
                    this.$nextTick(()=>{
                        this.addCopyClick()
                    })

                }else{
                    if(Object.keys(this.sseResponse).length !== 0 && this.sseResponse.code !== 7){
                        this.runResponse ="本次回答已被终止"
                        this.setStoreSessionStatus(-1)
                    }else{
                        this.stopEventSource();
                        this.setStoreSessionStatus(-1)
                        this.$refs['session-com'].stopPending();
                    }
                }
            },15)
        },
        stopEventSource() {
            this.ctrlAbort && this.ctrlAbort.abort()
            this.eventSource = null
        },
        refreshLastSession(){
            let endResponse = this._print.getAllworld()
            let history_list =[]
            let lastIndex = history_list.length - 1
            let lastRQ = history_list[lastIndex]
            // this.$refs['session-com'].replaceLastData(lastIndex, {
            //     ...lastRQ,
            //     response: endResponse
            // })
        },
        setPrompt(data) {
            this.$refs['editable'].setPrompt(data)
        },
        clearInput() {
            this.$refs.editable.clearInput()
            this.$refs.editable.clearFile()
            this.inputVal = ''
            this.fileId = ''
        },
        clearPageHistory(){
            this.$refs['session-com'] && this.$refs['session-com'].clearData()
            this.$refs.editable && this.clearInput()
        },
        clearHistory() {
            this.stopBtShow = false
            this.clearPageHistory()
            // this.doDeleteHistory()
        },
        refresh() {
            let history_list = this.$refs['session-com'].getList();
            let _history = history_list[history_list.length - 1];
            let inputVal = _history.query;
            let fileInfo = _history.fileInfo ? _history.fileInfo : [];
            let fileList = _history.fileList ? _history.fileList : [];
            this.preSend(inputVal,fileList,fileInfo);
        }
    }
};