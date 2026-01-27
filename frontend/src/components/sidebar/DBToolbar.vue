<template>
  <Toolbar class="icon-toolbar">
    <template #start>
      <div class="relative">
        <Button icon="pi pi-plus" text severity="secondary" @click="toggle" aria-haspopup="true" aria-controls="db_menu" />
        <Menu ref="menu" id="db_menu" :model="menuItems" :popup="true" class="db-selection-menu">
          <template #item="{ item }">
            <div class="db-menu-item">
              <SvgIcon :name="item.icon" class="db-icon" />
              <span>{{ item.label }}</span>
            </div>
          </template>
        </Menu>
      </div>
      <Button v-if="!showSearch" icon="pi pi-search" text severity="secondary" class="ml-2" @click="toggleSearch" />
      <template v-if="showSearch">
        <InputText v-model="localSearch" placeholder="搜索表名" class="ml-2 search-input" @input="emitSearch" />
        <Button icon="pi pi-times" text severity="secondary" class="ml-1" @click="handleClearClick" />
      </template>
    </template>
  </Toolbar>

  <DBConnectionDialog 
      v-model:visible="showConnectionDialog"
      :db-type="selectedDbType"
      @connect="handleConnect"
    />
</template>

<script setup>
import { ref, nextTick } from 'vue';
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';
import Menu from 'primevue/menu';
import InputText from 'primevue/inputtext';
import DBConnectionDialog from './DBConnectionDialog.vue';
import { useDatabaseStore } from '@/stores/databaseStore';
import SvgIcon from '@/components/SvgIcon.vue';
import { useToast } from 'primevue/usetoast';
import { defineEmits, defineProps } from 'vue';

const menu = ref();
const showConnectionDialog = ref(false);
const selectedDbType = ref('');
const databaseStore = useDatabaseStore();
const toast = useToast();
const emit = defineEmits(['update:searchTerm']);
defineProps({ searchTerm: { type: String, required: false } });
const showSearch = ref(false);
const localSearch = ref('');

const menuItems = ref([
  {
    label: '选择数据库类型',
    items: [
      {
        label: 'MySQL',
        icon: 'db_mysql',
        command: () => {
          selectedDbType.value = 'mysql';
          showConnectionDialog.value = true;
        }
      },
      {
        label: 'MariaDB',
        icon: 'db_mariadb',
        command: () => {
          selectedDbType.value = 'mariadb';
          showConnectionDialog.value = true;
        }
      },
      {
        label: 'PostgreSQL',
        icon: 'db_postgresql',
        command: () => {
          selectedDbType.value = 'postgresql';
          showConnectionDialog.value = true;
        }
      },
      {
        label: 'Oracle',
        icon: 'db_oracle',
        command: () => {
          selectedDbType.value = 'oracle';
          showConnectionDialog.value = true;
        }
      }
    ]
  }
]);

const handleConnect = async (connectionData) => {
  console.log('Connection data:', connectionData);
  try {
    await databaseStore.addDatabase(connectionData);
    showConnectionDialog.value = false;
  } catch (err) {
    console.error('添加数据库失败:', err);
    toast.add({ severity: 'error', summary: '错误', detail: '添加数据库失败: ' + (err?.message || err) });
  }
};

const toggle = (event) => {
  menu.value.toggle(event);
};

const toggleSearch = () => {
  showSearch.value = true;
  nextTick(() => {
    const el = document.querySelector('.icon-toolbar .search-input')
    const fn = el ? Reflect.get(el, 'focus') : null
    if (typeof fn === 'function') fn.call(el)
  })
};

const emitSearch = () => {
  emit('update:searchTerm', localSearch.value);
};

const handleClearClick = () => {
  localSearch.value = ''
  emit('update:searchTerm', '')
  showSearch.value = false
}
</script>

<style scoped>
.icon-toolbar {
  border: none;
  padding: 0.2rem;
  background: transparent;
  /* border-top: 0.1rem solid var(--surface-border); */
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

/* Database selection menu styling */
:deep(.db-selection-menu .p-menu) {
  min-width: 220px;
  border-radius: 6px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

:deep(.db-selection-menu .p-submenu-header) {
  padding: 0.75rem 1rem;
  font-weight: 600;
  background-color: var(--surface-50);
  border-bottom: 1px solid var(--surface-200);
  margin: 0;
}

:deep(.db-selection-menu .p-menuitem) {
  margin: 0.25rem 0;
}

:deep(.db-selection-menu .p-menuitem-link) {
  padding: 0.75rem 1rem;
  transition: background-color 0.2s;
}

:deep(.db-selection-menu .p-menuitem-link:hover) {
  background-color: var(--surface-100);
}

.db-menu-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  width: 100%;
  padding: 0.25rem 0;
}

.db-icon {
  width: 1.5rem;
  height: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
