<template>
  <div class="dark:bg-neutral-light dark:text-white min-h-screen">
    <div class="container mx-auto lg:p-4 dark:bg-dark-900">
      <div class="sm:mx-3 lg:mx-60">
        <VideoGrid
          :videos
          :isLoading="isLoading || isFetchingNextPage"
          :isPending
          :isFetchingNextPage
          :isSearch="false"
          :fetchNextPage
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'

import { usePersonalizedVideos } from '@/api'
import VideoGrid from '@/components/VideoGrid.vue'

const {
  data: pages,
  fetchNextPage,
  isLoading,
  isPending,
  isFetchingNextPage,
} = usePersonalizedVideos()

// Computed property to flatten the list of videos
const videos = computed(() => pages.value?.pages.flat() ?? [])

// Handle scrolling
const onScroll = () => {
  const scrollPosition = window.scrollY + window.innerHeight
  const documentHeight = document.documentElement.scrollHeight

  if (scrollPosition >= documentHeight - 100 && !isFetchingNextPage.value) {
    // Trigger loading when near bottom
    fetchNextPage()
  }
}

// Initial load on component mount
onMounted(() => {
  window.addEventListener('scroll', onScroll)
})

// Cleanup event listener
onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
})
</script>
