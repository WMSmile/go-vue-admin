<template>
  <div class="settings">
    <el-card>
      <template #header>外观设置</template>
      <el-form label-width="120px">
        <el-form-item label="主题">
          <el-radio-group :model-value="theme.mode" @change="theme.setMode">
            <el-radio label="light">浅色</el-radio>
            <el-radio label="dark">深色</el-radio>
            <el-radio label="system">跟随系统</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="说明">
          <span class="tip">深色 / 浅色实时切换并自动保存；"跟随系统" 根据操作系统的外观设置自动调整。</span>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { useThemeStore } from '@/store/theme'
const theme = useThemeStore()
</script>

<style scoped>
.settings { max-width: 720px; }
.tip { color: #999; font-size: 13px; }
</style>
