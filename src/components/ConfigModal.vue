<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useImageStore } from '@/stores/imageStore'

const imageStore = useImageStore()
const emit = defineEmits(['close'])
const imageId = ref(imageStore.imageId || '')
const colorsQuery = ref(localStorage.getItem('colorsQuery') || '4')
const isMobile = ref(window.innerWidth <= 768)

const updateScreenSize = () => {
  isMobile.value = window.innerWidth <= 768
}

const closeModal = () => {
  emit('close')
}

const process = async () => {
  if (!imageId.value) {
    alert('Please enter an image ID')
    return
  }

  localStorage.setItem('colorsQuery', colorsQuery.value) 
  await imageStore.processImage(imageId.value, colorsQuery.value)
  closeModal()
}

const handleOutsideClick = (event) => {
  if (!event.target.closest('.modal-content')) {
    closeModal()
  }
}

onMounted(() => {
  window.addEventListener('resize', updateScreenSize)
  window.addEventListener('click', handleOutsideClick)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateScreenSize)
  window.removeEventListener('click', handleOutsideClick)
  emit('close') // Ensure proper cleanup when unmounted
})

watch(
  () => imageStore.imageId,
  (newId) => {
    if (newId) imageId.value = newId
  },
)
</script>

<template>
  <!-- Content for Desktop Sidebar -->
  <div v-if="!isMobile" class="p-4 modal-content">
    <h3 class="text-lg font-bold">Image Processing</h3>
    <p class="text-sm">Configure and process the selected image.</p>

    <div class="mt-4 space-y-3">
      <div>
        <label for="image-id" class="block font-semibold my-2">Image ID</label>
        <input
          type="text"
          id="image-id"
          v-model="imageId"
          placeholder="Enter Image ID"
          class="input input-bordered w-full"
          :readonly="imageStore.imageId"
        />
      </div>

      <div>
        <label for="colors-query" class="block font-semibold my-2">Colors Query</label>
        <input
          type="text"
          id="colors-query"
          v-model="colorsQuery"
          placeholder="Enter Colors Query"
          class="input input-bordered w-full"
        />
      </div>
    </div>

    <button @click="process" :disabled="imageStore.processing" class="btn btn-primary mt-4 w-full">
      <span v-if="!imageStore.processing">Process</span>
      <span v-else class="loading loading-spinner"></span>
    </button>
  </div>

  <!-- Drawer for Mobile -->
  <div v-else class="fixed inset-0 bg-black/50 flex items-end sm:hidden">
    <div class="w-full bg-[#2e3440]/30 border-t border-gray-200/20 backdrop-blur-2xl shadow-lg ring-1 ring-black/5 isolate p-6 rounded-t-2xl modal-content">
      <h3 class="text-lg font-bold">Image Processing</h3>
      <p class="text-sm">Configure and process the selected image.</p>

      <div class="mt-4 space-y-3">
        <div>
          <label for="image-id" class="block font-semibold my-2">Image ID</label>
          <input
            type="text"
            id="image-id"
            v-model="imageId"
            placeholder="Enter Image ID"
            class="input input-bordered w-full"
            :readonly="imageStore.imageId"
          />
        </div>

        <div>
          <label for="colors-query" class="block font-semibold my-2">Colors Query</label>
          <input
            type="text"
            id="colors-query"
            v-model="colorsQuery"
            placeholder="Enter Colors Query"
            class="input input-bordered w-full"
          />
        </div>
      </div>

      <button
        @click="process"
        :disabled="imageStore.processing"
        class="btn btn-primary mt-4 w-full"
      >
        <span v-if="!imageStore.processing">Process</span>
        <span v-else class="loading loading-spinner"></span>
      </button>
    </div>
  </div>
</template>
