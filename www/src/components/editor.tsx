import React, { useLayoutEffect } from 'react'
import Editor, { loader } from '@monaco-editor/react'
import { useStore } from '@nanostores/react'
import { $code } from '@/store/workers'

loader.config({
  paths: {
    vs: 'monaco-assets/vs',
  },
})

export function MonacoEditor({ uid }: { uid: string }) {
  const code = useStore($code)
  useLayoutEffect(() => {}, [])
  return (
    <div className="flex-1">
      <Editor
        height="60vh"
        onChange={(v) => $code.set(v || '')}
        value={code}
        defaultLanguage="javascript"
      />
    </div>
  )
}
