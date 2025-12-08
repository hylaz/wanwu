<template>
  <div>
    <!-- 搜索区 -->
    <slot name="search"></slot>

    <!-- 表格区 -->
    <!--<el-scrollbar>
    <div style="max-height: 5rem;">-->
    <div>
      <el-table
        :data="table.data"
        :row-class-name="tableRowClassName"
        @select="handleSelect"
        @select-all="handleSelectAll"
        @cell-click="handCellClick"
      >
        <slot name="checkbox"></slot>

        <!-- 多选区 -->
        <slot name="selection"></slot>

        <!-- 索引区 -->
        <el-table-column
          v-if="index"
          type="index"
          :index="indexMethod"
          width="62"
          :label="$t('common.table.num')"
        ></el-table-column>

        <!-- 头部 -->
        <slot name="haed"></slot>

        <!-- 表头区 -->
        <el-table-column
          v-for="(column, index) in table.column"
          :key="index"
          :label="column.label"
          :prop="column.prop"
          :width="column.width"
          :formatter="column.formatter"
          :fixed="column.fixed"
        >
          <span v-if="column.html">{{ column.formatter() }}</span>
        </el-table-column>

        <!-- 设置区 -->
        <slot name="other1"></slot>
        <slot name="other2"></slot>
        <slot name="other3"></slot>

        <!-- 按钮区 -->
        <slot name="button"></slot>
      </el-table>
    </div>
    <!--</el-scrollbar>-->
    <!-- 分页区 -->
    <el-pagination
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="page.pageNo"
      :page-sizes="[5, 10, 15, 20]"
      :page-size="page.pageSize"
      layout="total, sizes, prev, pager, next, jumper"
      :total="table.total"
      class="pagination"
    ></el-pagination>
  </div>
</template>

<script>
export default {
  name: 'tableWithPagination',
  props: ['noCreate', 'table', 'index', 'select', 'pageSize', 'noborder'],
  data() {
    return {
      // 表格分页参数
      page: {
        pageNo: 1,
        pageSize: 5,
      },
      // 表格多选集合
      selection: [],
    };
  },
  created: function () {
    // 初始化加载
    this.page.pageSize = this.pageSize || 5;
    if (!this.noCreate) {
      this.handlePagination();
    }
  },
  methods: {
    // 序号索引算法
    indexMethod(index) {
      return (this.page.pageNo - 1) * this.page.pageSize + index + 1;
    },
    // 把每一行的索引放进row
    tableRowClassName({ row, rowIndex }) {
      row.rowIndex = rowIndex;
    },
    // 执行分页逻辑
    handlePagination() {
      this.selection = [];
      this.$emit('handlePagination', this.page);
    },
    handleSizeChange(val) {
      this.page.pageSize = val;
      this.handlePagination();
    },
    handleCurrentChange(val) {
      this.page.pageNo = val;
      this.handlePagination();
    },
    handleSelect(selection, row) {
      this.selection = selection;
    },
    handleSelectAll(selection) {
      this.selection = selection;
    },
    handCellClick(row, column, cell, event) {
      var params = {
        row: row,
        column: column,
        cell: cell,
        event: event,
      };
      this.$emit('handCellClick', params);
    },
    handelInitPage() {
      this.page.pageNo = 1;
      this.page.pageSize = this.pageSize || 5;
    },
  },
};
</script>

<style>
.pagination {
  margin-top: 20px;
  text-align: right;
}
.el-pager li.active {
  color: #d33a3a;
}
</style>
