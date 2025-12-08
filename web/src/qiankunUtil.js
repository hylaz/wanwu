import { registerMicroApps } from 'qiankun';
import { store } from '@/store';
import { basePath } from '@/utils/config';

export const useQiankun = () => {
  const apps = [
    {
      name: 'test', // 子应用名称
      entry:
        process.env.NODE_ENV === 'development'
          ? 'http://localhost:8081/'
          : window.location.origin + basePath + '/sub/test/', // 子应用入口
      container: '#container', // 子应用所在容器
      props: () => ({ user: store.state.user }), // 传参给子应用
      activeRule: basePath + '/aibase/portal/test', // 子应用触发规则（路径）
    },
  ];

  registerMicroApps(apps, {
    beforeLoad: [
      app => {
        console.log(`${app.name}的beforeLoad阶段`);
      },
    ],
    beforeMount: [
      app => {
        console.log(`${app.name}的beforeMount阶段`);
      },
    ],
    afterMount: [
      app => {
        console.log(`${app.name}的afterMount阶段`);
      },
    ],
    beforeUnmount: [
      app => {
        console.log(`${app.name}的beforeUnmount阶段`);
      },
    ],
    afterUnmount: [
      app => {
        console.log(`${app.name}的afterUnmount阶段`);
      },
    ],
  });
};
