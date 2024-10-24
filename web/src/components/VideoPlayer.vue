<template>
  <div
    ref="plyrVue"
    class="lg:w-64 lg:h-80 rounded-md"
    @mouseover="handleMouseOver"
    @mouseleave="handleMouseLeave"
  >
    <div v-if="!canPlay || isLoading">
      <div
        class="bg-gray-300 dark:bg-gray-800 animate-pulse lg:h-80 sm:h-[90vh] rounded"
      ></div>
    </div>
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
import { ref, onMounted, useTemplateRef, unref } from 'vue'
import { usePlyrVue, PlyrVue } from 'plyr-vue'

import { useUpdaUserProfile, useFetchVideoLink } from '@/api/index'
import type { TikTokItem } from '@/types'
import 'plyr-vue/dist/plyr-vue.css'

const { video } = defineProps<{
  video: TikTokItem
}>()

const canPlay = ref(false)
const plyrVue = useTemplateRef('plyrVue')

const [registerVideoPlayer, videoPlayerInstance] = usePlyrVue({
  controls: [],
})

const { mutateAsync: updateProfile } = useUpdaUserProfile()

const { data: videoLink, isLoading } = useFetchVideoLink(video.video_url)

const handleMouseOver = () => {
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

    // Catching error and hiding video player
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    videoPlayerInstance.value.once('error', e => {
      console.log(e.detail.plyr.source)

      videoWrapper.children[1].classList.add('hidden')
      canPlay.value = false

      setTimeout(() => {
        videoWrapper.children[1].classList.remove('hidden')
        canPlay.value = true
      }, 5 * 1000)
    })

    videoPlayerInstance.value.on('canplay ', () => {
      canPlay.value = true
    })

    videoPlayerInstance.value.on('play', () => {
      updateProfile(video)
    })
  } catch (e) {
    console.error(e)
  }
})
</script>
