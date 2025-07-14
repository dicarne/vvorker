import axios from 'axios';
import { getToken, getUrl } from './config';
import { encryptData } from '../encrypt';

export const apiClient = axios.create();

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
    return config;
});
