<template>
  <div class="page-wrapper">
    <div class="page-title">
      <img class="page-title-img" src="@/assets/imgs/knowledge.svg" alt=""/>
      <span class="page-title-name">{{ $t("knowledgeManage.knowledge") }}</span>
    </div>
    <div style="padding: 20px">
      <div class="knowledge-tabs">
        <div
          :class="['knowledge-tab', { active: category === 0 }]"
          @click="tabClick(0)"
        >
          {{ $t("menu.knowledge") }}
        </div>
        <div
          :class="['knowledge-tab', { active: category === 1 }]"
          @click="tabClick(1)"
        >
          {{ $t("knowledgeManage.qaDatabase.name") }}
        </div>
      </div>
      <div class="search-box">
        <div class="no-border-input">
          <search-input
            class="cover-input-icon"
            :placeholder="
              category === 0
                ? $t('knowledgeManage.searchPlaceholder')
                : $t('knowledgeManage.searchPlaceholderQa')
            "
            ref="searchInput"
            @handleSearch="getTableData"
          />
          <el-select
            v-model="tagIds"
            :placeholder="$t('knowledgeManage.selectTag')"
            multiple
            @visible-change="tagChange"
            @remove-tag="removeTag"
            v-if="category === 0"
          >
            <el-option
              v-for="item in tagOptions"
              :key="item.tagId"
              :label="item.tagName"
              :value="item.tagId"
            >
            </el-option>
          </el-select>
        </div>
        <div>
          <el-button
            size="mini"
            type="primary"
            @click="$router.push('/knowledge/keyword')"
            v-if="category === 0"
          >
            {{ $t("knowledgeManage.keyWordManage") }}
          </el-button>
          <el-button
            size="mini"
            type="primary"
            @click="showCreate()"
            icon="el-icon-plus"
          >
            {{ $t("common.button.create") }}
          </el-button>
        </div>
      </div>
      <knowledgeList
        :appData="knowledgeData"
        @editItem="editItem"
        @reloadData="tabClick"
        ref="knowledgeList"
        v-loading="tableLoading"
        :category="category"
      />
      <createKnowledge
        ref="createKnowledge"
        @reloadData="tabClick"
        :category="category"
      />
    </div>
  </div>
</template>
<script>
import {getKnowledgeList, tagList} from "@/api/knowledge";
import SearchInput from "@/components/searchInput.vue";
import knowledgeList from "./component/knowledgeList.vue";
import createKnowledge from "./component/create.vue";

export default {
  components: {SearchInput, knowledgeList, createKnowledge},
  provide() {
    return {
      reloadKnowledgeData: this.getTableData,
    };
  },
  data() {
    return {
      knowledgeData: [],
      tableLoading: false,
      tagOptions: [],
      tagIds: [],
      category: 0,
    };
  },
  beforeRouteEnter(to, from, next) {
    next(vm => {
      vm.handleRouteFrom(from);
    });
  },
  mounted() {
    this.getTableData();
    this.getList();
  },
  methods: {
    handleRouteFrom(from) {
      if (from.path.includes('/qa/docList')) {
        this.category = 1
      } else {
        this.category = 0
      }
    },
    tabClick(status) {
      this.category = status;
      this.getTableData();
    },
    getList() {
      tagList({knowledgeId: "", tagName: ""}).then((res) => {
        if (res.code === 0) {
          this.tagOptions = res.data.knowledgeTagList || [];
        }
      });
    },
    tagChange(val) {
      if (!val && this.tagIds.length > 0) {
        this.getTableData();
      } else {
        this.getList();
      }
    },
    removeTag() {
      this.getTableData();
    },
    getTableData() {
      const searchInput = this.$refs.searchInput.value;
      this.tableLoading = true;
      getKnowledgeList({
        name: searchInput,
        tagId: this.tagIds,
        category: this.category,
      })
        .then((res) => {
          this.knowledgeData = res.data.knowledgeList || [];
          this.tableLoading = false;
        })
        .catch((error) => {
          this.tableLoading = false;
          this.$message.error(error);
        });
    },
    clearIptValue() {
      this.$refs.searchInput.clearValue();
    },
    editItem(row) {
      this.$refs.createKnowledge.showDialog(row);
    },
    showCreate(row) {
      this.$refs.createKnowledge.showDialog(row);
    },
  },
};
</script>
<style lang="scss" scoped>
.search-box {
  display: flex;
  justify-content: space-between;
}

/deep/ {
  .el-loading-mask {
    background: none !important;
  }
}

.active {
  background: #333;
  color: #fff;
  font-weight: bold;
}

.knowledge-tabs {
  margin-bottom: 20px;

  .knowledge-tab {
    display: inline-block;
    vertical-align: middle;
    width: 160px;
    height: 40px;
    border-bottom: 1px solid #333;
    line-height: 40px;
    text-align: center;
    cursor: pointer;
  }
}
</style>
