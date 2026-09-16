<template>
  <div class="crontab">
    <el-tabs type="border-card">
      <el-tab-pane label="分钟">
        <CrontabMin @update="update" :check="checkNumber" :cron="obj" />
      </el-tab-pane>
      <el-tab-pane label="小时">
        <CrontabHour @update="update" :check="checkNumber" :cron="obj" />
      </el-tab-pane>
      <el-tab-pane label="日">
        <CrontabDay @update="update" :check="checkNumber" :cron="obj" />
      </el-tab-pane>
      <el-tab-pane label="月">
        <CrontabMonth @update="update" :check="checkNumber" :cron="obj" />
      </el-tab-pane>
      <el-tab-pane label="周">
        <CrontabWeek @update="update" :check="checkNumber" :cron="obj" />
      </el-tab-pane>
    </el-tabs>

    <div class="result-box">
      <p class="title">时间表达式</p>
      <table>
        <thead>
          <th v-for="t in titles" :key="t">{{ t }}</th>
          <th>Cron 表达式</th>
        </thead>
        <tbody>
          <td>{{ obj.min }}</td>
          <td>{{ obj.hour }}</td>
          <td>{{ obj.day }}</td>
          <td>{{ obj.month }}</td>
          <td>{{ obj.week }}</td>
          <td class="expr">{{ crontabValueString }}</td>
        </tbody>
      </table>
    </div>

    <CrontabResult :ex="crontabValueString" />

    <div class="pop_btn">
      <el-button type="primary" @click="submit">确定</el-button>
      <el-button type="warning" @click="clearCron">重置</el-button>
      <el-button @click="$emit('cancel')">取消</el-button>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import CrontabMin from './min.vue'
import CrontabHour from './hour.vue'
import CrontabDay from './day.vue'
import CrontabMonth from './month.vue'
import CrontabWeek from './week.vue'
import CrontabResult from './result.vue'

const emit = defineEmits(['fill', 'cancel'])
const props = defineProps({
  expression: { type: String, default: '' }
})

const titles = ['分钟', '小时', '日', '月', '周']
const expression = ref('')
const obj = ref({ min: '*', hour: '*', day: '*', month: '*', week: '?' })

const crontabValueString = computed(
  () => `${obj.value.min} ${obj.value.hour} ${obj.value.day} ${obj.value.month} ${obj.value.week}`
)

function update(name, value) {
  obj.value[name] = value
}

function checkNumber(value, minLimit, maxLimit) {
  value = Math.floor(value)
  if (value < minLimit) value = minLimit
  else if (value > maxLimit) value = maxLimit
  return value
}

function clearCron() {
  obj.value = { min: '*', hour: '*', day: '*', month: '*', week: '?' }
}

function resolveExp() {
  if (expression.value) {
    const arr = expression.value.trim().split(/\s+/)
    if (arr.length >= 5) {
      obj.value = {
        min: arr[0],
        hour: arr[1],
        day: arr[2],
        month: arr[3],
        week: arr[4]
      }
      return
    }
  }
  clearCron()
}

function submit() {
  emit('fill', crontabValueString.value)
}

onMounted(() => {
  expression.value = props.expression
  resolveExp()
})
watch(
  () => props.expression,
  () => {
    expression.value = props.expression
    resolveExp()
  }
)
</script>

<style scoped>
.pop_btn { text-align: center; margin-top: 16px; }
.result-box {
  box-sizing: border-box;
  margin: 18px auto 8px;
  padding: 14px 10px 10px;
  border: 1px solid #dcdfe6;
  position: relative;
}
.result-box .title {
  position: absolute;
  top: -13px;
  left: 50%;
  width: 140px;
  font-size: 13px;
  margin-left: -70px;
  text-align: center;
  line-height: 24px;
  background: #fff;
}
.result-box table { text-align: center; width: 100%; }
.result-box table th { font-weight: 600; color: #606266; }
.result-box table td {
  width: 3.2rem;
  min-width: 3.2rem;
  max-width: 3.2rem;
  font-family: arial;
  line-height: 28px;
  height: 28px;
  white-space: nowrap;
  overflow: hidden;
  border: 1px solid #ebeef5;
}
.result-box table td.expr { max-width: 100%; width: auto; }
.el-input-number--small,
.el-select { margin: 0 0.2rem; }
</style>
