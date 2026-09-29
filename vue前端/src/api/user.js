import { get, post, postJSON } from '../utility/request'

// 获取所有学生
export const getAllStudents = (success, failure , error) => {
    get('/api/user/all-student', null, success, failure)
} 