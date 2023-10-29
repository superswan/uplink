<template>
    <div>
      <UForm :validate="validate" :state="state" @submit="submit">
        <UFormGroup label="Host IP" name="hostIp">
          <UInput v-model="state.hostIp" placeholder="Enter Host IP" />
        </UFormGroup>
        <UFormGroup label="Host Port" name="hostPort">
          <UInput v-model="state.hostPort" placeholder="Enter Host Port" type="number" />
        </UFormGroup>
        <UFormGroup label="Payload Type" name="payloadType">
          <USelectMenu v-model="state.payloadType" :options="payloadTypes" />
        </UFormGroup>
        <UButton type="submit">Submit</UButton>
      </UForm>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue'
  import type { FormError, FormSubmitEvent } from '@nuxt/ui/dist/runtime/types'
  
  const payloadTypes = ['Python (All Platforms)', 'Go', 'Exe (Windows x86)', 'Exe (Windows x64)', 'Dll (Windows x86)', 'ELF (Linux)', 'Shell Script (Linux)'] 
  const state = ref({
    hostIp: undefined,
    hostPort: undefined,
    payloadType: payloadTypes[0] 
  })
  
  const validate = (state: any): FormError[] => {
    const errors = []
    if (!state.hostIp) errors.push({ path: 'hostIp', message: 'Host IP is required' })
    if (!state.hostPort) errors.push({ path: 'hostPort', message: 'Host Port is required' })
    if (!state.payloadType) errors.push({ path: 'payloadType', message: 'Payload Type is required' })
    return errors
  }
  
  async function submit (event: FormSubmitEvent<any>) {
     try {
        const response = await fetch('http://127.0.0.1:5000/generate_script', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                user_values: {
                    server: state.value.hostIp,
                    port: state.value.hostPort
                }
            }),
        });

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const data = await response.json();
        const scriptId = data.script_id;
        const downloadLink = `http://127.0.0.1:5000/download_script/${scriptId}`;
        
        // Handle the response, for example, show the download link to the user
        console.log(downloadLink);

    } catch (error) {
        // Handle errors, e.g., show an error message to the user
        console.error('Error submitting form', error);
    }
  }
  </script>

<style scoped>

</style>