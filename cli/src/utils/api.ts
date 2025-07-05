import axios from 'axios';
import { getToken, getUrl } from './config';
import { createPublicKey, publicEncrypt } from 'crypto';

export async function encryptData(data: any, vk: string) {
    const publicKey = createPublicKey({
        key: Buffer.from(vk, 'base64'),
        format: 'der',
        type: 'spki',
    });
    const encryptedData = publicEncrypt(
        {
            key: publicKey,
            padding: 4, // RSA_PKCS1_OAEP_PADDING
        },
        Buffer.from(JSON.stringify(data))
    );
    return encryptedData.toString('base64');
}

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
