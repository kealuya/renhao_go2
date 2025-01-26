// 导入必要的依赖和类型
import { App } from 'vue';
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import { RedirectRoute } from '@/router/base';
import { PageEnum } from '@/enums/pageEnum';
import { createRouterGuards } from './guards';
import type { IModuleType } from './types';

// 自动导入路由模块
// 使用 Vite 的 import.meta.glob 特性自动导入 modules 目录下的所有路由模块
const modules = import.meta.glob<IModuleType>('./modules/**/*.ts', { eager: true });

// 处理路由模块，将所有模块的路由配置合并到一个数组中
const routeModuleList: RouteRecordRaw[] = Object.keys(modules).reduce((list, key) => {
  const mod = modules[key].default ?? {};
  const modList = Array.isArray(mod) ? [...mod] : [mod];
  return [...list, ...modList];
}, []);

// 路由排序函数
// 根据路由的 meta.sort 属性对路由进行排序
function sortRoute(a, b) {
  return (a.meta?.sort ?? 0) - (b.meta?.sort ?? 0);
}

routeModuleList.sort(sortRoute);

// 根路由配置
// 定义应用的根路由，默认重定向到首页
export const RootRoute: RouteRecordRaw = {
  path: '/',
  name: 'Root',
  redirect: PageEnum.BASE_HOME,
  meta: {
    title: 'Root',
  },
};

// 登录路由配置
// 定义登录页面的路由，采用懒加载方式
export const LoginRoute: RouteRecordRaw = {
  path: '/login',
  name: 'Login',
  component: () => import('@/views/login/index.vue'),
  meta: {
    title: '登录',
  },
};

// 需要进行权限验证的路由
// 这些路由需要用户登录后根据权限动态加载
export const asyncRoutes = [...routeModuleList];

// 基础路由
// 这些路由无需权限验证，所有用户都可以访问
export const constantRouter: RouteRecordRaw[] = [LoginRoute, RootRoute, RedirectRoute];

// 创建路由实例
const router = createRouter({
  // 使用 HTML5 历史模式
  history: createWebHistory(),
  // 初始化基础路由
  routes: constantRouter,
  // 启用严格模式
  strict: true,
  // 配置滚动行为
  scrollBehavior: () => ({ left: 0, top: 0 }),
});

export function setupRouter(app: App) {
  app.use(router);
  // 创建路由守卫
  createRouterGuards(router);
}

export default router;
