<script setup lang="ts">
import { ref, toRaw, reactive } from 'vue'
import SvgIcon from '@/components/SvgIcon.vue'
import { ListDatabasesByConfig, SyncTableFieldsByTableID, SyncSchemaByID, SyncDatabaseByID, GetFieldsVOByTableID, GetTablesVOByDatabaseID } from '@/../wailsjs/go/api/MetadatasAPI'
import useDragAndDrop from './useDragAndDrop'
import { connect, metadata, models } from '@/../wailsjs/go/models'
import { useDatabaseStore } from '@/stores/databaseStore'
import { useConfirm } from 'primevue/useconfirm'
import ConfirmDialog from 'primevue/confirmdialog'
import DBConnectionDialog from '@/components/sidebar/DBConnectionDialog.vue'

type TableMeta = models.TableInfoVO
type DatabaseMeta = models.DatabaseInfoVO
type SchemaMeta = models.SchemaVO
type DbLink = metadata.Credentials & { dbs?: DatabaseMeta[] }

const props = defineProps<{ dbLinks: DbLink[]; canDrag: boolean; searchTerm?: string }>()

const databaseStore = useDatabaseStore()
const confirm = useConfirm()
const showDatabaseConnection = ref(false)
const editingDBLink = ref<any>(null)
const editingDbType = ref('')

const expandedMenus = ref<Set<string | number>>(new Set())
const selectedItem = ref<any>(null)
const showContextMenu = ref(false)
const contextMenuPosition = ref<{ x: number; y: number }>({ x: 0, y: 0 })
const contextMenuTarget = ref<{ level: 'link' | 'database' | 'schema' | 'table'; item: DbLink | DatabaseMeta | SchemaMeta | TableMeta; link?: DbLink; database?: DatabaseMeta } | null>(null)

// 刷新时的加载状态：支持多个并发刷新目标（数据库与表独立记录）
const refreshingTables = reactive(new Set<number | string>())
const refreshingDatabases = reactive(new Set<number | string>())
const refreshingSchemas = reactive(new Set<number | string>())
const getTableId = (t: any) => t?.id ?? t?.tableId
const isTableRefreshing = (t: any) => {
  const tid = getTableId(t)
  return tid !== undefined && tid !== null && refreshingTables.has(tid)
}
const isDatabaseRefreshing = (d: any) => {
  const did = d?.id
  return did !== undefined && did !== null && refreshingDatabases.has(did)
}

const isSchemaRefreshing = (s: any) => {
  const sid = s?.id
  return sid !== undefined && sid !== null && refreshingSchemas.has(sid)
}

const getFilteredTables = (tables: any[], term?: string) => {
  const q = (term || '').trim().toLowerCase()
  if (!q) return tables || []
  return (tables || []).filter((t: any) => String((t?.name || '')).toLowerCase().includes(q))
}

const collectAllTables = (database: any) => {
  const base = Array.isArray(database?.tables) ? database.tables : []
  const schemas = Array.isArray(database?.schemas) ? database.schemas : []
  const merged: any[] = [...base]
  for (const sch of schemas) {
    const st = Array.isArray(sch?.tables) ? sch.tables : []
    merged.push(...st)
  }
  const dedup: any[] = []
  const seen = new Set<string | number>()
  for (const t of merged) {
    const id = (t as any)?.id ?? (t as any)?.tableId ?? (t as any)?.name
    if (id !== undefined && id !== null && !seen.has(id)) {
      seen.add(id)
      dedup.push(t)
    }
  }
  return dedup
}

const getFilteredDatabases = (dbs: any[], term?: string) => {
  const q = (term || '').trim().toLowerCase()
  if (!q) return dbs || []
  return (dbs || []).filter((db: any) => {
    const tables = collectAllTables(db)
    return getFilteredTables(tables, q).length > 0
  })
}

const isOracleLink = (link: any) => String(link?.type || '').toLowerCase() === 'oracle'
const collectSchemas = (link: any) => {
  const res: any[] = []
  const dbs = Array.isArray(link?.dbs) ? link.dbs : []
  for (const db of dbs) {
    const schemas = Array.isArray((db as any)?.schemas) ? (db as any).schemas : []
    for (const sch of schemas) {
      res.push({ schema: sch, database: db })
    }
  }
  return res
}
const getFilteredSchemas = (pairs: any[], term?: string) => {
  const q = (term || '').trim().toLowerCase()
  if (!q) return pairs || []
  return (pairs || []).filter((p: any) => {
    const tables = Array.isArray(p?.schema?.tables) ? p.schema.tables : []
    return getFilteredTables(tables, q).length > 0
  })
}

