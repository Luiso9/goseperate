import { defineStore } from 'pinia';
import axios from 'axios';

export const useUploadStore = defineStore('upload', {
  state: () => ({
    file: null,
    uploadProgress: 0,
    uploadStatus: null,
    uploadedId: null,
    uploadedImage: null,
  }),

  actions: {
    setFile(newFile) {
      this.file = newFile;
    },

    async uploadFile() {
      if (!this.file) {
        this.uploadStatus = 'No file selected';
        return;
      }

      const formData = new FormData();
      formData.append('image', this.file);

      try {
        const response = await axios.post('/api/upload', formData, {
          headers: { 'Content-Type': 'multipart/form-data' },
          onUploadProgress: (progressEvent) => {
            this.uploadProgress = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          }
        });

        this.uploadedImage = '/api/' + response.data.path;
        this.uploadStatus = response.data.message;
        this.uploadedId = response.data.id;
        console.log(this.uploadedId)

      } catch (error) {
        this.uploadStatus = response.data.message;
      }
    }
  }
});
