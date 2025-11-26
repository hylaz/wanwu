<template>
  <div class="page-wrapper full-content">
    <div class="page-title">
      <i
        class="el-icon-arrow-left"
        @click="goBack"
        style="margin-right: 10px; font-size: 20px; cursor: pointer"
      >
      </i>
      {{ knowledgeName }}
    </div>
    <div class="block table-wrap list-common wrap-fullheight">
      <el-container class="konw_container">
        <el-main class="noPadding">
          <el-container>
            <el-header class="classifyTitle">
              <div class="searchInfo">
                <el-select
                  @change="changeOption($event)"
                  v-model="docQuery.status"
                  clearable
                  :placeholder="$t('knowledgeManage.please')"
                  style="width: 150px"
                  class="marginRight no-border-select cover-input-icon"
                >
                  <el-option
                    v-for="item in knowLegOptions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
                <search-input
                  class="cover-input-icon"
                  :placeholder="$t('knowledgeManage.docPlaceholder')"
                  ref="searchInput"
                  @handleSearch="handleSearch"
                />
              </div>

              <div class="content_title">
                <el-button
                  size="mini"
                  type="primary"
                  icon="el-icon-refresh"
                  @click="reload"
                >
                </el-button>
                <el-button
                  size="mini"
                  type="primary"
                  @click="
                    $router.push(
                      `/knowledge/hitTest?knowledgeId=${docQuery.knowledgeId}&type=qa`
                    )
                  "
                >
                  {{ $t("knowledgeManage.hitTest.name") }}
                </el-button>
                <el-button
                  size="mini"
                  type="primary"
                  @click="showMeta"
                  v-if="hasManagePerm"
                >
                  {{ $t("knowledgeManage.docList.metaDataManagement") }}
                </el-button>
                <template v-if="hasManagePerm">
                  <el-dropdown
                    v-for="(group,index) in dropdownGroups"
                    :key="group.label"
                    @command="handleCommand"
                    :style="{ margin: index === 0 ? '0 10px' : '' }"
                  >
                    <el-button size="mini" type="primary">
                      {{ group.label }}
                      <i :class="['el-icon--right', group.icon]"></i>
                    </el-button>
                    <el-dropdown-menu slot="dropdown">
                      <el-dropdown-item
                        v-for="item in group.items"
                        :key="item.command"
                        :command="item.command"
                      >
                        {{ item.label }}
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </el-dropdown>
                </template>
              </div>
            </el-header>
            <el-main class="noPadding" v-loading="tableLoading">
              <el-table
                ref="dataTable"
                :data="tableData"
                style="width: 100%"
                :row-key="'qaPairId'"
                :header-cell-style="{ background: '#F9F9F9', color: '#999999' }"
                @selection-change="handleSelectionChange"
              >
                <el-table-column
                  type="selection"
                  reserve-selection
                  v-if="hasManagePerm"
                  width="55"
                >
                </el-table-column>
                <el-table-column
                  prop="question"
                  :label="$t('knowledgeManage.qaDatabase.question')"
                  min-width="180"
                >
                  <template slot-scope="scope">
                    <el-popover
                      placement="bottom-start"
                      :content="scope.row.question"
                      trigger="hover"
                      width="300"
                    >
                      <span slot="reference">
                        {{
                          scope.row.question.length > 20
                            ? scope.row.question.slice(0, 20) + "..."
                            : scope.row.question
                        }}
                      </span>
                    </el-popover>
                  </template>
                </el-table-column>
                <el-table-column
                  prop="answer"
                  :label="$t('knowledgeManage.qaDatabase.answer')"
                >
                </el-table-column>
                <el-table-column
                  prop="metaDataList"
                  :label="$t('knowledgeManage.qaDatabase.metaData')"
                >
                  <template slot-scope="scope">
                    <span>
                      {{ getMetaDataText(scope.row.metaDataList) }}
                    </span>
                    <span
                      class="el-icon-edit-outline edit-icon"
                      @click="handleEditMetaData(scope.row)"
                    ></span>
                  </template>
                </el-table-column>
                <el-table-column
                  prop="author"
                  :label="$t('knowledgeManage.author')"
                >
                </el-table-column>
                <el-table-column
                  prop="switch"
                  :label="$t('user.table.status')"
                  width="150"
                >
                  <template slot-scope="scope">
                    <el-switch
                      v-model="scope.row.switch"
                      :active-value="true"
                      :inactive-value="false"
                      @change="handleSwitch(scope.row)"
                    ></el-switch>
                  </template>
                </el-table-column>
                <el-table-column
                  prop="status"
                  :label="$t('knowledgeManage.importStatus')"
                  width="200"
                >
                  <template slot-scope="scope">
                    <span>
                      {{
                        qaImportStatus &&
                        scope.row &&
                        scope.row.status !== undefined
                          ? qaImportStatus[Number(scope.row.status)]
                          : "--"
                      }}
                    </span>
                  </template>
                </el-table-column>
                <el-table-column
                  prop="uploadTime"
                  :label="$t('knowledgeManage.importTime')"
                  width="200"
                >
                </el-table-column>
                <el-table-column
                  :label="$t('knowledgeManage.operate')"
                  width="260"
                >
                  <template slot-scope="scope">
                    <el-button
                      size="mini"
                      round
                      :type="
                        scope.row &&
                        scope.row.status &&
                        [0, 1, 3].includes(Number(scope.row.status))
                          ? 'info'
                          : ''
                      "
                      :disabled="
                        scope.row &&
                        scope.row.status &&
                        [0, 1, 3].includes(Number(scope.row.status))
                      "
                      @click="handleEdit(scope.row)"
                    >
                      {{ $t("common.button.edit") }}
                    </el-button>
                    <el-button
                      size="mini"
                      round
                      @click="handleDel(scope.row)"
                      :disabled="
                        scope.row &&
                        scope.row.status &&
                        [0, 1, 3].includes(Number(scope.row.status))
                      "
                      v-if="hasManagePerm"
                      :type="
                        scope.row &&
                        scope.row.status &&
                        [0, 1, 3].includes(Number(scope.row.status))
                          ? 'info'
                          : ''
                      "
                    >
                      {{ $t("common.button.delete") }}
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
              <!-- 分页 -->
              <Pagination
                class="pagination table-pagination"
                ref="pagination"
                :listApi="listApi"
                :page_size="10"
                @refreshData="refreshData"
              />
            </el-main>
          </el-container>
        </el-main>
      </el-container>
    </div>
    <!-- 元数据管理 -->
    <el-dialog
      :title="$t('knowledgeManage.docList.metaDataManagement')"
      :visible.sync="metaVisible"
      width="550px"
      :before-close="handleClose"
    >
      <mataData
        ref="mataData"
        @updateMeata="updateMeata"
        type="create"
        :knowledgeId="docQuery.knowledgeId"
        class="mataData"
      />
      <span slot="footer" class="dialog-footer">
        <el-button @click="handleClose">
          {{ $t("common.button.cancel") }}
        </el-button>
        <el-button type="primary" @click="createMeta">
          {{ $t("common.button.create") }}
        </el-button>
        <el-button type="primary" @click="submitMeta" :disabled="isDisabled">
          {{ $t("common.button.confirm") }}
        </el-button>
      </span>
    </el-dialog>

    <!-- 批量编辑元数据值弹窗 -->
    <batchMetaData
      ref="batchMetaData"
      :selectedDocIds="selectedDocIds"
      @reLoadDocList="reLoadDocList"
      :type="batchMetaType"
    />
    <!-- 批量编辑元数据值操作框 -->
    <BatchMetatButton
      ref="BatchMetatButton"
      :selectedCount="selectedTableData.length"
      @showBatchMeta="showBatchMeta"
      @handleMetaCancel="handleMetaCancel"
      :type="batchMetaType"
    />
    <!-- 新建/编辑问答对 -->
    <createQa
      ref="createQa"
      @updateData="updateData"
      :knowledgeId="docQuery.knowledgeId"
    />
    <!-- 文件上传 -->
    <fileUpload
      ref="fileUpload"
      @updateData="updateData"
      :knowledgeId="docQuery.knowledgeId"
    />
    <!-- 导出记录 -->
    <exportRecord ref="exportRecord" :knowledgeId="docQuery.knowledgeId" />
  </div>
