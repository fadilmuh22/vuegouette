<template>
  <div id="main-wrapper" class="dark dark:bg-gray-900">
    <SideBar />
    <main class="sm:ml-64" data-mode="dark">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

import SideBar from './components/SideBar.vue'
import { useCreateGuestUser } from '@/api'

const { mutate: createGuestUser } = useCreateGuestUser()

onMounted(() => {
  createGuestUser(null)
})
</script>
