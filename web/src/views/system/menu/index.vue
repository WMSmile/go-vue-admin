<template>
  <div>
    <el-table :data="list" border stripe row-key="id" :tree-props="{ children: 'children' }" default-expand-all>
      <el-table-column prop="title" label="名称" />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.type === 1">目录</el-tag>
          <el-tag v-else-if="row.type === 2" type="success">菜单</el-tag>
          <el-tag v-else type="warning">按钮</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="path" label="路径" />
      <el-table-column prop="component" label="组件" />
      <el-table-column prop="permission" label="权限标识" />
      <el-table-column prop="api" label="接口" />
      <el-table-column prop="method" label="方法" width="90" />
      <el-table-column prop="sort" label="排序" width="80" />
    </el-table>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { listMenus } from '@/api'

const list = ref([])

function toTree(flat) {
  const map = {}
  flat.forEach((m) => (map[m.id] = { ...m, children: [] }))
  const tree = []
  flat.forEach((m) => {
    if (m.parentId && map[m.parentId]) map[m.parentId].children.push(map[m.id])
    else tree.push(map[m.id])
  })
  return tree
}

onMounted(async () => {
  const flat = await listMenus()
  list.value = toTree(flat)
})
</script>
