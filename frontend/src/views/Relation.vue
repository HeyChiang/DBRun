<script setup lang="ts">
import {ref, provide} from 'vue'
import LeftSidebar from '@/components/sidebar/LeftSidebar.vue'
import TopMenu from '@/components/menu/TopMenu.vue'
import TableViewDialog from '@/components/common/TableViewDialog.vue'
import FlowTabs from '@/components/flow/FlowTabs.vue'
import { usePageStore, type Tab } from '@/stores/pageStore'
import { computed } from 'vue'
import { TableNodeInfo } from '@/types/tableTypes'
import { models } from '@/../wailsjs/go/models'

// TableViewDialog相关状态
const dialogVisible = ref(false)
const currentTableData = ref<TableNodeInfo>({
  dbId: '',
  dbName: '',
  schemaName: '',
  table: {} as models.TableInfoVO
})

const showTableDialog = (tableData: TableNodeInfo) => {
  currentTableData.value = tableData
  dialogVisible.value = true
}

// 提供showTableDialog方法给后代组件使用
provide('showTableDialog', showTableDialog)

const pageStore = usePageStore()
</script>

<template>
  <div class="app-container">
    <TopMenu />
    <div class="dnd-flow">
      <LeftSidebar/>
      <div class="flow-container">
        <FlowTabs />
        <router-view style="background-color: var(--card-bg);" />
      </div>
      <TableViewDialog
        v-model="dialogVisible"
        :tableData="currentTableData"
      />
    </div>
  </div>
</template>

<style>
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  /* 防止内容溢出，确保主页面不显示滚动条 */
  overflow: hidden;
}

.dnd-flow {
  flex: 1;
  display: flex;
  position: relative;
}

.dnd-flow {
  flex-direction: row;
  display: flex;
  height: 100%
}

.flow-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-color-secondary);
}
</style>