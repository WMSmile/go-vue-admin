import request from '@/utils/request'

export const login = (data) => request({ url: '/auth/login', method: 'post', data })
export const logout = () => request({ url: '/auth/logout', method: 'post' })
export const getMe = () => request({ url: '/auth/me' })
export const getMenuTree = () => request({ url: '/menus/tree' })
export const getDashboard = () => request({ url: '/dashboard' })

export const listUsers = () => request({ url: '/users' })
export const createUser = (d) => request({ url: '/users', method: 'post', data: d })
export const updateUser = (id, d) => request({ url: `/users/${id}`, method: 'put', data: d })
export const deleteUser = (id) => request({ url: `/users/${id}`, method: 'delete' })

export const listRoles = () => request({ url: '/roles' })
export const createRole = (d) => request({ url: '/roles', method: 'post', data: d })
export const updateRole = (id, d) => request({ url: `/roles/${id}`, method: 'put', data: d })
export const deleteRole = (id) => request({ url: `/roles/${id}`, method: 'delete' })
export const assignMenus = (d) => request({ url: '/roles/menus', method: 'post', data: d })

export const listMenus = () => request({ url: '/menus' })
export const createMenu = (d) => request({ url: '/menus', method: 'post', data: d })
export const updateMenu = (id, d) => request({ url: `/menus/${id}`, method: 'put', data: d })
export const deleteMenu = (id) => request({ url: `/menus/${id}`, method: 'delete' })

export const listApiKeys = () => request({ url: '/apikeys' })
export const createApiKey = (d) => request({ url: '/apikeys', method: 'post', data: d })
export const deleteApiKey = (id) => request({ url: `/apikeys/${id}`, method: 'delete' })

export const listOperationLogs = (params) => request({ url: '/operation-logs', method: 'get', params })
export const deleteOperationLog = (id) => request({ url: `/operation-logs/${id}`, method: 'delete' })
export const clearOperationLogs = () => request({ url: '/operation-logs', method: 'delete' })

export const listDictTypes = (params) => request({ url: '/dict-types', method: 'get', params })
export const createDictType = (d) => request({ url: '/dict-types', method: 'post', data: d })
export const updateDictType = (id, d) => request({ url: `/dict-types/${id}`, method: 'put', data: d })
export const deleteDictType = (id) => request({ url: `/dict-types/${id}`, method: 'delete' })
export const listDictData = (params) => request({ url: '/dict-data', method: 'get', params })
export const createDictData = (d) => request({ url: '/dict-data', method: 'post', data: d })
export const updateDictData = (id, d) => request({ url: `/dict-data/${id}`, method: 'put', data: d })
export const deleteDictData = (id) => request({ url: `/dict-data/${id}`, method: 'delete' })

export const listTasks = (params) => request({ url: '/tasks', method: 'get', params })
export const createTask = (d) => request({ url: '/tasks', method: 'post', data: d })
export const updateTask = (id, d) => request({ url: `/tasks/${id}`, method: 'put', data: d })
export const deleteTask = (id) => request({ url: `/tasks/${id}`, method: 'delete' })
export const toggleTask = (id) => request({ url: `/tasks/${id}/toggle`, method: 'post' })
export const runTask = (id) => request({ url: `/tasks/${id}/run`, method: 'post' })
export const listTaskLogs = (params) => request({ url: '/task-logs', method: 'get', params })
