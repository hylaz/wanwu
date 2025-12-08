<template>
  <div>
    <el-dialog
      :title="titleMap[type]"
      :visible.sync="dialogVisible"
      width="750"
      append-to-body
      :close-on-click-modal="false"
    >
      <el-form ref="form" :model="form" label-width="130px" :rules="rules">
        <el-form-item :label="$t('tempSquare.promptPic') + ':'" prop="avatar">
          <el-upload
            class="avatar-uploader"
            action=""
            name="files"
            :disabled="type === detail"
            :show-file-list="false"
            :http-request="handleUploadImage"
            :on-error="handleUploadError"
            accept=".png,.jpg,.jpeg"
          >
            <img
              class="upload-img"
              :src="
                form.avatar && form.avatar.path
                  ? basePath + '/user/api/' + form.avatar.path
                  : defaultIcon || defaultLogo
              "
            />
            <p class="upload-hint" v-if="type !== detail">
              {{ $t('common.fileUpload.clickUploadImg') }}
            </p>
          </el-upload>
        </el-form-item>
        <el-form-item :label="$t('tempSquare.promptName') + ':'" prop="name">
          <el-input
            :placeholder="$t('tempSquare.namePlaceholder')"
            :disabled="type === detail"
            v-model="form.name"
            maxlength="30"
            show-word-limit
          ></el-input>
        </el-form-item>
        <el-form-item :label="$t('tempSquare.promptDesc') + ':'" prop="desc">
          <el-input
            type="textarea"
            :placeholder="$t('tempSquare.descPlaceholder')"
            :disabled="type === detail"
            v-model="form.desc"
            show-word-limit
            maxlength="50"
          ></el-input>
        </el-form-item>
        <el-form-item
          v-if="type !== 'copy'"
          :label="$t('tempSquare.promptText') + ':'"
          prop="prompt"
        >
          <el-tooltip
            effect="dark"
            :content="$t('tempSquare.promptOptimize')"
            placement="top-start"
            popper-class="prompt-optimize-tooltip"
          >
            <i
              v-if="type !== detail"
              class="el-icon-s-help prompt-optimize-icon"
              @click="showPromptOptimize"
            />
          </el-tooltip>
          <el-input
            type="textarea"
            :rows="10"
            :placeholder="$t('tempSquare.promptPlaceholder')"
            :disabled="type === detail"
            v-model="form.prompt"
          ></el-input>
        </el-form-item>
      </el-form>
      <span v-if="type !== detail" slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">
          {{ $t('common.button.cancel') }}
        </el-button>
        <el-button type="primary" @click="doSubmit">
          {{ $t('common.button.confirm') }}
        </el-button>
      </span>
    </el-dialog>
    <PromptOptimize ref="promptOptimize" @promptSubmit="promptSubmit" />
  </div>
</template>

<script>
import { uploadAvatar } from '@/api/user';
import {
  copyPromptTemplate,
  createCustomPrompt,
  editCustomPrompt,
} from '@/api/templateSquare';
import { PROMPT } from '@/views/tool/constants';
import PromptOptimize from '@/components/promptOptimize.vue';

