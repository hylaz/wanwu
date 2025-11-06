<template>
  <el-dialog
    :visible.sync="dialogVisible"
    width="60%"
    :before-close="handleClose"
    class="prompt-dialog"
  >
    <div slot="title" class="dialog-title">
      <span class="title-text">{{$t('tempSquare.prompt')}}</span>
      <el-input
        v-model="searchKeyword"
        :placeholder="$t('agent.promptTemplate.searchPlaceholder')"
        prefix-icon="el-icon-search"
        class="title-search-input"
        clearable
        @clear="handleSearchClear"
      ></el-input>
    </div>
    <div class="prompt-library-content">
      <div class="tab-buttons">
        <div 
          class="tab-button" 
          :class="{ active: activeTab === 'builtIn' }"
          @click="activeTab = 'builtIn'"
        >
          {{ $t('agent.promptTemplate.builtIn') }}
        </div>
        <div 
          class="tab-button" 
          :class="{ active: activeTab === 'custom' }"
          @click="activeTab = 'custom'"
        >
          {{ $t('agent.promptTemplate.custom') }}
        </div>
      </div>
      <div class="library-main">
        <div class="template-list">
          <div
            v-for="item in filteredTemplateList"
            :key="item.templateId || item.customPromptId"
            class="template-item"
            :class="{ active: isTemplateSelected(item) }"
            @click="selectTemplate(item)"
          >
            <div class="template-content">
              <div class="template-logo">
                <img :src="'/user/api' + item.avatar.path">
              </div>
              <div class="template-info">
                <div class="template-name">{{ item.name }}</div>
                <div class="template-desc">{{ item.desc }}</div>
              </div>
            </div>
            <div class="template-actions" @click.stop="handleInsertPrompt(item)">
              <el-button type="text" size="mini">{{ $t('agent.promptTemplate.insertPrompt') }}</el-button>
            </div>
          </div>
        </div>

        <div class="template-detail" v-if="selectedTemplate">
          <div class="detail-content markdown-body" v-html="formatTemplateContent(selectedTemplate.prompt)"></div>
        </div>
        <div class="template-detail empty" v-else>
          <div class="empty-text">{{ $t('agent.promptTemplate.selectTemplate')}}</div>
        </div>
      </div>
    </div>

  </el-dialog>
</template>

<script>
import { md } from "@/mixins/marksown-it.js";

export default {
  name: 'PromptDialog',
  inject: ['getPrompt'],
  data() {
    return {
      dialogVisible: false,
      searchKeyword: '',
      activeTab: 'builtIn',
      selectedTemplate: null,
      templateList: []
    }
  },
  watch:{
    activeTab(newVal, oldVal) {
      if (this.dialogVisible && newVal !== oldVal) {
        this.$emit('tabChange', newVal);
      }
    }
  },
  computed: {
    filteredTemplateList() {
      if (!this.searchKeyword.trim()) {
        return this.templateList;
      }
      const keyword = this.searchKeyword.toLowerCase().trim();
      return this.templateList.filter(item => {
        return item.name && item.name.toLowerCase().includes(keyword);
      });
    }
  },
  methods: {
    showDiglog(data,type) {
      this.dialogVisible = true;
      this.activeTab = type;
      if (data && data.length) {
        this.templateList = data;
      }
      this.searchKeyword = '';
    },
    handleClose() {
      this.dialogVisible = false;
      this.selectedTemplate = null;
      this.searchKeyword = '';
    },
    handleSearchClear() {
      this.searchKeyword = '';
    },
    isTemplateSelected(item) {
      if (!this.selectedTemplate || !item) {
        return false;
      }
      if (this.selectedTemplate.templateId && item.templateId) {
        return this.selectedTemplate.templateId === item.templateId;
      }
      if (this.selectedTemplate.customPromptId && item.customPromptId) {
        return this.selectedTemplate.customPromptId === item.customPromptId;
      }
      return false;
    },
    selectTemplate(template) {
      this.selectedTemplate = template;
    },
    handleInsertPrompt(item) {
      this.getPrompt(item.prompt)
      this.$message.success(this.$t('agent.promptTemplate.insertSuccess'));
      this.dialogVisible = false;
    },
    formatTemplateContent(content) {
      if (!content) return '';
      return md.render(content);
    }
  }
}
</script>

