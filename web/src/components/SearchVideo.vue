<template>
  <div class="container mx-auto p-4 dark:bg-dark-900">
    <div class="sm:mx-3 lg:mx-80">
      <div class="flex justify-center">
        <h1 class="text-2xl font-bold mb-4 dark:text-white">TikTok Search</h1>
      </div>
      <input
        type="text"
        v-model="searchTerm"
        @keyup.enter="onSearch"
        placeholder="Search for videos..."
        class="w-full p-2 border border-gray-300 rounded mb-4 dark:bg-gray-700 dark:text-white dark:border-gray-600"
      />

      <VideoGrid :videos="videos" :isLoading="isLoading" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import VideoGrid from './VideoGrid.vue'
import { useSearchVideoQuery } from '@/api'

const searchTerm = ref('')

const { data: videos, isLoading, refetch } = useSearchVideoQuery(searchTerm)

const onSearch = () => {
  if (searchTerm.value) {
    refetch()
  }
}
</script>

<style>
.skeleton-loading {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
