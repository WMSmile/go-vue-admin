<template>
  <div>
    <el-row :gutter="16">
      <el-col :span="8" v-for="c in cards" :key="c.label">
        <el-card shadow="hover">
          <div class="num">{{ stats[c.key] ?? 0 }}</div>
          <div class="label">{{ c.label }}</div>
        </el-card>
      </el-col>
    </el-row>
    <el-card class="mt">
      <h3>欢迎使用 go-vue-admin</h3>
      <p>基于 Gin + GORM + Vue3 + Casbin RBAC + JWT 的权限管理脚手架。</p>
    </el-card>
  </div>
</template>

<script setup>
import { onMounted, reactive } from 'vue'
import { getDashboard } from '@/api'

const stats = reactive({ users: 0, roles: 0, menus: 0 })
const cards = [
  { key: 'users', label: '用户数' },
  { key: 'roles', label: '角色数' },
  { key: 'menus', label: '菜单数' }
]

onMounted(async () => {
  try {
    const data = await getDashboard()
    Object.assign(stats, data)
  } catch (e) {}
})
</script>

<style scoped>
.num { font-size: 32px; font-weight: 600; color: #409eff; }
.label { color: #888; margin-top: 6px; }
.mt { margin-top: 16px; }
</style>
