import { createRequire } from 'module'
import { defineConfig } from 'vitepress'
import { generateSidebarChapter } from './side_bar.js'

const require = createRequire(import.meta.url)

const chapters = generateSidebarChapter('en_US', new Map([
  ['introduction', 'Introduction'],
]))

export default defineConfig({
  lang: 'en-US',

  description: 'A cli tool in Go.',

  themeConfig: {
    nav: nav(),

    lastUpdatedText: 'Last updated at',

    sidebar: chapters,

    socialLinks: [
      { icon: 'github', link: 'https://github.com/easysoft/quickon_cli' },
    ],

    editLink: {
      pattern: 'https://github.com/easysoft/quickon_cli/edit/master/docs/:path',
      text: 'Edit this page on GitHub'
    },

    outline: {
      level: 'deep',
      label: 'On this page',
    },

  }
})

function nav() {
  return [
    { text: 'Home', link: '/' },
    { text: 'Configuration', link: '/configuration/configuration-reference' },
    {
      text: 'Download',
      items: [
        { text: 'Open-source Edition', link: 'https://github.com/easysoft/quickon_cli/releases/' },
      ]
    }
  ]
}
