<script setup lang="ts">
import { useRouter } from 'vue-router'
import { usePageStore } from '@/stores/pageStore'
import { ref, onMounted, onUnmounted } from 'vue'
import { eventBus } from '@/utils/eventBus'

const pageStore = usePageStore()
const router = useRouter()

// Context menu state
const contextMenu = ref({ visible: false, x: 0, y: 0, tabKey: '' })

// 打开或切换到指定标签页
const openTab = (pageKey: string) => {
  pageStore.openTab(pageKey)
  // 触发页面切换事件
  eventBus.emit('page-switch', pageKey)
  // 跳转到子页面路由
  router.push({ name: 'flowPage', params: { pageId: pageKey }})
}

// 关闭标签页
const closeTab = (pageKey: string) => {
  pageStore.closeTab(pageKey)
  
  // 根据当前活动标签页跳转路由
  if (pageStore.activeTab) {
    // 触发页面切换事件
    eventBus.emit('page-switch', pageStore.activeTab)
    router.push({ name: 'flowPage', params: { pageId: pageStore.activeTab }})
  } else {
    router.push({ name: 'app' })
  }
}

// 显示上下文菜单
const showContextMenu = (event: MouseEvent, tabKey: string) => {
  event.preventDefault()
  event.stopPropagation()
  contextMenu.value = {
    visible: true,
    x: event.clientX,
    y: event.clientY,
    tabKey
  }
}

// 隐藏上下文菜单
const hideContextMenu = () => {
  contextMenu.value.visible = false
}

// 全局点击事件处理
const handleGlobalClick = (event: MouseEvent) => {
  if (contextMenu.value.visible) {
    hideContextMenu()
  }
}

// 挂载和卸载全局事件监听器
onMounted(() => {
  document.addEventListener('click', handleGlobalClick)
})

onUnmounted(() => {
  document.removeEventListener('click', handleGlobalClick)
})

// 获取tab的索引
const getTabIndex = (tabKey: string) => {
  return pageStore.openTabs.findIndex(tab => tab.key === tabKey)
}

// 关闭左侧标签页
const closeTabsToLeft = (tabKey: string) => {
  const currentIndex = getTabIndex(tabKey)
  const tabsToClose = pageStore.openTabs.slice(0, currentIndex)
  tabsToClose.forEach(tab => closeTab(tab.key))
  hideContextMenu()
}

// 关闭右侧标签页
const closeTabsToRight = (tabKey: string) => {
  const currentIndex = getTabIndex(tabKey)
  const tabsToClose = pageStore.openTabs.slice(currentIndex + 1)
  tabsToClose.forEach(tab => closeTab(tab.key))
  hideContextMenu()
}

// 关闭所有标签页
const closeAllTabs = () => {
  // Create a copy of the array，避免遍历的时候导致关闭错误
  const tabsToClose = [...pageStore.openTabs]  
  tabsToClose.forEach(tab => closeTab(tab.key))
  hideContextMenu()
}

// 检查是否有左侧标签页
const hasLeftTabs = (tabKey: string) => {
  return getTabIndex(tabKey) > 0
}

// 检查是否有右侧标签页
const hasRightTabs = (tabKey: string) => {
  return getTabIndex(tabKey) < pageStore.openTabs.length - 1
}
</script>

<template>
  <div class="tabs-container flow-tabs">
    <div v-for="tab in pageStore.openTabs" 
         :key="tab.key" 
         class="tab" 
         :class="{ active: tab.active }"
         @click="openTab(tab.key)"
         @contextmenu="showContextMenu($event, tab.key)">
      <span>{{ tab.label }}</span>
      <button v-if="tab.active" class="close-btn" @click.stop="closeTab(tab.key)">
        <i class="pi pi-times"></i>
      </button>
    </div>

    <!-- Context Menu -->
    <div v-if="contextMenu.visible" 
         class="context-menu" 
         :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
         @click.stop>
      <div class="menu-item" @click.stop="closeTab(contextMenu.tabKey); hideContextMenu()">
        Close
      </div>
      <div v-if="hasLeftTabs(contextMenu.tabKey)" 
           class="menu-item" 
           @click.stop="closeTabsToLeft(contextMenu.tabKey); hideContextMenu()">
        Close to the Left
      </div>
      <div v-if="hasRightTabs(contextMenu.tabKey)" 
           class="menu-item" 
           @click.stop="closeTabsToRight(contextMenu.tabKey); hideContextMenu()">
        Close to the Right
      </div>
      <div class="menu-item" @click.stop="closeAllTabs(); hideContextMenu()">
        Close All
      </div>
    </div>
  </div>
</template>

<style scoped>
.tabs-container {
  display: flex;
  background-color: var(--surface-section);
  border-bottom: 0.1rem solid var(--surface-border);
}

.flow-tabs {
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
}

.tab {
  display: flex;
  align-items: center;
  padding: 0.5rem 1rem;
  background-color: var(--surface-ground);
  border: 1px solid var(--surface-border);
  border-left: none;
  border-top: none;
  border-bottom: none;
  cursor: pointer;
  gap: 0.5rem;
  position: relative;
  top: 1px;
  transition: all 0.2s ease;
  color: var(--text-color-secondary);
  font-weight: normal;
}

.tab:hover:not(.active) {
  background-color: var(--surface-hover);
}

.tab.active {
  position: relative;
  background-color: var(--card-bg);
  border-bottom: none;
  color: var(--text-color-primary);
  font-weight: bold;
}

.tab.active::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background-color: #2196F3;
}

.close-btn {
  background: none;
  border: none;
  padding: 0.25rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-color-secondary);
  border-radius: 50%;
  width: 1.5rem;
  height: 1.5rem;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background-color: var(--surface-hover);
  color: var(--text-color);
}

.context-menu {
  position: fixed;
  background-color: var(--body-bg);
  border: 1px solid var(--surface-border);
  border-radius: 6px;
  padding: 4px 0;
  min-width: 160px;
  box-shadow: 0 2px 4px -1px rgba(0,0,0,.2), 0 4px 5px 0 rgba(0,0,0,.14), 0 1px 10px 0 rgba(0,0,0,.12);
  z-index: 1000;
  transform-origin: top left;
  animation: menuAppear 0.15s ease-out;
}

.menu-item {
  padding: 8px 16px;
  cursor: pointer;
  color: var(--text-color);
  transition: all 0.2s ease;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
  white-space: nowrap;
  margin: 2px 4px;
  border-radius: 4px;
}

.menu-item:hover {
  background-color: var(--surface-hover);
  transform: translateX(2px);
}

@keyframes menuAppear {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
</style>
