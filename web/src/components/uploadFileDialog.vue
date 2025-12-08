<template>
  <div>
    <el-dialog
      :title="title || $t('uploadDialog.title')"
      :visible.sync="dialogVisible"
      append-to-body
      :close-on-click-modal="false"
      width="500px"
    >
      <el-form
        label-width="100px"
        :model="uploadForm"
        :rules="rules"
        ref="uploadForm"
      >
        <el-form-item :label="$t('uploadDialog.file')" prop="file">
          <el-upload
            class="avatar-uploader"
            :action="basePath + '/service/api/v1/inferpub/upload'"
            name="file"
            :headers="headers"
            :show-file-list="false"
            :http-request="handleUpload"
            :on-error="handleAvatarError"
            accept=".json"
          >
            <i
              style="font-size: 26px; margin-right: 16px; margin-top: 5px"
              class="el-icon-upload"
            />
            <span>{{ $t('uploadDialog.hint') }}</span>
            <div style="text-align: left; line-height: normal">
              {{ uploadForm.file ? uploadForm.file.name : '' }}
            </div>
          </el-upload>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="handleClose">{{
          $t('common.button.cancel')
        }}</el-button>
        <el-button :loading="uploading" type="primary" @click="handleSubmit">
          {{ $t('common.button.confirm') }}
        </el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import Pagination from '@/components/pagination.vue';
import { importWorkflow } from '@/api/workflow';
export default {
  props: {
    title: '',
    appType: '',
  },
  components: { Pagination },
  data() {
    const user = this.$store.state.user || {};
    const { uid, orgId } = user.userInfo || {};
    return {
      basePath: this.$basePath,
      dialogVisible: false,
      tableData: [],
      headers: {
        Authorization: 'Bearer ' + user.token,
        'x-user-id': uid,
        'x-org-id': orgId,
      },
      uploadForm: {
        file: '',
      },
      rules: {
        file: [
          {
            required: true,
            message: this.$t('uploadDialog.noUpload'),
            trigger: 'change',
          },
        ],
      },
      uploading: false,
    };
  },
  methods: {
    openDialog() {
      this.dialogVisible = true;
    },
    handleUpload(res) {
      if (res.file) {
        this.uploadForm.file = res.file;
        this.$refs.uploadForm.clearValidate('file');
      }
    },
    handleAvatarError() {
      this.$message.error(this.$t('uploadDialog.uploadError'));
    },
    handleClose() {
      this.dialogVisible = false;
      for (let key in this.uploadForm) {
        this.uploadForm[key] = '';
      }
      this.$refs.uploadForm.resetFields();
    },
    handleSubmit() {
      this.$refs.uploadForm.validate(valid => {
        if (valid) {
          const formData = new FormData();
          const config = { headers: { 'Content-Type': 'multipart/form-data' } };
          for (let key in this.uploadForm) {
            formData.append(key, this.uploadForm[key]);
          }
          this.uploading = true;
          importWorkflow(formData, config, this.appType)
            .then(() => {
              this.uploading = false;
              this.$message.success(this.$t('common.message.success'));
              this.handleClose();
              this.$emit('reloadData');
            })
            .catch(() => (this.uploading = false));
        }
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.avatar-uploader {
  /deep/ .el-upload:focus {
    border-color: #606266 !important;
    color: #606266 !important;
  }
}
</style>
