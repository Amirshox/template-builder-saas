import { Node, mergeAttributes } from '@tiptap/core'
import { ReactNodeViewRenderer, NodeViewWrapper } from '@tiptap/react'

export default Node.create({
    name: 'variable',

    group: 'inline',

    inline: true,

    atom: true,

    addAttributes() {
        return {
            id: {
                default: null,
            },
            label: {
                default: 'variable',
            },
        }
    },

    parseHTML() {
        return [
            {
                tag: 'span[data-variable]',
            },
        ]
    },

    renderHTML({ HTMLAttributes }) {
        return ['span', mergeAttributes(HTMLAttributes, { 'data-variable': '' }), 0]
    },

    addNodeView() {
        return ReactNodeViewRenderer(VariableComponent)
    },
})

function VariableComponent(props: any) {
    return (
        <NodeViewWrapper className="inline-block mx-1 align-middle">
            <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-violet-100 text-violet-800 border border-violet-200 select-none">
                {`{{ ${props.node.attrs.label} }}`}
            </span>
        </NodeViewWrapper>
    )
}
