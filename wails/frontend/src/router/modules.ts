import type {RouteRecordRaw} from 'vue-router'
import Layout from "../layout/Layout.vue";

const BASE_HOME: string = "/dashboard"

const RootRoute: RouteRecordRaw = {
    path: '/',
    name: 'Root',
    redirect: BASE_HOME,
    meta: {
        title: 'Root',
    },
}


const routeModuleList: Array<RouteRecordRaw> = [
    RootRoute,
    {
        path: '/dashboard',
        name: 'Dashboard',
        redirect: '/dashboard/index',
        component: Layout,
        meta: {
            title: '主控台',
            icon: 'i-simple-icons:atlassian',
        },
        children: [
            {
                path: 'index',
                name: 'HomePage',
                meta: {},
                component: () => import('@/views/Home.vue'),
            },
        ],
    },
    {
        path: '/info',
        name: 'Info',
        redirect: '/info/index',
        component: Layout,
        meta: {
            title: '图表',
            icon: 'i-simple-icons:soundcharts',
        },
        children: [
            {
                path: 'index',
                name: 'MessagePage',
                meta: {
                },
                component: () => import('@/views/info/Info.vue'),
            },
            {
                path: 'index2',
                name: 'MessagePage2',
                meta: {
                },
                component: () => import('@/views/info/Info2.vue'),
            },
        ],
    },

]

export default routeModuleList
