<template>
  <div>
    <div class="toolbar">
      <el-input v-model="filters.name" placeholder="任务名称" clearable style="width:160px" />
      <el-button type="primary" @click="search">查询</el-button>
      <el-button @click="resetFilters">重置</el-button>
      <el-button type="primary" @click="openDialog()">新增任务</el-button>
    </div>

    <el-table :data="list" border stripe v-loading="loading">
      <el-table-column prop="name" label="名称" min-width="140" />
      <el-table-column prop="type" label="类型" width="90">
        <template #default="{ row }">
          <el-tag :type="row.type === 1 ? 'primary' : 'success'">{{ row.type === 1 ? 'HTTP' : '函数' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="spec" label="CRON 表达式" width="150" />
      <el-table-column label="目标" min-width="220">
        <template #default="{ row }">
          <span v-if="row.type === 1">{{ row.method || 'GET' }} {{ row.url }}</span>
          <span v-else>{{ row.handler }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '运行中' : '已停用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="lastRunAt" label="上次执行" width="180">
        <template #default="{ row }">{{ fmt(row.lastRunAt) }}</template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" show-overflow-tooltip />
      <el-table-column label="操作" width="290" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="run(row)">执行</el-button>
          <el-button size="small" @click="toggle(row)">{{ row.status === 1 ? '停用' : '启用' }}</el-button>
          <el-button size="small" @click="openLogs(row)">日志</el-button>
          <el-button size="small" @click="openDialog(row)">编辑</el-button>
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

    <el-dialog v-model="visible" :title="form.id ? '编辑任务' : '新增任务'" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="类型">
          <el-radio-group v-model="form.type">
            <el-radio :value="1">HTTP 请求</el-radio>
            <el-radio :value="2">函数</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="CRON">
          <el-input v-model="form.spec" placeholder="如 0 0 * * * （每天0点）" />
        </el-form-item>
        <template v-if="form.type === 1">
          <el-form-item label="方法">
            <el-select v-model="form.method" style="width:100%">
              <el-option label="GET" value="GET" />
              <el-option label="POST" value="POST" />
              <el-option label="PUT" value="PUT" />
              <el-option label="DELETE" value="DELETE" />
            </el-select>
          </el-form-item>
          <el-form-item label="URL"><el-input v-model="form.url" placeholder="http(s)://..." /></el-form-item>
        </template>
        <template v-else>
          <el-form-item label="Handler">
            <el-select v-model="form.handler" filterable allow-create style="width:100%">
              <el-option label="heartbeat" value="heartbeat" />
              <el-option label="clear_operation_logs（清理30天前操作日志）" value="clear_operation_logs" />
            </el-select>
          </el-form-item>
        </template>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" type="textarea" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" @click="submit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="logsVisible" title="执行日志" width="640px">
      <el-table :data="logs" border stripe v-loading="logsLoading">
        <el-table-column prop="status" label="结果" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">{{ row.status === 1 ? '成功' : '失败' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="durationMs" label="耗时(ms)" width="100" />
        <el-table-column prop="createdAt" label="时间" width="180">
          <template #default="{ row }">{{ fmt(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="输出">
          <template #default="{ row }"><pre class="log-out">{{ row.output }}</pre></template>
        </el-table-column>
      </el-table>
      <div class="pager">
        <el-pagination
          background
          layout="prev, pager, next"
          :total="logsTotal"
          :page-size="logsFilters.pageSize"
          :current-page="logsFilters.page"
          @current-change="onLogsPage"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listTasks, createTask, updateTask, deleteTask, toggleTask, runTask, listTaskLogs
} from '@/api'

const list = ref([])
const total = ref(0)
const loading = ref(false)
const filters = reactive({ page: 1, pageSize: 10, name: '' })

const visible = ref(false)
const form = reactive({ id: null, name: '', type: 1, spec: '', method: 'GET', url: '', handler: '', status: 1, remark: '' })

const logsVisible = ref(false)
const logs = ref([])
const logsTotal = ref(0)
const logsLoading = ref(false)
const logsFilters = reactive({ page: 1, pageSize: 10, taskId: null })

function fmt(t) {
  if (!t) return '-'
  const d = new Date(t)
  const p = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
}

async function load() {
  loading.value = true
  try {
    const res = await listTasks({ ...filters })
    list.value = res.list || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

function search() { filters.page = 1; load() }
function resetFilters() { filters.name = ''; filters.page = 1; load() }
function onPage(p) { filters.page = p; load() }

function resetForm() {
  Object.assign(form, { id: null, name: '', type: 1, spec: '', method: 'GET', url: '', handler: '', status: 1, remark: '' })
}
function openDialog(row) {
  if (row) Object.assign(form, { ...row })
  else resetForm()
  visible.value = true
}

async function submit() {
  const payload = {
    name: form.name, type: form.type, spec: form.spec,
    method: form.method, url: form.url, handler: form.handler,
    status: form.status, remark: form.remark
  }
  if (form.id) await updateTask(form.id, payload)
  else await createTask(payload)
  ElMessage.success('保存成功')
  visible.value = false
  load()
}

async function toggle(row) {
  await toggleTask(row.id)
  ElMessage.success('已切换')
  load()
}

async function run(row) {
  await runTask(row.id)
  ElMessage.success('已触发执行')
  load()
}

async function remove(row) {
  await ElMessageBox.confirm('确认删除该定时任务？', '提示', { type: 'warning' })
  await deleteTask(row.id)
  ElMessage.success('已删除')
  load()
}

async function openLogs(row) {
  logsFilters.taskId = row.id
  logsFilters.page = 1
  logsVisible.value = true
  loadLogs()
}
async function loadLogs() {
  logsLoading.value = true
  try {
    const res = await listTaskLogs({ ...logsFilters })
    logs.value = res.list || []
    logsTotal.value = res.total || 0
  } finally {
    logsLoading.value = false
  }
}
function onLogsPage(p) { logsFilters.page = p; loadLogs() }

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; gap: 8px; margin-bottom: 12px; flex-wrap: wrap; }
.pager { margin-top: 12px; display: flex; justify-content: flex-end; }
.log-out { margin: 0; white-space: pre-wrap; word-break: break-all; font-size: 12px; max-height: 120px; overflow: auto; }
</style>