</template>

<script>
import Pagination from "@/components/pagination.vue";
import SearchInput from "@/components/searchInput.vue";
import mataData from "../metadata.vue";
import batchMetaData from "../meta/batchMetaData.vue";
import BatchMetatButton from "../meta/batchMetatButton.vue";
import createQa from "./createQa.vue";
import fileUpload from "./fileUpload.vue";
import exportRecord from "./exportRecord.vue";
import { updateDocMeta } from "@/api/knowledge";
import {
  getQaPairList,
  delQaPair,
  switchQaPair,
  qaDocExport,
} from "@/api/qaDatabase";
import { mapGetters } from "vuex";
import { COMMUNITY_IMPORT_STATUS, DROPDOWN_GROUPS } from "../../config";
export default {
  components: {
    Pagination,
    SearchInput,
    mataData,
    batchMetaData,
    BatchMetatButton,
    createQa,
    fileUpload,
    exportRecord,
  },
  data() {
    return {
      batchMetaType: "single",
      knowledgeName: "",
      loading: false,
      tableLoading: false,
      docQuery: {
        name: "",
        knowledgeId: this.$route.params.id,
        status: -1,
      },
      fileList: [],
      listApi: getQaPairList,
      tableData: [],
      knowLegOptions: this.getKnowOptions(),
      knowledgeData: [],
      currentKnowValue: null,
      timer: null,
      refreshCount: 0,
      tagList: [],
      metaVisible: false,
      metaData: [],
      isDisabled: false,
      selectedTableData: [],
      selectedDocIds: [],
      qaImportStatus: COMMUNITY_IMPORT_STATUS,
      dropdownGroups: DROPDOWN_GROUPS,
    };
  },
  watch: {
    $route: {
      handler(val) {
        if (val.query.done) {
          this.startTimer();
        }
      },
      immediate: true,
    },
    metaData: {
      handler(val) {
        if (
          val.some((item) => !item.metaKey || !item.metaValueType) ||
          !val.length
        ) {
          this.isDisabled = true;
        } else {
          this.isDisabled = false;
        }
      },
    },
  },
  computed: {
    ...mapGetters("app", ["permissionType"]),
    hasManagePerm() {
      return [10, 20, 30].includes(this.permissionType);
    }
  },
  mounted() {
    this.getTableData(this.docQuery);
    if (
      this.permissionType === -1 ||
      this.permissionType === null ||
      this.permissionType === undefined
    ) {
      const savedData = localStorage.getItem("permission_data");
      if (savedData) {
        try {
          const parsed = JSON.parse(savedData);
          const savedPermissionType =
            parsed && parsed.app && parsed.app.permissionType;
          if (savedPermissionType !== undefined && savedPermissionType !== -1) {
            this.$store.dispatch("app/setPermissionType", savedPermissionType);
          }
        } catch (e) {}
      }
    }
  },
  beforeDestroy() {
    this.clearTimer();
  },
  methods: {
    handleCommand(command) {
      const actions = {
        exportData: this.exportData,
        exportRecord: this.exportRecord,
        createQaPair: this.handleCreateQaPair,
        fileUpload: this.handleUpload,
      };
      (actions[command] || this.handleUpload)();
    },
    exportRecord() {
      this.$refs.exportRecord.showDialog();
    },
    updateData() {
      this.getTableData(this.docQuery);
    },
    exportData() {
      if (!this.docQuery.knowledgeId) {
        this.$message.warning(this.$t("common.noData"));
        return;
      }
      if (this.loading) return;
      const params = {
        knowledgeId: this.docQuery.knowledgeId,
      };
      if (
        this.docQuery.status !== undefined &&
        this.docQuery.status !== null &&
        this.docQuery.status !== "" &&
        this.docQuery.status !== -1
      ) {
        params.status = this.docQuery.status;
      }
      if (this.docQuery.name) {
        params.name = this.docQuery.name.trim();
      }
      this.loading = true;
      qaDocExport(params)
        .then((res) => {
          if (res.code === 0) {
            this.$message.success(this.$t("common.message.success"));
            const data = res.data || {};
            const url = data.fileUrl || data.downloadUrl;
            if (url) {
              window.open(url, "_blank");
            } else if (data.recordCreated) {
              this.exportRecord();
            }
          }
        })
        .catch(() => {})
        .finally(() => {
          this.loading = false;
        });
    },
    handleCreateQaPair() {
      this.$refs.createQa.showDialog();
    },
    handleEditMetaData(row) {
      this.$refs.batchMetaData.showDialog(row);
      this.batchMetaType = "single";
      this.selectedTableData = [row];
      this.selectedDocIds = [row.qaPairId];
    },
    handleSwitch(row) {
      switchQaPair({ qaPairId: row.qaPairId, switch: row.switch }).then(
        (res) => {
          if (res.code === 0) {
            this.$message.success(this.$t("common.message.success"));
            this.getTableData(this.docQuery);
          }
        }
      );
    },
    handleMetaCancel() {
      this.selectedTableData = [];
      this.selectedDocIds = [];
      // 取消所有表格数据的选中状态
      this.$nextTick(() => {
        const table = this.$refs.dataTable;
        if (table) {
          table.clearSelection();
        }
      });
    },
    reLoadDocList() {
      this.getTableData(this.docQuery);
      this.selectedTableData = [];
      this.selectedDocIds = [];

      // 取消所有表格数据的选中状态
      this.$nextTick(() => {
        const table = this.$refs.dataTable;
        if (table) {
          table.clearSelection();
        }
      });
    },
    showBatchMeta() {
      if (!this.selectedTableData || this.selectedTableData.length === 0) {
        this.$message.warning(
          this.$t("knowledgeManage.docList.pleaseSelectDocFirst")
        );
        return;
      }
      this.$refs.batchMetaData.showDialog();
    },
    handleSelectionChange(val) {
      if (val.length > 100) {
        this.$message.warning(
          this.$t("knowledgeManage.docList.maxSelect100Files")
        );
        return;
      }
      this.selectedTableData = val;
      this.batchMetaType = "multiple";
      this.selectedDocIds = val.map((item) => item.qaPairId);
    },
    getMetaDataText(list) {
      if (!list || !Array.isArray(list) || list.length === 0) {
        return "";
      }
      return list.map((item) => item.metaKey).join(", ");
    },
    createMeta() {
      this.$refs.mataData.createMetaData();
      this.scrollToBottom();
    },
    scrollToBottom() {
      this.$nextTick(() => {
        const container = this.$refs.mataData;
        if (container) {
          container.scrollTop = container.scrollHeight;
        }
      });
    },
    submitMeta() {
      this.isDisabled = true;
      const metaList = this.metaData
        .filter((item) => item.option !== "")
        .map(({ metaId, metaKey, metaValueType, option }) => ({
          metaKey,
          ...(option === "add" ? { metaValueType } : {}),
          option,
          ...(option === "update" || option === "delete" ? { metaId } : {}),
        }));
      const data = {
        docId: "",
        knowledgeId: this.docQuery.knowledgeId,
        metaDataList: metaList,
      };
      updateDocMeta(data)
        .then((res) => {
          if (res.code === 0) {
            this.$message.success(this.$t("common.message.success"));
            this.$refs.mataData.getList();
            this.metaVisible = false;
            this.isDisabled = false;
          }
        })
        .catch(() => {
          this.isDisabled = false;
        });
    },
    showMeta() {
      this.metaVisible = true;
    },
    updateMeata(data) {
      this.metaData = data;
    },
    handleClose() {
      this.metaVisible = false;
    },
    startTimer() {
      this.clearTimer();
      if (this.refreshCount >= 2) {
        return;
      }
      const delay = this.refreshCount === 0 ? 1000 : 3000;
      this.timer = setTimeout(() => {
        this.getTableData(this.docQuery);
        this.refreshCount++;
        this.startTimer();
      }, delay);
    },
    clearTimer() {
      if (this.timer) {
        clearInterval(this.timer);
        this.timer = null;
      }
    },
    goBack() {
      this.$router.push({ path: "/knowledge" });
    },
    reload() {
      this.getTableData(this.docQuery);
    },
    handleSearch(val) {
      this.docQuery.name = val;
      this.getTableData(this.docQuery);
    },
    getKnowOptions() {
      const commonOptions = [
        { label: this.$t("knowledgeManage.all"), value: -1 },
        {
          label: this.$t("knowledgeManage.communityReport.taskPending"),
          value: 0,
        },
        {
          label: this.$t("knowledgeManage.communityReport.taskProcessing"),
          value: 1,
        },
        {
          label: this.$t("knowledgeManage.communityReport.taskFinished"),
          value: 2,
        },
        {
          label: this.$t("knowledgeManage.communityReport.taskFailed"),
          value: 3,
        },
      ];
      return commonOptions;
    },
    handleDel(data) {
      this.$confirm(
        this.$t("knowledgeManage.deleteTips"),
        this.$t("knowledgeManage.tip"),
        {
          confirmButtonText: this.$t("common.button.confirm"),
          cancelButtonText: this.$t("common.button.cancel"),
          type: "warning",
        }
      )
        .then(async () => {
          this.loading = true;
          let res = await delQaPair({ qaPairId: data.qaPairId });
          if (res.code === 0) {
            this.$message.success(this.$t("common.info.delete"));
            this.getTableData(this.docQuery); //获取知识分类数据
          }
          this.loading = false;
        })
        .catch((error) => {
          this.getTableData(this.docQuery);
        });
    },
    async getTableData(data) {
      this.tableLoading = true;
      this.tableData = await this.$refs["pagination"].getTableData(data);
      this.tableLoading = false;
    },
    changeOption(data) {
      //通过文档状态查找
      this.docQuery.status = data;
      this.getTableData({ ...this.docQuery, pageNo: 1 });
    },

    handleEdit(row) {
      this.$refs.createQa.showDialog(row);
    },
    handleUpload() {
      this.$refs.fileUpload.showDialog();
    },
    refreshData(data, tableInfo) {
      this.tableData = data;
      if (tableInfo && tableInfo.qaKnowledgeInfo) {
        this.knowledgeName = tableInfo.qaKnowledgeInfo.knowledgeName;
      }
    },
  },
};
</script>
<style lang="scss" scoped>
.mataData {
  max-height: 400px;
  overflow-y: auto;
}
.edit-icon {
  color: $color;
  cursor: pointer;
  font-size: 16px;
  margin-left: 5px;
}
.doc_tag {
  margin: 0 2px;
}
/deep/ {
  .el-button.is-disabled,
  .el-button--info.is-disabled {
    color: #c0c4cc !important;
    background-color: #fff !important;
    border-color: #ebeef5 !important;
  }
  .el-tree--highlight-current
    .el-tree-node.is-current
    > .el-tree-node__content {
    background: #ffefef;
  }
  .el-tabs__item.is-active {
    color: #e60001 !important;
  }
  .el-tabs__active-bar {
    background-color: #e60001 !important;
  }
  .el-tabs__content {
    width: 100%;
    height: calc(100% - 40px);
  }
  .el-tab-pane {
    width: 100%;
    height: 100%;
  }
  .el-tree .el-tree-node__content {
    height: 40px;
  }
  .custom-tree-node {
    padding: 0 10px;
  }
  .el-tree .el-tree-node__content:hover {
    background: #ffefef;
  }
  .el-tree-node__expand-icon {
    display: none;
  }
  .el-button.is-round {
    border-color: #dcdfe6;
    color: #606266;
  }
  .el-upload-list {
    max-height: 200px;
    overflow-y: auto;
  }
  .el-dialog__body {
    padding: 10px 20px;
  }
}
.fileNumber {
  margin-left: 10px;
  display: inline-block;
  padding: 0 20px;
  line-height: 2;
  background: rgb(243, 243, 243);
  border-radius: 8px;
}
.defalutColor {
  color: #e7e7e7 !important;
}
.border {
  border: 1px solid #e4e7ed;
}
.noPadding {
  padding: 0 10px;
}
.activeColor {
  color: #e60001;
}
.error {
  color: #e60001;
}
.marginRight {
  margin-right: 10px;
}
.full-content {
  //padding: 20px 20px 30px 20px;
  margin: auto;
  overflow: auto;
  //background: #fafafa;
  .title {
    font-size: 18px;
    font-weight: bold;
    color: #333;
    padding: 10px 0;
  }
  .tips {
    font-size: 14px;
    color: #aaabb0;
    margin-bottom: 10px;
  }
  .block {
    width: 100%;
    height: calc(100% - 58px);
    .el-tabs {
      width: 100%;
      height: 100%;
      .konw_container {
        width: 100%;
        height: 100%;
        .tree {
          height: 100%;
          background: none;
          .custom-tree-node {
            width: 100%;
            display: flex;
            justify-content: space-between;
            .icon {
              font-size: 16px;
              transform: rotate(90deg);
              color: #aaabb0;
            }
            .nodeLabel {
              color: #e60001;
              display: flex;
              align-items: center;
              .tag {
                display: block;
                width: 5px;
                height: 5px;
                border-radius: 50%;
                background: #e60001;
                margin-right: 5px;
              }
            }
          }
        }
      }
    }
    .classifyTitle {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0 10px;
      h2 {
        font-size: 16px;
      }
      .content_title {
        display: flex;
        align-items: center;
        justify-content: flex-end;
      }
    }
  }
  .uploadTips {
    color: #aaabb0;
    font-size: 12px;
    height: 30px;
  }
  .document_lise {
    list-style: none;
    li {
      display: flex;
      justify-content: space-between;
      font-size: 12px;
      padding: 7px;
      border-radius: 3px;
      line-height: 1;
      .el-icon-success {
        display: block;
      }
      .el-icon-error {
        display: none;
      }
      &:hover {
        cursor: pointer;
        background: #eee;

        .el-icon-success {
          display: none;
        }
        .el-icon-error {
          display: block;
        }
      }
      &.document_loading {
        &:hover {
          cursor: pointer;
          background: #eee;

          .el-icon-success {
            display: none;
          }
          .el-icon-error {
            display: none;
          }
        }
      }
      .el-icon-success {
        color: #67c23a;
      }

      .result_icon {
        float: right;
      }
      .size {
        font-weight: bold;
      }
    }
    .document_error {
      color: red;
    }
  }
}
</style>
<style lang="scss">
.custom-tooltip.is-light {
  border-color: #eee; /* 设置边框颜色 */
  background-color: #fff; /* 设置背景颜色 */
  color: #666; /* 设置文字颜色 */
}
.custom-tooltip.el-tooltip__popper[x-placement^="top"] .popper__arrow::after {
  border-top-color: #fff !important;
}
.custom-tooltip.el-tooltip__popper.is-light[x-placement^="top"] .popper__arrow {
  border-top-color: #ccc !important;
}
</style>