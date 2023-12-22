/**
 * Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
 * Use of this source code is covered by the following dual licenses:
 * (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
 * (2) Affero General Public License 3.0 (AGPL 3.0)
 * license that can be found in the LICENSE file.
 */

import { defineConfig } from 'vitepress'
import locales from './locales'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Quickon Cli",
  description: "命令行",

  locales: locales.locales,

  lastUpdated: true,

  themeConfig: {
    search: {
      provider: 'local',
      options: {
        locales: {
          zh_CN: {
            translations: {
              button: {
                buttonText: '搜索文档',
                buttonAriaLabel: '搜索文档'
              },
              modal: {
                noResultsText: '无法找到相关结果',
                resetButtonTitle: '清除查询条件',
                footer: {
                  selectText: '选择',
                  navigateText: '切换'
                }
              }
            }
          }
        },
    },
  },
  }
})
