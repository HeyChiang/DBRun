<template>
  <div class="table-view" >
    <h2 class="table-name">{{ props.data.table.name }}</h2>
    <div class="table-content">
      <div v-for="(field, index) in tableFields"
           :key="field.name"
           :class="['table-row', index % 2 === 0 ? 'even' : 'odd']"
      >
        <div class="field-item1">{{ field.name }}</div>
        <div class="field-item2">{{  field.remark || field.comment }}</div>
        <div class="field-item3">
          <span :class="['key-type', getKeyTypeClass(field.key)]">
            {{ getKeyTypeText(field.key) }}
          </span>
        </div>

        <Handle v-for="mark in handleMarks"
                :key="mark.id"
                :type="mark.type"
                :position="mark.position"
                :id="`${mark.id}${field.name}`"
                :style="{ top: `${(index + 0.5) * 100 / tableFields.length}%`, transform: 'translateY(-50%)', height: '2rem', width: '2rem', opacity: '0' }"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Handle, Position, type HandleType, XYPosition } from '@vue-flow/core'
import { getKeyTypeText, getKeyTypeClass } from './tableUtils'
import { models } from '@/../wailsjs/go/models'
import { onMounted, computed, ref } from 'vue';
import { ParseViewSQL, GetFieldsVOByTableID } from '@/../wailsjs/go/api/MetadatasAPI'
import { TableNodeInfo } from '@/types/tableTypes'

// 定义props类型
const props = defineProps<{
  id: string,
  type: string,
  position: XYPosition,
  data: TableNodeInfo
}>()

// 判断是否为视图类型
const isView = computed(() => {
  return props.data.table && 'definition' in props.data.table;
});

// 获取表字段的函数
const tableFields = computed(() => {
  if (!props.data.table) return [];
  
  // 如果是TableInfoVO，直接返回fields属性
  if (!isView.value) {
    const fields = (props.data.table as models.TableInfoVO).fields || [];
    // 过滤掉display为false的字段
    return fields.filter(field => field.display !== false);
  }
  
  // 如果是ViewInfoVO类型并且已解析过的字段存在，则返回解析后的字段
  // 过滤掉display为false的字段
  return (parsedViewFields.value || []).filter(field => field.display !== false);
})

// 存储解析后的视图字段
const parsedViewFields = ref<models.FieldInfoVO[]>([]);

// 解析视图SQL定义的函数
const parseViewDefinition = async () => {
  if (!isView.value) return;
  
  try {
    const viewInfo = props.data.table as models.ViewInfoVO;
    if (viewInfo.definition) {
      // 调用后端API解析视图SQL定义
      const fields = await ParseViewSQL(viewInfo.definition);
      parsedViewFields.value = fields;
    }
  } catch (error) {
    console.error('解析视图SQL失败:', error);
  }
};

onMounted(async () => {
  try {
    // 对于视图类型，解析SQL定义获取字段信息
    if (isView.value) {
      await parseViewDefinition();
    }
    
    // 仅使用表ID进行缓存查询；如果没有ID则跳过并提示
    const initialTableId: number | undefined = (props.data.table as any)?.id ?? (props.data as any)?.tableId;
    if (!initialTableId) {
      console.warn('未发现表ID，跳过缓存加载');
      return;
    }

    // 写回以便后续复用
    (props.data as any).tableId = initialTableId;

    // 刷新该表的字段VO，确保排序与显示配置为最新
    try {
      if (!isView.value) {
        const latestFields = await GetFieldsVOByTableID(initialTableId as number)
        if (Array.isArray(latestFields) && latestFields.length > 0) {
          ;(props.data.table as models.TableInfoVO).fields = latestFields as any
          console.log('[TableNodeMounted] 已刷新字段VO', {
            tableId: initialTableId,
            count: latestFields.length,
            order: latestFields.map((f: any) => ({ name: f?.name, sort: (f as any)?.sort }))
          })
        }
      }
    } catch (e) {
      console.warn('[TableNodeMounted] 刷新字段VO失败，沿用现有字段', e)
    }

    // 仅依赖 GetFieldsVOByTableID 返回的字段显示与备注，无需额外缓存合并

    // 打印一次节点挂载时的字段顺序，辅助排查排序问题
    try {
      const currentTableId: number | undefined = (props.data.table as any)?.id ?? (props.data as any)?.tableId;
      const currentTableName: string | undefined = (props.data.table as any)?.name;
      const order = (tableFields.value || []).map((f: any) => ({ name: f?.name, sort: (f as any)?.sort }))
      console.log('[TableNodeMounted] TableNode已挂载:', {
        tableId: currentTableId,
        tableName: currentTableName,
        fieldOrder: order
      })
    } catch (e) {
      console.debug('[TableNodeMounted] 打印字段顺序失败:', e)
    }
  } catch (error) {
    console.error('获取表和字段信息失败:', error);
  }
})

const handleMarks = [
  { type: 'source' as HandleType, position: Position.Left, id: 'sl-' },
  { type: 'target' as HandleType, position: Position.Left, id: 'tl-' },
  { type: 'source' as HandleType, position: Position.Right, id: 'sr-' },
  { type: 'target' as HandleType, position: Position.Right, id: 'tr-' },
]
</script>

<style scoped>
.table-view {
  width: 16rem;
  margin: 0 auto;
  background-color: #ffffff;
  color: black;
  font-size: 0.8rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.table-name {
  background-color: #3498db;
  color: #ffffff;
  padding: 0.5rem;
  font-size: 0.8rem;
  margin: 0;
  text-align: center;
}

.table-content {
  position: relative;
}

.table-row {
  display: flex;
  flex-wrap: nowrap;
  padding: 0.6rem;
  transition: background-color 0.3s ease;
}

.table-row:hover {
  background-color: #f1f3f5;
}

.table-row.even {
  background-color: #F0F8FF;
}

.table-row.odd {
  background-color: #ffffff;
}

.field-item1 {
  flex: 4 1 0;
  text-align: left;
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
  padding: 0 5px;
}

.field-item2 {
  flex: 4 1 0;
  text-align: left;
  overflow: hidden;
  color: #777777;
  white-space: nowrap;
  text-overflow: ellipsis;
  padding: 0 5px;
}

.field-item3 {
  flex: 1 1 0;
}

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
</style>
