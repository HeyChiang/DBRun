<script setup lang="ts">
import { ref, watch, nextTick, computed } from 'vue'

const props = defineProps<{ visible: boolean; message?: string; options?: { id: string; name: string }[] }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'search', query: string): void; (e: 'select', id: string): void }>()

const query = ref('')
const inputRef = ref<HTMLInputElement | null>(null)
const selectedIndex = ref(0)

watch(() => props.visible, async (v) => {
  if (v) {
    await nextTick()
    inputRef.value?.focus()
    inputRef.value?.select()
    selectedIndex.value = 0
  }
})

watch(query, () => {
  selectedIndex.value = 0
})

const filteredOptions = computed(() => {
  const term = query.value.trim().toLowerCase()
  if (!term) return props.options || []
  return (props.options || []).filter(o => o.name?.toLowerCase().includes(term))
})

const onKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') emit('close')
  if (e.key === 'Enter') {
    const opt = filteredOptions.value[selectedIndex.value]
    if (opt) emit('select', opt.id)
    else emit('search', query.value)
  }
  if (e.key === 'ArrowDown') {
    if (filteredOptions.value.length) {
      selectedIndex.value = Math.min(selectedIndex.value + 1, filteredOptions.value.length - 1)
    }
    e.preventDefault()
  }
  if (e.key === 'ArrowUp') {
    if (filteredOptions.value.length) {
      selectedIndex.value = Math.max(selectedIndex.value - 1, 0)
    }
    e.preventDefault()
  }
}

const onBackdropClick = (e: MouseEvent) => {
  const target = e.target as HTMLElement
  if (target.classList.contains('flow-search-backdrop')) emit('close')
}
</script>

<template>
  <div v-if="visible" class="flow-search-backdrop" @click="onBackdropClick">
    <div class="flow-search-panel" @keydown.stop="onKeydown">
      <input ref="inputRef" v-model="query" class="flow-search-input" type="text" placeholder="搜索表名，回车确认" />
      <div v-if="filteredOptions.length" class="flow-search-suggestions">
        <div
          v-for="(opt, idx) in filteredOptions"
          :key="opt.id"
          class="flow-search-suggestion"
          :class="{ active: idx === selectedIndex }"
          @mousedown.prevent="emit('select', opt.id)"
        >
          {{ opt.name }}
        </div>
      </div>
      <div v-if="message" class="flow-search-message">{{ message }}</div>
    </div>
  </div>
  </template>

<style scoped>
.flow-search-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.flow-search-panel {
  width: 420px;
  max-width: 90vw;
  background: var(--surface-card, #fff);
  border-radius: 8px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.flow-search-input {
  width: 100%;
  height: 36px;
  border: 1px solid var(--surface-border, #ddd);
  border-radius: 6px;
  padding: 0 10px;
  outline: none;
  background: var(--surface-ground, #fff);
  color: var(--text-color, #222);
}
.flow-search-input:focus {
  border-color: var(--primary-color, #4f46e5);
}
.flow-search-suggestions {
  max-height: 220px;
  overflow-y: auto;
  border: 1px solid var(--surface-border, #ddd);
  border-radius: 6px;
  background: var(--card-bg);
}
.flow-search-suggestion {
  padding: 8px 10px;
  cursor: pointer;
  color: var(--text-color);
}
.flow-search-suggestion:hover,
.flow-search-suggestion.active {
  background: var(--surface-hover);
}
.flow-search-message {
  color: var(--red-500, #ef4444);
  font-size: 12px;
}
</style>