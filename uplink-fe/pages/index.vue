<template>

    <div class="container columns-4">
        <div class="relative flex w-auto max-w-md flex-col items-start gap-2 overflow-hidden rounded-lg p-4 shadow-lg">
            <h2 class="text-xl font-semibold">Clients</h2>
            <span class="text-4xl mx-auto text-slate-400"> {{ clients.items.length }}</span>
        </div>
                <div class="relative flex w-auto max-w-md flex-col items-start gap-2 overflow-hidden rounded-lg p-4 shadow-lg">
            <h2 class="text-xl font-semibold">Online</h2>
            <transition name="flash" mode="out-in">
            <span class="text-4xl mx-auto text-slate-400" :class="{ blink: blink }"> {{ activeClients.length }}</span>
            </transition>
        </div>
    </div>

    <div class="py-8 w-full">
        <div class="shadow overflow-hidden rounded border-b border-gray-200">
            <table class="min-w-full bg-white">
                <thead class="bg-slate-600 text-white">
                    <tr>
                        <th class="w-1/3 text-left py-3 px-4 uppercase font-semibold text-sm">UUID</th>
                        <th class="w-1/3 text-left py-3 px-4 uppercase font-semibold text-sm">IP Address</th>
                        <th class="w-1/3 text-left py-3 px-4 uppercase font-semibold text-sm">PORT</th>
                        <th class="w-1/3 text-left py-3 px-4 uppercase font-semibold text-sm">Sessions</th>
                        <th class="w-1/3 text-left py-3 px-4 uppercase font-semibold text-sm">First</th>
                        <th class="w-1/3 text-left py-3 px-4 uppercase font-semibold text-sm">Last</th>
                    </tr>
                </thead>
                <tbody class="text-gray-700">
                    <tr v-for="c in clients.items">
                        <td class="w-1/3 text-left py-3 px-4">
                            <NuxtLink 
                            :to="`/client/${c.client_id}`" 
                            class="hover:text-blue-500" 
                            :class="{ 'text-gray-500': !isActive(c.client_id), 'cursor-not-allowed': !isActive(c.client_id) }">
                            <span :class="{'green-dot': activeClients.includes(c.client_id), 'red-dot': !activeClients.includes(c.client_id)}"></span>        {{ c.client_id }}</NuxtLink>
                        </td>
                        <td class="w-1/3 text-left py-3 px-4"> {{ c.ipaddress }}</td>
                        <td class="w-1/3 text-left py-3 px-4">{{ c.port}}</td>
                        <td class="w-1/3 text-left py-3 px-4"> {{ c.sessionCount }}</td>
                        <td class="w-1/3 text-left py-3 px-4">{{ c.created}}</td>
                        <td class="w-1/3 text-left py-3 px-4"> {{ c.updated }}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<script setup>
import { useIntervalFn } from '@vueuse/core' // VueUse helper, install it
import { ref, watch, computed } from 'vue'

let blink = ref(false); // Controls whether the span should blink
    const {
        data: clients,
        refresh: refreshClients
    } = await useFetch('http://127.0.0.1:8090/api/collections/clients/records');

    const {
        pending,
        data: activeClients,
        error,
        refresh: refreshActiveClients
    } = await useFetch('http://127.0.0.1:8081/status/')

    useIntervalFn(() => {
  console.log(`fetching client status from server @ ${new Date().toISOString()}`)
  refreshActiveClients() // will call the 'todos' endpoint, just above
  console.log(`fetching client data from db backend @ ${new Date().toISOString()}`)
  refreshClients()
}, 3500)

const activeClientsLength = computed(() => activeClients.value.length);
watch(activeClientsLength, (newVal, oldVal) => {
    if (newVal !== oldVal) {
    blink.value = true; // Blink when activeClients changes
    setTimeout(() => blink.value = false, 3000); // Stop blinking after 1.5 seconds
    }
});

    const isActive = (clientId) => {
  return activeClients.value.includes(clientId);
};

    console.log(activeClients.value)
</script>

<style scoped>
.green-dot, .red-dot {
  display: inline-block;
  width: 10px;
  height: 10px;
  border-radius: 50%; /* To make the dot circular */
}
.green-dot {
  background-color: darkseagreen;
}

.red-dot {
  background-color: firebrick;
}

.cursor-not-allowed {
  cursor: not-allowed !important;
    pointer-events: none;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}

.blink {
  animation: blink 1.5s 3; /* blink 3 times */
}

</style>