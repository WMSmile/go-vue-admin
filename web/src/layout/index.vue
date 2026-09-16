<template>
  <el-container class="app-wrapper">
    <el-aside :width="collapse ? '64px' : '220px'" class="sidebar" :class="{ collapsed: collapse }">
      <div class="logo">go-vue-admin</div>
      <Sidebar :collapse="collapse" />
    </el-aside>
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-button text class="collapse-btn" :title="collapse ? '展开菜单' : '收起菜单'" @click="toggleCollapse">
            <el-icon><Fold v-if="!collapse" /><Expand v-else /></el-icon>
          </el-button>
          <span class="title">管理后台</span>
        </div>
        <el-dropdown @command="onCommand">
          <span class="user-info">
            <el-avatar :size="32" class="avatar">{{ initial }}</el-avatar>
            <span class="username">{{ userStore.userInfo?.username || '用户' }}</span>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="settings">设置</el-dropdown-item>
              <el-dropdown-item command="profile">个人中心</el-dropdown-item>
              <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
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
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import Sidebar from './components/Sidebar.vue'

const router = useRouter()
const userStore = useUserStore()

const collapse = ref(false)
function toggleCollapse() {
  collapse.value = !collapse.value
}

const initial = (userStore.userInfo?.username || '?').charAt(0).toUpperCase()

function onCommand(cmd) {
  if (cmd === 'settings') {
    router.push('/settings')
  } else if (cmd === 'profile') {
    router.push('/profile')
  } else if (cmd === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.app-wrapper { height: 100vh; }
.header { display: flex; align-items: center; justify-content: space-between; background: #fff; border-bottom: 1px solid #eee; }
.header-left { display: flex; align-items: center; gap: 8px; }
.collapse-btn { font-size: 18px; }
.user-info { display: flex; align-items: center; gap: 8px; cursor: pointer; outline: none; }
.user-info .avatar { background: #409eff; }
.user-info .username { font-size: 14px; }
</style>

<style>
.sidebar { background: var(--el-bg-color, #fff); color: var(--el-text-color-primary, #303133); overflow-y: auto; transition: width .2s; scrollbar-width: none; -ms-overflow-style: none; }
.sidebar::-webkit-scrollbar { width: 0; height: 0; display: none; }
.logo { height: 56px; line-height: 56px; text-align: center; font-weight: 600; color: var(--el-text-color-primary, #303133); white-space: nowrap; }
.sidebar.collapsed .logo { display: none; }
html.dark .sidebar { background: #001529; color: #fff; }
html.dark .logo { color: #fff; }
html.dark .sidebar .el-menu {
  --el-menu-bg-color: #001529;
  --el-menu-text-color: #fff;
  --el-menu-active-color: #409eff;
  --el-menu-hover-bg-color: #000c17;
  background-color: var(--el-menu-bg-color);
  color: var(--el-menu-text-color);
}
html.dark .sidebar .el-menu .el-menu-item.is-active { color: var(--el-menu-active-color); }
</style>
