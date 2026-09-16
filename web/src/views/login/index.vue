<template>
  <div class="login-container">
    <el-card class="login-card">
      <h2 class="title">{{ sysConfig.get('site.loginTitle', 'go-vue-admin') }}</h2>
      <el-form :model="form" label-width="0" @submit.prevent="onLogin">
        <el-form-item>
          <el-input v-model="form.username" placeholder="用户名" size="large" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.password" type="password" placeholder="密码" size="large" show-password />
        </el-form-item>
        <el-button type="primary" size="large" style="width: 100%" :loading="loading" @click="onLogin">
          登录
        </el-button>
      </el-form>
      <p class="tip">默认账号：admin / admin123</p>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { useSysConfigStore } from '@/store/sysConfig'

const router = useRouter()
const userStore = useUserStore()
const sysConfig = useSysConfigStore()
const form = reactive({ username: '', password: '' })
const loading = ref(false)

async function onLogin() {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    await userStore.login({ ...form })
    ElMessage.success('登录成功')
  } catch (e) {
    // error already shown by interceptor
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1e3c72, #2a5298);
}
.login-card { width: 360px; padding: 10px 20px; }
.title { text-align: center; margin-bottom: 20px; }
.tip { text-align: center; color: #999; font-size: 12px; margin-top: 10px; }
</style>
