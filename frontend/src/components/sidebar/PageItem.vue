<template>
  <li class="mb-1">
    <!-- Page Item -->
    <template v-if="item.type === 'page'">
      <router-link 
        :to="{ name: 'flowPage', params: { pageId: item.key }}"
        custom
        v-slot="{ navigate, isActive }"
      >
        <div 
          class="flex align-items-center cursor-pointer p-2 border-round text-color hover:surface-200"
          :class="{ 'surface-200': isActive }"
          @click="navigate(); handlers?.handlePageClick(item)"
          @contextmenu.prevent.stop="handlers?.handleContextMenu($event, item)"
          draggable="true"
          @dragstart="handlers?.handleDragStart(item)"
          @dragover.prevent="handlers?.handleDragOver($event, item)"
          @drop.prevent="handlers?.handleDrop($event, item)"
        >
          <span class="mr-2 ml-2"><i :class="item.icon || 'pi pi-file'"></i></span>
          <span class="item_label">{{ item.label }}</span>
        </div>
      </router-link>
    </template>

    <!-- Group Item -->
    <template v-else>
      <div 
        class="flex align-items-center cursor-pointer p-2 border-round text-color hover:surface-200"
        @click="handlers?.toggleGroup(item.key)"
        @contextmenu.prevent.stop="handlers?.handleContextMenu($event, item)"
        draggable="true"
        @dragstart="handlers?.handleDragStart(item)"
        @dragover.prevent="handlers?.handleDragOver($event, item)"
        @drop.prevent="handlers?.handleDrop($event, item)"
      >
        <span class="mr-2"><i style="font-size: 0.8rem;" :class="[item.expanded ? 'pi pi-chevron-down' : 'pi pi-chevron-right']"></i></span>
        <span class="mr-2"><i :class="item.icon || 'pi pi-folder'"></i></span>
        <span class="flex-grow-1 item_label">{{ item.label }}</span>
      </div>
      
      <div v-show="item.expanded">
        <ul class="list-none pl-3 py-0 m-0">
          <PageItem
            v-for="child in item.children"
            :key="child.key"
            :item="child"
          />
        </ul>
      </div>
    </template>
  </li>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import type { Page } from '@/stores/pageStore'
import { PAGE_HANDLERS_KEY } from './constants'

defineProps<{
  item: Page
}>()

interface PageHandlers {
  handlePageClick: (page: Page) => void
  handleContextMenu: (event: MouseEvent, item: Page) => void
  toggleGroup: (key: string) => void
  handleDragStart: (item: Page) => void
  handleDragOver: (event: DragEvent, item: Page) => void
  handleDrop: (event: DragEvent, item: Page) => void
}

const handlers = inject<PageHandlers>(PAGE_HANDLERS_KEY)

if (!handlers) {
  console.warn('PageHandlers not provided to PageItem component')
}
</script>

<style>
.item_label {
  color: var(--text-color);
}

.pi {
  color: var(--text-color);
}

.text-color {
  color: var(--text-color) !important;
}

.hover\:surface-200:hover {
  background-color: var(--surface-hover-light);
  border-radius: 0.3rem;
}

.surface-200 {
  background-color: var(--surface-hover);
  border-radius: 0.3rem;
}
</style>
