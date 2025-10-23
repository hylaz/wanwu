<template>
  <CommonLayout
    :aside-title="asideTitle"
    :isButton="true"
    :asideWidth="asideWidth"
    @handleBtnClick="handleBtnClick"
    :isBtnDisabled="sessionStatus === 0"
    :class="[chatType==='webChat'?'chatBg':'']"
  >
    <template #aside-content>
      <div class="explore-aside-app">
        <div
          v-for="n in historyList "
          class="appList"
          :class="['appList',{'disabled':sessionStatus === 0},{'active':n.active}]"
          @click="historyClick(n)"
          @mouseenter="mouseEnter(n)"
          @mouseleave="mouseLeave(n)"
        >
          <span class="appName">
            <span class="appTag"></span>
            {{n.title}}
          </span>
          <span
            class="el-icon-delete appDelete"
            v-if="n.hover || n.active"
            @click.stop="deleteConversation(n)"
          ></span>
        </div>
      </div>
    </template>
    <template #main-content>
      <div class="app-content">
        <!-- <div
          class="app-header-api"
          v-if="chatType === 'agentChat'"
        >
          <div class="header-api-box">
            <div class="header-api-url">
              <el-tag
                effect="plain"
                class="root-url"
              >API根地址</el-tag>
              {{apiURL}}
            </div>
            <el-button
              size="small"
              @click="openApiDialog"
              plain
              class="apikeyBtn"
            >
              <img :src="require('@/assets/imgs/apikey.png')" />
              API密钥
            </el-button>
          </div>
        </div> -->
        <Chat
          :chatType="'chat'"
          :editForm="editForm"
          :appUrlInfo="appUrlInfo"
          :type="chatType"
          ref="agentChat"
          @reloadList="reloadList"
          @setHistoryStatus="setHistoryStatus"
        />
        <!-- <ApiKeyDialog
          ref="apiKeyDialog"
          :appId="editForm.assistantId"
          :appType="'agent'"
          :type="'webChat'"
        /> -->
      </div>
    </template>
  </CommonLayout>
