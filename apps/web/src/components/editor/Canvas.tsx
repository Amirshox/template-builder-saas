'use client'

import { Stage, Layer, Rect, Text, Image as KonvaImage, Transformer } from 'react-konva'
import useImage from 'use-image'
import { useEditorStore } from '@/store/editor'
import { useEffect, useRef } from 'react'

const URLImage = ({ src, ...props }: any) => {
    const [image] = useImage(src || 'https://via.placeholder.com/150');
    return <KonvaImage image={image} {...props} />;
};

export default function Canvas() {
    const elements = useEditorStore((state) => state.elements)
    const selectedId = useEditorStore((state) => state.selectedId)
    const selectElement = useEditorStore((state) => state.selectElement)
    const updateElement = useEditorStore((state) => state.updateElement)

    const trRef = useRef<any>(null)

    // Selection Logic
    const checkDeselect = (e: any) => {
        const clickedOnEmpty = e.target === e.target.getStage();
        if (clickedOnEmpty) {
            selectElement(null);
        }
    };

    // Transformer (Selection Box) Logic
    useEffect(() => {
        if (selectedId && trRef.current) {
            // Find the node by id from the stage
            const stage = trRef.current.getStage();
            const selectedNode = stage.findOne('#' + selectedId);
            if (selectedNode) {
                trRef.current.nodes([selectedNode]);
                trRef.current.getLayer().batchDraw();
            }
        }
    }, [selectedId, elements]); // Re-run when selectedId or elements change (to update pos)

    return (
        <div className="flex-1 bg-gray-100 p-8 overflow-auto flex items-center justify-center">
            <div className="bg-white shadow-lg shadow-gray-200/50" style={{ width: 595, height: 842 }}> {/* A4 size at 72dpi approx, simplified */}
                <Stage
                    width={595}
                    height={842}
                    onMouseDown={checkDeselect}
                    onTouchStart={checkDeselect}
                >
                    <Layer>
                        {elements.map((el) => {
                            const isSelected = selectedId === el.id
                            const commonProps = {
                                id: el.id,
                                x: el.x,
                                y: el.y,
                                width: el.width,
                                height: el.height,
                                draggable: true,
                                onClick: () => selectElement(el.id),
                                onTap: () => selectElement(el.id),
                                onDragEnd: (e: any) => {
                                    updateElement(el.id, {
                                        x: e.target.x(),
                                        y: e.target.y(),
                                    })
                                },
                                onTransformEnd: (e: any) => {
                                    // Transformer changes scale, we want to normalize to width/height
                                    const node = e.target;
                                    const scaleX = node.scaleX();
                                    const scaleY = node.scaleY();

                                    // Reset scale
                                    node.scaleX(1);
                                    node.scaleY(1);

                                    updateElement(el.id, {
                                        x: node.x(),
                                        y: node.y(),
                                        width: Math.max(5, node.width() * scaleX),
                                        height: Math.max(5, node.height() * scaleY),
                                    });
                                }
                            }

                            if (el.type === 'text') {
                                return <Text
                                    key={el.id}
                                    {...commonProps}
                                    text={el.text}
                                    fontSize={el.style?.fontSize || 16}
                                    fontFamily={el.style?.fontFamily || 'Arial'}
                                    fill={el.style?.color || 'black'}
                                />
                            }

                            if (el.type === 'field') {
                                return <Text
                                    key={el.id}
                                    {...commonProps}
                                    text={`{{${el.fieldKey || 'field'}}}`}
                                    fontSize={el.style?.fontSize || 14}
                                    fontFamily="monospace"
                                    fill="blue"
                                />
                            }

                            if (el.type === 'image') {
                                return <URLImage key={el.id} {...commonProps} src={el.src} />
                            }

                            return null
                        })}

                        {/* Transformer Layer */}
                        {selectedId && (
                            <Transformer
                                ref={trRef}
                                boundBoxFunc={(oldBox, newBox) => {
                                    // Limit minimum size
                                    if (newBox.width < 5 || newBox.height < 5) {
                                        return oldBox;
                                    }
                                    return newBox;
                                }}
                            />
                        )}

                    </Layer>
                </Stage>
            </div>
        </div>
    )
}
