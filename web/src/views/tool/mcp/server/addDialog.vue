<template>
  <div class="add-dialog">
    <el-dialog
      :title="title"
      :visible.sync="dialogVisible"
      width="50%"
      :show-close="false"
      :close-on-click-modal="false"
    >
      <div>
        <el-form
          :model="ruleForm"
          :rules="rules"
          ref="ruleForm"
          label-width="130px"
        >
          <el-form-item :label="$t('tool.server.avatar')" prop="avatar">
            <upload-avatar
              :avatar="ruleForm.avatar"
              :default-avatar="defaultAvatar"
              @update-avatar="handleUpdateAvatar"
            />
          </el-form-item>
          <el-form-item :label="$t('tool.server.name')" prop="name">
            <el-input
              v-model="ruleForm.name"
              :placeholder="$t('common.hint.modelName')"
            ></el-input>
          </el-form-item>
          <el-form-item :label="$t('tool.server.desc')" prop="desc">
            <el-input v-model="ruleForm.desc"></el-input>
          </el-form-item>
        </el-form>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button @click="handleClose" size="mini">
          {{ $t('common.button.cancel') }}
        </el-button>
        <el-button
          type="primary"
          size="mini"
          @click="submitForm('ruleForm')"
          :loading="publishLoading"
        >
          {{ $t('common.button.confirm') }}
        </el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { addServer, editServer } from '@/api/mcp';
import uploadAvatar from '@/components/uploadAvatar.vue';

export default {
  components: {
    uploadAvatar,
  },
  data() {
    return {
      dialogVisible: false,
      mcpServerId: '',
      title: '',
      defaultAvatar: require('@/assets/imgs/mcp_active.svg'),
      ruleForm: {
        MCPServerId: '',
        avatar: {
          key: '',
          path: '',
        },
        name: '',
        desc: '',
      },
      rules: {
        name: [
          {
            pattern: /^(?!_)[a-zA-Z0-9-_.\u4e00-\u9fa5]+$/,
            message: this.$t('common.hint.modelName'),
            trigger: 'blur',
          },
          {
            min: 2,
            max: 50,
            message: this.$t('common.hint.modelNameLimit'),
            trigger: 'blur',
          },
          {
            required: true,
            message: this.$t('common.input.placeholder'),
            trigger: 'blur',
          },
        ],
        desc: [
          {
            required: true,
            message:
              this.$t('common.input.placeholder') + this.$t('tool.server.desc'),
            trigger: 'blur',
          },
        ],
      },
      publishLoading: false,
    };
  },
  methods: {
    showDialog(item) {
      this.dialogVisible = true;
      if (item) {
        this.mcpServerId = item.mcpServerId;
        this.ruleForm = item;
        this.title = this.$t('tool.server.editTitle');
      } else this.title = this.$t('tool.server.addTitle');
    },
    handleUpdateAvatar(avatar) {
      this.ruleForm = { ...this.ruleForm, avatar: avatar };
    },
    handleClose() {
      this.dialogVisible = false;
      this.$emit('handleClose', false);
      this.$refs.ruleForm.resetFields();
      this.$refs.ruleForm.clearValidate();
      this.mcpServerId = '';
      this.ruleForm = {
        MCPServerId: '',
        avatar: {
          key: '',
          path: '',
        },
        name: '',
        desc: '',
      };
    },
    submitForm(formName) {
      this.$refs[formName].validate(valid => {
        if (valid) {
          this.publishLoading = true;
          const params = {
            ...this.ruleForm,
          };
          if (this.mcpServerId)
            editServer(params)
              .then(res => {
                if (res.code === 0) {
                  this.$message.success(this.$t('common.info.publish'));
                  this.$emit('handleFetch', false);
                  this.handleClose();
                }
              })
              .finally(() => (this.publishLoading = false));
          else
            addServer(params)
              .then(res => {
                if (res.code === 0) {
                  this.$message.success(this.$t('common.info.publish'));
                  this.$emit('handleFetch', false);
                  this.handleClose();
                  this.$router.push({
                    path: `/tool/detail/server?mcpServerId=${res.data.mcpServerId}`,
                  });
                }
              })
              .finally(() => (this.publishLoading = false));
        }
      });
    },
  },
};
</script>
<style lang="scss" scoped>
.add-dialog {
  .el-button.is-disabled {
    &:active {
      background: transparent !important;
      border-color: #ebeef5 !important;
    }
  }
}
</style>
