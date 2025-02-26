<template>
  <main>
    <div class="input_file">
      <input type="file" name="image" id="upload_image" @change="handleFileChange" />
    </div>
    <div class="image_link">
      <input type="text" />
    </div>
    <button @click="upload" :disabled="!uploadStore.file">Submit</button>
    <button @click="processImage" :disabled="!uploadStore.file">Process</button>

    <p v-if="uploadStore.uploadedId">File ID: {{ uploadStore.uploadedId }}</p>

    <img :src="uploadStore.uploadedImage" alt="Uploaded image">

    <p v-if="uploadStore.uploadProgress > 0">Progress: {{ uploadStore.uploadProgress }}%</p>
    <p v-if="processStore.processStats"> Color: {{ processStore.totalLayer }}</p>
    <p v-if="uploadStore.uploadStatus">{{ uploadStore.uploadStatus }}</p>
  </main>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUploadStore } from '@/stores/upload'
import { useProcessStore } from '@/stores/process'

const uploadStore = useUploadStore()
const fileInput = ref(null)
const processStore = useProcessStore()

const handleFileChange = (event) => {
  const file = event.target.files[0]
  if (file) {
    uploadStore.setFile(file)
  }
}

onMounted(() => {
  processStore.setFileId()
})

const upload = () => {
  uploadStore.uploadFile()
}

const processImage = async () => {
  await processStore.processFile()
}
</script>
