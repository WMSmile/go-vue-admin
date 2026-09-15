import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

const service = axios.create({
  baseURL: '/api/v1',
  timeout: 10000
})

service.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

service.interceptors.response.use(
  (resp) => {
    const data = resp.data
    if (data.code === 0) return data.data
    if (data.code === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('menus')
      router.push('/login')
      return Promise.reject(data.msg)
    }
    ElMessage.error(data.msg || '请求失败')
    return Promise.reject(data.msg)
  },
  (error) => {
    ElMessage.error(error.message || '网络错误')
    return Promise.reject(error)
  }
)

export default service
