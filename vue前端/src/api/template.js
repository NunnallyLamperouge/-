import { get, post, postJSON } from '../utility/request'

// 获取模板列表
export const getTemplateList = (success, failure) => {
    get('/api/template/list', null, success, failure)
}
export const getTemplateListStudent = (success, failure) => {
    get('/api/template/list-student', null, success, failure)
}

// 获取模板详情
export const getTemplateDetail = (id, success, failure) => {
    get(`/api/template/detail/${id}`, null, success, failure)
}

// 创建模板
export const createTemplate = (data, success, failure) => {
    postJSON('/api/template/create', data, success, failure)
}

// 更新模板
export const updateTemplate = (id, data, success, failure) => {
    postJSON(`/api/template/update/${id}`, data, success, failure)
}

// 删除模板
export const deleteTemplate = (id, success, failure) => {
    post(`/api/template/delete/${id}`, null, success, failure)
}

// 添加组件
export const addComponent = (templateId, data, success, failure) => {
    postJSON(`/api/template/${templateId}/components`, data, success, failure)
}

// 更新组件
export const updateComponent = (componentId, data, success, failure) => {
    postJSON(`/api/template/components/${componentId}/update`, data, success, failure)
}

// 删除组件
export const deleteComponent = (componentId, success, failure) => {
    post(`/api/template/components/${componentId}/delete`, null, success, failure)
}

// 调整组件顺序
export const reorderComponents = (templateId, data, success, failure) => {
    postJSON(`/api/template/${templateId}/components/reorder`, data, success, failure)
}

// 导入模板
export const importTemplate = (data, success, failure) => {
    postJSON('/api/template/import', data, success, failure)
}

// 导出模板
export const exportTemplate = (id, success, failure) => {
    get(`/api/template/detail/${id}`, null, success, failure)
}

// 发送模板给学生
export const sendTemplateToStudents = (code, studentIds, success, failure) => {
    postJSON('/api/template/send-template-out', { code: code, studentIds: studentIds }, success, failure)
}

// Fork模板
export const forkTemplate = (data, success, failure) => {
    postJSON('/api/template/fork', data, success, failure)
} 