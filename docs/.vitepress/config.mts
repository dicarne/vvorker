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
          { text: 'JS SDK', link: '/sdk/js/vvorker-sdk' },
        ]
      },
      {
        text: "CLI",
        items: [
          { text: 'CLI', link: '/cli/vvorker-cli' },
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
