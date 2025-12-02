<template>
  <div class="batch-operation-toolbar" v-if="selectedCount > 0 && type === 'multiple'">
    <div class="toolbar-container">
      <div class="toolbar-item selected-item">
        <span class="selected-badge">{{ selectedCount }}</span>
        <span class="selected-text">{{ $t('metaData.selected') }}</span>
      </div>
      <div class="toolbar-divider"></div>
      <div class="toolbar-item" @click="handleMetadata">
        <div class="icon-wrapper">
          <i class="el-icon-document metadata-icon"></i>
          <i class="el-icon-price-tag tag-icon"></i>
        </div>
        <span class="item-text">{{ $t('metaData.batchEdit') }}</span>
      </div>
      <div class="toolbar-divider"></div>
      <div class="toolbar-item" @click="handleBatchDelete">
        <div class="icon-wrapper">
          <i class="el-icon-delete metadata-icon"></i>
        </div>
        <span class="item-text">{{ $t('metaData.batchDelete') }}</span>
      </div>
      <div class="toolbar-divider"></div>
      <div class="toolbar-item" @click="handleCancel">
        <span class="item-text">{{ $t('common.confirm.cancel') }}</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'BatchMetatButton',
  props: {
    selectedCount: {
      type: Number,
      default: 0
    },
    type:{
      type:String,
      default:'multiple'
    }
  },
  methods: {
    handleMetadata() {
      this.$emit('showBatchMeta');
    },
    handleBatchDelete() {
      this.$emit('handleBatchDelete');
    },
    handleCancel() {
      this.$emit('handleMetaCancel');
    }
  }
};
</script>

<style lang="scss" scoped>
.batch-operation-toolbar {
  position:absolute;
  bottom:20px;
  left:50%;
  transform: translateX(-50%);
  z-index:1000;
  padding: 12px 20px;
  margin-bottom: 16px;
  border-radius: 8px;
  
  .toolbar-container {
    display: flex;
    align-items: center;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid $color; 
    border-radius: 8px;
    padding: 6px 26px;
    flex-wrap: wrap;
    
    .toolbar-item {
      display: flex;
      align-items: center;
      cursor: pointer;
      padding: 6px 12px;
      transition: all 0.2s;
      border-radius: 4px;
      
      &:hover:not(.selected-item) {
        background: rgba(255, 255, 255, 0.1);
      }
      
      &.selected-item {
        cursor: default;
        
        .selected-badge {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: 20px;
          height: 20px;
          background: $color;
          color: #fff;
          border-radius: 4px;
          font-size: 12px;
          font-weight: bold;
          margin-right: 8px;
        }
        
        .selected-text {
          color: $color;
          font-size: 14px;
          font-weight: 500;
        }
      }
      
      .icon-wrapper {
        position: relative;
        width: 20px;
        height: 20px;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        margin-right: 8px;
        
        i {
          font-size: 16px;
          color: $color;
          
          &.metadata-icon {
            font-size: 18px;
            color: $color;
          }
          
          &.tag-icon {
            position: absolute;
            top: -2px;
            right: -4px;
            font-size: 10px;
            color: $color;
          }
        }
      }
      
      .item-text {
        color: $color;
        font-size: 14px;
        white-space: nowrap;
      }
    }
    
    .toolbar-divider {
      width: 1px;
      height: 14px;
      background: $color;
      margin: 0 8px;
    }
  }
}
</style>