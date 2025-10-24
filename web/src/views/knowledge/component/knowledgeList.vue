<template>
  <div class="app-card-container">
    <div class="app-card">
      <div class="smart rl smart-create">
        <div class="app-card-create" @click="showCreate">
          <div class="create-img-wrap">
            <img class="create-type" src="@/assets/imgs/create_knowledge.svg" alt="" />
            <img class="create-img" src="@/assets/imgs/create_icon.png" alt="" />
            <div class="create-filter"></div>
          </div>
          <span>创建知识库</span>
        </div>
      </div>
      <template v-if="listData && listData.length">
        <div class="smart rl" 
        v-for="(n,i) in listData" 
        :key="`${i}sm`" 
        @click.stop="toDocList(n)">
          <div>
              <img  class="logo" :src="require('@/assets/imgs/knowledgeIcon.png')" />
              <p :class="['smartDate']">{{n.docCount || 0}}个文档</p>
          </div>
          <div class="info rl">
            <p class="name" :title="n.name">
              {{n.name}}
            </p>
            <el-tooltip
              v-if="n.description"
              popper-class="instr-tooltip tooltip-cover-arrow"
              effect="dark"
              :content="n.description"
              placement="bottom-start"
            >
              <p class="desc">{{n.description}}</p>
            </el-tooltip>
          </div>
          <div class="tags">
            <!-- <span :class="['smartDate','tagList']">{{n.docCount || 0}}个文档</span> -->
            <span :class="['smartDate','tagList']" v-if="formattedTagNames(n.knowledgeTagList).length === 0" @click.stop="addTag(n.knowledgeId)">
              <span class="el-icon-price-tag icon-tag"></span>
              添加标签
            </span>
            <span v-else @click.stop="addTag(n.knowledgeId)">{{formattedTagNames(n.knowledgeTagList) }}</span>
          </div>
          <div class="editor">
            <el-dropdown @command="handleClick($event, n)" placement="top">
              <span class="el-dropdown-link">
                <i class="el-icon-more icon edit-icon" @click.stop></i>
              </span>
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item command="edit">{{$t('common.button.edit')}}</el-dropdown-item>
                <el-dropdown-item command="delete">{{$t('common.button.delete')}}</el-dropdown-item>
                <el-dropdown-item command="power">权限</el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
          </div>
        </div>
      </template>
    </div>
    <el-empty class="noData" v-if="!(listData && listData.length)" :description="$t('common.noData')"></el-empty>
    <tagDialog ref="tagDialog" @relodaData="relodaData" type="knowledge" :title="title"/>
    <PowerManagement ref="powerManagement" />
  </div>
</template>

<script>
import { delKnowledgeItem } from "@/api/knowledge";
import { AppType } from "@/utils/commonSet"
import tagDialog from './tagDialog.vue';
import PowerManagement from './power/index.vue';
export default {
  components:{tagDialog, PowerManagement},
  props:{
    appData:{
      type:Array,
      required:true,
      default:[]
    }
  },
  watch:{
    appData:{
      handler:function(val){
        this.listData = val
      },
      immediate:true,
      deep:true
    }
  },
  data(){
    return{
      apptype:AppType,
      basePath: this.$basePath,
      listData:[],
      title:'创建标签'
    }
  },
  methods:{
  formattedTagNames(data){
    if(data.length === 0){
      return [];
    }
    const tags = data.filter(item => item.selected).map(item =>  item.tagName ).join(', ');
    if (tags.length > 30) {
        return tags.slice(0, 30) + '...';
    }
    return tags;
  },
  addTag(id){
    this.$refs.tagDialog.showDiaglog(id);
  },
  showCreate(){
    this.$parent.showCreate();
  },
    handleClick(command,n){
      switch (command){
        case 'edit':
          this.editItem(n);
          break;
        case 'delete':
          this.deleteItem(n.knowledgeId)
          break;
        case 'power':
          this.showPowerManagement(n);
          break;
      }
    },
    editItem(row) {
      this.$emit('editItem', row)
    },
    relodaData(){
      this.$emit('reloadData');
    },
    deleteItem(knowledgeId){
      this.$confirm(this.$t('knowledgeManage.delKnowledgeTips'), this.$t('knowledgeManage.tip'), {
        confirmButtonText: this.$t('common.confirm.confirm'),
        cancelButtonText: this.$t('common.confirm.cancel'),
        type: "warning",
        beforeClose:(action, instance, done) =>{
          if(action === 'confirm'){
            instance.confirmButtonLoading = true;
            delKnowledgeItem({knowledgeId})
              .then(res =>{
                if(res.code === 0){
                  this.$message.success(this.$t('knowledgeManage.operateSuccess'));
                  this.$emit('reloadData')
                }
              })
              .catch(() => {})
              .finally(() =>{
                done();
                setTimeout(() => {
                  instance.confirmButtonLoading = false;
                }, 300);
              })
          }else{
            done()
          }
        }
      }).then(() => {})
    },
    toDocList(n){
      this.$router.push({path:`/knowledge/doclist/${n.knowledgeId}`,query:{name:n.name,permission:n.permissionType}});
    },
    showPowerManagement(knowledgeItem) {
      this.$refs.powerManagement.knowledgeId = knowledgeItem.knowledgeId;
      this.$refs.powerManagement.knowledgeName = knowledgeItem.knowledgeName;
      this.$refs.powerManagement.showDialog();
    },
  }
}
</script>

<style lang="scss" scoped>
@import "@/style/appCard.scss";
.app-card {
  .smart {
    height: 152px;
    .smartDate{
      // text-align:center;
      padding-top:3px;
      color:#888888;
    }
    .info {
      padding-right: 0;
    }
    .desc{
      padding-top: 5px;
    }
    .logo{
      border-radius:50%;
      background: #F1F4FF;
      box-shadow: none;
      padding: 5px !important;
      width: 65px !important;
      height: 65px !important;
    }
    .tagList{
      cursor: pointer;
      .icon-tag{
        transform: rotate(-40deg);
        margin-right:3px;
      }
    }
    .tagList:hover{
        color:#384BF7;
    }
  }
}
</style>
