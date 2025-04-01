<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useImageStore } from '@/stores/imageStore'

const imageStore = useImageStore()
const emit = defineEmits(['close'])

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

const downloadId = ref(imageStore.imageId || '')

watch(
  () => imageStore.imageId,
  (newId) => {
    if (newId) downloadId.value = newId
  },
)

const closeModal = () => {
  emit('close')
}

const download = async () => {
  if (!downloadId.value) {
    alert('Please enter an image ID')
    return
  }

  await imageStore.downloadZip(downloadId.value)
  closeModal()
}

const handleOutsideClick = (event) => {
  if (!event.target.closest('.modal-content')) {
    closeModal()
  }
}

onMounted(() => {
  window.addEventListener('click', handleOutsideClick)
})

onUnmounted(() => {
  window.removeEventListener('click', handleOutsideClick)
  emit('close') // Ensure proper cleanup when unmounted
})
</script>

<template>
  <!-- Content for Desktop Sidebar -->
  <div v-if="!isMobile" class="p-4 modal-content">
    <h3 class="text-lg font-bold">Download Processed Image</h3>
    <p class="text-sm">Download your processed image as a zip file.</p>

    <div class="mt-4 space-y-3">
      <!-- Image ID -->
      <div>
        <label for="download-id" class="font-semibold hidden">Image ID</label>
        <input
          type="text"
          id="download-id"
          v-if="!imageStore.imageId"
          v-model="downloadId"
          placeholder="Enter Image ID"
          class="input input-bordered w-full hidden"
        />
        <input
          type="text"
          id="download-id-readonly"
          v-else
          :value="imageStore.imageId"
          readonly
          class="input input-bordered w-full bg-gray-100 hidden"
        />
      </div>
    </div>

    <!-- Download Button -->
    <button
      @click="download"
      :disabled="imageStore.downloading"
      class="btn btn-primary mt-4 w-full"
    >
      <span v-if="!imageStore.downloading">Download</span>
      <span v-else class="loading loading-spinner"></span>
    </button>
    <!-- Removed Close Button -->
  </div>

  <!-- Drawer for Mobile -->
  <div v-else class="fixed inset-0 bg-black/50 flex items-end">
    <div
      class="w-full p-6 rounded-t-2xl bg-[#2e3440]/30 border-t border-gray-200/20 backdrop-blur-2xl shadow-lg ring-1 ring-black/5 isolate modal-content"
    >
      <h3 class="text-lg font-bold">Download Processed Image</h3>
      <p>Download your processed image as a zip file.</p>

      <div class="mt-4 space-y-3">
        <!-- Image ID -->
        <div>
          <label for="download-id" class="font-semibold hidden">Image ID</label>
          <input
            type="text"
            id="download-id"
            v-if="!imageStore.imageId"
            v-model="downloadId"
            placeholder="Enter Image ID"
            class="input input-bordered w-full hidden"
          />
          <input
            type="text"
            id="download-id-readonly"
            v-else
            :value="imageStore.imageId"
            readonly
            class="input input-bordered w-full bg-gray-100 hidden"
          />
        </div>
      </div>

      <!-- Download Button -->
      <button
        @click="download"
        :disabled="imageStore.downloading"
        class="btn btn-primary mt-4 w-full"
      >
        <span v-if="!imageStore.downloading">Download</span>
        <span v-else class="loading loading-spinner"></span>
      </button>
      <!-- Removed Close Button -->
    </div>
  </div>
</template>
