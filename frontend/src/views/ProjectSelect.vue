<template>
  <div class="project-select-container">
    <!-- 顶部标题区域 -->
    <div class="header-section">
      <div class="header-content">
        <div class="title-section">
          <h1 class="main-title">
            <i class="pi pi-folder-open mr-2"></i>
            项目管理
          </h1>
        </div>

        <div class="action-buttons">
          <div class="search-container">
            <div class="search-input-wrapper">
              <i class="pi pi-search search-icon"></i>
              <InputText 
                v-model="searchQuery"
                placeholder="搜索项目名称..."
                class="search-input"
                @input="onSearchInput"
              />
            </div>
          </div>
          <Button 
            label="打开项目" 
            icon="pi pi-folder-open" 
            class="p-button-outlined p-button-lg open-project-btn"
            @click="onOpenProject" 
          />
          <Button 
            label="新建项目" 
            icon="pi pi-plus" 
            class="p-button-primary p-button-lg"
            @click="onCreateProject" 
          />
        </div>
      </div>
    </div>

    <!-- 项目列表区域 -->
    <div class="content-section">
      <div class="projects-container">
        <!-- 空状态 -->
        <div v-if="filteredProjects.length === 0" class="empty-state-card">
          <Card class="empty-card">
            <template #content>
              <div class="empty-content">
                <div class="empty-icon-wrapper">
                  <i class="pi pi-folder-open empty-icon"></i>
                </div>
                <h3 class="empty-title">
                  {{ searchQuery ? '未找到匹配的项目' : '暂无项目' }}
                </h3>
                <p class="empty-description">
                  {{ searchQuery ? 
                    `没有找到包含 "${searchQuery}" 的项目，请尝试其他关键词` : 
                    '您还没有任何项目，点击上方按钮创建或打开一个项目开始使用' 
                  }}
                </p>
                <div class="empty-actions" v-if="!searchQuery">
                  <Button 
                    label="创建第一个项目" 
                    icon="pi pi-plus" 
                    class="p-button-primary"
                    @click="onCreateProject"
                  />
                </div>
              </div>
            </template>
          </Card>
        </div>

        <!-- 项目列表 -->
        <div v-else class="projects-list">
          <Card class="projects-card">
            <template #content>
              <div class="projects-list-content">
                <div 
                  v-for="project in filteredProjects" 
                  :key="project.id"
                  class="project-item"
                  :class="{ 'selected': selectedProject?.id === project.id }"
                  @click="selectProject(project)"
                >
                  <!-- 项目信息 -->
                  <div class="project-info">
                    <h3 class="project-name">{{ project.name }}</h3>
                    <p class="project-path" :title="project.path">
                      <i class="pi pi-folder mr-1"></i>
                      {{ project.path }}
                    </p>
                    <div class="project-meta">
                      <span class="created-date">
                        <i class="pi pi-calendar mr-1"></i>
                        创建于 {{ formatDate(project.created_at) }}
                      </span>
                    </div>
                  </div>

                  <!-- 删除按钮 -->
                  <div class="project-actions">
                    <Button 
                      icon="pi pi-trash" 
                      class="p-button-text p-button-danger p-button-sm delete-btn"
                      @click.stop="removeProject(project.id)"
                      v-tooltip="'删除项目'"
                    />
                  </div>
                </div>
              </div>
            </template>
          </Card>
        </div>
      </div>
    </div>
  </div>
  
  <!-- 确认对话框 -->
  <ConfirmDialog />
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useConfirm } from 'primevue/useconfirm'
import { useThemeStore } from '@/stores/themeStore'
import { usePageStore } from '@/stores/pageStore'
import { useDatabaseStore } from '@/stores/databaseStore'
import { GetAllProjects, DeleteProject } from '@/../wailsjs/go/api/SQLiteAPI'
import { WindowSetSize, WindowSetLightTheme, WindowSetDarkTheme } from '@/../wailsjs/runtime/runtime'
import { InitCache } from '@/../wailsjs/go/api/MetadatasAPI'
import { PathExists, CreateDirectory, OpenDirectory } from '@/../wailsjs/go/api/SystemServiceAPI'
import { Init as InitAppCache } from '@/../wailsjs/go/api/AppCacheApi'
import { flowDataStore } from '@/stores/flowDataStore'

interface Project {
  id: number
  name: string
  path: string
  created_at: any
}

const router = useRouter()
const confirm = useConfirm()
const projects = ref<Project[]>([])
const selectedProject = ref<Project | null>(null)
const themeStore = useThemeStore()
const pageStore = usePageStore()
const databaseStore = useDatabaseStore()
const searchQuery = ref('')
const hasWailsRuntime = typeof window !== 'undefined' && !!(window as any).runtime

