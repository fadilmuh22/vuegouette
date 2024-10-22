<template>
  <div v-if="videoUrl">
    <vue-plyr>
      <video :src="videoUrl" controls></video>
    </vue-plyr>
  </div>
</template>

<script lang="ts">
import VuePlyr from 'vue-plyr'
import 'vue-plyr/dist/vue-plyr.css'

export default {
  components: {
    VuePlyr,
  },
  props: {
    videoUrl: String,
  },
}
</script>

<style scoped>
video {
  width: 100%;
  height: auto;
}
</style>
