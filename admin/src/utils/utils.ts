export const decodeBase64 = (base64String: string | undefined) => {
  if (!base64String) {
    return ''
  }
  try {
    // 使用 atob 解码 Base64 字符串
    const binaryString = atob(base64String)
    // 将二进制字符串转换为 Uint8Array
    const bytes = new Uint8Array(binaryString.length)
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i)
    }
    // 使用 TextDecoder 处理 Unicode 字符
    return new TextDecoder().decode(bytes)
  } catch (error) {
    console.error('Base64 解码失败:', error)
    return ''
  }
}

export const copyContent = (content: string) => {
  const input = document.createElement('input')
  input.value = content
  document.body.appendChild(input)
  input.select()
  document.execCommand('copy')
  document.body.removeChild(input)
}

export const formatDate = (date: Date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}
