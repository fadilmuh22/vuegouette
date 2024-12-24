<template>
  <div v-click-out-side="handleClickOutside">
    <!-- Toggle Button -->
    <button
      type="button"
      @click="toggleSidebar"
      class="inline-flex items-center p-2 mt-2 ms-3 text-sm text-gray-500 rounded-lg sm:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600"
    >
      <span class="sr-only">Open sidebar</span>
      <svg
        class="w-6 h-6"
        aria-hidden="true"
        fill="currentColor"
        viewBox="0 0 20 20"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          clip-rule="evenodd"
          fill-rule="evenodd"
          d="M2 4.75A.75.75 0 012.75 4h14.5a.75.75 0 010 1.5H2.75A.75.75 0 012 4.75zm0 10.5a.75.75 0 01.75-.75h7.5a.75.75 0 010 1.5h-7.5a.75.75 0 01-.75-.75zM2 10a.75.75 0 01.75-.75h14.5a.75.75 0 010 1.5H2.75A.75.75 0 012 10z"
        ></path>
      </svg>
    </button>

    <!-- Sidebar -->
    <aside
      :class="[
        'fixed top-0 left-0 z-40 w-64 h-screen transition-transform',
        sidebarOpen ? 'translate-x-0' : '-translate-x-full',
        'sm:translate-x-0',
      ]"
      aria-label="Sidebar"
    >
      <div class="h-full px-3 py-4 overflow-y-auto bg-primary">
        <ul class="space-y-2 font-medium">
          <li>
            <RouterLink
              to="/"
              class="flex items-center p-2 text-neutral-light rounded-lg dark:text-white hover:bg-gray-500 dark:hover:bg-gray-700 group"
            >
              <span class="ms-3">Personalized</span>
            </RouterLink>
          </li>
          <li>
            <RouterLink
              to="/search"
              class="flex items-center p-2 text-neutral-light rounded-lg dark:text-white hover:bg-gray-500 dark:hover:bg-gray-700 group"
            >
              <span class="ms-3">Search</span>
            </RouterLink>
          </li>
          <li>
            <div class="text-neutral-light font-medium mt-4">Manage Keywords</div>
            <div class="flex flex-wrap gap-2 mt-2">
              <span
                v-for="(keyword, index) in keywords"
                :key="index"
                class="flex items-center justify-between bg-gray-300 text-black py-1 px-2 rounded-lg"
              >
                <span>{{ keyword }}</span>
                <button
                  type="button"
                  class="ml-2 text-red-500"
                  @click="removeKeyword(keyword)"
                >
                  &#10005;
                </button>
              </span>
            </div>

            <!-- Save and Cancel Buttons -->
            <div class="mt-4">
              <button
                type="button"
                class="ml-2 px-4 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600"
                @click="saveKeywords"
              >
                Save
              </button>
              <button
                type="button"
                class="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
                @click="cancelChanges"
              >
                Cancel
              </button>
            </div>
          </li>
        </ul>
      </div>
    </aside>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { clickOutSide as vClickOutSide } from '@mahdikhashan/vue3-click-outside'

import { useDeleteUserProfileKeyword } from '@/api'

const sidebarOpen = ref(false)

// Data store for the sidebar keywords
const keywords = ref<string[]>([])

// To store the original keywords before editing
const originalKeywords = ref<string[]>([])

const { mutate: deleteUserProfileKeyword } = useDeleteUserProfileKeyword()

function toggleSidebar() {
  sidebarOpen.value = !sidebarOpen.value
}

const handleClickOutside = () => {
  sidebarOpen.value = false
}

// Function to fetch keywords from localStorage or initialize with 'trending'
const loadKeywords = () => {
  const storedKeywords = localStorage.getItem('keyword')
  if (storedKeywords) {
    keywords.value = storedKeywords.split('+')
    originalKeywords.value = [...keywords.value] // Save the original keywords
  } else {
    keywords.value = ['trending']
    originalKeywords.value = ['trending'] // Save the original keywords
  }
}

// Function to remove a keyword
const removeKeyword = (keyword: string) => {
  keywords.value = keywords.value.filter(k => k !== keyword)
}

// Function to update keywords on the server and in localStorage
const updateKeywords = () => {
  const keywordToDelete = originalKeywords.value.filter(k => !keywords.value.includes(k))

  deleteUserProfileKeyword(keywordToDelete.join('+'))
}

// Mutation to update user profile keywords on the server

// Save the updated keywords
const saveKeywords = () => {
  updateKeywords()
}

// Cancel changes and restore original keywords
const cancelChanges = () => {
  keywords.value = [...originalKeywords.value]
}

onMounted(() => {
 loadKeywords()
})

</script>
