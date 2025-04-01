import { defineStore } from 'pinia'
import axios from 'axios'

export const useImageStore = defineStore('image', {
  state: () => ({
    uploading: false,
    processing: false,
    downloading: false,
    imageUrl: null,
    previewUrl: null,
    imageId: null,
    error: null,
  }),
  actions: {
    async uploadImage(fileOrUrl) {
      if (!fileOrUrl) {
        this.error = 'No file or URL provided'
        return
      }

      this.uploading = true
      this.error = null

      try {
        if (typeof fileOrUrl === 'string') {
          if (!fileOrUrl.startsWith('http') || !/^https?:\/\/[^\s/$.?#].[^\s]*$/.test(fileOrUrl)) {
            this.error = 'Invalid URL'
            return
          }
          const response = await axios.post(
            '/api/upload',
            { image_url: fileOrUrl },
            {
              headers: { 'Content-Type': 'application/json' },
            },
          )
          this.imageId = response.data.id
          this.imageUrl = `/api/uploads/${response.data.path}`
          this.previewUrl = `/api/preview/${this.imageId}?colors=4`
        } else if (fileOrUrl instanceof File) {
          if (!fileOrUrl.type.includes('image/png')) {
            this.error = 'Only PNG images are allowed'
            return
          }
          if (fileOrUrl.size > 5 * 1024 * 1024) {
            // 5MB limit
            this.error = 'File size exceeds the 5MB limit'
            return
          }
          const formData = new FormData()
          formData.append('image', fileOrUrl)
          const response = await axios.post('/api/upload', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
          })
          this.imageId = response.data.id
          this.imageUrl = `/api/uploads/${response.data.path}`
          this.previewUrl = `/api/preview/${this.imageId}?colors=4`
        } else {
          this.error = 'Invalid file format'
        }
      } catch (err) {
        this.error = 'Failed to upload image'
      } finally {
        this.uploading = false
      }
    },

    async processImage(id, colors = '4') {
      if (!id) {
        this.error = 'Invalid image ID'
        return
      }

      this.processing = true
      this.error = null

      try {
        await axios.post(`/api/process/${id}?colors=${colors}`)
        this.previewUrl = `/api/preview/${id}?colors=${colors}`
      } catch (err) {
        this.error = 'Failed to process image'
      } finally {
        this.processing = false
      }
    },

    async downloadZip(id, colors = '4') {
      if (!id) {
        this.error = 'Invalid image ID'
        return
      }

      this.downloading = true
      this.error = null

      try {
        const response = await axios.get(`/api/download/${id}?colors=${colors}`, {
          responseType: 'blob',
        })
        const url = window.URL.createObjectURL(new Blob([response.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', `processed_${id}.zip`)
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
      } catch (err) {
        this.error = 'Failed to download zip'
      } finally {
        this.downloading = false
      }
    },
  },
})
