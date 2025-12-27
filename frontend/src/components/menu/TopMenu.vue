<script setup lang="ts">
import { ref, onMounted } from 'vue';
import Menubar from 'primevue/menubar';
import { useThemeStore } from '../../stores/themeStore';
import { WindowSetLightTheme, WindowSetDarkTheme } from '@/../wailsjs/runtime/runtime';
import FeedbackDialog from './FeedbackDialog.vue';
import AboutDialog from './AboutDialog.vue';
import { useRouter } from 'vue-router';
import { usePageStore } from '@/stores/pageStore';
import { useDatabaseStore } from '@/stores/databaseStore';
import { OpenDirectory } from '@/../wailsjs/go/api/SystemServiceAPI';
import { Init as InitAppCache } from '@/../wailsjs/go/api/AppCacheApi';
import { InitCache } from '@/../wailsjs/go/api/MetadatasAPI';
import { flowDataStore } from '@/stores/flowDataStore';

const themeStore = useThemeStore();
const router = useRouter();
const pageStore = usePageStore();
const databaseStore = useDatabaseStore();
const feedbackDialog = ref();
const aboutDialog = ref();

const updateMenuItems = () => {
    const viewMenu = items.value.find(item => item.label === '视图');
    if (viewMenu && viewMenu.items) {
        const themeItem = viewMenu.items.find(item => item.label.includes('模式'));
        if (themeItem) {
            themeItem.icon = themeStore.currentTheme.icons.theme;
            themeItem.label = themeStore.currentTheme.icons.themeLabel;
        }
        
        console.log('themeStore.isDark', themeStore.isDark)
        // 设置系统标题栏主题
        themeStore.isDark ? WindowSetDarkTheme() : WindowSetLightTheme();
    }
};

const openProject = async () => {
    const dir = await OpenDirectory();
    if (!dir) return;
    await InitCache(dir);
    await InitAppCache(dir);
    flowDataStore.reset();
    await pageStore.initializeState();
    await databaseStore.refreshDatabases();
    await router.push('/app');
};

onMounted(() => {
    themeStore.initTheme();
    updateMenuItems();
});

const items = ref([
    {
        label: '文件',
        icon: themeStore.currentTheme.icons.file,
        items: [
            {
                label: '新建',
                icon: themeStore.currentTheme.icons.new,
                command: () => {
                    router.push('/new-project');
                }
            },
            {
                label: '打开',
                icon: themeStore.currentTheme.icons.open,
                command: () => {
                    openProject();
                }
            },
            {
                label: '保存',
                icon: themeStore.currentTheme.icons.save,
                command: () => {
                    console.log('点击菜单：保存');
                }
            }
        ]
    },
    {
        label: '视图',
        icon: themeStore.currentTheme.icons.view,
        items: [
            {
                label: themeStore.currentTheme.icons.themeLabel,
                icon: themeStore.currentTheme.icons.theme,
                command: () => {
                    themeStore.toggleTheme();
                    updateMenuItems();
                }
            }
        ]
    },
    {
        label: '帮助',
        icon: themeStore.currentTheme.icons.help,
        items: [
            {
                label: '关于',
                icon: 'pi pi-info-circle',
                command: () => {
                    aboutDialog.value.show();
                }
            },
            {
                label: '反馈建议',
                icon: 'pi pi-comment',
                command: () => {
                    feedbackDialog.value.open();
                }
            }
        ]
    }
]);
</script>

<template>
    <div class="top-menu">
        <div class="menu-container">
            <Menubar 
                :model="items" 
                class="menu-bar" 
                :autoZIndex="true" 
                :dismissable="true"
            />
        </div>
        <FeedbackDialog ref="feedbackDialog" />
        <AboutDialog ref="aboutDialog" />
    </div>
</template>

<style scoped>
.top-menu {
    width: 100%;
    display: flex;
    align-items: center;
    background-color: var(--surface-card);
    border-bottom: 1px solid var(--surface-border);
}

.menu-container {
    flex: 1;
    display: flex;
}

.menu-bar {
    width: 100%;
    border: none;
    border-radius: 0;
}
</style>

<style>
.p-menubar {
    --p-menubar-background: var(--surface-card);
    --p-menubar-border-color: var(--surface-border);
    --p-menubar-color: var(--text-color);
    
    /* 菜单项样式 */
    --p-menubar-item-color: var(--text-color);
    --p-menubar-item-focus-color: var(--primary-color);
    --p-menubar-item-active-color: var(--primary-color);
    --p-menubar-item-focus-background: var(--surface-hover);
    --p-menubar-item-active-background: var(--surface-hover);
    
    /* 子菜单样式 */
    --p-menubar-submenu-background: var(--card-bg);
    --p-menubar-submenu-border-color: var(--surface-border);
    --p-menubar-submenu-padding: 0.5rem;
    
    /* 图标颜色 */
    --p-menubar-item-icon-color: var(--text-color);
    --p-menubar-item-icon-focus-color: var(--primary-color);
    --p-menubar-item-icon-active-color: var(--primary-color);
    
    background-color: var(--p-menubar-background);
    border: none;
}

/* 确保子菜单的样式正确 */
.p-menubar .p-submenu-list {
    background-color: var(--p-menubar-submenu-background);
    border: 1px solid var(--p-menubar-submenu-border-color);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
    backdrop-filter: blur(8px);
}

/* 确保菜单项的悬停效果正确 */
.p-menubar .p-menuitem-link:hover {
    background-color: var(--p-menubar-item-focus-background);
}

/* 确保菜单项文字和图标颜色正确 */
.p-menubar .p-menuitem-text,
.p-menubar .p-menuitem-icon {
    color: var(--p-menubar-item-color);
}

/* 确保子菜单图标颜色正确 */
.p-menubar .p-submenu-icon {
    color: var(--p-menubar-submenu-icon-color, var(--text-color));
}
</style>
