<template>
  <el-dialog
    title="新增分段"
    :visible.sync="dialogVisible"
    width="40%"
    :before-close="handleClose"
  >
    <el-form
      :model="ruleForm"
      ref="ruleForm"
      label-width="100px"
      class="demo-ruleForm"
    >
      <el-form-item
        class="itemCenter"
        v-if="!isChildChunk"
      >
        <el-radio-group
          v-model="createType"
          @input="typeChange($event)"
        >
          <el-radio-button :label="'single'">单条新增</el-radio-button>
          <el-radio-button :label="'file'">批量上传</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item
        label="批量上传"
        v-if="createType === 'file' && !isChildChunk"
        prop="fileUploadId"
        :rules="[{ required: true, message: '请上传文件', trigger: 'blur' }]"
      >
        <fileUpload
          ref="fileUpload"
          :templateUrl="templateUrl"
          @uploadFile="uploadFile"
          :accept="accept"
        />
      </el-form-item>
      <template v-if="createType === 'single' || isChildChunk">
        <el-form-item
          label="分段内容"
          prop="content"
          :rules="[{ required: true, message: '请输入分段内容', trigger: 'blur' }]"
        >
          <el-input
            placeholder="请输入分段内容"
            v-model="ruleForm.content"
            type="textarea"
            :rows="4"
          ></el-input>
        </el-form-item>
        <el-form-item
          label="关键词"
          prop="labels"
          v-if="!isChildChunk"
        >
          <el-tag
            :key="tag"
            v-for="(tag,index) in ruleForm.labels"
            closable
            :disable-transitions="false"
            @close="handleTagClose(index)"
          >
            {{tag}}
          </el-tag>
          <el-input
            class="input-new-tag"
            v-if="inputVisible"
            v-model="inputValue"
            ref="saveTagInput"
            size="small"
            @keyup.enter.native="handleInputConfirm"
            @blur="handleInputConfirm"
          >
          </el-input>
          <el-button
            v-else
            class="button-new-tag"
            size="small"
            @click="showInput"
          >+ 新建关键词</el-button>
        </el-form-item>
        <el-form-item label="新增方式">
          <el-checkbox-group v-model="checkType">
            <el-checkbox
              label="more"
              name="type"
            >连续新增</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </template>
    </el-form>
    <span
      slot="footer"
      class="dialog-footer"
    >
      <el-button @click="dialogVisible = false">取 消</el-button>
      <el-button
        type="primary"
        @click="submit('ruleForm')"
        :loading="btnLoading"
      >确 定</el-button>
    </span>
  </el-dialog>
</template>
<script>
import fileUpload from "@/components/fileUpload";
import {
  createSegment,
  createBatchSegment,
  createSegmentChild,
} from "@/api/knowledge";
export default {
  components: { fileUpload },
  props: {
    parentId: {
      type: String,
      default: ""
    }
  },
  data() {
    return {
      btnLoading: false,
      accept: ".csv",
      checkType: [],
      inputVisible: false,
      inputValue: "",
      createType: "single",
      ruleForm: {
        content: "",
        docId: "",
        labels: [],
        fileUploadId: "",
      },
      dialogVisible: false,
      templateUrl: "/user/api/v1/static/docs/segment.csv",
      isChildChunk: false,
    };
  },
  methods: {
    typeChange(val) {
      if (this.isChildChunk) {
        this.createType = "single";
        return;
      }
      if (val === "single") {
        this.ruleForm.fileUploadId = "";
        this.$refs.fileUpload.clearFileList();
      } else {
        this.clearForm();
        this.$refs.ruleForm.clearValidate();
      }
    },
    uploadFile(fileUploadId) {
      this.ruleForm.fileUploadId = fileUploadId;
    },
    handleClose() {
      this.dialogVisible = false;
    },
    showDiglog(docId, isChildChunk = false) {
      this.dialogVisible = true;
      this.isChildChunk = isChildChunk;
      this.ruleForm.docId = docId;
      this.clearForm();
    },
    showInput() {
      this.inputVisible = true;
      this.$nextTick((_) => {
        this.$refs.saveTagInput.$refs.input.focus();
      });
    },
    handleTagClose(index) {
      this.ruleForm.labels.splice(index, 1);
    },
    handleInputConfirm() {
      if (this.inputValue) {
        this.ruleForm.labels.push(this.inputValue);
        this.inputVisible = false;
        this.inputValue = "";
      } else {
        this.$message.warning("请输入关键词");
      }
    },
    submit(formName) {
      if (this.createType === "single" || this.isChildChunk) {
        this.handleSingle(formName);
      } else {
        this.handleFile();
      }
    },
    handleSingle(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.btnLoading = true;
          if (this.isChildChunk) {
            this.createChildChunk();
          } else {
            this.createParentChunk();
          }
        } else {
          return false;
        }
      });
    },
    createParentChunk() {
      const data = this.isChildChunk
        ? { content: this.ruleForm.content, docId: this.ruleForm.docId }
        : {
            content: this.ruleForm.content,
            docId: this.ruleForm.docId,
            labels: this.ruleForm.labels,
          };
      createSegment(data)
        .then((res) => {
          if (res.code === 0) {
            this.$message.success("创建成功");
            if (!this.checkType.length) {
              this.dialogVisible = false;
              this.$emit("updateDataBatch");
            } else {
              this.clearForm();
              this.$emit("updateData");
            }
            this.btnLoading = false;
          }
        })
        .catch(() => {
          this.btnLoading = false;
        });
    },
    createChildChunk() {
      const data = {
        content: [this.ruleForm.content],
        docId: this.ruleForm.docId,
        parentId: this.parentId
      };
      createSegmentChild(data).then((res) => {
        if (res.code === 0) {
            this.$message.success("创建成功");
            if (!this.checkType.length) {
              this.dialogVisible = false;
            } else {
              this.clearForm();
            }
            this.$emit("updateChildData");
            this.btnLoading = false;
          }
      }).catch(() => {
        this.btnLoading = false;
      });
    },
    handleFile() {
      this.btnLoading = true;
      const data = {
        fileUploadId: this.ruleForm.fileUploadId,
        docId: this.ruleForm.docId,
      };
      createBatchSegment(data)
        .then((res) => {
          if (res.code === 0) {
            this.$message.success("创建成功");
            this.dialogVisible = false;
            this.btnLoading = false;
            this.$emit("updateDataBatch");
          }
        })
        .catch(() => {
          this.btnLoading = false;
        });
    },
    clearForm() {
      this.ruleForm.content = "";
      if (!this.isChildChunk) {
        this.ruleForm.labels = [];
      }
      this.checkType = [];
    },
  },
};
</script>
<style lang="scss" scoped>
.itemCenter {
  display: flex;
  justify-content: center;
  /deep/.el-form-item__content {
    margin-left: 0 !important;
  }
}
.el-tag {
  margin-right: 5px;
  color: #3848f7;
  border-color: #3848f7;
  background: $color_opacity;
}
/deep/ {
  .el-tag .el-tag__close {
    color: #3848f7 !important;
  }
  .el-tag .el-tag__close:hover {
    color: #fff !important;
    background: #3848f7;
  }
  .el-checkbox__input.is-checked + .el-checkbox__label {
    color: #3848f7;
  }
}
</style>
