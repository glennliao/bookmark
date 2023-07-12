export function toCateTree (list: any[]): any[] {
  // 首先按级别进行排序（不需要请删除）
  list.sort(function (a, b) {
    return a.sorts - b.sorts
  })

  const map: Record<string, any[]> = {}

  list.forEach(item => {
    map[item.parentId] = map[item.parentId] || []
    map[item.parentId].push(item)
  })

  list.forEach(item => {
    if (map[item.cateId] && map[item.cateId].length > 0) {
      item.children = map[item.cateId]
    }
  })

  return map.root
}

export function treeEach (tree: any[], call: (any) => any): void {
  tree = tree || []
  tree.forEach(item => {

    if (item.children) {
      treeEach(item.children, call)
    }

    call(item)
  })
}
