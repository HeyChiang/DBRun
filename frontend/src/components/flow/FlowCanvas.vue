<script setup lang="ts">
import { computed, onMounted, inject, onBeforeUnmount, watch, ref } from 'vue'
import { VueFlow, useVueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import TableNode from '@/components/flow/TableNode.vue'
import useDragAndDrop from '@/components/sidebar/useDragAndDrop'
import { usePageStore } from '@/stores/pageStore'
import { flowDataStore } from '@/stores/flowDataStore'
import EdgeWithButton from './edge/EdgeWithButton.vue'
import RelationshipSelector from './edge/RelationshipSelector.vue'
import FlowSearch from './search/FlowSearch.vue'
import { useRoute } from 'vue-router'
import { TableNodeInfo } from '@/types/tableTypes'
import { eventBus } from '@/utils/eventBus'

// 接收路由传入的 pageId，以避免“非 props 属性”警告
const props = defineProps<{ pageId?: string }>()

// 注入由App.vue提供的showTableDialog方法
const showTableDialog = inject<(tableData: TableNodeInfo) => void>('showTableDialog')

const route = useRoute()
const pageStore = usePageStore()

// 保存当前页面节点数据
const saveNodeData = (): void => {
  const pageId = route.params.pageId as string || pageStore.activeTab
  if (!pageId) return

  flowDataStore.saveNodeData(pageId, visibleNodes.value)
}


// 保存当前页面边缘数据
const saveEdgeData = (): void => {
  const pageId = route.params.pageId as string || pageStore.activeTab
  if (!pageId) return

  flowDataStore.saveEdgeData(pageId, visibleEdges.value)
}

const { onDragOver, onDrop } = useDragAndDrop(saveNodeData)


// VueFlow 实例
const {
  onConnect,
  addEdges,
  getEdges,
  getNodes,
  setNodes,
  setEdges,
  onNodeDragStop,
  onNodeDoubleClick,
  removeNodes,
  removeEdges,
  onEdgeClick,
  onEdgeDoubleClick,
  onPaneClick,
  onMove,
  viewport,
  setCenter,
} = useVueFlow()

// 获取所有可见的节点和边
const visibleNodes = computed(() => getNodes.value)
const visibleEdges = computed(() => getEdges.value)
const searchOptions = computed(() => {
  return getNodes.value
    .filter(n => n.type === 'table-node' && (n.data as any)?.table?.name)
    .map(n => ({ id: n.id, name: (n.data as any).table.name as string }))
})

// 选择器状态（由父组件管理并通过 props 传递）
const selectorVisible = ref(false)
const selectorPos = ref<{ x: number; y: number }>({ x: 0, y: 0 })
const selectorEdgeId = ref<string>('')
const selectorSourceX = ref(0)
const selectorTargetX = ref(0)
const zoom = ref(1)
const searchVisible = ref(false)
const searchMessage = ref('')

// 监听表备注更新事件
const handleTableRemarkUpdated = (data: any) => {
  const nodes = getNodes.value;

  // 查找并更新匹配的节点
  const updatedNodes = nodes.map(node => {
    // 检查节点是否匹配更新的表
    const nodeTableId = (node.data as any)?.tableId ?? (node.data?.table as any)?.id;
    const matchById = data.tableId && node.data && nodeTableId === data.tableId;
    const matchByName = node.data.dbId === data.dbId && node.data.table.name === data.tableName;
    if (node.type === 'table-node' && (matchById || matchByName)) {

      // 深拷贝节点以避免引用问题
      const updatedNode = JSON.parse(JSON.stringify(node));

      // 更新表和字段的备注信息
      updatedNode.data.table.remark = data.tableInfo.remark;
      updatedNode.data.table.fields = data.tableInfo.fields;

      return updatedNode;
    }
    return node;
  });

  // 只有当有节点更新时才重新设置节点
  if (JSON.stringify(nodes) !== JSON.stringify(updatedNodes)) {
    setNodes(updatedNodes);
    // 只保存更新后的节点数据
    saveNodeData();
  }
};

// 节点拖动相关事件处理
onNodeDragStop((event) => {
  const movedNode = event.node
  const edges = getEdges.value
  const nodes = getNodes.value
  let edgesChanged = false;

  edges.forEach(edge => {
    if (edge.source === movedNode.id || edge.target === movedNode.id) {
      const sourceNode = nodes.find(n => n.id === edge.source)
      const targetNode = nodes.find(n => n.id === edge.target)

      if (sourceNode && targetNode) {
        // 获取字段名（移除前缀）
        const sourceField = edge.sourceHandle?.replace(/^(sl|sr|tl|tr)-/, '')
        const targetField = edge.targetHandle?.replace(/^(sl|sr|tl|tr)-/, '')

        if (sourceField && targetField) {
          let newSourceHandle = '';
          let newTargetHandle = '';
          
          if (sourceNode.position.x < targetNode.position.x) {
            newSourceHandle = `sr-${sourceField}`;
            newTargetHandle = `tl-${targetField}`;
          } else {
            newSourceHandle = `sl-${sourceField}`;
            newTargetHandle = `tr-${targetField}`;
          }
          
          // 只有当处理器变化时才标记为已更改
          if (edge.sourceHandle !== newSourceHandle || edge.targetHandle !== newTargetHandle) {
            edge.sourceHandle = newSourceHandle;
            edge.targetHandle = newTargetHandle;
            edgesChanged = true;
          }
        }
      }
    }
  })
  
  // 保存节点数据
  saveNodeData()
  
  // 只有当边缘发生变化时才保存边缘数据
  if (edgesChanged) {
    setEdges([...edges])
    saveEdgeData()
  }
})

// 节点交互事件处理 - 使用inject注入的方法直接调用
onNodeDoubleClick((event) => {
  if (event.node.type === 'table-node' && showTableDialog) {
    // 直接调用注入的方法
    showTableDialog(event.node.data)
  }
})

// 当前选中的边缘ID
let selectedEdgeId: string | null = null;

// 选择器触发的关系更新处理（顶层作用域，便于 on/off 引用一致）
const handleRelationUpdate = ({ edgeId, source, target }: { edgeId: string, source: string, target: string }) => {
  const edges = getEdges.value.map(e => {
    if (e.id === edgeId) {
      const newData = { ...(e.data || {}), relationship: { source, target } }
      return { ...e, data: newData }
    }
    return e
  })                                                                                                                 
  setEdges(edges)
  saveEdgeData()
}

const deleteSelectedEdge = () => {
  if (selectedEdgeId) {
    const edges = getEdges.value.filter(edge => edge.id !== selectedEdgeId);
    setEdges(edges);
    selectedEdgeId = null;
    saveEdgeData();
  }
};

// 清除选中并隐藏选择器
const clearSelectionAndHide = () => {
  selectedEdgeId = null
  selectorVisible.value = false
  // 清除所有边的选中状态
  const edges = getEdges.value
  const updated = edges.map(e => (e.selected ? { ...e, selected: false } : e))
  setEdges(updated)
}

// 选择器关闭时，清理选中高亮（edge.selected）
const handleSelectorClose = () => {
  selectorVisible.value = false
  selectedEdgeId = null
}

// 键盘事件处理函数
const handleKeyDown = (event: KeyboardEvent) => {
  // 只处理 Delete 和 Backspace 键
  if (event.key !== 'Delete' && event.key !== 'Backspace') {
    return;
  }

  // 检查当前是否有输入框处于焦点状态
  const activeElement = document.activeElement;
  const isInputFocused = activeElement && (activeElement.tagName === 'INPUT' || activeElement.tagName === 'TEXTAREA');
  if (isInputFocused) {
    return;
  }

  // 判断是否为节点 DOM（class 包含 vue-flow__node）
  // 这里假设节点 DOM 上有 vue-flow__node 且 data-id 存储节点 id
  const target = event.target as HTMLElement;
  if (target && target.classList.contains('vue-flow__node')) {
    const nodeId = target.getAttribute('data-id');
    if (nodeId) {
      // 先删除与该节点相关的所有边
      const edges = getEdges.value;
      const relatedEdgeIds = edges.filter(edge => edge.source === nodeId || edge.target === nodeId).map(edge => edge.id);
      if (relatedEdgeIds.length > 0) {
        removeEdges(relatedEdgeIds);
        saveEdgeData(); // 先保存边数据
      }
      removeNodes([nodeId]); // 删除节点
      saveNodeData(); // 保存节点数据
      return;
    }
  }

  // 有选中的边缘时，才执行删除
  deleteSelectedEdge();
};

const handleSearchKeyDown = (event: KeyboardEvent) => {
  const key = event.key.toLowerCase()
  if (event.ctrlKey && key === 'f') {
    event.preventDefault()
    searchMessage.value = ''
    searchVisible.value = true
  }
}

const performSearch = (query: string) => {
  const q = (query || '').trim().toLowerCase()
  if (!q) return
  const nodes = getNodes.value
  const matches = nodes.filter(n => n.type === 'table-node' && n.data && (n.data as any).table && typeof (n.data as any).table.name === 'string' && (n.data as any).table.name.toLowerCase().includes(q))
  if (!matches.length) {
    searchMessage.value = '未找到匹配的表'
    return
  }
  const target = matches[0]
  const cx = target.position.x + ((target.dimensions?.width as number) || 0) / 2
  const cy = target.position.y + ((target.dimensions?.height as number) || 0) / 2
  try {
    setCenter(cx, cy, { zoom: 1.2 })
  } catch {}
  const updated = nodes.map(n => ({ ...n, selected: n.id === target.id }))
  setNodes(updated)
  searchVisible.value = false
  searchMessage.value = ''
}

const handleSearchSelect = (id: string) => {
  const nodes = getNodes.value
  const target = nodes.find(n => n.id === id)
  if (!target) return
  const cx = target.position.x + ((target.dimensions?.width as number) || 0) / 2
  const cy = target.position.y + ((target.dimensions?.height as number) || 0) / 2
  try {
    setCenter(cx, cy, { zoom: 1.2 })
  } catch {}
  const updated = nodes.map(n => ({ ...n, selected: n.id === target.id }))
  setNodes(updated)
  searchVisible.value = false
  searchMessage.value = ''
}

// 统一的数据加载函数（异步，使用 AppCache）
const loadPageData = async (pageId: string) => {
  console.log("loadPageData :", pageId)
  if (pageId) {
    try {
      const flowData = await flowDataStore.getFlowDataAsync(pageId)
      const nodes = flowData?.nodes || []
      const edges = flowData?.edges || []
      setNodes(nodes)
      setEdges(edges)
    } catch (e) {
      console.warn('加载页面数据失败:', e)
      setNodes([])
      setEdges([])
    }
  }
}

// 监听路由参数变化
watch(() => route.params.pageId, async (newPageId) => {
  if (newPageId) {
    await loadPageData(newPageId as string)
  }
}, { immediate: false })

// 监听activeTab变化
watch(() => pageStore.activeTab, async (newActiveTab) => {
  if (newActiveTab && !route.params.pageId) {
    await loadPageData(newActiveTab)
  }
}, { immediate: false })

// 组件挂载时注册事件监听
onMounted(async () => {
  const pageId = (route.params.pageId as string) || pageStore.activeTab
  if (pageId) {
    await loadPageData(pageId)
  }

  // 注册表备注更新事件监听
  eventBus.on('table-remark-updated', handleTableRemarkUpdated);
  
  // 注册边缘数据更新事件监听
  eventBus.on('edge-data-updated', saveEdgeData);
  
  // 简化缩放逻辑：不持续监听缩放，仅在选择器显示时取一次 zoom
  
  // 监听边缘删除事件
  document.addEventListener('keydown', handleKeyDown);
  document.addEventListener('keydown', handleSearchKeyDown);

  // 点击空白区域清除选中与选择器
  onPaneClick(() => {
    clearSelectionAndHide()
  })

  // 缩放/移动时清除选中与选择器
  onMove(() => {
    clearSelectionAndHide()
  })
})

// 边缘双击事件处理
onEdgeDoubleClick((event) => {
  const mouseEvent = event.event as MouseEvent
  eventBus.emit('showEdgeInput', {
    edge: event.edge,
    position: { x: mouseEvent.clientX, y: mouseEvent.clientY }
  })
})

// 边缘点击事件处理
onEdgeClick(({ edge, event }) => {
  // 设置当前选中的边缘
  selectedEdgeId = edge.id;
  const mouseEvent = event as MouseEvent

  // 计算左右端点X坐标：基于节点位置（简化为比较X确定左右）
  const sourceNode = getNodes.value.find(n => n.id === edge.source)
  const targetNode = getNodes.value.find(n => n.id === edge.target)
  const sourceX = sourceNode?.position?.x || 0
  const targetX = targetNode?.position?.x || 0

  // 更新选择器状态并显示（通过 props 传递）
  // 仅在显示时取一次当前缩放
  zoom.value = (viewport as any)?.zoom ?? (viewport as any)?.value?.zoom ?? 1
  selectorEdgeId.value = edge.id
  selectorPos.value = { x: mouseEvent.clientX, y: mouseEvent.clientY }
  selectorSourceX.value = sourceX
  selectorTargetX.value = targetX
  selectorVisible.value = true

  // 切换选中效果：仅当前点击的边为选中，其他取消
  const edges = getEdges.value
  const updated = edges.map(e => {
    const shouldSelect = e.id === edge.id
    if (e.selected === shouldSelect) return e
    return { ...e, selected: shouldSelect }
  })
  setEdges(updated)
})

// 组件卸载前移除事件监听
onBeforeUnmount(() => { 
  eventBus.off('table-remark-updated', handleTableRemarkUpdated);
  eventBus.off('edge-data-updated', saveEdgeData);
  eventBus.off('relation-selector:update', handleRelationUpdate);


  // 卸载键盘事件监听
  document.removeEventListener('keydown', handleKeyDown);
  document.removeEventListener('keydown', handleSearchKeyDown);
})

// 支持手动连接表字段的线
onConnect((params) => {
  if(params.source === params.target) {
    console.warn('自连接不支持')
    return
  }

  addEdges([{
    ...params,
    type: 'button',
    // data: { text: 'styled custom edge label',
    //   relationship: {
    //     source: RelationshipType.ONE,
    //     target: RelationshipType.MANY
    //   }
    //  },
  }])
  // 只保存边缘数据
  saveEdgeData()
})

// 关系更新逻辑已迁移至 RelationshipSelector，通过事件总线触发
// 由子组件 emit('update') 回传，沿用原处理逻辑
</script>

<template>
  <!-- 包裹为单一根元素，并将非 props 的属性（如 style）显式绑定到容器 -->
  <div class="flow-root" v-bind="$attrs">
    <div v-if="pageStore.activeTab" class="flow-canvas" @dragover="onDragOver" @drop="onDrop">
      <VueFlow :nodes="visibleNodes" :edges="visibleEdges" 
        :default-viewport="{ x: 0, y: 0, zoom: 1 }" :fit-view-on-init="true">

        <!-- bind your custom node type to a component by using slots, slot names are always `node-<type>` -->
        <template #node-table-node="specialNodeProps">
          <TableNode v-bind="specialNodeProps" />
        </template>

        <template #edge-button="buttonEdgeProps">
          <EdgeWithButton :id="buttonEdgeProps.id" :source-x="buttonEdgeProps.sourceX" :source-y="buttonEdgeProps.sourceY"
            :target-x="buttonEdgeProps.targetX" :target-y="buttonEdgeProps.targetY"
            :source-position="buttonEdgeProps.sourcePosition" :target-position="buttonEdgeProps.targetPosition"
            :marker-end="buttonEdgeProps.markerEnd" :data="buttonEdgeProps.data" :style="buttonEdgeProps.style" :selected="buttonEdgeProps.selected" />
        </template>

        <!--  定制背景的样式 -->
        <Background :size="2" :gap="20" pattern-color="transparent" />
      </VueFlow>
    </div>
    <div v-else class="empty-state">
      <p>未选择活动页面</p>
    </div>

    <!-- 全局浮层渲染关系选择器，定位在点击位置上方 -->
    <RelationshipSelector
      :visible="selectorVisible"
      :position="selectorPos"
      :edge-id="selectorEdgeId"
      :source-x="selectorSourceX"
      :target-x="selectorTargetX"
      :zoom="zoom"
      @close="handleSelectorClose"
      @update="handleRelationUpdate"
    />
    <FlowSearch
      :visible="searchVisible"
      :message="searchMessage"
      :options="searchOptions"
      @close="() => { searchVisible = false; searchMessage = '' }"
      @search="performSearch"
      @select="handleSearchSelect"
    />
  </div>
</template>

<style scoped>
.flow-root {
  width: 100%;
  height: 100%;
}
.flow-canvas {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.empty-state {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
</style>
