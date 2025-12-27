import { service, models } from '@/../wailsjs/go/models'

// 表节点信息接口定义
export interface TableNodeInfo {
  dbId: string;
  dbName: string;
  schemaName: string;
  table: models.TableInfoVO | models.ViewInfoVO;
  tableId?: number; // 新增：表的稳定ID，优先用于缓存查询与更新
}
