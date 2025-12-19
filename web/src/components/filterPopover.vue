<template>
  <div class="filter-popover-wrapper">
    <el-popover
      v-model="showPopover"
      placement="bottom-start"
      trigger="hover"
      popper-class="filter-popover-popper"
      @hide="onPopoverHide"
    >
      <!-- 弹出内容区域 -->
      <div class="popover-content">
        <div class="options">
          <div class="option-item">
            <el-checkbox
              :indeterminate="isIndeterminate"
              v-model="checkAll"
              @change="handleCheckAllChange"
            >
              {{ $t('common.button.all') }}
            </el-checkbox>
          </div>
          <el-checkbox-group v-model="selectedStatuses">
            <div
              v-for="(item, index) in options"
              :key="index"
              class="option-item"
            >
              <el-checkbox :label="item.value">
                {{ item.label }}
              </el-checkbox>
            </div>
          </el-checkbox-group>
        </div>

        <div class="actions">
          <el-button size="mini" type="text" @click="applyFilter">
            {{ $t('common.button.filter') }}
          </el-button>
          <el-button size="mini" type="text" @click="resetFilter">
            {{ $t('common.button.reset') }}
          </el-button>
        </div>
      </div>

      <!-- 触发按钮 -->
      <div slot="reference" class="trigger">
        <i class="el-icon-arrow-down"></i>
      </div>
    </el-popover>
  </div>
</template>

<script>
export default {
  name: 'FilterPopover',
  props: {
    options: {
      type: Array,
      required: true,
      default: () => [],
    },
  },
  data() {
    return {
      selectedStatuses: [],
      showPopover: false,
      checkAll: false,
      isIndeterminate: false,
    };
  },
  watch: {
    selectedStatuses: {
      handler(val) {
        const checkedCount = val.length;
        this.checkAll = checkedCount === this.options.length;
        this.isIndeterminate =
          checkedCount > 0 && checkedCount < this.options.length;
      },
      immediate: true,
    },
  },
  methods: {
    handleCheckAllChange(val) {
      this.selectedStatuses = val ? this.options.map(item => item.value) : [];
      this.isIndeterminate = false;
    },
    applyFilter() {
      this.$emit('applyFilter', this.selectedStatuses);
      this.showPopover = false;
    },
    resetFilter() {
      this.selectedStatuses = [];
      this.$emit('applyFilter', this.selectedStatuses);
      this.checkAll = false;
      this.isIndeterminate = false;
      this.showPopover = false;
    },
    onPopoverHide() {
      this.$nextTick(() => {
        const popover = document.querySelector('.filter-popover-popper');
        if (popover && popover.contains(document.activeElement)) {
          document.activeElement.blur();
        }
      });
    },
  },
};
</script>

<style lang="scss">
.filter-popover-wrapper {
  .trigger {
    display: flex;
    align-items: center;
    cursor: pointer;
    color: #666;
    font-size: 14px;
    border: none;
    background: transparent;
    user-select: none;

    &:hover {
      background-color: transparent;
    }
  }
}

.filter-popover-popper {
  padding: 0 !important;

  .options {
    padding: 12px 0;

    .option-item {
      padding: 8px 16px;
      display: flex;
      align-items: center;
      cursor: pointer;

      .el-checkbox__input {
        margin-right: 8px;
      }

      .el-checkbox__label {
        cursor: pointer;
        font-size: 14px;
        color: #333;

        .el-checkbox__inner {
          accent-color: #409eff;
        }
      }
    }
  }

  .actions {
    display: flex;
    justify-content: space-between;
    padding: 8px 16px;
    border-top: 1px solid #ebeef5;
  }
}
</style>
