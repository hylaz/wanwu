<template>
  <el-dialog
    :visible.sync="dialogVisible"
    width="60%"
    :before-close="handleClose"
    class="batch-meta-dialog"
  >
    <div slot="title" class="custom-title">
      <h1>{{ type !== 'single' ? $t("metaData.batchEdit") : $t("knowledgeManage.meta.editMetaData") }}</h1>
    </div>
    <div class="dialog-content">
      <div class="create-section">
        <el-button type="primary" @click="addMetaData" class="create-btn">
          <i class="el-icon-plus"></i>
          {{ $t("common.button.add") }}
        </el-button>
      </div>

      <div
        class="meta-list"
        v-if="metaDataList.length > 0"
        :loading="docLoading"
      >
        <div
          v-for="(item, index) in metaDataList"
          :key="index"
          class="meta-item"
        >
          <div class="meta-row">
            <div class="field-group">
              <label class="field-label">Key:</label>
              <el-select
                v-model="item.metaKey"
                :placeholder="$t('common.select.placeholder')"
                class="field-select"
                :disabled="item.metaId && item.metaId !== '' ? true : false"
                @change="handleMetaKeyChange($event, item)"
              >
                <el-option
                  v-for="meta in keyOptions"
                  :key="meta.metaKey"
                  :label="meta.metaKey"
                  :value="meta.metaKey"
                />
              </el-select>
            </div>

            <div class="field-group type-group">
              <span class="type-label">
                {{ $t("knowledgeManage.meta.type") }}:
              </span>
              <span class="type-value">[{{ item.metaValueType }}]</span>
            </div>

            <el-divider direction="vertical" class="field-divider" />

            <div class="field-group">
              <label class="field-label">Value:</label>
              <el-tag
                v-if="isJsonArray(item.metaValue)"
                type="info"
                closable
                @close="handleCloseArray(item)"
              >
                {{ $t("knowledgeManage.meta.multipleValue") }}
              </el-tag>
              <template v-else>
                <el-input
                  v-if="item.metaValueType === 'string'"
                  v-model="item.metaValue"
                  :placeholder="$t('common.input.placeholder')"
                  class="field-input"
                  @change="handleMetaValueChange(item, index)"
                />
                <el-input
                  v-if="item.metaValueType === 'number'"
                  v-model="item.metaValue"
                  :placeholder="$t('common.input.placeholder')"
                  class="field-input"
                  type="number"
                  @change="handleMetaValueChange(item, index)"
                />
                <el-date-picker
                  v-if="item.metaValueType === 'time'"
                  v-model="item.metaValue"
                  type="datetime"
                  :placeholder="$t('common.select.placeholder')"
                  class="field-input"
                  format="yyyy-MM-dd HH:mm:ss"
                  value-format="timestamp"
                  @change="handleMetaValueChange(item, index)"
                />
              </template>
            </div>

            <el-divider direction="vertical" class="field-divider" />

            <div class="field-group delete-group">
              <el-button
                type="text"
                @click="removeMetaData(item, index)"
                class="delete-btn"
                icon="el-icon-delete"
              />
            </div>
          </div>
        </div>
      </div>

      <div class="apply-section" v-if="type !== 'single'">
        <el-checkbox v-model="applyToSelected" class="apply-checkbox">
          {{ $t("knowledgeManage.meta.applyAll") }}
        </el-checkbox>
        <el-tooltip
          :content="$t('knowledgeManage.meta.applyAllTips')"
          placement="right"
        >
          <i class="el-icon-question question-icon"></i>
        </el-tooltip>
      </div>
    </div>

    <span slot="footer" class="dialog-footer">
      <el-button @click="handleClose" class="cancel-btn">
        {{ $t("common.confirm.cancel") }}
      </el-button>
      <el-button
        type="primary"
        @click="handleConfirm"
        class="confirm-btn"
        :loading="loading"
      >
        {{ $t("common.confirm.confirm") }}
      </el-button>
    </span>
  </el-dialog>
</template>

