<template>
  <el-dialog
    :visible.sync="dialogVisible"
    title="文生图名字"
    width="600px"
    :before-close="handleClose"
    class="text-to-image-dialog"
  >
    <div class="dialog-content">
      <!-- 标题和描述 -->
      <div class="header-section">
        <h2 class="dialog-title">文生图名字</h2>
        <p class="dialog-subtitle">MaaS-简要描述</p>
      </div>

      <!-- API Key 部分 -->
      <div class="api-key-section">
        <div class="api-key-label">API Key</div>
        <div class="api-key-input-group">
          <el-input
            v-model="apiKey"
            :placeholder="apiKeyPlaceholder"
            class="api-key-input"
            :disabled="hasApiKey"
          />
          <div class="api-key-buttons">
            <el-button 
              type="primary" 
              size="small" 
              @click="handleConfirm"
              class="confirm-btn"
            >
              确认
            </el-button>
            <el-button 
              type="primary" 
              size="small" 
              @click="handleUpdate"
              class="update-btn"
            >
              更新
            </el-button>
          </div>
        </div>
      </div>

      <!-- 参数表格 -->
      <div class="parameters-section">
        <el-table
          :data="parametersData"
          border
          class="parameters-table"
        >
          <el-table-column prop="parameter" label="参数" width="120" />
          <el-table-column prop="type" label="类型" width="100" />
          <el-table-column prop="description" label="描述" />
          <el-table-column prop="required" label="是否必填" width="100" align="center" />
        </el-table>
      </div>
    </div>

    <!-- 底部确认按钮 -->
    <div slot="footer" class="dialog-footer">
      <el-button 
        type="danger" 
        @click="handleFinalConfirm"
        class="final-confirm-btn"
      >
        确认
      </el-button>
    </div>
  </el-dialog>
</template>

<script>
export default {
  name: 'TextToImageDialog',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    initialApiKey: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      dialogVisible: false,
      apiKey: '',
      parametersData: [
        {
          parameter: '网页链接',
          type: '字符串',
          description: '用于链接到网页',
          required: '是'
        }
      ]
    }
  },
  computed: {
    hasApiKey() {
      return this.apiKey && this.apiKey.trim() !== '';
    },
    apiKeyPlaceholder() {
      return this.hasApiKey 
        ? '若没有添加过API Key,则显示输入框;若添加过,直接展示....'
        : '请输入API Key';
    }
  },
  watch: {
    visible: {
      handler(newVal) {
        this.dialogVisible = newVal;
        if (newVal) {
          this.apiKey = this.initialApiKey;
        }
      },
      immediate: true
    },
    dialogVisible(newVal) {
      this.$emit('update:visible', newVal);
    }
  },
  methods: {
    handleClose() {
      this.dialogVisible = false;
      this.$emit('close');
    },
    handleConfirm() {
      if (!this.apiKey.trim()) {
        this.$message.warning('请输入API Key');
        return;
      }
      this.$message.success('API Key 确认成功');
      this.$emit('api-key-confirm', this.apiKey);
    },
    handleUpdate() {
      if (!this.apiKey.trim()) {
        this.$message.warning('请输入API Key');
        return;
      }
      this.$message.success('API Key 更新成功');
      this.$emit('api-key-update', this.apiKey);
    },
    handleFinalConfirm() {
      if (!this.apiKey.trim()) {
        this.$message.warning('请先输入API Key');
        return;
      }
      this.$message.success('配置确认成功');
      this.$emit('confirm', {
        apiKey: this.apiKey,
        parameters: this.parametersData
      });
      this.dialogVisible = false;
    }
  }
}
</script>

<style lang="scss" scoped>
.text-to-image-dialog {
  .dialog-content {
    padding: 20px 0;
  }

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
      
      ::v-deep .el-table__header {
        background-color: #f5f7fa;
      }
      
      ::v-deep .el-table__header th {
        background-color: #f5f7fa;
        color: #333;
        font-weight: 500;
      }
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
}

// 响应式设计
@media (max-width: 768px) {
  .text-to-image-dialog {
    ::v-deep .el-dialog {
      width: 90% !important;
      margin: 5vh auto !important;
    }
    
    .api-key-input-group {
      flex-direction: column;
      
      .api-key-buttons {
        flex-direction: row;
        justify-content: flex-end;
      }
    }
  }
}
</style>

