import { useVueFlow } from '@vue-flow/core'
import { ref, Ref } from 'vue'
import { TableNodeInfo } from '@/types/tableTypes'

interface State {
  draggedType: Ref<string | null>;
  isDragging: Ref<boolean>;
  nodeInfo: Ref<TableNodeInfo | null>;
}

const state: State = {
  draggedType: ref(null),
  isDragging: ref(false),
  nodeInfo: ref<TableNodeInfo | null>(null)
}

export default function useDragAndDrop(saveCurrentPageData?: () => void) {
  const { draggedType, isDragging,  nodeInfo } = state
  const { addNodes, screenToFlowCoordinate, onNodesInitialized, updateNode } = useVueFlow()

  function onDragStart(event: DragEvent, type: string, tableNodeInfo: TableNodeInfo): void {
    if (event.dataTransfer) {
      event.dataTransfer.setData('application/vueflow', type)
      event.dataTransfer.effectAllowed = 'move'
    }

    draggedType.value = type
    nodeInfo.value = tableNodeInfo
    isDragging.value = true

    document.addEventListener('drop', onDragEnd)
  }

  function onDragOver(event: DragEvent): void {
    event.preventDefault()

    if (draggedType.value && event.dataTransfer) {
      event.dataTransfer.dropEffect = 'move'
    }
  }

  function onDragEnd(): void {
    isDragging.value = false
    nodeInfo.value = null
    draggedType.value = null
    document.removeEventListener('drop', onDragEnd)
  }

  function onDrop(event: DragEvent): void {
    const position = screenToFlowCoordinate({
      x: event.clientX,
      y: event.clientY,
    })

    if (!nodeInfo.value) return;

    // 生成唯一ID，确保每次拖拽都创建新节点
    const uniqueId = `${nodeInfo.value.table.name}_${Date.now()}`

    // 创建深拷贝，确保每个节点有独立的数据副本
    const nodeInfoCopy = JSON.parse(JSON.stringify(nodeInfo.value)) as TableNodeInfo

    // 创建新节点
    const newNode = {
      id: uniqueId,
      type: 'table-node',
      position,
      data: nodeInfoCopy,
      dimensions: { width: 0, height: 0 },
      computedPosition: position,
      handleBounds: {
        source: [],
        target: []
      },
      isParent: false,
      selected: false,
      dragging: false,
      initialized: true
    }

    // 确保节点被添加到 VueFlow 实例
    try {
      addNodes(newNode)
      // 保存流程图数据
      saveCurrentPageData && saveCurrentPageData()
    } catch (error) {
      console.error('Failed to add node:', error)
    }

    // 调整节点位置
    const { off } = onNodesInitialized(() => {
      try {
        updateNode(newNode.id, (node) => ({
          position: {
            x: node.position.x - node.dimensions.width / 2,
            y: node.position.y - node.dimensions.height / 2
          },
        }))
      } catch (error) {
        console.error('Failed to update node position:', error)
      }
      off()
    })

    onDragEnd()
  }

  return {
    onDragStart,
    onDragOver,
    onDrop,
    isDragging,
  }
}
