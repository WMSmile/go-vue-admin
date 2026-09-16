<template>
  <template v-for="item in items" :key="item.id">
    <el-sub-menu v-if="isCatalog(item)" :index="resolve(item)">
      <template #title>
        <el-icon v-if="item.icon"><component :is="item.icon" /></el-icon>
        <span>{{ item.title }}</span>
      </template>
      <sidebar-menu :items="item.children || []" :parent="resolve(item)" />
    </el-sub-menu>
    <el-menu-item v-else-if="item.type === 2" :index="resolve(item)">
      <el-icon v-if="item.icon"><component :is="item.icon" /></el-icon>
      <span>{{ item.title }}</span>
    </el-menu-item>
  </template>
</template>

<script setup>
const props = defineProps({
  items: { type: Array, default: () => [] },
  parent: { type: String, default: '' }
})

const isCatalog = (i) => i.type === 1 && i.children && i.children.length

function resolve(item) {
  if (!item.path) return props.parent
  if (item.path.startsWith('/')) return item.path
  return props.parent ? `${props.parent}/${item.path}` : item.path
}
</script>
