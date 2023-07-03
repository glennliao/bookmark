export const camelCaseToLine = (v: string): string => {
  return v.replace(/([A-Z])/g, '-$1').toLowerCase()
}

function getHostFromUrl (url: string): string {
  let host = ''
  if (url.indexOf('://') > -1) {
    host = url.split('/')[2]
  } else {
    host = url.split('/')[0]
  }
  host = host.split(':')[0]
  host = host.split('?')[0]
  return host
}

export function colorByURL (url: string):string {

  const host = getHostFromUrl(url)

  let hash = 0
  for (let i = 0; i < host.length; i++) {
    hash = host.charCodeAt(i) + ((hash << 5) - hash)
  }
  const c = (hash & 0x00FFFFFF)
    .toString(16)
    .toUpperCase()

  return '#' + '00000'.substring(0, 6 - c.length) + c
}
