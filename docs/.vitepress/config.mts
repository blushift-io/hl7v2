import { defineConfig } from 'vitepress'

export default defineConfig({
  title: "hl7v2",
  description: "High-Performance Go HL7 v2 Parser, Encoder, Query Engine, & MLLP Toolkit",
  base: "/hl7v2/",
  themeConfig: {
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Examples', link: '/guide/examples' },
      {
        text: 'v0.3.0',
        items: [
          { text: 'v0.3.0 (Latest)', link: '/guide/getting-started' },
          { text: 'Release Notes', link: '/guide/versioning' }
        ]
      }
    ],

    sidebar: [
      {
        text: 'Getting Started',
        items: [
          { text: 'Introduction', link: '/guide/getting-started' }
        ]
      },
      {
        text: 'Core Guides',
        items: [
          { text: 'Raw Messages (Fast Parser)', link: '/guide/raw-message' },
          { text: 'Structured Messages (Object Model)', link: '/guide/message' },
          { text: 'Query & Grammar Selection', link: '/guide/query-grammar' },
          { text: 'MLLP & Networking', link: '/guide/mllp' }
        ]
      },
      {
        text: 'Resources',
        items: [
          { text: 'Real-World Examples', link: '/guide/examples' },
          { text: 'Versioning & Releases', link: '/guide/versioning' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/blushift-io/hl7v2' }
    ],

    search: {
      provider: 'local'
    },

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2026 blushift-io'
    }
  }
})
