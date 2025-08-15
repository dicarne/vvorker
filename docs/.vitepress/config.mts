// import { defineConfig } from 'vitepress'
import { withMermaid } from 'vitepress-plugin-mermaid'

// https://vitepress.dev/reference/site-config
export default withMermaid({
  title: "VVORKER DOCS",
  description: "VVORKER DOCS",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' }
    ],

    sidebar: [
      {
        text: 'SDK',
        items: [
          { text: '基础', link: '/sdk/vvorker-sdk' },
          { text: '最佳实践', link: '/sdk/best-practice' },
        ]
      },
      {
        text: "CLI",
        items: [
          { text: 'CLI', link: '/cli/vvorker-cli' },
        ]
      },
      {
        text: "CONFIG",
        items: [
          { text: '环境变量', link: '/config/env' },
          { text: '节点配置', link: '/config/node_config' },
        ]
      },
      {
        text: "DESIGN",
        items: [
          { text: '钦定', link: '/design/must' },
          { text: '网络', link: '/design/network' },
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/dicarne/vvorker' }
    ]
  }
})
