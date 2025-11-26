<template>
  <div class="safety-config-item">
    <div class="safety-config-content">
      <div class="safety-config-left">
        <span class="safety-config-label">{{ $t("app.chiChat.title") }}</span>
        <el-tooltip
          class="item"
          effect="dark"
          :content="$t('app.chiChat.tips')"
          placement="right-start"
        >
          <span class="el-icon-question question-icon"></span>
        </el-tooltip>
      </div>
      <div class="safety-config-right">
        <el-switch v-model="isEnabled" @change="handleSwitchChange"></el-switch>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "ChiChat",
  props: {
    chiChatSwitch: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      isEnabled: false,
    };
  },
  watch: {
    chiChatSwitch: {
      handler(newVal) {
        this.isEnabled = newVal;
      },
      immediate: true,
    },
  },
  methods: {
    handleSwitchChange(value) {
      // 处理开关变化事件
      this.$emit("chiSwitchChange", value);
    },
  },
};
</script>

<style lang="scss" scoped>
.safety-config-item {
  width: 100%;

  .safety-config-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: #f7f8fa;
    border-radius: 8px;
    padding: 10px 15px;
    box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.15);

    .safety-config-left {
      display: flex;
      align-items: center;

      .safety-config-label {
        font-size: 15px;
        font-weight: bold;
      }

      .question-icon {
        font-size: 13px;
        line-height: 1;
        cursor: pointer;
        margin-left: 5px;
      }
    }

    .safety-config-right {
      display: flex;
      align-items: center;
      gap: 12px;

      .safety-config-status {
        font-size: 14px;
        color: #333333;
      }
    }
  }
}
</style>
