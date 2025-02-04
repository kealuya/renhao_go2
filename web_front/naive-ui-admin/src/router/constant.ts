/*
    renhao
    特点：
    懒加载组件
    组件会在实际需要时才被加载
    适用于路由配置，可以减少首屏加载时间
    返回的是一个 Promise
*/

export const RedirectName = 'Redirect';

export const ErrorPage = () => import('@/views/exception/404.vue');

export const Layout = () => import('@/layout/index.vue');
// renhao 多级菜单结构：表格和表单模块需要展示二级子菜单， 所以需要使用ParentLayout，作为中间过渡，用来在menu侧边栏展示二级子菜单
export const ParentLayout = () => import('@/layout/parentLayout.vue');
