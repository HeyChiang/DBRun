import { defineStore } from 'pinia'
import { ref } from 'vue'
import { Get as CacheGet, Set as CacheSet } from '@/../wailsjs/go/api/AppCacheApi'

export interface Page {
  key: string;
  label: string;
  icon?: string;
  type: 'page' | 'group';
  expanded?: boolean;
  children?: Page[];
}

export interface Tab {
  key: string;
  label: string;
  active: boolean;
}

export const usePageStore = defineStore('page', () => {
  // State
  const pages = ref<Page[]>([])
  const openTabs = ref<Tab[]>([])
  const activeTab = ref<string | null>(null)

  // 初始化状态（从应用缓存）
  const initializeState = async () => {
    try {
      const result = await CacheGet('pageStore')
      const savedState = typeof result === 'string' ? result : null
      console.log('[PageStore] Loading saved state from AppCache:', savedState)

      if (savedState) {
        try {
          const parsed = JSON.parse(savedState)
          console.log('[PageStore] Parsed state:', parsed)
          const migratePages = (pagesData: any[]): Page[] => {
            return pagesData.map(page => ({
              key: page.key,
              label: page.label,
              icon: page.icon,
              type: page.type || (page.children ? 'group' : 'page'),
              expanded: page.expanded,
              children: page.children ? migratePages(page.children) : undefined
            }))
          }
          pages.value = migratePages(parsed.pages || [])
          openTabs.value = parsed.openTabs || []
          activeTab.value = parsed.activeTab || null
        } catch (e) {
          console.warn('Failed to parse saved state:', e)
          pages.value = []
          openTabs.value = []
          activeTab.value = null
        }
      } else {
        pages.value = []
        openTabs.value = []
        activeTab.value = null
      }
    } catch (e) {
      console.warn('Failed to load state from AppCache:', e)
      pages.value = []
      openTabs.value = []
      activeTab.value = null
    }
  }

  // Actions
  async function savePageTab() {
    try {
      // Save page structure
      const payload = JSON.stringify({
        pages: pages.value,
        openTabs: openTabs.value,
        activeTab: activeTab.value
      })
      console.log('[PageStore] Saving state to AppCache:', payload)
      await CacheSet('pageStore', payload)
    } catch (e) {
      console.warn('Failed to save state:', e)
    }
  }

  function generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
      var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
      return v.toString(16);
    });
  }

  function findPage(key: string): Page | null {
    const search = (pagesArray: Page[]): Page | null => {
      for (const page of pagesArray) {
        if (page.key === key) return page
        if (page.children) {
          const found = search(page.children)
          if (found) return found
        }
      }
      return null
    }
    return search(pages.value)
  }

  // 查找目标项的父级数组与索引（根级的父级为 pages.value，parentKey 为 null）
  function findParentInfo(targetKey: string): { parentArray: Page[] | null, parentKey: string | null, index: number } | null {
    // 根级处理
    const rootIndex = pages.value.findIndex(p => p.key === targetKey)
    if (rootIndex !== -1) {
      return { parentArray: pages.value, parentKey: null, index: rootIndex }
    }

    // 递归查找子级
    const stack: { parent: Page; children: Page[] }[] = []
    const pushChildren = (parent: Page) => {
      if (parent.children && parent.children.length) {
        stack.push({ parent, children: parent.children })
        parent.children.forEach(child => pushChildren(child))
      }
    }
    pages.value.forEach(p => pushChildren(p))

    // 遍历所有可能的 children 引用，找到包含 targetKey 的父数组与索引
    for (const { parent, children } of stack) {
      const idx = children.findIndex(c => c.key === targetKey)
      if (idx !== -1) {
        return { parentArray: parent.children || null, parentKey: parent.key, index: idx }
      }
    }
    return null
  }

  // 判断 descendantKey 是否为 ancestorKey 的后代
  function isDescendant(ancestorKey: string, descendantKey: string): boolean {
    const ancestor = findPage(ancestorKey)
    if (!ancestor) return false
    const dfs = (node: Page | undefined): boolean => {
      if (!node || !node.children) return false
      for (const child of node.children) {
        if (child.key === descendantKey) return true
        if (dfs(child)) return true
      }
      return false
    }
    return dfs(ancestor)
  }

  // 移动：支持同级前/后插入，以及放入分组（inside）
  function moveItem(sourceKey: string, targetKey: string, position: 'before' | 'after' | 'inside' = 'before') {
    if (sourceKey === targetKey) return

    const sourceInfo = findParentInfo(sourceKey)
    const targetInfo = findParentInfo(targetKey)
    const sourceItem = findPage(sourceKey)
    const targetItem = findPage(targetKey)

    if (!sourceInfo || !targetInfo || !sourceInfo.parentArray || !sourceItem || !targetItem) {
      console.warn('[PageStore] moveItem failed: parent info or items not found')
      return
    }

    // 防止将某个节点/分组移动到其后代中（形成循环）
    if (position === 'inside' && isDescendant(sourceKey, targetKey)) {
      console.warn('[PageStore] moveItem ignored: cannot move a node into its own descendant')
      return
    }

    // 先从原父级移除
    const sourceParentArray = sourceInfo.parentArray
    const sourceIndex = sourceInfo.index
    const [moved] = sourceParentArray.splice(sourceIndex, 1)

    if (position === 'inside') {
      if (targetItem.type !== 'group') {
        console.warn('[PageStore] moveItem ignored: inside only valid for group target')
        // 放回原位置
        sourceParentArray.splice(sourceIndex, 0, moved)
        return
      }
      if (!targetItem.children) targetItem.children = []
      targetItem.children.push(moved)
      savePageTab()
      return
    }

    // before/after：插入到目标所在父级的兄弟位置（可跨父级）
    const targetParentArray = targetInfo.parentArray
    if (!targetParentArray) {
      console.warn('[PageStore] moveItem failed: target parent array missing')
      // 放回原位置
      sourceParentArray.splice(sourceIndex, 0, moved)
      return
    }

    let originalTargetIndex = targetInfo.index
    // 如果源与目标同父级，删除源后目标索引可能变化
    if (sourceParentArray === targetParentArray && sourceIndex < originalTargetIndex) {
      originalTargetIndex -= 1
    }

    let insertIndex = position === 'before' ? originalTargetIndex : originalTargetIndex + 1
    if (insertIndex < 0) insertIndex = 0
    if (insertIndex > targetParentArray.length) insertIndex = targetParentArray.length

    targetParentArray.splice(insertIndex, 0, moved)
    savePageTab()
  }

  function addItem(newItem: { 
    label: string; 
    type: 'page' | 'group';
    icon?: string; 
    parentKey?: string 
  }): Page {
    const item: Page = {
      key: generateUUID(),
      label: newItem.label,
      icon: newItem.icon || (newItem.type === 'page' ? 'pi pi-file' : 'pi pi-folder'),
      type: newItem.type,
    }

    if (newItem.type === 'group') {
      item.expanded = true
      item.children = []
    }

    if (newItem.parentKey) {
      const parent = findPage(newItem.parentKey)
      if (parent && parent.type === 'group') {
        if (!parent.children) parent.children = []
        parent.children.push(item)
      }
    } else {
      pages.value.push(item)
    }

    savePageTab()
    return item
  }

  //  更新页面标签的名字
  function updatePage(key: string, updates: Partial<Page>) {
    const page = findPage(key)
    if (page) {
      Object.assign(page, updates)
      // 如果更新了标签，同时更新打开的标签页
      if (updates.label) {
        const idx = openTabs.value.findIndex(t => t.key === key)
        if (idx !== -1) {
          const current = openTabs.value[idx]
          openTabs.value[idx] = { ...current, label: updates.label }
        }
      }
      savePageTab()
    }
  }

  function deletePage(key: string) {
    console.log('[PageStore] Deleting page with key:', key);
    console.log('[PageStore] Pages before deletion:', JSON.parse(JSON.stringify(pages.value)));
    
    const deleteRecursive = (pages: Page[], targetKey: string): Page[] => {
      return pages.filter(page => {
        if (page.key === targetKey) {
          // Delete all children if it's a group
          if (page.type === 'group' && page.children) {
            page.children.forEach(child => {
              deleteRecursive([child], child.key);
            });
          }
          return false;
        }
        if (page.children) {
          page.children = deleteRecursive(page.children, targetKey);
        }
        return true;
      });
    };

    pages.value = deleteRecursive(pages.value, key);
    console.log('[PageStore] Pages after deletion:', JSON.parse(JSON.stringify(pages.value)));
    
    // Cleanup tabs，里面同时有保存数据
    closeTab(key)
    
    // Explicitly save state after deletion
    savePageTab();
  }

  function toggleGroup(key: string) {
    const group = findPage(key)
    if (group && group.type === 'group') {
      group.expanded = !group.expanded
      savePageTab()
    }
  }

  // Tab management
  const openTab = (pageKey: string) => {
    console.log('[PageStore] Opening tab:', pageKey)
    const page = findPage(pageKey)
    if (!page || page.type !== 'page') {
      console.warn('[PageStore] Cannot open tab: invalid page or type')
      return
    }

    const existingTab = openTabs.value.find(tab => tab.key === pageKey)
    if (!existingTab) {
      console.log('[PageStore] Creating new tab')
      openTabs.value.push({
        key: pageKey,
        label: page.label,
        active: true
      })
    }

    // Deactivate all tabs
    openTabs.value.forEach(tab => {
      if (tab.key === pageKey) {
        console.log('[PageStore] Activating tab:', pageKey)
        tab.active = true
      } else {
        tab.active = false
      }
    })

    activeTab.value = pageKey
    console.log('[PageStore] Active tab set to:', activeTab.value)
    savePageTab()
  }

  function closeTab(pageKey: string) {
    const index = openTabs.value.findIndex(tab => tab.key === pageKey)
    if (index !== -1) {
      openTabs.value.splice(index, 1)
      
      if (activeTab.value === pageKey && openTabs.value.length > 0) {
        let newActiveTab;
        // 如果有下一个tab就打开下一个，否则打开前一个
        if (index < openTabs.value.length) {
          newActiveTab = openTabs.value[index];
        } else if (index > 0) {
          newActiveTab = openTabs.value[index - 1];
        } else {
          newActiveTab = openTabs.value[0];
        }
        newActiveTab.active = true
        activeTab.value = newActiveTab.key
      } else if (openTabs.value.length === 0) {
        activeTab.value = null
      }
      
      savePageTab()
    }
  }

  
  
  return {
    pages,
    openTabs,
    activeTab,
    // 暴露初始化方法，便于视图在启动时显式等待
    initializeState,
    findPage,
    // 暴露移动函数用于拖拽排序
    moveItem,
    addItem,
    updatePage,
    deletePage,
    toggleGroup,
    openTab,
    closeTab
  }
})
