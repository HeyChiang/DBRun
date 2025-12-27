<template>
  <Dialog :visible="modelValue" @update:visible="$emit('update:modelValue', $event)" :header="tableData.table.name"
    :modal="true" class="p-fluid" :style="{ width: '70vw' }">
    <DataTable :value="filteredFields" :scrollable="true" scrollHeight="400px" :reorderableRows="true" @rowReorder="onRowReorder">
      <Column rowReorder header="排序" :style="{ width: '90px' }"></Column>
      <Column field="name" :style="{ width: '150px' }">
        <template #header>
          <div style="display: flex; align-items: center;">
            <template v-if="showNameFilterInput">
              <InputText id="name-filter-input" v-model="nameFilterText" placeholder="筛选字段名" style="width: 180px;"
                @blur="showNameFilterInput = false"
                @keydown.esc.stop.prevent="showNameFilterInput = false" />
            </template>
            <template v-else>
              <span style="font-weight: bold;">字段名</span>
              <Button icon="pi pi-filter" class="p-button-text p-button-sm" style="margin-left:4px;"
                @click="handleFilterClick" />
            </template>
          </div>
        </template>
      </Column>
      <Column field="comment" header="注释" :style="{ width: '150px' }"></Column>
      <Column field="key" header="键类型" :style="{ width: '100px' }">
        <template #body="slotProps">
          <span :class="['key-type', getKeyTypeClass(slotProps.data.key)]">
            {{ getKeyTypeText(slotProps.data.key) }}
          </span>
        </template>
      </Column>
      <Column field="display" header="展示" :style="{ width: '120px' }">
        <template #header>
          <div class="flex align-items-center justify-content-center gap-2">
            <Checkbox v-model="allDisplay" @change="toggleAllDisplay" :binary="true" />
          </div>
        </template>
        <template #body="slotProps">
          <Checkbox v-model="slotProps.data.display" :binary="true" style="transform: translateX(50%);" />
        </template>
      </Column>
      <Column field="remark" header="说明">
        <template #body="slotProps">
          <InputText v-model="slotProps.data.remark" placeholder="请输入说明" class="w-full" />
        </template>
      </Column>
    </DataTable>
    <template #footer>
      <div class="flex justify-content-end">
        <Button label="取消" icon="pi pi-times" @click="handleCancel" class="p-button-text" />
        <Button label="确认" icon="pi pi-check" @click="handleConfirm" autofocus />
      </div>
    </template>
  </Dialog>
</template>

<script setup lang="ts">
import { computed, ref, watch, nextTick } from 'vue'
import Dialog from 'primevue/dialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Checkbox from 'primevue/checkbox'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import { getKeyTypeText, getKeyTypeClass } from '@/components/flow/tableUtils'
import { service, models } from '@/../wailsjs/go/models'
import { UpdateFieldsSortByTableID } from '@/../wailsjs/go/api/MetadatasAPI'
import { useDatabaseStore } from '@/stores/databaseStore'
import { TableNodeInfo } from '@/types/tableTypes'
import { eventBus } from '@/utils/eventBus'

const props = defineProps<{
  modelValue: boolean;
  tableData: TableNodeInfo;
}>()

// 筛选字段名相关变量
const showNameFilterInput = ref(false)
const nameFilterText = ref('')

// autofocus在vue3里没效果
// 就不再使用ref引用输入框，改为直接使用DOM选择器
const handleFilterClick = () => {
  showNameFilterInput.value = true;
  // 等待DOM更新
  setTimeout(() => {
    const input = document.getElementById('name-filter-input');
    if (input) {
      // 如果是InputText工件，实际的输入元素在它内部
      const realInput = input.querySelector('input') || input;
      if (realInput instanceof HTMLElement) {
        realInput.focus();
        console.log('成功设置焦点到输入框');
      }
    }
  }, 100);
}

// 筛选后的字段
const filteredFields = computed(() => {
  if (!nameFilterText.value) return fields.value
  return fields.value.filter(f => f.name && f.name.includes(nameFilterText.value))
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
}>()

const allDisplay = ref(true)
const fields = ref<models.FieldInfoVO[]>([])

// 初始化数据
const initFields = () => {
  console.log('传入参数 TableViewDialog init', props.tableData)

  // 检查是否为TableInfoVO类型
  const tableInfo = props.tableData.table as models.TableInfoVO;
  fields.value = tableInfo.fields || [];
  fields.value.forEach(field => {
    if (field.display !== false) {
      field.display = true
    }
    if (!field.remark) {
      field.remark = ''
    }
  })
  updateAllDisplayState()
}

// 更新全选状态
const updateAllDisplayState = () => {
  allDisplay.value = fields.value.every(field => field.display)
}

// 监听visible变化，当打开对话框时初始化数据
watch(() => props.modelValue, (newVal) => {
  if (newVal) {
    initFields()
    nameFilterText.value = ''
    showNameFilterInput.value = false
  }
})

// 全选/取消全选
const toggleAllDisplay = () => {
  fields.value.forEach(field => {
    field.display = allDisplay.value
  })
}

const handleCancel = () => {
  emit('update:modelValue', false)
}

// 拖拽排序后更新本地排序并保存到数据库
const onRowReorder = async (event: any) => {
  try {
    const newOrder: models.FieldInfoVO[] = event.value || []
    // 更新显示顺序
    fields.value = newOrder
    fields.value.forEach((f, i) => {
      (f as any).sort = i + 1
    })

    // 构建ID数组并调用后端保存排序
    const tableInfo = props.tableData.table as models.TableInfoVO
    const tableId = (tableInfo as any)?.id ?? props.tableData?.tableId
    const fieldIDs = fields.value.map(f => (f as any)?.id).filter((id): id is number => typeof id === 'number')
    if (tableId && fieldIDs.length > 0) {
      await UpdateFieldsSortByTableID(tableId, fieldIDs)
    }
  } catch (err) {
    console.error('更新字段排序失败:', err)
  }
}

  const handleConfirm = async () => {
    try {
      const databaseStore = useDatabaseStore();
      console.log('修改表信息 handleConfirm fields', fields.value)

      // 更新tableData对象的字段信息
      const tableInfo = props.tableData.table as models.TableInfoVO;
      tableInfo.fields = fields.value;

      // 调用新的updateTableRemark方法保存remark信息，传入可选tableId
      await databaseStore.updateTableRemark(
      tableInfo,
      (tableInfo as any)?.id ?? props.tableData?.tableId
      );

      // 使用事件总线发送表更新事件，包含tableId
      eventBus.emit('table-remark-updated', {
        tableId: (tableInfo as any)?.id ?? props.tableData?.tableId,
        dbId: props.tableData.dbId,
        dbName: props.tableData.dbName,
        schemaName: props.tableData.schemaName,
        tableName: props.tableData.table.name,
        tableInfo: tableInfo
      });

    emit('update:modelValue', false);
  } catch (error) {
    console.error('Failed to update fields:', error);
  }
}

// 监听单个字段display的变化来更新全选状态
watch(() => fields.value.map(f => f.display), () => {
  updateAllDisplayState()
}, { deep: true })
</script>

<style scoped>
.key-type {
  font-weight: bold;
  padding: 2px 6px;
  border-radius: 4px;
}

.primary-key {
  background-color: #2ecc71;
  color: #ffffff;
}

.foreign-key {
  background-color: #3498db;
  color: #ffffff;
}

.unique-key {
  background-color: #f39c12;
  color: #ffffff;
}

.index-key {
  background-color: #95a5a6;
  color: #ffffff;
}

:deep(.p-inputtext) {
  width: 100%;
}
</style>
