<template>
    <div>
        <el-dialog
            :visible.sync="dialogVisible"
            width="40%"
            :before-close="handleClose">
            <template slot="title">
                <div class="diglog_title">
                  <h3>{{$t('knowledgeSelect.title')}}</h3>
                   <el-input v-model="toolName" :placeholder="$t('knowledgeSelect.searchPlaceholder')" class="tool-input" suffix-icon="el-icon-search" @keyup.enter.native="searchTool" clearable></el-input>
                </div>
            </template>
            <div class="toolContent">
                <div 
                    v-for="(item,i) in knowledgeData"
                    :key="item['knowledgeId']"
                    class="toolContent_item"
                >
                    <div class="knowledge-info">
                        <span class="knowledge-name">{{ item.name }}</span>
                        <span class="knowledge-desc">{{item.description}}</span>
                        <div class="knowledge-meta">
                            <span class="meta-text">{{item.share ? $t('knowledgeSelect.public') : $t('knowledgeSelect.private')}}</span>
                            <span v-if="item.share" class="meta-text">{{item.orgName}}</span>
                        </div>
                        <span class="knowledge-createAt">{{ $t('knowledgeSelect.createTime')}} {{item.createAt}}</span>
                    </div>
                    <el-button type="primary" @click="openTool($event,item)" v-if="!item.checked" size="small">{{ $t('knowledgeSelect.add')}}</el-button>
                    <el-button type="primary" v-else  size="small">{{ $t('knowledgeSelect.added')}}</el-button>
                </div>
            </div>
            <span slot="footer" class="dialog-footer">
                <el-button @click="handleClose">{{$t('common.button.cancel')}}</el-button>
                <el-button type="primary" @click="submit">{{$t('common.button.confirm')}}</el-button>
            </span>
        </el-dialog>
    </div>
</template>
<script>
import { getKnowledgeList } from "@/api/knowledge";
export default {
    data(){
        return {
            dialogVisible:false,
            knowledgeData:[],
            knowledgeList:[],
            checkedData:[],
            toolName:''
        }
    },
    created(){
        this.getKnowledgeList('');
    },
    methods:{
        getKnowledgeList(name) {
            getKnowledgeList({name}).then((res) => {
                if (res.code === 0) {
                this.knowledgeData = (res.data.knowledgeList || []).map(m => ({
                    ...m,
                    checked:this.knowledgeList.some(item => item.id === m.knowledgeId)
                }));
                }
            }).catch(() =>{});
        },
        openTool(e,item){
            if(!e) return;
            item.checked = !item.checked
        },
        searchTool(){
            this.getKnowledgeList(this.toolName);
        },
        showDialog(data){
            this.dialogVisible = true;
            this.setKnowledge(data || []);
            this.knowledgeList = data || [];
        },
        setKnowledge(data){
           this.knowledgeData = this.knowledgeData.map(m => ({
            ...m,
            checked: data.some(item => item.id === m.knowledgeId)
            }));
        },
        handleClose(){
            this.dialogVisible = false;
        },
        submit(){
            const data = this.knowledgeData.filter(item => item.checked).map(item =>({
                id:item.knowledgeId,
                name:item.name
            }));
            this.$emit('getKnowledgeData',data);
            this.dialogVisible = false;
        }
    }
}
</script>
<style lang="scss" scoped>
@import "@/style/markdown.scss";
/deep/{
    .el-dialog__body{
        padding:10px 20px;
    }
    .el-dialog__header{
        display:flex;
        align-items:center;
        .el-dialog__headerbtn{
            top:unset!important;
        }
    }
}
.diglog_title{
    display:flex;
    justify-content:space-between;
    align-items:center;
    flex:1;
    h3{
        font-size:16px;
        font-weight:bold;
    }
    .tool-input{
        width:250px;
        margin-right:28px;
    }
}
.tool-typ{
    display:flex;
    justify-content:space-between;
    padding:10px 0;
    border-bottom: 1px solid #dbdbdb;
    .toolbtn{
        display:flex;
        justify-content:flex-start;
        gap:20px;
        div{
            text-align: center;
            padding:5px 20px;
            border-radius:6px;
            border:1px solid #ddd;
            cursor: pointer;
        }
    }
    .tool-input{
        width:200px;
    }
}
.toolContent{
    padding:10px 0;
    max-height:300px;
    overflow-y:auto;
    .toolContent_item{
        padding:5px 20px;
        border-bottom:1px solid $color_opacity;
        border-radius:6px;
        margin-bottom:10px;
        cursor: pointer;
        display: flex;
        align-items:center;
        justify-content:space-between;
        /deep/{
            .el-button--primary{
                background:#fff!important;
                border:1px solid #eee!important;
                padding:8px 16px;
                border-radius:6px;
                span{
                    color:$color !important;
                    font-size:14px;
                }
            }
        }
        .knowledge-info{
            display: flex;
            flex-direction: column;
            gap: 4px;
            .knowledge-name{
                font-size: 14px;
                font-weight: 600;
                color:#1c1d23;
            }
            .knowledge-desc,.knowledge-createAt{
                font-size:12px;
            }
            .knowledge-desc{
                color:rgba(28,29,35,.8);
            }
            .knowledge-createAt{
                color:rgba(28,29,35,.35);
                margin-top:5px;
            }
            .knowledge-meta{
                display: flex;
                gap: 8px;
                margin-top:5px;
                span{
                    padding:2px 8px;
                    background:rgba(139,139,149,.15);
                    color:#4b4a58;
                    font-size: 12px;
                    border-radius:6px;
                }
            }
        }
    }
    .toolContent_item:hover{
        background:$color_opacity;
    }
}
.active{
    border:1px solid $color !important;
    color: #fff;
    background:$color;
}
</style>