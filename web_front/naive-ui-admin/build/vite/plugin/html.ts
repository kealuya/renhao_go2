/**
 * Vite HTML 插件配置
 * 用于处理 index.html 的模板语法和最小化
 * 插件地址: https://github.com/anncwb/vite-plugin-html
 */
import type { PluginOption } from 'vite';

import { createHtmlPlugin } from 'vite-plugin-html';

import pkg from '../../../package.json';
import { GLOB_CONFIG_FILE_NAME } from '../../constant';

/**
 * 配置 HTML 插件
 * @param env - Vite 环境变量对象
 * @param isBuild - 是否是生产环境构建
 * @returns HTML 插件配置数组
 */
export function configHtmlPlugin(env: ViteEnv, isBuild: boolean) {
  // 从环境变量中解构需要的配置
  const { VITE_GLOB_APP_TITLE, VITE_PUBLIC_PATH } = env;

  // 确保公共路径以 '/' 结尾
  const path = VITE_PUBLIC_PATH.endsWith('/') ? VITE_PUBLIC_PATH : `${VITE_PUBLIC_PATH}/`;

  /**
   * 获取应用配置文件的 URL
   * 添加版本号和时间戳以防止缓存
   * @returns 完整的配置文件 URL
   */
  const getAppConfigSrc = () => {
    return `${path || '/'}${GLOB_CONFIG_FILE_NAME}?v=${pkg.version}-${new Date().getTime()}`;
  };

  /**
   * 创建 HTML 插件配置
   * 配置包括：
   * 1. minify: 是否压缩 HTML（仅在生产环境）
   * 2. inject: 注入数据到 EJS 模板
   *    - data: 注入的数据对象，包含页面标题
   *    - tags: 在生产环境下注入 script 标签，用于加载应用配置
   */
  const htmlPlugin: PluginOption[] = createHtmlPlugin({
    minify: isBuild,
    inject: {
      // 向 EJS 模板注入数据
      data: {
        title: VITE_GLOB_APP_TITLE, // 注入页面标题
      },
      // 在生产环境下注入生成的应用配置脚本
      tags: isBuild
        ? [
            {
              tag: 'script',
              attrs: {
                src: getAppConfigSrc(), // 配置文件的 URL
              },
            },
          ]
        : [], // 开发环境下不注入脚本
    },
  });
  return htmlPlugin;
}
