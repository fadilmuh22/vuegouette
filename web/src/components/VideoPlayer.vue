<template>
  <PlyrVue @register="registerVideoPlayer">
    <video
      class="max-h-80"
      controls
      playsinline
      :src="getVideoUrl(video.video_url)"
    ></video>
  </PlyrVue>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { usePlyrVue, PlyrVue } from 'plyr-vue'

import type { TikTokItem } from '@/types'
import 'plyr-vue/dist/plyr-vue.css'

const { video } = defineProps<{
  video: TikTokItem
}>()

const [registerVideoPlayer, videoPlayerInstance] = usePlyrVue({})

onMounted(() => {
  try {
    videoPlayerInstance.value.once('error', e => {
      console.log(e.detail.plyr.source)
      videoPlayerInstance.value.play()
    })
  } catch (e) {
    console.error(e)
  }
})

const getVideoUrl = (videoUrl: string) => {
  return `https://tikcdn.io/ssstik/${videoUrl.split('/').pop()}`
}
</script>
