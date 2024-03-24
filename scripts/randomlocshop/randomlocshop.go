package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/brianvoe/gofakeit/v6"
)

type Shop struct {
	ID         int
	Name       string
	DistrictID int
}

type Province struct {
	ProvinceID   int
	ProvinceName string
}

type City struct {
	CityID     int
	CityName   string
	ProvinceID int
}

type District struct {
	CityID       int
	DistrictID   int
	DistrictName string
	Latitude     float64
	Longitude    float64
}

func main() {
	// Create a new CSV file
	shopNames := [...]string{
		"The Crafty Corner",
		"Savvy Styles",
		"Treasure Trove Emporium",
		"Gourmet Delights",
		"Vintage Finds",
		"Tech Haven",
		"Eco-Friendly Emporium",
		"Chic Boutique",
		"Golden Gadgetry",
		"Mystic Marketplace",
		"Harmony Haven",
		"Artisanal Attic",
		"Gadget Gallery",
		"Serenity Shoppe",
		"Luxury Loft",
		"Wholesome Wonders",
		"Dreamy Designs",
		"Enchanted Emporium",
		"Paws and Claws Corner",
		"Curio Cabinet",
		"Floral Finesse",
		"Adorned Alcove",
		"Gourmet Galleria",
		"Curious Corner",
		"Urban Umbrella",
		"Funky Fusion",
		"Bountiful Bazaar",
		"Cozy Corner",
		"Radiant Raiment",
		"Nautical Nook",
		"Sole Sanctuary",
		"Ethereal Elegance",
		"Dapper Den",
		"Fireside Finds",
		"Trendy Trinkets",
		"Sweet Serenade",
		"Eclectic Emporium",
		"Radiant Rugs",
		"Velvet Vault",
		"Nifty Nook",
		"Petal Pushers",
		"Galactic Goods",
		"Zen Zone",
		"Rustic Retreat",
		"Luxe Lounge",
		"Sunny Soiree",
		"Hearty Hearth",
		"Urban Oasis",
		"Whimsical Wonders",
		"Azure Abode",
		"Blossom Boutique",
		"Vibrant Visions",
		"Coastal Corner",
		"Chic Charm",
		"Tranquil Treasures",
		"Elite Essentials",
		"Bamboo Bazaar",
		"Crystal Cove",
		"Radiant Rose",
		"Gleaming Gallery",
		"Curious Curio",
		"Opulent Orchard",
		"Hidden Haven",
		"Swanky Soiree",
		"Ethereal Essence",
		"Festive Finds",
		"Infinite Ideas",
		"Coastal Creations",
		"Whimsy Workshop",
		"Luminous Lounge",
		"Crimson Creations",
		"Feline Fancies",
		"Artful Abode",
		"Rustic Rendezvous",
		"Plush Palace",
		"Glowing Gallery",
		"Starry Showcase",
		"Vintage Vault",
		"Zen Zephyr",
		"Wholesome Haven",
		"Radiant Retreat",
		"Golden Garden",
		"Majestic Manor",
		"Urban Utopia",
		"Frosted Flair",
		"Twilight Terrace",
		"Sapphire Sanctuary",
		"Blissful Boutique",
		"Emerald Emporium",
		"Sleek Sanctuary",
		"Rustic Rhythm",
		"Artisan Alley",
		"Fancy Fiesta",
		"Opulent Oasis",
		"Breezy Bazaar",
		"Sunny Side Store",
		"Mystical Mirage",
		"Dreamy Depot",
		"Vintage Vistas",
		"Pristine Pantry",
		"Tranquil Trinkets",
		"Eco Echo",
		"Chic Chicane",
		"Grand Gadgetry",
		"Bamboo Boutique",
		"Crimson Cove",
		"Crystal Cavern",
		"Gleaming Grove",
		"Dapper Depot",
		"Blissful Bazaar",
		"Sunny Side Shop",
		"Majestic Mosaic",
		"Wholesome Warehouse",
		"Adorned Arcade",
		"Emerald Enclave",
		"Ethereal Enclave",
		"Gilded Gazebo",
		"Hidden Halls",
		"Serene Stash",
		"Sleek Showcase",
		"Tranquil Trove",
		"Whimsical Wardrobe",
		"Zesty Zenith",
		"Sapphire Summit",
		"Azure Arcade",
		"Twilight Treasure",
		"Opulent Outpost",
		"Fancy Forum",
		"Pristine Pavilion",
		"Velvet Vista",
		"Vibrant Vestige",
		"Radiant Realm",
		"Golden Grove",
		"Mystical Mingle",
		"Urban Underground",
		"Savvy Sanctuary",
		"Pristine Palette",
		"Vintage Vantage",
		"Eco Enclave",
		"Whimsy Wharf",
		"Sleek Salon",
		"Sunny Sanctuary",
		"Breezy Boutique",
		"Cozy Cove",
		"Enchanted Ensemble",
		"Majestic Market",
		"Wholesome Wharf",
		"Artisan Arcade",
		"Jade Jewelers",
		"Kaleidoscope Corner",
		"Lush Lagoon",
		"Quaint Quarters",
		"Xanadu Expanse",
		"Yellow Yonder",
		"Artisan's Attic",
		"Elegant Emporium",
		"Golden Gallery",
		"Ivory Isles",
		"Jubilant Jamboree",
		"Kaleidoscope Keep",
		"Nautical Niche",
		"Quirky Quarters",
		"Sapphire Sanctum",
		"Vintage Vista",
		"Whimsy Warehouse",
		"Xanadu Exchange",
		"Yellow Yarn",
		"Zesty Zephyr",
		"Artisan's Alley",
		"Curious Cove",
		"Dazzling Depot",
	}

	provincies := []Province{
		{ProvinceID: 1, ProvinceName: "DKI Jakarta"},
		{ProvinceID: 2, ProvinceName: "Jawa Timur"},
		{ProvinceID: 3, ProvinceName: "Jawa Barat"},
		{ProvinceID: 4, ProvinceName: "Sumatera Utara"},
		{ProvinceID: 5, ProvinceName: "Jawa Tengah"},
		{ProvinceID: 6, ProvinceName: "Sulawesi Selatan"},
		{ProvinceID: 7, ProvinceName: "Sumatera Selatan"},
		{ProvinceID: 3, ProvinceName: "Jawa Barat"},
		{ProvinceID: 1, ProvinceName: "Banten"},
		{ProvinceID: 8, ProvinceName: "Kepulauan Riau"},
	}

	cities := []City{
		{CityID: 1, CityName: "Jakarta", ProvinceID: 1},
		{CityID: 2, CityName: "Surabaya", ProvinceID: 2},
		{CityID: 3, CityName: "Bandung", ProvinceID: 3},
		{CityID: 4, CityName: "Medan", ProvinceID: 4},
		{CityID: 5, CityName: "Semarang", ProvinceID: 5},
		{CityID: 6, CityName: "Makassar", ProvinceID: 6},
		{CityID: 7, CityName: "Palembang", ProvinceID: 7},
		{CityID: 8, CityName: "Depok", ProvinceID: 3},
		{CityID: 9, CityName: "Tangerang", ProvinceID: 1},
		{CityID: 10, CityName: "Batam", ProvinceID: 8},
	}

	// Define district data
	districts := []District{
		{CityID: 1, DistrictID: 101, DistrictName: "Central Jakarta", Latitude: -6.2088, Longitude: 106.8456},
		{CityID: 1, DistrictID: 102, DistrictName: "North Jakarta", Latitude: -6.1622, Longitude: 106.8163},
		{CityID: 1, DistrictID: 103, DistrictName: "South Jakarta", Latitude: -6.2615, Longitude: 106.8106},
		{CityID: 1, DistrictID: 104, DistrictName: "East Jakarta", Latitude: -6.2665, Longitude: 106.8865},
		{CityID: 1, DistrictID: 105, DistrictName: "West Jakarta", Latitude: -6.2088, Longitude: 106.8456},
		{CityID: 1, DistrictID: 106, DistrictName: "Thousand Islands", Latitude: -5.8245, Longitude: 106.5135},
		{CityID: 1, DistrictID: 107, DistrictName: "Jakarta Barat", Latitude: -6.1618, Longitude: 106.7439},
		{CityID: 1, DistrictID: 108, DistrictName: "Jakarta Timur", Latitude: -6.2297, Longitude: 106.8947},
		{CityID: 1, DistrictID: 109, DistrictName: "Jakarta Selatan", Latitude: -6.2615, Longitude: 106.8106},
		{CityID: 1, DistrictID: 110, DistrictName: "Jakarta Utara", Latitude: -6.1622, Longitude: 106.8163},
		// Add more districts for Jakarta
		{CityID: 2, DistrictID: 201, DistrictName: "Surabaya Utara", Latitude: -7.2794, Longitude: 112.7480},
		{CityID: 2, DistrictID: 202, DistrictName: "Surabaya Selatan", Latitude: -7.3313, Longitude: 112.7327},
		{CityID: 2, DistrictID: 203, DistrictName: "Surabaya Barat", Latitude: -7.2591, Longitude: 112.7413},
		{CityID: 2, DistrictID: 204, DistrictName: "Surabaya Timur", Latitude: -7.2794, Longitude: 112.7480},
		{CityID: 2, DistrictID: 205, DistrictName: "Surabaya Pusat", Latitude: -7.2588, Longitude: 112.7414},
		{CityID: 2, DistrictID: 206, DistrictName: "Surabaya Selatan", Latitude: -7.3313, Longitude: 112.7327},
		{CityID: 2, DistrictID: 207, DistrictName: "Surabaya Barat Daya", Latitude: -7.2940, Longitude: 112.7034},
		{CityID: 2, DistrictID: 208, DistrictName: "Surabaya Gunung Anyar", Latitude: -7.3149, Longitude: 112.7259},
		{CityID: 2, DistrictID: 209, DistrictName: "Surabaya Lakarsantri", Latitude: -7.3066, Longitude: 112.7491},
		{CityID: 2, DistrictID: 210, DistrictName: "Surabaya Sambikerep", Latitude: -7.2916, Longitude: 112.7467},
		// Add more districts for Surabaya
		// Add more districts for other cities
		{CityID: 3, DistrictID: 301, DistrictName: "Bandung Utara", Latitude: -6.9038, Longitude: 107.6186},
		{CityID: 3, DistrictID: 302, DistrictName: "Bandung Selatan", Latitude: -7.0526, Longitude: 107.6605},
		{CityID: 3, DistrictID: 303, DistrictName: "Bandung Barat", Latitude: -6.8235, Longitude: 107.5624},
		{CityID: 3, DistrictID: 304, DistrictName: "Bandung Timur", Latitude: -6.9345, Longitude: 107.6231},
		{CityID: 3, DistrictID: 305, DistrictName: "Bandung Tengah", Latitude: -6.9175, Longitude: 107.6191},
		// Add more districts for Bandung if needed
		{CityID: 4, DistrictID: 401, DistrictName: "Medan Barat", Latitude: 3.5903, Longitude: 98.6757},
		{CityID: 4, DistrictID: 402, DistrictName: "Medan Timur", Latitude: 3.6015, Longitude: 98.6611},
		{CityID: 4, DistrictID: 403, DistrictName: "Medan Selatan", Latitude: 3.5864, Longitude: 98.6623},
		{CityID: 4, DistrictID: 404, DistrictName: "Medan Utara", Latitude: 3.6052, Longitude: 98.6751},
		{CityID: 5, DistrictID: 501, DistrictName: "Semarang Barat", Latitude: -6.9847, Longitude: 110.4169},
		{CityID: 5, DistrictID: 502, DistrictName: "Semarang Timur", Latitude: -6.9919, Longitude: 110.4262},
		{CityID: 5, DistrictID: 503, DistrictName: "Semarang Selatan", Latitude: -7.0364, Longitude: 110.4474},
		{CityID: 5, DistrictID: 504, DistrictName: "Semarang Utara", Latitude: -6.9922, Longitude: 110.4262},
		{CityID: 6, DistrictID: 601, DistrictName: "Makassar Barat", Latitude: -5.1477, Longitude: 119.4428},
		{CityID: 6, DistrictID: 602, DistrictName: "Makassar Timur", Latitude: -5.1316, Longitude: 119.4506},
		{CityID: 6, DistrictID: 603, DistrictName: "Makassar Selatan", Latitude: -5.1428, Longitude: 119.4171},
		{CityID: 6, DistrictID: 604, DistrictName: "Makassar Utara", Latitude: -5.1321, Longitude: 119.4373},
		{CityID: 7, DistrictID: 701, DistrictName: "Palembang Barat", Latitude: -2.9814, Longitude: 104.7569},
		{CityID: 7, DistrictID: 702, DistrictName: "Palembang Timur", Latitude: -2.9922, Longitude: 104.7614},
		{CityID: 7, DistrictID: 703, DistrictName: "Palembang Selatan", Latitude: -3.0018, Longitude: 104.7418},
		{CityID: 7, DistrictID: 704, DistrictName: "Palembang Utara", Latitude: -2.9710, Longitude: 104.7563},
		{CityID: 8, DistrictID: 801, DistrictName: "Depok Barat", Latitude: -6.3918, Longitude: 106.7673},
		{CityID: 8, DistrictID: 802, DistrictName: "Depok Timur", Latitude: -6.3954, Longitude: 106.8214},
		{CityID: 8, DistrictID: 803, DistrictName: "Depok Selatan", Latitude: -6.4225, Longitude: 106.8141},
		{CityID: 8, DistrictID: 804, DistrictName: "Depok Utara", Latitude: -6.3646, Longitude: 106.8295},
		{CityID: 9, DistrictID: 901, DistrictName: "Tangerang Barat", Latitude: -6.1843, Longitude: 106.6416},
		{CityID: 9, DistrictID: 902, DistrictName: "Tangerang Selatan", Latitude: -6.2854, Longitude: 106.7121},
		{CityID: 9, DistrictID: 903, DistrictName: "Tangerang Timur", Latitude: -6.1914, Longitude: 106.7696},
		{CityID: 9, DistrictID: 904, DistrictName: "Tangerang Utara", Latitude: -6.1787, Longitude: 106.6306},
		{CityID: 10, DistrictID: 1001, DistrictName: "Batam Barat", Latitude: 1.0761, Longitude: 104.0260},
		{CityID: 10, DistrictID: 1002, DistrictName: "Batam Timur", Latitude: 1.1106, Longitude: 104.0407},
		{CityID: 10, DistrictID: 1003, DistrictName: "Batam Selatan", Latitude: 1.0606, Longitude: 104.0216},
		{CityID: 10, DistrictID: 1004, DistrictName: "Batam Utara", Latitude: 1.1496, Longitude: 104.0246},
		// Add more districts for Medan if needed
	}

	enableLocGen := true
	enableShopGen := true
	if enableLocGen {
		csvLoc, err := os.Create("../../files/loc.csv")
		if err != nil {
			panic(err)
		}
		defer csvLoc.Close()

		// Create a new CSV writer
		writerLoc := csv.NewWriter(csvLoc)

		pMap := make(map[int]Province, len(provincies))
		for _, p := range provincies {
			pMap[p.ProvinceID] = p
		}

		citiesMap := make(map[int]City, len(cities))
		for _, c := range cities {
			citiesMap[c.CityID] = c
		}

		districtArr := make([]int, 0, len(districts))
		for _, d := range districts {
			districtArr = append(districtArr, int(d.DistrictID))
		}

		// Function to convert Shop struct to a slice of strings for CSV row
		toCSVRow := func(d District) []string {
			return []string{
				fmt.Sprintf("%d", d.DistrictID),
				d.DistrictName,
				fmt.Sprintf("%v", d.Latitude),
				fmt.Sprintf("%v", d.Longitude),
				fmt.Sprintf("%d", d.CityID),
				citiesMap[d.CityID].CityName,
				fmt.Sprintf("%d", citiesMap[d.CityID].ProvinceID),
				pMap[citiesMap[d.CityID].ProvinceID].ProvinceName,
			}
		}

		// Write data to CSV
		var csvData [][]string
		for _, d := range districts {
			csvData = append(csvData, toCSVRow(d))
		}
		err = writerLoc.WriteAll(csvData)
		if err != nil {
			panic(err)
		}

		// Ensure all data is written to the file
		writerLoc.Flush()

		fmt.Println("Locations written to CSV successfully!")

		if enableShopGen {
			csvFile, err := os.Create("../../files/shops.csv")
			if err != nil {
				panic(err)
			}
			defer csvFile.Close()

			// Create a new CSV writer
			writer := csv.NewWriter(csvFile)
			shops := make([]Shop, 0, len(shopNames))
			for i, sName := range shopNames {
				shops = append(shops, Shop{
					ID:         (i + 1),
					Name:       sName,
					DistrictID: gofakeit.RandomInt(districtArr),
				})
			}

			// Function to convert Shop struct to a slice of strings for CSV row
			toCSVRow := func(shop Shop) []string {
				return []string{fmt.Sprintf("%d", shop.ID), shop.Name, fmt.Sprintf("%d", shop.DistrictID)}
			}

			// Write data to CSV
			var csvData [][]string
			for _, shop := range shops {
				csvData = append(csvData, toCSVRow(shop))
			}
			err = writer.WriteAll(csvData)
			if err != nil {
				panic(err)
			}

			// Ensure all data is written to the file
			writer.Flush()

			fmt.Println("Shops written to CSV successfully!")
		}
	}

	fmt.Println("Generator Done")

}
