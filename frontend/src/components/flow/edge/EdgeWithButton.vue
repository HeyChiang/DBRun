<script setup>
import { BaseEdge, getBezierPath, useEdge } from '@vue-flow/core'
import { computed, ref } from 'vue'
import { eventBus } from '@/utils/eventBus'
import CustomMarker from './CustomMarker.vue'
import { RelationshipType } from '@/types/relationshipTypes'

const props = defineProps({
  id: { type: String, required: true },
  sourceX: { type: Number, required: true },
  sourceY: { type: Number, required: true },
  targetX: { type: Number, required: true },
  targetY: { type: Number, required: true },
  sourcePosition: { type: String, required: true },
  targetPosition: { type: String, required: true },
  data: { type: Object, required: true },
  markerEnd: { type: String, required: false },
  style: { type: Object, required: false },
  selected: { type: Boolean, required: false, default: false },
})

// 使用useEdge获取edge对象
const { edge } = useEdge()

// 使用官方示例的居中定位方案
const path = computed(() => getBezierPath(props))

// 取消所有基于像素的偏移，完全以中点为基准进行居中定位

// 标记ID
const sourceMarkerId = computed(() => `${props.id}-source-marker`)
const targetMarkerId = computed(() => `${props.id}-target-marker`)

// 标记类型
const getMarkerType = (type) => {
  if (type === RelationshipType.MANY) return 'many'
  if (type === RelationshipType.ONE) return 'one'
  if (type === RelationshipType.NONE) return null
  return null
}

const sourceMarkerType = computed(() => getMarkerType(props.data?.relationship?.source))
const targetMarkerType = computed(() => getMarkerType(props.data?.relationship?.target))

// 标记颜色（根据选中状态决定）
const markerColor = computed(() => {
  if (props.selected) {
    const isDark = document.body.classList.contains('dark')
    return isDark ? '#36CFC9' : '#409EFF'
  }
  return (props.style && props.style.stroke) ? props.style.stroke : '#b1b1b7'
})

// 自定义标记URL
const customMarkerStart = computed(() => {
  return sourceMarkerType.value ? `url(#${sourceMarkerId.value})` : undefined
})

const customMarkerEnd = computed(() => {
  return targetMarkerType.value ? `url(#${targetMarkerId.value})` : props.markerEnd
})

// 文本编辑
const isEditing = ref(false)
const inputText = ref('')

const showInput = () => {
  inputText.value = edge.data.text || ''
  isEditing.value = true
  setTimeout(() => {
    document.getElementById(`edge-input-${props.id}`)?.focus()
  }, 10)
}

const handleInputConfirm = () => {
  if (inputText.value.trim() !== edge.data.text) {
    edge.data = {
      ...edge.data,
      text: inputText.value,
    }
    eventBus.emit('edge-data-updated')
  }
  isEditing.value = false
}

const handleBlur = () => {
  isEditing.value = false
}


// 旧的关系更新与全局事件监听逻辑已移除，改为由父组件统一管理
</script>

<script>
export default {
  inheritAttrs: false,
}
</script>

<template>
  <!-- 使用自定义标记的边 -->
  <BaseEdge
    :id="id"
    :style="{ ...style, stroke: markerColor, strokeWidth: style?.strokeWidth || 2 }"
    :path="path[0]"
    :marker-end="customMarkerEnd"
    :marker-start="customMarkerStart"
  />

  <!-- 自定义标记组件 -->
  <CustomMarker
    v-if="sourceMarkerType"
    :id="sourceMarkerId"
    :type="sourceMarkerType"
    :stroke="markerColor"
    :fill="markerColor"
    :stroke-width="2"
    :width="11"
    :height="11"
  />
  <CustomMarker
    v-if="targetMarkerType"
    :id="targetMarkerId"
    :type="targetMarkerType"
    :stroke="markerColor"
    :fill="markerColor"
    :stroke-width="2"
    :width="11"
    :height="11"
  />

  <!-- 居中渲染标签（仅文本编辑/显示），避免与关系选择器耦合 -->
  <foreignObject
    :x="path[1] - 60"
    :y="path[2] - 16"
    :width="120"
    :height="32"
    class="nodrag nopan"
    style="overflow: visible; pointer-events: all;"
  >
    <div
      xmlns="http://www.w3.org/1999/xhtml"
      style="width: 100%; height: 100%; display: flex; align-items: center; justify-content: center;"
    >
      <!-- 文本编辑/显示按钮 -->
      <div v-if="isEditing" class="p-inputgroup">
        <input
          type="text"
          :id="`edge-input-${id}`"
          class="p-inputtext p-component edge-input"
          v-model="inputText"
          @keyup.enter="handleInputConfirm"
          @blur="handleBlur"
          autocomplete="off"
        />
      </div>
      <button
        v-else-if="data.text"
        class="edgebutton"
        @click="showInput"
        :title="data.text.length > 10 ? data.text : ''"
      >
        {{ data.text }}
      </button>
    </div>
  </foreignObject>
</template>

<style scoped>
.edge-input { 
  min-width: 60px;
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid var(--primary-color, #3B82F6);
  font-size: 12px;
  background-color: var(--card-bg, #ffffff);
  color: var(--body-text-color, #333333);
  caret-color: var(--primary-color, #3B82F6);
  box-shadow: 0 0 2px rgba(0, 0, 0, 0.1);
}

.edgebutton {
  cursor: pointer;
  background-color: var(--card-bg);
  color: var(--body-text-color);
  padding: 2px 5px;
  font-size: 12px;
  text-align: center;
  border-radius: 3px;
  border: 1px solid var(--border-color);
  min-width: 40px;
  max-width: 150px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  user-select: none;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.edgebutton:hover {
  background-color: var(--surface-hover);
  box-shadow: 0 0 3px rgba(0, 0, 0, 0.2);
}

.p-inputgroup {
  background-color: white;
  border: 1px solid #ccc;
  border-radius: 3px;
  min-width: 100px;
}

.edge-input {
  width: 100%;
  padding: 2px 5px;
  font-size: 12px;
  outline: none;
}

</style>