// 计算过滤后的项目列表
const filteredProjects = computed(() => {
  if (!searchQuery.value.trim()) {
    return projects.value
  }
  
  const query = searchQuery.value.toLowerCase().trim()
  return projects.value.filter(project => 
    project.name.toLowerCase().includes(query)
  )
})

// 监听主题变化，同步更新标题栏
watch(() => themeStore.isDark, async (isDark) => {
  if (hasWailsRuntime) {
    if (isDark) {
      await WindowSetDarkTheme()
    } else {
      await WindowSetLightTheme()
    }
  }
})

// 页面加载时获取项目列表并设置小窗口
onMounted(async () => {
  // 初始化主题（需等待）
  await themeStore.initTheme()
  // 初始化页面状态（需等待，以便拿到 activeTab/openTabs 等）
  await pageStore.initializeState()
  
  // 检查是否在Wails环境中运行，然后设置小窗口尺寸
  if (hasWailsRuntime) {
    WindowSetSize(800, 600)
  }
  
  await loadProjectList()
  
  // 设置标题栏主题（基于已初始化的主题状态）
  if (hasWailsRuntime) {
    if (themeStore.isDark) {
      await WindowSetDarkTheme()
    } else {
      await WindowSetLightTheme()
    }
  }
})

// 搜索输入处理
const onSearchInput = () => {
  // 搜索逻辑已通过计算属性实现，这里可以添加额外的处理逻辑
  console.log('搜索关键词:', searchQuery.value)
}

// 加载项目列表
const loadProjectList = async () => {
  try {
    const result = await GetAllProjects();
    console.log('Raw project data from backend:', result);
    
    if (result && Array.isArray(result)) {
      result.forEach((project, index) => {
        console.log(`Project ${index}:`, project);
        console.log(`  - created_at type:`, typeof project.created_at);
        console.log(`  - created_at value:`, project.created_at);
        console.log(`  - created_at JSON:`, JSON.stringify(project.created_at));
      });
    }
    
    projects.value = result || [];
  } catch (error) {
    console.error('Error loading projects:', error);
    projects.value = [];
  }
};

// 格式化日期
const formatDate = (dateValue: any) => {
  try {
    let date: Date;
    
    // 处理不同的日期格式
    if (typeof dateValue === 'string') {
      // 如果是字符串，尝试解析
      date = new Date(dateValue);
    } else if (dateValue instanceof Date) {
      // 如果已经是Date对象
      date = dateValue;
    } else if (dateValue && typeof dateValue === 'object') {
      // 如果是Go time.Time对象，可能包含时间戳或其他格式
      if (dateValue.seconds) {
        // Unix时间戳（秒）
        date = new Date(dateValue.seconds * 1000);
      } else if (dateValue.nanos) {
        // 纳秒时间戳
        date = new Date(dateValue.nanos / 1000000);
      } else {
        // 尝试直接转换
        date = new Date(dateValue);
      }
    } else {
      return '未知';
    }
    
    // 检查日期是否有效
    if (isNaN(date.getTime())) {
      return '未知';
    }
    
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  } catch (error) {
    console.error('日期格式化错误:', error, '原始值:', dateValue);
    return '未知';
  }
}

const onOpenProject = async () => {
  try {
    const dir = await OpenDirectory()
    if (!dir) return
    await InitCache(dir)
    await InitAppCache(dir)
    flowDataStore.reset()
    await pageStore.initializeState()
    await databaseStore.refreshDatabases()
    await router.push('/app')
    if (hasWailsRuntime) {
      WindowSetSize(1280, 720)
    }
  } catch (error) {
    console.error('Error opening project:', error)
  }
}

const onCreateProject = async () => {
  try {
    // 直接跳转到新建项目页面
    router.push('/new-project')
  } catch (error: any) {
     console.error('跳转到新建项目页面失败:', error)
   }
}

// 路由跳转：优先跳到已选中的页面标签，否则进入主页面
const navigateToActiveOrApp = async (projectId: number) => {
  const active = pageStore.activeTab
  if (active) {
    await router.push({ name: 'flowPage', params: { pageId: active } })
  } else {
    await router.push({ name: 'app', query: { projectId: projectId.toString() } })
  }
  if (hasWailsRuntime) {
    WindowSetSize(1280, 720)
  }
}

// 处理路径存在的项目：初始化并跳转
const handleProjectWithExistingPath = async (project: Project) => {
  await InitCache(project.path)
  await InitAppCache(project.path)
  flowDataStore.reset()
  await pageStore.initializeState()
  await databaseStore.refreshDatabases()
  await navigateToActiveOrApp(project.id)
  // 清除选中状态，避免视觉残留
  selectedProject.value = null
}

