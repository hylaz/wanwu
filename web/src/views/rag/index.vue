<template>
  <CommonLayout :aside-title="asideTitle" :isButton="false" :showAside="false">
    <template #main-content>
      <div class="app-content">
        <Chat :editForm="editForm" :chatType="'chat'" />
      </div>
    </template>
  </CommonLayout>
</template>
<script>
import CommonLayout from '@/components/exploreContainer.vue';
import Chat from './components/chat.vue';
import { getApiKeyRoot } from '@/api/appspace';
import { getRagInfo } from '@/api/rag';
export default {
  components: { CommonLayout, Chat },
  data() {
    return {
      apiURL: '',
      editForm: {
        appId: '',
        avatar: {},
        name: '',
        desc: '',
      },
      asideTitle: '文本问答名称',
      historyList: [
        {
          appId: '122249e8-c986-4c02-a731-c4c338c0683a',
          conversationId: '39ecc738-eb39-4812-93bf-3280746082ca',
          createdAt: '2025-06-03 14:57:44',
          hover: false,
          title:
            '你是一个问答助手，主要任务是汇总参考信息回答用户问题。请根据参考信息中提供的上下文信息回答用户问题，注意仅用提供的上下文作答不要根据自己已经有的先验知识来回答问题。',
        },
      ],
    };
  },
  created() {
    if (this.$route.query.id) {
      this.editForm.appId = this.$route.query.id;
      this.getDetail();
    }
  },
  methods: {
    getDetail() {
      getRagInfo({ ragId: this.editForm.appId }).then(res => {
        if (res.code === 0) {
          this.editForm.avatar = res.data.avatar;
          this.editForm.name = res.data.name;
          this.editForm.desc = res.data.desc;
        }
      });
    },
    goBack() {
      this.$router.go(-1);
    },
    openApiDialog() {
      this.$refs.apiKeyDialog.showDialog();
    },
    apiKeyRootUrl() {
      const data = { appId: this.editForm.appId, appType: 'rag' };
      getApiKeyRoot(data).then(res => {
        if (res.code === 0) {
          this.apiURL = res.data || '';
        }
      });
    },
  },
};
</script>
<style lang="scss" scoped>
/deep/ {
  .apikeyBtn {
    padding: 11px 10px;
    border: 1px solid $btn_bg;
    color: $btn_bg;
    display: flex;
    align-items: center;
    img {
      height: 14px;
    }
  }
}
.app-content {
  width: 100%;
  height: 100%;
  position: relative;
  .app-header-api {
    width: 100%;
    padding: 10px;
    position: absolute;
    z-index: 999;
    top: 0;
    left: 0;
    border-bottom: 1px solid #eaeaea;
    display: flex;
    justify-content: space-between;
    align-content: center;
    .app_name {
      font-size: 18px;
      font-weight: bold;
      color: $color_title;
      display: flex;
      align-items: center;
      .goBack {
        font-weight: bold;
        font-size: 16px;
        cursor: pointer;
        margin-right: 15px;
        color: #333;
      }
    }
    .header-api-box {
      display: flex;
      .header-api-url {
        padding: 6px 10px;
        background: #fff;
        margin: 0 10px;
        border-radius: 6px;
        .root-url {
          background-color: #eceefe;
          color: $color;
          border: none;
        }
      }
    }
  }
}
</style>