const toggleMenu = (key: string | number) => {
  if (expandedMenus.value.has(key)) {
    expandedMenus.value.delete(key)
  } else {
    expandedMenus.value.add(key)
  }
}

const handleItemClick = (item: any) => {
  selectedItem.value = item
}

const handleContextMenu = (event: MouseEvent, item: DbLink | DatabaseMeta | SchemaMeta | TableMeta, link?: DbLink, database?: DatabaseMeta) => {
  event.preventDefault()
  showContextMenu.value = true
  contextMenuPosition.value = { x: event.clientX, y: event.clientY }
  const level: 'link' | 'database' | 'schema' | 'table' = database ? (((item as any)?.fields !== undefined) ? 'table' : 'schema') : (link ? 'database' : 'link')
  contextMenuTarget.value = { level, item, link, database }
}

// 通过表ID一次性同步并拉取字段VO，避免逐字段请求
const refreshTableFieldsById = async (table: TableMeta) => {
  const tableId = (table as any)?.id ?? (table as any)?.tableId
  if (!tableId) {
    console.warn('无法同步：未发现表ID')
    return
  }
  console.log('[SidebarRefresh] 刷新表开始', {
    tableId,
    tableName: (table as any)?.name,
  })
  // 同步表的字段元数据
  await SyncTableFieldsByTableID(tableId as number)
  console.log('[SidebarRefresh] SyncTableFieldsByTableID 完成', { tableId })
  // 同步后拉取该表的最新字段VO并替换
  const fieldsVO = await GetFieldsVOByTableID(tableId as number)
  ;(table as any).fields = fieldsVO
  console.log('[SidebarRefresh] 已获取最新字段VO', { tableId, fieldCount: Array.isArray(fieldsVO) ? fieldsVO.length : 0 })
}

const handleEdit = () => {
  if (contextMenuTarget.value && contextMenuTarget.value.level === 'link') {
    const item = contextMenuTarget.value.item as any
    editingDBLink.value = { ...toRaw(item) }
    editingDbType.value = item.type
    showDatabaseConnection.value = true
    showContextMenu.value = false
  }
}

  const handleDelete = () => {
    if (contextMenuTarget.value && contextMenuTarget.value.level === 'link') {
      const item = contextMenuTarget.value.item as any
      confirm.require({
        message: '确定要删除这个数据库连接吗？',
        header: '删除确认',
        icon: 'pi pi-exclamation-triangle',
        acceptClass: 'p-button-danger',
        accept: async () => {
          try {
            if (item?.id) {
              await databaseStore.removeDatabase(item.id)
            }
          } catch (error) {
            console.error('删除失败:', error);
          }
        }
      })
      showContextMenu.value = false
    }
  }

