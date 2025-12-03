<!--http流式 前端保存点赞效果，后端不保存-->
<template>
  <div class="session rl">
    <div class="session-setting">
      <el-dropdown
        class="right-setting"
        @command="gropdownClick"
      >
        <i
          class="el-icon-more"
          trigger="click"
          style="color:var(--color);"
        ></i>
        <el-dropdown-menu
          :append-to-body="false"
          placement="bottom-end"
          slot="dropdown"
        >
          <el-dropdown-item command="clear">清空对话</el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>

    <div
      class="history-box showScroll"
      id="timeScroll"
      v-loading="loading"
    >
      <div
        v-for="(n,i) in session_data.history"
        :key="`${i}sdhs`"
      >
        <!--问题-->
        <div
          v-if="n.query"
          class="session-question"
        >
          <div :class="['session-item','rl']">
            <img
              class="logo"
              :src="'/user/api/'+ userAvatar"
            />
            <div class="answer-content">
              <div class="answer-content-query">
                <!-- <span class="session-setting-id" v-if="$route.params && $route.params.id">ragID: {{$route.params.id}}</span> -->
                <el-popover
                  placement="bottom-start"
                  trigger="click"
                  :visible-arrow="false"
                  popper-class="query-copy-popover"
                  content=""
                >
                  <p
                    class="query-copy"
                    @click="queryCopy(n.query)"
                    style="cursor: pointer"
                  >
                    <i class="el-icon-s-order"></i>
                    &nbsp;{{$t('agent.copyToInput')}}
                  </p>
                  <span
                    slot="reference"
                    class="answer-text"
                  >
                    {{n.query}}
                  </span>
                </el-popover>
              </div>
            </div>
          </div>
        </div>
        <!--loading-->
        <div
          v-if="n.responseLoading"
          class="session-answer"
        >
          <div class="session-answer-wrapper">
            <img
              class="logo"
              :src="'/user/api/'+ defaultUrl"
            />
            <div class="answer-content"><i class="el-icon-loading"></i></div>
          </div>
        </div>
        <!--pending-->
        <div
          v-if="n.pendingResponse"
          class="session-answer"
        >
          <div class="session-answer-wrapper">
            <img
              class="logo"
              :src="'/user/api/'+ defaultUrl"
            />
            <div
              class="answer-content"
              style="padding:0 10px;color:#E6A23C;"
            >
              {{n.pendingResponse}}
            </div>
          </div>
        </div>
        <!-- 回答故障  code:7-->
        <div
          class="session-error"
          v-if="n.error"
        >
          <i class="el-icon-warning"></i>
          &nbsp;{{n.response}}
        </div>

        <!--回答 文字+图片-->
        <div
          v-if="(!n.error) && (n.response || n.msg_type)"
          class="session-answer"
          :id="'message-container'+i"
        >
          <div class="session-answer-wrapper">
            <img
              class="logo"
              :src="'/user/api/'+ defaultUrl"
            />
            <div
              class="session-wrap"
              style="width:calc(100% - 30px);"
            >
              <div class="deepseek" v-if="n.msg_type && ['qa_start','qa_finish','knowledge_start'].includes(n.msg_type)">
                <img
                  :src="require('@/assets/imgs/think-icon.png')"
                  class="think_icon"
                />
                {{getTitle(n.msg_type) }}
              </div>
              <template v-else>
                <img
                  :src="require('@/assets/imgs/think-icon.png')"
                  class="think_icon"
                />
                <div
                  v-if="showDSBtn(n.response)"
                  class="deepseek"
                  @click="toggle($event,i)"
                >
                  {{n.thinkText}}
                  <i v-bind:class="{'el-icon-arrow-down': !n.isOpen,'el-icon-arrow-up': n.isOpen}"></i>
                </div>
                <span v-else class="deepseek">{{$t('menu.knowledge')}}</span>
              </template>
              <!--内容-->
              <div
                v-if="n.response"
                class="answer-content"
                :id="i"
                v-bind:class="{'ds-res':showDSBtn(n.response)}"
                v-html="showDSBtn(n.response)?replaceHTML(n.response,n):n.response"
              ></div>
              <!--loading-->
              <div
                v-if="n.finish === 0 && sessionStatus == 0 && i === session_data.history.length - 1"
                class="text-loading"
              >
                <div></div>
                <div></div>
                <div></div>
              </div>
              <!--出处-->
              <div
                v-if="n.searchList && n.searchList.length && n.finish === 1"
                class="search-list"
              >
                <h2 class="recommended-question-title" v-if="n.msg_type && ['qa_finish'].includes(n.msg_type)">{{$t('app.recommendedQuestion')}}</h2>
                <div
                  v-for="(m,j) in n.searchList"
                  :key="`${j}sdsl`"
                  class="search-list-item"
                >
                  <div v-if="m.content_type && m.content_type === 'qa'" class="qa_content" @click="handleRecommendedQuestion(m)">
                    <span>{{j+1}}. {{m.question}}</span>
                  </div>
                  <template v-else>
                    <div
                      class="serach-list-item"
                      v-if="n.citations && n.citations.includes(j+1)"
                    >
                      <span @click="collapseClick(n,m,j)"><i :class="['',m.collapse?'el-icon-caret-bottom':'el-icon-caret-right']"></i>出处：</span>
                      <a
                        v-if="m.link"
                        :href="m.link"
                        target="_blank"
                        rel="noopener noreferrer"
                        class="link"
                      >{{m.link}}</a>
                      <span v-if="m.title">
                        <sub
                          class="subTag"
                          :data-parents-index="i"
                          :data-collapse="m.collapse?'true':'false'"
                        >{{j + 1}}</sub> {{m.title}}
                      </span>
                      <!-- <span @click="goPreview($event,m)" class="search-doc">查看全文</span> -->
                    </div>
                    <el-collapse-transition>
                      <div
                        v-show="m.collapse?true:false"
                        class="snippet"
                      >
                        <p v-html="m.snippet"></p>
                      </div>
                    </el-collapse-transition>
                  </template>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { marked } from "marked";
