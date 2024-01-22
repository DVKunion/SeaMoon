export default [
  {
    path: '/user',
    layout: false,
    routes: [
      {
        name: 'login',
        path: '/user/login',
        component: './user/Login',
      },
      {
        component: './404',
      },
    ],
  },
  {
    path: '/dashboard',
    name: '仪表盘',
    icon: 'dashboard',
    component: './Dashboard',
  },
  {
    path: '/network',
    name: '网络',
    icon: 'cluster',
    routes: [
      {
        path: 'proxy',
        name: '节点',
        component: './network/Proxy',
      },
      {
        component: './404',
      },
    ],
  },
  {
    path: '/application',
    name: '应用',
    icon: 'AppstoreAdd',
    routes: [
      {
        component: './404',
      },
    ],
  },
  {
    path: '/setting',
    name: '配置',
    icon: 'setting',
    routes: [
      {
        path: 'system',
        name: '系统', // 本地代理状态、web相关选项、管理员账户
        component: './setting/System',
      },
      {
        path: 'server',
        name: '服务', // 当前服务运行状态
        component: './setting/Server',
      },
      {
        component: './404',
      },
    ]
  },
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    component: './404',
  },
];
