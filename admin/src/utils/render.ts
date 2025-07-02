import { h, type Component } from 'vue'
import { RouterLink } from 'vue-router'
import { NIcon } from 'naive-ui'
export function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

export function renderMenuRouterLink(label: string, routerName: string) {
  return () =>
    h(
      RouterLink,
      {
        to: {
          name: routerName,
        },
      },
      { default: () => label },
    )
}
