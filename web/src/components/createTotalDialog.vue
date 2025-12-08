<template>
  <div>
    <el-dialog
      :title="$t('createTotal.title')"
      :visible.sync="dialogVisible"
      width="460px"
      append-to-body
      :close-on-click-modal="false"
      :before-close="handleClose"
      class="createTotalDialog"
    >
      <div>
        <div class="create-card-total" v-if="list && list.length">
          <div
            class="create-card-item"
            @click="showCreate(item.key)"
            v-for="(item, index) in list"
            :key="item.name + index"
          >
            <img
              class="create-card-img"
              v-if="item.img"
              :src="item.img"
              alt=""
            />
            <div class="create-card-right">
              <div class="create-card-name">{{ item.name }}</div>
              <div class="create-card-desc">{{ item.desc }}</div>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
    <CreateIntelligent ref="createIntelligentDialog" />
    <CreateTxtQues ref="createTxtQuesDialog" />
    <CreateWorkflow ref="createWorkflowDialog" />
    <CreateChat ref="createChatDialog" />
  </div>
</template>

<script>
import CreateWorkflow from './createApp/createWorkflow.vue';
import CreateChat from './createApp/createChat.vue';
import CreateIntelligent from './createApp/createIntelligent';
import CreateTxtQues from './createApp/createRag.vue';
import { CHAT, RAG, AGENT } from '@/utils/commonSet';

export default {
  components: { CreateWorkflow, CreateChat, CreateIntelligent, CreateTxtQues },
  data() {
    return {
      dialogVisible: false,
      list: this.$t('createTotal.list'),
    };
  },
  mounted() {},
  methods: {
    openDialog() {
      this.dialogVisible = true;
    },
    handleClose() {
      this.dialogVisible = false;
    },
    showCreate(key) {
      this.handleClose();
      this.$nextTick(() => {
        this.showCreateDialog(key);
      });
    },
    showCreateIntelligent() {
      // 显示创建智能体
      this.$refs.createIntelligentDialog.openDialog();
    },
    showCreateTxtQues() {
      // 显示创建文本问答
      this.$refs.createTxtQuesDialog.openDialog();
    },
    showCreateWorkflow() {
      // 显示创建工作流
      this.$refs.createWorkflowDialog.openDialog();
    },
    showCreateChat() {
      // 显示创建对话流
      this.$refs.createChatDialog.openDialog();
    },
    showCreateDialog(key) {
      switch (key) {
        case AGENT:
          this.showCreateIntelligent();
          break;
        case RAG:
          this.showCreateTxtQues();
          break;
        case CHAT:
          this.showCreateChat();
          break;
        default:
          this.showCreateWorkflow();
          break;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.create-card-total {
  margin-top: -15px;
  .create-card-item {
    margin-bottom: 10px;
    border-radius: 6px;
    padding: 10px;
    border: 1px solid #d9d9d9;
    display: flex;
    align-items: center;
    justify-content: space-between;
    cursor: pointer;
    .create-card-right {
      width: calc(100% - 60px);
      color: $color_title;
    }
    .create-card-img {
      width: 50px;
    }
    .create-card-name {
      font-size: 14px;
      font-weight: bold;
    }
    .create-card-desc {
      font-size: 13px;
      margin-top: 5px;
      color: #666666;
    }
  }
  .create-card-item:hover {
    box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.15);
    border: 1px solid $border_color;
    .create-card-name {
      color: $color;
    }
  }
}
.createTotalDialog /deep/ {
  .el-dialog {
    background: linear-gradient(
      1deg,
      #ffffff 42%,
      #ffffff 42%,
      #ebedfe 98%,
      #eef0ff 98%
    );
  }
  .el-dialog__title {
    color: $color_title;
    font-weight: bold;
  }
}
</style>