// 刷新（数据库或表）
const handleRefresh = async () => {
  // 记录本次刷新开始的目标，用于 finally 精确清理
  let startedLevel: 'table' | 'database' | 'schema' | null = null
  let startedId: number | string | undefined
  try {
    const ctx = contextMenuTarget.value
    console.log('[SidebarRefresh] 触发刷新', {
      level: ctx?.level,
      linkId: (ctx as any)?.link?.id,
      databaseId: (ctx as any)?.item?.id,
      databaseName: (ctx as any)?.item?.name,
      tableName: (ctx as any)?.item?.name
    })
    if (!ctx || ctx.level === 'link') return

    // 进入刷新：设置加载状态，驱动左侧图标展示为加载中（支持并发）
    if (ctx.level === 'table') {
      const tid = getTableId(ctx.item as any)
      if (tid !== undefined && tid !== null) {
        refreshingTables.add(tid)
        startedLevel = 'table'
        startedId = tid
      }
    } else if (ctx.level === 'database') {
      const did = (ctx.item as any)?.id
      if (did !== undefined && did !== null) {
        refreshingDatabases.add(did)
        startedLevel = 'database'
        startedId = did
      }
    } else if (ctx.level === 'schema') {
      const sid = (ctx.item as any)?.id
      if (sid !== undefined && sid !== null) {
        refreshingSchemas.add(sid)
        startedLevel = 'schema'
        startedId = sid
      }
    }

    if (ctx.level === 'table' && ctx.link && ctx.database && ctx.item) {
      // 基于表ID同步字段，仅更新当前表数据（一次性按表ID获取字段）
      try {
        await refreshTableFieldsById(ctx.item as TableMeta)
      } catch (e) {
        console.warn('[SidebarRefresh] 获取最新字段失败，回退为全量刷新该连接', e)
        const dbInfo = await ListDatabasesByConfig((ctx.link as unknown as connect.Config))
        if (ctx.link) {
          ;(ctx.link as any).dbs = (dbInfo as any).dbs
          console.log('[SidebarRefresh] 全量刷新完成', { linkId: (ctx.link as any)?.id, dbCount: Array.isArray((dbInfo as any)?.dbs) ? (dbInfo as any).dbs.length : 0 })
        }
      }
    } else if (ctx.level === 'database' && ctx.link && ctx.item) {
      // 刷新当前数据库（或schema）下的表，不做全量连接刷新
      const database = ctx.item as any
      const databaseId = database?.id as number | undefined
      if (!databaseId) {
        console.warn('刷新跳过：未发现数据库ID')
        return
      }
      console.log('[SidebarRefresh] 刷新数据库/Schema开始', {
        linkId: (ctx.link as any)?.id,
        databaseId,
        dbName: database?.name
      })

      // 若包含schemas，尝试按schema逐个同步并获取表；否则仅获取数据库直挂的表
      let mergedTables: any[] = []
      try {
        if (Array.isArray(database.schemas) && database.schemas.length > 0) {
          // 先同步每个schema（避免数据过期），再精确拉取对应表
          for (const sch of database.schemas) {
            if (sch?.id) {
              try {
                console.log('[SidebarRefresh] 同步Schema开始', { schemaId: sch.id, schemaName: sch.name })
                await SyncSchemaByID(sch.id as number)
                console.log('[SidebarRefresh] 同步Schema完成', { schemaId: sch.id })
              } catch (err) {
                console.warn('[SidebarRefresh] 同步Schema失败，继续', { schemaId: sch.id, error: err })
              }
              const schemaId = sch.id as number
              const tablesBySchema = await GetTablesVOByDatabaseID(databaseId as number, schemaId as any)
              console.log('[SidebarRefresh] 获取Schema表完成', { schemaId, count: Array.isArray(tablesBySchema) ? tablesBySchema.length : 0 })
              if (Array.isArray(tablesBySchema)) mergedTables = mergedTables.concat(tablesBySchema)
            }
          }
        } else {
          // 无schema型数据库（如 MySQL/MariaDB）：直接按数据库ID进行一次完整同步
          console.log('[SidebarRefresh] 无Schema，执行数据库级同步', { databaseId })
          await SyncDatabaseByID(databaseId as number)
        }
        // 同时获取无schema归属的表
        const plainTables = await GetTablesVOByDatabaseID(databaseId as number, null as any)
        console.log('[SidebarRefresh] 获取无Schema表完成', { count: Array.isArray(plainTables) ? plainTables.length : 0 })
        if (Array.isArray(plainTables)) mergedTables = mergedTables.concat(plainTables)

        // 写回当前数据库的表列表
        database.tables = mergedTables
        console.log('[SidebarRefresh] 刷新数据库/Schema完成', { databaseId, totalTableCount: mergedTables.length })
      } catch (e) {
        console.warn('[SidebarRefresh] 按库/Schema刷新失败，回退为全量刷新该连接', e)
        const dbInfo = await ListDatabasesByConfig((ctx.link as unknown as connect.Config))
        if (ctx.link) {
          (ctx.link as any).dbs = (dbInfo as any).dbs
          console.log('[SidebarRefresh] 全量刷新完成', { linkId: (ctx.link as any)?.id, dbCount: Array.isArray((dbInfo as any)?.dbs) ? (dbInfo as any).dbs.length : 0 })
        }
      }
    } else if (ctx.level === 'schema' && ctx.link && ctx.database && ctx.item) {
      const schema = ctx.item as any
      const schemaId = schema?.id as number | undefined
      const databaseId = (ctx.database as any)?.id as number | undefined
      if (!schemaId || !databaseId) {
        console.warn('刷新跳过：未发现Schema或数据库ID')
        return
      }
      try {
        await SyncSchemaByID(schemaId as number)
        const tablesBySchema = await GetTablesVOByDatabaseID(databaseId as number, schemaId as any)
        schema.tables = Array.isArray(tablesBySchema) ? tablesBySchema : []
      } catch (e) {
        console.warn('[SidebarRefresh] 刷新Schema失败，回退为全量刷新该连接', e)
        const dbInfo = await ListDatabasesByConfig((ctx.link as unknown as connect.Config))
        if (ctx.link) {
          (ctx.link as any).dbs = (dbInfo as any).dbs
        }
      }
    }
  } catch (err) {
    console.error('刷新失败:', err)
  } finally {
    showContextMenu.value = false
    // 刷新完成：仅移除当前项的加载状态，恢复原图标
    try {
      if (startedLevel === 'table' && startedId !== undefined && startedId !== null) {
        refreshingTables.delete(startedId)
      } else if (startedLevel === 'database' && startedId !== undefined && startedId !== null) {
        refreshingDatabases.delete(startedId)
      } else if (startedLevel === 'schema' && startedId !== undefined && startedId !== null) {
        refreshingSchemas.delete(startedId)
      }
    } catch {}
  }
}

