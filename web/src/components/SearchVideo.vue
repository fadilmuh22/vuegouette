<template>
  <div class="container mx-auto p-4">
    <div class="sm:mx-3 lg:mx-24">
      <h1 class="text-2xl font-bold mb-4">TikTok Search</h1>
      <input
        type="text"
        v-model="searchTerm"
        @keyup.enter="onSearch"
        placeholder="Search for videos..."
        class="w-full p-2 border border-gray-300 rounded mb-4"
      />

      <div
        class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
        v-if="isLoading"
      >
        <div class="skeleton-loading" v-for="i in 6" :key="i">
          <!-- Skeleton loading effect -->
          <div class="bg-gray-300 animate-pulse h-40 rounded"></div>
          <div class="bg-gray-300 animate-pulse h-6 rounded mt-2"></div>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="video in videos"
          :key="video.video_url"
          class="border rounded-lg shadow-md overflow-hidden"
        >
          <div class="w-full h-80 object-cover">
            <VideoPlayer :video="video" />
          </div>
          <div class="p-4">
            <h2 class="font-semibold">{{ video.user_name }}</h2>
            <p class="text-gray-500">{{ video.video_count }} views</p>
            <a
              :href="video.video_url"
              target="_blank"
              class="text-blue-500 hover:underline"
              >Watch Video</a
            >
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import axios from 'axios'
import VideoPlayer from './VideoPlayer.vue'
import type { TikTokItem } from '@/types'

const searchTerm = ref('')

const fetchVideos = async (): Promise<TikTokItem[]> => {
  console.log('Fetching videos...')
  const response = await axios.get(
    `http://localhost:1323/api/video?keyword=${searchTerm.value}`,
  )
  return response.data.data as TikTokItem[] // Adjust this based on your API response structure
}

const {
  data: videos,
  isLoading,
  refetch,
} = useQuery({
  queryKey: ['videos', searchTerm],
  queryFn: fetchVideos,
  enabled: false, // Disable auto-fetching on component mount
  refetchOnWindowFocus: false,
})

const onSearch = () => {
  console.log('Search term:', searchTerm.value)
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
