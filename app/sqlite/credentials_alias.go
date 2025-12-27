package sqlite

import meta "dbrun/app/sqlite/metadata"

// Credentials 是 metadata.Credentials 的类型别名，用于保持对前端的兼容
// 这样 API 仍可使用 sqlite.Credentials，而底层通过 MetadataManager（GORM）进行操作。
type Credentials = meta.Credentials