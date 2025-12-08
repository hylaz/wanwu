<template>
  <div class="hello" style="color: red">
    <iframe
      :src="`${pdfUrl}?file=${pdfOps.downLink}#page=${pdfOps.page}`"
      width="100%"
      height="100%"
      frameborder="0"
      id="myIframe"
    ></iframe>
  </div>
</template>
<script>
export default {
  name: 'pdf-view',
  data() {
    return {
      pdfUrl:
        window.location.origin + this.$basePath + '/pdfjs/web/viewer.html',
      pdfOps: {
        downLink: '',
        page: 1,
        question: '',
      },
    };
  },
  props: {},
  watch: {
    $route: {
      deep: true,
      handler: function (to) {
        if (to.query) {
          let fileUrl = to.query.fileUrl;
          let question = to.query.question;
          let page = to.query.page;
          this.pdfOps.downLink = fileUrl;
          this.pdfOps.question = question;
          this.pdfOps.page = page;
          if (fileUrl) {
            this.inpage();
          }
        }
      },
      immediate: true,
    },
  },
  methods: {
    inpage() {
      this.$nextTick(() => {
        const iframe = document.getElementById('myIframe');
        iframe.addEventListener('load', () => {
          let postMessage = iframe.contentWindow.postMessage;
          postMessage({ type: 'HIGHT_LINE', value: this.pdfOps.question }, '*');
        });
      });
    },
  },
  mounted() {
    this.$nextTick(() => {
      this.inpage();
    });
  },
};
</script>

<style lang="scss" scoped>
.hello {
  height: 100%;
  overflow: hidden;
}
</style>
