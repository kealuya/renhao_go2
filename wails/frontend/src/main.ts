import {createApp} from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import './style.css';
import router, {setupRouter} from "./router";


async function bootstrap() {
    const app = createApp(App)
    app.use(ElementPlus)
    // 挂载路由
    setupRouter(app)
    await router.isReady()
    // 路由准备就绪后挂载APP实例
    app.mount('#app', true)
}

void bootstrap()