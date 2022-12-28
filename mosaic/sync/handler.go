package sync

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/draw"
	"image/jpeg"
	"mosaic/tilesdb"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Mosaic(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	// 从 POST 表单获取用户上传内容
	r.ParseMultipartForm((2 << 20) * 10) // 最大上传内容大小是 10 MB
	// 通过 image 字段读取用户上传图片
	file, _, _ := r.FormFile("image")
	defer file.Close()
	// 通过 tile_size 字段读取用户设置的区块尺寸
	tileSize, _ := strconv.Atoi(r.FormValue("tile_size"))
	// 解码获取原始图片
	original, _, _ := image.Decode(file)
	bounds := original.Bounds()
	// 克隆目标图片用于绘制马赛克图片
	newimage := image.NewNRGBA(image.Rect(bounds.Min.X, bounds.Min.X, bounds.Max.X, bounds.Max.Y))
	// 克隆嵌入图片数据库
	db := tilesdb.CloneTilesDB()
	// 从 (0, 0) 坐标开始将目标图片按照 tile_size 值均分成多个小区块
	sp := image.Point{0, 0}
	for y := bounds.Min.Y; y < bounds.Max.Y; y = y + tileSize {
		for x := bounds.Min.X; x < bounds.Max.X; x = x + tileSize {
			// 使用图片区块左上角像素颜色作为该区块的平均颜色
			r, g, b, _ := original.At(x, y).RGBA()
			color := [3]float64{float64(r), float64(g), float64(b)}
			// 从嵌入图片数据库获取平均颜色与之最接近的嵌入图片
			nearest := tilesdb.Nearest(color, &db)
			file, err := os.Open(nearest)
			if err == nil {
				img, _, err := image.Decode(file)
				if err == nil {
					// 将嵌入图片调整到当前区块大小
					t := tilesdb.Resize(img, tileSize)
					tile := t.SubImage(t.Bounds())
					tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
					// 然后将调整大小后的嵌入图片绘制到当前区块位置，从而真正嵌入到马赛克图片中
					draw.Draw(newimage, tileBounds, tile, sp, draw.Src)
				} else {
					fmt.Println("error:", err, nearest)
				}
			} else {
				fmt.Println("error:", nearest)
			}
			file.Close()
		}
	}

	buf1 := new(bytes.Buffer)
	jpeg.Encode(buf1, original, nil)
	// 原始图片 base64 值
	originalStr := base64.StdEncoding.EncodeToString(buf1.Bytes())

	buf2 := new(bytes.Buffer)
	jpeg.Encode(buf2, newimage, nil)
	// 马赛克图片 base64 值
	mosaic := base64.StdEncoding.EncodeToString(buf2.Bytes())

	t1 := time.Now()
	// 构建一个包含原始图片、马赛克图片和处理时间的字典类型变量 images
	images := map[string]string{
		"original": originalStr,
		"mosaic":   mosaic,
		"duration": fmt.Sprintf("%v ", t1.Sub(t0)),
	}
	// 将 images 值渲染到响应 HTML 文件返回给用户
	t, _ := template.ParseFiles("views/results.html")
	t.Execute(w, images)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/upload.html")
	t.Execute(w, nil)
}
