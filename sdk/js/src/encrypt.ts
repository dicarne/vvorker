import { config } from "./common/common";

const encoder = new TextEncoder();

/**
 * Encrypts data using AES-GCM
 * @param data The data to encrypt (string or object)
 * @param key The encryption key as a string (16, 24, or 32 bytes for AES-128/192/256)
 * @returns Base64 encoded encrypted data
 */
export async function encryptData<T>(data: T, key: string): Promise<string> {
  const dataStr = typeof data === "string" ? data : JSON.stringify(data);

  const iv = crypto.getRandomValues(new Uint8Array(12));

  const cryptoKey = await crypto.subtle.importKey(
    "raw",
    encoder.encode(key),
    { name: "AES-GCM" },
    false,
    ["encrypt"]
  );

  const encrypted = await crypto.subtle.encrypt(
    { name: "AES-GCM", iv },
    cryptoKey,
    encoder.encode(dataStr)
  );

  const result = new Uint8Array(iv.length + encrypted.byteLength);
  result.set(new Uint8Array(iv), 0);
  result.set(new Uint8Array(encrypted), iv.length);

  return btoa(
    Array.from(result, (byte) => String.fromCharCode(byte)).join("")
  );
}

/**
 * Encrypted fetch wrapper.
 * Reads the encryption key from VVORKER_ENCRYPTION_KEY environment variable,
 * encrypts the request body, and sets the X-Encrypted-Data header.
 *
 * Drop-in replacement for fetch() — use this to transparently encrypt requests.
 *
 * @example
 * ```ts
 * // Before:
 * const res = await fetch("/api/data", { method: "POST", body: JSON.stringify(payload) });
 *
 * // After:
 * const res = await encryptedFetch("/api/data", { method: "POST", body: JSON.stringify(payload) });
 * ```
 */
function getEncryptionKey(): string | undefined {
  const g = globalThis as any;
  return (
    g.VITE_VVORKER_ENCRYPTION_KEY ??
    g.process?.env?.VITE_VVORKER_ENCRYPTION_KEY ??
    (import.meta as any).env?.VITE_VVORKER_ENCRYPTION_KEY
  );
}

/**
 * Debug endpoint fetch with optional encryption.
 * When VVORKER_ENCRYPTION_KEY is set, automatically encrypts the request body.
 */
export async function debugFetch(
  body: Record<string, any>
): Promise<Response> {
  const key = getEncryptionKey();
  const url = `${config().url}/__vvorker__debug`;
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${config().token}`,
  };

  let payload: string;
  const jsonBody = JSON.stringify(body);

  if (key) {
    payload = await encryptData(jsonBody, key);
    headers["X-Encrypted-Data"] = "true";
  } else {
    payload = jsonBody;
  }

  return fetch(url, { method: "POST", headers, body: payload });
}

export async function encryptedFetch(
  url: string,
  init?: RequestInit
): Promise<Response> {
  const key = getEncryptionKey();

  if (!key) {
    throw new Error(
      "VITE_VVORKER_ENCRYPTION_KEY is not set. Set it as an environment variable."
    );
  }

  let body = init?.body;

  if (body !== undefined && body !== null) {
    let text: string;
    if (typeof body === "string") {
      text = body;
    } else if (body instanceof ArrayBuffer) {
      text = new TextDecoder().decode(body);
    } else if (ArrayBuffer.isView(body)) {
      text = new TextDecoder().decode(body.buffer);
    } else {
      // FormData, Blob, URLSearchParams, ReadableStream etc. — cannot encrypt
      throw new Error(
        "encryptedFetch only supports string, ArrayBuffer, or ArrayBufferView bodies"
      );
    }
    body = await encryptData(text, key);
  }

  return fetch(url, {
    ...init,
    body,
    headers: {
      ...init?.headers,
      "X-Encrypted-Data": "true",
    },
  });
}
