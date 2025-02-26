package services

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"

	"gonum.org/v1/gonum/mat"
	"gocv.io/x/gocv"
)

func ExtractColors(imagePath, outputDir string, k int) error {
	img := gocv.IMRead(imagePath, gocv.IMReadColor)
	if img.Empty() {
		return fmt.Errorf("could not read image")
	}
	defer img.Close()

	rows, cols := img.Rows(), img.Cols()
	data := make([]float64, 0, rows*cols*3)

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			pixel := img.GetVecbAt(y, x)
			data = append(data, float64(pixel[2]), float64(pixel[1]), float64(pixel[0]))
		}
	}

	points := mat.NewDense(len(data)/3, 3, data)
	centroids, labels := KMeans(points, k)

	os.MkdirAll(outputDir, os.ModePerm)

	for i, center := range centroids {
		layer := image.NewRGBA(image.Rect(0, 0, cols, rows))
		hexColor := fmt.Sprintf("%02X%02X%02X", uint8(center[0]), uint8(center[1]), uint8(center[2]))

		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				index := y*cols + x
				if index < len(labels) && labels[index] == i {
					layer.Set(x, y, color.RGBA{
						R: uint8(center[0]), G: uint8(center[1]), B: uint8(center[2]), A: 255,
					})
				}
			}
		}

		filePath := filepath.Join(outputDir, fmt.Sprintf("%s.png", hexColor))
		file, _ := os.Create(filePath)
		png.Encode(file, layer)
		file.Close()
	}

	return nil
}

func KMeans(data *mat.Dense, k int) ([][]float64, []int) {
	rows, _ := data.Dims()
	labels := make([]int, rows)
	centroids := make([][]float64, k)

	for i := range centroids {
		idx := rand.Intn(rows)
		centroids[i] = mat.Row(nil, idx, data)
	}

	for iter := 0; iter < 10; iter++ {
		for i := 0; i < rows; i++ {
			sample := mat.Row(nil, i, data)
			labels[i] = nearestCentroid(sample, centroids)
		}

		newCenter := mat.NewDense(k, 3, nil)
		counts := make([]int, k)

		for i := 0; i < rows; i++ {
			label := labels[i]
			pixel := mat.Row(nil, i, data)

			for j := 0; j < 3; j++ {
				newCenter.Set(label, j, newCenter.At(label, j)+pixel[j])
			}
			counts[label]++
		}

		for i := 0; i < k; i++ {
			if counts[i] > 0 {
				for j := 0; j < 3; j++ {
					newCenter.Set(i, j, newCenter.At(i, j)/float64(counts[i]))
				}
			}
		}

		for i := range centroids {
			centroids[i] = mat.Row(nil, i, newCenter)
		}
	}

	return centroids, labels
}

func nearestCentroid(sample []float64, centroids [][]float64) int {
	bestIdx := 0
	bestDist := euclideanDistance(sample, centroids[0])

	for i, center := range centroids[1:] {
		dist := euclideanDistance(sample, center)
		if dist < bestDist {
			bestIdx = i + 1
			bestDist = dist
		}
	}

	return bestIdx
}

func euclideanDistance(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return sum
}
