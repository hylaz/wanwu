<template>
  <el-dialog
    title=""
    :visible.sync="dialogVisible"
    width="400px"
    append-to-body
    :close-on-click-modal="false"
    :before-close="handleClose"
  >
    <div class="about-wrap">
      <div v-if="about.logoPath">
        <img
          style="height: 60px; margin: 0 auto"
          :src="basePath + '/user/api' + about.logoPath"
        />
      </div>
      <div class="about-version">
        {{ $t('about.version') }}: {{ about.version || '1.0' }}
      </div>
      <div>
        {{ about.copyright }}
      </div>
    </div>
  </el-dialog>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  data() {
    return {
      dialogVisible: false,
      basePath: this.$basePath,
      about: {},
    };
  },
  watch: {
    commonInfo: {
      handler(val) {
        const { about } = val.data || {};
        this.about = about || {};
      },
      deep: true,
    },
  },
  computed: {
    ...mapGetters('user', ['commonInfo']),
  },
  mounted() {},
  methods: {
    openDialog() {
      this.dialogVisible = true;
    },
    handleClose() {
      this.dialogVisible = false;
    },
  },
};
</script>

<style lang="scss" scoped>
.about-wrap {
  div {
    text-align: center;
  }
  .about-version {
    font-size: 14px;
    color: #121212;
    margin-bottom: 30px;
  }
}
</style>
