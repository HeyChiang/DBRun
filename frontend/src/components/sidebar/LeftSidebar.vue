<script setup>
import { onMounted, ref, toRaw } from 'vue'
import DBToolbar from '@/components/sidebar/DBToolbar.vue'
import { useDatabaseStore } from '@/stores/databaseStore'
// 刷新逻辑已移至 SidebarTree.vue
import PageMenuManager from './PageMenuManager.vue'
import { usePageStore } from '@/stores/pageStore'
import SidebarTree from './SidebarTree.vue'

const pageStore = usePageStore();

const databaseStore = useDatabaseStore()
// 连接编辑与删除逻辑已迁移到 SidebarTree.vue

const searchTerm = ref('')

const onSearchTermUpdate = (val) => {
  searchTerm.value = val || ''
}

onMounted(() => {
  console.log('Sidebar mounted')
  databaseStore.refreshDatabases()
})

// 刷新（数据库或表）逻辑由 SidebarTree 组件直接处理

</script>

<template> 
  <div class="sidebar">
    <PageMenuManager />
    <DBToolbar :search-term="searchTerm" @update:searchTerm="onSearchTermUpdate" />
    <SidebarTree
      :db-links="databaseStore.dbLinks"
      :can-drag="!!pageStore.activeTab"
      :search-term="searchTerm"
    />
  </div>
</template>

<style>
.sidebar {
  display: flex;
  flex-direction: column;
  background-color: var(--surface-section);
  height: 100%;
  color: var(--text-color);
  width: 15rem;
  min-width: 15rem;
}

.sidebar-content {
  flex: 1;
  overflow-y: auto;
}

.menu-item {
  padding: 0.5rem;
  cursor: pointer;
  border-radius: 0.2rem;
  transition: background-color 0.2s;
}

.menu-item:hover {
  background-color: var(--surface-hover);
}

.context-menu {
  position: fixed;
  background: var(--overlay-background, #fff);
  border: 1px solid var(--border-color, #e0e0e0);
  border-radius: 0.2rem;
  padding: 0.5rem 0;
  min-width: 150px;
  z-index: 1000;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.08);
  color: var(--body-text-color, #222);
}


.context-menu-item {
  padding: 0.5rem 1rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  color: var(--body-text-color, #222);
  background: transparent;
}


.context-menu-item:hover {
  background-color: var(--surface-hover);
  color: var(--body-text-color);
}


.dark .context-menu-item:hover {
  background-color: var(--p-surface-700, #222);
  color: #fff;
}
</style>