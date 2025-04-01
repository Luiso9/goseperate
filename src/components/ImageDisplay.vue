<script setup>
import { useImageStore } from '@/stores/imageStore'
import { computed } from 'vue'

const imageStore = useImageStore()

const imagePreviewUrl = computed(() => imageStore.previewUrl)
const isLoading = computed(() => imageStore.loading)
const error = computed(() => imageStore.error)
</script>

<template>
  <div class="container mx-auto my-8">
    <div v-if="isLoading" class="text-center">
      <p class="text-gray-100">Loading...</p>
    </div>
    <div v-else-if="error" class="text-center">
      <p class="text-red-500">Failed to load image.</p>
    </div>
    <div v-else-if="imagePreviewUrl" class="flex justify-center items-center">
      <img
        :src="imagePreviewUrl"
        alt="Image Preview"
        class="rounded-lg object-contain w-9/12 max-h-[calc(100vh-4rem)]"
        style="padding: 2rem 0"
      />
    </div>
    <div v-else class="text-center">
      <p class="text-gray-100">No image selected</p>
    </div>
  </div>
</template>
