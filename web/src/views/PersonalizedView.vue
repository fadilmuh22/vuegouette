<template>
  <div class="dark:bg-gray-900 dark:text-white min-h-screen">
    <div class="container mx-auto p-4 dark:bg-dark-900">
      <div class="sm:mx-3 lg:mx-80">
        <VideoGrid :videos :isLoading="isFetching" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onActivated } from 'vue'

import { usePresonalizedVideoQuery } from '@/api'
import VideoGrid from '@/components/VideoGrid.vue'

const { data: videos, isFetching, refetch } = usePresonalizedVideoQuery()

onActivated(() => {
  refetch()
})
</script>