// 处理路径不存在的项目：确认后创建目录、初始化并跳转
const handleProjectWithMissingPath = async (project: Project) => {
  confirm.require({
    message: `当前项目路径不存在：${project.path}\n是否要重新创建并初始化？`,
    header: '项目路径不存在',
    icon: 'pi pi-exclamation-triangle',
    acceptLabel: '重新创建',
    rejectLabel: '取消',
    accept: async () => {
      try {
        await CreateDirectory(project.path)
        await InitCache(project.path)
        await InitAppCache(project.path)
        flowDataStore.reset()
        await pageStore.initializeState()
        await databaseStore.refreshDatabases()
        await navigateToActiveOrApp(project.id)
        // 清除选中状态，避免视觉残留
        selectedProject.value = null
      } catch (err) {
        console.error('创建并初始化项目失败:', err)
      }
    },
    reject: () => {
      console.log('用户取消重新创建项目目录')
      // 清除选中状态，恢复正常背景
      selectedProject.value = null
    }
  })
}

const selectProject = async (project: Project) => {
  try {
    selectedProject.value = project
    console.log('选择项目:', project)
    
    // 先检查项目路径是否存在
    const exists = await PathExists(project.path)
    if (!exists) {
      await handleProjectWithMissingPath(project)
      return
    }

    await handleProjectWithExistingPath(project)
  } catch (error) {
    console.error('Error selecting project:', error)
  }
}

const removeProject = async (projectId: number) => {
  try {
    confirm.require({
      message: '确定要删除这个项目吗？此操作不可撤销。',
      header: '删除确认',
      icon: 'pi pi-exclamation-triangle',
      acceptClass: 'p-button-danger',
      acceptLabel: '删除',
      rejectLabel: '取消',
      accept: async () => {
        try {
          await DeleteProject(projectId)
          await loadProjectList()
          // 如果删除的是当前选中的项目，清空选择
          if (selectedProject.value?.id === projectId) {
            selectedProject.value = null
          }
        } catch (error) {
          console.error('Error removing project:', error)
        }
      }
    })
  } catch (error) {
    console.error('Error removing project:', error)
  }
}
</script>

<style scoped>
.project-select-container {
  min-height: 100vh;
  background: var(--project-select-bg);
  display: flex;
  flex-direction: column;
}

/* 顶部标题区域 */
.header-section {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  background: var(--header-bg);
  backdrop-filter: blur(15px);
  border-bottom: 1px solid var(--surface-border);
  padding: 1.5rem 0;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.08);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  flex-wrap: wrap;
  gap: 2rem;
}

.title-section {
  flex: 1;
  text-align: left;
}

.main-title {
  font-size: 2.2rem;
  font-weight: 600;
  color: var(--text-color);
  margin: 0;
  display: flex;
  align-items: center;
  letter-spacing: -0.02em;
}

.main-title i {
  font-size: 1.6rem;
  color: var(--body-text-color);
  opacity: 0.8;
  background: var(--card-bg);
  padding: 0.5rem;
  border-radius: 8px;
  margin-right: 1rem;
  box-shadow: 0 4px 12px var(--box-shadow-color);
  transition: all 0.3s ease;
  border: 1px solid var(--surface-border);
}

.main-title i:hover {
  transform: translateY(-1px);
  box-shadow: 0 3px 8px rgba(0, 0, 0, 0.1);
  opacity: 1;
  color: var(--body-text-color);
}

/* 打开项目按钮样式优化 */
.open-project-btn {
  background: var(--card-bg) !important;
  color: var(--text-color) !important;
  border: 2px solid var(--border-color) !important;
  font-weight: 600 !important;
  text-shadow: none !important;
}

.open-project-btn:hover {
  background: var(--surface-hover-light) !important;
  color: var(--text-color) !important;
  border-color: var(--surface-hover-light) !important;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px var(--box-shadow-color);
}

.action-buttons {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  align-items: center;
}

.search-container {
  flex: 1;
  min-width: 250px;
  max-width: 400px;
}

.search-input-wrapper {
  position: relative;
  width: 100%;
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #6b7280;
  z-index: 1;
  font-weight: 600;
}

.search-input {
  width: 100%;
  padding-left: 40px;
}

/* 内容区域 */
.content-section {
  flex: 1;
  padding: 3rem 0;
  margin-top: 90px; /* 调整顶部边距适应新的header高度 */
}

.projects-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

