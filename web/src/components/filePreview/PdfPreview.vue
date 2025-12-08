<template>
  <iframe
    id="pdf-view-iframe"
    :src="fileUrl"
    width="100%"
    height="100%"
  ></iframe>
</template>
<script>
export default {
  name: 'PdfPeview',
  components: {},
  data() {
    return {
      fileUrl: '',
    };
  },
  watch: {
    $route: {
      deep: true,
      handler: function (to) {
        if (to.query) {
          let fileUrl = to.query.fileUrl;
          this.fileUrl = fileUrl;
        }
      },
      immediate: true,
    },
  },
  created() {},
  mounted() {
    this.$nextTick(() => {});
  },
  methods: {},
};
</script>
<style lang="scss" scoped></style>
