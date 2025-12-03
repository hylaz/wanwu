<template>
    <!-- 远景大模型 -->
    <div class="full-content flex">
        <el-main class="scroll">
            <div class="smart-center">
                <!--基础配置回显-->
                <div v-show="echo" class="session rl echo">
                    <Prologue  :editForm="editForm" @setProloguePrompt="setProloguePrompt" :isBigModel="true" />
                </div>
                <!--对话-->
                <div  v-show="!echo" class="center-session">
                    <SessionComponentSe
                            ref="session-com"
                            class="component"
                            :sessionStatus="sessionStatus"
                            @clearHistory="clearHistory"
                            @refresh="refresh"
                            @queryCopy="queryCopy"
                            @handleRecommendedQuestion="handleRecommendedQuestion"
                            :defaultUrl="editForm.avatar.path"
                    />
                </div>
                <!--输入框-->
                <div class="center-editable">
                    <div v-show="stopBtShow" class="stop-box">
                        <span v-show="sessionStatus === 0" class="stop" @click="preStop">
                            <img class="stop-icon mdl" :src="require('@/assets/imgs/stop.png')"/>
                            <span class="mdl">停止生成</span>
                        </span>
                        <!-- <span v-show="sessionStatus !== 0" class="stop" @click="refresh">
                            <img class="stop-icon mdl" :src="require('@/assets/imgs/refresh.png')"/>
                            <span class="mdl">{{$t('yuanjing.refresh')}}</span>
                        </span> -->
                    </div>
                    <EditableDivV3
                            ref="editable"
                            source="perfectReminder"
                            :fileTypeArr="fileTypeArr"
                            :currentModel="currentModel"
                            :isModelDisable="isModelDisable"
                            :showModelSelect="false"
                            @preSend="preSend"
                            @modelChange="modelChange"
                            @getModelType="getModelType"
                            @setSessionStatus="setSessionStatus"
                    />
                </div>
            </div>
        </el-main>
    </div>
</template>