</template>
<script>
import CommonLayout from "@/components/exploreContainer.vue";
import Chat from "./components/chat.vue";
import { mapGetters,mapActions } from "vuex";
// import ApiKeyDialog from "./components/ApiKeyDialog.vue";
import {
  getAgentInfo,
  getOpenurlInfo,
  OpenurlConverList,
  getConversationlist,
} from "@/api/agent";
import { getApiKeyRoot } from "@/api/appspace";
import sseMethod from "@/mixins/sseMethod";
export default {
  components: { CommonLayout, Chat },
  mixins: [sseMethod],
  provide(){
    return{
      getHeaderConfig:this.headerConfig
    }
  },
  data() {
    return {
      asideWidth: "260px",
      apiURL: "",
      asideTitle: "新建对话",
      assistantId: "",
      historyList: [],
      appUrlInfo:{},
      editForm: {
        assistantId: "",
        avatar: {},
        name: "",
        desc: "",
        prologue: "",
        recommendQuestion: [],
      },
      chatType: "agentChat",
      apiStrategies: {
        agentChat_info: getAgentInfo,
        webChat_info: getOpenurlInfo,
        agentChat_converstionList: getConversationlist,
        webChat_converstionList: OpenurlConverList,
      },
      uuid: "",
      STORAGE_KEY: "chatUUID",
    };
  },
  computed: {
    ...mapGetters("app", ["sessionStatus"]),
  },
  created() {
    const id = this.$route.query.id || this.$route.params.id;
    if (id) {
      this.assistantId = id;
      this.editForm.assistantId = id;
    }
    if (this.$route.path.includes("/webChat")) {
      this.chatType = "webChat";
      this.initUUID();
    } else {
      this.chatType = "agentChat";
      // this.apiKeyRootUrl();
    }
    this.getDetail();
    this.getList();
  },
  mounted() {
    if (!localStorage.getItem(this.STORAGE_KEY)) {
      localStorage.setItem(this.STORAGE_KEY, "active");
    }
    window.addEventListener("storage", this.handleStorageEvent);
  },
  beforeDestroy() {
    window.removeEventListener("storage", this.handleStorageEvent);
    this.clearMaxPicNum();
  },
  methods: {
    ...mapActions("app", ["setMaxPicNum","clearMaxPicNum"]),
    initUUID() {
      const storedUUID = localStorage.getItem("chatUUID");
      this.uuid = storedUUID || this.$guid();
      if (!storedUUID) {
        localStorage.setItem("chatUUID", this.uuid);
      }
    },
    handleStorageEvent(event) {
      if (event.key === this.STORAGE_KEY && !event.newValue) {
        this.clearUUID();
      }
    },
    clearUUID() {
      localStorage.removeItem("chatUUID");
      this.uuid = this.$guid();
      localStorage.setItem("chatUUID", this.uuid);
    },
    reloadList(val) {
      this.getList(val);
    },
    headerConfig() {
      if(!this.uuid){
        return { headers: {"X-Client-ID": ''} }
      }
      const config = { 
            headers: { "X-Client-ID": this.uuid}
         }
      return config
    },
    async getDetail() {
      let res = null;
      let data = null;
      if (this.chatType === "agentChat") {
        res = await getAgentInfo({ assistantId: this.editForm.assistantId });
        data = res.data;
      } else {
        const config = this.headerConfig();
        res = await getOpenurlInfo(this.assistantId, config);
        data = res.data.assistant;
        this.appUrlInfo = res.data.appUrlInfo;
      }
      if (res.code === 0) {
        this.editForm.avatar = data.avatar;
        this.editForm.name = data.name;
        this.editForm.desc = data.desc;
        this.editForm.prologue = data.prologue;
        this.setMaxPicNum(data.visionConfig.picNum);
        this.editForm.recommendQuestion = data.recommendQuestion.map(
          (item) => ({ value: item })
        );
      }
    },
    async getList(noInit) {
      let res = null;
      if (this.chatType === "agentChat") {
        res = await getConversationlist({
          assistantId: this.assistantId,
          pageNo: 1,
          pageSize: 1000,
        });
      } else {
        const config = this.headerConfig();
        res = await OpenurlConverList(this.assistantId, config);
      }
      if (res.code === 0) {
        if (res.data.list && res.data.list.length > 0) {
          this.historyList = res.data.list.map((n) => {
            return { ...n, hover: false, active: false };
          });
          if (noInit) {
            this.historyList[0].active = true; //noInit 是true时，左侧默认选中第一个,但是不要调接口刷新详情
          } else {
            this.historyClick[this.historyList[0]];
          }
        } else {
          this.historyList = [];
        }
      } else {
        this.historyList = [];
      }
    },
    setHistoryStatus() {
      this.historyList.forEach((m) => {
        m.active = false;
      });
    },
    historyClick(n) {
      //切换对话
      n.hover = true;
      this.$refs["agentChat"].conversionClick(n);
    },
    deleteConversation(n) {
      this.$refs["agentChat"].preDelConversation(n);
    },
    handleBtnClick() {
      //新建对话
      this.$refs["agentChat"].createConversion();
    },
    mouseEnter(n) {
      n.hover = true;
    },
    mouseLeave(n) {
      n.hover = false;
    },
    apiKeyRootUrl() {
      const data = { appId: this.editForm.assistantId, appType: "agent" };
      getApiKeyRoot(data).then((res) => {
        if (res.code === 0) {
          this.apiURL = res.data || "";
        }
      });
    },
    openApiDialog() {
      this.$refs.apiKeyDialog.showDialog();
    },
  },
};
</script>
<style lang="scss" scoped>
@import "@/style/chat.scss";
.chatBg {
  background: linear-gradient(
    1deg,
    rgb(255, 255, 255) 42%,
    rgb(255, 255, 255) 42%,
    rgb(235, 237, 254) 98%,
    rgb(238, 240, 255) 98%
  );
}
.active {
  background-color: $color_opacity !important;
  .appTag {
    background-color: #384bf7 !important;
  }
}
.explore-aside-app {
  .appList:hover {
    background-color: $color_opacity !important;
  }
  .appList {
    margin: 10px 20px;
    padding: 10px;
    border-radius: 6px;
    margin-bottom: 6px;
    display: flex;
    gap: 8px;
    align-items: center;
    justify-content: space-between;
    cursor: pointer;
    position: relative;
    .appDelete {
      color: #384bf7;
      margin-right: -5px;
      cursor: pointer;
    }
    .appName {
      display: block;
      max-width: 130px;
      overflow: hidden;
      white-space: nowrap;
      pointer-events: none;
      text-overflow: ellipsis;
      .appTag {
        display: inline-block;
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: #ccc;
      }
    }
  }
}
.app-content {
  width: 100%;
  height: 100%;
  .app-header-api {
    width: 50%;
    padding: 10px;
    position: absolute;
    z-index: 999;
    top: 0;
    right: 0;
    display: flex;
    justify-content: flex-end;
    align-content: center;
    .header-api-box {
      display: flex;
      .header-api-url {
        padding: 6px 10px;
        background: #fff;
        margin: 0 10px;
        border-radius: 6px;
        .root-url {
          background-color: #eceefe;
          color: #384bf7;
          border: none;
        }
      }
    }
  }
}
</style>