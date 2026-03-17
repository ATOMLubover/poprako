<template>
  <div class="login-page">
    <div class="glow glow-left"></div>
    <div class="glow glow-right"></div>
    <a-card class="login-card" :bordered="false">
      <template #title>
        <div class="title-wrap">
          <span class="title-main">Poprako 指挥台</span>
        </div>
      </template>
      <a-form layout="vertical" :model="loginForm" @finish="handleSubmit">
        <a-form-item
          label="QQ"
          name="qq"
          :rules="[{ required: true, message: '请输入 QQ 号' }]"
        >
          <a-input
            v-model:value="loginForm.qq"
            placeholder="请输入 QQ"
            size="large"
          />
        </a-form-item>
        <a-form-item
          label="密码"
          name="password"
          :rules="[{ required: true, message: '请输入密码' }]"
        >
          <a-input-password
            v-model:value="loginForm.password"
            placeholder="请输入密码"
            size="large"
          />
        </a-form-item>
        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            size="large"
            block
            :loading="submitting"
          >
            登录并进入工作台
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { message } from "ant-design-vue";
import { loginUser, type LoginUserArgs } from "../api/modules";
import { useAuthStore } from "../stores/auth";

const router = useRouter();
const authStore = useAuthStore();
const submitting = ref(false);
const loginForm = reactive<LoginUserArgs>({
  qq: "",
  password: "",
});

/**
 * 处理登录提交。
 * 登录成功后会缓存令牌并跳转到仪表盘页面。
 */
async function handleSubmit(): Promise<void> {
  submitting.value = true;
  try {
    const loginUserResult = await loginUser(loginForm);
    authStore.setAccessToken(loginUserResult.access_token);
    message.success("登录成功，正在进入仪表盘");
    await router.push("/dashboard");
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "登录失败";
    message.error(errorMessage);
  } finally {
    submitting.value = false;
  }
}
</script>

<style scoped lang="scss">
/* 登录页容器与背景层。 */
.login-page {
  position: relative;
  min-height: 100vh;
  display: grid;
  place-items: center;
  overflow: hidden;
  background: linear-gradient(
    120deg,
    var(--login-grad-primary),
    var(--login-grad-accent)
  );
}

/* 登录卡片主体。 */
.login-card {
  width: min(460px, calc(100vw - 32px));
  border-radius: 20px;
  box-shadow: var(--shadow-heavy);
  backdrop-filter: blur(8px);
}

.title-wrap {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.title-main {
  font-size: 24px;
  font-weight: 700;
  letter-spacing: 0.02em;
}

.title-sub {
  color: var(--text-muted);
  font-size: 13px;
}

.glow {
  position: absolute;
  width: 380px;
  height: 380px;
  border-radius: 999px;
  filter: blur(44px);
}

.glow-left {
  left: -120px;
  top: -100px;
  background: var(--login-glow-primary);
}

.glow-right {
  right: -140px;
  bottom: -120px;
  background: var(--login-glow-accent);
}

:global(html[data-theme="dark"]) {
  /* 暗黑模式下登录卡片使用全局暗色 token。 */
  .login-card {
    background: var(--login-card-bg, var(--surface));
    color: var(--text-primary);
  }

  .title-sub {
    color: var(--text-muted);
  }
}
</style>
