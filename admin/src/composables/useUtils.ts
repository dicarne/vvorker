import { useMessage } from 'naive-ui'
import { copyContent as _copyContent } from '@/utils/utils'

export function useCopyContent() {
  const message = useMessage()
  const copyContent = (content: string) => {
    try {
      _copyContent(content)
      message.success('复制成功')
    } catch (error) {
      message.error('复制失败')
    }
  }
  return {
    copyContent,
  }
}
