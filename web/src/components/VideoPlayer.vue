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

import { useUpdaUserProfile } from '@/api/index.ts'
import type { TikTokItem } from '@/types'
import 'plyr-vue/dist/plyr-vue.css'

const { video } = defineProps<{
  video: TikTokItem
}>()

const [registerVideoPlayer, videoPlayerInstance] = usePlyrVue({})

const { mutateAsync: updateProfile } = useUpdaUserProfile()

onMounted(() => {
  try {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    videoPlayerInstance.value.once('error', e => {
      console.log(e.detail.plyr.source)
    })

    videoPlayerInstance.value.on('play', () => {
      updateProfile(video)
    })
  } catch (e) {
    console.error(e)
  }
})

const getVideoUrl = (videoUrl: string) => {
  return `https://tikcdn.io/ssstik/${videoUrl.split('/').pop()}`
}
</script>
