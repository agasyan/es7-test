package docgen

import (
	"encoding/hex"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type Document struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Genre       string `json:"genre"`
	WidthImage  int    `json:"width_image"`
	HeightImage int    `json:"height_image"`
	ImageURL    string `json:"image"`
	CreatedUnix int64  `json:"created_unix"`
	Price       int    `json:"price"`
}

func (d *Docgen) generate() Document {
	var id int
	id = gofakeit.IntRange(1, math.MaxInt32)
	d.mutex.Lock()
	for d.mapID[id] {
		id = gofakeit.IntRange(1, math.MaxInt)
	}
	d.mapID[id] = true
	d.mutex.Unlock()
	wImg := gofakeit.Number(200, 1800)
	hImg := gofakeit.Number(200, 1800)
	cAt := gofakeit.DateRange(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Now())
	return Document{
		ID:          id,
		Title:       fmt.Sprintf("%s %s", gofakeit.BookTitle(), generateRandomString(15)),
		Author:      fmt.Sprintf("%s %s", gofakeit.FirstName(), gofakeit.LastName()),
		Genre:       gofakeit.BookGenre(),
		WidthImage:  wImg,
		HeightImage: hImg,
		ImageURL:    gofakeit.ImageURL(wImg, hImg),
		CreatedUnix: cAt.Unix(),
		Price:       gofakeit.Number(1000, 10000000),
	}
}

func (d *Docgen) BulkGenerate(count int) []Document {
	arr := make([]Document, 0, count)
	for i := 0; i < count; i++ {
		arr = append(arr, d.generate())
	}
	return arr
}

func (d *Docgen) UpdateArr() {

	arr := make([]int, 0)

	for k := range d.mapID {
		arr = append(d.arr, k)
	}
	d.mutex.Lock()
	d.arr = arr
	d.mutex.Unlock()
}

func (d *Docgen) GetExistKey(action string) int {
	if len(d.arr) == 0 {
		return 0
	}
	i := gofakeit.IntRange(0, len(d.arr))

	d.mutex.Lock()
	val := d.arr[i]
	d.mutex.Unlock()

	if action == "DELETE" {
		d.mutex.Lock()
		delete(d.mapID, val)
		d.mutex.Unlock()
		d.UpdateArr()
	}

	return d.arr[i]
}

type Docgen struct {
	mapID map[int]bool
	mutex *sync.Mutex
	arr   []int
}

func Init() *Docgen {
	mapID := make(map[int]bool, 0)
	return &Docgen{
		mapID: mapID,
		mutex: &sync.Mutex{},
		arr:   make([]int, 0),
	}
}

func generateRandomString(length int) string {
	bytes := make([]byte, length/2)
	return hex.EncodeToString(bytes)[:length]
}