const { onDragStart } = useDragAndDrop()

const handleDragStart = (event: DragEvent, dbId: number | string | undefined, dbName: string, schemaName: string, table: TableMeta) => {
  try {
    const fields = (table as any)?.fields || []
    console.log('[SidebarDragStart] 拖拽创建TableNode:', {
      dbId: dbId,
      dbName,
      schemaName,
      tableId: (table as any)?.id,
      tableName: (table as any)?.name,
      fieldOrder: Array.isArray(fields) ? fields.map((f: any) => ({ name: f?.name, sort: (f as any)?.sort })) : []
    })
  } catch (e) {
    console.debug('[SidebarDragStart] 打印字段顺序失败:', e)
  }
  onDragStart(event, 'table-node', {
    dbId: String(dbId ?? ''),
    dbName,
    schemaName,
    table: table as unknown as models.TableInfoVO
  })
}

document.addEventListener('click', () => {
  showContextMenu.value = false
})
</script>

<template>
  <div class="sidebar-content no-select">
    <ul class="list-none px-1 m-0">
      <li v-for="(dbLink, dbIndex) in props.dbLinks" :key="dbIndex">
        <div
          class="menu-item flex items-center justify-between"
          @click="toggleMenu(dbIndex); handleItemClick(dbLink)"
          @contextmenu="handleContextMenu($event, dbLink)"
          :class="{ 'selected': selectedItem === dbLink }"
        >
          <div class="flex items-center">
            <i class="pi pi-chevron-right text-gray-400" style="font-size: 0.8rem;" :class="{ 'rotate-90': expandedMenus.has(dbIndex) }"></i>
            <SvgIcon :name="'db_' + dbLink.type.toLowerCase()" class="mx-2" />
            <span class="font-medium text-sm whitespace-nowrap">{{ dbLink.label || (dbLink.host + ':' + String(dbLink.port)) }}</span>
          </div>
        </div>

        <ul v-show="expandedMenus.has(dbIndex)" class="list-none py-1 m-0 overflow-hidden">
          <template v-if="isOracleLink(dbLink)">
            <li v-for="(pair, pairIndex) in getFilteredSchemas(collectSchemas(dbLink), props.searchTerm)" :key="pairIndex">
              <div
                class="menu-item flex items-center justify-between"
                @click="toggleMenu(`${dbIndex}-schema-${pairIndex}`); handleItemClick(pair.schema)"
                @contextmenu="handleContextMenu($event, pair.schema, dbLink, pair.database)"
                :class="{ 'selected': selectedItem === pair.schema }"
              >
                <div class="flex items-center">
                  <i class="pi pi-chevron-right text-gray-400 ml-5" style="font-size: 0.8rem;" :class="{ 'rotate-90': expandedMenus.has(`${dbIndex}-schema-${pairIndex}`) }"></i>
                  <template v-if="!isSchemaRefreshing(pair.schema)">
                    <SvgIcon name="db_schema" class="mx-2" :width="'1.1rem'" :height="'1.1rem'" />
                  </template>
                  <template v-else>
                    <i class="pi pi-spinner pi-spin mx-2" style="font-size: 1.1rem;"></i>
                  </template>
                  <span class="font-medium text-sm">{{ pair.schema.name }}</span>
                </div>
              </div>

              <ul v-if="expandedMenus.has(`${dbIndex}-schema-${pairIndex}`)" class="list-none py-2 pl-10 pr-4 overflow-hidden">
                <li v-for="(table, tableIndex) in getFilteredTables(pair.schema.tables, props.searchTerm)" :key="tableIndex" :draggable="props.canDrag"
                    @dragstart="handleDragStart($event, dbLink.id, pair.database.name, pair.schema.name, table)"
                    @click="handleItemClick(table)"
                    @contextmenu="handleContextMenu($event, table, dbLink, pair.database)"
                    :class="{ 'selected': selectedItem === table }" class="hover:bg-surface-700 cursor-pointer">
                  <a v-ripple class="menu-item flex items-center rounded">
                    <template v-if="!isTableRefreshing(table)">
                      <SvgIcon name="db_table" class="mr-2" :width="'1.1rem'" :height="'1.1rem'" />
                    </template>
                    <template v-else>
                      <i class="pi pi-spinner pi-spin mr-2" style="font-size: 1.1rem;"></i>
                    </template>
                    <span class="font-medium text-sm">{{ table.name }}</span>
                  </a>
                </li>
              </ul>
            </li>
          </template>
          <template v-else>
            <li v-for="(database, dbNameIndex) in getFilteredDatabases(dbLink.dbs || [], props.searchTerm)" :key="dbNameIndex">
              <div
                class="menu-item flex items-center justify-between"
                @click="toggleMenu(`${dbIndex}-${dbNameIndex}`); handleItemClick(database)"
                @contextmenu="handleContextMenu($event, database, dbLink)"
                :class="{ 'selected': selectedItem === database }"
              >
                <div class="flex items-center">
                  <i class="pi pi-chevron-right text-gray-400 ml-5" style="font-size: 0.8rem;" :class="{ 'rotate-90': expandedMenus.has(`${dbIndex}-${dbNameIndex}`) }"></i>
                  <template v-if="!isDatabaseRefreshing(database)">
                    <SvgIcon name="db_schema" class="mx-2" :width="'1.1rem'" :height="'1.1rem'" />
                  </template>
                  <template v-else>
                    <i class="pi pi-spinner pi-spin mx-2" style="font-size: 1.1rem;"></i>
                  </template>
                  <span class="font-medium text-sm">{{ database.name }}</span>
                </div>
              </div>

              <ul v-if="expandedMenus.has(`${dbIndex}-${dbNameIndex}`)" class="list-none py-2 pl-10 pr-4 overflow-hidden">
                <template v-if="Array.isArray(database.schemas) && database.schemas.length > 0">
                  <li v-for="(schema, schemaIndex) in (database.schemas || [])" :key="schemaIndex">
                    <div
                      class="menu-item flex items-center justify-between"
                      @click="toggleMenu(`${dbIndex}-${dbNameIndex}-schema-${schemaIndex}`); handleItemClick(schema)"
                      @contextmenu="handleContextMenu($event, schema, dbLink, database)"
                      :class="{ 'selected': selectedItem === schema }"
                    >
                      <div class="flex items-center">
                        <i class="pi pi-chevron-right text-gray-400 ml-5" style="font-size: 0.8rem;" :class="{ 'rotate-90': expandedMenus.has(`${dbIndex}-${dbNameIndex}-schema-${schemaIndex}`) }"></i>
                        <template v-if="!isSchemaRefreshing(schema)">
                          <SvgIcon name="db_schema" class="mx-2" :width="'1.1rem'" :height="'1.1rem'" />
                        </template>
                        <template v-else>
                          <i class="pi pi-spinner pi-spin mx-2" style="font-size: 1.1rem;"></i>
                        </template>
                        <span class="font-medium text-sm">{{ schema.name }}</span>
                      </div>
                    </div>

                    <ul v-if="expandedMenus.has(`${dbIndex}-${dbNameIndex}-schema-${schemaIndex}`)" class="list-none py-2 pl-10 pr-4 overflow-hidden">
                      <li v-for="(table, tableIndex) in getFilteredTables(schema.tables, props.searchTerm)" :key="tableIndex" :draggable="props.canDrag"
                          @dragstart="handleDragStart($event, dbLink.id, database.name, schema.name, table)"
                          @click="handleItemClick(table)"
                          @contextmenu="handleContextMenu($event, table, dbLink, database)"
                          :class="{ 'selected': selectedItem === table }" class="hover:bg-surface-700 cursor-pointer">
                        <a v-ripple class="menu-item flex items-center rounded">
                          <template v-if="!isTableRefreshing(table)">
                            <SvgIcon name="db_table" class="mr-2" :width="'1.1rem'" :height="'1.1rem'" />
                          </template>
                          <template v-else>
                            <i class="pi pi-spinner pi-spin mr-2" style="font-size: 1.1rem;"></i>
                          </template>
                          <span class="font-medium text-sm">{{ table.name }}</span>
                        </a>
                      </li>
                    </ul>
                  </li>

                  <li v-if="Array.isArray(database.tables) && database.tables.length > 0">
                    <div
                      class="menu-item flex items-center justify-between"
                      @click="toggleMenu(`${dbIndex}-${dbNameIndex}-schema-plain`); handleItemClick({ name: 'plain' })"
                      @contextmenu="handleContextMenu($event, database, dbLink)"
                    >
                      <div class="flex items-center">
                        <i class="pi pi-chevron-right text-gray-400 ml-5" style="font-size: 0.8rem;" :class="{ 'rotate-90': expandedMenus.has(`${dbIndex}-${dbNameIndex}-schema-plain`) }"></i>
                        <SvgIcon name="db_schema" class="mx-2" :width="'1.1rem'" :height="'1.1rem'" />
                        <span class="font-medium text-sm">无Schema</span>
                      </div>
                    </div>
                    <ul v-if="expandedMenus.has(`${dbIndex}-${dbNameIndex}-schema-plain`)" class="list-none py-2 pl-10 pr-4 overflow-hidden">
                      <li v-for="(table, tableIndex) in getFilteredTables(database.tables, props.searchTerm)" :key="tableIndex" :draggable="props.canDrag"
                          @dragstart="handleDragStart($event, dbLink.id, database.name, '', table)"
                          @click="handleItemClick(table)"
                          @contextmenu="handleContextMenu($event, table, dbLink, database)"
                          :class="{ 'selected': selectedItem === table }" class="hover:bg-surface-700 cursor-pointer">
                        <a v-ripple class="menu-item flex items-center rounded">
                          <template v-if="!isTableRefreshing(table)">
                            <SvgIcon name="db_table" class="mr-2" :width="'1.1rem'" :height="'1.1rem'" />
                          </template>
                          <template v-else>
                            <i class="pi pi-spinner pi-spin mr-2" style="font-size: 1.1rem;"></i>
                          </template>
                          <span class="font-medium text-sm">{{ table.name }}</span>
                        </a>
                      </li>
                    </ul>
                  </li>
                </template>
                <template v-else>
                  <li v-for="(table, tableIndex) in getFilteredTables(database.tables, props.searchTerm)" :key="tableIndex" :draggable="props.canDrag"
                      @dragstart="handleDragStart($event, dbLink.id, database.name, '', table)"
                      @click="handleItemClick(table)"
                      @contextmenu="handleContextMenu($event, table, dbLink, database)"
                      :class="{ 'selected': selectedItem === table }" class="hover:bg-surface-700 cursor-pointer">
                    <a v-ripple class="menu-item flex items-center rounded">
                      <template v-if="!isTableRefreshing(table)">
                        <SvgIcon name="db_table" class="mr-2" :width="'1.1rem'" :height="'1.1rem'" />
                      </template>
                      <template v-else>
                        <i class="pi pi-spinner pi-spin mr-2" style="font-size: 1.1rem;"></i>
                      </template>
                      <span class="font-medium text-sm">{{ table.name }}</span>
                    </a>
                  </li>
                </template>
              </ul>
            </li>
            <li v-if="!(dbLink.dbs && dbLink.dbs.length > 0)" class="px-4 py-2 text-gray-400 text-sm">暂无数据</li>
          </template>
        </ul>
      </li>
    </ul>

    <div v-if="showContextMenu" class="context-menu" :style="{ top: contextMenuPosition.y + 'px', left: contextMenuPosition.x + 'px' }">
      <template v-if="contextMenuTarget?.level === 'link'">
        <div class="context-menu-item" @click="handleEdit">
          <i class="pi pi-pencil mr-2"></i>编辑
        </div>
        <div class="context-menu-item" @click="handleDelete">
          <i class="pi pi-trash mr-2"></i>删除
        </div>
      </template>
      <template v-else>
        <div class="context-menu-item" @click="handleRefresh">
          <i class="pi pi-refresh mr-2"></i>刷新
        </div>
      </template>
    </div>

    <ConfirmDialog />

    <!-- Database Connection Dialog -->
    <DBConnectionDialog
      v-if="showDatabaseConnection"
      v-model:visible="showDatabaseConnection"
      :editing-data="editingDBLink"
      :db-type="editingDbType"
      @saved="databaseStore.refreshDatabases()"
    />
  </div>
</template>

<style scoped>
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