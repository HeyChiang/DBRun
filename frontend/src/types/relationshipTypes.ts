// 关系类型枚举定义
export enum RelationshipType {
  ONE = 'one',
  MANY = 'many',
  NONE = 'none'
}

// 表间关系接口定义
export interface RelationshipData {
  source: RelationshipType;
  target: RelationshipType;
}
