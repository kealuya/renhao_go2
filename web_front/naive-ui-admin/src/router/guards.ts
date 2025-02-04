// 导入必要的类型和依赖
import type { RouteRecordRaw } from 'vue-router';
import { isNavigationFailure, Router } from 'vue-router';
import { useUser } from '@/store/modules/user';
import { useAsyncRoute } from '@/store/modules/asyncRoute';
import { ACCESS_TOKEN } from '@/store/mutation-types';
import { storage } from '@/utils/Storage';
import { PageEnum } from '@/enums/pageEnum';
import { ErrorPageRoute } from '@/router/base';

// 定义登录页路径
const LOGIN_PATH = PageEnum.BASE_LOGIN;

// 白名单路径列表，这些路径可以直接访问而无需权限验证
const whitePathList = [LOGIN_PATH];

export function createRouterGuards(router: Router) {
  // 获取用户和路由store
  const userStore = useUser();
  const asyncRouteStore = useAsyncRoute();

  // 全局前置守卫
  router.beforeEach(async (to, from, next) => {
    // 获取全局loading实例（如果存在）
    const Loading = window['$loading'] || null;
    Loading && Loading.start();

    // 特殊情况处理：从登录页跳转到错误页时，重定向到首页
    if (from.path === LOGIN_PATH && to.name === 'errorPage') {
      next(PageEnum.BASE_HOME);
      return;
    }

    // 白名单路径直接放行
    if (whitePathList.includes(to.path as PageEnum)) {
      next();
      return;
    }

    // 获取存储的token
    const token = storage.get(ACCESS_TOKEN);
    // 没有token的情况
    if (!token) {
      // 如果路由配置了ignoreAuth，则允许无token访问
      if (to.meta.ignoreAuth) {
        next();
        return;
      }
      // 否则重定向到登录页，并携带原目标路径
      const redirectData: { path: string; replace: boolean; query?: Recordable<string> } = {
        path: LOGIN_PATH,
        replace: true,
      };
      if (to.path) {
        redirectData.query = {
          ...redirectData.query,
          redirect: to.path,
        };
        console.debug(redirectData, 'redirectData');
      }
      next(redirectData);
      return;
    }

    // 如果动态路由已添加，直接放行
    if (asyncRouteStore.getIsDynamicRouteAdded) {
      next();
      return;
    }

    // 获取用户信息
    const userInfo = await userStore.getInfo();

    // 根据用户信息生成可访问的路由表
    const routes = await asyncRouteStore.generateRoutes(userInfo);

    // 动态添加可访问路由表
    routes.forEach((item) => {
      router.addRoute(item as unknown as RouteRecordRaw);
    });

    // 确保404页面路由存在
    const isErrorPage = router.getRoutes().findIndex((item) => item.name === ErrorPageRoute.name);
    if (isErrorPage === -1) {
      router.addRoute(ErrorPageRoute as unknown as RouteRecordRaw);
    }

    // 处理重定向
    const redirectPath = (from.query.redirect || to.path) as string;
    const redirect = decodeURIComponent(redirectPath);
    const nextData = to.path === redirect ? { ...to, replace: true } : { path: redirect };
    // 标记动态路由已添加
    asyncRouteStore.setDynamicRouteAdded(true);
    next(nextData);
    Loading && Loading.finish();
  });

  // 全局后置守卫
  router.afterEach((to, _, failure) => {
    // 设置页面标题
    document.title = (to?.meta?.title as string) || document.title;
    
    // 处理导航失败的情况
    if (isNavigationFailure(failure)) {
      //console.log('failed navigation', failure)
    }

    // 处理组件缓存逻辑
    const asyncRouteStore = useAsyncRoute();
    // 获取需要缓存的组件名称列表
    const keepAliveComponents = asyncRouteStore.keepAliveComponents;
    // 获取当前路由组件名称
    const currentComName: any = to.matched.find((item) => item.name == to.name)?.name;

    // 处理组件缓存逻辑
    if (currentComName && !keepAliveComponents.includes(currentComName) && to.meta?.keepAlive) {
      // 需要缓存但未缓存的组件，添加到缓存列表
      keepAliveComponents.push(currentComName);
    } else if (!to.meta?.keepAlive || to.name == 'Redirect') {
      // 不需要缓存或是重定向页面，从缓存列表中移除
      const index = asyncRouteStore.keepAliveComponents.findIndex((name) => name == currentComName);
      if (index != -1) {
        keepAliveComponents.splice(index, 1);
      }
    }
    // 更新缓存组件列表
    asyncRouteStore.setKeepAliveComponents(keepAliveComponents);
    
    // 结束loading
    const Loading = window['$loading'] || null;
    Loading && Loading.finish();
  });

  // 路由错误处理
  router.onError((error) => {
    console.log(error, '路由错误');
  });
}
