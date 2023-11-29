import { createRequire } from 'module'
import { defineConfig } from 'vitepress'
import { generateSidebarChapter } from './side_bar.js'

const require = createRequire(import.meta.url)

const chapters = generateSidebarChapter('en_US', new Map([
  ['install', 'Install'],
  ['init', 'Init'],
  ['backup', 'BackUP'],
  ['cluster', 'Cluster'],
  ['platform', 'Platform'],
  ['status', 'Status'],
  ['upgrade', 'Upgrade'],
  ['bugreport', 'Bug Report'],
  ['uninstall', 'Uninstall'],
  ['experimental', 'Experimental'],
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
    { text: 'Init', link: '/init/init' },
    { text: 'Download',link: 'https://github.com/easysoft/quickon_cli/releases/' },
  ]
}
