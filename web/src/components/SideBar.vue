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
              <svg
                class="flex-shrink-0 w-5 h-5 text-neutral-light transition duration-75 dark:text-gray-400 group-hover:neutral-light dark:group-hover:text-white"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                fill="currentColor"
                viewBox="0 0 18 18"
              >
                <path
                  d="M6.143 0H1.857A1.857 1.857 0 0 0 0 1.857v4.286C0 7.169.831 8 1.857 8h4.286A1.857 1.857 0 0 0 8 6.143V1.857A1.857 1.857 0 0 0 6.143 0Zm10 0h-4.286A1.857 1.857 0 0 0 10 1.857v4.286C10 7.169 10.831 8 11.857 8h4.286A1.857 1.857 0 0 0 18 6.143V1.857A1.857 1.857 0 0 0 16.143 0Zm-10 10H1.857A1.857 1.857 0 0 0 0 11.857v4.286C0 17.169.831 18 1.857 18h4.286A1.857 1.857 0 0 0 8 16.143v-4.286A1.857 1.857 0 0 0 6.143 10Zm10 0h-4.286A1.857 1.857 0 0 0 10 11.857v4.286c0 1.026.831 1.857 1.857 1.857h4.286A1.857 1.857 0 0 0 18 16.143v-4.286A1.857 1.857 0 0 0 16.143 10Z"
                ></path>
              </svg>
              <span class="ms-3">Personalized</span>
            </RouterLink>
          </li>
          <li>
            <RouterLink
              to="/search"
              class="flex items-center p-2 text-neutral-light rounded-lg dark:text-white hover:bg-gray-500 dark:hover:bg-gray-700 group"
            >
              <svg
                class="flex-shrink-0 w-5 h-5 text-neutral-light transition duration-75 dark:text-gray-400 group-hover:neutral-light dark:group-hover:text-white"
                width="20"
                height="20"
                viewBox="0 0 20 20"
              >
                <path
                  d="M14.386 14.386l4.0877 4.0877-4.0877-4.0877c-2.9418 2.9419-7.7115 2.9419-10.6533 0-2.9419-2.9418-2.9419-7.7115 0-10.6533 2.9418-2.9419 7.7115-2.9419 10.6533 0 2.9419 2.9418 2.9419 7.7115 0 10.6533z"
                  stroke="currentColor"
                  fill="none"
                  fill-rule="evenodd"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                ></path>
              </svg>
              <span class="ms-3">Search</span>
            </RouterLink>
          </li>
        </ul>
      </div>
    </aside>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { clickOutSide as vClickOutSide } from '@mahdikhashan/vue3-click-outside'

const sidebarOpen = ref(false)

function toggleSidebar() {
  sidebarOpen.value = !sidebarOpen.value
}

const handleClickOutside = () => {
  sidebarOpen.value = false
}
</script>
