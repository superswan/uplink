<template>
    <div>
        <header class="shadow-sm bg-gray-200">
            <nav class="container mx-auto p-4 flex justify-between">
                <NuxtLink class="font-bold"><img src=""/>UPLINK UI</NuxtLink>
                <ul class="flex gap-4">
                    <li><NuxtLink to="/">Dashboard</NuxtLink></li>
                    <li><NuxtLink to="/builder">Builder</NuxtLink></li>
                </ul>
            </nav>
        </header>

        <!-- Page Content -->
        <div class="container mx-auto p-4">
            <slot />
        </div>
    </div>
</template>

<style scoped>
.router-link-exact-active {
    color:cornflowerblue;
}

nav {
    color: text-gray-900;
}

tr {
  cursor: pointer;
}

.green-dot, .red-dot {
  display: inline-block;
  width: 10px;
  height: 10px;
  border-radius: 50%; /* To make the dot circular */
}

</style>

