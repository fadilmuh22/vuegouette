<template>
  <div class="dark:bg-neutral-light dark:text-white min-h-screen">
    <div class="container lg:p-4 mx-auto">
      <div class="sm:mx-3 lg:mx-60">
        <div class="flex justify-center">
          <h1 class="text-2xl font-bold mb-4 text-onbackground">
            TikTok Search
          </h1>
        </div>
        <input
          type="text"
          v-model="searchTerm"
          @keyup.enter="onSearch"
          placeholder="Search for videos..."
          class="w-full p-2 border border-gray-300 rounded mb-4 dark:bg-neutral-dark dark:text-white dark:border-gray-600"
        />

        <VideoGrid
          :videos="videos"
          :isLoading="isFetchingNextPage"
          :isPending
          :isFetchingNextPage
          :isSearch="true"
          :fetchNextPage
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted, watch } from 'vue'

import { useSearchVideos } from '@/api'
import VideoGrid from '@/components/VideoGrid.vue'

const searchTerm = ref('')
const {
  data: pages,
  fetchNextPage,
  isPending,
  isFetchingNextPage,
} = useSearchVideos(searchTerm)

// Computed property to flatten the list of videos
const videos = computed(() => pages.value?.pages.flat() ?? [])

// Handle scrolling
const onScroll = () => {
  const scrollPosition = window.scrollY + window.innerHeight
  const documentHeight = document.documentElement.scrollHeight

  if (
    scrollPosition >= documentHeight - 100 &&
    !isFetchingNextPage.value &&
    searchTerm.value.length > 3
  ) {
    // Trigger loading when near bottom
    fetchNextPage()
  }
}

watch(searchTerm, () => {
  if (searchTerm.value.length > 3) {
    window.location.hash = searchTerm.value
  }
})

onMounted(() => {
  if (window.location.hash) {
    searchTerm.value = decodeURIComponent(window.location.hash).slice(1)
    onSearch()
  }
  window.addEventListener('scroll', onScroll)
})

// Cleanup event listener
onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
})

const onSearch = () => {
  if (searchTerm.value && searchTerm.value.length > 3) {
    fetchNextPage()
  }
}
</script>
