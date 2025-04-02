<script setup>
import { useImageStore } from '@/stores/imageStore'
import { computed } from 'vue'

const imageStore = useImageStore()

const imagePreviewUrl = computed(() => imageStore.previewUrl)
const isLoading = computed(() => imageStore.loading || imageStore.processing)
const error = computed(() => imageStore.error)
</script>

<template>
  <div class="container mx-auto my-8">
    <div v-if="isLoading" class="flex justify-center items-center">
      <div class="skeleton opacity-30 w-96 h-96"></div>
    </div>
    <div v-else-if="error" class="text-center">
      <p class="text-red-500">Failed to load image.</p>
    </div>
    <div v-else-if="imagePreviewUrl" class="flex justify-center items-center">
      <img
        :src="imagePreviewUrl"
        alt="Image Preview"
        :class="{ 'blur-sm': isLoading }"
        class="rounded-lg object-contain w-9/12 max-h-[calc(100vh-4rem)]"
        style="padding: 2rem 0"
      />
    </div>
    <div v-else class="flex justify-center items-center">
      <div class="skeleton opacity-30 w-96 h-96 flex justify-center items-center">
        <p class="text-gray-500">No image selected</p>
      </div>
    </div>
  </div>
</template>
