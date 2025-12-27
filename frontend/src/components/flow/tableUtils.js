export const getKeyTypeText = (key) => {
  switch (key) {
    case 'PRI':
      return 'PK'
    case 'UNI':
      return 'UK'
    case 'FOR':
      return 'FK'
    case 'IDX':
      return 'IDX'
    default:
      return ''
  }
}

export const getKeyTypeClass = (key) => {
  switch (key) {
    case 'PRI':
      return 'primary-key'
    case 'FOR':
      return 'foreign-key'
    case 'UNI':
      return 'unique-key'
    case 'IDX':
      return 'index-key'
    default:
      return ''
  }
}
