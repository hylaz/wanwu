<template>
  <div class="add-permission-content">
    <div class="content-wrapper" :class="{ 'transfer-mode': transferMode }">
      <div class="left-panel">
        <div class="search-section">
          <span style="margin-top: 3px;">
            <el-checkbox
              v-if="!transferMode"
              v-model="checkAll"
              :indeterminate="isIndeterminate"
              :disabled="treeData.length === 0"
              @change="handleCheckAllChange"
            />
            {{ $t('knowledgeManage.power.all') }}
          </span>
          <el-select
            v-model="selectedOrganization"
            :placeholder="$t('knowledgeManage.power.org')"
            filterable
            clearable
            class="org-select"
            @change="handleOrgChange"
          >
            <el-option
              v-for="org in organizationList"
              :key="org.orgId"
              :label="org.orgName"
              :value="org.orgId"
            >
            </el-option>
          </el-select>
          <el-input
            v-model="searchKeyword"
            :placeholder="$t('knowledgeManage.power.user')"
            class="search-input"
            :disabled="!selectedOrganization"
            @focus="handleInputFocus"
          >
          </el-input>
        </div>

        <div class="selection-tree">
          <el-tree
            :data="treeData"
            :props="treeProps"
            :show-checkbox="!transferMode"
            node-key="id"
            highlight-current
            :default-expand-all="true"
            @check="handleTreeCheck"
            @node-click="handleNodeClick"
            :filter-node-method="filterNode"
            class="permission-tree"
            ref="tree"
          >
            <span class="custom-tree-node" slot-scope="{ node, data }"
                  :class="{ 'selected-node': transferMode && data.type === 'user' && isNodeSelected(data.id) }">
              <span class="node-label">{{ node.label }}</span>
              <span v-if="transferMode && data.type === 'user' && isNodeSelected(data.id)" class="selected-icon">
                <i class="el-icon-check"></i>
              </span>
            </span>
          </el-tree>
        </div>
      </div>

      <div class="right-panel" v-if="!transferMode">
        <div class="permission-section">
          <div class="permission-label">{{ $t('knowledgeManage.power.permission.label') }}</div>
          <el-select v-model="selectedPermission" :placeholder="$t('knowledgeManage.power.permission.placeholder')"
                     class="permission-select">
            <el-option :label="$t('knowledgeManage.power.permission.read')" :value="0"></el-option>
            <el-option :label="$t('knowledgeManage.power.permission.edit')" :value="10"></el-option>
            <el-option :label="$t('knowledgeManage.power.permission.admin')" :value="20"></el-option>
          </el-select>
        </div>

        <div class="selected-users-section">
          <div class="selected-users-list">
            <div
              v-for="orgGroup in groupedSelectedUsers"
              :key="orgGroup.organization"
              class="org-group"
            >
              <div class="org-group-header">
                <span class="org-name">{{ orgGroup.organization }}</span>
                <span class="user-count">({{ orgGroup.users.length }})</span>
              </div>
              <div class="org-users">
                <div
                  v-for="(user, index) in orgGroup.users"
                  :key="orgGroup.organization + '-' + user.id + '-' + index"
                  class="selected-user-item"
                >
                  <span class="user-info">{{ user.name }}</span>
                  <i class="el-icon-close remove-icon" @click="removeSelectedUser(user)"></i>
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
import {getOrgList, getOrgUser} from '@/api/knowledge'

