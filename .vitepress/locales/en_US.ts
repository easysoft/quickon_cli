/**
 * Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
 * Use of this source code is covered by the following dual licenses:
 * (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
 * (2) Affero General Public License 3.0 (AGPL 3.0)
 * license that can be found in the LICENSE file.
 */

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
