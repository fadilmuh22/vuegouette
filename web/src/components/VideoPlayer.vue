<template>
  <div
    ref="plyrVue"
    class="lg:w-64 lg:h-80 rounded-md"
    @mouseover="handleMouseOver"
    @mouseleave="handleMouseLeave"
  >
    <PlyrVue
      class="h-full w-full cursor-pointer"
      @register="registerVideoPlayer"
    >
      <video
        playsinline
        class="h-full w-full object-cover rounded-md"
        :src="videoLink"
      ></video>
    </PlyrVue>
  </div>
</template>

<script setup lang="ts">
import { onMounted, useTemplateRef, unref } from 'vue'
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { usePlyrVue, PlyrVue } from 'plyr-vue'

import { useUpdaUserProfile, useFetchVideoLink } from '@/api/index'
import type { TikTokItem } from '@/types'
import 'plyr-vue/dist/plyr-vue.css'

const { video } = defineProps<{
  video: TikTokItem
}>()

const plyrVue = useTemplateRef('plyrVue')

const [registerVideoPlayer, videoPlayerInstance] = usePlyrVue({
  controls: [],
})

const { mutateAsync: updateProfile } = useUpdaUserProfile()

const { data: videoLink, isLoading } = useFetchVideoLink(video.video_url)

const handleMouseOver = () => {
  if (isLoading.value) {
    return
  }
  videoPlayerInstance.value.play()
}

const handleMouseLeave = () => {
  videoPlayerInstance.value.stop()
}

onMounted(() => {
  try {
    // Getting the video wrapper element and set it to match parent container
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    const plyrElement = unref(plyrVue)?.children[0]
    const videoWrapper = plyrElement.firstElementChild
    videoWrapper.classList.add('h-full')

    videoPlayerInstance.value.on('play', () => {
      updateProfile(video)
    })
  } catch (e) {
    console.error(e)
  }
})
</script>
