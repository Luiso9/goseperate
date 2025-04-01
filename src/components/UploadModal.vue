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
  <div class="p-4 modal-content">
    <h3>Upload PNG Image</h3>
    <input type="file" @change="onFileChange" accept="image/png" />
    <p v-if="errorMessage">{{ errorMessage }}</p>
    <div v-if="selectedFile">
      <span>{{ selectedFile.name }}</span>
      <button @click="removeFile">Remove</button>
    </div>
    <button @click="upload" :disabled="imageStore.uploading">
      <span v-if="!imageStore.uploading">Upload</span>
      <span v-else>Uploading...</span>
    </button>
  </div>
</template>
