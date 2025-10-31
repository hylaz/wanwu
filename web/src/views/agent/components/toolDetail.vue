<template>
  <el-dialog
    :visible.sync="dialogVisible"
    width="600px"
    :before-close="handleClose"
  >
    <!-- 标题和描述 -->
    <template slot="title">
      <div class="custom-title">
        <div class="header-section">
          <h2 class="dialog-title">{{actionDetail.action && actionDetail.action.name}}</h2>
          <!-- <p class="dialog-subtitle">{{actionDetail.action && actionDetail.action.description}}</p> -->
        </div>
      </div>
    </template>

    <!-- API Key 部分 -->
    <div class="api-key-section">
      <div class="api-key-label">API Key</div>
      <div class="api-key-input-group">
        <el-input
          v-model="apiKey"
          placeholder="请输入apikey"
          class="api-key-input"
          showPassword
          style="width:300px;"
        />
        <div class="api-key-buttons">
          <el-button style="width: 100px" size="mini" type="primary"  @click="changeApiKey">
            {{actionDetail.apiKey ? '更新' : '确认'}}
          </el-button>
        </div>
      </div>
    </div>

    <div class="parameters-section">
      <el-table :data="parametersData" border class="parameters-table">
        <el-table-column prop="key" label="参数" width="120" />
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column prop="description" label="描述" />
        <el-table-column
          prop="required"
          label="是否必填"
          width="100"
          align="center"
          :formatter="(row, column, cellValue) => (cellValue ? '是' : '否')"
        />
      </el-table>
    </div>
  </el-dialog>
</template>

<script>
import {toolActionDetail} from "@/api/agent";
import { changeApiKey } from "@/api/mcp";
export default {
  data() {
    return {
      dialogVisible: false,
      actionDetail:{
        action: {
          name: '',
          description: ''
        },
        apiKey: ''
      },
      apiKey: '',
      parametersData: [],
      currentItem:null
    }
  },
  methods: {
    changeApiKey(){
      changeApiKey({
        apiKey: this.apiKey,
        toolSquareId: this.currentItem.toolId
      }).then((res) => {
        if (res.code === 0) {
          this.$message.success(this.$t('common.message.success'))
          this.getDeatil(this.currentItem)
        }
      })
    },
    handleClose() {
      this.dialogVisible = false;
    },
    showDiaglog(n){
      this.dialogVisible = true;
      this.currentItem = n
      this.getDeatil(n)
    },
    getDeatil(n){
      toolActionDetail({actionName:n.actionName,toolId:n.toolId,toolType:n.toolType}).then(res =>{
        if(res.code === 0){
          this.actionDetail = res.data || {}
          const base = { action: { name: '', description: '' }, apiKey: '' }
          const payload = res.data || {}
          this.actionDetail = Object.assign({}, base, payload, {
            action: Object.assign({}, base.action, (payload.action || {})),
            apiKey: typeof payload.apiKey === 'string' ? payload.apiKey : ''
          })
          this.apiKey = res.data.apiKey
          const {properties,required} = (this.actionDetail.action && this.actionDetail.action.inputSchema) ? this.actionDetail.action.inputSchema : { properties: {}, required: [] };
          this.parametersData = this.toParametersData(properties,required);
        }
      }).catch(() =>{})
    },
    toParametersData(schemaObject, requiredKeys){
      return Object.entries(schemaObject).map(([key, def]) => ({
        key,
        type: def.type,
        description: def.description || '',
        required: requiredKeys.includes(key)
      }));
    }
  }
}
</script>

<style lang="scss" scoped>
/deep/{
  .el-dialog__body {
    padding:10px 20px!important;
  }
  .el-dialog__header{
    padding:10px 20px !important; 
    height:45px !important;
    border-bottom:1px solid #dbdbdb;
  }
}

.header-section {
  .dialog-title {
    font-size: 20px;
    font-weight: bold;
    color: #333;
    margin: 0 0 8px 0;
  }
  
  .dialog-subtitle {
    font-size: 14px;
    color: #666;
    margin: 0;
  }
}

.api-key-section {
  margin-bottom:15px;
  display:flex;
  .api-key-label {
    font-size: 14px;
    font-weight: 500;
    color: #333;
    margin-right: 12px;
    white-space: nowrap;
  }
  
  .api-key-input-group {
    display: flex;
    align-items: center;
    gap: 12px;
    .api-key-input {
      flex: 1;
    }
    
    .api-key-buttons {
      display: flex;
      flex-direction: row;
      gap: 8px;
      .confirm-btn,
      .update-btn {
        width: 60px;
        height: 32px;
      }
    }
  }
}

.parameters-section {
  .parameters-table {
    width: 100%;
  }
}

.dialog-footer {
  text-align: center;
  padding: 20px 0 0 0;
  
  .final-confirm-btn {
    width: 120px;
    height: 40px;
    font-size: 16px;
  }
}
</style>