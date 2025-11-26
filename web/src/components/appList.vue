<template>
  <div class="app-card-container">
    <div class="app-card">
      <div class="smart rl smart-create" v-if="isShowTool">
        <div class="app-card-create" @click="showCreate">
          <div class="create-img-wrap">
            <img
              v-if="imgObj[type]"
              class="create-type"
              :src="imgObj[type]"
              alt=""
            />
            <img
              class="create-img"
              src="@/assets/imgs/create_icon.png"
              alt=""
            />
            <div class="create-filter"></div>
          </div>
          <span>{{ `${$t("common.button.add")}${apptype[type] || ""}` }}</span>
        </div>
      </div>
      <div
        v-if="listData && listData.length"
        class="smart rl"
        v-for="(n, i) in listData"
        :key="`${i}sm`"
        :style="`cursor: ${isCanClick(n) ? 'pointer' : 'default'} !important;`"
        @click.stop="isCanClick(n) && toEdit(n)"
      >
        <el-image
          v-if="n.avatar && n.avatar.path"
          class="logo"
          lazy
          :src="basePath + '/user/api/' + n.avatar.path"
          :key="`${i}-${n.appId}-avatar`"
        ></el-image>
        <span :class="['tag-app', `${n.appType}-tag`]">
          {{ apptype[n.appType] || "" }}
        </span>
        <img
          v-if="apptype[n.appType]"
          class="tag-img"
          src="@/assets/imgs/rectangle.png"
          alt=""
        />
        <div class="info rl">
          <p class="name-wrap" :title="n.name">
            <span class="name">{{ n.name }}</span>
            <i
              v-if="isShowPublished && n.publishType"
              class="el-icon-success published-icon"
            />
          </p>
          <el-tooltip
            v-if="n.desc"
            popper-class="instr-tooltip tooltip-cover-arrow"
            effect="dark"
            :content="n.desc"
            placement="bottom-start"
          >
            <p class="desc">{{ n.desc }}</p>
          </el-tooltip>
        </div>
        <div class="tags">
          <span :class="['smartDate']">{{ n.createdAt }}</span>
          <div v-if="!isShowTool" class="favorite-wrap">
            <el-tooltip
              class="item"
              effect="dark"
              :content="n.user.userName"
              placement="top-start"
            >
              <span class="user-name">
                {{
                  n.user
                    ? n.user.userName.length > 6
                      ? n.user.userName.substring(0, 6) + "..."
                      : n.user.userName
                    : ""
                }}
              </span>
            </el-tooltip>
            <img
              v-if="!n.isFavorite"
              class="favorite"
              src="@/assets/imgs/like.png"
              alt=""
              @click="handelMark($event, n, i)"
            />
            <img
              v-else
              class="favorite"
              src="@/assets/imgs/like_active.png"
              alt=""
              @click="handelMark($event, n, i)"
            />
          </div>
        </div>
        <div v-if="isShowPublished && n.publishType" class="publishType">
          <span v-if="n.publishType === 'private'" class="publishType-tag">
            <span class="el-icon-lock"></span> {{ $t("appSpace.private") }}
          </span>
          <span v-else class="publishType-tag">
            <span class="el-icon-unlock"></span> {{ $t("appSpace.public") }}
          </span>
        </div>
        <div class="editor" v-if="isShowTool">
          <el-dropdown @command="handleClick($event, n)" placement="top">
            <span class="el-dropdown-link">
              <i class="el-icon-more icon edit-icon" @click.stop />
            </span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item command="edit" v-if="isCanClick(n)">
                {{ $t("common.button.edit") }}
              </el-dropdown-item>
              <el-dropdown-item command="delete">
                {{ $t("common.button.delete") }}
              </el-dropdown-item>
              <el-dropdown-item command="copy">
                {{ $t("common.button.copy") }}
              </el-dropdown-item>
              <!--不在卡片进行发布-->
              <!--<el-dropdown-item
                command="publish"
                v-if="n.appType === workflow && !n.publishType"
              >
                {{$t('common.button.publish')}}
              </el-dropdown-item>-->
              <el-dropdown-item command="cancelPublish" v-if="n.publishType">
                {{ $t("common.button.cancelPublish") }}
              </el-dropdown-item>
              <el-dropdown-item command="publishSet">
                {{ $t("appSpace.publishSet") }}
              </el-dropdown-item>
              <el-dropdown-item
                command="export"
                v-if="[workflow, chat].includes(n.appType)"
              >
                {{ $t("common.button.export") }}
              </el-dropdown-item>
              <el-dropdown-item
                command="transform"
                v-if="[workflow, chat].includes(n.appType)"
              >
                {{
                  $t("common.button.transform") +
                  (n.appType === workflow
                    ? $t("appSpace.chat")
                    : $t("appSpace.workflow"))
                }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
      </div>
    </div>
    <el-empty
      class="noData"
      v-if="!(listData && listData.length)"
      :description="$t('common.noData')"
    ></el-empty>
    <el-dialog
      :title="$t('list.tips')"
      :visible.sync="dialogVisible"
      width="400px"
      append-to-body
      :close-on-click-modal="false"
      :before-close="handleClose"
      class="createTotalDialog"
    >
      <div style="margin-top: -20px">
        <div
          v-for="item in publishList"
          :key="item.key"
          style="margin-bottom: 5px"
        >
          <el-radio :label="item.key" v-model="publishType">
            {{ item.value }}
          </el-radio>
        </div>
        <div style="text-align: right; margin-top: 15px; margin-bottom: -10px">
          <el-button size="mini" type="primary" @click="doPublish">
            {{ $t("common.button.confirm") }}
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { AppType } from "@/utils/commonSet";
import {
  deleteApp,
  appCancelPublish,
  appPublish,
  copyTextQues,
  copyAgentApp,
} from "@/api/appspace";
import {
  copyWorkFlow,
  exportWorkflow,
  transformWorkflow,
} from "@/api/workflow";
import { setFavorite } from "@/api/explore";
import { AGENT, RAG, CHAT, WORKFLOW } from "@/utils/commonSet";

export default {
  props: {
    type: String,
    showCreate: Function,
    appData: {
      type: Array,
      required: true,
      default: [],
    },
    isShowTool: false,
    isShowPublished: false,
    appFrom: {
      type: String,
      default: "",
    },
  },
  watch: {
    appData: {
      handler: function (val) {
        this.listData = val;
      },
      immediate: true,
      deep: true,
    },
  },
  data() {
    return {
      apptype: AppType,
      basePath: this.$basePath,
      workflow: WORKFLOW,
      chat: CHAT,
      listData: [],
      row: {},
      publishType: "private",
      dialogVisible: false,
      publishList: [
        { key: "private", value: this.$t("workflow.publishText") },
        { key: "organization", value: this.$t("workflow.publicOrgText") },
        { key: "public", value: this.$t("workflow.publicTotalText") },
      ],
      imgObj: {
        [WORKFLOW]: require(`@/assets/imgs/create_workflow.png`),
        [CHAT]: require(`@/assets/imgs/create_chatflow.svg`),
        [AGENT]: require(`@/assets/imgs/create_agent.png`),
        [RAG]: require(`@/assets/imgs/create_rag.png`),
      },
    };
  },
  methods: {
    handleClose() {
      this.dialogVisible = false;
    },
    isCanClick(n) {
      return this.isShowTool
        ? ([WORKFLOW, CHAT].includes(n.appType) && !n.publishType) ||
            ![WORKFLOW, CHAT].includes(n.appType)
        : true;
    },
    // 公用删除方法
    async handleDelete() {
      const params = {
        appId: this.row.appId,
        appType: this.row.appType,
      };
      const res = await deleteApp(params);
      if (res.code === 0) {
        this.$message.success(this.$t("list.delSuccess"));
        this.$emit("reloadData");
      }
    },
    workflowEdit(row) {
      const querys = {
        id: row.appId,
      };
      this.$router.push({ path: "/workflow", query: querys });
    },
    workflowDelete(row) {
      this.row = row;
      this.$alert(this.$t("list.deleteTips"), this.$t("list.tips"), {
        confirmButtonText: this.$t("list.confirm"),
        callback: (action) => {
          if (action === "confirm") {
            this.handleDelete();
          }
        },
      });
    },
    async workflowCopy(row) {
      const params = { workflow_id: row.appId };
      const res = await copyWorkFlow(params, row.appType);

      if (res.code === 0) {
        this.$router.push({
          path: "/workflow",
          query: { id: res.data.workflow_id },
        });
      }
    },
    workflowPublish(row) {
      this.row = row;
      this.dialogVisible = true;
      this.publishType = "private";
    },
    async doPublish() {
      const params = {
        appId: this.row.appId,
        appType: this.row.appType,
        publishType: this.publishType,
      };
      const res = await appPublish(params);
      if (res.code === 0) {
        this.$message.success(this.$t("list.publicSuccess"));
        this.handleClose();
        this.$emit("reloadData");
      }
    },
    async cancelPublish(row) {
      let confirmed = true;
      const params = {
        appId: row.appId,
        appType: row.appType,
      };

      //工作流取消发布，需弹窗提示
      if (row.appType === WORKFLOW) {
        confirmed = await this.showDeleteConfirm(this.$t("list.cancelHint"));
      }

      if (confirmed) {
        const res = await appCancelPublish(params);
        if (res.code === 0) {
          this.$message.success(this.$t("common.message.success"));
          this.$emit("reloadData");
        }
      }
    },
    workflowExport(row) {
      exportWorkflow({ workflow_id: row.appId }, row.appType).then(
        (response) => {
          const blob = new Blob([response], { type: response.type });
          const url = URL.createObjectURL(blob);
          const link = document.createElement("a");
          link.href = url;
          link.download = row.name + ".json";
          link.click();
          window.URL.revokeObjectURL(link.href);
        }
      );
    },
    jumpToWorkflowPublicSet(row) {
      this.$router.push({
        path: `/workflow/publishSet`,
        query: { appId: row.appId, appType: row.appType, name: row.name },
      });
    },
    workflowTransform(row) {
      transformWorkflow({ workflow_id: row.appId }, row.appType).then(() => {
        this.$emit("reloadData");
      });
    },
    workflowOperation(method, row) {
      switch (method) {
        case "edit":
          this.workflowEdit(row);
          break;
        case "delete":
          this.workflowDelete(row);
          break;
        case "copy":
          this.workflowCopy(row);
          break;
        case "publish":
          this.workflowPublish(row);
          break;
        case "cancelPublish":
          this.cancelPublish(row);
          break;
        case "publishSet":
          this.jumpToWorkflowPublicSet(row);
          break;
        case "export":
          this.workflowExport(row);
          break;
        case "transform":
          this.workflowTransform(row);
          break;
      }
    },
    chatDelete(row) {
      this.row = row;
      this.$alert(this.$t("list.deleteChatTips"), this.$t("list.tips"), {
        confirmButtonText: this.$t("list.confirm"),
        callback: (action) => {
          if (action === "confirm") {
            this.handleDelete();
          }
        },
      });
    },
    chatOperation(method, row) {
      switch (method) {
        case "edit":
          this.workflowEdit(row);
          break;
        case "delete":
          this.chatDelete(row);
          break;
        case "copy":
          this.workflowCopy(row);
          break;
        case "cancelPublish":
          this.cancelPublish(row);
          break;
        case "publishSet":
          this.jumpToWorkflowPublicSet(row);
          break;
        case "export":
          this.workflowExport(row);
          break;
        case "transform":
          this.workflowTransform(row);
          break;
      }
    },
    async showDeleteConfirm(tips) {
      try {
        await this.$alert(tips, this.$t("list.tips"), {
          confirmButtonText: this.$t("list.confirm"),
        });
        return true;
      } catch (err) {
        return false;
      }
    },
    intelligentEdit(row) {
      this.$router.push({
        path: "/agent/test",
        query: {
          id: row.appId,
          ...(row.publishType !== "" && { publish: true }),
        },
      });
    },
    intelligentDelete(row) {
      this.row = row;
      this.handleDelete();
    },
    intelligentCopy(row) {
      copyAgentApp({ assistantId: row.appId })
        .then((res) => {
          if (res.code === 0) {
            const id = res.data.assistantId;
            this.$message.success(this.$t("list.copySuccess"));
            this.$router.push({ path: `/agent/test?id=${id}` });
          }
        })
        .catch(() => {});
    },
    intelligentOperation(method, row) {
      switch (method) {
        case "edit":
          // 智能体编辑
          this.intelligentEdit(row);
          break;
        case "delete":
          // 智能体删除
          this.intelligentDelete(row);
          break;
        case "copy":
          // 智能体复制
          this.intelligentCopy(row);
          break;
        case "cancelPublish":
          this.cancelPublish(row);
          break;
        case "publishSet":
          //发布设置
          this.$router.push({
            path: `/agent/publishSet`,
            query: { appId: row.appId, appType: row.appType, name: row.name },
          });
          break;
      }
    },
    txtQuesEdit(row) {
      this.$router.push({
        path: "/rag/test",
        query: {
          id: row.appId,
          ...(row.publishType !== "" && { publish: true }),
        },
      });
    },
    txtQuesDelete(row) {
      this.row = row;
      this.handleDelete();
    },
    txtQuesCopy(row) {
      copyTextQues({ ragId: row.appId })
        .then((res) => {
          if (res.code === 0) {
            const id = res.data.ragId;
            this.$message.success(this.$t("list.copySuccess"));
            this.$router.push({ path: `/rag/test?id=${id}` });
          }
        })
        .catch(() => {});
    },
    txtQuesOperation(method, row) {
      switch (method) {
        case "edit":
          // 文本问答编辑
          this.txtQuesEdit(row);
          break;
        case "delete":
          // 文本问答删除
          this.txtQuesDelete(row);
          break;
        case "copy":
          // 文本问答复制
          this.txtQuesCopy(row);
          break;
        case "cancelPublish":
          this.cancelPublish(row);
          break;
        case "publishSet":
          this.$router.push({
            path: `/rag/publishSet`,
            query: { appId: row.appId, appType: row.appType, name: row.name },
          });
          break;
      }
    },
    commonToChat(row) {
      const type = row.appType;
      switch (type) {
        case AGENT:
          this.$router.push({
            path: "/explore/agent",
            query: { id: row.appId },
          });
          break;
        case RAG:
          this.$router.push({ path: "/explore/rag", query: { id: row.appId } });
          break;
        case WORKFLOW:
          this.$router.push({
            path: "/explore/workflow",
            query: { id: row.appId },
          });
          break;
      }
    },
    commonMethods(method, row) {
      const type = row.appType;
      switch (type) {
        case AGENT:
          this.intelligentOperation(method, row);
          break;
        case RAG:
          this.txtQuesOperation(method, row);
          break;
        case WORKFLOW:
          this.workflowOperation(method, row);
          break;
        case CHAT:
          this.chatOperation(method, row);
          break;
      }
    },
    handleClick(command, row) {
      this.commonMethods(command, row);
    },
    toEdit(row) {
      if (this.appFrom === "explore") {
        this.commonToChat(row);
      } else {
        this.commonMethods("edit", row);
      }
    },
    handelMark(e, n, i) {
      e.stopPropagation();
      this.$confirm(
        n.isFavorite
          ? this.$t("explore.unFavorite")
          : this.$t("explore.favorite"),
        this.$t("common.confirm.title"),
        {
          confirmButtonText: this.$t("common.confirm.confirm"),
          cancelButtonText: this.$t("common.confirm.cancel"),
          type: "warning",
        }
      )
        .then(() => {
          setFavorite({
            appId: n.appId,
            appType: n.appType,
            isFavorite: !n.isFavorite,
          }).then((res) => {
            if (res.code === 0) {
              this.$message.success(
                n.isFavorite
                  ? this.$t("explore.delSuccess")
                  : this.$t("explore.setSuccess")
              );
              const list = [...this.listData];
              list[i].isFavorite = !n.isFavorite;
              this.listData = [...list];
              // this.getHistoryList();
            }
          });
        })
        .catch(() => {});
    },
  },
};
</script>

<style lang="scss" scoped>
@import "@/style/appCard.scss";
.noData {
  padding: 30px 0;
}
</style>
