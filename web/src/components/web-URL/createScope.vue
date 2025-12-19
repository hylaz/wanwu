<template>
  <div class="version-detail-container">
    <div class="version-detail">
      <h3>{{ $t('list.version.detail') }}</h3>
      <div class="detail-form">
        <el-form ref="publishForm" :model="publishForm" :rules="publishRules">
          <el-form-item :label="$t('list.version.no')" prop="version">
            <el-input v-model="publishForm.version" :disabled="true"></el-input>
          </el-form-item>
          <el-form-item :label="$t('list.version.desc')" prop="desc">
            <el-input
              v-model="publishForm.desc"
              :placeholder="$t('list.version.descPlaceholder')"
            ></el-input>
          </el-form-item>
          <el-form-item
            :label="$t('list.version.publishType')"
            prop="publishType"
          >
            <el-radio-group v-model="publishForm.publishType">
              <div>
                <el-radio label="private">
                  {{ $t('agent.form.publishType') }}
                </el-radio>
              </div>
              <div>
                <el-radio label="organization">
                  {{ $t('agent.form.publishType1') }}
                </el-radio>
              </div>
              <div>
                <el-radio label="public">
                  {{ $t('agent.form.publishType2') }}
                </el-radio>
              </div>
            </el-radio-group>
          </el-form-item>

          <div class="saveBtn">
            <el-button size="mini" type="primary" @click="savePublish">
              {{ $t('common.button.save') }}
            </el-button>
          </div>
        </el-form>
      </div>
    </div>

    <div class="version-history">
      <h3>{{ $t('list.version.history') }}</h3>
      <VersionTimeLine :appId="appId" :appType="appType" />
    </div>
  </div>
</template>

<script>
import VersionTimeLine from '@/components/versionTimeLine.vue';
import { getAppLatestVersion, updateAppVersion } from '@/api/appspace';

export default {
  props: {
    appType: {
      type: String,
      required: true,
    },
    appId: {
      type: String,
      required: true,
    },
  },
  components: {
    VersionTimeLine,
  },
  data() {
    return {
      publishForm: {
        publishType: 'private',
        version: '',
        desc: '',
      },
      publishRules: {
        version: [
          {
            required: true,
            message: this.$t('list.version.noPlaceholder'),
            trigger: 'blur',
          },
          {
            pattern: /^v\d+\.\d+\.\d+$/,
            message: this.$t('list.version.versionMsg'),
            trigger: 'blur',
          },
        ],
        desc: [
          {
            required: true,
            message: this.$t('list.version.descPlaceholder'),
            trigger: 'blur',
          },
        ],
        publishType: [
          {
            required: true,
            message: this.$t('common.select.placeholder'),
            trigger: 'change',
          },
        ],
      },
    };
  },
  created() {
    getAppLatestVersion({
      appId: this.appId,
      appType: this.appType,
    }).then(res => {
      this.publishForm = res.data;
    });
  },
  methods: {
    savePublish() {
      this.$refs.publishForm.validate(valid => {
        if (valid) {
          updateAppVersion({
            appId: this.appId,
            appType: this.appType,
            desc: this.publishForm.desc,
            publishType: this.publishForm.publishType,
          }).then(res => {
            if (res.code === 0) {
              this.$message.success(this.$t('common.info.save'));
            }
          });
        }
      });
    },
  },
};
</script>

<style scoped>
.version-detail-container {
  padding: 20px;
  overflow-y: auto;
  height: calc(100vh - 120px);
}

.version-detail,
.version-history {
  border: 1px solid #e6e6e6;
  border-radius: 8px;
  margin-bottom: 20px;
  padding: 20px;
}

.version-detail h3,
.version-history h3 {
  font-size: 18px;
  margin-bottom: 20px;
  border-bottom: 1px solid #e6e6e6;
  padding-bottom: 10px;
}

.detail-form {
  margin-top: 20px;
}

.version-history {
  margin-top: 20px;
}
</style>