export default {
  components: { PromptOptimize },
  props: {
    type: {
      type: String,
      default: 'create',
    },
    isCustom: false,
  },
  data() {
    return {
      basePath: this.$basePath,
      dialogVisible: false,
      defaultLogo: require('@/assets/imgs/bg-logo.png'),
      defaultIcon: '',
      form: {
        name: '',
        desc: '',
        avatar: {
          key: '',
          path: '',
        },
        prompt: '',
      },
      detail: 'detail',
      titleMap: {
        edit: this.$t('tempSquare.editPrompt'),
        create: this.$t('tempSquare.createPrompt'),
        copy: this.$t('tempSquare.copyPrompt'),
        detail: this.$t('tempSquare.promptDetail'),
      },
      templateId: '',
      customPromptId: '',
      rules: {
        name: [
          {
            required: true,
            message: this.$t('tempSquare.nameRules'),
            trigger: 'change',
          },
          {
            max: 30,
            message: this.$t('tempSquare.promptNameRules'),
            trigger: 'change',
          },
          {
            validator: (rule, value, callback) => {
              if (/^[A-Za-z0-9.\u4e00-\u9fa5_-]+$/.test(value)) {
                callback();
              } else {
                callback(new Error(this.$t('tempSquare.namePlaceholder')));
              }
            },
            trigger: 'change',
          },
        ],
        desc: [
          {
            required: true,
            message: this.$t('tempSquare.descRules'),
            trigger: 'blur',
          },
          {
            max: 50,
            message: this.$t('tempSquare.promptLimitRules'),
            trigger: 'blur',
          },
        ],
        prompt: [
          {
            required: true,
            message: this.$t('tempSquare.promptRules'),
            trigger: 'blur',
          },
        ],
      },
    };
  },
  created() {
    const { defaultIcon = {} } = this.$store.state.user.commonInfo.data || {};
    this.defaultIcon = defaultIcon.promptIcon
      ? this.$basePath + '/user/api/' + defaultIcon.promptIcon
      : '';
  },
  methods: {
    uploadAvatar(file, key) {
      const formData = new FormData();
      const config = { headers: { 'Content-Type': 'multipart/form-data' } };
      formData.append(key, file);
      return uploadAvatar(formData, config);
    },
    handleUploadImage(data) {
      if (data.file) {
        this.uploadAvatar(data.file, 'avatar').then(res => {
          const { key, path } = res.data || {};
          this.form.avatar = { key, path };
        });
      }
    },
    handleUploadError() {
      this.$message.error(this.$t('common.message.uploadError'));
    },
    openDialog(row) {
      if (row) {
        const { templateId, name, desc, avatar, prompt, customPromptId } = row;
        this.templateId = templateId;
        this.customPromptId = customPromptId;
        this.form = { name, desc, avatar, prompt };
      } else {
        this.clearForm();
      }
      this.dialogVisible = true;
      this.$nextTick(() => {
        this.$refs['form'].clearValidate();
      });
    },
    clearForm() {
      this.form = {
        name: '',
        desc: '',
        avatar: {
          key: '',
          path: '',
        },
        prompt: '',
      };
    },
    showPromptOptimize() {
      if (!this.form.prompt) {
        this.$message.warning(this.$t('tempSquare.promptOptimizeHint'));
        return;
      }
      this.$refs.promptOptimize.openDialog(this.form);
    },
    promptSubmit(prompt) {
      this.form.prompt = prompt;
    },
    async doSubmit() {
      await this.$refs.form.validate(async valid => {
        if (valid) {
          if (this.type === 'copy') {
            const form = { ...this.form, templateId: this.templateId };
            delete form.prompt;
            const res = await copyPromptTemplate(form);
            if (res.code === 0) {
              this.$message.success(this.$t('tempSquare.copySuccess'));
              this.dialogVisible = false;
              this.$router.push({ path: '/tool', query: { type: PROMPT } });
            }
            return;
          }

          const res =
            this.type === 'create'
              ? await createCustomPrompt(this.form)
              : await editCustomPrompt({
                  ...this.form,
                  customPromptId: this.customPromptId,
                });
          if (res.code === 0) {
            this.$message.success(this.$t('common.message.success'));
            this.dialogVisible = false;
            this.$emit('reload');
          }
        }
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.avatar-uploader {
  position: relative;
  width: 98px;
  .upload-img {
    object-fit: cover;
    width: 100%;
    height: 98px;
    background: #eee;
    border-radius: 8px;
    border: 1px solid #dcdfe6;
    display: inline-block;
    vertical-align: middle;
  }
  .upload-hint {
    position: absolute;
    width: 100%;
    bottom: 0;
    background: $color_opacity;
    color: $color;
    font-size: 12px;
    line-height: 26px;
    z-index: 10;
    border-radius: 0 0 8px 8px;
  }
}
.prompt-optimize-icon {
  font-size: 16px;
  color: $color;
  margin-top: -10px;
  margin-bottom: 2px;
  float: right;
  cursor: pointer;
}
.prompt-optimize-tooltip {
  z-index: 2100;
}
</style>
