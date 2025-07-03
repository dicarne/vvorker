// Text encoder/decoder for string conversion
const encoder = new TextEncoder();
const decoder = new TextDecoder();

/**
 * Encrypts data using AES-GCM
 * @param data The data to encrypt (string or object)
 * @param key The encryption key as a string
 * @returns Base64 encoded encrypted data
 */
export async function encryptData<T>(data: T, key: string): Promise<string> {
  try {
    // Convert data to string if it's an object
    const dataStr = typeof data === 'string' ? data : JSON.stringify(data);

    // Generate a random IV (Initialization Vector)
    const iv = crypto.getRandomValues(new Uint8Array(12));

    // Import the key
    const cryptoKey = await window.crypto.subtle.importKey(
      'raw',
      encoder.encode(key),
      { name: 'AES-GCM' },
      false,
      ['encrypt']
    );

    // Encrypt the data
    const encrypted = await window.crypto.subtle.encrypt(
      { name: 'AES-GCM', iv },
      cryptoKey,
      encoder.encode(dataStr)
    );

    // Combine IV and encrypted data
    const result = new Uint8Array(iv.length + encrypted.byteLength);
    result.set(new Uint8Array(iv), 0);
    result.set(new Uint8Array(encrypted), iv.length);

    // Convert to base64
    return btoa(Array.from(result, byte => String.fromCharCode(byte)).join(''));
  } catch (error) {
    console.error('Encryption error:', error);
    throw new Error('Failed to encrypt data');
  }
}

/**
 * Decrypts data using AES-GCM
 * @param encryptedData Base64 encoded encrypted data
 * @param key The decryption key as a string
 * @returns Decrypted data as string or parsed object if it was an object
 */
export async function decryptData<T = any>(encryptedData: string, key: string): Promise<T> {
  try {
    // Convert from base64
    const binaryString = atob(encryptedData);
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i);
    }

    // Extract IV (first 12 bytes for GCM)
    const iv = bytes.slice(0, 12);
    const data = bytes.slice(12);

    // Import the key
    const cryptoKey = await window.crypto.subtle.importKey(
      'raw',
      encoder.encode(key),
      { name: 'AES-GCM' },
      false,
      ['decrypt']
    );

    // Decrypt the data
    const decrypted = await window.crypto.subtle.decrypt(
      { name: 'AES-GCM', iv },
      cryptoKey,
      data
    );

    // Convert to string
    const decryptedStr = decoder.decode(decrypted);

    // Try to parse as JSON, return as string if it fails
    try {
      return JSON.parse(decryptedStr) as T;
    } catch {
      return decryptedStr as unknown as T;
    }
  } catch (error) {
    console.error('Decryption error:', error);
    throw new Error('Failed to decrypt data');
  }
}

/**
 * Creates an encrypted request interceptor for axios
 * @param axiosInstance The axios instance to add the interceptor to
 * @param encryptionKey The encryption key
 */
export function setupAxiosEncryption(axiosInstance: any, encryptionKey: string) {
  if (!encryptionKey || encryptionKey == "") return
  // Request interceptor for encrypting request data
  axiosInstance.interceptors.request.use(async (config: any) => {
    // Only encrypt POST, PUT, PATCH requests with data
    if (['post', 'put', 'patch'].includes(config.method?.toLowerCase()) && config.data && config.headers['X-Encrypted-Data'] === 'true') {
      try {
        const encryptedData = await encryptData(config.data, encryptionKey);
        config.data = encryptedData;
      } catch (error) {
        console.error('Request encryption failed:', error);
        throw error;
      }
    }
    return config;
  }, (error: any) => {
    return Promise.reject(error);
  });

  // Response interceptor for decrypting response data
  axiosInstance.interceptors.response.use(
    async (response: any) => {
      // Only try to decrypt if the response is encrypted
      if (response.headers['x-encrypted-data'] === 'true' && response.data) {
        try {
          const decryptedData = await decryptData(response.data, encryptionKey);
          response.data = decryptedData;
        } catch (error) {
          console.error('Response decryption failed:', error);
          throw error;
        }
      }
      return response;
    },
    (error: any) => {
      return Promise.reject(error);
    }
  );
}

/**
 * Gets the encryption key from the user's data
 * @param user The user object containing the encryption key
 * @returns The encryption key or null if not available
 */
export function getEncryptionKey(user: { vk?: string } | null): string | null {
  return user?.vk || null;
}