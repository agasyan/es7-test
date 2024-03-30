package docgen

import (
	"crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type Document struct {
	ID               int     `json:"id"`
	ProductTitle     string  `json:"product_title"`
	ItemID           int     `json:"item_id"`
	CreatedUnix      int64   `json:"created_unix"`
	CreatedWeek      string  `json:"created_week"`
	CategoryID       int64   `json:"category_id"`
	CategoryName     string  `json:"category_name"`
	ProductCondition int     `json:"product_condition"`
	WidthImage       int     `json:"width_image"`
	HeightImage      int     `json:"height_image"`
	ImageURL         string  `json:"image"`
	Price            float64 `json:"price"`
	FinalPrice       float64 `json:"final_price"`
	DiscountPercent  float64 `json:"discount_percentage"`
	ScoreDoc         float64 `json:"score_doc"`
	CTR              float64 `json:"ctr"`
	PriceBid         int     `json:"price_bid"`
	IsCashback       bool    `json:"is_cashback"`
	IsFreeshipping   bool    `json:"is_freeshipping"`
	IsLocal          bool    `json:"is_local"`
	MinimumOrder     int     `json:"minimum_order"`
	ProvinceID       int     `json:"province_id"`
	ProvinceName     string  `json:"province_name"`
	CityID           int     `json:"city_id"`
	CityName         string  `json:"city_name"`
	DistrictID       int     `json:"district_id"`
	DistrictName     string  `json:"district_name"`
	Latitude         float64 `json:"lat"`
	Longitude        float64 `json:"long"`
	ShopName         string  `json:"shop_name"`
	ShopID           int64   `json:"shop_id"`
	ShopTier         int     `json:"shop_tier"`
	Rating           float64 `json:"rating"`
	Review           int     `json:"review"`
	WeightGram       int     `json:"weight_gram"`
}

func (d *Docgen) generate() Document {
	var id int
	id = gofakeit.IntRange(1, math.MaxInt)
	d.mutex.Lock()
	// check until not exists
	for d.mapID[id] {
		id = gofakeit.IntRange(1, math.MaxInt)
	}
	d.mapID[id] = true
	d.mutex.Unlock()
	wImg := gofakeit.Number(200, 1800)
	hImg := gofakeit.Number(200, 1800)
	cAt := gofakeit.DateRange(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Now())
	y, m := cAt.ISOWeek()
	p := d.GetRandomProduct()
	s := d.arrShops[gofakeit.IntRange(0, len(d.arrShops))]
	price := gofakeit.Number(1000, 10000000)
	disc := gofakeit.IntRange(0, 90)
	fp := float64(price) * float64(100-disc) / 100
	ctr := gofakeit.Float64Range(0, 5)
	bid := gofakeit.Number(400, 10000)
	return Document{
		ID:               id,
		ProductTitle:     fmt.Sprintf("%s %s", p.Name, generateRandomString(16)),
		ItemID:           id,
		CreatedUnix:      cAt.Unix(),
		CreatedWeek:      fmt.Sprintf("%d%d", y, m),
		CategoryID:       p.CategoriesID,
		CategoryName:     p.CategoriesName,
		ProductCondition: gofakeit.RandomInt([]int{0, 1}),
		WidthImage:       wImg,
		HeightImage:      hImg,
		ImageURL:         gofakeit.ImageURL(wImg, hImg),
		Price:            float64(price),
		FinalPrice:       fp,
		DiscountPercent:  float64(disc),
		ScoreDoc:         (ctr * float64(bid)) + 250,
		CTR:              gofakeit.Float64Range(0, 5),
		PriceBid:         bid,
		IsCashback:       gofakeit.Bool(),
		IsFreeshipping:   gofakeit.Bool(),
		IsLocal:          gofakeit.Bool(),
		MinimumOrder:     gofakeit.IntRange(1, 3),
		ProvinceID:       s.Loc.ProvinceID,
		ProvinceName:     s.Loc.ProvinceName,
		CityID:           s.Loc.CityID,
		CityName:         s.Loc.CityName,
		DistrictID:       s.Loc.DistrictID,
		DistrictName:     s.Loc.DistrictName,
		Latitude:         s.Loc.Latitude,
		Longitude:        s.Loc.Longitude,
		ShopName:         s.Name,
		ShopID:           int64(s.ID),
		ShopTier:         gofakeit.IntRange(0, 4),
		Rating:           gofakeit.Float64Range(3, 5),
		Review:           gofakeit.IntRange(0, 1000),
		WeightGram:       gofakeit.IntRange(100, 10000),
	}
}

func (d *Docgen) BulkGenerate(count int) []Document {
	arrDocID := make([]Document, 0, count)
	for i := 0; i < count; i++ {
		arrDocID = append(arrDocID, d.generate())
	}
	return arrDocID
}

func (d *Docgen) UpdateArr() {

	arrDocID := make([]int, 0)

	d.mutex.Lock()
	for k := range d.mapID {
		arrDocID = append(d.arrDocID, k)
	}
	d.arrDocID = arrDocID
	d.mutex.Unlock()
}

