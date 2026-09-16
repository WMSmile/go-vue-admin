<template>
  <el-config-provider :locale="zhCn">
    <router-view />
  </el-config-provider>
</template>

<script setup>
import { ElConfigProvider } from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import { useThemeStore } from '@/store/theme'

const theme = useThemeStore()
theme.init()
</script>

<style>
html, body, #app { height: 100%; margin: 0; }

/* dark-mode overrides for the custom layout surfaces (Element Plus handles its own) */
html.dark .header { background: #1f1f1f; border-bottom-color: #333; }
html.dark .header .title,
html.dark .header .username { color: #e5e5e5; }
html.dark body { background: #141414; }
</style>
