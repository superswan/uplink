import { defineStore } from 'pinia'

const STORE_NAME = 'main'
const SETTINGS_LOCAL_STORAGE_KEY = 'settings'

const getDefaultSettings = () => ({
  apiEndpoint: 'http://127.0.0.1:8090',
  listenerEndpoint: 'http://127.0.0.1:8081',
})

const getSettings = () => {
  const settings = localStorage.getItem(STORE_NAME)

  return settings ? JSON.parse(settings) : getDefaultSettings()
}

export const useStore = defineStore(STORE_NAME, {
  state: () => ({
    settings: getSettings(),
  }),
  actions: {
    updateSettings(partialSettings) {
      this.settings = {
        ...this.settings,
        ...partialSettings,
      }
    },
  },
})