import smoothscroll from "smoothscroll-polyfill";
var highlight = require("highlight.js");
import "highlight.js/styles/atom-one-dark.css";
import commonMixin from "@/mixins/common";
import { mapGetters } from "vuex";

marked.setOptions({
  renderer: new marked.Renderer(),
  gfm: true,
  tables: true,
  breaks: false,
  pedantic: false,
  sanitize: false,
  smartLists: true,
  smartypants: false,
  highlight: function (code) {
    return highlight.highlightAuto(code).value;
  },
});


export default {
  props: ["sessionStatus", "defaultUrl"],
  mixins: [commonMixin],
  data() {
    return {
      autoScroll: true,
      scrollTimeout: null,
      isDs:
        [
          "txt2txt-002-001",
          "txt2txt-002-002",
          "txt2txt-002-004",
          "txt2txt-002-005",
          "txt2txt-002-006",
          "txt2txt-002-007",
          "txt2txt-002-008",
        ].indexOf(this.$route.params.id) != -1,
      loading: false,
      marked: marked,
      session_data: {
        tool: "",
        searchList: [],
        history: [],
        response: "",
      },
      basePath: this.$basePath,
      current_data: [],
      //标注相关
      c: null,
      ctx: null,
      canvasShow: false,
      cv: null,
      currImg: {
        url: "",
        width: 0, // 原始宽高
        height: 0,
        w: 0, // 压缩后的宽高
        h: 358,
        roteX: 0, // 压缩后的比例
        roteY: 0,
      },
      imgConfig: ["jpeg", "PNG", "png", "JPG", "jpg", "bmp", "webp"],
      audioConfig: ["mp3", "wav"],
    };
  },
  computed: {
    ...mapGetters('user', ['userAvatar'])
  },
  watch: {
    sessionStatus: {
      handler(val, oldVal) {},
      immediate: true,
    },
  },
  mounted() {
    this.setupScrollListener();
    // this.listenerImg();
    smoothscroll.polyfill();
    document.addEventListener('click', this.handleCitationClick);
  },
  beforeDestroy() {
    if(this.handleCitationClick) {
      document.removeEventListener('click', this.handleCitationClick);
    }
    const container = document.getElementById("timeScroll");
    if (container) {
      container.removeEventListener("scroll", this.handleScroll);
    }
    clearTimeout(this.scrollTimeout);
  },
  methods: {
    handleRecommendedQuestion(m){
      this.$emit("handleRecommendedQuestion", m.question);
    },
    handleCitationClick(e) {
      // 调用 common.js 中的通用方法
      this.$handleCitationClick(e, {
        sessionStatus: this.sessionStatus,
        sessionData: this.session_data,
        citationSelector: '.citation',
        scrollElementId: 'timeScroll',
        onToggleCollapse: (item, collapse) => {
          // 使用 Vue.set 确保响应式更新
          this.$set(item, 'collapse', collapse);
        }
      });
    },
    setCitations(index) {
      let citation = `#message-container${index} .citation`;
      const allCitations = document.querySelectorAll(citation);
      const citationsSet = new Set();

      allCitations.forEach((element) => {
        const text = element.textContent.trim();
        if (text) {
          citationsSet.add(Number(text));
        }
      });

      return Array.from(citationsSet);
    },
    goPreview(event, item) {
      event.stopPropagation(); // 阻止事件冒泡
      let { meta_data } = item;
      let { file_name, download_link, page_num, row_num, sheet_name } =
        meta_data;
      var index = file_name.lastIndexOf(".");
      var ext = file_name.substr(index + 1);
      let openUrl = "";
      let fileUrl = encodeURIComponent(download_link);
      const fileType = ["docx", "doc", "txt", "pdf", "xlsx"];
      if (fileType.includes(ext)) {
        switch (ext) {
          case "docx" || "doc":
            openUrl = `${window.location.origin}/aibase/doc?fileUrl=` + fileUrl;
            break;
          case "txt":
            openUrl =
              `${window.location.origin}/aibase/txtView?fileUrl=` + fileUrl;
            break;
          case "pdf":
            if (page_num.length > 0) {
              openUrl =
                `${window.location.origin}/aibase/pdfView?fileUrl=` +
                fileUrl +
                "&page=" +
                page_num[0];
            }
            break;
          case "xlsx":
            openUrl =
              `${window.location.origin}/aibase/jsExcel?url=` +
              fileUrl +
              "&rownum=" +
              row_num +
              "&sheetName=" +
              sheet_name;
            break;
          default:
            this.$message.warning("暂不支持此格式查看");
        }
      }
      if (openUrl !== "") {
        window.open(openUrl, "_blank","noopener,noreferrer");
      } else {
        this.$message.warning("暂不支持此格式查看");
      }
    },
    listenerImg() {
      //捕获图片加载错误
      this.imageErrorHandler = (e) => {
        if (e.target.tagName === "IMG") {
          this.handleImageError(e.target);
        }
      };
      document.body.addEventListener("error", this.imageErrorHandler, true);
    },
    handleImageError(img) {
      // 防止重复处理
      if (img.classList.contains("failed")) {
        return;
      }
      img.classList.add("failed");

      // // 设置图片为不可见，避免闪烁
      img.style.visibility = "hidden";
      img.style.display = "none";
    },
    setupScrollListener() {
      const container = document.getElementById("timeScroll");
      container.addEventListener("scroll", this.handleScroll);
    },
    handleScroll(e) {
      const container = document.getElementById("timeScroll");
      const { scrollTop, clientHeight, scrollHeight } = container;
      // 检测是否接近底部（5px容差）
      const nearBottom = scrollHeight - (scrollTop + clientHeight) < 5;
      // 用户手动滚动时取消自动置底
      if (!nearBottom) {
        this.autoScroll = false;
      }
      // 清除之前的定时器
      clearTimeout(this.scrollTimeout);
      // 设置新的定时器检测滚动停止
      this.scrollTimeout = setTimeout(() => {
        // 如果停止时接近底部，恢复自动置底
        if (nearBottom) {
          this.autoScroll = true;
          this.scrollBottom();
        }
      }, 500); // 500ms内没有新滚动视为停止
    },
    getTitle(type){
      if(type === 'qa_start'){
        return this.$t('app.qaSearching')
      }else if(type === 'knowledge_start'){
        return this.$t('app.knowledgeSearch')
      }else if(type === 'qa_finish'){
        return this.$t('knowledgeManage.qaDatabase.name')
      }else{
        return this.$t('menu.knowledge')
      }
    },
    replaceHTML(data, n) {
      let _data = data;
      var a = new RegExp("<think>");
      var b = new RegExp("</think>");
      if (b.test(data)) {
        n.thinkText = this.$t('agent.alreadyThink');
      }
      // 如果没有返回前缀，则补上
      if (b.test(data) && !a.test(data)) {
        _data = "<think>\n" + data;
      }
      return _data.replace(/think>/g, "section>");
    },
    showDSBtn(data) {
      const pattern = /<\/?think>/;
      const matches = data.match(pattern);
      if (!matches) {
        return false;
      }
      return true;
    },
    toggle(event, index) {
      const name = event.target.className;
      if (
        name === "deepseek" ||
        name === "el-icon-arrow-up" ||
        name === "el-icon-arrow-down"
      ) {
        this.session_data.history[index].isOpen =
          !this.session_data.history[index].isOpen;
        this.$set(
          this.session_data.history,
          index,
          this.session_data.history[index]
        );
        let elm = null;
        if (name === "el-icon-arrow-up" || name === "el-icon-arrow-down") {
          elm = event.target.parentNode.parentNode
            .getElementsByClassName("answer-content")[0]
            .getElementsByTagName("section")[0];
        } else {
          elm = event.target.parentNode
            .getElementsByClassName("answer-content")[0]
            .getElementsByTagName("section")[0];
        }
        if (!Boolean(this.session_data.history[index].isOpen)) {
          elm.className = "hideDs";
        } else {
          elm.className = "";
        }
      }
    },
    queryCopy(text) {
      this.$emit("queryCopy", text);
    },
    copy(text) {
      text = text.replaceAll("<br/>", "\n");
      var textareaEl = document.createElement("textarea");
      textareaEl.setAttribute("readonly", "readonly"); // 防止手机上弹出软键盘
      textareaEl.value = text;
      document.body.appendChild(textareaEl);
      textareaEl.select();
      var res = document.execCommand("copy");
      document.body.removeChild(textareaEl);
      return res;
    },
    copycb() {
      this.$message.success(this.$t('agent.copyTips'));
    },
    collapseClick(n, m, j) {
      if (!m.collapse) {
        this.$set(n.searchList, j, { ...m, collapse: true });
      } else {
        this.$set(n.searchList, j, { ...m, collapse: false });
      }
    },
    doLoading() {
      this.loading = true;
    },
    scrollBottom() {
      if (!this.autoScroll) return;
      this.$nextTick(() => {
        this.loading = false;
        document.getElementById("timeScroll").scrollTop =
          document.getElementById("timeScroll").scrollHeight;
      });
    },
    pushHistory(data) {
      this.session_data.history.push(data);
      this.scrollBottom();
    },
    replaceLastData(index, data) {
      if (!data.response && data.finish === 1 ) {
        data.response = this.$t('app.noResponse');
      }
      this.scrollBottom();
      this.$set(this.session_data.history, index, data);
      if (data.finish === 1) {
        const setCitations = this.setCitations(index);
        this.$set(this.session_data.history[index], "citations", setCitations);
      }
    },
    getFileSizeDisplay(fileSize) {
      if (!fileSize || typeof fileSize !== "number" || isNaN(fileSize)) {
        return "...";
      }
      return fileSize > 1024
        ? `${(fileSize / (1024 * 1024)).toFixed(2)} MB`
        : `${fileSize} bytes`;
    },
    //websocket 替换全部数据
    replaceData(data) {
      this.session_data = data;
      this.scrollBottom();
    },
    //http 只替换history
    replaceHistory(data) {
      this.session_data.history = data;
      this.scrollBottom();
      //this.loadAllImg()
    },
    replaceHistoryWithImg(data) {
      this.session_data.history = data;
      this.$nextTick(() => {
        this.preTagging(data[0].annotation);
      });
    },
    clearData() {
      this.session_data = {
        tool: "",
        searchList: [],
        history: [],
        response: "",
      };
    },
    loadAllImg() {
      this.session_data.history.forEach((n, i) => {
        n.gen_file_url_list.forEach((m, j) => {
          setTimeout(() => {
            this.$set(this.session_data.history[i].gen_file_url_list, j, {
              ...m,
              loadedUrl: m.url,
              loading: false,
            });
          }, 2000);
        });
      });
    },
    gropdownClick() {
      this.$emit("clearHistory");
    },
    getSessionData() {
      return this.session_data;
    },
    getList() {
      return JSON.parse(
        JSON.stringify(
          this.session_data.history.filter((item) => {
            delete item.operation;
            return item;
          })
        )
      );
      // return JSON.parse(JSON.stringify(this.session_data.history.filter((item)=>{ delete item.operation ; return !item.pending})))
    },
    getAllList() {
      return JSON.parse(JSON.stringify(this.session_data.history));
    },
    stopLoading() {
      this.session_data.history = this.session_data.history.filter((item) => {
        return !item.pending;
      });
    },
    stopPending() {
      // this.session_data.history = this.session_data.history.filter(item =>{
      this.session_data.history = this.session_data.history.map((item) => {
        if (item.pending) {
          return {
            ...item,
            responseLoading: false,
            pendingResponse: this.$t('app.stopStream'),
            pending: false, // 标记为已完成
            finish: 1,
          };
        }
        return item;
      });
    },
    refresh() {
      if (this.sessionStatus === 0) {
        return;
      }
      this.$emit("refresh");
    },
    preZan(index, item) {
      if (this.sessionStatus === 0) {
        return;
      }
      this.$set(this.session_data.history, index, { ...item, evaluate: 1 });
    },
    preCai(index, item) {
      if (this.sessionStatus === 0) {
        return;
      }
      this.$set(this.session_data.history, index, { ...item, evaluate: 2 });
    },
    doScore(index, evaluate) {},

    //=================标注相关===============
    initCanvasUtil() {
      this.canvasShow = true;
      this.$nextTick(() => {
        // 开始画图 canvas, 2d, 宽高，形状
        this.cv &&
          this.cv.destroy() &&
          this.cv.clearPre() &&
          this.cv.clearLabels() &&
          (this.cv = null);
        this.cv = new CanvasUtil(this);
      });
    },
    preTagging(response) {
      // canvas大小重置
      this.currImg = {
        url: "",
        width: 0,
        height: 0,
        w: 0,
        h: 358,
        roteX: 0,
        roteY: 0,
        dx: 0,
        dy: 0,
      };
      // 图片原始宽高
      var image = new Image();
      image.src = response.annotationImg;
      image.onload = () => {
        this.currImg.width = image.width;
        this.currImg.height = image.height;
        //if (!this.c) {
        this.c = document.getElementById("mycanvas");
        this.ctx = this.c.getContext("2d");
        //}
        this.resizeCanvas();
        this.initCanvasUtil();

        this.$nextTick(() => {
          this.echoLabels(response);
        });
      };
    },
    echoLabels(response) {
      this.cv.echoLabels(response);
    },
    resizeCanvas() {
      this.currImg.w = 0;
      this.currImg.h = 358;
      this.currImg.dx = 0;
      this.currImg.dy = 0;
      this.currImg.roteX = 0;
      this.currImg.roteY = 0;

      let currImg = this.currImg;
      let contain = document.getElementById("mycantain");
      if (currImg.width > contain.offsetWidth) {
        // 宽度大于容器
        this.currImg.roteX = currImg.width / contain.offsetWidth;
        currImg.w = contain.offsetWidth;
        currImg.h = (currImg.height * contain.offsetWidth) / currImg.width;
        // 压缩后高度大于cantain
        if (currImg.h > contain.offsetHeight) {
          currImg.h = contain.offsetHeight;
          currImg.w = (currImg.width * currImg.h) / currImg.height;
          currImg.roteX = currImg.width / currImg.w;
          currImg.dx = (contain.offsetWidth - currImg.w) / 2;
        } else {
          currImg.roteY = currImg.height / currImg.h;
          currImg.dy = (contain.offsetHeight - currImg.h) / 2;
        }
      } else {
        // 高度压缩比例
        currImg.roteY = currImg.height / currImg.h;
        // 压缩后宽度
        currImg.w = (currImg.width * currImg.h) / currImg.height;
        currImg.roteX = currImg.width / currImg.w;
        currImg.dx = (contain.offsetWidth - currImg.w) / 2;
      }

      this.canvasShow = true;
      this.c.width = currImg.w;
      this.c.height = currImg.h;
      this.$nextTick(() => {
        this.cv && this.cv.resizeCurrImg(currImg);
      });
    },
  },
};
</script>

