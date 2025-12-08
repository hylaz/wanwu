<template>
  <img
    style="height: 14px; margin-left: 5px; margin-top: 1px; cursor: pointer"
    src="@/assets/imgs/linkIcon.png"
    @click="handleClick"
  />
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  props: {
    type: '',
  },
  data() {
    return {
      basePath: this.$basePath,
      linkList: {},
    };
  },
  mounted() {
    const { linkList } = this.commonInfo.data || {};
    this.linkList = linkList || {};
  },
  computed: {
    ...mapGetters('user', ['commonInfo']),
  },
  methods: {
    handleClick() {
      const docUrl = this.linkList[this.type];
      if (docUrl) window.open(docUrl);
    },
  },
};
</script>

<style lang="scss" scoped></style>
