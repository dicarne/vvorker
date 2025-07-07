import axios from 'axios';
import { getToken, getUrl } from './config';


export const apiClient = axios.create();

apiClient.interceptors.request.use(config => {
    const token = getToken();
    const url = getUrl();
    if (token) {
        config.headers['Authorization'] = `Bearer ${token}`;
    }
    if (url) {
        config.baseURL = url;
    }
    return config;
});
