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
          <p class="dialog-subtitle">{{actionDetail.action && actionDetail.action.description}}</p>
        </div>
      </div>
    </template>

    <!-- API Key 部分 -->
    <div class="api-key-section">
      <div class="api-key-label">API Key</div>
      <div class="api-key-input-group">
        <el-input
          v-model="actionDetail.apiKey"
          placeholder="若没有添加过API Key,则显示输入框;若添加过,直接展示...."
          class="api-key-input"
        />
        <div class="api-key-buttons">
          <el-button type="primary" size="small" class="confirm-btn" v-if="!actionDetail.apiKey || !actionDetail.apiKey.length">
            确认
          </el-button>
          <el-button type="primary" size="small" class="update-btn" else>
            更新
          </el-button>
        </div>
      </div>
    </div>

    <div class="parameters-section">
      <el-table :data="parametersData" border class="parameters-table">
        <el-table-column prop="key" label="参数" width="120" />
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="required" label="是否必填" width="100" align="center" />
      </el-table>
    </div>

    <div slot="footer" class="dialog-footer">
      <el-button type="danger" class="final-confirm-btn">
        确认
      </el-button>
    </div>
  </el-dialog>
</template>

<script>
import {toolActionDetail} from "@/api/agent";
export default {
  data() {
    return {
      dialogVisible: false,
      actionDetail:{},
      apiKey: '',
      parametersData: []
    }
  },
  methods: {
    handleClose() {
      this.dialogVisible = false;
    },
    showDiaglog(n){
      this.dialogVisible = true;
      this.getDeatil(n)
    },
    getDeatil(n){
      toolActionDetail({actionName:n.actionName,toolId:n.toolId,toolType:n.toolType}).then(res =>{
        if(res.code === 0){
          this.actionDetail = res.data || {}
          const {properties,required} = res.data.action.inputSchema;
          this.parametersData = this.toParametersData(properties,required);
          console.log(this.parametersData)
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
.header-section {
  margin-bottom: 24px;
  
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
  margin-bottom: 24px;
  
  .api-key-label {
    font-size: 14px;
    font-weight: 500;
    color: #333;
    margin-bottom: 8px;
  }
  
  .api-key-input-group {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    
    .api-key-input {
      flex: 1;
    }
    
    .api-key-buttons {
      display: flex;
      flex-direction: column;
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