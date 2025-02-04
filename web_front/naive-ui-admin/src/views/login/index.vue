<template>
  <!-- 登录页面主容器 -->
  <div class="view-account">
    <div class="view-account-header"></div>
    <div class="view-account-container">
      <!-- 顶部 Logo 和描述区域 -->
      <div class="view-account-top">
        <div class="view-account-top-logo">
          <img :src="websiteConfig.loginImage" alt="" />
        </div>
        <div class="view-account-top-desc">{{ websiteConfig.loginDesc }}</div>
      </div>
      
      <!-- 登录表单区域 -->
      <div class="view-account-form">
        <n-form
          ref="formRef"
          label-placement="left"
          size="large"
          :model="formInline"
          :rules="rules"
        >
          <!-- 用户名输入框 -->
          <n-form-item path="username">
            <n-input v-model:value="formInline.username" placeholder="请输入用户名">
              <template #prefix>
                <n-icon size="18" color="#808695">
                  <PersonOutline />
                </n-icon>
              </template>
            </n-input>
          </n-form-item>

          <!-- 密码输入框 -->
          <n-form-item path="password">
            <n-input
              v-model:value="formInline.password"
              type="password"
              showPasswordOn="mousedown"
              placeholder="请输入密码"
            >
              <template #prefix>
                <n-icon size="18" color="#808695">
                  <LockClosedOutline />
                </n-icon>
              </template>
            </n-input>
          </n-form-item>

          <!-- 自动登录和忘记密码选项 -->
          <n-form-item class="default-color">
            <div class="flex justify-between">
              <div class="flex-initial">
                <n-checkbox v-model:checked="autoLogin">自动登录</n-checkbox>
              </div>
              <div class="flex-initial order-last">
                <a href="javascript:">忘记密码</a>
              </div>
            </div>
          </n-form-item>

          <!-- 登录按钮 -->
          <n-form-item>
            <n-button type="primary" @click="handleSubmit" size="large" :loading="loading" block>
              登录
            </n-button>
          </n-form-item>

          <!-- 其他登录方式 -->
          <n-form-item class="default-color">
            <div class="flex view-  -other">
              <div class="flex-initial">
                <span>其它登录方式</span>
              </div>
              <div class="flex-initial mx-2">
                <a href="javascript:">
                  <n-icon size="24" color="#2d8cf0">
                    <LogoGithub />
                  </n-icon>
                </a>
              </div>
              <div class="flex-initial mx-2">
                <a href="javascript:">
                  <n-icon size="24" color="#2d8cf0">
                    <LogoFacebook />
                  </n-icon>
                </a>
              </div>
              <div class="flex-initial" style="margin-left: auto">
                <a href="javascript:">注册账号</a>
              </div>
            </div>
          </n-form-item>
        </n-form>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  // 导入所需的组件和工具
  import { reactive, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import { useUserStore } from '@/store/modules/user';
  import { useMessage } from 'naive-ui';
  import { ResultEnum } from '@/enums/httpEnum';
  import { PersonOutline, LockClosedOutline, LogoGithub, LogoFacebook } from '@vicons/ionicons5';
  import { PageEnum } from '@/enums/pageEnum';
  import { websiteConfig } from '@/config/website.config';

  // 定义表单数据接口
  interface FormState {
    username: string;
    password: string;
  }

  // 初始化表单引用和其他响应式变量
  const formRef = ref(); // 表单引用
  const message = useMessage(); // 消息提示
  const loading = ref(false); // 加载状态
  const autoLogin = ref(true); // 自动登录选项
  const LOGIN_NAME = PageEnum.BASE_LOGIN_NAME;

  // 表单数据对象
  const formInline = reactive({
    username: 'admin', // 默认用户名
    password: '123456', // 默认密码
    isCaptcha: true, // 是否需要验证码
  });

  // 表单验证规则
  const rules = {
    username: { required: true, message: '请输入用户名', trigger: 'blur' },
    password: { required: true, message: '请输入密码', trigger: 'blur' },
  };

  // 获取用户存储实例和路由实例
  const userStore = useUserStore();
  const router = useRouter();
  const route = useRoute();

  // 表单提交处理函数
  const handleSubmit = (e) => {
    e.preventDefault();
    // 表单验证
    formRef.value.validate(async (errors) => {
      if (!errors) {
        // 获取表单数据
        const { username, password } = formInline;
        message.loading('登录中...');
        loading.value = true;

        const params: FormState = {
          username,
          password,
        };

        try {
          // 调用登录接口
          // 使用解构赋值，并重命名 message 为 msg！！
          const { code, message: msg } = await userStore.login(params);
          message.destroyAll();
          if (code == ResultEnum.SUCCESS) {
            // 登录成功处理
            const toPath = decodeURIComponent((route.query?.redirect || '/') as string);
            message.success('登录成功，即将进入系统');
            // 路由跳转
            if (route.name === LOGIN_NAME) {
              router.replace('/');
            } else router.replace(toPath);
          } else {
            // 登录失败处理
            message.info(msg || '登录失败');
          }
        } finally {
          loading.value = false;
        }
      } else {
        // 表单验证失败提示
        message.error('请填写完整信息，并且进行验证码校验');
      }
    });
  };
</script>

<style lang="less" scoped>
  // 登录页面样式
  .view-account {
    display: flex;
    flex-direction: column;
    height: 100vh;
    overflow: auto;

    &-container {
      flex: 1;
      padding: 32px 12px;
      max-width: 384px;
      min-width: 320px;
      margin: 0 auto;
    }

    &-top {
      padding: 32px 0;
      text-align: center;

      &-desc {
        font-size: 14px;
        color: #808695;
      }
    }

    &-other {
      width: 100%;
    }

    .default-color {
      color: #515a6e;

      .ant-checkbox-wrapper {
        color: #515a6e;
      }
    }
  }

  @media (min-width: 768px) {
    .view-account {
      background-image: url('../../assets/images/login.svg');
      background-repeat: no-repeat;
      background-position: 50%;
      background-size: 100%;
    }

    .page-account-container {
      padding: 32px 0 24px 0;
    }
  }
</style>
