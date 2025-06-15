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

export const deleteUser = (userId: string) => {
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

export const changePassword = (userId: number, password: string) => {
    return http.post(`/api/admin/users/${userId}`, {
        password
    })
}