<template>
  <div
    v-if="isPending && !isFetchingNextPage && isSearch"
    class="lg:min-h-[80vh] flex items-center justify-center text-onbackground"
  >
    Search for videos by using the search bar above.
  </div>

  <div
    v-if="!isPending && !videos?.length"
    class="lg:min-h-[80vh] flex items-center justify-center text-onbackground"
  >
    No videos found.
  </div>

  <div
    v-if="videos?.length"
    class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
  >
    <div
      v-for="(video, index) in videos"
      :key="video.video_url"
      class="rounded-md shadow-md overflow-hidden"
    >
      <VideoPlayer
        :video="video"
        :isFullscreen="false"
        @click="() => handleClickVideo(video, index)"
      />
      <div class="p-1">
        <h2 class="font-semibold dark:text-white text-onbackground">
          @{{ video.user_name }}
        </h2>
        <p class="text-gray-500">{{ video.video_count }} views</p>
      </div>
    </div>
  </div>

  <div
    class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
    v-if="isLoading"
  >
    <div class="skeleton-loading" v-for="i in 6" :key="i">
      <!-- Skeleton loading effect -->
      <div class="bg-gray-300 animate-pulse lg:h-80 sm:h-[90vh] rounded"></div>
      <div class="bg-gray-300 animate-pulse lg:h-6 sm:h-0 rounded mt-2"></div>
    </div>
  </div>

  <div
    v-if="videos?.length && isFetchingNextPage"
    class="flex items-center justify-center p-3"
  >
    Loading more videos...
  </div>

  <VideoFocused
    v-if="selectedVideo"
    :video="selectedVideo"
    :onClose="handleVideoFocusedClose"
    @wheel="handleScroll"
    @touchstart="handleTouchStart"
    @touchend="handleTouchEnd"
  />
</template>

<script setup lang="ts">
import { defineProps, ref, onMounted, onBeforeUnmount, watch } from 'vue'

import type { TikTokItem } from '@/types'
import VideoFocused from './VideoFocused.vue'
import VideoPlayer from './VideoPlayer.vue'

const {
  videos,
  isLoading,
  isPending,
  isFetchingNextPage,
  isSearch,
  fetchNextPage,
} = defineProps<{
  videos: TikTokItem[] | undefined
  isLoading: boolean
  isPending: boolean
  isFetchingNextPage: boolean
  isSearch: boolean
  fetchNextPage: () => void
}>()

const selectedVideo = ref<TikTokItem | undefined>(undefined)
const selectedVideoIndex = ref<number>(-1)
const fullscreenMode = ref<boolean>(false)

let startX = 0
let endX = 0
const scrollThreshold = 100 // Threshold for scroll detection

// Add keyboard listener when component is mounted and remove on unmount
onMounted(() => {
  window.addEventListener('keydown', handleShortcut)
})
onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleShortcut)
})

const handleClickVideo = (video: TikTokItem, index: number) => {
  selectedVideo.value = video
  selectedVideoIndex.value = index
  fullscreenMode.value = true
}

const handleVideoFocusedClose = () => {
  selectedVideo.value = undefined
  selectedVideoIndex.value = -1
  fullscreenMode.value = false
}

const nextVideo = () => {
  if (selectedVideoIndex.value < videos!.length - 1) {
    handleClickVideo(
      videos![selectedVideoIndex.value + 1],
      selectedVideoIndex.value + 1,
    )
  }

  if (selectedVideoIndex.value === videos!.length - 1) {
    fetchNextPage()
  }
}

const prevVideo = () => {
  if (selectedVideoIndex.value > 0) {
    handleClickVideo(
      videos![selectedVideoIndex.value - 1],
      selectedVideoIndex.value - 1,
    )
  }
}

const handleShortcut = (e: KeyboardEvent) => {
  if (!videos || !videos.length) return
  if (!fullscreenMode.value) return

  if (e.key === 'Escape') {
    handleVideoFocusedClose()
  }
  // page up or arraow up for previous video.
  if (
    (e.key === 'PageUp' || e.key === 'ArrowUp') &&
    selectedVideoIndex.value > 0
  ) {
    prevVideo()
  }

  // page down or arrow down for next video.
  if (
    (e.key === 'PageDown' || e.key === 'ArrowDown') &&
    selectedVideoIndex.value < videos.length - 1
  ) {
    nextVideo()
  }
}

// Touch event handlers for swipe functionality
const handleTouchStart = (event: TouchEvent) => {
  startX = event.touches[0].clientX
}

const handleTouchEnd = (event: TouchEvent) => {
  endX = event.changedTouches[0].clientX
  handleSwipe()
}

const handleSwipe = () => {
  if (!fullscreenMode.value) return

  if (startX - endX > 50) {
    // Swiped left
    nextVideo()
  } else if (endX - startX > 50) {
    // Swiped right
    prevVideo()
  }
}

// Handle scroll event for next and previous video
const handleScroll = (event: WheelEvent) => {
  if (!fullscreenMode.value) return

  if (event.deltaY > scrollThreshold) {
    // Scrolled down
    nextVideo()
  } else if (event.deltaY < -scrollThreshold) {
    // Scrolled up
    prevVideo()
  }
}

watch(fullscreenMode, isFullscreen => {
  document.body.style.overflow = isFullscreen ? 'hidden' : ''
})
</script>