<script>
import { metaSelect, getDocMetaList, updateMetaData } from "@/api/knowledge";
export default {
  name: "BatchMetaData",
  props: ["selectedDocIds", "type"],
  data() {
    return {
      dialogVisible: false,
      loading: false,
      docLoading: false,
      applyToSelected: false,
      metaDataList: [],
      keyOptions: [],
    };
  },
  created() {},
  methods: {
    handleMetaKeyChange(val, item) {
      item.metaValueType = this.keyOptions
        .filter((i) => i.metaKey === val)
        .map((e) => e.metaValueType)[0];
    },
    isJsonArray(val) {
      try {
        return Array.isArray(val) && val.length > 1;
      } catch (e) {
        return false;
      }
    },
    handleMetaValueChange(item, index) {
      if (item.metaId && item.originalMetaValue !== item.metaValue) {
        item.option = "update";
      }
    },
    handleCloseArray(item) {
      item.metaValue = "";
    },
    getMetaList() {
      this.docLoading = true;
      getDocMetaList({
        docIdList: this.selectedDocIds,
        knowledgeId: this.$route.params.id,
      })
        .then((res) => {
          if (res.code === 0) {
            this.handleMataData(res.data.knowledgeMetaValues || [])
            this.docLoading = false;
          }
        })
        .catch(() => {
          this.docLoading = false;
        });
    },
    handleMataData(data){
      this.metaDataList = data.map(
        (item) => ({
          ...item,
          metaValue:
            item.metaValue.length > 1
              ? item.metaValue
              : item.metaValueType === "time"
              ? Number(item.metaValue[0])
              : item.metaValue[0],
          option: "existing",
          originalMetaValue: item.metaValue,
        })
      );
    },
    getList() {
      const knowledgeId = this.$route.params.id;
      metaSelect({ knowledgeId })
        .then((res) => {
          if (res.code === 0) {
            this.keyOptions = res.data.knowledgeMetaList || [];
          }
        })
        .catch(() => {});
    },
    showDialog(row=null) {
      this.dialogVisible = true;
      this.applyToSelected = false;
      this.getList();
      if (row) {
        this.handleMataData(row.metaDataList);
      } else {
        this.getMetaList();
      }
    },
    handleClose() {
      this.dialogVisible = false;
      this.$emit("reLoadDocList");
    },
    addMetaData() {
      this.metaDataList.push({
        metaKey: "",
        metaValueType: "string",
        metaValue: "",
        option: "add",
      });
    },

    removeMetaData(item, index) {
      this.$confirm(
        this.$t("knowledgeManage.deleteTips"),
        this.$t("common.confirm.title"),
        {
          confirmButtonText: this.$t("common.confirm.confirm"),
          cancelButtonText: this.$t("common.confirm.cancel"),
          type: "warning",
        }
      )
        .then(() => {
          this.metaDataList.splice(index, 1);
          const data = {
            applyToSelected: this.applyToSelected,
            docIdList: this.selectedDocIds,
            knowledgeId: this.$route.params.id,
            metaValueList: [
              {
                metaId: item.metaId,
                metaKey: item.metaKey,
                option: "delete",
              },
            ],
          };
          this.unpdateMetaApi(data, "delete");
        })
        .catch(() => {
          this.$message.info(this.$t("common.noData"));
        });
    },
    unpdateMetaApi(data, type) {
      this.loading = true;
      updateMetaData(data)
        .then((res) => {
          if (res.code === 0) {
            this.$message.success(this.$t("common.message.success"));
            this.getMetaList();
            if (type === "submit") {
              this.handleClose();
            }
            this.loading = false;
          }
        })
        .catch(() => {});
    },
    handleConfirm() {
      if (this.metaDataList.length === 0) {
        this.$message.warning(this.$t("common.noData"));
        return;
      }

      for (let i = 0; i < this.metaDataList.length; i++) {
        const item = this.metaDataList[i];
        if (!item.metaKey) {
          this.$message.warning(
            `${this.$t("knowledgeManage.meta.rowKey")} ${i + 1} ${this.$t(
              "knowledgeManage.meta.rowKeyTips"
            )}`
          );
          return;
        }
        if (!item.metaValue) {
          this.$message.warning(
            `${this.$t("knowledgeManage.meta.rowKey")} ${i + 1} ${this.$t(
              "knowledgeManage.meta.rowValueTips"
            )}`
          );
          return;
        }
      }

      const updateData = this.metaDataList.filter(
        (item) => item.option === "update" || item.option === "add"
      );

      if (updateData.length === 0) {
        this.$message.info(this.$t("common.noData"));
        return;
      }
      const processedUpdateData = updateData.map((item) => ({
        ...item,
        metaValue: String(item.metaValue),
      }));

      const data = {
        applyToSelected: this.applyToSelected,
        docIdList: this.selectedDocIds,
        metaValueList: processedUpdateData,
        knowledgeId: this.$route.params.id,
      };
      this.unpdateMetaApi(data, "submit");
    },
  },
};
</script>

