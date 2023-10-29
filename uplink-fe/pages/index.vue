<script setup>
import { useIntervalFn } from '@vueuse/core' // VueUse helper, install it
import { ref, watch, computed } from 'vue'

const POCKETBASE_URL = import.meta.env.VITE_POCKETBASE_URL;
const UPLINK_URL = import.meta.env.VITE_STATUS_URL;

const columns = [{
    key: 'client_id',
    label: 'UUID'
    }, {
    key: 'ipaddress',
    label: 'IP Address'
    }, {
    key: 'port',
    label: 'Port' 
    }, {
    key: 'sessionCount',
    label: 'Sessions'
    }, {
    key: 'created',
    label: 'First'
    }, {
    key: 'updated',
    label: 'Last'
    }
]

let blink = ref(false); // Controls whether the span should blink
    const {
        data: clients,
        refresh: refreshClients
    } = await useFetch(`${POCKETBASE_URL}/api/collections/clients/records`); //Pocketbase DB

    const {
        pending,
        data: activeClients,
        error,
        refresh: refreshActiveClients
    } = await useFetch(`${UPLINK_URL}/status/`) //UPLINK server

    useIntervalFn(() => {
  console.log(`fetching client status from server @ ${new Date().toISOString()}`)
  refreshActiveClients() // will call the 'todos' endpoint, just above
  console.log(`fetching client data from db backend @ ${new Date().toISOString()}`)
  refreshClients()
  console.log(clients.value.items)
}, 3500)

const activeClientsLength = computed(() => activeClients.value?.length || 0);
watch(activeClientsLength, (newVal, oldVal) => {
    if (newVal !== oldVal) {
    blink.value = true; // Blink when activeClients changes
    setTimeout(() => blink.value = false, 3000); // Stop blinking after 1.5 seconds
    }
});

    const isActive = (clientId) => {
  return activeClients.value.includes(clientId);
};

const tableRows = computed(() => {
  if (!clients.value || !clients.value.items) return [];

  return clients.value.items.map(c => {
    const isActiveClient = activeClients.value.includes(c.client_id);

    return {
      client_id: c.client_id,
      ipaddress: c.ipaddress,
      port: c.port,
      sessionCount: c.sessionCount,
      created: c.created,
      updated: c.updated,
      isActive: isActiveClient,
      link: `/client/${c.client_id}`,
      linkClass: `hover:text-blue-500 ${isActiveClient ? '' : 'text-gray-500 cursor-not-allowed'}`,
      dotClass: isActiveClient ? 'green-dot' : 'red-dot',
    };
  });
});

const query = ref('')

const filteredRows = computed(() => {
  if (!query.value) {
    return tableRows
  }
  return people.filter((person) => {
    return Object.values(person).some((value) => {
      return String(value).toLowerCase().includes(q.value.toLowerCase())
    })
  })
})

const selected = ref([tableRows[1]])
</script>

<template>

    <div class="container columns-4">
        <div class="relative flex w-auto max-w-md flex-col items-start gap-2 overflow-hidden rounded-lg p-4 shadow-lg">
            <h2 class="text-xl font-semibold">Clients</h2>
            <span class="text-4xl mx-auto text-slate-400"> {{ clients.items?.length || 0 }}</span>
        </div>
                <div class="relative flex w-auto max-w-md flex-col items-start gap-2 overflow-hidden rounded-lg p-4 shadow-lg">
            <h2 class="text-xl font-semibold">Online</h2>
            <transition name="flash" mode="out-in">
            <span class="text-4xl mx-auto text-slate-400" :class="{ blink: blink }"> {{ activeClients.value?.length || 0 }}</span>
            </transition>
        </div>
    </div>

    <div class="py-8 w-full">
        <div class="shadow overflow-hidden rounded border-b border-gray-200">
            <UInput v-model="query" placeholder="Search..."/>
            <UTable class="min-w-full bg-white" v-model="selected" :columns="columns" :rows="tableRows" :empty-state="{ icon: 'i-heroicons-circle-stack-20-solid', label: 'No items.' }">
               <template #client_id-data="{ row }">
                <NuxtLink 
                    :to="`/client/${row.client_id}`" 
                    class="hover:text-blue-500" 
                    :class="{ 'text-gray-500': !isActive(row.client_id), 'cursor-not-allowed': !isActive(row.client_id) }">
                    <span :class="{'green-dot': activeClients.includes(row.client_id), 'red-dot': !activeClients.includes(row.client_id)}"></span>        
                    {{ row.client_id }}
                </NuxtLink>
        </template> 
            </UTable> 
        </div>
    </div>
</template>



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