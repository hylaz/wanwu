<template>
  <div>
    <el-dialog
      :title="title"
      :visible.sync="dialogVisible"
      width="45%"
      :before-close="handleClose"
    >
      <el-form ref="form" :model="form" label-width="130px" :rules="rules">
        <el-form-item :label="$t('keyword.quesKeyword')" prop="name">
          <el-input v-model="form.name"></el-input>
        </el-form-item>
        <el-form-item :label="$t('keyword.docWord')" prop="alias">
          <el-input v-model="form.alias"></el-input>
        </el-form-item>
        <el-form-item
          :label="$t('keyword.chooseKnowledge')"
          prop="knowledgeBaseIds"
        >
          <el-select
            v-model="form.knowledgeBaseIds"
            :placeholder="$t('common.select.placeholder')"
            multiple
            clearable
            filterable
            style="width: 100%"
            @visible-change="visibleChange($event)"
          >
            <el-option
              v-for="item in knowledgeOptions"
              :key="item.knowledgeId"
              :label="item.name"
              :value="item.knowledgeId"
            >
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">{{
          $t('common.button.cancel')
        }}</el-button>
        <el-button type="primary" @click="submit('form')">{{
          $t('common.button.confirm')
        }}</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
import { getKnowledgeList } from '@/api/knowledge';
import { addKeyWord, editKeyWord, keyWordDetail } from '@/api/keyword';

export default {
  data() {
    return {
      form: {
        name: '',
        alias: '',
        knowledgeBaseIds: [],
      },
      rules: {
        name: [
          {
            required: true,
            message: this.$t('keyword.quesKeywordMsg'),
            trigger: 'blur',
          },
        ],
        alias: [
          {
            required: true,
            message: this.$t('keyword.docWordMsg'),
            trigger: 'blur',
          },
        ],
        knowledgeBaseIds: [
          {
            required: true,
            message: this.$t('keyword.chooseKnowledgeMsg'),
            trigger: 'blur',
          },
        ],
      },
      knowledgeOptions: [],
      title: this.$t('keyword.create'),
      dialogVisible: false,
      id: '',
    };
  },
  created() {
    this.getKnowledgeList();
  },
  methods: {
    submit(formName) {
      this.$refs[formName].validate(valid => {
        if (valid) {
          if (this.id !== '') {
            this.editItem();
          } else {
            this.addItem();
          }
        } else {
          return false;
        }
      });
    },
    editItem() {
      const data = {
        ...this.form,
        id: this.id,
      };
      editKeyWord(data).then(res => {
        if (res.code === 0) {
          this.$message.success('success');
          this.dialogVisible = false;
          this.$parent.updateData();
        }
      });
    },
    addItem() {
      addKeyWord(this.form).then(res => {
        if (res.code === 0) {
          this.$message.success('success');
          this.dialogVisible = false;
          this.$parent.updateData();
        }
      });
    },
    async getKnowledgeList() {
      //获取文档知识分类
      const res = await getKnowledgeList({});
      if (res.code === 0) {
        this.knowledgeOptions = res.data.knowledgeList || [];
      } else {
        this.$message.error(res.message);
      }
    },
    visibleChange(val) {
      if (val) {
        this.getKnowledgeList();
      }
    },
    showDialog(row = null) {
      this.dialogVisible = true;
      if (row !== null) {
        this.title = this.$t('keyword.edit');
        this.id = row.id;
        this.form.name = row.name;
        this.form.alias = row.alias;
        this.form.knowledgeBaseIds = row.knowledgeBaseIds;
      } else {
        this.clearForm();
      }
    },
    clearForm() {
      this.form.name = '';
      this.form.alias = '';
      this.form.knowledgeBaseIds = [];
      this.id = '';
      this.title = this.$t('keyword.add');
      this.$refs.form.clearValidate();
      this.$refs.form.validateField();
    },
    handleClose() {
      this.dialogVisible = false;
    },
  },
};
</script>
