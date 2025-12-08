const defaultNo = 1;
const defaultSize = 10;
export default {
  data() {
    return {
      total: 0,
      pageNo: defaultNo,
      pageSize: defaultSize,
      tableData: [],
      searchInfo: {},
    };
  },
  methods: {
    reset() {
      this.searchInfo = {};
    },
    handleSizeChange(val) {
      this.pageSize = val;
      this.getTableData({
        ...this.searchInfo,
        pageSize: val,
        pageNo: defaultNo,
      });
    },
    handleCurrentChange(val) {
      this.pageNo = val;
      this.getTableData({ ...this.searchInfo, pageNo: val });
    },
    justifyIsLastNo(total, pageNo) {
      if (!total && pageNo > 1) {
        this.getTableData({ ...this.searchInfo, pageNo: pageNo - 1 });
      }
    },
    async getTableData(searchInfo, resKey) {
      this.searchInfo = searchInfo;
      const params = {
        pageNo: this.pageNo || defaultNo,
        pageSize: this.pageSize || defaultSize,
        ...this.searchInfo,
      };
      const table = await this.listApi(params);
      const key = resKey || 'list';
      const list = table.data ? table.data[key] || [] : [];
      const tableInfo = table.data ? table.data || {} : {};

      if (table.code === 0) {
        this.tableData = list;
        this.total = table.data.total * 1;
        this.pageNo = table.data.pageNo * 1;
        this.pageSize = table.data.pageSize * 1;
      }
      this.justifyIsLastNo(list ? list.length : 0, params.pageNo);
      this.$emit('refreshData', this.tableData, tableInfo);
      return this.tableData;
    },
  },
};