export default {
  name: 'AddPermission',
  props: {
    knowledgeId: {
      type: String,
      default: ''
    },
    transferMode: {
      type: Boolean,
      default: false
    },
    transferData: {
      type: Object,
      default: null
    }
  },
  computed: {
    defaultPermission() {
      return this.transferMode ? this.$t('knowledgeManage.power.permission.admin') : this.$t('knowledgeManage.power.permission.read')
    },
    groupedSelectedUsers() {
      // 按组织分组选中的用户
      const groups = {};
      this.selectedUsers.forEach(user => {
        if (!groups[user.organization]) {
          groups[user.organization] = {
            organization: user.organization,
            users: []
          };
        }
        groups[user.organization].users.push(user);
      });

      // 转换为数组并按组织名称排序
      return Object.values(groups).sort((a, b) => a.organization.localeCompare(b.organization));
    }
  },
  data() {
    return {
      searchKeyword: '',
      selectedOrganization: '',
      selectedPermission: 0,
      organizationList: [],
      originalTreeData: null,
      treeProps: {
        children: 'children',
        label: 'userName'
      },
      treeData: [],
      selectedUsers: [],
      isSettingChecked: false,
      checkAll: false,
      isIndeterminate: false,
    }
  },
  watch: {
    transferMode: {
      handler(newVal) {
        if (newVal) {
          this.selectedPermission = this.$t('knowledgeManage.power.permission.admin')
        }
      },
      immediate: true
    },
    searchKeyword(val) {
      // 只有在选择了组织时才进行搜索
      if (this.selectedOrganization) {
        this.$refs.tree.filter(val);
      }
    }
  },
  created() {
    this.getOrgList()
  },
  methods: {
    getOrgList() {
      getOrgList({knowledgeId: this.knowledgeId, transfer: this.transferMode}).then(res => {
        if (res.code === 0) {
          this.organizationList = res.data.knowOrgInfoList || []
        }
      })
    },
    getResults() {
      return {node: this.groupedSelectedUsers, selectedPermission: this.selectedPermission}
    },
    isNodeSelected(nodeId) {
      return this.selectedUsers.some(user => user.id === nodeId)
    },
    filterNode(value, data) {
      if (!value) return true;
      // 搜索用户名，需要判断属性是否存在
      if (data.userName) {
        return data.userName.indexOf(value) !== -1;
      }
      return false;
    },
    handleCheckAllChange(checkAll) {
      this.isIndeterminate = false;
      if (checkAll) {
        this.$refs.tree.setCheckedNodes(this.treeData);

        // 触发check事件更新selectedUsers
        this.handleTreeCheck(null, {
          checkedNodes: this.treeData,
          checkedKeys: this.treeData.map(node => node.id)
        });
      } else {
        this.$refs.tree.setCheckedKeys([]);

        // 清空当前组织的选中用户
        const currentOrgId = this.selectedOrganization;
        this.selectedUsers = this.selectedUsers.filter(user => user.orgId !== currentOrgId);
      }
    },
    handleOrgChange(orgId) {
      // 清空全选
      this.checkAll = false;
      this.isIndeterminate = false;

      // 当组织选择改变时，过滤树形数据
      this.getOrgUser(orgId);

      // 如果清空了组织选择，同时清空用户名搜索
      if (!orgId) {
        this.searchKeyword = '';
        this.treeData = [];
      }
    },
    getOrgUser(orgId) {
      if (!orgId) return
      var self = this;
      getOrgUser({knowledgeId: this.knowledgeId, orgId}).then(res => {
        if (res.code === 0) {
          var userList = res.data.userInfoList || [];
          var orgIdValue = res.data.orgId;
          // 给每一项添加 orgId 和 id 字段
          self.treeData = userList.map(function (item) {
            item.orgId = orgIdValue;
            item.id = item.userId;  // 确保有 id 字段，与 userId 保持一致，用于树节点的 key
            return item;
          });
          // 加载完数据后，设置当前组织已选中的用户
          self.$nextTick(function () {
            self.setCheckedUsersForCurrentOrg();
          });
        }
      })
    },
    setCheckedUsersForCurrentOrg() {
      var self = this;
      //获取当前组织id
      var currentOrgId = this.selectedOrganization;

      // 确保有组织选择
      if (!currentOrgId) {
        if (this.$refs.tree) {
          this.$refs.tree.setCheckedKeys([]);
        }
        return;
      }

      // 找出当前组织ID下已选中的用户ID列表
      // 必须同时匹配用户的 orgId 和当前选择的组织 ID
      var checkedUserIds = this.selectedUsers
        .filter(function (user) {
          return user.orgId === currentOrgId;
        })
        .map(function (user) {
          // 兼容 id 和 userId 两种字段
          return user.id || user.userId;
        })
        .filter(function (id) {
          return id != null && id !== undefined && id !== '';
        });

      // 设置树形控件的选中状态
      if (this.$refs.tree) {
        // 设置标志位，防止触发 handleTreeCheck
        this.isSettingChecked = true;
        if (checkedUserIds.length > 0) {
          this.$refs.tree.setCheckedKeys(checkedUserIds);
        } else {
          this.$refs.tree.setCheckedKeys([]);
        }
        // 延迟重置标志位
        this.$nextTick(function () {
          self.isSettingChecked = false;
        });
      }
    },
    handleInputFocus() {
      // 当用户名输入框获得焦点时，如果没有选择组织，给出提示
      if (!this.selectedOrganization) {
        this.$message.warning(this.$t('knowledgeManage.power.warning'));
      }
    },
    filterTreeByOrganization(orgId) {
      if (!orgId) {
        this.$refs.tree.filter('');
        return;
      }

      // 应用过滤
      this.treeData = filterData(this.originalTreeData);
    },
    handleTreeCheck(data, checkedInfo) {
      // 如果是程序设置的选中状态，不处理
      if (this.isSettingChecked) {
        return;
      }

      const checkedNodes = checkedInfo.checkedNodes || [];

      // 更新全选状态
      const totalLeafCount = this.treeData.length;
      const checkedLeafCount = checkedNodes.length;

      if (checkedLeafCount === 0) {
        // 无选中
        this.checkAll = false;
        this.isIndeterminate = false;
      } else if (checkedLeafCount === totalLeafCount) {
        // 全选
        this.checkAll = true;
        this.isIndeterminate = false;
      } else {
        // 部分选中
        this.checkAll = false;
        this.isIndeterminate = true;
      }

      const currentOrg = this.organizationList.find(function (org) {
        return org.orgId === this.selectedOrganization;
      }.bind(this));
      const currentOrgId = this.selectedOrganization;
      const currentOrgName = currentOrg ? currentOrg.orgName : '';

      // 先移除当前组织的所有用户
      var otherOrgUsers = this.selectedUsers.filter(function (user) {
        return user.orgId !== currentOrgId;
      });

      // 收集当前选中的用户（去重）
      var currentOrgUsers = [];
      var addedUserIds = {};

      checkedNodes.forEach(function (node) {
        // 使用 id 字段（已经与 userId 保持一致）
        const nodeId = node.id || node.userId;
        if (nodeId && !addedUserIds[nodeId]) {
          addedUserIds[nodeId] = true;
          currentOrgUsers.push({
            id: nodeId,
            name: node.userName || node.name,
            orgId: node.orgId,
            organization: currentOrgName
          });
        }
      });
      // 合并其他组织的用户和当前组织的用户
      var mergedUsers = otherOrgUsers.concat(currentOrgUsers);

      // 最终去重：使用 userId + orgId 作为唯一标识
      var uniqueUsers = [];
      var uniqueKeys = {};

      mergedUsers.forEach(function (user) {
        var key = user.id + '_' + user.orgId;
        if (!uniqueKeys[key]) {
          uniqueKeys[key] = true;
          uniqueUsers.push(user);
        }
      });

      this.selectedUsers = uniqueUsers;
    },
    handleNodeClick(data, node) {
      if (this.transferMode) {
        this.selectedUsers = [{
          userId: data.userId,
          orgId: data.orgId,
        }]
      }
    },
    getTransferData() {
      return {
        knowledgeId: this.knowledgeId,
        knowledgeUser: this.selectedUsers[0] || []
      }
    },
    removeUser(user) {
      user.selected = false
      this.selectedUsers = this.selectedUsers.filter(u => u.id !== user.id)
    },
    removeSelectedUser(user) {
      // 使用 userId 或 id 都能删除
      const userId = user.userId || user.id;
      this.selectedUsers = this.selectedUsers.filter(u => {
        const uId = u.userId || u.id;
        return uId !== userId;
      });

      this.updateTreeSelection(userId, false);

      this.$nextTick(() => {
        if (this.$refs.tree) {
          if (this.transferMode) {
            this.$refs.tree.setCheckedKeys([]);
          } else {
            const checkedKeys = this.$refs.tree.getCheckedKeys();
            const newCheckedKeys = checkedKeys.filter(key => key !== userId);
            this.$refs.tree.setCheckedKeys(newCheckedKeys);
          }
        }
      });
    },
    updateTreeSelection(userId, selected) {
      const updateNode = (nodes) => {
        nodes.forEach(node => {
          // 兼容 id 和 userId 两种字段
          const nodeId = node.id || node.userId;
          if (nodeId === userId) {
            node.selected = selected;
          }
          if (node.children) {
            updateNode(node.children);
          }
        });
      };
      updateNode(this.treeData);
    },
    updateSelectedNodeBackground() {
      this.$nextTick(() => {
        const allContents = document.querySelectorAll('.permission-tree .el-tree-node__content')
        allContents.forEach(content => {
          content.classList.remove('selected-content')
        })

        this.selectedUsers.forEach(user => {
          const nodeContent = document.querySelector(`[data-key="${user.id}"] .el-tree-node__content`)
          if (nodeContent) {
            nodeContent.classList.add('selected-content')
          }
        })
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.add-permission-content {
  background: #fff;
  border-radius: 8px;

  .content-wrapper {
    display: flex;
    gap: 15px;
    height: 400px;

    &.transfer-mode {
      .left-panel {
        flex: 1;
        width: 100%;
      }
    }

    .left-panel {
      flex: 1;
      display: flex;
      flex-direction: column;
      border: 1px solid #e4e7ed;
      border-radius: 4px;
      padding: 15px;

      .search-section {
        margin-bottom: 15px;
        display: flex;
        gap: 10px;

        .org-select {
          flex: 1;
        }

        .search-input {
          flex: 1;
        }
      }

      .selection-tree {
        flex: 1;
        overflow-y: auto;

        .permission-tree {
          /deep/ .el-tree-node__content {
            height: 32px;
            line-height: 32px;

            &:hover {
              background-color: #f5f7fa;
            }

            &.selected-content {
              background-color: #f5f7fa;
            }
          }

          /deep/ .el-checkbox {
            margin-right: 8px;
          }

          .custom-tree-node {
            display: flex;
            align-items: center;
            justify-content: space-between;
            width: 100%;

            .node-label {
              flex: 1;
              font-size: 14px;
              color: #606266;
            }

            .selected-icon {
              color: $color;
              font-size: 16px;
              margin-right: 8px;
            }

            &.selected-node {
              .node-label {
                color: $color;
              }
            }
          }
        }
      }

    }

    .right-panel {
      flex: 1;
      display: flex;
      flex-direction: column;
      border: 1px solid #e4e7ed;
      border-radius: 4px;
      padding: 15px;

      .permission-section {
        margin-bottom: 20px;
        display: flex;
        align-items: center;

        .permission-label {
          font-size: 14px;
          color: #606266;
          margin-right: 10px;
          white-space: nowrap;
        }

        .permission-select {
          flex: 1;
        }
      }

      .selected-users-section {
        flex: 1;

        .selected-users-list {
          max-height: 300px;
          overflow-y: auto;

          .org-group {
            margin-bottom: 16px;

            &:last-child {
              margin-bottom: 0;
            }

            .org-group-header {
              display: flex;
              align-items: center;
              margin-bottom: 8px;
              padding: 4px 0;
              border-bottom: 1px solid #e4e7ed;

              .org-name {
                font-size: 14px;
                font-weight: 600;
                color: $color;
              }

              .user-count {
                font-size: 12px;
                color: #909399;
                margin-left: 8px;
              }
            }

            .org-users {
              .selected-user-item {
                display: flex;
                align-items: center;
                justify-content: space-between;
                padding: 6px 8px;
                cursor: pointer;
                border-radius: 4px;
                background-color: #f5f7fa;
                border: 1px solid transparent;
                transition: all 0.3s ease;
                margin-bottom: 6px;

                &:last-child {
                  margin-bottom: 0;
                }

                &:hover {
                  background-color: #f0f2ff;
                  border-color: $color;
                }

                .user-info {
                  font-size: 13px;
                  color: #606266;
                }

                .remove-icon {
                  color: $color;
                  cursor: pointer;
                  font-size: 12px;
                  padding: 2px;
                  border-radius: 2px;
                  opacity: 0;
                  transition: opacity 0.3s ease;
                }
              }

              .selected-user-item:hover .remove-icon {
                opacity: 1;
              }
            }
          }
        }
      }
    }
  }
}
</style>
