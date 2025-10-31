<template>
    <div>
        <el-dialog
            title="选择知识库"
            :visible.sync="dialogVisible"
            width="40%"
            :before-close="handleClose">
            <div class="tool-typ">
                <el-input v-model="toolName" placeholder="输入知识库名称搜索" class="tool-input" suffix-icon="el-icon-search" @keyup.enter.native="searchTool" clearable></el-input>
            </div>
            <div class="toolContent">
                <div 
                    v-for="(item,i) in knowledgeData"
                    :key="item['knowledgeId']"
                    class="toolContent_item"
                >
                    <div class="knowledge-info">
                        <span class="knowledge-name">{{ item.name }}</span>
                        <div class="knowledge-meta">
                            <span class="meta-text">{{item.share ? '公开' : '私密'}}</span>
                            <span v-if="item.share" class="meta-text">{{item.orgName}}</span>
                        </div>
                    </div>
                    <el-button type="text" @click="openTool($event,item)" v-if="!item.checked">添加</el-button>
                    <el-button type="text" v-else >已添加</el-button>
                </div>
            </div>
            <span slot="footer" class="dialog-footer">
                <el-button @click="handleClose">取 消</el-button>
                <el-button type="primary" @click="submit">确 定</el-button>
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
/deep/{
    .el-dialog__body{
        padding:10px 20px;
    }
}
.createTool{
    padding:10px;
    cursor: pointer;
    .add{
        padding-right:5px;
    }
}
.createTool:hover{
    color: $color;
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
        border:1px solid #dbdbdb;
        border-radius:6px;
        margin-bottom:10px;
        cursor: pointer;
        display: flex;
        align-items:center;
        justify-content:space-between;
        .knowledge-info{
            display: flex;
            flex-direction: column;
            gap: 4px;
            .knowledge-name{
                font-size: 14px;
                font-weight: 500;
            }
            .knowledge-meta{
                display: flex;
                gap: 8px;
                .meta-text{
                    color: #384BF7;
                    font-size: 12px;
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