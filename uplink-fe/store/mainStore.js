import { defineStore } from 'pinia';

export const useMainStore = defineStore('main', {
  state: () => ({
    apiEndpoint: 'http://127.0.0.1:8090',
    listenerEndpoint: 'http://127.0.0.1:8081',
  }),
  actions: {
    updateSettings({ apiEndpoint, listenerEndpoint }) {
      if (apiEndpoint) this.apiEndpoint = apiEndpoint;
      if (listenerEndpoint) this.listenerEndpoint = listenerEndpoint;
    },
  },
});