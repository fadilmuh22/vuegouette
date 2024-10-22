<template>
  <div class="search-container">
    <input
      v-model="searchQuery"
      @keyup.enter="searchVideos"
      placeholder="Search TikTok videos"
      type="text"
    />
    <button @click="searchVideos">Search</button>
    <div v-if="loading">Searching...</div>
    <div v-if="videos.length > 0">
      <h2>Search Results:</h2>
      <div class="video-list">
        <div v-for="(video, index) in videos" :key="index" class="video-item">
          <img :src="video.thumbnail" @click="playVideo(video.url)" />
          <p>{{ video.title }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import axios from 'axios'

export default {
  data() {
    return {
      searchQuery: '',
      videos: [],
      loading: false,
    }
  },
  methods: {
    async searchVideos() {
      this.loading = true
      try {
        const response = await axios.get(`http://localhost:1323/api/video`, {
          params: { keyword: this.searchQuery },
        })
        this.videos = response.data.videos.map(video => ({
          title: video.title,
          thumbnail: video.thumbnailUrl,
          url: video.videoUrl,
        }))
      } catch (error) {
        console.error('Error fetching videos:', error)
      } finally {
        this.loading = false
      }
    },
    playVideo(url) {
      this.$emit('play-video', url) // Trigger video player component
    },
  },
}
</script>

<style scoped>
.search-container {
  text-align: center;
  margin-top: 20px;
}
.video-list {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
}
.video-item {
  margin: 10px;
  cursor: pointer;
}
.video-item img {
  width: 200px;
  height: auto;
}
</style>
