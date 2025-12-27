<template>
  <div class="page-manager" @contextmenu.prevent="handleEmptyAreaContextMenu($event)">
    <Toolbar class="icon-toolbar">
      <template #start>
        <div class="relative">
          <Button icon="pi pi-plus" text severity="secondary" @click="toggle" aria-haspopup="true" aria-controls="page_add_menu" />
          <Menu ref="menu" id="page_add_menu" :model="menuItems" :popup="true" />
        </div>
        <Button v-if="!showSearch" icon="pi pi-search" text severity="secondary" class="ml-2" @click="toggleSearch" />
        <template v-if="showSearch">
          <InputText v-model="localSearch" placeholder="搜索页面" class="ml-2 search-input" @input="onSearchInput" />
          <Button icon="pi pi-times" text severity="secondary" class="ml-1" @click="handleClearClick" />
        </template>
      </template>
    </Toolbar>
    <div class="page-manager-content no-select">
      <ul class="list-none px-1 m-0">
        <PageItem
          v-for="item in filteredPages"
          :key="item.key"
          :item="item"
          @page-click="handlePageClick"
          @context-menu="handleContextMenu"
          @toggle-group="toggleGroup"
        />
      </ul>
    </div>
    <ContextMenu ref="contextMenu" :model="contextMenuItems" />
    
    <Dialog v-model:visible="dialogVisible" :header="getDialogHeader" modal>
      <div class="flex flex-column gap-2">
        <div class="field">
          <label for="label">{{ dialogMode === 'edit' ? '名称' : '名称' }}</label>
          <InputText id="label" v-model="editingItem.label" maxlength="20" class="w-full" autocomplete="off" />
        </div>
        <div class="field" v-if="dialogMode !== 'edit'">
          <label for="group">所属分组（可选）</label>
          <Select id="group" v-model="editingItem.parentKey" :options="groupOptions" optionLabel="label" 
                 optionValue="key" class="w-full" placeholder="默认分组" :showClear="true" />
        </div>
      </div>
      <template #footer>
        <Button label="取消" icon="pi pi-times" @click="dialogVisible = false" text />
        <Button label="确定" icon="pi pi-check" @click="handleDialogConfirm" autofocus />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, provide, nextTick } from 'vue'
import { usePageStore } from '@/stores/pageStore'
import type { Page } from '@/stores/pageStore'
import PageItem from './PageItem.vue'
import { PAGE_HANDLERS_KEY } from './constants'
import { useRouter } from 'vue-router'
import Toolbar from 'primevue/toolbar'
import Button from 'primevue/button'
import Menu from 'primevue/menu'
import InputText from 'primevue/inputtext'

interface MenuItem {
  label: string;
  icon: string;
  command: () => void;
  disabled?: boolean;
}

const router = useRouter()
const pageStore = usePageStore()
const contextMenu = ref()
const dialogVisible = ref(false)
const dialogMode = ref<'add' | 'edit'>('add')
const rightClickedItem = ref<Page | null>(null)
// 正在拖拽的条目 key
const draggingKey = ref<string | null>(null)
const editingItem = ref<{
  label: string;
  parentKey?: string | null;
  type: 'page' | 'group';
  key?: string;
}>({
  label: '',
  parentKey: null,
  type: 'page'
})

const menu = ref()
const showSearch = ref(false)
const localSearch = ref('')

const getDialogHeader = computed(() => {
  switch (dialogMode.value) {
    case 'add': return editingItem.value.type === 'page' ? '新增页面' : '新增分组'
    case 'edit': return '编辑'
    default: return ''
  }
})

const groupOptions = computed(() => {
  const collectGroups = (items: Page[]): { label: string, key: string }[] => {
    
    let result: { label: string, key: string }[] = []

    items.forEach(item => {
      if (item.type === 'group') {
        result.push({
          label: item.label,
          key: item.key
        })
        if (item.children) {
          result = result.concat(collectGroups(item.children))
        }
      }
    })
    return result
  }
  return collectGroups(pageStore.pages)
})

const contextMenuItems = computed(() => {
  const baseItems: MenuItem[] = [
    {
      label: '新增页面',
      icon: 'pi pi-file',
      command: () => {
        showAddDialog('page', rightClickedItem.value?.type === 'group' ? rightClickedItem.value.key : null)
      }
    },
    {
      label: '新增分组',
      icon: 'pi pi-folder',
      command: () => {
        showAddDialog('group', rightClickedItem.value?.type === 'group' ? rightClickedItem.value.key : null)
      }
    },
    ...(rightClickedItem.value ? [
      {
        label: '编辑',
        icon: 'pi pi-pencil',
        command: () => showEditDialog(rightClickedItem.value as Page)
      },
      {
        label: '删除',
        icon: 'pi pi-trash',
        command: () => handleDelete()
      }
    ] : [])
  ]

  return baseItems
})

const menuItems = ref([
  {
    label: '新增页面',
    icon: 'pi pi-file',
    command: () => {
      showAddDialog('page', null)
    }
  },
  {
    label: '新增分组',
    icon: 'pi pi-folder',
    command: () => {
      showAddDialog('group', null)
    }
  }
])

function toggle(event: MouseEvent) {
  menu.value?.toggle(event)
}

function toggleSearch() {
  showSearch.value = true
  nextTick(() => {
    const el = document.querySelector('.icon-toolbar .search-input') as HTMLInputElement | null
    const fn = el ? Reflect.get(el, 'focus') : null
    if (typeof fn === 'function') fn.call(el)
  })
}

function onSearchInput() {}

