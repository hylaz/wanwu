<template>
    <div class="fileUpload">
        <el-dialog
                custom-class="upload-dialog"
                :visible.sync="dialogVisible"
                width="800px"
                append-to-body
                :before-close="handleClose">
                <div  v-loading="loading" element-loading-background="rgba(255, 255, 255, 0.5)">
                    <div class="dialog-body">
                        <p class="upload-title">{{$t('common.fileUpload.uploadFile')}}</p>
                        <el-upload
                                :class="['upload-box']"
                                drag
                                action=""
                                :show-file-list="false"
                                :auto-upload="false"
                                :limit="fileType === 'image/*' ? maxPicNum : 2"
                                :accept="tipsArr"
                                :file-list="fileList"
                                :on-change="uploadOnChange"
                                >
                                <div v-if="fileUrl" class="echo-img-box">
                                    <div class="echo-img">
                                        <!-- '/user/api'+fileUrl -->
                                        <video v-if="fileType === 'video/*'" id="video" muted loop playsinline>
                                            <source :src='fileUrl' type="video/mp4">
                                            {{$t('common.fileUpload.videoTips')}}
                                        </video>
                                        <audio v-if="fileType === 'audio/*'" id="audio" controls>
                                            <source :src="fileUrl" type="video/mp3">
                                            <source :src="fileUrl" type="audio/ogg">
                                            <source :src="fileUrl" type="audio/mpeg">
                                            {{$t('common.fileUpload.audioTips')}}
                                        </audio>
                                        <div v-if="fileType === 'doc/*'" class="docFile">
                                            <img :src="require('@/assets/imgs/fileicon.png')" />
                                        </div>
                                        <div v-if="fileType === 'image/*'" class="type-img-container">
                                            <el-button v-show="canScroll" icon="el-icon-arrow-left " @click="prev($event)" circle class="scroll-btn left" size="mini" type="primary"></el-button>
                                            <div class="type-img" ref="imgList" :style="{justifyContent: !canScroll ? 'center':'unset'}">
                                                <div v-for="(f, idx) in fileList" :key="f.uid || idx"  class="type-img-item">
                                                    <img :src="f.fileUrl" />
                                                    <p class="type-img-info">
                                                        <el-tooltip class="item" effect="dark" :content="f.name" placement="top-start">
                                                            <span>{{f.name.length>6?f.name.slice(0,6)+'...':f.name}}</span>
                                                        </el-tooltip>
                                                        <span>[ {{f.size > 1024 ? (f.size / (1024 * 1024 )).toFixed(2) + ' MB' : f.size + ' bytes' }} ]</span>
                                                    </p>
                                                </div>
                                            </div>
                                            <el-button v-show="canScroll" icon="el-icon-arrow-right" @click="next($event)" circle class="scroll-btn right" size="mini" type="primary"></el-button>
                                        </div>
                                        <div v-else>
                                            <p>文件名称: {{fileList[0]['name']}}</p>
                                            <p>文件大小: {{fileList[0]['size'] > 1024 ? (fileList[0]['size'] / (1024 * 1024 )).toFixed(2) + ' MB' : fileList[0]['size'] + ' bytes' }}</p>
                                        </div>
                                    </div>
                                    <!--<i  class="el-icon-close" @click.stop="clearFile"></i>-->
                                    <div class="tips">
                                        <el-progress
                                            :percentage="file.percentage"
                                            v-if="file.percentage !== 100"
                                            :status="file.progressStatus"
                                            max="100"
                                            style="width:360px;margin:0 auto;"
                                        ></el-progress>
                                        <p>图片类型限制{{maxPicNum}}个文件，其它类型限制1个文件<span style="color:var(--color);"> {{$t('common.fileUpload.click')}} </span>非图片类型文件会替换已有文件</p>
                                    </div>
                                </div>
                                <div v-else>
                                    <i class="el-icon-upload"></i><p>{{ $t('common.fileUpload.uploadText') + $t('common.fileUpload.uploadClick') }}</p>
                                    <div class="tips">
                                        <p>{{$t('common.fileUpload.typeFileTip1')}}
                                            <span>{{tipsArr}}</span>
                                            {{$t('common.fileUpload.typeFileTip')}}
                                        </p>
                                        <p style="padding-top: 5px;color:#dc6803!important;">*若该智能体基于大语言模型创建，则上传图片暂时无法进行解析</p>
                                    </div>
                                </div>
                        </el-upload>
                    </div>
                    <div class="dialog-footer">
                        <el-button type="primary" :disabled="!fileUrl || !(file && file.percentage === 100 )" @click="doBatchUpload">{{$t('common.fileUpload.submitBtn')}}</el-button>
                    </div>
                </div>
        </el-dialog>

    </div>
