<template>
  <!-- dark for default,cv -->
  <div id="main-wrapper" class="bg-background">
    <SideBar />
    <main class="sm:ml-64">
      <RouterView />
    </main>

    <!-- TODO: Create tags customizer for personalized videos -->
    <!-- Could be chip temporary for filtering or user can save their personlization  -->
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

import SideBar from './components/SideBar.vue'
import { useCreateGuestUser, useGetUserProfileKeyword } from '@/api'

const { mutate: createGuestUser } = useCreateGuestUser()
const { mutate: getUserProfileKeyword } = useGetUserProfileKeyword()

onMounted(() => {
  createGuestUser()
  getUserProfileKeyword()
})
</script>
