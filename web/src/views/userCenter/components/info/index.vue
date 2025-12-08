<template>
  <div class="routerview-container rl">
    <div class="info page-wrapper hide-loading-bg">
      <p class="page-title form-header rl">
        <i class="el-icon-arrow-left" @click="$router.go(-1)" />
        <img class="page-title-img" src="@/assets/imgs/userInfo.png" alt="" />
        <span class="page-title-name">{{ $t('userInfo.title') }}</span>
      </p>
      <userDetailInfo v-loading="loading" ref="info" />
    </div>
  </div>
</template>

<script>
import { getUserDetail } from '@/api/user';
import userDetailInfo from '../common/userDetailInfo.vue';

export default {
  components: {
    userDetailInfo,
  },
  data() {
    return {
      loading: true,
      userId: this.$store.state.user.userInfo.uid,
    };
  },
  async created() {
    try {
      let res = await getUserDetail({ userId: this.userId });
      if (res.code === 0) {
        this.$refs['info'].setData(res.data);
      }
    } finally {
      this.loading = false;
    }
  },
};
</script>

<style lang="scss" scoped>
.info {
  height: 100%;
  margin: auto;
  border-radius: 4px;
  padding-bottom: 20px;
  .form-header {
    .el-icon-back {
      position: absolute;
      left: 40px;
      top: 35px;
      font-size: 20px;
      color: #999;
      cursor: pointer;
    }
  }
}
.page-title {
  .el-icon-arrow-left {
    margin-right: 10px;
    font-size: 15px;
    cursor: pointer;
    color: $color_title;
  }
}
</style>
