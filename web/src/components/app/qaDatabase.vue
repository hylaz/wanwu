<template>
  <div class="qa-database-container">
    <div class="qa-database-header">
      <div class="header-left">
        <span class="header-icon el-icon-star-off"></span>
        <span class="header-title">
          {{ $t("knowledgeManage.qaDatabase.linkQaDatabase") }}
        </span>
      </div>
      <div class="header-right">
        <span class="common-add" @click="handleAdd">
          <span class="el-icon-plus"></span>
          <span class="handleBtn">{{ $t("knowledgeSelect.add") }}</span>
        </span>
      </div>
    </div>
    <div class="qa-database-content">
      <div class="action-list" v-if="knowledgeList && knowledgeList.length">
        <div
          v-for="(item, index) in knowledgeList"
          :key="index"
          class="action-item"
        >
          <div class="name">
            <span>{{ item.name }}</span>
          </div>
          <div class="bt">
            <el-tooltip
              class="item"
              effect="dark"
              :content="$t('agent.form.metaDataFilter')"
              placement="top-start"
            >
              <span
                class="el-icon-setting del"
                @click="handleSetting(item, index)"
                style="margin-right: 10px"
              ></span>
            </el-tooltip>
            <span
              class="el-icon-delete del"
              @click="handleDelete(item, index)"
            ></span>
          </div>
        </div>
      </div>
    </div>
    <knowledgeSelect ref="knowledgeSelect" />
  </div>
</template>
<script>
import knowledgeSelect from "@/components/knowledgeSelect.vue";
export default {
  name: "QaDatabase",
  components:{knowledgeSelect},
  props: {
    knowledgeList: {
      type: Array,
      default: () => [],
    },
    category:{
      type:Number,
      default:0
    }
  },
  methods: {
    handleAdd() {
      this.$emit("add");
    },
    handleSetting(item, index) {
      this.$emit("setting", item, index);
    },
    handleDelete(item, index) {
      this.$emit("delete", item, index);
    },
  },
};
</script>

<style lang="scss" scoped>
.qa-database-container {
  width: 100%;

  .qa-database-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;

    .header-left {
      display: flex;
      align-items: center;
      gap: 8px;

      .header-icon {
        font-size: 16px;
        color: #999;
      }

      .header-title {
        font-size: 15px;
        font-weight: bold;
        color: #333;
      }
    }

    .header-right {
      .common-add {
        display: flex;
        align-items: center;
        gap: 4px;
        cursor: pointer;
        color: #333;
        font-size: 14px;

        .el-icon-plus {
          font-size: 14px;
        }

        .handleBtn {
          cursor: pointer;
        }

        &:hover {
          color: $color;
        }
      }
    }
  }

  .qa-database-content {
    .action-list {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 10px;
      width: 100%;

      .action-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        background: #fafafa;
        border: 1px solid #e8e8e8;
        border-radius: 8px;
        padding: 12px 16px;
        min-height: 48px;
        box-sizing: border-box;

        .name {
          flex: 1;
          color: #333;
          font-size: 14px;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
          margin-right: 12px;
        }

        .bt {
          display: flex;
          align-items: center;
          justify-content: flex-end;
          flex-shrink: 0;

          .del {
            color: $color;
            font-size: 16px;
            cursor: pointer;
            transition: opacity 0.2s;

            &:hover {
              opacity: 0.7;
            }
          }
        }
      }
    }

    .empty-state {
      text-align: center;
      padding: 40px 0;
      color: #999;
      font-size: 14px;
    }
  }
}
</style>

