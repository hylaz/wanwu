## HTTP Header
| Header        | 说明      |
| ------------- | --------- |
| Authorization | JWT token |
| X-Language    | 语言Code  |
| X-Org-Id      | 组织ID    |
| X-Client-Id   | 客户端标识|

## HTTP Status
| HTTP Status             | 说明                   |
| ----------------------- | ---------------------- |
| 200, StatusOK           | 请求返回成功           |
| 400, StatusBadRequest   | 请求返回失败，用于业务 |
| 401, StatusUnauthorized | JWT认证失败            |
| 403, StatusForbidden    | 没有权限               |

## 权限-菜单对应表
| 一级权限        | 二级权限  | 三级权限 | 一级菜单 | 二级菜单 | 三级菜单 |
|-------------|-------|------|------|------|------|
| guest       |       |      | 【访客】 |      |      |
| common      |       |      | 【通用】 |      |      |
| permission  |       |      | 权限管理 |      |      |
| permission  | user  |      | 权限管理 | 用户管理 |      |
| permission  | org   |      | 权限管理 | 组织管理 |      |
| permission  | role  |      | 权限管理 | 角色管理 |      |

## `/v1/user/permission`返回用例
```json
{
  "code": 0,
  "data": {
    "orgPermission": {
      "org": {"id": "test-org-id", "name": "test-org-name"},
      "permissions": [
        {"perm": "permission"},
        {"perm": "permission.user"},
        {"perm": "permission.org"},
        {"perm": "permission.role"}
      ]
    }
  },
  "msg": "操作成功"
}
```