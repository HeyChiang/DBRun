<template>
  <div class="new-project-page">
    <div class="form-container p-card">
    <div class="form-body">
      <h2 class="title">新建项目</h2>
      <div class="p-fluid">
        <div class="field-row">
          <label for="projectName" class="field-label">项目名称</label>
          <div class="field-control">
            <InputText id="projectName" v-model="projectName" placeholder="请输入项目名称" />
          </div>
        </div>
        <div class="field-row">
          <label for="projectDir" class="field-label">项目目录</label>
          <div class="field-control">
            <div class="directory-input-group">
              <InputText id="projectDir" v-model="projectDir" placeholder="请选择项目目录" readonly style="cursor:pointer;" @click="selectDirectory" />
              <Button label="选择" icon="pi pi-folder-open" class="p-button-secondary" @click="selectDirectory" />
            </div>
          </div>
        </div>
        <div class="field-row">
          <div class="field-label"></div>
          <div class="field-control">
            <div class="form-actions">
              <Button label="取消" class="p-button-text" @click="cancel" />
              <Button label="创建" class="p-button-primary" @click="createProject" />
            </div>
          </div>
        </div>
            </div>
    </div>
  </div>
  </div>
  <Toast />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { InsertProject } from '@/../wailsjs/go/api/SQLiteAPI'
import { OpenDirectory, CreateDirectory } from '@/../wailsjs/go/api/SystemServiceAPI'
import { InitCache } from '@/../wailsjs/go/api/MetadatasAPI'
import { Init as InitAppCache } from '@/../wailsjs/go/api/AppCacheApi'
import { WindowSetSize } from '@/../wailsjs/runtime/runtime'
import { usePageStore } from '@/stores/pageStore'
import Toast from 'primevue/toast'
import { useToast } from 'primevue/usetoast'

const router = useRouter()
const projectName = ref('')
const projectDir = ref('')
const pageStore = usePageStore()
const toast = useToast()

const selectDirectory = async () => {
  try {
      OpenDirectory().then((path) => {
      console.log(path)
      projectDir.value = path
    })
  } catch (err) {
    console.error('选择目录失败:', err)
  }
}

const createProject = async () => {
  try {
    if (!projectName.value || !projectDir.value) {
      if (!projectName.value) {
        toast.add({ severity: 'warn', summary: '提示', detail: '请输入项目名称', life: 2000, closable: false })
      }
      if (!projectDir.value) {
        toast.add({ severity: 'warn', summary: '提示', detail: '请选择项目目录', life: 2000, closable: false })
      }
      return
    }
    // 0. 构建项目完整路径：在用户选择的目录下创建以项目名称命名的文件夹
    const projectPath = `${projectDir.value}\\${projectName.value}`
    
    // 1. 创建项目目录
    await CreateDirectory(projectPath)
    console.log('项目目录创建成功:', projectPath)
    
    // 2. 在SQLite中创建项目记录，使用完整的项目路径
    const newProject = await InsertProject(projectName.value, projectPath)
    if (newProject && newProject.id) {
      console.log('项目创建成功:', newProject)
      
      // 3. 初始化项目的 relation.db
      await InitCache(projectPath)
      console.log('relation.db 初始化完成:', projectPath)
      
      // 3.1 初始化应用缓存到项目目录
      await InitAppCache(projectPath)
      console.log('app_cache.bolt 初始化完成:', projectPath)
      await pageStore.initializeState()
      
      // 4. 设置大窗口尺寸
      WindowSetSize(1280, 720)
      
      // 5. 跳转到主页面
      router.push('/app')                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   
    } else {
      throw new Error('创建项目失败，未返回有效的项目信息。')
    }
  } catch (error: any) {
    console.error('创建项目失败:', error)
    toast.add({ severity: 'error', summary: '创建失败', detail: String(error), life: 2500 })
  }
}

const cancel = () => {
  router.back()
}
</script>

<style scoped>
.new-project-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 100px); /* 根据您的布局进行调整 */
  background: linear-gradient(135deg, rgba(0, 120, 212, 0.06), rgba(0, 120, 212, 0.03)), var(--surface-ground);
  color: var(--text-color);
  padding: 2rem;
}

.form-container {
  width: 100%;
  max-width: 600px;
  padding: 2rem 2.5rem;
  background-color: var(--card-bg);
  border: 1px solid var(--surface-border);
  border-radius: 0.75rem;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(8px);
}

.form-body {
  width: 100%;
  max-width: 400px;
  margin: 0 auto;
}


.title {
  text-align: center;
  font-size: 1.75rem;
  font-weight: 600;
  margin-bottom: 2.5rem;
  color: var(--text-color);
}

.field-row {
  display: flex;
  align-items: center;
  margin-bottom: 1.5rem;
}

.p-fluid .field-row:last-of-type {
  margin-bottom: 0;
}

/* 保证输入框左对齐 */
.field-control .p-inputtext {
  text-align: left;
}


.field-label {
  width: 100px;
  text-align: right;
  margin-right: 1.5rem;
  font-weight: 600;
  color: var(--text-color-secondary);
  flex-shrink: 0;
}

.field-control {
  flex: 1 1 auto;
  display: flex;
  align-items: center;
  justify-content: flex-start;
}

.directory-input-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  justify-content: flex-start;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}

/* 使只读输入框看起来更好 */
.p-inputtext:read-only {
    cursor: pointer !important;
}

.form-actions :deep(.p-button) {
  min-width: 90px;
  border-radius: 0.5rem;
  transition: background-color .2s, box-shadow .2s, color .2s, border-color .2s;
}

.form-actions :deep(.p-button-text) {
  color: var(--text-color-secondary);
  background: transparent;
  border: 1px solid transparent;
}

.form-actions :deep(.p-button-text:hover) {
  background-color: var(--surface-hover);
  color: var(--text-color);
  border-color: var(--surface-border);
}

.directory-input-group :deep(.p-button.p-button-secondary) {
  background-color: var(--surface-hover);
  color: var(--text-color);
  border: 1px solid var(--surface-border);
}

.directory-input-group :deep(.p-button.p-button-secondary:hover) {
  background-color: var(--surface-hover-light);
  color: var(--text-color);
  border-color: var(--surface-border);
}

.form-actions :deep(.p-button-primary) {
  background-image: linear-gradient(135deg, #0A84FF, #0078D4);
  color: #ffffff;
  border: none;
  box-shadow: 0 6px 18px rgba(10, 132, 255, 0.35);
  transition: transform .12s ease, box-shadow .12s ease, filter .12s ease;
}

.form-actions :deep(.p-button-primary:hover) {
  filter: brightness(1.04);
  box-shadow: 0 8px 22px rgba(10, 132, 255, 0.45);
  transform: translateY(-1px);
}

.form-actions :deep(.p-button:focus) {
  outline: none;
  box-shadow: 0 0 0 2px rgba(10, 132, 255, 0.35);
}

.form-actions :deep(.p-button-primary:disabled) {
  background-image: linear-gradient(135deg, #0A84FF, #0078D4);
  color: #ffffff;
  opacity: 1;
  cursor: not-allowed;
  box-shadow: none;
  filter: none;
}

.hint {
  margin-left: 0.75rem;
  font-size: 0.85rem;
}

.hint.error {
  color: var(--primary-color);
}

</style>
