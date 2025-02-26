package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"backend/services"
)

func PreviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imagePath := fmt.Sprintf("backend/uploads/%s", vars["filename"]) 
	numColors, err := strconv.Atoi(r.URL.Query().Get("color"))
	if err != nil || numColors <= 0 {
		http.Error(w, "Invalid color parameter", http.StatusBadRequest)
		return
	}

	previewData, err := services.GeneratePreview(imagePath, numColors)
	if err != nil {
		http.Error(w, "Error generating preview: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(previewData)
}
