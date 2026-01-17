import { create } from 'zustand'
import { immer } from 'zustand/middleware/immer'

export interface EditorElement {
    id: string
    type: 'text' | 'image' | 'field'
    x: number
    y: number
    width: number
    height: number
    text?: string
    fieldKey?: string
    style?: Record<string, any>
    src?: string // for images
}

interface EditorState {
    elements: EditorElement[]
    selectedId: string | null
    scale: number

    addElement: (element: EditorElement) => void
    updateElement: (id: string, updates: Partial<EditorElement>) => void
    selectElement: (id: string | null) => void
    setElements: (elements: EditorElement[]) => void
}

export const useEditorStore = create<EditorState>()(
    immer((set) => ({
        elements: [],
        selectedId: null,
        scale: 1,

        addElement: (element) =>
            set((state) => {
                state.elements.push(element)
                state.selectedId = element.id
            }),

        updateElement: (id, updates) =>
            set((state) => {
                const index = state.elements.findIndex((el) => el.id === id)
                if (index !== -1) {
                    state.elements[index] = { ...state.elements[index], ...updates }
                }
            }),

        selectElement: (id) =>
            set((state) => {
                state.selectedId = id
            }),

        setElements: (elements) =>
            set((state) => {
                state.elements = elements
            })
    }))
)
