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
