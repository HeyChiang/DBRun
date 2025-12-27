import { ref, readonly } from 'vue';

// 定义表更新事件接口
export interface TableUpdateEvent {
  dbId: string | number;
  dbName: string;
  schemaName: string;
  tableName: string;
  tableInfo: any;
}

// 简单的事件总线
export const createEventBus = () => {
  // 存储所有事件处理函数
  const listeners: Record<string, Function[]> = {};
  
  // 订阅事件
  const on = (event: string, callback: Function) => {
    if (!listeners[event]) {
      listeners[event] = [];
    }
    listeners[event].push(callback);
  };
  
  // 取消订阅
  const off = (event: string, callback: Function) => {
    if (!listeners[event]) return;
    listeners[event] = listeners[event].filter(
      listener => listener !== callback
    );
  };
  
  // 触发事件
  const emit = (event: string, data?: any) => {
    if (!listeners[event]) return;
    listeners[event].forEach(callback => callback(data));
  };
  
  return {
    on,
    off,
    emit
  };
};

// 创建一个全局事件总线实例
export const eventBus = createEventBus();
