<template>
  <div>
    <div class="toolbar">
      <el-select v-model="filters.method" clearable placeholder="请求方法" style="width: 120px">
        <el-option label="GET" value="GET" />
        <el-option label="POST" value="POST" />
        <el-option label="PUT" value="PUT" />
        <el-option label="DELETE" value="DELETE" />
      </el-select>
      <el-input v-model="filters.username" placeholder="操作人" clearable style="width: 140px" />
      <el-input v-model="filters.path" placeholder="接口路径" clearable style="width: 220px" />
      <el-button type="primary" @click="search">查询</el-button>
      <el-button @click="resetFilters">重置</el-button>
      <el-button type="danger" plain @click="clearAll">清空日志</el-button>
    </div>

    <el-table :data="list" border stripe v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="username" label="操作人" width="120">
        <template #default="{ row }">{{ row.username || '-' }}</template>
      </el-table-column>
      <el-table-column prop="method" label="方法" width="100">
        <template #default="{ row }">
          <el-tag :type="methodType(row.method)" disable-transitions>{{ row.method }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="path" label="接口路径" min-width="220" show-overflow-tooltip />
      <el-table-column prop="ip" label="IP" width="130" />
      <el-table-column prop="status" label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)" disable-transitions>{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="latency" label="耗时(ms)" width="100" />
      <el-table-column prop="createdAt" label="操作时间" width="180">
        <template #default="{ row }">{{ fmt(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="90" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pager">
      <el-pagination
        background
        layout="total, prev, pager, next"
        :total="total"
        :page-size="filters.pageSize"
        :current-page="filters.page"
        @current-change="onPage"
      />
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listOperationLogs, deleteOperationLog, clearOperationLogs } from '@/api'

const list = ref([])
const total = ref(0)
const loading = ref(false)
const filters = reactive({ page: 1, pageSize: 10, method: '', username: '', path: '' })

const methodType = (m) => ({ GET: 'success', POST: 'primary', PUT: 'warning', DELETE: 'danger' }[m] || 'info')
const statusType = (s) => (s >= 200 && s < 400 ? 'success' : 'danger')

function fmt(t) {
  if (!t) return '-'
  const d = new Date(t)
  const p = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
}

async function load() {
  loading.value = true
  try {
    const res = await listOperationLogs({ ...filters })
    list.value = res.list || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

function search() {
  filters.page = 1
  load()
}

function resetFilters() {
  filters.method = ''
  filters.username = ''
  filters.path = ''
  filters.page = 1
  load()
}

function onPage(p) {
  filters.page = p
  load()
}

async function remove(row) {
  await ElMessageBox.confirm('确认删除该日志记录？', '提示', { type: 'warning' })
  await deleteOperationLog(row.id)
  ElMessage.success('已删除')
  load()
}

async function clearAll() {
  await ElMessageBox.confirm('确认清空所有操作日志？此操作不可恢复。', '提示', { type: 'warning' })
  await clearOperationLogs()
  ElMessage.success('已清空')
  load()
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; gap: 8px; margin-bottom: 12px; flex-wrap: wrap; }
.pager { margin-top: 12px; display: flex; justify-content: flex-end; }
</style>
