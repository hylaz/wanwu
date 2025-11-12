import service from "@/utils/request"
import {USER_API} from "@/utils/requestConstants"

// 获取OAuth应用列表
export const fetchOAuthList = (data) => {
    return service({
        url: `${USER_API}/oauth/app/list`,
        method: "get",
        params: data,
    })
}

// 创建OAuth应用
export const createOAuth = (data) => {
    return service({
        url: `${USER_API}/oauth/app`,
        method: "post",
        data: data,
    })
}

// 更新OAuth应用
export const editOAuth = (data) => {
    return service({
        url: `${USER_API}/oauth/app`,
        method: "put",
        data: data,
    })
}

// 删除OAuth应用
export const deleteOAuth = (data) => {
    return service({
        url: `${USER_API}/oauth/app`,
        method: "delete",
        data: data,
    })
}

// 修改OAuth应用状态
export const changeOAuthStatus = (data) => {
    return service({
        url: `${USER_API}/oauth/app/status`,
        method: "put",
        data: data,
    })
}

// 授权码认证
export const codeOAuth = (data) => {
    return service({
        url: `${USER_API}/oauth/code/authorize`,
        method: "get",
        params: data,
    })
}