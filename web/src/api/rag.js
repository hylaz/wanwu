import service from '@/utils/request';
import { USER_API } from '@/utils/requestConstants';

export const getRagInfo = params => {
  return service({
    url: `${USER_API}/appspace/rag`,
    method: 'get',
    params,
  });
};
export const updateRag = data => {
  return service({
    url: `${USER_API}/appspace/rag`,
    method: 'put',
    data,
  });
};
export const updateRagConfig = data => {
  return service({
    url: `${USER_API}/appspace/rag/config`,
    method: 'put',
    data,
  });
};
export const createRag = data => {
  return service({
    url: `${USER_API}/appspace/rag`,
    method: 'post',
    data,
  });
};
export const delRag = data => {
  return service({
    url: `${USER_API}/appspace/rag`,
    method: 'delete',
    data,
  });
};
export const ragChat = data => {
  return service({
    url: `${USER_API}/rag/chat`,
    method: 'post',
    data,
  });
};
