<template>
    <div>
      <el-dialog
        top="10vh"
        title="添加敏感词"
        :close-on-click-modal="false"
        :visible.sync="dialogVisible"
        width="50%"
        :before-close="handleClose"
        >
        <el-form
            :model="ruleForm"
            ref="ruleForm"
            label-width="140px"
            class="demo-ruleForm"
            :rules="rules"
            @submit.native.prevent
        >
            <el-form-item class="itemCenter">
              <el-radio-group v-model="ruleForm.importType">
                <el-radio-button :label="'single'">单条添加</el-radio-button>
                <el-radio-button :label="'file'">批量上传</el-radio-button>
              </el-radio-group>  
            </el-form-item>
            <el-form-item
            label="敏感词"
            prop="word"
            v-if="ruleForm.importType === 'single'"
            >
            <el-input
                v-model="ruleForm.word"
                placeholder="您可添加一个词"
            ></el-input>
            </el-form-item>
            <el-form-item
            label="敏感词类型"
            prop="sensitiveType"
            v-if="ruleForm.importType === 'single'"
            >
            <el-select v-model="ruleForm.sensitiveType" placeholder="请选择" style="width:100%;">
                <el-option
                v-for="item in sensitiveTypeOptions"
                :key="item.value"
                :label="item.name"
                :value="item.value">
                </el-option>
            </el-select>
            </el-form-item>
            <el-form-item
            label="批量上传"
            prop="fileName"
            v-if="ruleForm.importType === 'file'"
            >
            <el-upload
                class="upload-box"
                drag
                action=""
                :show-file-list="false"
                :auto-upload="false"
                accept=".xlsx"
                :file-list="fileList"
                :on-change="uploadOnChange"
              >
              <div>
                <div>
                    <img :src="require('@/assets/imgs/uploadImg.png')" class="upload-img" />
                    <p class="click-text">
                        将文件拖到此处，或
                        <span class="clickUpload">点击上传</span>
                        <a class="clickUpload template" :href="`/user/api/v1/static/docs/sensitive.xlsx`" download @click.stop>模版下载</a>
                    </p>
                </div>
              </div>
              </el-upload>
               <!-- 上传文件的列表 -->
                <div class="file-list" v-if="fileList.length > 0">
                    <transition name="el-zoom-in-top">
                    <ul class="document_lise">
                    <li
                        v-for="(file, index) in fileList"
                        :key="index"
                        class="document_lise_item"
                    >
                        <div style="padding:8px 0;" class="lise_item_box">
                        <span class="size">
                            <img :src="require('@/assets/imgs/fileicon.png')" />
                            {{ file.name }}
                            <span class="file-size">
                            {{ filterSize(file.size) }}
                            </span>
                            <el-progress 
                            :percentage="file.percentage" 
                            v-if="file.percentage !== 100"
                            :status="file.progressStatus"
                            max="100"
                            class="progress"
                            ></el-progress>
                        </span>
                        <span class="handleBtn">
                            <span>
                            <span v-if="file.percentage === 100">
                                <i class="el-icon-check check success" v-if="file.progressStatus === 'success'"></i>
                                <i class="el-icon-close close fail" v-else></i>
                            </span>
                            <i class="el-icon-loading" v-else-if="file.percentage !== 100 && index === fileIndex"></i>
                            </span>
                            <span style="margin-left:30px;">
                            <i class="el-icon-error error" @click="handleRemove(file,index)"></i>
                            </span>
                        </span>
                        </div>
                    </li>
                    </ul>
                </transition>
                </div>
            </el-form-item>
        </el-form>
        <span
            slot="footer"
            class="dialog-footer"
        >
            <el-button 
                @click="handleClose()">
                {{$t('common.confirm.cancel')}}
            </el-button>
            <el-button
                type="primary"
                @click="submitForm('ruleForm')"
            >{{$t('common.confirm.confirm')}}</el-button>
        </span>
        </el-dialog>
    </div>
