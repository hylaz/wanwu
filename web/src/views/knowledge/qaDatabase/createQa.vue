<template>
  <el-dialog
    top="10vh"
    :title="
      isEdit
        ? $t('knowledgeManage.qaDatabase.editInfo')
        : $t('knowledgeManage.createQaDatabase')
    "
    :close-on-click-modal="false"
    :visible.sync="dialogVisible"
    width="60%"
    :before-close="handleClose"
    class="qa-create-dialog"
  >
    <el-form
      :model="ruleForm"
      ref="ruleForm"
      label-width="100px"
      class="qa-form"
      :rules="rules"
      @submit.native.prevent
    >
      <el-form-item
        :label="$t('knowledgeManage.qaDatabase.question') + '：'"
        prop="question"
      >
        <el-input
          v-model="ruleForm.question"
          :placeholder="$t('common.input.placeholder')"
          maxlength="200"
          show-word-limit
        ></el-input>
      </el-form-item>
      <el-form-item
        :label="$t('knowledgeManage.qaDatabase.answer') + '：'"
        prop="answer"
      >
        <el-input
          v-model="ruleForm.answer"
          type="textarea"
          :rows="6"
          :placeholder="$t('common.input.placeholder')"
          maxlength="5000"
          show-word-limit
        ></el-input>
      </el-form-item>
    </el-form>
    <span slot="footer" class="dialog-footer">
      <el-button @click="handleClose">
        {{ $t('common.confirm.cancel') }}
      </el-button>
      <el-button type="primary" @click="submitForm" :loading="loading">
        {{ $t('common.confirm.confirm') }}
      </el-button>
    </span>
  </el-dialog>
</template>

<script>
import { addQaPair, editQaPair } from '@/api/qaDatabase';

export default {
  name: 'CreateQa',
  props: {
    knowledgeId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      dialogVisible: false,
      loading: false,
      isEdit: false,
      editId: null,
      ruleForm: {
        question: '',
        answer: '',
      },
      rules: {
        question: [
          {
            required: true,
            message: this.$t('common.input.placeholder'),
            trigger: 'blur',
          },
          {
            max: 200,
            message:
              this.$t('knowledgeManage.qaDatabase.question') +
              this.$t('common.hint.descLimit'),
            trigger: 'blur',
          },
        ],
        answer: [
          {
            required: true,
            message: this.$t('common.input.placeholder'),
            trigger: 'blur',
          },
          {
            max: 5000,
            message:
              this.$t('knowledgeManage.qaDatabase.answer') +
              this.$t('common.hint.descLimit'),
            trigger: 'blur',
          },
        ],
      },
    };
  },
  methods: {
    showDialog(row = null) {
      this.dialogVisible = true;
      this.isEdit = !!row;

      if (this.isEdit) {
        // 编辑模式，填充数据
        this.editId = row.qaPairId;
        this.ruleForm = {
          question: row.question || '',
          answer: row.answer || '',
        };
      } else {
        // 创建模式，重置表单
        this.resetForm();
      }

      // 清除验证状态
      this.$nextTick(() => {
        if (this.$refs.ruleForm) {
          this.$refs.ruleForm.clearValidate();
        }
      });
    },
    resetForm() {
      this.ruleForm = {
        question: '',
        answer: '',
      };
      this.editId = null;
      this.isEdit = false;
    },
    handleClose() {
      this.dialogVisible = false;
      this.resetForm();
      if (this.$refs.ruleForm) {
        this.$refs.ruleForm.clearValidate();
      }
    },
    submitForm() {
      this.$refs.ruleForm.validate(valid => {
        if (valid) {
          this.loading = true;
          const baseData = {
            question: this.ruleForm.question.trim(),
            answer: this.ruleForm.answer.trim(),
          };

          const apiCall = this.isEdit
            ? editQaPair({ ...baseData, qaPairId: this.editId })
            : addQaPair({ ...baseData, knowledgeId: this.knowledgeId });

          apiCall
            .then(res => {
              if (res.code === 0) {
                this.$message.success(
                  this.isEdit
                    ? this.$t('common.info.edit')
                    : this.$t('common.info.create'),
                );
                this.handleClose();
                this.$emit('updateData');
              }
            })
            .catch(() => {})
            .finally(() => {
              this.loading = false;
            });
        } else {
          return false;
        }
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.qa-create-dialog {
  /deep/ .el-dialog__header {
    padding: 20px 20px 10px;
    border-bottom: 1px solid #f0f0f0;
  }

  /deep/ .el-dialog__body {
    padding: 20px;
  }

  /deep/ .el-dialog__footer {
    padding: 10px 20px 20px;
    border-top: 1px solid #f0f0f0;
  }
}

.qa-form {
  .el-form-item {
    margin-bottom: 20px;
  }

  /deep/ .el-form-item__label {
    font-weight: 500;
    color: #333;
  }

  /deep/ .el-input__inner,
  /deep/ .el-textarea__inner {
    border-radius: 4px;
  }

  /deep/ .el-textarea {
    .el-input__count {
      background: transparent;
      color: #909399;
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;

  .el-button {
    border-radius: 4px;
    padding: 8px 16px;
  }

  .el-button--primary {
    background: $btn_bg;
    border-color: $btn_bg;

    &:hover {
      background: #2a3cc7;
      border-color: #2a3cc7;
    }
  }
}
</style>
