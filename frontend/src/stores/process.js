import { defineStore } from 'pinia'
import axios from 'axios'
import { useUploadStore } from './upload'

export const useProcessStore = defineStore('process', {
  state: () => ({
    fileId: null,
    processProgress: 0,
    processStatus: null,
    totalLayer: null,
  }),
  actions: {
    setFileId() {
      const uploadStore = useUploadStore()
      if (uploadStore.uploadedId) {
        this.fileId = uploadStore.uploadedId
      }
    },

    async processFile() {
      this.setFileId()

      if (!this.fileId) {
        this.processStatus = 'Image with represent Id not found.'
        return
      }

      try {
        const response = await axios.post(`/api/process/${this.fileId}`)

        this.processStatus = response.data.message
        this.totalLayer = response.data.layers
      } catch (error) {
        this.processStatus = 'Something went wrong.'
      }
    },
  },
})
