<template>
  <div class="routerview-container user-center">
    <div class="full-content">
      <div class="router-box rl">
        <router-view style="z-index: 100"></router-view>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {};
  },
};
</script>

<style lang="scss" scoped>
.user-center {
  height: 100%;
  background-color: #f9f9f9;
  .full-content {
    background-color: #fff;
    display: flex;
    .menu-box {
      width: 220px;
      .el-menu {
        border-right: 1px solid transparent !important;
        padding: 20px 0 60px;
      }
    }
    .router-box {
      width: 100%;
      height: 100%;
      border-left: solid 1px #e6e6e6;
      min-height: 100%;
      border-left: solid 1px #e6e6e6;
      overflow: auto;
    }
  }
}
</style>
