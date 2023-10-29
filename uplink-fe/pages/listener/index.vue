<script setup>
import { computed, ref } from 'vue';
import { useMainStore } from '@/store/mainStore'; // Adjust the import as per your file location
import { storeToRefs } from 'pinia';

// Initialize store
const mainStore = useMainStore();

// Convert store state to refs
const { apiEndpoint, listenerEndpoint } = storeToRefs(mainStore);

// Function to parse URL to get IP and port
const parseUrl = (url) => {
  try {
    const urlObj = new URL(url);
    return { ip: urlObj.hostname, port: urlObj.port || '80' };
  } catch (e) {
    console.error('Invalid URL', e);
    return { ip: '', port: '80' };
  }
};

// Function to format URL from IP and port
const formatUrl = (ip, port, path) => `http://${ip}:${port}${path}`;

// Get the current values from the store and parse them to get IP and Port
const { ip: apiIpValue, port: apiPortValue } = parseUrl(apiEndpoint.value);
const { ip: listenerIpValue, port: listenerPortValue } = parseUrl(listenerEndpoint.value);

// Input refs with default values from the store
const apiIp = ref(apiIpValue);
const apiPort = ref(apiPortValue);
const listenerIp = ref(listenerIpValue);
const listenerPort = ref(listenerPortValue);

// Function to update the endpoints in the store
const updateEndpoints = () => {
  mainStore.updateSettings({
    apiEndpoint: formatUrl(apiIp.value, apiPort.value),
    listenerEndpoint: formatUrl(listenerIp.value, listenerPort.value)
  });
};
</script>

<template>
  <div>
    <div>
      <label>API IP: <input v-model="apiIp" @change="updateEndpoints" /></label>
      <label>API Port: <input v-model="apiPort" @change="updateEndpoints" /></label>
    </div>
    <div>
      <label>Listener IP: <input v-model="listenerIp" @change="updateEndpoints" /></label>
      <label>Listener Port: <input v-model="listenerPort" @change="updateEndpoints" /></label>
    </div>
    <p>API Endpoint is: {{ apiEndpoint }}</p>
    <p>Listener Endpoint is: {{ listenerEndpoint }}</p>
  </div>
</template>