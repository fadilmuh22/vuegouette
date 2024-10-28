<template>
  <div class="z-[800] fixed inset-0 0 h-[100dvh] w-[100dvw]">
    <!-- Blur background for VideoPlayer -->
    <div
      class="absolute inset-0 bg-black bg-opacity-70 backdrop-blur-3xl"
    ></div>

    <!-- Close button -->
    <button
      @click="onClose(video)"
      class="absolute top-4 left-4 z-[1000] p-2 rounded-full bg-white bg-opacity-50"
    >
      <svg
        class="w-3 h-3"
        aria-hidden="true"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 14 14"
      >
        <path
          stroke="currentColor"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"
        />
      </svg>
    </button>

    <div class="z-[1000] h-full w-full flex items-center justify-center">
      <VideoPlayer :key="video.video_url" :video :isFullscreen="true" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'

import type { TikTokItem } from '@/types'
import VideoPlayer from './VideoPlayer.vue'

const { video } = defineProps<{
  video: TikTokItem
  onClose: (video: TikTokItem) => void
}>()
</script>
