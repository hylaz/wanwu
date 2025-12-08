import service from '@/utils/request';
import { USER_API } from '@/utils/requestConstants';

//问答库文档导出
export const qaDocExport = data => {
  return service({
    url: `${USER_API}/knowledge/qa/export`,
    method: 'get',
    params: data,
  });
};

//删除问答库记录
export const delQaRecord = data => {
  return service({
    url: `${USER_API}/knowledge/export/record`,
    method: 'delete',
    data,
  });
};

//获取问答库导出记录列表
export const getQaExportRecordList = data => {
  return service({
    url: `${USER_API}/knowledge/export/record/list`,
    method: 'get',
    params: data,
  });
};

//编辑问答对
export const editQaPair = data => {
  return service({
    url: `${USER_API}/knowledge/qa/pair`,
    method: 'put',
    data,
  });
};

//新增问答对
export const addQaPair = data => {
  return service({
    url: `${USER_API}/knowledge/qa/pair`,
    method: 'post',
    data,
  });
};

//删除问答对
export const delQaPair = data => {
  return service({
    url: `${USER_API}/knowledge/qa/pair`,
    method: 'delete',
    data,
  });
};

//问答库文档导入
export const qaDocImport = data => {
  return service({
    url: `${USER_API}/knowledge/qa/pair/import`,
    method: 'post',
    data,
  });
};

//获取问答对列表
export const getQaPairList = data => {
  return service({
    url: `${USER_API}/knowledge/qa/pair/list`,
    method: 'get',
    params: data,
  });
};

//启动问答对
export const switchQaPair = data => {
  return service({
    url: `${USER_API}/knowledge/qa/pair/switch`,
    method: 'put',
    data,
  });
};

//问答库命中测试
export const qaHitTest = data => {
  return service({
    url: `${USER_API}/knowledge/qa/hit`,
    method: 'post',
    data,
  });
};

//问答库异步上传任务提示
export const qaTips = data => {
  return service({
    url: `${USER_API}/knowledge/qa/pair/import/tip`,
    method: 'get',
    params: data,
  });
};
