<script setup>
import { ref, shallowRef, markRaw, defineAsyncComponent, onMounted, onUnmounted } from 'vue'
import MenuIcon from 'vue-material-design-icons/FolderWrench.vue'
import UploadCircleOutline from 'vue-material-design-icons/CloudUpload.vue'
import ExportIcon from 'vue-material-design-icons/Export.vue'

import { VTooltip } from 'floating-vue'
import 'floating-vue/dist/style.css'
import ImageDisplay from './components/ImageDisplay.vue'

const modalComponent = shallowRef(null)

const modals = {
  upload: defineAsyncComponent(() => import('./components/UploadModal.vue')),
  config: defineAsyncComponent(() => import('./components/ConfigModal.vue')),
  export: defineAsyncComponent(() => import('./components/ExportModal.vue')),
}

let isModalTransitioning = false

const openModal = (type) => {
  if (isModalTransitioning) return
  isModalTransitioning = true
  modalComponent.value = markRaw(modals[type])
  setTimeout(() => {
    isModalTransitioning = false
  }, 300) // Adjust timeout to match modal transition duration
}

const closeModal = () => {
  if (isModalTransitioning) return
  isModalTransitioning = true
  modalComponent.value = null
  setTimeout(() => {
    isModalTransitioning = false
  }, 300) // Adjust timeout to match modal transition duration
}

const isMobile = ref(window.innerWidth <= 768)

const updateScreenSize = () => {
  isMobile.value = window.innerWidth <= 768
}

onMounted(() => {
  window.addEventListener('resize', updateScreenSize)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateScreenSize)
})
</script>

<template>
  <div class="flex flex-col min-h-screen">
    <aside
      class="fixed top-0 left-0 h-full bg-[#2e3440]/30 border-r border-gray-200/20 shadow-lg ring-1 ring-black/5 isolate w-18 md:flex lg:flex flex-col items-center py-4 rounded-tr-2xl rounded-br-2xl hidden"
    >
      <nav>
        <ul class="flex flex-col items-center space-y-4 px-4">
          <li>
            <button @click="openModal('upload')" class="btn btn-ghost" v-tooltip="'Upload'">
              <UploadCircleOutline />
            </button>
          </li>
          <li>
            <button @click="openModal('config')" class="btn btn-ghost" v-tooltip="'Config'">
              <MenuIcon />
            </button>
          </li>
          <li>
            <button @click="openModal('export')" class="btn btn-ghost" v-tooltip="'Export'">
              <ExportIcon />
            </button>
          </li>
        </ul>
      </nav>
    </aside>

    <!-- Mobile Navigation -->
    <div
      class="fixed bottom-0 left-0 w-full border-t rounded-t-2xl sm:hidden bg-[#2e3440]/30 border-r border-gray-200/20 shadow-lg ring-1 ring-black/5 isolate"
    >
      <div class="flex justify-around p-2">
        <button @click="openModal('upload')" class="btn btn-ghost" v-tooltip="'Upload'">
          <UploadCircleOutline />
        </button>
        <button @click="openModal('config')" class="btn btn-ghost" v-tooltip="'Config'">
          <MenuIcon />
        </button>
        <button @click="openModal('export')" class="btn btn-ghost" v-tooltip="'Export'">
          <ExportIcon />
        </button>
      </div>
    </div>

    <!-- Sidebar for Desktop -->
    <aside
      v-if="!isMobile && modalComponent"
      class="fixed top-0 right-0 h-full py-2 border-l w-64 rounded-tl-2xl rounded-bl-2xl bg-[#2e3440]/30 border-r border-gray-200/20 shadow-lg ring-1 ring-black/5 isolate"
    >
      <component :is="modalComponent" @close="closeModal" />
    </aside>

    <!-- Mobile Modal -->
    <Suspense v-if="isMobile && modalComponent">
      <component :is="modalComponent" @close="closeModal" />
      <template #fallback>
        <div
          class="fixed inset-0 flex items-center justify-center border-t bg-[#2e3440]/30 border-r border-gray-200/20 shadow-lg ring-1 ring-black/5 isolate"
        >
          <span class="loading loading-spinner"></span>
        </div>
      </template>
    </Suspense>

    <div class="flex-1 overflow-hidden">
      <ImageDisplay />
    </div>
  </div>
</template>