<script>
    import SessionComponentSe from './SessionComponentSe'
    import EditableDivV3 from './EditableDivV3'
    import {getConversationList, createConversation, getConversationDetail, deleteConversation,deleteConversationHistory} from '@/api/cubm'
    import Prologue from './Prologue'
    import sseMethod from '@/mixins/sseMethod'
    import {md} from '@/mixins/marksown-it'
    import {mapActions, mapGetters} from 'vuex'
    // import { getTemplateList } from '@/api/prompt'

    export default {
        props:{
            chatType:{
                type: String,
                default:''
            },
            editForm:{
                type:Object,
                default:null
            }
        },
        components: {
            SessionComponentSe,
            EditableDivV3,
            Prologue
        },
        mixins: [sseMethod],
        computed: {
            ...mapGetters('app', ['sessionStatus']),
            ...mapGetters('menu', ['basicInfo']),
            ...mapGetters('user', ['commonInfo']),
        },
        data() {
            return {
                amswerNum:0,
                isModelDisable:false,
                currentModel:null,
                echo: true,
                basicForm: {
                    avatar: '123',//img/tr/deepseek-icon.png-使用接口获取的图片
                    instructions: '456',//我是您的智能小助手，可以帮您思考文案，与您聊天，还可以答疑结果。比如您可以问我：
                    name: 'ffff',//你好，我是DeepSeek-使用接口获取的文案
                    description: 'fsdfdggfh'
                },
                expandForm: {
                    starterPrompts: [
                        // {value: '如何打粉底不卡粉?'},
                        // {value: '我想问怎么化妆皮肤不干?'},
                        // {value: '写一个故宫一日游攻略'}
                    ]
                },
                // fileTypeArr: ['image/*','doc/*'],
                fileTypeArr: ['doc/*'],
                hasDrawer: false,
                drawer: true,
                sseApi: this.$basePath + '/use/model/api/v1/chatllm/stream',
            }
        },
        watch: {
            $route: {
                handler(val, oldval) {

                },
                deep: true
            }
        },
        created() {
            // this.getReminderList(() => {
            //     this.hasDrawer = true
            // })
            // this.getConversationList()
        },
        methods: {
            handleRecommendedQuestion(question){
                this.inputVal = question
                this.preSend(question)
            },
            getModelType(type){
                const dataInfo = this.commonInfo.data.useModel;
                if(type ==='deepseek'){
                   this.expandForm.starterPrompts = dataInfo.useModels[1]['welcomeQuestions']
                   this.basicForm.instructions = dataInfo.useModels[1]['welcomeDesc'];
                   this.basicForm.name = dataInfo.useModels[1]['welcomeText'];
                   this.basicForm.avatar = this.$basePath + '/use/model/api' + dataInfo.useModels[1]['welcomeLogoPath'];
                }else{
                    this.expandForm.starterPrompts = dataInfo.useModels[0]['welcomeQuestions']
                   this.basicForm.instructions = dataInfo.useModels[0]['welcomeDesc'];
                   this.basicForm.name = dataInfo.useModels[0]['welcomeText'];
                   this.basicForm.avatar = this.$basePath + '/use/model/api' + dataInfo.useModels[0]['welcomeLogoPath'];
                }
            },
            //对话列表
            async getConversationList(noInit) {
                let res = await getConversationList({assistantId: this.assistantId, pageSize: 1000, pageNo: 1})
                if (res.code === 0) {
                    if(res.data.list && res.data.list.length > 0){
                        this.chatList = res.data.list.map(n => {
                            return {...n, hover: false, active: false}
                         })
                        if (noInit) {
                            this.chatList[0].active = true  //noInit 是true时，左侧默认选中第一个,但是不要调接口刷新详情
                        } else {
                            this.conversionClick[this.chatList[0]]
                        }
                    }else{
                        this.chatList = []
                    }
                }
            },
            //新建对话
            preCreateConversation() {
                if (this.echo) {
                    this.$message({
                        type: 'info',
                        message: this.$t('yuanjing.changeDialogMsg'),
                        customClass: 'dark-message',
                        iconClass: 'none',
                        duration: 1500
                    })
                    return
                }
                this.isModelDisable = false
                this.conversationId = ''
                this.currentModel = null
                this.echo = true
                this.clearPageHistory()
                this.chatList.forEach(m => {
                    m.active = false
                })
            },
            //切换对话
            async conversionClick(n) {
                this.isModelDisable = true;
                if (this.sessionStatus === 0) {
                    //this.$message.warning('上个问题未答完')
                    return
                }else{
                    this.stopBtShow = false
                }

                this.chatList.forEach(m => {
                    m.active = false
                })
                this.amswerNum = 0
                n.active = true
                this.clearPageHistory()
                this.echo = false

                this.conversationId = n.conversationId
                this.getConversationDetail(this.conversationId,true)
            },
            async getConversationDetail(id,loading){
                loading && this.$refs['session-com'].doLoading()
                let res = await getConversationDetail({conversationId: id, pageSize: 1000, pageNo: 1})
                if (res.code === 0) {
                    let history = res.data.list ? res.data.list.map(n => {
                        return {
                            ...n,
                            query: n.prompt,
                            //response:n.qa_type===4?(marked(n.response)).replaceAll('\\n','<br/>'):n.response.replaceAll('\n-','<br/>•').replaceAll('\n','<br/>'),
                            response:[0,1,2,3,4,5,6,20,21,10].includes(n.qa_type)?md.render(n.response):n.response.replaceAll('\n-','\n•'),
                            oriResponse:n.response,
                            searchList: n.searchList ? JSON.parse(n.searchList) : [],
                            filepath: n.responseFileUrls,
                            "gen_file_url_list":n.responseFileUrls || [],
                            "isOpen":true
                        }
                    }) : []

                    //切换历史记录，选择对应模型
                    if(res.data.list && res.data.list !== null){
                        this.currentModel = {
                        modelId:res.data.list[0]['modelId'],
                        modelVersion:res.data.list[0]['modelVersion']
                        }
                    }else{
                        this.currentModel = null
                    }
                    this.$refs['session-com'].replaceHistory(history)
                    this.$nextTick(()=>{
                        this.addCopyClick()
                    })
                }
            },
            //删除对话
            async preDelConversation(n) {
                //todo 给所有的点击事件统一添加拦截
                if (this.sessionStatus === 0) {
                    //this.$message.warning('上个问题未答完')
                    return
                }
                let res = await deleteConversation({conversationId: n.conversationId})
                if (res.code === 0) {
                    this.getConversationList()
                    if(this.conversationId === n.conversationId){
                        this.conversationId = ''
                        this.$refs['session-com'].clearData()
                    }
                    this.echo = true
                }
            },
            /*------会话------*/
            async preSend(val,fileId,fileInfo) {
                this.inputVal = val || this.$refs['editable'].getPrompt()
                if (!this.inputVal) {
                    this.$message.warning('请输入内容');
                    return
                }
                if (!this.verifiyFormParams()) {
                    return;
                }
                // this.setParams()
                this.setSseParams({ragId: this.editForm.appId, question: this.inputVal})
                this.doragSend()
                this.echo = false
            },
            verifiyFormParams(){
                if (this.chatType === 'chat') return true;
                const { matchType, priorityMatch, rerankModelId } = this.editForm.knowledgeBaseConfig.config;
                const qArerankModelId  = this.editForm.qaKnowledgeBaseConfig.config.rerankModelId;
                const isMixPriorityMatch = matchType === 'mix' && priorityMatch;
                const conditions = [
                    { check: !this.editForm.modelParams, message: this.$t('knowledgeManage.create.selectModel') },
                    { check: !isMixPriorityMatch && !rerankModelId, message:  this.$t('knowledgeManage.hitTest.selectRerankModel')},
                    { check: this.editForm.qaKnowledgeBaseConfig.knowledgebases.length === 0 && this.editForm.knowledgeBaseConfig.knowledgebases.length === 0, message: this.$t('app.selectKnowledge') },
                    { check: this.editForm.qaKnowledgeBaseConfig.knowledgebases.length > 0 && !qArerankModelId,
                      message: this.$t('knowledgeManage.hitTest.selectRerankModel')
                    },
                ];
                for (const condition of conditions) {
                    if (condition.check) {
                    this.$message.warning(condition.message);
                    return false;
                    }
                }
                return true;
            },
            modelChange(){//切换模型新建对话
                this.preCreateConversation()
            },
            setParams() {
                ++this.amswerNum;
                if(this.amswerNum > 0){
                    this.isModelDisable = true
                }
                let fileId = this.getFileIdList() || this.fileId;
                this.useSearch = this.$refs['editable'].sendUseSearch()
                this.modelParams = this.$refs['editable'].getModelInfo()
                this.isBigModel = true;
                this.setSseParams({conversationId: this.conversationId, fileId})
                this.doSend()
                this.echo = false
            },
            /*--右侧提示词--*/
            showDrawer() {
                this.drawer = true
            },
            hideDrawer() {
                this.drawer = false
            },
            async getReminderList(cb) {
                let res = await getTemplateList({pageNo:0,pageSize:0,title:''})
                if (res.code === 0) {
                    this.reminderList = res.data.list||[]
                    cb && cb()
                    console.log(new Date().getTime())
                }
            },
            reminderClick(n) {
                this.$refs['editable'].setPrompt(n.prompt)
            },
            async doDeleteHistory(){
                let res = await deleteConversationHistory({conversationId:this.conversationId})
                if(res.code === 0){
                    this.$message.success(this.$t('yuanjing.deleteTips'))
                }
            },
        }
    }
</script>

<style lang="scss" scoped>
@import '@/style/chat.scss';
</style>
