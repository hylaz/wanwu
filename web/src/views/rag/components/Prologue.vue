<!--开场白-->
<template>
  <div class="history-box">
    <div class="session-answer">
      <div :class="['session-item', 'rl']">
        <img
          class="logo"
          :src="
            editForm.avatar.path
              ? `/user/api` + editForm.avatar.path
              : '@/assets/imgs/bg-logo.png'
          "
        />
        <div class="answer-content">
          <p class="name">{{ editForm.name }}</p>
          <p class="systemPrompt">{{ editForm.desc }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  props: ['editForm'],
  data() {
    return {
      basePath: this.$basePath,
    };
  },
  methods: {
    setProloguePrompt(val) {
      this.$emit('setProloguePrompt', val);
    },
  },
};
</script>
<style lang="scss" scoped>
.history-box {
  height: calc(100% - 46px);
  overflow-y: auto;
  display: flex;
  align-items: center;
  justify-content: center;
}
.echo {
  word-break: break-all;
  height: 100%;
  overflow-y: auto;
  padding: 0 20px;
  width: auto;
  .session-item {
    min-height: 80px;
    display: flex;
    justify-content: center;
    flex-wrap: wrap;
    padding: 20px;
    line-height: 28px;
    img {
      width: 100px;
      height: 100px;
      object-fit: cover;
      border-radius: 50%;
      border: 1px solid #eee;
    }
    .logo {
      border-radius: 50%;
    }
    .answer-content {
      width: calc(100% - 30px);
      position: relative;
      margin-left: 14px;
      color: #333;
      margin: 0 auto;
      .name {
        font-size: 18px;
        font-weight: bold;
        text-align: center;
        padding: 12px 0;
      }
      .systemPrompt {
        line-height: 26px;
        text-align: center;
        font-size: 16px;
        font-weight: bold;
      }
      .recommend {
        color: #2c7eea;
        text-align: center;
        cursor: pointer;
        .recommend-p {
          margin: 10px 0;
          margin-left: -10px;
          padding: 0 10px;
        }
        .recommend-p:hover {
          background-color: #fafafa;
          border-radius: 6px;
        }
      }
    }
  }
}
</style>
