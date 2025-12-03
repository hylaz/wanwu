<template>
  <div>
    <el-dialog
      :title="$t('app.recallParameterSet')"
      :visible.sync="dialogVisible"
      width="50%"
      :before-close="handleClose">
      <span v-if="dialogVisible">
        <searchConfig 
          ref='searchConfig' 
          @sendConfigInfo="sendConfigInfo" 
          :setType="'agent'" 
          :config="config" 
          :showGraphSwitch="showGraphSwitch"
          :category="category"
        />
      </span>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">{{$t('common.button.cancel')}}</el-button>
        <el-button type="primary" @click="submit">{{$t('common.button.confirm')}}</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
import searchConfig from '@/components/searchConfig.vue';
export default {
  props:['showGraphSwitch','config','category'],
  components:{
    searchConfig
  },
  data(){
    return{
      dialogVisible:false,
      knowledgeConfig:{},
    }
  },
  watch:{
    config:{
      handler(val){
        if(val){
          this.knowledgeConfig = {...val}
        }
      },
      deep:true
    }
  },
  methods:{
    sendConfigInfo(data){
      this.knowledgeConfig = { ...data.knowledgeMatchParams };
    },
    showDialog(){
      this.dialogVisible = true;
    },
    handleClose(){
      this.dialogVisible = false;
    },
    submit(){
      // 验证模型选择
      const { matchType, priorityMatch, rerankModelId } = this.knowledgeConfig;
      const needRerankModel = matchType === 'vector' || 
                               matchType === 'text' || 
                               (matchType === 'mix' && priorityMatch === 0);
      
      if (needRerankModel && !rerankModelId) {
        this.$message.error(this.$t('agent.form.selectModel'));
        return;
      }
      
      if(matchType === 'mix' && priorityMatch === 1){
        this.knowledgeConfig.rerankModelId = '';
      }
      this.dialogVisible = false;
      this.$emit('setKnowledgeSet',this.knowledgeConfig)
    }
  }
}
</script>
<style lang="scss" scoped>
/deep/{
  .el-input-number--small{
    line-height: 28px!important;
  }
}
.question{
  cursor: pointer;
  color:#ccc;
}
</style>
