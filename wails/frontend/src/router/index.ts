import type {App} from 'vue'
import {createRouter, createWebHashHistory} from 'vue-router'
import routeModuleList from './modules'


const router = createRouter({
    history: createWebHashHistory(''),
    routes: routeModuleList,
    strict: true,
    scrollBehavior: () => ({left: 0, top: 0}),
})

export function setupRouter(app: App) {
    app.use(router)
}

export default router