</template>

<script>
    import { mapGetters } from "vuex";
    import { batchUpload,confirmPath } from '@/api/chat';
    import uploadChunk from "@/mixins/uploadChunk";
    export default {
        props:['fileTypeArr','sessionId'],
        mixins: [uploadChunk],
        data(){
            return{
                canScroll:false,
                fileIdList:[],
                fileList:[],
                fileType:'',
                //上传文件弹框
                loading:false,
                dialogVisible:false,
                fileUrl:'',
                //上传文件
                imgConfig:["jpeg", "PNG", "png", "JPG", "jpg",'bmp','webp'],
                audioConfig:['mp3','wav'],
                tipsArr:'',
                tipsObj:{
                    'image/*':['.jpg', '.jpeg', '.png'],//'.webp'
                    'audio/*':['.wav', '.mp3'],
                    'doc/*':['.txt','.csv','.xlsx','.docx','.html','.pptx','.pdf']
                },
                chunkFileName:'',
                fileInfo:[],
                lastFileType:'',
                imgUrl:''
            }
        },
        watch:{
            fileTypeArr:{
                handler(val,oldVal) {
                    this.setFileType(val)
                },
                immediate:true
            }
        },
        computed: {
            ...mapGetters("app", ["maxPicNum"]),
        },
        created(){
            this.sessionId = this.sessionId || this.$route.query.sessionId
        },
        methods:{
            checkScrollable() {
                this.$nextTick(() => {
                    const container = this.$refs.imgList
                    if (container) {
                        this.canScroll = container.scrollWidth > container.clientWidth
                    }
                })
            },
            prev(e){
                e.stopPropagation()
                this.$refs.imgList.scrollBy({
                    left: -200,
                    behavior: "smooth",
                });
            },
            next(e){
                e.stopPropagation()
                this.$refs.imgList.scrollBy({
                    left: 200,
                    behavior: "smooth",
                });
            },
            setFileType(fileTypeArr){
                if(fileTypeArr.length){
                    this.tipsArr = ''
                    let tips_arr=[]
                    fileTypeArr.forEach(item=>{
                        tips_arr = tips_arr.concat(this.tipsObj[item])
                    })
                    this.tipsArr = tips_arr.join(', ')
                }
            },
            openDialog(){
                this.dialogVisible = true
            },
            clearFile(){
                this.fileIdList = []
                this.fileList = []
                this.fileType = ''
                this.fileUrl = ''
                this.imgUrl = ''
                this.fileInfo = []
                this.canScroll = false
            },
            handleClose(){
                this.clearFile()
                this.dialogVisible = false
            },
            uploadOnChange(file, fileList) {
                const prevFileType = this.fileType; // 保存上一次的文件类型
                let filename= file.name
                //通过上传的文件名判断文件类型，用于回显
                let fileType = filename.split('.')[filename.split('.').length-1]
                // 重置图片URL
                this.imgUrl = '';
                
                if(["jpeg", "PNG", "png", "JPG", "jpg",'bmp','webp'].includes(fileType)){
                    this.fileType = 'image/*'
                    // 获取图片预览URL
                    if (file.url) {
                        this.imgUrl = file.url;
                    } else if (file.raw) {
                        this.imgUrl = URL.createObjectURL(file.raw);
                    }
                }
                if(["mp3", "wav"].includes(fileType)){
                    this.fileType = 'audio/*'
                }
                if(['txt','csv','xlsx','docx','html','pptx','pdf'].includes(fileType)){
                    this.fileType = 'doc/*'
                }
                
                // 创建文件预览URL
                this.fileUrl = URL.createObjectURL(file.raw);
                
                if (this.fileType === 'image/*') {
                    // 图片类型可累加至6个
                    if (fileList.length > 6) {
                        this.$message.warning('只能上传6个图片文件');
                        return;
                    }
                    if (prevFileType && prevFileType !== this.fileType) {
                        this.fileList = [];
                        this.canScroll = false;
                        this.fileList.push(file);
                    }else{
                        this.fileList = fileList;
                    }
                    const currentFileIndex = this.fileList.length - 1; // 当前文件在列表中的索引
                    if (file.raw) {
                        this.fileList[currentFileIndex].fileUrl = URL.createObjectURL(file.raw);
                    }
                    this.checkScrollable();
                } else {
                    // 非图片类型只保留最新一个
                    this.fileList = [];
                    this.fileList.push(file);
                }
                
                if(this.fileList.length > 0){
                    this.maxSizeBytes = 0;
                    this.isExpire = true;
                    //this.startUpload();
                    // 为每个文件启动上传，而不是只上传索引0的文件
                    for(let i = 0; i < this.fileList.length; i++) {
                        if (!this.fileList[i].uploaded) { // 添加标记避免重复上传
                            this.startUpload(i);
                            this.fileList[i].uploaded = true;
                        }
                    }
                }
            },
            uploadFile(fileName,oldFileName,fiePath){//文件上传完之后
                if (this.lastFileType && this.lastFileType !== this.fileType) {
                    this.fileInfo = [];
                }
                this.lastFileType = this.fileType;

                this.fileInfo.push({
                    fileName,
                    fileSize:this.fileList[this.fileIndex]['size'],
                    fileUrl:fiePath,
                })
            },
            doBatchUpload(){
                this.$emit('setFileId',this.fileInfo)
                this.$emit('setFile',this.fileList)
                this.clearFile()
                this.handleClose()
            },
            getFileIdList(){
                return this.fileIdList
            },
        }
    }
