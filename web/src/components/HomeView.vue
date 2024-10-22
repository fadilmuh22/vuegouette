<template>
  <div id="app">
    <h1>Custom TikTok Video Player</h1>
    <SearchVideo @play-video="setVideoUrl" />
    <VideoPlayer v-if="videoUrl" :videoUrl="videoUrl" />
  </div>
</template>

<script lang="ts">
import SearchVideo from './SearchVideo.vue'
import VideoPlayer from './VideoPlayer.vue'

export default {
  components: {
    SearchVideo,
    VideoPlayer,
  },
  data() {
    return {
      videoUrl: null,
    }
  },
  methods: {
    setVideoUrl(url) {
      this.videoUrl = url
    },
  },
}
</script>

<style>
#app {
  text-align: center;
  margin-top: 20px;
}
</style>