func (d *Docgen) GetExistKey(action string) int {
	d.mutex.Lock()
	if len(d.arrDocID) == 0 {
		return 0
	}
	i := gofakeit.IntRange(0, len(d.arrDocID)-1)
	val := d.arrDocID[i]
	d.mutex.Unlock()

	if action == "DELETE" {
		d.mutex.Lock()
		delete(d.mapID, val)
		d.mutex.Unlock()

		d.UpdateArr()
	}

	return val
}

type Docgen struct {
	productMap map[string]Product
	mapID      map[int]bool
	mutex      *sync.Mutex
	arrDocID   []int
	arrShops   []Shop
}

func Init() *Docgen {
	mapID := make(map[int]bool, 0)
	productArr, shopArr := readFromCSVFiles()
	pMap := make(map[string]Product, len(productArr))
	for _, p := range productArr {
		pMap[p.Name] = p
	}

	return &Docgen{
		productMap: pMap,
		mapID:      mapID,
		mutex:      &sync.Mutex{},
		arrDocID:   make([]int, 0),
		arrShops:   shopArr,
	}
}

func (d *Docgen) InitMapArr(ids []int) {
	mapID := make(map[int]bool, len(ids))
	for _, id := range ids {
		mapID[id] = true
	}
	d.mutex.Lock()
	d.mapID = mapID
	d.arrDocID = ids
	d.mutex.Unlock()
}

func (d *Docgen) GetRandomProduct() Product {
	pName := gofakeit.RandomMapKey(d.productMap)
	return d.productMap[pName.(string)]
}

func generateRandomString(length int) string {
	bytes := make([]byte, length/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)[:length]
}

// Product struct represents the structure of each product
type Product struct {
	Name           string
	CategoriesName string
	CategoriesID   int64
}

type Location struct {
	DistrictID   int
	DistrictName string
	Latitude     float64
	Longitude    float64
	CityID       int
	CityName     string
	ProvinceID   int
	ProvinceName string
}

type Shop struct {
	ID   int
	Name string
	Loc  Location
}

func readFromCSVFiles() ([]Product, []Shop) {
	// Open the CSV file
	// path, err := os.Getwd()
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(path) // for example /home/user
	file, err := os.Open("../../files/name.csv")
	if err != nil {
		log.Fatalln("Error opening the file:", err)
	}
	defer file.Close()

	// Create a CSV reader with comma as the field delimiter
	reader := csv.NewReader(file)
	reader.Comma = ';'
	// Read all records from CSV
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Error reading the CSV:", err)
	}

	// Initialize slices
	var products []Product

	// Parse CSV records into Product struct and populate slices
	for i, record := range records {
		if i == 1 {
			continue
		}
		cID, _ := strconv.ParseInt(record[0], 10, 64)
		product := Product{
			Name:           record[2], // Assuming the name is in the first column
			CategoriesName: record[1], // Assuming the category name is in the second column
			CategoriesID:   cID,
		}

		// Append to slices
		products = append(products, product)

	}

	// Open the CSV file
	fileLoc, err := os.Open("../../files/loc.csv")
	if err != nil {
		log.Fatalln("Error opening the file:", err)
	}
	defer fileLoc.Close()

	// Create a CSV reader with comma as the field delimiter
	readerLoc := csv.NewReader(fileLoc)
	readerLoc.Comma = ','
	// Read all records from CSV
	recordsLoc, err := readerLoc.ReadAll()
	if err != nil {
		log.Fatalln("Error reading the CSV:", err)
	}

	// Initialize slices
	locMap := make(map[int]Location)

	// Parse CSV records into Product struct and populate slices
	for _, record := range recordsLoc {
		districtID, _ := strconv.Atoi(record[0])
		latitude, _ := strconv.ParseFloat(record[2], 64)
		longitude, _ := strconv.ParseFloat(record[3], 64)
		cityID, _ := strconv.Atoi(record[4])
		provinceID, _ := strconv.Atoi(record[6])

		location := Location{
			DistrictID:   districtID,
			DistrictName: record[1],
			Latitude:     latitude,
			Longitude:    longitude,
			CityID:       cityID,
			CityName:     record[5],
			ProvinceID:   provinceID,
			ProvinceName: record[7],
		}

		// Append location to slice
		locMap[districtID] = location

	}

	// Open the CSV file
	fileShop, err := os.Open("../../files/shops.csv")
	if err != nil {
		log.Fatalln("Error opening the file:", err)
	}
	defer fileShop.Close()

	// Create a CSV reader with comma as the field delimiter
	readerShop := csv.NewReader(fileShop)
	readerShop.Comma = ','
	// Read all records from CSV
	recordsShop, err := readerShop.ReadAll()
	if err != nil {
		log.Fatalln("Error reading the CSV:", err)
	}

	// Initialize slices
	var shops []Shop

	// Parse CSV records into Product struct and populate slices
	for _, record := range recordsShop {
		id, _ := strconv.Atoi(record[0])

		districtID, _ := strconv.Atoi(record[2])
		product := Shop{
			ID:   id,
			Name: record[1],
			Loc:  locMap[districtID],
		}

		// Append to slices
		shops = append(shops, product)
	}

	return products, shops
}
