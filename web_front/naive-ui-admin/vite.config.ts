// 导入必要的类型和函数
import type { UserConfig, ConfigEnv } from 'vite';
import { loadEnv } from 'vite';
import { resolve } from 'path';
import { wrapperEnv } from './build/utils';
import { createVitePlugins } from './build/vite/plugin';
import { OUTPUT_DIR } from './build/constant';
import { createProxy } from './build/vite/proxy';
import pkg from './package.json';
import { format } from 'date-fns';

// 从 package.json 中获取项目信息
const { dependencies, devDependencies, name, version } = pkg;

// 定义应用信息，包含依赖信息和构建时间
const __APP_INFO__ = {
  pkg: { dependencies, devDependencies, name, version },
  lastBuildTime: format(new Date(), 'yyyy-MM-dd HH:mm:ss'),
};

// 路径解析工具函数
function pathResolve(dir: string) {
  return resolve(process.cwd(), '.', dir);
}

// Vite 配置主函数 - 接收命令类型和环境模式作为参数
export default ({ command, mode }: ConfigEnv): UserConfig => {
  // 获取项目根目录的绝对路径
  const root = process.cwd();
  // 根据模式加载对应的环境变量文件（.env、.env.development 等）
  const env = loadEnv(mode, root);
  // 处理环境变量，转换类型并进行必要的处理
  const viteEnv = wrapperEnv(env);
  // 解构出需要的环境变量
  const { VITE_PUBLIC_PATH, VITE_PORT, VITE_PROXY } = viteEnv;
  // 判断当前是否为构建模式（区分开发环境和生产环境）
  const isBuild = command === 'build';

  return {
    // base: 部署时的基础路径，影响资源加载的URL前缀
    base: VITE_PUBLIC_PATH,
    // esbuild: JavaScript/TypeScript 的编译器配置
    esbuild: {},
    // resolve: 路径解析配置
    resolve: {
      alias: [
        // 配置类型文件路径别名，使用 /#/ 可以直接访问 types 目录
        {
          find: /\/#\//,
          replacement: pathResolve('types') + '/',
        },
        // 配置 @ 指向 src 目录，实现简化导入路径
        {
          find: '@',
          replacement: pathResolve('src') + '/',
        },
      ],
      // dedupe: 强制使用单一版本的包，避免同一个包多个版本共存
      dedupe: ['vue'],
    },
    // plugins: Vite 插件配置，用于扩展 Vite 的功能
    plugins: createVitePlugins(viteEnv, isBuild),
    // define: 定义全局常量，在代码中可以直接使用
    define: {
      // 定义环境变量
      __APP_ENV__: JSON.stringify(env.APP_ENV),
      // 定义应用信息，包含依赖版本和构建时间
      __APP_INFO__: JSON.stringify(__APP_INFO__),
      // 关闭 Vue 的水合不匹配警告
      __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: false,
    },
    // server: 开发服务器配置
    server: {
      // host: true 表示监听所有地址，包括局域网和公网地址
      host: true,
      // port: 开发服务器端口号
      port: VITE_PORT,
      // proxy: 代理配置，用于解决跨域问题
      proxy: createProxy(VITE_PROXY),
    },
    // optimizeDeps: 依赖优化选项，用于控制预构建行为
    optimizeDeps: {
      // include: 需要强制预构建的依赖
      include: [],
      // exclude: 排除预构建的依赖
      exclude: ['vue-demi'],
    },
    // build: 构建配置选项
    build: {
      // target: 构建目标的浏览器版本
      target: 'es2015',
      // cssTarget: CSS 的目标浏览器版本
      cssTarget: 'chrome80',
      // outDir: 构建输出目录
      outDir: OUTPUT_DIR,
      // reportCompressedSize: 禁用压缩大小报告，提高构建性能
      reportCompressedSize: false,
      // chunkSizeWarningLimit: 规定触发警告的 chunk 大小（单位：KB）
      chunkSizeWarningLimit: 2000,
    },
  };
};