<style lang="scss" scoped>
.batch-meta-dialog {
  /deep/ .el-dialog__header {
    padding: 20px 20px 10px;
    border-bottom: 1px solid #f0f0f0;
  }
}

.custom-title {
  display: flex;
  align-items: center;
  gap: 10px;

  h1 {
    font-size: 18px;
    font-weight: bold;
    line-height: 24px;
    margin: 0;
  }

  span {
    color: $color;
  }
}

.batch-meta-dialog {
  /deep/ .el-dialog__body {
    padding: 20px;
  }

  /deep/ .el-dialog__footer {
    padding: 10px 20px 20px;
    border-top: 1px solid #f0f0f0;
  }
}

.dialog-content {
  .create-section {
    margin-bottom: 20px;

    .create-btn {
      background: $btn_bg;
      border-color: $btn_bg;
      border-radius: 6px;
      padding: 8px 16px;

      &:hover {
        background: #2a3cc7;
        border-color: #2a3cc7;
      }

      i {
        margin-right: 4px;
      }
    }
  }

  .meta-list {
    .meta-item {
      background: #f7f8fa;
      border-radius: 8px;
      padding: 10px;
      margin-bottom: 12px;

      .meta-row {
        display: flex;
        align-items: center;
        gap: 16px;
        width: 100%;

        .field-divider {
          height: 20px;
          margin: 0 8px;
        }

        .field-group {
          display: flex;
          align-items: center;
          flex: 1;

          &.type-group {
            flex: 0 0 10%;
            justify-content: center;
          }

          &.delete-group {
            flex: 0 0 3%;
            justify-content: center;
          }

          .field-label {
            color: #606266;
            font-size: 14px;
            margin-right: 8px;
            min-width: 40px;
          }

          .field-select {
            flex: 1;
            min-width: 120px;

            /deep/ .el-input__inner {
              border: 1px solid #dcdfe6;
              border-radius: 4px;
              height: 32px;
              line-height: 32px;
            }
          }

          .field-input {
            flex: 1;
            min-width: 120px;

            /deep/ .el-input__inner {
              border: 1px solid #dcdfe6;
              border-radius: 4px;
              height: 32px;
              line-height: 32px;
            }
          }

          .type-label {
            color: #606266;
            font-size: 14px;
            margin-right: 8px;
          }

          .type-value {
            color: $color;
            font-size: 14px;
            font-weight: 500;
          }

          .string-input {
            flex: 1;
            min-width: 120px;

            /deep/ .el-input__inner {
              border: 1px solid #dcdfe6;
              border-radius: 4px;
              height: 32px;
              line-height: 32px;
            }
          }

          .delete-btn {
            color: #f56c6c;
            font-size: 16px;
            padding: 4px;

            &:hover {
              color: #f78989;
            }
          }
        }
      }
    }
  }

  .apply-section {
    display: flex;
    align-items: center;
    margin-top: 20px;
    padding: 0 10px 10px 0;
    .apply-checkbox {
      /deep/ .el-checkbox__label {
        color: #606266;
        font-size: 14px;
      }

      /deep/ .el-checkbox__input.is-checked .el-checkbox__inner {
        background-color: $color;
        border-color: $color;
      }
    }

    .question-icon {
      color: #909399;
      font-size: 16px;
      margin-left: 8px;
      cursor: pointer;

      &:hover {
        color: #606266;
      }
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;

  .cancel-btn {
    background: #fff;
    border: 1px solid #dcdfe6;
    color: #606266;
    border-radius: 4px;
    padding: 8px 16px;

    &:hover {
      color: $color;
      border-color: #c6e2ff;
      background-color: #ecf5ff;
    }
  }

  .confirm-btn {
    background: #f56c6c;
    border-color: #f56c6c;
    border-radius: 4px;
    padding: 8px 16px;

    &:hover {
      background: #f78989;
      border-color: #f78989;
    }

    &:focus {
      background: #f56c6c;
      border-color: #f56c6c;
    }
  }
}
</style>
