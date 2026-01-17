'use client'

import { useEditorStore } from '@/store/editor'

// Helper for input components
const InputGroup = ({ label, children }: { label: string; children: React.ReactNode }) => (
    <div className="flex flex-col gap-1.5">
        <label className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest">{label}</label>
        {children}
    </div>
)

const Input = (props: React.InputHTMLAttributes<HTMLInputElement>) => (
    <input
        {...props}
        className="w-full text-sm bg-zinc-50 border border-zinc-200 rounded-lg px-3 py-2 text-zinc-700 placeholder:text-zinc-400 focus:outline-none focus:ring-2 focus:ring-violet-500/20 focus:border-violet-500 transition-all duration-200"
    />
)

const TextArea = (props: React.TextareaHTMLAttributes<HTMLTextAreaElement>) => (
    <textarea
        {...props}
        className="w-full text-sm bg-zinc-50 border border-zinc-200 rounded-lg px-3 py-2 text-zinc-700 placeholder:text-zinc-400 focus:outline-none focus:ring-2 focus:ring-violet-500/20 focus:border-violet-500 transition-all duration-200"
    />
)

export default function PropertiesPanel() {
    const selectedId = useEditorStore((state) => state.selectedId)
    const elements = useEditorStore((state) => state.elements)
    const updateElement = useEditorStore((state) => state.updateElement)

    const selectedElement = elements.find((el) => el.id === selectedId)

    if (!selectedElement) {
        return (
            <div className="w-80 bg-white border-l border-zinc-200 p-8 shrink-0 flex flex-col items-center justify-center text-center gap-4">
                <div className="w-16 h-16 bg-zinc-50 rounded-full flex items-center justify-center">
                    <div className="w-8 h-8 rounded bg-zinc-200/50" />
                </div>
                <div>
                    <h3 className="font-medium text-zinc-900 mb-1">No Selection</h3>
                    <p className="text-sm text-zinc-500">Select an element on the canvas to edit its properties.</p>
                </div>
            </div>
        )
    }

    const handleChange = (key: string, value: any) => {
        updateElement(selectedElement.id, { [key]: value })
    }

    const handleStyleChange = (key: string, value: any) => {
        // For MVP assume style is flat
        updateElement(selectedElement.id, {
            style: { ...selectedElement.style, [key]: value }
        })
    }

    return (
        <div className="w-80 bg-white border-l border-zinc-200 flex flex-col h-full">
            <div className="p-4 border-b border-zinc-100 bg-zinc-50/30">
                <h2 className="font-semibold text-zinc-900">Properties</h2>
                <div className="flex items-center gap-2 mt-1">
                    <span className="text-[10px] font-mono bg-zinc-100 text-zinc-500 px-1.5 py-0.5 rounded border border-zinc-200">
                        {selectedElement.type.toUpperCase()}
                    </span>
                    <span className="text-[10px] font-mono text-zinc-400 truncate max-w-[150px]">
                        {selectedElement.id}
                    </span>
                </div>
            </div>

            <div className="p-6 flex flex-col gap-6 overflow-y-auto flex-1">
                {/* Position Group */}
                <div className="space-y-4">
                    <h3 className="text-sm font-medium text-zinc-900 pb-2 border-b border-zinc-100">Layout</h3>
                    <div className="grid grid-cols-2 gap-3">
                        <InputGroup label="X Position">
                            <Input
                                type="number"
                                value={selectedElement.x}
                                onChange={(e) => handleChange('x', Number(e.target.value))}
                            />
                        </InputGroup>
                        <InputGroup label="Y Position">
                            <Input
                                type="number"
                                value={selectedElement.y}
                                onChange={(e) => handleChange('y', Number(e.target.value))}
                            />
                        </InputGroup>
                    </div>

                    <div className="grid grid-cols-2 gap-3">
                        <InputGroup label="Width">
                            <Input
                                type="number"
                                value={selectedElement.width}
                                onChange={(e) => handleChange('width', Number(e.target.value))}
                            />
                        </InputGroup>
                        <InputGroup label="Height">
                            <Input
                                type="number"
                                value={selectedElement.height}
                                onChange={(e) => handleChange('height', Number(e.target.value))}
                            />
                        </InputGroup>
                    </div>
                </div>

                {/* Content Group */}
                <div className="space-y-4">
                    <h3 className="text-sm font-medium text-zinc-900 pb-2 border-b border-zinc-100">Content</h3>

                    {selectedElement.type === 'text' && (
                        <>
                            <InputGroup label="Text Content">
                                <TextArea
                                    value={selectedElement.text}
                                    onChange={(e) => handleChange('text', e.target.value)}
                                    rows={4}
                                />
                            </InputGroup>

                            <InputGroup label="Font Size">
                                <Input
                                    type="number"
                                    value={selectedElement.style?.fontSize || 16}
                                    onChange={(e) => handleStyleChange('fontSize', Number(e.target.value))}
                                />
                            </InputGroup>
                        </>
                    )}

                    {selectedElement.type === 'field' && (
                        <InputGroup label="Field Key">
                            <Input
                                type="text"
                                value={selectedElement.fieldKey || ''}
                                onChange={(e) => handleChange('fieldKey', e.target.value)}
                                placeholder="e.g. user.name"
                            />
                            <p className="text-[10px] text-zinc-400 mt-1">Key used for data binding during generation</p>
                        </InputGroup>
                    )}

                    {selectedElement.type === 'image' && (
                        <InputGroup label="Image Source">
                            <Input
                                type="text"
                                value={selectedElement.src || ''}
                                onChange={(e) => handleChange('src', e.target.value)}
                                placeholder="https://..."
                            />
                            <p className="text-[10px] text-zinc-400 mt-1">URL or Asset ID</p>
                        </InputGroup>
                    )}
                </div>
            </div>
        </div>
    )
}
