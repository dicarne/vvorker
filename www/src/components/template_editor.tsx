import React, { useLayoutEffect } from 'react'
import Editor, { loader } from '@monaco-editor/react'

loader.config({
    paths: {
        vs: 'monaco-assets/vs',
    },
})

export interface TemplateEditorProps {
    setContent: (content: string) => void
    content: string
}


export function TemplateEditor({ setContent, content }: TemplateEditorProps) {
    useLayoutEffect(() => { }, [])
    return (
        <div className="flex-1">
            <Editor
                height="60vh"
                onChange={(v) => setContent(v || '')}
                value={content}
                defaultLanguage="json"
            />
        </div>
    )
}
