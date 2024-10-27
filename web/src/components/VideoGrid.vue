<template>
  <div
    class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
    v-if="isLoading"
  >
    <div class="skeleton-loading" v-for="i in 6" :key="i">
      <!-- Skeleton loading effect -->
      <div
        class="bg-gray-300 dark:bg-gray-800 animate-pulse lg:h-80 sm:h-[90vh] rounded"
      ></div>
      <div
        class="bg-gray-300 dark:bg-gray-800 animate-pulse lg:h-6 sm:h-0 rounded mt-2"
      ></div>
    </div>
  </div>

  <div
    v-if="!isLoading"
    class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
  >
    <div
      v-for="video in videos"
      :key="video.video_url"
      class="rounded-md shadow-md overflow-hidden"
    >
      <VideoPlayer :video="video" />
      <div class="p-1">
        <h2 class="font-semibold">@{{ video.user_name }}</h2>
        <p class="text-gray-500">{{ video.video_count }} views</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'

import type { TikTokItem } from '@/types'
import VideoPlayer from './VideoPlayer.vue'

const { videos, isLoading } = defineProps<{
  videos: TikTokItem[] | undefined
  isLoading: boolean
}>()
</script>
