<script setup lang="ts">
import { computed } from 'vue'
import { RelationshipType } from '@/types/relationshipTypes'

const props = defineProps<{ 
  visible: boolean,
  position: { x: number, y: number },
  edgeId: string,
  sourceX: number,
  targetX: number,
  zoom?: number,
}>()

const emit = defineEmits<{ 
  (e: 'close'): void,
  (e: 'update', payload: { edgeId: string, source: string, target: string }): void,
}>()

const isVisible = computed(() => props.visible)
const posX = computed(() => props.position?.x ?? 0)
const posY = computed(() => props.position?.y ?? 0)
const zoom = computed(() => props.zoom ?? 1)

// 不使用全局点击监听，外部关闭交由画布事件处理

// 更新关系类型
const updateRelationship = (source: string, target: string) => {
  const isSourceOnLeft = props.sourceX <= props.targetX
  let leftNodeType = source
  let rightNodeType = target

  if (source === RelationshipType.ONE && target === RelationshipType.MANY) {
    leftNodeType = RelationshipType.ONE
    rightNodeType = RelationshipType.MANY
  } else if (source === RelationshipType.MANY && target === RelationshipType.ONE) {
    leftNodeType = RelationshipType.MANY
    rightNodeType = RelationshipType.ONE
  }

  const finalSource = isSourceOnLeft ? leftNodeType : rightNodeType
  const finalTarget = isSourceOnLeft ? rightNodeType : leftNodeType

  emit('update', {
    edgeId: props.edgeId,
    source: finalSource,
    target: finalTarget,
  })
  emit('close')
}

// 无需 onMounted/onBeforeUnmount 的全局事件管理
</script>

<template>
  <Teleport to="body">
    <div v-show="isVisible" class="relationship-selector-overlay"
         :style="{ position: 'fixed', left: posX + 'px', top: posY + 'px', zIndex: 9999, transform: `translate(-50%, -100%) scale(${zoom})`, transformOrigin: 'bottom center', pointerEvents: 'none' }">
      <!-- 关系类型选择器 -->
      <div class="relationship-selector" @click.stop @mousedown.stop style="pointer-events: auto;">
        <div class="relationship-options">
          <button class="relationship-option" @click="updateRelationship(RelationshipType.ONE, RelationshipType.ONE)" title="一对一关系 (1:1)">
            <svg class="icon" viewBox="0 0 800 800" width="24" height="24">
              <path d="M170 260v280M110 400h120" stroke="#8A8A8A" stroke-width="75" stroke-linecap="round" />
              <path d="M630 260v280M570 400h120" stroke="#8A8A8A" stroke-width="75" stroke-linecap="round" />
              <path d="M230 400h340" stroke="#8A8A8A" stroke-width="75" stroke-linecap="round" />
            </svg>
          </button>
          <button class="relationship-option" @click="updateRelationship(RelationshipType.ONE, RelationshipType.MANY)" title="一对多关系 (1:N)">
            <svg class="icon" viewBox="0 0 1024 1024" width="20" height="20" xmlns="http://www.w3.org/2000/svg">
              <path d="M873.344 297.344a32 32 0 0 1 45.312 45.312L781.248 480H896a32 32 0 0 1 0 64h-114.752l137.408 137.344a32 32 0 0 1-45.312 45.312L690.752 544H256V704a32 32 0 0 1-64 0V544H128a32 32 0 0 1 0-64h64V320a32 32 0 0 1 64 0v160h434.752l182.592-182.656z" fill="#8A8A8A" stroke="#8A8A8A" stroke-width="40" />
            </svg>
          </button>
          <button class="relationship-option" @click="updateRelationship(RelationshipType.MANY, RelationshipType.ONE)" title="多对一关系 (N:1)">
            <svg class="icon" viewBox="0 0 1024 1024" width="20" height="20" xmlns="http://www.w3.org/2000/svg">
              <path d="M150.656 297.344a32 32 0 0 0-45.312 45.312L242.752 480H128a32 32 0 0 0 0 64h114.752L105.344 681.344a32 32 0 0 0 45.312 45.312L333.248 544H768V704a32 32 0 0 0 64 0V544h64a32 32 0 0 0 0-64h-64V320a32 32 0 0 0-64 0v160H333.248l-182.592-182.656z" fill="#8A8A8A" stroke="#8A8A8A" stroke-width="40" />
            </svg>
          </button>
          <button class="relationship-option" @click="updateRelationship(RelationshipType.MANY, RelationshipType.MANY)" title="多对多关系 (N:N)">
            <svg class="icon" viewBox="0 0 1024 1024" width="20" height="20" xmlns="http://www.w3.org/2000/svg">
              <path d="M150.656 297.344a32 32 0 0 0-45.312 45.312L242.752 480H128a32 32 0 0 0 0 64h114.752L105.344 681.344a32 32 0 0 0 45.312 45.312L333.248 544h357.504l182.592 182.656a32 32 0 0 0 45.312-45.312L781.248 544H896a32 32 0 0 0 0-64h-114.752l137.408-137.344a32 32 0 0 0-45.312-45.312L690.752 480H333.248l-182.592-182.656z" fill="#8A8A8A" stroke="#8A8A8A" stroke-width="40" />
            </svg>
          </button>
          <button class="relationship-option" @click="updateRelationship(RelationshipType.NONE, RelationshipType.NONE)" title="直线">
            <svg class="icon" viewBox="0 0 1024 1024" width="20" height="20" xmlns="http://www.w3.org/2000/svg">
              <path d="M128 480a32 32 0 0 0 0 64h768a32 32 0 0 0 0-64H128z" fill="#8A8A8A" stroke="#8A8A8A" stroke-width="40" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
/* 关系选择器样式 */
.relationship-selector {
  position: relative; /* 由父容器负责绝对定位与居中 */
  background-color: var(--card-bg, #ffffff);
  border: 1px solid var(--border-color, #ddd);
  border-radius: 4px;
  padding: 1px;
  box-shadow: 0 2px 8px var(--box-shadow-color);
  z-index: 1000;
}

.relationship-selector-overlay {
  pointer-events: none;
}

.relationship-options {
  display: flex;
  flex-direction: row;
  gap: 3px;
}

.relationship-option {
  background: none;
  border: 1px solid transparent;
  cursor: pointer;
  border-radius: 3px;
  padding: 3px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.relationship-option:hover {
  background-color: var(--hover-bg, #f0f0f0);
  border-color: var(--border-color, #ddd);
}

.relationship-option svg {
  width: 20px;
  height: 20px;
}
</style>
