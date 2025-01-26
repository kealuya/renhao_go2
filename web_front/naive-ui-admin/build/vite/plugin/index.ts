import type { Plugin, PluginOption } from 'vite';
import Components from 'unplugin-vue-components/vite';
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers';

import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';

import { configHtmlPlugin } from './html';
import { configCompressPlugin } from './compress';

export function createVitePlugins(viteEnv: ViteEnv, isBuild: boolean) {
  // 从环境变量中解构出压缩相关的配置
  const { VITE_BUILD_COMPRESS, VITE_BUILD_COMPRESS_DELETE_ORIGIN_FILE } = viteEnv;

  // 定义 Vite 插件数组，包含基础必需插件
  const vitePlugins: (Plugin | Plugin[] | PluginOption[])[] = [
    // Vue 3 支持
    vue(),
    // Vue JSX 支持
    vueJsx(),

    // Naive UI 组件按需引入插件
    Components({
      dts: true,  // 自动生成组件的类型声明文件
      resolvers: [NaiveUiResolver()],  // 使用 Naive UI 的解析器
    }),
  ];

  // 添加 HTML 处理插件，用于处理 index.html
  vitePlugins.push(configHtmlPlugin(viteEnv, isBuild));

  // 仅在构建模式下添加压缩插件
  if (isBuild) {
    // 添加 Gzip 压缩插件
    vitePlugins.push(
      configCompressPlugin(VITE_BUILD_COMPRESS, VITE_BUILD_COMPRESS_DELETE_ORIGIN_FILE)
    );
  }

  return vitePlugins;
}
