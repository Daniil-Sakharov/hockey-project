import { useState, useCallback } from 'react'
import { useMediaQuery, breakpoints } from '@/shared/hooks'

interface SidebarState {
  isOpen: boolean
  isCollapsed: boolean
  toggle: () => void
  open: () => void
  close: () => void
  collapse: () => void
  expand: () => void
}

export function useSidebar(): SidebarState {
  const isDesktop = useMediaQuery(breakpoints.lg)
  const [isOpen, setIsOpen] = useState(false)
  const [isCollapsed, setIsCollapsed] = useState(!isDesktop)

  const toggle = useCallback(() => {
    if (isDesktop) {
      setIsCollapsed((prev) => !prev)
    } else {
      setIsOpen((prev) => !prev)
    }
  }, [isDesktop])

  const open = useCallback(() => setIsOpen(true), [])
  const close = useCallback(() => setIsOpen(false), [])
  const collapse = useCallback(() => setIsCollapsed(true), [])
  const expand = useCallback(() => setIsCollapsed(false), [])

  return {
    isOpen: isDesktop || isOpen,
    isCollapsed: !isDesktop ? false : isCollapsed,
    toggle,
    open,
    close,
    collapse,
    expand,
  }
}