<style scoped lang="scss">
.serach-list-item {
  .link:hover {
    color: $color !important;
  }
  .search-doc {
    margin-left: 10px;
    cursor: pointer;
    color: $color !important;
  }
  .subTag {
    display: inline-flex;
    color: $color;
    border-radius: 50%;
    width: 18px;
    height: 18px;
    border: 1px solid $color;
    line-height: 18px;
    vertical-align: middle;
    margin-left: 2px;
    justify-content: center;
    align-items: center;
    font-size: 14px;
    overflow: hidden;
    white-space: nowrap;
    margin-bottom: 2px;
    transform: scale(0.8);
  }
}

/* 图片加载失败时的样式 */
img.failed {
  position: relative;
  border: 2px dashed #ff6b6b;
  background-color: #fff5f5;
  opacity: 0.5;
}

img.failed::after {
  content: "图片加载失败";
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #ff6b6b;
  font-size: 12px;
  background: rgba(255, 255, 255, 0.9);
  padding: 4px 8px;
  border-radius: 4px;
  white-space: nowrap;
}

/deep/ {
  pre {
    white-space: pre-wrap !important;
  }
  .answer-content {
    margin-top:5px !important;
    img {
      // height:100px;
      width: 100%;
      display: block;
    }
    section li {
      list-style-position: inside; /* 将标记符号放在内容框内 */
    }
    .citation {
      display: inline-flex;
      color: $color;
      border-radius: 50%;
      width: 18px;
      height: 18px;
      border: 1px solid $color;
      cursor: pointer;
      line-height: 18px;
      vertical-align: middle;
      margin-left: 2px;
      justify-content: center;
      align-items: center;
      font-size: 14px;
      overflow: hidden;
      white-space: nowrap;
      margin-bottom: 2px;
      transform: scale(0.8);
    }
    .citation-active {
      background: $color;
      color: #fff;
    }
  }
  .search-list {
    img {
      width: 80% !important;
    }
  }
}
.session {
  word-break: break-all;
  height: 100%;
  overflow-y: auto;
  .session-item {
    min-height: 80px;
    display: flex;
    padding: 20px;
    line-height: 28px;
    img {
      width: 30px;
      height: 30px;
      object-fit: cover;
    }
    .logo {
      border-radius: 6px;
    }
    .answer-content {
      width:100%;
      position: relative;
      margin-left: 14px;
      color: #333;
      .answer-content-query {
        width:100%;
        display: flex;
        flex-wrap: wrap;
        flex-direction: column;
        align-items: flex-end;
        .answer-text {
          background: #7288fa;
          color: #fff;
          border-radius: 10px 0 10px 10px;
          padding: 10px 10px 10px 20px;
          margin:0!important;
          display:inline-block;
          line-height:1.5;
        }
        .session-setting-id {
          color: rgba(98, 98, 98, 0.5);
          font-size: 12px;
          margin-top: -8px;
        }
        .echo-doc-box {
          margin-bottom: 10px;
          background: #fff;
          width: auto;
          border: 1px solid #dcdfe6;
          border-radius: 5px;
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 2px 20px 5px 5px;
          .docIcon {
            width: 30px;
            height: 30px;
            margin-right: 10px;
          }
          .docInfo {
            .docInfo_name {
              color: #333;
            }
            .docInfo_size {
              color: #bbbbbb;
            }
          }
        }
      }
      li {
        display: revert !important;
      }
    }
  }
  .session-answer {
    // background-color: #eceefe;
    border-radius: 10px;
    
    .session-answer-wrapper {
      display: flex;
      align-items: flex-start;
      gap: 10px; /* 头像和内容之间10px距离 */
      padding: 20px 20px 0 20px;
      min-height: 80px;
      background: none; /* 确保外层容器无背景色 */
      
      .logo {
        width: 30px;
        height: 30px;
        border-radius: 6px;
        object-fit: cover;
        flex-shrink: 0; /* 防止头像被压缩 */
        background: none; /* 头像无背景色 */
      }
      
      .answer-content {
        flex: 1;
        background-color: #eceefe; /* 只有内容区域有背景色 */
        border-radius: 0 10px 10px 10px;
        padding: 20px;
        line-height: 1.6;
      }
    }
  }
  
  /* 问题在右侧，答案在左侧 */
  .session-question {
    .session-item {
      flex-direction: row-reverse;
      margin-left: auto;
      width: auto;
      gap: 10px; /* 问题和问题图标之间10px间距 */
      display: flex;
      align-items: flex-start;
    }
  }
  
  .session-answer {
    .answer-annotation {
      line-height: 0 !important;
      .annotation-img {
        width: 460px;
        object-fit: contain;
        height: 358px;
      }
      .tagging-canvas {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        margin: auto;
      }
    }

    .no-response {
      margin: 15px 0;
    }
    /*出处*/
    .search-list {
      padding: 10px 20px 3px 0;
      .qa_content{
        display: flex;
        gap: 10px;
        margin-top:5px;
      }
      .recommended-question-title{
        border-bottom: 1px solid #e5e5e5;
        padding: 5px 0;
      }
      .search-list-item {
        margin-bottom: 5px;
        line-height: 22px;
        p:nth-child(1) {
          white-space: normal;
        }
        a,
        span {
          color: #666;
          cursor: pointer;
          white-space: normal;
          overflow-wrap: break-word;
        }
        a {
          text-decoration: underline;
        }
        a:hover {
          color: deepskyblue;
        }
        .snippet {
          padding: 5px 14px;
        }
      }
    }
    /*操作*/
    .answer-operation {
      display: flex;
      justify-content: space-between;
      padding: 15px 20px 15px 53px;
      color: #777;
      .opera-left {
        flex: 8;
        .restart {
          cursor: pointer;
        }
      }
      .opera-right {
        flex: 1;
        display: inline-flex;
        img {
          width: 20px;
          height: 20px;
          padding: 2px;
        }
        .split-icon {
          background: rgba(195, 197, 217, 0.65);
          height: 22px;
          margin: 0 10px;
          width: 1px;
        }
        .copy-icon {
          font-size: 17px;
          padding: 3px 6px;
          margin: 0 15px;
          cursor: pointer;
        }
        .copy-icon:hover {
          color: #33a4df;
        }
      }
    }
  }

  /*图片*/
  .file-path {
    .el-image {
      height: 200px !important;
      background-color: #f9f9f9;
      /deep/.el-image__inner,
      img {
        width: 100%;
        height: 100%;
        object-fit: contain;
      }
    }
    audio {
      width: 300px !important;
    }
  }
  .query-file {
    padding: 10px 0;
  }
  .response-file {
    margin: 0 0 0 66px;
    width: 400px;
    font-size: 0;
    .img {
      display: inline-block;
      width: 200px;
      height: 200px;
      img {
        width: 100%;
        height: 100%;
      }
    }
  }

  .session-error {
    background-color: #fef0f0;
    border-color: #fde2e2;
    color: #f56c6c !important;
    margin-top: 10px;
    padding: 10px;
    border-radius: 4px;
    .el-icon-warning {
      font-size: 16px;
    }
  }

  .history-box {
    height: calc(100% - 46px);
    overflow-y: auto;
    padding: 20px;
  }
  /*删除历史...*/
  .session-setting {
    position: relative;
    height: 36px;
    .right-setting {
      position: absolute;
      right: 10px;
      top: -5px;
      color: #ff2324;
      font-size: 16px;
      cursor: pointer;
      /deep/ {
        .el-dropdown-menu {
          width: 100px;
        }
        .el-dropdown-menu__item {
          padding: 0 15px !important;
        }
      }
    }
  }

  .think_icon {
    width: 12px !important;
    height: 12px !important;
    margin-right: 3px;
  }
  .ds-res {
    /deep/ section {
      color: #8b8b8b;
      position: relative;
      font-size: 12px;
      * {
        font-size: 12px;
      }
    }
    /deep/ section::before {
      content: "";
      position: absolute;
      height: 100%;
      width: 1px;
      background: #ddd;
      left: -8px;
    }
    /deep/.hideDs {
      display: none;
    }
  }

  .deepseek {
    font-size: 13px;
    color: #8b8b8b;
    font-weight: bold;
    margin-left: 6px;
    cursor: pointer;
  }
}

.text-loading,
.text-loading > div {
  position: relative;
  box-sizing: border-box;
}

.text-loading {
  display: block;
  font-size: 0;
  color: #c8c8c8;
}

.text-loading.la-dark {
  color: #e8e8e8;
}

.text-loading > div {
  display: inline-block;
  float: none;
  background-color: currentColor;
  border: 0 solid currentColor;
}

.text-loading {
  width: 54px;
  height: 18px;
  margin-top: 6px;
}

.text-loading > div {
  width: 8px;
  height: 8px;
  margin: 4px;
  border-radius: 100%;
  animation: ball-beat 0.7s -0.15s infinite linear;
}

.text-loading > div:nth-child(2n-1) {
  animation-delay: -0.5s;
}

@keyframes ball-beat {
  50% {
    opacity: 0.2;
    transform: scale(0.75);
  }

  100% {
    opacity: 1;
    transform: scale(1);
  }
}
</style>