</script>

<style lang="scss" scoped>
    .upload-dialog {
        .dialog-body{
            padding:0 20px;
            .upload-title{
                text-align: center;
                font-size: 18px;
                margin-bottom: 20px;
            }
            .upload-box{
                height: 190px;
                width: 100% !important;
                background-color: #fff;
                .el-upload-dragger {
                    .el-icon-upload{
                        margin: 46px 0 10px 0 !important;
                        font-size: 32px !important;
                        line-height: 36px!important;
                        color: $color;
                    }
                    .el-upload__text{
                        margin-top: -10px;
                    }
                }
            }

            .echo-img-box{
                background-color: transparent!important;
                .echo-img{
                    .type-img-container{
                        width:100%;
                        position:relative;
                        .scroll-btn{
                            position:absolute;
                            top:50%;
                            transform: translateY(-32px);
                            &.left{
                                left:5px;
                            }
                            &.right{
                                right:5px;
                            }
                        }
                    .type-img{
                        display: flex;
                        gap: 10px;
                        width:100%;
                        overflow-x: hidden;
                        scroll-behavior: smooth;
                        .type-img-item{
                            width: auto !important;
                            flex-shrink: 0;
                            margin-bottom: 10px;
                        }
                        .type-img-info{
                            display: flex;
                            gap: 5px;
                            justify-content: center;
                            span{
                                color: $color;
                            }
                        }
                    }
                    }
                    img,video{
                        width: auto;
                        height: 80px;
                        margin: 10px auto;
                        border-radius: 4px;
                        background-color: transparent;
                    }
                    audio{
                        width: 300px;
                        height: 54px;
                        margin: 50px auto;
                    }
                }
                .docFile{
                    img{
                        margin:0;
                        width:60px;
                        height: 100px;
                    }
                }
            }
            .tips{
                position: absolute;
                bottom: 16px;
                left: 0;
                right: 0;
                p{
                    color: #9d8d8d!important;
                }
            }
        }
        .dialog-footer{
            text-align: center;
            margin: 30px 0 20px 0;
        }
    }
</style>
