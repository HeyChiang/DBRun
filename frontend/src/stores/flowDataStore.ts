export interface FlowData {
  nodes: any[];
  edges: any[];
}

import { Get as CacheGet, Set as CacheSet } from '@/../wailsjs/go/api/AppCacheApi'

class FlowDataManager {
  private nodeDataMap = new Map<string, any[]>();
  private edgeDataMap = new Map<string, any[]>();

  // 保存节点数据
  async saveNodeData(pageId: string, nodes: any[]) {
    // 保存当前节点状态
    this.nodeDataMap.set(pageId, nodes);
    try {
      await CacheSet(`flow_node_${pageId}`, JSON.stringify(nodes))
    } catch (e) {
      console.warn('保存节点数据到AppCache失败:', e)
    }
  }

  // 保存边缘数据
  async saveEdgeData(pageId: string, edges: any[]) {
    // 保存当前边缘状态
    this.edgeDataMap.set(pageId, edges);
    try {
      await CacheSet(`flow_edge_${pageId}`, JSON.stringify(edges))
    } catch (e) {
      console.warn('保存边缘数据到AppCache失败:', e)
    }
  }

  // 同时保存节点和边缘数据
  async saveFlowData(pageId: string, nodes: any[], edges: any[]) {
    await this.saveNodeData(pageId, nodes);
    await this.saveEdgeData(pageId, edges);
  }

  // 获取节点数据（同步：仅返回内存中的数据）
  getNodeData(pageId: string): any[] {
    return this.nodeDataMap.get(pageId) || [];
  }

  // 异步获取节点数据（从AppCache）
  async getNodeDataAsync(pageId: string): Promise<any[]> {
    if (!this.nodeDataMap.has(pageId)) {
      try {
        const result = await CacheGet(`flow_node_${pageId}`)
        const savedData = typeof result === 'string' ? result : null
        if (savedData) {
          const parsedData = JSON.parse(savedData)
          this.nodeDataMap.set(pageId, parsedData)
        } else {
          this.nodeDataMap.set(pageId, [])
        }
      } catch (e) {
        console.warn('从AppCache读取节点数据失败:', e)
        this.nodeDataMap.set(pageId, [])
      }
    }
    return this.nodeDataMap.get(pageId) || []
  }

  // 获取边缘数据（同步：仅返回内存中的数据）
  getEdgeData(pageId: string): any[] {
    return this.edgeDataMap.get(pageId) || [];
  }

  // 异步获取边缘数据（从AppCache）
  async getEdgeDataAsync(pageId: string): Promise<any[]> {
    if (!this.edgeDataMap.has(pageId)) {
      try {
        const result = await CacheGet(`flow_edge_${pageId}`)
        const savedData = typeof result === 'string' ? result : null
        if (savedData) {
          const parsedData = JSON.parse(savedData)
          this.edgeDataMap.set(pageId, parsedData)
        } else {
          this.edgeDataMap.set(pageId, [])
        }
      } catch (e) {
        console.warn('从AppCache读取边缘数据失败:', e)
        this.edgeDataMap.set(pageId, [])
      }
    }
    return this.edgeDataMap.get(pageId) || []
  }

  // 获取完整的流程图数据（异步）
  async getFlowDataAsync(pageId: string): Promise<FlowData> {
    console.log('flowDataStore.getFlowDataAsync 被调用，pageId:', pageId)
    const nodes = await this.getNodeDataAsync(pageId)
    const edges = await this.getEdgeDataAsync(pageId)
    console.log('获取到的节点数据:', nodes)
    console.log('获取到的边缘数据:', edges)
    return { nodes, edges }
  }

  reset(pageId?: string) {
    if (pageId) {
      this.nodeDataMap.delete(pageId)
      this.edgeDataMap.delete(pageId)
    } else {
      this.nodeDataMap.clear()
      this.edgeDataMap.clear()
    }
  }
}

// 创建单例实例
export const flowDataStore = new FlowDataManager();