<style lang="scss" scoped>
@import "@/style/markdown.scss";

.prompt-dialog {
  /deep/ .el-dialog__body {
    padding: 5px 20px 20px 20px;
  }
  
  /deep/ .el-dialog__header {
    padding: 20px 20px 10px;
    display: flex;
    align-items: center;
  }
  
  /deep/ .el-dialog__headerbtn {
    position:unset!important;
  }
}

.dialog-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  
  .title-text {
    font-size: 18px;
    font-weight: 500;
    color: #303133;
  }
  
  .title-search-input {
    width: 300px;
    margin-right:10px;
    /deep/ .el-input__inner {
      height: 32px;
      line-height: 32px;
    }
  }
}

.prompt-library-content {
  .tab-buttons {
    display: flex;
    margin-bottom: 16px;
    
    .tab-button {
      padding:5px 10px;
      cursor: pointer;
      font-size: 14px;
      color: #606266;
      transition: all 0.3s;
      
      &:hover {
        color: $color;
      }
      
      &.active {
        color: $color;
        font-weight: 500;
      }
    }
  }
  
  .library-main {
    display: flex;
    gap: 16px;
    height: 45vh;
    margin-top: 16px;
    
    .template-list {
      width:30vw;
      flex-shrink: 0;
      border: 1px solid #EBEEF5;
      border-radius: 4px;
      overflow-y: auto;
      
      .template-item {
        padding: 12px 10px;
        border-bottom: 1px solid #EBEEF5;
        cursor: pointer;
        transition: all 0.3s;
        display: flex;
        align-items: center;
        justify-content: space-between;
        
        &:last-child {
          border-bottom: none;
        }
        
        &:hover {
          background-color: #F5F7FA;
          
          .template-actions {
            opacity: 1;
          }
        }
        
        &.active {
          background-color: #ECF5FF;
          border-left: 3px solid $color;
          
          .template-name {
            color: $color;
          }
        }
        
        .template-content {
          flex: 1;
          display: flex;
          align-items: center;
          min-width: 0;
          
          .template-logo {
            width: 40px;
            height: 40px;
            border-radius:50%;
            background:#eee;
            flex-shrink: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-right: 12px;
            
            img {
              width:100%;
              height:100%;
              border-radius:50%;
              object-fit: cover;
            }
          }
          
          .template-info {
            flex: 1;
            min-width: 0;
            
            .template-name {
              font-size: 14px;
              font-weight: 500;
              color: #303133;
              margin-bottom: 4px;
              white-space: nowrap;
              overflow: hidden;
              text-overflow: ellipsis;
            }
            
            .template-desc {
              font-size: 12px;
              color: #909399;
              line-height: 1.5;
              display: -webkit-box;
              -webkit-line-clamp: 2;
              -webkit-box-orient: vertical;
              overflow: hidden;
            }
          }
        }
        
        .template-actions {
          flex-shrink: 0;
          opacity: 0;
          transition: opacity 0.3s;
          
          /deep/ .el-button {
            padding: 4px 0;
            font-weight: bold;
          }
        }
      }
    }
    
    .template-detail {
      flex: 1;
      border: 1px solid #EBEEF5;
      border-radius: 4px;
      padding: 20px;
      overflow-y: auto;
      background-color: #FAFAFA;
      
      &.empty {
        display: flex;
        align-items: center;
        justify-content: center;
        
        .empty-text {
          font-size: 14px;
          color: #909399;
        }
      }
      
      .detail-content {
        font-size: 14px;
        line-height: 1.8;
        color: #303133;
        
        /deep/ h1, /deep/ h2, /deep/ h3 {
          margin-top: 16px;
          margin-bottom: 8px;
          
          &:first-child {
            margin-top: 0;
          }
        }
      }
    }
  }
}
</style>