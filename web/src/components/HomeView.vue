<template>
  <div id="app">
    <div class="flex content-center justify-center">
      <img
        alt="Vue logo"
        class="logo p-4"
        src="@/assets/logo.svg"
        width="125"
        height="125"
      />
    </div>
    <SearchVideo />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import SearchVideo from './SearchVideo.vue'

export default defineComponent({
  components: {
    SearchVideo,
  },
})
</script>

<style>
#app {
  text-align: center;
  margin-top: 20px;
}
</style>
