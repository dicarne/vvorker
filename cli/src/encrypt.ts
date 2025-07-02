// Text encoder/decoder for string conversion
const encoder = new TextEncoder();

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
    const cryptoKey = await crypto.subtle.importKey(
      'raw',
      encoder.encode(key),
      { name: 'AES-GCM' },
      false,
      ['encrypt']
    );

    // Encrypt the data
    const encrypted = await crypto.subtle.encrypt(
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
