<template>
  <div class="page-wrapper">
    <div class="page-title">
      <img
        class="page-title-img"
        :src="
          typeObj[type] ? typeObj[type].img : require('@/assets/imgs/task.png')
        "
        alt=""
      />
      <span class="page-title-name">
        {{ typeObj[type] ? typeObj[type].title : $t('appSpace.title') }}
      </span>
    </div>
    <div class="hide-loading-bg" style="padding: 20px" v-loading="loading">
      <search-input
        :placeholder="$t('appSpace.search')"
        ref="searchInput"
        @handleSearch="handleSearch"
      />
      <div class="workflow-tabs" v-if="[workflow, chat].includes(type)">
        <div
          :class="['workflow-tab', { active: tabActive === workflow }]"
          @click="tabClick(workflow)"
        >
          {{ $t('appSpace.workflow') }}
        </div>
        <div
          :class="['workflow-tab', { active: tabActive === chat }]"
          @click="tabClick(chat)"
        >
          {{ $t('appSpace.chat') }}
        </div>
      </div>
      <div class="header-right">
        <el-button
          size="mini"
          type="primary"
          @click="showImport"
          v-if="[workflow, chat].includes(type)"
        >
          {{ $t('common.button.import') }}
        </el-button>
        <el-button
          size="mini"
          type="primary"
          @click="showCreate"
          icon="el-icon-plus"
        >
          {{ $t('common.button.create') }}
        </el-button>
      </div>
      <AppList
        :type="type"
        :showCreate="showCreate"
        :appData="listData"
        :isShowPublished="true"
        :isShowTool="true"
        @reloadData="getTableData"
      />
      <CreateTotalDialog ref="createTotalDialog" />
      <UploadFileDialog
        @reloadData="getTableData"
        :appType="type"
        :title="$t('appSpace.workflowExport')"
        ref="uploadFileDialog"
      />
    </div>
  </div>
</template>

<script>
import SearchInput from '@/components/searchInput.vue';
import AppList from '@/components/appList.vue';
import CreateTotalDialog from '@/components/createTotalDialog.vue';
import UploadFileDialog from '@/components/uploadFileDialog.vue';
import { getAppSpaceList, agnetTemplateList } from '@/api/appspace';
import { CHAT, WORKFLOW, RAG, AGENT } from '@/utils/commonSet';
import { mapGetters } from 'vuex';
import { fetchPermFirPath } from '@/utils/util';

export default {
  components: { SearchInput, CreateTotalDialog, UploadFileDialog, AppList },
  data() {
    return {
      type: '',
      chat: CHAT,
      workflow: WORKFLOW,
      tabActive: WORKFLOW,
      loading: false,
      listData: [],
      typeObj: {
        [WORKFLOW]: {
          title: this.$t('appSpace.workflow'),
          img: require('@/assets/imgs/workflow_icon.svg'),
        },
        [CHAT]: {
          title: this.$t('appSpace.workflow'),
          img: require('@/assets/imgs/workflow_icon.svg'),
        },
        [RAG]: {
          title: this.$t('appSpace.rag'),
          img: require('@/assets/imgs/rag.svg'),
        },
        [AGENT]: {
          title: this.$t('appSpace.agent'),
          img: require('@/assets/imgs/agent.svg'),
        },
      },
      currentTypeObj: {},
    };
  },
  watch: {
    $route: {
      handler(val) {
        const { type } = val ? val.params || {} : {};
        const { type: flowType } = val.query || {};
        this.type = flowType || type;
        this.tabActive = flowType || WORKFLOW;

        this.listData = [];
        this.$refs.searchInput.value = '';
        this.justifyRenderPage(type);
        this.getTableData();
      },
      // 深度观察监听
      deep: true,
    },
    fromList: {
      handler(val) {
        if (val !== '') {
          this.type = val;
          this.getTableData();
        }
      },
    },
  },
  computed: {
    ...mapGetters('app', ['fromList']),
  },
  mounted() {
    const { type } = this.$route.params || {};
    const { type: flowType } = this.$route.query || {};
    this.type = flowType || type;
    this.tabActive = flowType || WORKFLOW;
    this.justifyRenderPage(type);
    this.getTableData();
  },
  methods: {
    justifyRenderPage(type) {
      if (![WORKFLOW, AGENT, RAG].includes(type)) {
        const { path } = fetchPermFirPath();
        this.$router.push({ path });
      }
    },
    handleSearch() {
      this.getTableData();
    },
    getTableData() {
      this.loading = true;
      const searchInput = this.$refs.searchInput;
      const searchInfo = {
        appType: this.type === 'all' ? '' : this.type,
        ...(searchInput.value && { name: searchInput.value }),
      };
      getAppSpaceList(searchInfo)
        .then(res => {
          this.loading = false;
          this.listData = res.data ? res.data.list || [] : [];
        })
        .catch(() => {
          this.loading = false;
          this.listData = [];
        });
    },
    tabClick(type) {
      this.tabActive = type;
      this.type = type;
      if (type === CHAT) {
        this.$router.replace({ query: { type } });
      } else if (type === WORKFLOW) {
        this.$router.replace({ query: {} });
      }
      this.getTableData();
    },
    showImport() {
      this.$refs.uploadFileDialog.openDialog();
    },
    showCreate() {
      switch (this.type) {
        case AGENT:
          this.$refs.createTotalDialog.showCreateIntelligent();
          break;
        case RAG:
          this.$refs.createTotalDialog.showCreateTxtQues();
          break;
        case CHAT:
          this.$refs.createTotalDialog.showCreateChat();
          break;
        case WORKFLOW:
          this.$refs.createTotalDialog.showCreateWorkflow();
          break;
        default:
          this.$refs.createTotalDialog.openDialog();
          break;
      }
    },
  },
};
</script>
<style lang="scss" scoped>
.header-right {
  display: inline-block;
  float: right;
}
.workflow-tabs {
  margin-top: 20px;
  display: inline-block;
  width: calc(100% - 200px);

  .workflow-tab {
    display: inline-block;
    vertical-align: middle;
    width: 160px;
    height: 40px;
    border-bottom: 1px solid #333;
    line-height: 40px;
    text-align: center;
    cursor: pointer;
  }

  .active {
    background: #333;
    color: #fff;
    font-weight: bold;
  }
}
</style>
