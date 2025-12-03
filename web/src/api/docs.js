import service from "@/utils/request"
import {USER_API} from "@/utils/requestConstants"

// 获取文档中心 md 内容
export const getMarkdown = (params) => {
  return service({
    url: `${USER_API}/doc_center/markdown`,
    method: 'get',
    params
  });
};

// 获取文档中心目录
export const getDocMenu = () => {
  return service({
    url: `${USER_API}/doc_center/menu`,
    method: 'get',
  });
};

// 获取文档搜索内容
export const getDocSearchContent = (params) => {
  return service({
    url: `${USER_API}/doc_center/search`,
    method: 'get',
    params
  });
};