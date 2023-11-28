import { defineConfig } from 'vitepress'
import en_US from './en_US'

export default defineConfig({
  locales: {
    root: {
      label: 'English',
      lang: en_US.lang,
      themeConfig: en_US.themeConfig,
      description: en_US.description
    },
  }
})
