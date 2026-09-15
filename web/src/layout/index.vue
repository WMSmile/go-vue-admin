<template>
  <el-container class="app-wrapper">
    <el-aside width="220px" class="sidebar">
      <div class="logo">go-vue-admin</div>
      <Sidebar />
    </el-aside>
    <el-container>
      <el-header class="header">
        <span class="title">管理后台</span>
        <el-dropdown @command="onCommand">
          <span class="user">{{ userStore.userInfo?.username || '用户' }}</span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import Sidebar from './components/Sidebar.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

function onCommand(cmd) {
  if (cmd === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.app-wrapper { height: 100vh; }
.sidebar { background: #001529; color: #fff; overflow-y: auto; }
.logo { height: 56px; line-height: 56px; text-align: center; font-weight: 600; color: #fff; }
.header { display: flex; align-items: center; justify-content: space-between; background: #fff; border-bottom: 1px solid #eee; }
.user { cursor: pointer; }
</style>
