import http from './http'

interface CreateUserData {
    username: string
    email: string
    password: string
}

export const getUsers = () => {
    return http.get('/api/admin/users')
}

export const createUser = (data: CreateUserData) => {
    return http.post('/api/admin/users', data)
}

export const deleteUser = (userId: number) => {
    return http.delete(`/api/admin/users/${userId}`)
}

export const banUser = (userId: number) => {
    return http.post(`/api/admin/users/${userId}`, {
        status: 1
    })
}

export const unbanUser = (userId: number) => {
    return http.post(`/api/admin/users/${userId}`, {
        status: 2
    })
}

// 可以用于管理员修改密码或用户修改自己的密码。如果是用户自己修改，需要传入自己的id
export const changePassword = (userId: number, password: string) => {
    return http.post(`/api/admin/users/${userId}`, {
        password
    })
}