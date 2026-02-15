// import { defineConfig } from 'vitepress'
import { withMermaid } from "vitepress-plugin-mermaid";

// https://vitepress.dev/reference/site-config
export default withMermaid({
  title: "VVORKER DOCS",
  description: "VVORKER DOCS",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [{ text: "Home", link: "/" }],

    sidebar: [
      {
        text: "DESIGN",
        items: [
          { text: "基本概念", link: "/design/basic" },
          { text: "快速开始", link: "/design/quickstart" },
          { text: "网络", link: "/design/network" },
        ],
      },
      {
        text: "SDK",
        items: [
          { text: "基础", link: "/sdk/vvorker-sdk" },
          { text: "最佳实践", link: "/sdk/best-practice" },
          { text: "测试", link: "/sdk/test" },
          { text: "vvbind - KV", link: "/sdk/vvbind-kv" },
          { text: "vvbind - OSS", link: "/sdk/vvbind-oss" },
          { text: "vvbind - MySQL", link: "/sdk/vvbind-mysql" },
          { text: "vvbind - PostgreSQL", link: "/sdk/vvbind-pgsql" },
          { text: "vvbind - Service", link: "/sdk/vvbind-service" },
          { text: "vvbind - Vars", link: "/sdk/vvbind-vars" },
          { text: "vvbind - Task", link: "/sdk/vvbind-task" },
          { text: "vvbind - Assets", link: "/sdk/vvbind-assets" },
          { text: "vvbind - Proxy", link: "/sdk/vvbind-proxy" },
        ],
      },
      {
        text: "CLI",
        items: [{ text: "CLI", link: "/cli/vvorker-cli" }],
      },
      {
        text: "CONFIG",
        items: [
          { text: "CLI 环境", link: "/config/env" },
          { text: "项目配置", link: "/config/vvorker_project_config" },
          { text: "节点配置", link: "/config/node_config" },
          { text: "SSO配置", link: "/config/sso" },
        ],
      },
    ],

    socialLinks: [
      { icon: "github", link: "https://github.com/dicarne/vvorker" },
    ],
  },
});
