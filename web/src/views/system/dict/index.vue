<template>
  <div>
    <!-- 字典类型 -->
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>字典类型</span>
          <el-button type="primary" size="small" @click="openTypeDialog()">新增类型</el-button>
        </div>
      </template>
      <div class="toolbar">
        <el-input v-model="typeFilters.code" placeholder="编码" clearable style="width:140px" />
        <el-input v-model="typeFilters.name" placeholder="名称" clearable style="width:140px" />
        <el-button type="primary" @click="searchTypes">查询</el-button>
        <el-button @click="resetTypeFilters">重置</el-button>
      </div>
      <el-table
        :data="types"
        border
        stripe
        highlight-current-row
        :current-change="onTypeChange"
        v-loading="typeLoading"
      >
        <el-table-column prop="code" label="编码" width="160" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" show-overflow-tooltip />
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button size="small" @click="openTypeDialog(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="removeType(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pager">
        <el-pagination
          background
          layout="total, prev, pager, next"
          :total="typeTotal"
          :page-size="typeFilters.pageSize"
          :current-page="typeFilters.page"
          @current-change="onTypePage"
        />
      </div>
    </el-card>

    <!-- 字典数据 -->
    <el-card class="box-card" style="margin-top:12px">
      <template #header>
        <div class="card-header">
          <span>字典数据{{ selectedType ? '（' + selectedType.name + '）' : '' }}</span>
          <el-button type="primary" size="small" :disabled="!selectedType" @click="openDataDialog()">新增数据</el-button>
        </div>
      </template>
      <el-table :data="dataList" border stripe v-loading="dataLoading">
        <el-table-column prop="label" label="标签" width="200" />
        <el-table-column prop="value" label="值" width="200" />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" show-overflow-tooltip />
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button size="small" @click="openDataDialog(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="removeData(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 类型弹窗 -->
    <el-dialog v-model="typeVisible" :title="typeForm.id ? '编辑字典类型' : '新增字典类型'" width="480px">
      <el-form :model="typeForm" label-width="80px">
        <el-form-item label="编码">
          <el-input v-model="typeForm.code" :disabled="!!typeForm.id" />
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="typeForm.name" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="typeForm.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="typeForm.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="typeVisible = false">取消</el-button>
        <el-button type="primary" @click="submitType">确定</el-button>
      </template>
    </el-dialog>

    <!-- 数据弹窗 -->
    <el-dialog v-model="dataVisible" :title="dataForm.id ? '编辑字典数据' : '新增字典数据'" width="480px">
      <el-form :model="dataForm" label-width="80px">
        <el-form-item label="标签">
          <el-input v-model="dataForm.label" />
        </el-form-item>
        <el-form-item label="值">
          <el-input v-model="dataForm.value" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="dataForm.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="dataForm.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="dataForm.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dataVisible = false">取消</el-button>
        <el-button type="primary" @click="submitData">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listDictTypes, createDictType, updateDictType, deleteDictType,
  listDictData, createDictData, updateDictData, deleteDictData
} from '@/api'

const types = ref([])
const typeTotal = ref(0)
const typeLoading = ref(false)
const typeFilters = reactive({ page: 1, pageSize: 10, code: '', name: '' })

const selectedType = ref(null)
const dataList = ref([])
const dataLoading = ref(false)

const typeVisible = ref(false)
const typeForm = reactive({ id: null, code: '', name: '', status: 1, remark: '' })

const dataVisible = ref(false)
const dataForm = reactive({ id: null, typeId: null, label: '', value: '', sort: 0, status: 1, remark: '' })

async function loadTypes() {
  typeLoading.value = true
  try {
    const res = await listDictTypes({ ...typeFilters })
    types.value = res.list || []
    typeTotal.value = res.total || 0
  } finally {
    typeLoading.value = false
  }
}

function searchTypes() {
  typeFilters.page = 1
  loadTypes()
}
function resetTypeFilters() {
  typeFilters.code = ''
  typeFilters.name = ''
  typeFilters.page = 1
  loadTypes()
}
function onTypePage(p) {
  typeFilters.page = p
  loadTypes()
}

function onTypeChange(row) {
  selectedType.value = row || null
  dataForm.typeId = row ? row.id : null
  loadData()
}

async function loadData() {
  if (!selectedType.value) {
    dataList.value = []
    return
  }
  dataLoading.value = true
  try {
    const res = await listDictData({ typeId: selectedType.value.id })
    dataList.value = res.list || []
  } finally {
    dataLoading.value = false
  }
}

function openTypeDialog(row) {
  if (row) {
    Object.assign(typeForm, { id: row.id, code: row.code, name: row.name, status: row.status, remark: row.remark })
  } else {
    Object.assign(typeForm, { id: null, code: '', name: '', status: 1, remark: '' })
  }
  typeVisible.value = true
}

async function submitType() {
  const payload = { code: typeForm.code, name: typeForm.name, status: typeForm.status, remark: typeForm.remark }
  if (typeForm.id) await updateDictType(typeForm.id, payload)
  else await createDictType(payload)
  ElMessage.success('保存成功')
  typeVisible.value = false
  loadTypes()
  if (selectedType.value && selectedType.value.id === typeForm.id) selectedType.value.name = typeForm.name
}

async function removeType(row) {
  await ElMessageBox.confirm(`确认删除字典类型「${row.name}」及其下所有数据？`, '提示', { type: 'warning' })
  await deleteDictType(row.id)
  ElMessage.success('已删除')
  if (selectedType.value && selectedType.value.id === row.id) selectedType.value = null
  loadTypes()
}

function openDataDialog(row) {
  if (!selectedType.value) return
  if (row) {
    Object.assign(dataForm, {
      id: row.id, typeId: selectedType.value.id, label: row.label,
      value: row.value, sort: row.sort, status: row.status, remark: row.remark
    })
  } else {
    Object.assign(dataForm, { id: null, typeId: selectedType.value.id, label: '', value: '', sort: 0, status: 1, remark: '' })
  }
  dataVisible.value = true
}

async function submitData() {
  const payload = {
    typeId: dataForm.typeId, label: dataForm.label, value: dataForm.value,
    sort: dataForm.sort, status: dataForm.status, remark: dataForm.remark
  }
  if (dataForm.id) await updateDictData(dataForm.id, payload)
  else await createDictData(payload)
  ElMessage.success('保存成功')
  dataVisible.value = false
  loadData()
}

async function removeData(row) {
  await ElMessageBox.confirm('确认删除该字典数据？', '提示', { type: 'warning' })
  await deleteDictData(row.id)
  ElMessage.success('已删除')
  loadData()
}

onMounted(loadTypes)
</script>

<style scoped>
.card-header { display: flex; align-items: center; justify-content: space-between; }
.toolbar { display: flex; gap: 8px; margin-bottom: 12px; flex-wrap: wrap; }
.pager { margin-top: 12px; display: flex; justify-content: flex-end; }
</style>