/* 空状态 */
.empty-state-card {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.empty-card {
  max-width: 500px;
  width: 100%;
  box-shadow: var(--card-shadow);
  border-radius: 16px;
  border: none;
  background: var(--card-bg);
}

.empty-content {
  text-align: center;
  padding: 3rem 2rem;
}

.empty-icon-wrapper {
  margin-bottom: 2rem;
}

.empty-icon {
  font-size: 4rem;
  color: var(--body-text-color);
  opacity: 0.5;
}

.empty-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-color);
  margin: 0 0 1rem 0;
}

.empty-description {
  color: var(--body-text-color);
  opacity: 0.7;
  line-height: 1.6;
  margin: 0 0 2rem 0;
}

.empty-actions {
  display: flex;
  justify-content: center;
}

/* 项目列表 */
.projects-list {
  width: 100%;
}

.projects-card {
  border-radius: 16px;
  border: none;
  box-shadow: var(--card-shadow);
  background: var(--card-bg);
  backdrop-filter: blur(10px);
}

.projects-list-content {
  padding: 0;
}

.project-item {
  display: flex;
  align-items: center;
  padding: 1.5rem 2rem;
  border-bottom: 1px solid var(--surface-border);
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.project-item:last-child {
  border-bottom: none;
}

.project-item:hover {
  background: var(--surface-hover-subtle);
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.project-item.selected {
  background: var(--project-selected-bg);
  border-left: 3px solid var(--project-selected-color);
  border-radius: 12px;
}

.project-info {
  flex: 1;
  min-width: 0;
  text-align: left;
}

.project-name {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-color);
  margin: 0 0 0.5rem 0;
  line-height: 1.3;
  text-align: left;
}

.project-path {
  color: var(--text-color-tertiary);
  font-size: 0.85rem;
  margin: 0 0 0.5rem 0;
  display: flex;
  align-items: center;
  word-break: break-all;
  line-height: 1.4;
  text-align: left;
}

.project-meta {
  display: flex;
  align-items: center;
  gap: 1rem;
  text-align: left;
}

.created-date {
  color: var(--text-color-tertiary);
  font-size: 0.8rem;
  display: flex;
  align-items: center;
}

.project-actions {
  opacity: 0;
  transition: opacity 0.2s ease;
  margin-left: 1rem;
}

.project-item:hover .project-actions {
  opacity: 1;
}

.delete-btn {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: rgba(239, 68, 68, 0.1);
}

.delete-btn:hover {
  background: rgba(239, 68, 68, 0.2);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    text-align: center;
    gap: 1.5rem;
    align-items: center;
  }

  .title-section {
    text-align: center;
  }

  .main-title {
    font-size: 1.8rem;
    justify-content: center;
  }

  .main-title i {
    font-size: 1.5rem;
    padding: 0.5rem;
    margin-right: 0.8rem;
  }

  .content-section {
    padding: 2rem 0;
    margin-top: 110px; /* 调整移动端的顶部边距 */
  }

  .projects-container {
    padding: 0 1rem;
  }

  .project-item {
    padding: 1.25rem 1.5rem;
  }

  .search-container {
    width: 100%;
    max-width: none;
    order: -1; /* 将搜索框放在按钮上方 */
  }
}

@media (max-width: 480px) {
  .action-buttons {
    width: 100%;
    justify-content: center;
  }

  .main-title {
    font-size: 1.6rem;
  }

  .project-item {
    padding: 1rem;
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .project-actions {
    position: absolute;
    top: 1rem;
    right: 1rem;
    margin-left: 0;
  }
}
</style>

<!-- 全局样式，定义主题变量 -->
<style>
:root {
  --project-select-bg: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  --header-bg: var(--card-bg);
  --card-shadow: 0 8px 25px var(--box-shadow-color);
  --project-selected-color: rgba(0, 120, 212, 0.25);
  --project-selected-bg: rgba(0, 120, 212, 0.08);
  --text-color-tertiary: rgba(107, 114, 128, 0.8); 
  --surface-hover-subtle: rgba(0, 120, 212, 0.05); 
  --logo-color: var(--body-text-color); 
  --logo-bg: var(--surface-hover); 
}

:root.p-dark {
  --project-select-bg: linear-gradient(135deg, #2d3748 0%, #4a5568 100%);
  --header-bg: var(--card-bg);
  --card-shadow: 0 8px 25px var(--box-shadow-color);
  --project-selected-color: rgba(0, 120, 212, 0.35);
  --project-selected-bg: rgba(0, 120, 212, 0.12);
  --text-color-tertiary: rgba(156, 163, 175, 0.9); 
  --surface-hover-subtle: rgba(0, 120, 212, 0.08); 
  --logo-color: var(--body-text-color); 
  --logo-bg: var(--surface-hover); 
}
</style>