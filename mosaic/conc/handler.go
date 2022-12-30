// conc 即 concurrency，即使用并发
package conc

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
	"sync"
	"time"
)

// 并发版本 Mosaic 处理器实现
func Mosaic(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	r.ParseMultipartForm(10485760)
	file, _, _ := r.FormFile("image")
	defer file.Close()
	tileSize, _ := strconv.Atoi(r.FormValue("tile_size"))
	original, _, _ := image.Decode(file)
	bounds := original.Bounds()
	db := tilesdb.CloneTilesDB()
	// 调用 cut 方法切分原始图片为 4 等份（切分后的子图片还会各自进行马赛克处理，这些都会通过协程异步执行）
	c1 := cut(original, &db, tileSize, bounds.Min.X, bounds.Min.Y, bounds.Max.X/2, bounds.Max.Y/2)
	c2 := cut(original, &db, tileSize, bounds.Max.X/2, bounds.Min.Y, bounds.Max.X, bounds.Max.Y/2)
	c3 := cut(original, &db, tileSize, bounds.Min.X, bounds.Max.Y/2, bounds.Max.X/2, bounds.Max.Y)
	c4 := cut(original, &db, tileSize, bounds.Max.X/2, bounds.Max.Y/2, bounds.Max.X, bounds.Max.Y)
	// 将上述各自进行马赛克处理的子图片合并为最终的马赛克图片（通过协程异步执行）
	c := combine(bounds, c1, c2, c3, c4)
	buf1 := new(bytes.Buffer)
	jpeg.Encode(buf1, original, nil)
	originalStr := base64.StdEncoding.EncodeToString(buf1.Bytes())

	t1 := time.Now()
	// 由于切分后图片的马赛克处理后合并都是基于协程异步调用的，
	// 所以这里从 c 通道读取最终马赛克图片会阻塞，知道该通道写入值
	images := map[string]string{
		"original": originalStr,
		"mosaic":   <-c,
		"duration": fmt.Sprintf("%v ", t1.Sub(t0)),
	}
	t, _ := template.ParseFiles("views/results.html")
	t.Execute(w, images)
}

// 按照指定坐标切分原始图片
func cut(original image.Image, db *tilesdb.DB, tileSize, x1, y1, x2, y2 int) <-chan image.Image {
	c := make(chan image.Image)
	sp := image.Point{0, 0}
	// 由于对于每块切分后的图片会单独进行马赛克处理，所以为了提升处理性能，这里引入协程异步执行
	go func() {
		newimage := image.NewNRGBA(image.Rect(x1, y1, x2, y2))
		for y := y1; y < y2; y = y + tileSize {
			for x := x1; x < x2; x = x + tileSize {
				r, g, b, _ := original.At(x, y).RGBA()
				color := [3]float64{float64(r), float64(g), float64(b)}
				nearest := tilesdb.Nearest(color, db)
				file, err := os.Open(nearest)
				if err == nil {
					img, _, err := image.Decode(file)
					if err == nil {
						t := tilesdb.Resize(img, tileSize)
						tile := t.SubImage(t.Bounds())
						tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
						draw.Draw(newimage, tileBounds, tile, sp, draw.Src)
					} else {
						fmt.Println("error in decoding nearest", err, nearest)
					}
				} else {
					fmt.Println("error opening file when creating mosaic:", nearest)
				}
				file.Close()
			}
		}
		c <- newimage.SubImage(newimage.Rect)
	}()

	return c
}

// 合并经过马赛克处理的子图片并对最终马赛克图片进行 base64 编码
func combine(r image.Rectangle, c1, c2, c3, c4 <-chan image.Image) <-chan string {
	c := make(chan string)
	// 由于传入的每个子图片都是异步进行马赛克处理的，调用的时候可能尚未处理完成，所以这里也通过协程异步处理合并操作
	go func() {
		var wg sync.WaitGroup
		newimage := image.NewNRGBA(r)
		// 定义子图片嵌入最终目标图片的匿名函数，将其赋值给 copy
		copy := func(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
			draw.Draw(dst, r, src, sp, draw.Src)
			wg.Done()
		}
		wg.Add(4)
		var s1, s2, s3, s4 image.Image
		var ok1, ok2, ok3, ok4 bool
		for {
			// 通过选择通道读取每个子图片的异步马赛克处理结果，然后再通过协程异步执行子图片嵌入最终图片的 copy 函数
			select {
			case s1, ok1 = <-c1:
				go copy(newimage, s1.Bounds(), s1, image.Point{r.Min.X, r.Min.Y})
			case s2, ok2 = <-c2:
				go copy(newimage, s2.Bounds(), s2, image.Point{r.Max.X / 2, r.Min.Y})
			case s3, ok3 = <-c3:
				go copy(newimage, s3.Bounds(), s3, image.Point{r.Min.X, r.Max.Y / 2})
			case s4, ok4 = <-c4:
				go copy(newimage, s4.Bounds(), s4, image.Point{r.Max.X / 2, r.Max.Y / 2})
			}
			if ok1 && ok2 && ok3 && ok4 {
				break
			}
		}
		// 等待所有子图片都嵌入到目标图片（将此节点视为合并操作完成）
		wg.Wait()
		buf2 := new(bytes.Buffer)
		// 将合并后的最终马赛克图片进行 base64 编码并写入 c 通道返回
		jpeg.Encode(buf2, newimage, nil)
		c <- base64.StdEncoding.EncodeToString(buf2.Bytes())
	}()
	return c
}
