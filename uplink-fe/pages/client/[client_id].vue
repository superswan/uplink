<template>
    <div class="p-4">
        <h1>{{ client_id }}</h1>
        <hr>
    </div>

    <div class="shell-container">
        <div class="terminal" ref="terminalContainer">
            <pre ref="output"></pre>
        </div>
        <div class="input-area">
        <form @submit.prevent="submitCommand" autocomplete="off">
            <input v-model="command" type="text" placeholder="Enter command..." id="input">
            <button type="submit" id="submit">Submit</button>
            </form>
        </div>
    </div>


</template>

<script setup>
import { ref, nextTick, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const { client_id } = useRoute().params

let command = ref('')  // command to be entered by the user
let output = ref(null)  // ref to terminal output
let terminalContainer = ref(null)  // ref to terminal container

// Function to submit command
const submitCommand = async () => {
  if(output.value && command.value) {
    output.value.innerHTML += `\n> ${command.value}`

    // Construct fetch options
    const fetchOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ command: command.value })
    }

    // Call API
    const { data, error } = await useFetch(`http://localhost:8081/command/${client_id}`, fetchOptions)
    if (data.value) {
        const parsedData = JSON.parse(data.value);
        console.log(parsedData.output);
      output.value.innerHTML += `\n${parsedData.output}`
    } else if (error) {
      output.value.innerHTML += `\nError: ${error.message}`
    }

     // Wait for DOM to update
    await nextTick()

    // Scroll to the bottom
    terminalContainer.value.scrollTop = terminalContainer.value.scrollHeight

    command.value = ''
  }
}

// Ensure output is defined after component is mounted
onMounted(() => {
  if(!output.value) {
    console.error('Could not get output element')
  }
})

</script>


<style scoped>

.shell-container {
    width: 100%;
    height: 800px;
    background-color: black;
    border-radius: 10px;
    display: flex;
    flex-direction: column;
    box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.5);
}

.terminal {
    flex: 1;
    color: #77dd88;
    overflow-y: auto;
    padding: 10px;
    font-family: monospace;
    font-size: 0.8rem;
}

.input-area {
    display: flex;
    justify-content: space-between;
    padding: 10px;
    background-color: #555555;
}

form {
    display: flex;
    width: 100%;
    flex-grow: 1;
}

#input {
    width: 93%;
    flex: 1;
    color: white;
    background-color: #333;
    border: none;
    padding: 5px;
}

#input:focus {
    outline: none;
}

#submit {
    margin-left: 10px;
    background-color: gray;
    color: rgb(209, 209, 209);
    border: none;
    padding: 5px 10px;
}

</style>