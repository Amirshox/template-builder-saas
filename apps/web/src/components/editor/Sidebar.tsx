'use client'

import { useDraggable } from '@dnd-kit/core'
import { Type, Image, Keyboard } from 'lucide-react'

// Draggable Item Component
function SidebarItem({ type, icon: Icon, label, description }: { type: string; icon: any; label: string; description?: string }) {
    const { attributes, listeners, setNodeRef, transform } = useDraggable({
        id: `sidebar-${type}`,
        data: { type },
    })

    // Basic drag feedback style
    const style = transform ? {
        transform: `translate3d(${transform.x}px, ${transform.y}px, 0)`,
        zIndex: 1000,
        opacity: 0.9,
    } : undefined

    const handleClick = () => {
        // Check for TipTap instance
        const editor = (window as any).tiptapEditor
        if (editor && type === 'field') {
            editor.commands.insertContent({
                type: 'variable',
                attrs: { label: 'new_variable' }
            })
        }
    }

    return (
        <div
            ref={setNodeRef}
            {...listeners}
            {...attributes}
            onClick={handleClick}
            style={style}
            className="group flex flex-col gap-2 p-3 bg-white border border-zinc-200 rounded-xl cursor-grab hover:border-violet-300 hover:shadow-md hover:shadow-violet-100/50 active:cursor-grabbing transition-all duration-200"
        >
            <div className="flex items-center gap-3">
                <div className={`p-2 rounded-lg ${type === 'image' ? 'bg-emerald-50 text-emerald-600' : type === 'field' ? 'bg-amber-50 text-amber-600' : 'bg-violet-50 text-violet-600'}`}>
                    <Icon className="w-5 h-5" />
                </div>
                <div>
                    <span className="text-sm font-semibold text-zinc-700 block">{label}</span>
                </div>
                <div className="ml-auto opacity-0 group-hover:opacity-100 text-zinc-400">
                    :::
                </div>
            </div>
            {description && <p className="text-xs text-zinc-400 pl-1">{description}</p>}
        </div>
    )
}

export default function Sidebar() {
    return (
        <div className="w-72 bg-zinc-50/50 border-r border-zinc-200 p-5 shrink-0 flex flex-col gap-6 backdrop-blur-sm relative z-20">
            <div>
                <h2 className="text-xs font-bold text-zinc-400 uppercase tracking-widest mb-4 pl-1">Components</h2>
                <div className="flex flex-col gap-3">
                    <SidebarItem type="text" icon={Type} label="Text Block" description="Add static text content" />
                    <SidebarItem type="image" icon={Image} label="Image" description="Insert an image placeholder" />
                    <SidebarItem type="field" icon={Keyboard} label="Dynamic Field" description="Variable data from API" />
                </div>
            </div>

            <div className="mt-auto">
                <div className="bg-violet-50 rounded-xl p-4 border border-violet-100">
                    <p className="text-xs text-violet-700 text-center font-medium">
                        Drag items to the canvas to build your template.
                    </p>
                </div>
            </div>
        </div>
    )
}
