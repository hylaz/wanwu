import request from "@/utils/request";
import {SERVICE_API} from "@/utils/requestConstants"

export const uploadChunks = (data, config) => {//切片上传
  return request({
    url: `${SERVICE_API}/file/upload`,
    method: "post",
    headers: {"Content-Type": "application/x-www-form-urlencoded"},
    data,
    cancelToken: config,
  });
}
export const checkChunks = (data) => {//检测切片
  return request({
    url: `${SERVICE_API}/file/check`,
    method: "get",
    params: data,
  });
}
export const mergeChunks = (data) => {//合并切片
  return request({
    url: `${SERVICE_API}/file/merge`,
    method: "post",
    data
  });
}
export const clearChunks = (data) => {//清除切片
  return request({
    url: `${SERVICE_API}/file/clean`,
    method: "post",
    data
  });
}
export const delfile = (data) => {//清除切片
  return request({
    url: `${SERVICE_API}/file/delete`,
    method: "delete",
    data
  });
}
export const continueChunks = (data) => {//断点续传,获取已经上传的切片
  return request({
    url: `${SERVICE_API}/file/check/chunk/list`,
    method: "get",
    params: data
  });
}