</template>
<script>
import uploadChunk from "@/mixins/uploadChunk";
import { delfile } from "@/api/chunkFile";
import { uploadSensitiveWord } from "@/api/safety";
export default {
    mixins: [uploadChunk],
    data(){
        return{
            sensitiveTypeOptions:[
                {
                    value:'Political',
                    name:'涉政'
                },
                {
                    value:'Revile',
                    name:'辱骂'
                },
                {
                    value:'Pornography',
                    name:'涉黄'
                },
                {
                    value:'ViolentTerror',
                    name:'暴恐'
                },
                {
                    value:'Illegal',
                    name:'违禁'
                },
                {
                    value:'InformationSecurity',
                    name:'信息安全'
                },
                {
                    value:'Other',
                    name:'其他'
                }
            ],
            title:"新建词表",
            dialogVisible:false,
            ruleForm:{
                importType:'single',
                word:'',
                sensitiveType:'',
                fileName:'',
                tableId:''
            },
            fileList:[],
            rules: {
                word: [{ required: true, message:'请输入敏感词', trigger: "blur" }],
                sensitiveType:[{ required: true, message: '请输选择敏感词类型', trigger: "blur" }],
                fileName:[{ required: true, message: '请上传文件', trigger: "blur" }]
            }
        }
    },
    methods:{
        uploadOnChange(file, fileList){
            if (!fileList.length) return;
            this.fileList = fileList;
            if(this.fileList.length > 0){
                this.maxSizeBytes = 0;
                this.isExpire = true;
                this.startUpload();
            }
        },
        filterSize(size) {
            if (!size) return "";
            var num = 1024.0; //byte
            if (size < num) return size + "B";
            if (size < Math.pow(num, 2)) return (size / num).toFixed(2) + "KB"; //kb
            if (size < Math.pow(num, 3))
                return (size / Math.pow(num, 2)).toFixed(2) + "MB"; //M
            if (size < Math.pow(num, 4))
                return (size / Math.pow(num, 3)).toFixed(2) + "G"; //G
            return (size / Math.pow(num, 4)).toFixed(2) + "T"; //T
        },
        handleRemove(item,index){
            const data = {fileList:[this.resList[index]['name']],isExpired:true}
            delfile(data).then(res =>{
                if(res.code === 0){
                this.$message.success('删除成功')
                }
            })
            this.fileList = this.fileList.filter((files) => files.name !== item.name);
            if(this.fileList.length === 0){
                this.file = null
            }else{
                this.fileIndex--
            }
        },
        uploadFile(chunkFileName){
            this.ruleForm.fileName = chunkFileName;
        },
        handleClose(){
            this.dialogVisible = false;
            this.ruleForm.tableId = '';
            this.clearform()
        },
        clearform(){
            this.fileList = []
            this.$refs.ruleForm.resetFields()
            this.$refs.ruleForm.clearValidate()
        },
        submitForm(formName){
            this.$refs[formName].validate((valid) =>{
                if(valid){
                   uploadSensitiveWord(this.ruleForm).then(res =>{
                    if(res.code == 0){
                        this.$message.success('操作成功')
                        this.$emit('reload')
                        this.dialogVisible = false;
                    }
                   }).catch(err =>{
                     
                   })
                }else{
                    return false;
                }
            })
        },
        showDialog(tableId){
            this.dialogVisible = true;
            this.ruleForm.tableId = tableId;
            this.clearform();
        }
    }
}
</script>
<style lang="scss" scoped>
.itemCenter{
    display:flex;
    justify-content: center;
    /deep/.el-form-item__content{
        margin-left: 0 !important;
    }
}
.upload-box{
    .upload-img{
        width:56px;
        height:56px;
        margin-top: 10px;
    }
    .clickUpload,.template{
       color: $color;
       font-weight: bold;
    }
    .template{
        margin-left:10px;
    }
}
.file-list{
  padding: 20px 0;
  .document_lise_item{
    cursor: pointer;
    padding:0 10px;
    list-style: none;
    border-radius:4px;
    border:1px solid #7684fd;
    display:flex;
    align-items:center;
    margin-bottom:10px;
    .lise_item_box{
      width:100%;
      display:flex;
      align-items:center;
      justify-content:space-between;
      .size{
          display:flex;
          flex:1;
          align-items:center;
          .progress{
            width:200px;
            margin-left:30px;
          }
          img{
            width: 18px;
            height:18px;
            margin-bottom:-3px;
          }
          .file-size{
            margin-left:10px;
          }
      }
    }
  }
  .document_lise_item:hover{
    background: #ECEEFE;
  }
}
</style>