'use client'

import { useParams } from 'next/navigation'
import Sidebar from '@/components/editor/Sidebar'
import PropertiesPanel from '@/components/editor/PropertiesPanel'
import dynamic from 'next/dynamic'
import { DndContext, DragEndEvent } from '@dnd-kit/core'
import { useEditorStore } from '@/store/editor'

// Import Canvas dynamically to avoid SSR issues with Konva
const Canvas = dynamic(() => import('@/components/editor/Canvas'), {
    ssr: false,
    loading: () => <div className="flex-1 bg-gray-100 flex items-center justify-center">Loading Canvas...</div>
})

import { downloadPreview, createVersion, getTemplate } from '@/lib/api'
import { useState, useEffect } from 'react'
import DocumentEditor from '@/components/editor/DocumentEditor'

export default function EditorPage() {
    const { id } = useParams()
    const addElement = useEditorStore((state) => state.addElement)
    const elements = useEditorStore((state) => state.elements)
    const setElements = useEditorStore((state) => state.setElements)

    const [template, setTemplate] = useState<any>(null)
    const [documentContent, setDocumentContent] = useState<any>(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        if (typeof id === 'string') {
            getTemplate(id).then(t => {
                setTemplate(t)
                setLoading(false)
                // If we had a way to load version content, we'd do it here.
                // For now, new docs start empty or could load from a version logic.
            }).catch(err => {
                console.error("Failed to load template", err)
                setLoading(false)
            })
        }
    }, [id])

    const handleDragEnd = (event: DragEndEvent) => {
        const { active, over } = event

        if (over && active.data.current) {
            // Dropped on canvas
            const type = active.data.current.type

            // Simple positioning refinement:
            // dnd-kit gives us the transformed coordinates of the dragged overlay
            // We can try to use that if we knew the canvas offset.
            // For now, let's just randomize slightly so they don't stack perfectly effectively hiding new ones.
            addElement({
                id: crypto.randomUUID(),
                type: type as any,
                x: 100 + (elements.length * 10),
                y: 100 + (elements.length * 10),
                width: type === 'image' ? 100 : 200,
                height: type === 'image' ? 100 : 40,
                text: type === 'text' ? 'Double click to edit' : undefined,
                fieldKey: type === 'field' ? 'new.field' : undefined,
                src: type === 'image' ? 'https://via.placeholder.com/150' : undefined
            })
        }
    }

    const handleSave = async () => {
        if (typeof id === 'string') {
            try {
                let dataToSave = {}
                if (template?.type === 'docx' || template?.type === 'document') {
                    dataToSave = documentContent // Save TiTap JSON
                } else {
                    dataToSave = { elements } // Save Layout JSON
                }

                await createVersion(id, dataToSave)
                // Maybe show a toast
                console.log("Saved successfully")
            } catch (e) {
                console.error("Save failed", e)
                alert("Failed to save draft.")
            }
        }
    }

    const handlePreview = async () => {
        if (typeof id === 'string') {
            try {
                // Auto-save before preview
                await handleSave()
                await downloadPreview(id)
            } catch (e) {
                console.error("Preview failed", e)
                alert("Preview failed. Make sure backend and renderer are running.")
            }
        }
    }

    if (loading) {
        return <div className="flex h-screen w-screen items-center justify-center">Loading...</div>
    }

    const isDocumentMode = template?.type === 'docx' || template?.type === 'document';

    return (
        <DndContext id="editor-dnd-context" onDragEnd={handleDragEnd}>
            <div className="flex h-screen w-screen overflow-hidden">
                <Sidebar />

                {isDocumentMode ? (
                    <div className="flex-1 flex flex-col relative z-0 h-full bg-[#FAFAFA]">
                        <CanvasHeader title={template?.name} type="Document Editor" onSave={handleSave} onPreview={handlePreview} />
                        <DocumentEditor content={documentContent} onChange={setDocumentContent} />
                    </div>
                ) : (
                    <CanvasWrapper title={template?.name} onPreview={handlePreview} onSave={handleSave} />
                )}

                {!isDocumentMode && <PropertiesPanel />}
            </div>
        </DndContext>
    )
}

import { useDroppable } from '@dnd-kit/core'
import { Save, Eye, FileUp, ChevronLeft } from 'lucide-react'
import Link from 'next/link'

function CanvasHeader({ title, type = "Layout Editor", onPreview, onSave }: { title?: string, type?: string, onPreview: () => void, onSave: () => void }) {
    return (
        <header className="h-16 bg-white border-b border-zinc-200 flex items-center px-6 justify-between shrink-0 shadow-sm z-10">
            <div className="flex items-center gap-4">
                <Link href="/" className="p-2 -ml-2 text-zinc-400 hover:text-zinc-600 rounded-full hover:bg-zinc-100 transition-colors">
                    <ChevronLeft className="w-5 h-5" />
                </Link>
                <div>
                    <h1 className="font-semibold text-zinc-900 leading-tight">{title || 'Untitled Template'}</h1>
                    <p className="text-xs text-zinc-400">{type}</p>
                </div>
            </div>

            <div className="flex gap-3">
                <button
                    onClick={onSave}
                    className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-zinc-600 bg-white border border-zinc-200 rounded-lg hover:bg-zinc-50 hover:border-zinc-300 transition-all focus:ring-2 focus:ring-zinc-200 focus:outline-none"
                >
                    <Save className="w-4 h-4" />
                    Save Draft
                </button>
                <button
                    onClick={onPreview}
                    className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-violet-600 bg-violet-50 border border-violet-200 rounded-lg hover:bg-violet-100 hover:border-violet-300 transition-all focus:ring-2 focus:ring-violet-200 focus:outline-none"
                >
                    <Eye className="w-4 h-4" />
                    Preview
                </button>
                <button className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-gradient-to-r from-violet-600 to-indigo-600 rounded-lg hover:from-violet-500 hover:to-indigo-500 shadow-md shadow-violet-500/20 transition-all focus:ring-2 focus:ring-violet-500 focus:ring-offset-2 focus:outline-none">
                    <FileUp className="w-4 h-4" />
                    Publish
                </button>
            </div>
        </header>
    )
}

function CanvasWrapper({ title, onPreview, onSave }: { title?: string, onPreview: () => void, onSave: () => void }) {
    const { setNodeRef } = useDroppable({
        id: 'canvas-drop-zone',
    })

    return (
        <div ref={setNodeRef} className="flex-1 flex flex-col relative z-0 h-full bg-[#FAFAFA]">
            <CanvasHeader title={title} onPreview={onPreview} onSave={onSave} />

            {/* Canvas Container with Dot Pattern */}
            <div className="flex-1 overflow-auto flex items-center justify-center p-8 bg-[radial-gradient(#e5e7eb_1px,transparent_1px)] [background-size:20px_20px]">
                <div className="shadow-2xl shadow-zinc-900/10 ring-1 ring-zinc-900/5 transition-transform duration-300 ease-in-out">
                    <Canvas />
                </div>
            </div>
        </div>
    )
}
