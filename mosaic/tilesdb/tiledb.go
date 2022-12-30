package tilesdb

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
)

type DB map[string][3]float64

var TILESDB DB

// 获取整张图片的平均rgb值
func AverageColor(img image.Image) [3]float64 {
	bounds := img.Bounds()
	rsum, gsum, bsum := 0.0, 0.0, 0.0
	//遍历图片所有点，把每个点的rgb值累加起来
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rsum, gsum, bsum = rsum+float64(r), gsum+float64(g), bsum+float64(b)
		}
	}
	totalPixels := float64(bounds.Dx() * bounds.Dy())
	return [3]float64{rsum / totalPixels, gsum / totalPixels, bsum / totalPixels}
}

// 把图片缩放到指定的尺寸
func Resize(img image.Image, newWidth int) image.NRGBA {
	bounds := img.Bounds()
	ratio := bounds.Dx() / newWidth
	out := image.NewNRGBA(image.Rect(bounds.Min.X/ratio, bounds.Min.Y/ratio,
		bounds.Max.X/ratio, bounds.Max.Y/ratio))
	for y, j := bounds.Min.Y, bounds.Min.Y; y < bounds.Max.Y; y, j = y+ratio, j+1 {
		for x, i := bounds.Min.X, bounds.Min.X; x < bounds.Max.X; x, i = x+ratio, i+1 {
			r, g, b, a := img.At(x, y).RGBA()
			out.SetNRGBA(i, j, color.NRGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}

	return *out
}

// 读取并构建图片数据库
func TilesDB() DB {
	var dirname string = "tiles"
	fmt.Println("开始构建嵌入图片数据库...")
	db := make(DB)
	files, _ := os.ReadDir(dirname)
	for _, f := range files {
		name := "tiles/" + f.Name()
		file, err := os.Open(name)
		if err == nil {
			img, _, err := image.Decode(file)
			if err == nil {
				db[name] = AverageColor(img)
			} else {
				fmt.Println("构建嵌入图片数据库出错：", err, name)
			}
		} else {
			fmt.Println("构建嵌入图片数据库出错：", err, "无法打开文件", name)
		}
		file.Close()
	}
	fmt.Println("构建嵌入图片数据库完毕")
	return db
}

// 利用rgb值在素材库中查找最接近的tile
func Nearest(target [3]float64, db *DB) string {
	var filename string
	smallest := 10000000.0
	for k, v := range *db {
		dist := distance(target, v)
		if dist < smallest {
			filename, smallest = k, dist
		}
	}
	//INFO: 这里是否delete取决于是否允许出现重复的tile(素材库小的话就不用删了)
	//delete(*db, filename)
	return filename
}

func distance(p1 [3]float64, p2 [3]float64) float64 {
	return math.Sqrt(sq(p2[0]-p1[0]) + sq(p2[1]-p1[1]) + sq(p2[2]-p1[2]))
}

func sq(n float64) float64 {
	return n * n
}

// 复制TILESDB，因为读取文件过程很耗时
func CloneTilesDB() DB {
	db := make(DB)
	for k, v := range TILESDB {
		db[k] = v
	}
	return db
}
