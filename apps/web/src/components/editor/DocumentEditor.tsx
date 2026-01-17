'use client'

import { useEditor, EditorContent } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import VariableNode from './extensions/VariableNode'
import { Bold, Italic, List, ListOrdered, Heading1, Heading2 } from 'lucide-react'
import { useEffect } from 'react'

const MenuBar = ({ editor }: { editor: any }) => {
    if (!editor) {
        return null
    }

    return (
        <div className="flex items-center gap-1 p-2 bg-white border-b border-zinc-200 sticky top-0 z-10">
            <button
                onClick={() => editor.chain().focus().toggleBold().run()}
                className={`p-1.5 rounded hover:bg-zinc-100 ${editor.isActive('bold') ? 'bg-zinc-100 text-violet-600' : 'text-zinc-600'}`}
                title="Bold"
            >
                <Bold className="w-4 h-4" />
            </button>
            <button
                onClick={() => editor.chain().focus().toggleItalic().run()}
                className={`p-1.5 rounded hover:bg-zinc-100 ${editor.isActive('italic') ? 'bg-zinc-100 text-violet-600' : 'text-zinc-600'}`}
                title="Italic"
            >
                <Italic className="w-4 h-4" />
            </button>
            <div className="w-px h-4 bg-zinc-200 mx-1" />
            <button
                onClick={() => editor.chain().focus().toggleHeading({ level: 1 }).run()}
                className={`p-1.5 rounded hover:bg-zinc-100 ${editor.isActive('heading', { level: 1 }) ? 'bg-zinc-100 text-violet-600' : 'text-zinc-600'}`}
                title="Heading 1"
            >
                <Heading1 className="w-4 h-4" />
            </button>
            <button
                onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
                className={`p-1.5 rounded hover:bg-zinc-100 ${editor.isActive('heading', { level: 2 }) ? 'bg-zinc-100 text-violet-600' : 'text-zinc-600'}`}
                title="Heading 2"
            >
                <Heading2 className="w-4 h-4" />
            </button>
            <div className="w-px h-4 bg-zinc-200 mx-1" />
            <button
                onClick={() => editor.chain().focus().toggleBulletList().run()}
                className={`p-1.5 rounded hover:bg-zinc-100 ${editor.isActive('bulletList') ? 'bg-zinc-100 text-violet-600' : 'text-zinc-600'}`}
                title="Bullet List"
            >
                <List className="w-4 h-4" />
            </button>
            <button
                onClick={() => editor.chain().focus().toggleOrderedList().run()}
                className={`p-1.5 rounded hover:bg-zinc-100 ${editor.isActive('orderedList') ? 'bg-zinc-100 text-violet-600' : 'text-zinc-600'}`}
                title="Ordered List"
            >
                <ListOrdered className="w-4 h-4" />
            </button>
        </div>
    )
}

export default function DocumentEditor({ content, onChange }: { content: any, onChange: (content: any) => void }) {
    const editor = useEditor({
        extensions: [
            StarterKit,
            Placeholder.configure({
                placeholder: 'Start typing your document...',
            }),
            VariableNode,
        ],
        content: content,
        editorProps: {
            attributes: {
                class: 'prose prose-sm sm:prose-base lg:prose-lg xl:prose-xl m-5 focus:outline-none min-h-[500px] max-w-none',
            },
        },
        onUpdate: ({ editor }) => {
            onChange(editor.getJSON())
        },
    })

    // Expose editor instance to window for Sidebar to access (a bit hacky but effective for MVP)
    useEffect(() => {
        if (editor) {
            (window as any).tiptapEditor = editor
        }
        return () => {
            (window as any).tiptapEditor = undefined
        }
    }, [editor])

    return (
        <div className="flex flex-col h-full bg-white">
            <MenuBar editor={editor} />
            <div className="flex-1 overflow-y-auto bg-zinc-50 p-8 flex justify-center cursor-text" onClick={() => editor?.chain().focus().run()}>
                <div className="w-full max-w-[816px] min-h-[1056px] bg-white shadow-lg p-12 mb-8" onClick={(e) => e.stopPropagation()}>
                    <EditorContent editor={editor} />
                </div>
            </div>
        </div>
    )
}
