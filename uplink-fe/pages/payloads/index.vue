<template>
    <div>
        <UTable class="min-w-full bg-white" v-model="selected" :columns="columns" :rows="payloads" :empty-state="{ icon: 'i-heroicons-circle-stack-20-solid', label: 'No items.' }">
        </UTable> 
    </div>
</template>

<script setup>
const columns = [{
    key: 'client_id',
    label: 'Type'
    }, {
    key: 'ipaddress',
    label: 'Host IP'
    }, {
    key: 'port',
    label: 'Host Port' 
    }, {
    key: 'sessionCount',
    label: 'Implant ID'
    }, {
    key: 'created',
    label: 'Download'
    }
]
</script>

<style scoped>

</style>