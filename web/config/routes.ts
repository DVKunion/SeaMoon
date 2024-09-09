export default [
  {
    path: '/user',
    layout: false,
    routes: [
      {
        name: 'login',
        path: '/user/login',
        component: './user/',
      },
      {
        component: './404',
      },
    ],
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    icon: 'dashboard',
    component: './dashboard/',
  },
  {
    path: '/service',
    name: 'service',
    icon: 'Thunderbolt',
    routes: [
      {
        path: "/service",
        component: "./service",
      },
      {
        name: 'network',
        path: '/service/network',
        component: './service/network',
      },
      {
        name: 'application',
        path: '/service/application',
        component: './service/application',
      },
      {
        component: './404',
      },
    ],
  },
  {
    path: '/function',
    name: 'function',
    icon: 'cluster',
    routes: [
      {
        path: "/function",
        component: "./function",
      },
      {
        name: 'network',
        path: '/function/network',
        component: './function/network',
      },
      {
        name: 'application',
        path: '/function/application',
        component: './function/application',
      },
      {
        component: './404',
      },
    ],
  },
  {
    path: 'account',
    name: 'account',
    icon: 'user',
    component: './account/',
  },
  {
    path: '/setting',
    name: 'setting',
    icon: 'setting',
    component: './setting/',
  },
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    component: './404',
  },
];
