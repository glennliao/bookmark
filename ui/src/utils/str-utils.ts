export const camelCaseToLine = (v: string): string => {
  return v.replace(/([A-Z])/g, '-$1').toLowerCase()
}
