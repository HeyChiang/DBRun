<template>
  <Dialog 
    :visible="visible" 
    @update:visible="$emit('update:visible', $event)"
    :header="'Connect to ' + dbType" 
    :modal="true"
    class="database-connection-dialog no-select"
    :style="{ width: '50vw' }"
  >
    <div class="connection-form">
      <div class="field-row">
        <div class="field-label">Label</div>
        <div class="field-input">
          <InputText v-model="formData.label" placeholder="My Database Connection" maxlength="50" />
        </div>
      </div>

      <div class="field-row-group">
        <div class="field-row">
          <div class="field-label">Host</div>
          <div class="field-input">
            <InputText v-model="formData.host" placeholder="localhost" />
          </div>
        </div>
        <div class="field-row">
          <div class="field-label">Port</div>
          <div class="field-input">
            <InputNumber v-model="formData.port" 
                        :useGrouping="false"
                        :min="0"
                        :max="65535" />
          </div>
        </div>
      </div>

      <div class="field-row">
        <div class="field-label">User</div>
        <div class="field-input">
          <InputText v-model="formData.username" placeholder="root" maxlength="100" />
        </div>
      </div>

      <div class="field-row">
        <div class="field-label">Password</div>
        <div class="field-input">
          <InputText type="password" v-model="formData.password" maxlength="100" />
        </div>
      </div>

      <div class="field-row">
        <div class="field-label">{{ dbType === 'oracle' ? 'SID' : 'Database' }}</div>
        <div class="field-input">
          <InputText v-model="formData.database" :placeholder="dbType === 'oracle' ? 'ORCL' : 'mydatabase'" maxlength="100" />
        </div>
      </div>

      <div class="field-row">
        <div class="field-label">URL</div>
        <div class="field-input">
          <InputText :value="connectionUrl" readonly class="url-display" maxlength="500"/>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <div class="left-buttons">
          <Button label="Test Connection" @click="testConnection" class="p-button-secondary" />
        </div>
        <div class="right-buttons">
          <Button label="Cancel" @click="closeDialog" class="p-button-text" />
          <Button v-if="!props.editingData" label="Connect" @click="connect" />
          <Button v-else label="Save" @click="save" />
        </div>
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { ref, defineProps, defineEmits, computed, watch } from 'vue';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import { InputNumber } from 'primevue';
import { useToast } from 'primevue/usetoast';
import { useDatabaseStore } from '@/stores/databaseStore';
import { TestConnection } from '@/../wailsjs/go/api/MetadatasAPI';
import { connect as connectModels } from '@/../wailsjs/go/models';

const toast = useToast();
const databaseStore = useDatabaseStore();

const props = defineProps({
  visible: {
    type: Boolean,
    required: true
  },
  dbType: {
    type: String,
    required: true
  },
  editingData: {
    type: Object,
    default: null
  }
});

const emit = defineEmits(['update:visible', 'connect','saved']);

const getDefaultPort = (dbType) => {
  switch (dbType.toLowerCase()) {
    case 'mysql':
      return 3306;
    case 'mariadb':
      return 3306;
    case 'postgresql':
      return 5432;
    case 'oracle':
      return 1521;
    default:
      return 3306;
  }
};

const getDefaultValues = (dbType) => {
  return {
    label: 'My Database Connection',
    host: 'localhost',
    port: getDefaultPort(dbType),
    username: 'root',
    password: '',
    database: dbType === 'oracle' ? 'ORCL' : 'mydatabase'
  };
};

const getValueOrDefault = (value, key, dbType) => {
  const defaults = getDefaultValues(dbType);
  if (key === 'port') {
    return value ?? defaults[key];
  }
  return (value !== undefined && value !== null && value !== '') ? value : defaults[key];
};

const formData = ref({
  label: props.editingData?.label || '',
  host: props.editingData?.host || '',
  port: props.editingData?.port || getDefaultPort(props.dbType),
  username: props.editingData?.username || '',
  password: props.editingData?.password || '',
  database: props.editingData?.database || '',
});

// 当编辑数据变化时，更新表单
watch(() => props.editingData, (newVal) => {
  if (newVal) {
    formData.value = { ...newVal };
  }
}, { immediate: true });

watch(() => props.dbType, () => {
  if (props.editingData) {
    formData.value = {
      label: getValueOrDefault(props.editingData.label, 'label', props.dbType),
      host: getValueOrDefault(props.editingData.host, 'host', props.dbType),
      port: getValueOrDefault(props.editingData.port, 'port', props.dbType),
      username: getValueOrDefault(props.editingData.username, 'username', props.dbType),
      password: props.editingData.password || '',
      database: getValueOrDefault(props.editingData.database, 'database', props.dbType),
    };
  } else {
    formData.value = getDefaultValues(props.dbType);
  }
}, { immediate: true });

