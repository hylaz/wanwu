<template>
  <!--<el-input
    class="table-search-input"
    v-model="value"
    @change="searchInputChange"
    @keyup.enter.native="addByEnterKey"
    :placeholder="placeholder"
    clearable
  >
    <el-button slot="append" icon="el-icon-search" @click="getTableData"></el-button>
  </el-input>-->
  <el-input
    class="table-search-input no-border-input"
    v-model="value"
    prefix-icon="el-icon-search"
    :placeholder="placeholder"
    @keyup.enter.native="addByEnterKey"
    @change="searchInputChange"
    clearable
  />
</template>

<script>
export default {
  props: {
    placeholder: '',
  },
  data() {
    return {
      value: '',
    };
  },
  methods: {
    async getTableData() {
      this.$emit('handleSearch', this.value);
    },
    searchInputChange(val) {
      !val && this.getTableData();
    },
    addByEnterKey(e) {
      if (e.keyCode === 13) this.getTableData();
    },
    clearValue() {
      this.value = '';
    },
  },
};
</script>

<style lang="scss" scoped>
.table-search-input /deep/ {
  width: 260px;
  margin-right: 15px;
  /* .el-input__inner:focus,
    .el-input__inner:hover {
      border-color: #dcdfe6 !important;
    }*/
}
</style>
