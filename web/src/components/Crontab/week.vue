<template>
  <el-form>
    <el-form-item>
      <el-radio v-model="radioValue" :value="1">周，允许的通配符[, - * /]</el-radio>
    </el-form-item>
    <el-form-item>
      <el-radio v-model="radioValue" :value="2">不指定</el-radio>
    </el-form-item>
    <el-form-item>
      <el-radio v-model="radioValue" :value="3">
        周期从
        <el-select clearable v-model="cycle01" style="width:8rem">
          <el-option v-for="item in weekList" :key="item.key" :label="item.value" :value="item.key" :disabled="item.key === 6" />
        </el-select>
        -
        <el-select clearable v-model="cycle02" style="width:8rem">
          <el-option v-for="item in weekList" :key="item.key" :label="item.value" :value="item.key" :disabled="item.key <= cycle01" />
        </el-select>
      </el-radio>
    </el-form-item>
    <el-form-item>
      <el-radio v-model="radioValue" :value="4">
        指定
        <el-select clearable v-model="checkboxList" placeholder="可多选" multiple :multiple-limit="6" style="width:17.8rem">
          <el-option v-for="item in weekList" :key="item.key" :label="item.value" :value="item.key" />
        </el-select>
      </el-radio>
    </el-form-item>
  </el-form>
</template>

<script setup>
import { computed, ref, watch } from 'vue'

const emit = defineEmits(['update'])
const props = defineProps({
  cron: { type: Object, default: () => ({ min: '*', hour: '*', day: '*', month: '*', week: '?' }) },
  check: { type: Function, default: () => {} }
})

// 0=星期日 ... 6=星期六（与后端 time.Weekday() 一致）
const radioValue = ref(2)
const cycle01 = ref(1)
const cycle02 = ref(2)
const checkboxList = ref([])
const checkCopy = ref([1])
const weekList = ref([
  { key: 0, value: '星期日' }, { key: 1, value: '星期一' }, { key: 2, value: '星期二' },
  { key: 3, value: '星期三' }, { key: 4, value: '星期四' }, { key: 5, value: '星期五' },
  { key: 6, value: '星期六' }
])

const cycleTotal = computed(() => {
  cycle01.value = props.check(cycle01.value, 0, 5)
  cycle02.value = props.check(cycle02.value, cycle01.value + 1, 6)
  return cycle01.value + '-' + cycle02.value
})
const checkboxString = computed(() => checkboxList.value.join(','))

watch(() => props.cron.week, (value) => changeRadioValue(value))
watch([radioValue, cycleTotal, checkboxString], () => onRadioChange())

function changeRadioValue(value) {
  if (value === '*') radioValue.value = 1
  else if (value === '?') radioValue.value = 2
  else if (value.indexOf('-') > -1) {
    const a = value.split('-')
    cycle01.value = Number(a[0])
    cycle02.value = Number(a[1])
    radioValue.value = 3
  } else {
    checkboxList.value = [...new Set(value.split(',').map((i) => Number(i)))]
    radioValue.value = 4
  }
}

function onRadioChange() {
  // 日/周互斥：设置一方时，另一方置为“不指定”
  if (radioValue.value === 2 && props.cron.day === '?') emit('update', 'day', '*')
  if (radioValue.value !== 2 && props.cron.day !== '?') emit('update', 'day', '?')

  switch (radioValue.value) {
    case 1: emit('update', 'week', '*'); break
    case 2: emit('update', 'week', '?'); break
    case 3: emit('update', 'week', cycleTotal.value); break
    case 4:
      if (!checkboxList.value.length) checkboxList.value.push(checkCopy.value[0])
      else checkCopy.value = checkboxList.value
      emit('update', 'week', checkboxString.value)
      break
  }
}
</script>
