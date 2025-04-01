<script setup>
import { ref, defineEmits, onMounted, onUnmounted } from 'vue'
import { useImageStore } from '@/stores/imageStore'

const emit = defineEmits(['close'])
const imageStore = useImageStore()

const selectedFile = ref(null)
const errorMessage = ref(null)

const onFileChange = (event) => {
  const file = event.target.files[0]
  if (file) {
    if (!file.type.includes('image/png')) {
      errorMessage.value = 'Only PNG images are allowed!'
      selectedFile.value = null
    } else if (file.size > 5 * 1024 * 1024) {
      // 5MB limit
      errorMessage.value = 'File size exceeds the 5MB limit!'
      selectedFile.value = null
    } else {
      selectedFile.value = file
      errorMessage.value = null
    }
  }
}

const removeFile = () => {
  selectedFile.value = null
  errorMessage.value = null
}

const closeModal = () => {
  emit('close')
}

const upload = async () => {
  if (!selectedFile.value) {
    errorMessage.value = 'Please select a PNG image to upload.'
    return
  }

  errorMessage.value = null
  await imageStore.uploadImage(selectedFile.value)

  if (imageStore.error) {
    errorMessage.value = imageStore.error
  } else {
    selectedFile.value = null
    closeModal()
  }
}

onMounted(() => {
  window.addEventListener('click', (event) => {
    if (!event.target.closest('.modal-content')) {
      closeModal()
    }
  })
})

onUnmounted(() => {
  window.removeEventListener('click', closeModal)
})
</script>

<template>
  <!-- Content for Desktop Sidebar -->
  <div v-if="!isMobile" class="p-4 modal-content">
    <h3 class="text-lg font-bold">Upload PNG Image</h3>
    <p class="text-sm">Select a PNG image to upload.</p>
    <div class="mt-4">
      <input
        type="file"
        @change="onFileChange"
        accept="image/png"
        class="file-input file-input-bordered w-full active:border-none active:outline-none"
      />
    </div>
    <p v-if="errorMessage" class="text-red-500 mt-2">{{ errorMessage }}</p>
    <div v-if="selectedFile" class="mt-2 flex items-center gap-2">
      <span class="text-sm">{{
        selectedFile.name.length > 20 ? selectedFile.name.slice(0, 17) + '...' : selectedFile.name
      }}</span>
      <button @click="removeFile" class="btn btn-xs btn-ghost"><CloseIcon /></button>
    </div>
    <button @click="upload" :disabled="imageStore.uploading" class="btn btn-primary mt-4 w-full">
      <span v-if="!imageStore.uploading">Upload</span>
      <span v-else class="loading loading-spinner"></span>
    </button>
    <!-- Removed Close Button -->
  </div>

  <!-- Drawer for Mobile -->
  <div v-else class="fixed inset-0 bg-black/50 flex items-end">
    <div
      class="w-full bg-[#2e3440]/30 border-t border-gray-200/20 backdrop-blur-2xl shadow-lg ring-1 ring-black/5 isolate p-6 rounded-t-2xl modal-content"
    >
      <h3 class="text-lg font-bold">Upload PNG Image</h3>
      <p class="text-sm">Select a PNG image to upload.</p>
      <div class="mt-4">
        <input
          type="file"
          @change="onFileChange"
          accept="image/png"
          class="file-input file-input-bordered w-full active:border-none active:outline-none"
        />
      </div>
      <p v-if="errorMessage" class="text-red-500 mt-2">{{ errorMessage }}</p>
      <div v-if="selectedFile" class="mt-2 flex items-center gap-2">
        <span class="text-sm">{{ selectedFile.name }}</span>
        <button @click="removeFile" class="btn btn-xs btn-ghost"><CloseIcon /></button>
      </div>
      <button @click="upload" :disabled="imageStore.uploading" class="btn btn-primary mt-4 w-full">
        <span v-if="!imageStore.uploading">Upload</span>
        <span v-else class="loading loading-spinner"></span>
      </button>
      <!-- Removed Close Button -->
    </div>
  </div>
</template>
