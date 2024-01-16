import { ref, type Ref } from 'vue'
import { useWindowSize } from '@vueuse/core'

const {
  width,
  height
} = useWindowSize()

interface WindowStyle {
  width: number
  height: number
  x: number
  y: number
  zIndex: number
}

interface ProcessInfo {
  pid: number
  appId: string
  title: string
  icon: string
  url: string
  hidden: boolean
  focus: boolean
  windowStyle: WindowStyle
}

interface ForkAppParam {
  title: string
  icon: string
  url: string
  appId: string
  windowStyle: {
    width: number
    height: number
  }
}

let pid = 1

const DefaultZIndex = 1000

const focus = (pid: number, processList: Ref<ProcessInfo[]>) => {
  processList.value.forEach(p => {
    if (p.pid === pid) {
      p.windowStyle.zIndex = DefaultZIndex + 1
      p.focus = true
    } else {
      p.windowStyle.zIndex = DefaultZIndex
      p.focus = false
    }
  })
}

const close = (pid: number, processList: Ref<ProcessInfo[]>) => {
  processList.value.splice(processList.value.findIndex(p => p.pid === pid), 1)
}

const min = (pid: number, processList: Ref<ProcessInfo[]>) => {
  processList.value.find(p => p.pid === pid).hidden = true
}

const show = (pid: number, processList: Ref<ProcessInfo[]>) => {
  processList.value.find(p => p.pid === pid).hidden = false
}

const toggle = (pid: number, processList: Ref<ProcessInfo[]>) => {
  let item = processList.value.find(p => p.pid === pid)
  const isHidden = item.hidden
  if (isHidden) {
    item.hidden = false
    focus(pid, processList)
  } else {
    if (item.windowStyle.zIndex === DefaultZIndex + 1) {
      item.hidden = true
    } else {
      focus(pid, processList)
    }
  }
}

export function useDesktop () {
  const processList: Ref<ProcessInfo[]> = ref([])
  const fork = (app: ForkAppParam) => {
    processList.value.forEach(p => {
      p.windowStyle.zIndex = DefaultZIndex
    })

    pid += 1
    const newProcess = {
      pid,
      appId: app.appId,
      title: app.title,
      icon: app.icon,
      url: app.url,
      hidden: false,
      windowStyle: {
        zIndex: DefaultZIndex,
        ...app.windowStyle
      }
    }

    newProcess.windowStyle.x = width.value / 2 - newProcess.windowStyle.width / 2
    newProcess.windowStyle.y = height.value / 2 - newProcess.windowStyle.height / 2 - 60

    processList.value.push(newProcess)

    console.log(newProcess)

    focus(pid, processList)
  }

  const action = (pid: number, action: 'focus' | 'close' | 'min' | 'show' | 'toggle') => {
    const actions = {
      focus,
      close,
      min,
      show,
      toggle
    }

    actions[action](pid, processList)

  }

  return {
    processList,
    fork,
    action
  }
}
