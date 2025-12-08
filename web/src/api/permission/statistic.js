import service from '@/utils/request';
import { USER_API } from '@/utils/requestConstants';

// 获取客户端统计数据
export const getData = params => {
  return service({
    url: `${USER_API}/statistic/client`,
    method: 'get',
    params,
  });
};
