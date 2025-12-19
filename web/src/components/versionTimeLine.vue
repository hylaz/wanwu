<template>
  <div>
    <div class="version-container">
      <div class="version-scroll-wrapper">
        <el-timeline style="margin: 10px 0">
          <el-timeline-item
            v-for="(item, index) in versionList"
            :key="index"
            :type="item.isCurrent ? 'primary' : 'info'"
            :color="item.isCurrent ? '#409EFF' : '#E6A23C'"
          >
            <div
              class="version-status current"
              style="margin-left: 32px"
              v-if="item.isCurrent"
            >
              {{ $t('list.now') }}
            </div>
            <el-card v-else class="version-card">
              <div class="version-header">
                <div class="version-info">
                  <div
                    class="version-status"
                    :class="{
                      current: item.isCurrent,
                      published: !item.isCurrent,
                    }"
                  >
                    {{ $t('list.published') }}
                  </div>
                  <div>
                    <strong>{{ $t('list.version.version') }}:</strong>
                    {{ item.version }}
                  </div>
                  <div>
                    <strong>{{ $t('list.version.desc') }}:</strong>
                    {{ item.desc }}
                  </div>
                  <div>
                    <strong>{{ $t('list.version.createdAt') }}:</strong>
                    {{ item.createdAt }}
                  </div>
                </div>
                <el-dropdown trigger="click" @command="handleCommand">
                  <span class="el-dropdown-link">
                    <i class="el-icon-more"></i>
                  </span>
                  <el-dropdown-menu slot="dropdown">
                    <el-dropdown-item
                      v-if="showExportList.includes(appType)"
                      :command="{ action: 'export', index }"
                    >
                      {{ $t('common.button.export') }}
                    </el-dropdown-item>
                    <el-dropdown-item
                      v-if="!showExportList.includes(appType)"
                      :command="{ action: 'rollback', index }"
                      :divided="appType === 'workflow'"
                    >
                      {{ $t('list.version.rollback') }}
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </el-dropdown>
              </div>
            </el-card>
          </el-timeline-item>
        </el-timeline>
      </div>
    </div>
  </div>
</template>

<script>
import { getAppVersionList, rollbackAppVersion } from '@/api/appspace';
import { exportWorkflow } from '@/api/workflow';
import { WORKFLOW, CHAT } from '@/utils/commonSet';

export default {
  name: 'VersionPopover',
  props: {
    appId: {
      type: String,
      required: true,
      default: () => '',
    },
    appType: {
      type: String,
      required: true,
      default: () => '',
    },
  },
  data() {
    return {
      popoverVisible: false,
      showExportList: [WORKFLOW, CHAT],
      versionList: [
        {
          isCurrent: true,
        },
      ],
    };
  },
  created() {
    this.getAppVersionList();
  },
  methods: {
    getAppVersionList() {
      getAppVersionList({ appId: this.appId, appType: this.appType }).then(
        res => {
          if (res.code === 0 && res.data.list) {
            this.versionList = [
              { isCurrent: true },
              ...res.data.list.map(item => ({ ...item, isCurrent: false })),
            ];
          }
        },
      );
    },
    handleCommand(command) {
      const { action, index } = command;
      switch (action) {
        case 'export':
          this.exportVersion(index);
          break;
        case 'rollback':
          this.rollbackVersion(index);
          break;
      }
    },
    rollbackVersion(index) {
      rollbackAppVersion({
        appId: this.appId,
        appType: this.appType,
        version: this.versionList[index].version,
      }).then(res => {
        if (res.code === 0) {
          this.$message.success(this.$t('common.info.rollback'));
          this.$emit('reloadData');
        }
      });
    },
    exportVersion(index) {
      exportWorkflow(
        { workflow_id: this.appId, version: this.versionList[index].version },
        this.appType,
      ).then(response => {
        const blob = new Blob([response], { type: response.type });
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `${this.$route.query.name || ''}_${this.versionList[index].version}.json`;
        link.click();
        window.URL.revokeObjectURL(link.href);
      });
    },
  },
};
</script>

<style scoped>
.version-container {
  padding: 0 20px;
}

.version-scroll-wrapper {
  max-height: 400px;
  overflow-y: auto;
  padding-right: 10px;
}

.version-scroll-wrapper::-webkit-scrollbar {
  width: 6px;
}

.version-scroll-wrapper::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 10px;
}

.version-scroll-wrapper::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 10px;
}

.version-scroll-wrapper::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

.version-status {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: bold;
  margin-bottom: 8px;
}

.version-status.current {
  background-color: #ecf5ff;
  color: #409eff;
  border: 1px solid #b3d8ff;
}

.version-status.published {
  background-color: #fdf6ec;
  color: #e6a23c;
  border: 1px solid #f5dab1;
}

.version-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.version-info {
  line-height: 1.8;
}

.el-dropdown-link {
  cursor: pointer;
  color: #409eff;
  font-size: 16px;
  margin-left: 10px;
}

.el-dropdown-link:hover {
  color: #66b1ff;
}

.version-card {
  margin-bottom: 10px;
  padding: 12px;
}
</style>
