import axios from 'axios';
import { getToken, getUrl } from './config';
import { encryptData } from '../encrypt';
import inquirer from 'inquirer';

export const apiClient = axios.create();
let otptoken = ""
apiClient.interceptors.request.use(async config => {
    const token = getToken();
    const url = getUrl();
    if (token) {
        config.headers['Authorization'] = `Bearer ${token}`;
    }
    if (url) {
        config.baseURL = url;
    }
    if (config.headers['x-encrypted-data']) {
        config.data = await encryptData(config.data, config.headers['x-encrypted-data']);
        config.headers['x-encrypted-data'] = "true"
    }
    if (otptoken) {
        config.headers['vv-otp-token'] = otptoken
    }
    return config;
});

// catch all error
apiClient.interceptors.response.use(response => {
    return response;
}, error => {
    console.error('api Error')
    const body = error.response.data
    if (!body) {
        console.error(error)
        return Promise.reject(error)
    }
    console.error(body)
    return Promise.reject(body)
});

export async function requireOTP() {
    const response = await apiClient.post('/api/otp/is-enable')
    if (!response.data.data.enabled) return
    const { otp } = await inquirer.prompt([{
        type: 'input',
        name: 'otp',
        message: '请输入OTP代码:',
    }])
    const response2 = await apiClient.post('/api/otp/valid?code=' + otp)
    if (response2.data.code !== 0) {
        throw new Error('OTP验证失败')
    }
    otptoken = response2.data.data["vv-otp-token"]
}