const connectionUrl = computed(() => {
  const type = props.dbType.toLowerCase();
  const host = getValueOrDefault(formData.value.host, 'host', props.dbType);
  const port = getValueOrDefault(formData.value.port, 'port', props.dbType);
  const database = getValueOrDefault(formData.value.database, 'database', props.dbType);

  switch (type) {
    case 'mysql':
      return `jdbc:mysql://${host}:${port}/${database}`;
    case 'mariadb':
      return `jdbc:mariadb://${host}:${port}/${database}`;
    case 'postgresql':
      return `jdbc:postgresql://${host}:${port}/${database}`;
    case 'oracle':
      return `jdbc:oracle:thin:@${host}:${port}:${database}`;
    default:
      return '';
  }
});

const closeDialog = () => {
  emit('update:visible', false);
};

const connect = async () => {
  const connectionData = {
    type: props.dbType,
    label: getValueOrDefault(formData.value.label, 'label', props.dbType),
    host: getValueOrDefault(formData.value.host, 'host', props.dbType),
    port: getValueOrDefault(formData.value.port, 'port', props.dbType),
    username: getValueOrDefault(formData.value.username, 'username', props.dbType),
    password: formData.value.password,
    database: getValueOrDefault(formData.value.database, 'database', props.dbType)
  };

  const cfg = connectModels.Config.createFrom({
    id: 0,
    type: props.dbType,
    label: connectionData.label,
    host: connectionData.host,
    port: connectionData.port,
    username: connectionData.username,
    password: connectionData.password,
    database: connectionData.database,
    instance: '',
    options: ''
  });

  try {
    await TestConnection(cfg);
    emit('connect', connectionData);
    closeDialog();
  } catch (err) {
    const message = err?.message || String(err);
    toast.add({ severity: 'error', summary: '连接失败', detail: message });
  }
};

const save = async () => {
  try {
    const dbData = {
      id: props.editingData?.id,
      label: formData.value.label,
      host: formData.value.host,
      port: formData.value.port,
      username: formData.value.username,
      password: formData.value.password,
      database: formData.value.database,
      // 必须使用后端结构定义的字段名：type，而不是dbType
      type: props.dbType
    };
    
    console.log('Saved database:', dbData , "props.editingData", props.editingData);

    if (props.editingData) {
      // Update existing database
      await databaseStore.updateDBCredentials(dbData);
    } else {
      // Add new database
      await databaseStore.addDatabase(dbData);
    }
    emit('update:visible', false);
    emit('saved');
  } catch (err) {
    console.error('保存数据库连接错误:', err);
    toast.add({ severity: 'error', summary: '错误', detail: '保存数据库连接失败: ' + err.message });
  }
};

const testConnection = async () => {
  const cfg = connectModels.Config.createFrom({
    id: 0,
    type: props.dbType,
    label: getValueOrDefault(formData.value.label, 'label', props.dbType),
    host: getValueOrDefault(formData.value.host, 'host', props.dbType),
    port: getValueOrDefault(formData.value.port, 'port', props.dbType),
    username: getValueOrDefault(formData.value.username, 'username', props.dbType),
    password: formData.value.password,
    database: getValueOrDefault(formData.value.database, 'database', props.dbType),
    instance: '',
    options: ''
  });

  try {
    await TestConnection(cfg);
    toast.add({ severity: 'success', summary: '成功', detail: '数据库连接正常', life: 2000, closable: false });
  } catch (err) {
    console.error('测试数据库连接错误:', err);
    const message = err?.message || String(err);
    toast.add({ severity: 'error', summary: '连接失败', detail: message });
  }
};
</script>

<style scoped>
.database-connection-dialog {
  min-width: 400px;
}

.connection-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1rem 0;
}

.field-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.field-row-group {
  display: flex;
  flex-direction: row;
  gap: 1rem;
}

.field-label {
  min-width: 80px;
  color: var(--p-surface-900) !important;
  font-weight: 500;
  text-align: left;
}

.field-input {
  flex: 1;
}

.url-display {
  background-color: var(--surface-100);
  color: var(--text-color-secondary);
}

.form-group label {
  font-weight: 500;
  color: var(--text-color);
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.right-buttons {
  display: flex;
  gap: 0.5rem;
}

:deep(.p-inputtext) {
  width: 100%;
}

.field-input :deep(input),
.field-input :deep(.p-inputtext),
.field-input :deep(.p-inputnumber-input),
.url-display {
  text-align: left;
}

:deep(.p-dark) .database-connection-dialog .field-label {
  color: var(--p-surface-50) !important;
}
</style>
