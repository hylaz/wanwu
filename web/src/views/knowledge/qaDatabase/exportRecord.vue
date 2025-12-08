<template>
  <el-dialog
    top="10vh"
    :visible.sync="dialogVisible"
    :close-on-click-modal="false"
    width="70%"
    :before-close="handleClose"
    class="qa-export-dialog"
  >
    <template #title>
      <div class="custom-title">
        <h1>{{ $t('knowledgeManage.qaExport.title') }}</h1>
      </div>
    </template>
    <el-table
      :data="tableData"
      v-loading="tableLoading"
      border
      style="width: 100%"
      :header-cell-style="{ background: '#F9F9F9', color: '#999999' }"
    >
      <el-table-column
        :label="$t('knowledgeManage.qaExport.exportTime')"
        prop="exportTime"
        min-width="200"
      >
        <template slot-scope="scope">
          {{ $formatTime(scope.row.exportTime, undefined, '--') }}
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('knowledgeManage.qaExport.exportName')"
        prop="fileName"
        min-width="160"
      >
        <template slot-scope="scope">
          {{ scope.row.fileName }}
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('knowledgeManage.qaExport.exportUser')"
        prop="author"
        min-width="160"
      >
        <template slot-scope="scope">
          {{ scope.row.author || '--' }}
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('knowledgeManage.qaExport.status')"
        prop="status"
        min-width="140"
      >
        <template slot-scope="scope">
          <span class="status-text">
            {{ formatStatus(scope.row.status) }}
          </span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('knowledgeManage.qaExport.action')"
        width="180"
        fixed="right"
      >
        <template slot-scope="scope">
          <el-button
            type="text"
            size="mini"
            :disabled="
              [STATUS_PENDING, STATUS_PROCESSING, STATUS_FAILED].includes(
                scope.row.status,
              )
            "
            @click="handleDownload(scope.row)"
          >
            {{ $t('knowledgeManage.qaExport.download') }}
          </el-button>
          <el-divider direction="vertical"></el-divider>
          <el-button
            type="text"
            size="mini"
            :disabled="tableLoading"
            @click="handleDelete(scope.row)"
          >
            {{ $t('knowledgeManage.qaExport.delete') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="pagination-wrapper" v-if="pagination.total > 0">
      <el-pagination
        background
        layout="total, prev, pager, next"
        :current-page="pagination.pageNo"
        :page-size="pagination.pageSize"
        :total="pagination.total"
        @current-change="handlePageChange"
      />
    </div>
  </el-dialog>
</template>

<script>
import commonMixin from '@/mixins/common';
import { getQaExportRecordList, delQaRecord } from '@/api/qaDatabase';
import {
  STATUS_FAILED,
  STATUS_FINISHED,
  STATUS_PENDING,
  STATUS_PROCESSING,
} from '@/views/knowledge/constants';

export default {
  name: 'QaExportRecord',
  mixins: [commonMixin],
  data() {
    return {
      dialogVisible: false,
      tableLoading: false,
      tableData: [],
      knowledgeId: '',
      pagination: {
        pageNo: 1,
        pageSize: 10,
        total: 0,
      },
      statusMap: {
        [STATUS_PENDING]: this.$t('knowledgeManage.qaExportStatus.pending'),
        [STATUS_PROCESSING]: this.$t(
          'knowledgeManage.qaExportStatus.processing',
        ),
        [STATUS_FINISHED]: this.$t('knowledgeManage.qaExportStatus.finished'),
        [STATUS_FAILED]: this.$t('knowledgeManage.qaExportStatus.failed'),
      },
      STATUS_FAILED,
      STATUS_FINISHED,
      STATUS_PENDING,
      STATUS_PROCESSING,
    };
  },
  methods: {
    showDialog(knowledgeId) {
      this.knowledgeId = knowledgeId;
      this.dialogVisible = true;
      this.pagination.pageNo = 1;
      this.fetchRecordList();
    },
    handleClose() {
      this.dialogVisible = false;
      this.tableData = [];
      this.pagination.total = 0;
    },
    handlePageChange(page) {
      this.pagination.pageNo = page;
      this.fetchRecordList();
    },
    fetchRecordList() {
      if (!this.knowledgeId) return;
      this.tableLoading = true;
      const params = {
        knowledgeId: this.knowledgeId,
        pageNo: this.pagination.pageNo,
        pageSize: this.pagination.pageSize,
      };
      getQaExportRecordList(params)
        .then(res => {
          if (res.code === 0) {
            const data = res.data || {};
            this.tableData = data.list || [];
            // 为每个表格数据项添加fileName属性，提取自filePath
            this.tableData = this.tableData.map(item => {
              if (item.filePath) {
                const fileName = item.filePath.substring(
                  item.filePath.lastIndexOf('/') + 1,
                );
                return {
                  ...item,
                  fileName,
                };
              }
              return item;
            });
            this.pagination.total = data.total || 0;
          }
        })
        .catch(() => {})
        .finally(() => {
          this.tableLoading = false;
        });
    },
    formatStatus(status) {
      return this.statusMap[status] || this.$t('knowledgeManage.noStatus');
    },
    handleDownload(row) {
      const url = row.filePath;
      const fileName =
        url.substring(url.lastIndexOf('/') + 1) || `download_${Date.now()}`;
      const link = document.createElement('a');
      link.href = url;
      link.download = fileName;
      link.click();
      window.URL.revokeObjectURL(link.href);
    },
    handleDelete(row) {
      const data = {
        exportRecordId: row.exportRecordId,
        knowledgeId: this.knowledgeId,
      };
      this.$confirm(
        this.$t('knowledgeManage.deleteTips'),
        this.$t('knowledgeManage.tip'),
        {
          confirmButtonText: this.$t('common.confirm.confirm'),
          cancelButtonText: this.$t('common.confirm.cancel'),
          type: 'warning',
        },
      )
        .then(() => {
          this.tableLoading = true;
          delQaRecord(data)
            .then(res => {
              if (res.code === 0) {
                this.$message.success(this.$t('common.info.delete'));
                this.fetchRecordList();
              }
            })
            .catch(() => {})
            .finally(() => {
              this.tableLoading = false;
            });
        })
        .catch(() => {});
    },
  },
};
</script>

<style lang="scss" scoped>
.qa-export-dialog {
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

  /deep/ .el-button.is-disabled {
    background: transparent;
  }
}

.custom-title {
  display: flex;
  align-items: center;
  gap: 10px;

  h1 {
    font-size: 18px;
    font-weight: bold;
    margin: 0;
    color: #1f2329;
  }

  .title-tip {
    color: $color;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.status-text {
  &.success {
    color: #67c23a;
  }

  &.error {
    color: #f56c6c;
  }

  &.pending {
    color: #909399;
  }
}
</style>