function showContextMenu(event: MouseEvent) {
  contextMenu.value?.show(event)
}

function handleEmptyAreaContextMenu(event: MouseEvent) {
  rightClickedItem.value = null
  showContextMenu(event)
}

function showAddDialog(type: 'page' | 'group', parentKey: string | null = null) {
  dialogMode.value = 'add'
  editingItem.value = {
    label: '',
    parentKey,
    type
  }
  dialogVisible.value = true
}

function showEditDialog(item: Page) {
  dialogMode.value = 'edit'
  editingItem.value = {
    ...item,
    label: item.label,
    type: item.type
  }
  dialogVisible.value = true
}

function handleDialogConfirm() {
  if (!editingItem.value.label.trim()) {
    return
  }

  switch (dialogMode.value) {
    case 'add':
      const newItem = pageStore.addItem({
        label: editingItem.value.label.trim(),
        type: editingItem.value.type,
        parentKey: editingItem.value.parentKey ?? undefined
      })
      // 如果是新增页面，则自动打开该页面
      if (editingItem.value.type === 'page') {
        pageStore.openTab(newItem.key)
        router.push({ name: 'flowPage', params: { pageId: newItem.key }})
      }
      break
    case 'edit':
      if (rightClickedItem.value?.key) {
        pageStore.updatePage(rightClickedItem.value.key, {
          label: editingItem.value.label.trim()
        })
      }
      break
  }

  dialogVisible.value = false
}

function handleDelete() {
  if (rightClickedItem.value?.key) {
    pageStore.deletePage(rightClickedItem.value.key)
  }
  rightClickedItem.value = null
}

function handlePageClick(page: Page) {
  pageStore.openTab(page.key)
  // 跳转到子页面路由
  router.push({ name: 'flowPage', params: { pageId: page.key }})
}

function toggleGroup(key: string) {
  pageStore.toggleGroup(key)
}

function handleContextMenu(event: MouseEvent, item: Page) {
  rightClickedItem.value = item
  if (item.type === 'page') {
    handlePageContextMenu(event, item)
  } else {
    handleGroupContextMenu(event, item)
  }
}

function handlePageContextMenu(event: MouseEvent, page: Page) {
  rightClickedItem.value = page
  showContextMenu(event)
}

function handleGroupContextMenu(event: MouseEvent, group: Page) {
  rightClickedItem.value = group
  showContextMenu(event)
}

// 拖拽处理
function handleDragStart(item: Page) {
  draggingKey.value = item.key
}

function handleDragOver(event: DragEvent, target: Page) {
  // 允许放置
  event.preventDefault()
}

function handleDrop(event: DragEvent, target: Page) {
  event.preventDefault()
  const sourceKey = draggingKey.value
  if (!sourceKey || sourceKey === target.key) {
    draggingKey.value = null
    return
  }
  const el = event.currentTarget as HTMLElement
  const rect = el.getBoundingClientRect()
  const y = event.clientY
  const top = rect.top
  const h = rect.height

  // 三段式：上1/3 -> before；下1/3 -> after；中间 -> inside（仅分组）
  const thresholdTop = top + h / 3
  const thresholdBottom = top + (2 * h) / 3
  let position: 'before' | 'after' | 'inside'
  if (y < thresholdTop) {
    position = 'before'
  } else if (y > thresholdBottom) {
    position = 'after'
  } else {
    position = target.type === 'group' ? 'inside' : 'after'
  }

  pageStore.moveItem(sourceKey, target.key, position)
  draggingKey.value = null
}

provide(PAGE_HANDLERS_KEY, {
  handlePageClick,
  handleContextMenu,
  toggleGroup,
  handleDragStart,
  handleDragOver,
  handleDrop
})

const filteredPages = computed(() => {
  const term = localSearch.value.trim().toLowerCase()
  if (!term) return pageStore.pages

  const filterItems = (items: Page[]): Page[] => {
    const result: Page[] = []
    for (const item of items) {
      const match = item.label.toLowerCase().includes(term)
      if (item.type === 'group') {
        const children = item.children ? filterItems(item.children) : []
        if (match || children.length) {
          result.push({ ...item, children })
        }
      } else {
        if (match) result.push(item)
      }
    }
    return result
  }

  return filterItems(pageStore.pages)
})

function handleClearClick() {
  localSearch.value = ''
  showSearch.value = false
}
</script>

<style scoped>
.page-manager {
  height: 50%;
  overflow-y: auto;
}

.item_label {
  color: var(--text-color);
}

.page-manager-content {
  padding: 0.5rem;
}

:deep(.p-contextmenu) {
  min-width: 150px;
}

.no-select {
  user-select: none;
}

.icon-toolbar {
  border: none;
  padding: 0.2rem;
  background: transparent;
  border-top: 0.1rem solid var(--surface-border);
}

.icon-toolbar ::v-deep(.p-toolbar) {
  background: transparent;
  border: none;
  padding: 0;
}

.icon-toolbar ::v-deep(.p-button) {
  width: 2rem;
  height: 2rem;
}

.icon-toolbar ::v-deep(.p-button:focus) {
  box-shadow: none;
}

.icon-toolbar ::v-deep(.p-button:hover) {
  background-color: var(--surface-hover-light) !important;
  color: var(--text-color) !important;
}

.search-input {
  width: 10rem;
  height: 2rem;
  line-height: 2rem;
  padding: 0 0.5rem;
  font-size: 0.9rem;
  border-radius: 0.25rem;
  background: var(--surface-section);
  border: 1px solid var(--surface-border);
  color: var(--text-color);
}


</style>
