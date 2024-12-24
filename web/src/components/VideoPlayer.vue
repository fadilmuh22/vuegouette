<template>
  <div
    ref="plyrVue"
    :class="
      (isFullscreen ? 'lg:w-[36rem]' : 'lg:h-80 lg:w-64') +
      ' h-[90dvh] w-[90vw] rounded-md '
    "
    @mouseover="handleMouseOver"
    @mouseleave="handleMouseLeave"
  >
    <PlyrVue
      class="h-full w-full cursor-pointer"
      @register="registerVideoPlayer"
    >
      <video
        playsinline
        :autoplay="isFullscreen"
        :muted="!isFullscreen"
        :src="videoLink"
      ></video>
    </PlyrVue>
  </div>
</template>

<script setup lang="ts">
import { onMounted, useTemplateRef, unref, onBeforeUnmount } from 'vue'
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { usePlyrVue, PlyrVue } from 'plyr-vue'

import { useUpdaUserProfile, useFetchVideoLink } from '@/api/index'
import type { TikTokItem } from '@/types'
import 'plyr-vue/dist/plyr-vue.css'

const { video, isFullscreen = false } = defineProps<{
  video: TikTokItem
  isFullscreen: boolean | undefined
}>()

const plyrVue = useTemplateRef<HTMLDivElement>('plyrVue')

const [registerVideoPlayer, videoPlayerInstance] = usePlyrVue({
  controls: isFullscreen ? ['progress', 'volume', 'play'] : [],
})

const { mutateAsync: updateProfile } = useUpdaUserProfile()

const { data: videoLink, isLoading } = useFetchVideoLink(video.video_url)

const handleMouseOver = () => {
  if (isLoading.value) return
  if (!isFullscreen) {
    videoPlayerInstance.value.volume = 5
    videoPlayerInstance.value.play()
  }
}

const handleMouseLeave = () => {
  if (isFullscreen) return
  videoPlayerInstance.value.stop()
}

onMounted(() => {
  try {
    const plyrElement = unref(plyrVue)?.children[0]
    const videoWrapper = plyrElement?.firstElementChild
    videoWrapper?.classList.add('h-full')

    videoPlayerInstance.value.on('play', () => {
      updateProfile(video)
    })
  } catch (e) {
    console.error(e)
  }
})

onBeforeUnmount(() => {
  videoPlayerInstance.value.stop()
})
</script>
