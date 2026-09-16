<template>
  <div>
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>参数设置</span>
          <el-button type="primary" size="small" v-permission="'config:add'" @click="openDialog()">新增参数</el-button>
        </div>
      </template>
      <div class="toolbar">
        <el-input v-model="filters.key" placeholder="参数键" clearable style="width:160px" />
        <el-input v-model="filters.name" placeholder="名称" clearable style="width:160px" />
        <el-input v-model="filters.group" placeholder="分组" clearable style="width:140px" />
        <el-button type="primary" @click="search">查询</el-button>
        <el-button @click="resetFilters">重置</el-button>
      </div>
      <el-table :data="list" border stripe v-loading="loading">
        <el-table-column prop="key" label="参数键" width="180" />
        <el-table-column prop="name" label="名称" width="160" />
        <el-table-column prop="group" label="分组" width="120" />
        <el-table-column prop="value" label="参数值" show-overflow-tooltip />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" show-overflow-tooltip />
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button size="small" v-permission="'config:edit'" @click="openDialog(row)">编辑</el-button>
            <el-button size="small" type="danger" v-permission="'config:del'" @click="remove(row)">删除</el-button>
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
    </el-card>

    <el-dialog v-model="visible" :title="form.id ? '编辑参数' : '新增参数'" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="参数键">
          <el-input v-model="form.key" :disabled="!!form.id" placeholder="如 site.name" />
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="分组">
          <el-input v-model="form.group" placeholder="如 基础设置" />
        </el-form-item>
        <el-form-item label="参数值">
          <el-input v-model="form.value" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" @click="submit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listConfigs, createConfig, updateConfig, deleteConfig } from '@/api'

const list = ref([])
const total = ref(0)
const loading = ref(false)
const filters = reactive({ page: 1, pageSize: 10, key: '', name: '', group: '' })

const visible = ref(false)
const form = reactive({ id: null, key: '', name: '', group: '', value: '', sort: 0, status: 1, remark: '' })

async function load() {
  loading.value = true
  try {
    const res = await listConfigs({ ...filters })
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
  filters.key = ''
  filters.name = ''
  filters.group = ''
  filters.page = 1
  load()
}
function onPage(p) {
  filters.page = p
  load()
}

function openDialog(row) {
  if (row) {
    Object.assign(form, {
      id: row.id, key: row.key, name: row.name, group: row.group,
      value: row.value, sort: row.sort, status: row.status, remark: row.remark
    })
  } else {
    Object.assign(form, { id: null, key: '', name: '', group: '', value: '', sort: 0, status: 1, remark: '' })
  }
  visible.value = true
}

async function submit() {
  const payload = {
    key: form.key, name: form.name, group: form.group, value: form.value,
    sort: form.sort, status: form.status, remark: form.remark
  }
  if (form.id) await updateConfig(form.id, payload)
  else await createConfig(payload)
  ElMessage.success('保存成功')
  visible.value = false
  load()
}

async function remove(row) {
  await ElMessageBox.confirm(`确认删除参数「${row.name}」？`, '提示', { type: 'warning' })
  await deleteConfig(row.id)
  ElMessage.success('已删除')
  load()
}

onMounted(load)
</script>

<style scoped>
.card-header { display: flex; align-items: center; justify-content: space-between; }
.toolbar { display: flex; gap: 8px; margin-bottom: 12px; flex-wrap: wrap; }
.pager { margin-top: 12px; display: flex; justify-content: flex-end; }
</style>
