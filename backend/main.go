package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Strukture podataka
type User struct {
	ID    int
	Name  string
	Email string
	Goal  string // lose_weight ili hypertrophy
}

type Food struct {
	Name     string
	Calories float64
	Protein  float64
	Carbs    float64
	Fat      float64
}

// Open Food Facts API response
type OFFProduct struct {
	Product struct {
		ProductName string `json:"product_name"`
		Nutriments  struct {
			EnergyKcal    float64 `json:"energy-kcal_100g"`
			Proteins      float64 `json:"proteins_100g"`
			Carbohydrates float64 `json:"carbohydrates_100g"`
			Fat           float64 `json:"fat_100g"`
		} `json:"nutriments"`
	} `json:"product"`
}

// Dohvatanje hrane sa Open Food Facts
func getFoodFromOFF(barcode string) (*Food, error) {
	url := fmt.Sprintf("https://world.openfoodfacts.org/api/v2/product/%s.json", barcode)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var product OFFProduct
	if err := json.Unmarshal(body, &product); err != nil {
		return nil, err
	}

	return &Food{
		Name:     product.Product.ProductName,
		Calories: product.Product.Nutriments.EnergyKcal,
		Protein:  product.Product.Nutriments.Proteins,
		Carbs:    product.Product.Nutriments.Carbohydrates,
		Fat:      product.Product.Nutriments.Fat,
	}, nil
}

// Generisanje plana ishrane
func generateMealPlan(goal string) []Food {
	barcodes := []string{"3274080005003", "3017620425035"} // primer proizvoda
	var plan []Food
	for _, bc := range barcodes {
		if f, err := getFoodFromOFF(bc); err == nil {
			plan = append(plan, *f)
		}
	}

	if goal == "lose_weight" {
		fmt.Println("\nPlan za mršavljenje:")
	} else {
		fmt.Println("\nPlan za hipertrofiju:")
	}

	return plan
}

func main() {
	fmt.Println("=== PROVERA SISTEMA ===\n")

	// 1️⃣ Konekcija na bazu
	fmt.Println("📡 Korak 1: Konekcija na bazu")
	fmt.Println("----------------------------")

	dsn := "root:Vojislav123!@tcp(127.0.0.1:3306)/app_db?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Printf("DSN: %s\n", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("❌ Greška pri otvaranju konekcije:")
		fmt.Println(err)
		fmt.Println("\n🔧 Mogući uzroci:")
		fmt.Println("   - Pogrešan DSN string")
		fmt.Println("   - Lozinka sa specijalnim karakterima")
		return
	}
	defer db.Close()
	fmt.Println("✅ sql.Open() je prošao")

	// Ping provera
	if err := db.Ping(); err != nil {
		fmt.Println("❌ Greška pri Ping() na bazu:")
		fmt.Println(err)
		fmt.Println("\n🔧 Mogući uzroci:")
		fmt.Println("   1. MySQL servis NIJE POKRENUT")
		fmt.Println("   2. Korisnik 'root' nema tu lozinku")
		fmt.Println("   3. Baza 'app_db' ne postoji")
		fmt.Println("   4. Port 3306 je zauzet ili drugačiji")
		fmt.Println("\n💡 Rešenja:")
		fmt.Println("   - Pokreni: mysql -u root -p")
		fmt.Println("   - Kreiraj bazu: CREATE DATABASE app_db;")
		return
	}
	fmt.Println("✅ Baza je povezana (Ping OK)\n")

	// Provera tabele
	fmt.Println("📋 Korak 2: Provera tabele 'users'")
	fmt.Println("-----------------------------------")
	row := db.QueryRow("SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_NAME = 'users'")
	var tableExists int
	if err := row.Scan(&tableExists); err != nil {
		fmt.Println("❌ Greška pri proveri tabele:")
		fmt.Println(err)
		return
	}

	if tableExists == 0 {
		fmt.Println("❌ Tabela 'users' NE POSTOJI!")
		fmt.Println("\n💡 Kreiraj tabelu sa:")
		fmt.Println(`
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    goal VARCHAR(50) NOT NULL
);
        `)
		return
	}
	fmt.Println("✅ Tabela 'users' postoji\n")

	// 2️⃣ Login / registracija
	fmt.Println("👤 Korak 3: Login / Registracija")
	fmt.Println("--------------------------------")
	var user User
	fmt.Print("Unesi email: ")
	fmt.Scanln(&user.Email)

	if user.Email == "" {
		fmt.Println("❌ Email je obavezan!")
		return
	}
	fmt.Printf("✅ Email unesen: %s\n", user.Email)

	// Proveri da li korisnik postoji
	row = db.QueryRow("SELECT id, name, goal FROM users WHERE email = ?", user.Email)
	err = row.Scan(&user.ID, &user.Name, &user.Goal)
	if err == sql.ErrNoRows {
		// Registracija
		fmt.Println("ℹ️  Korisnik ne postoji, registracija...")
		fmt.Print("Unesi ime: ")
		fmt.Scanln(&user.Name)
		fmt.Printf("✅ Ime unešeno: %s\n", user.Name)

		for {
			fmt.Print("Izaberi cilj (lose_weight / hypertrophy): ")
			fmt.Scanln(&user.Goal)
			if user.Goal == "lose_weight" || user.Goal == "hypertrophy" {
				break
			}
			fmt.Println("❌ Nevalidan cilj! Unesi lose_weight ili hypertrophy.")
		}
		fmt.Printf("✅ Cilj odabran: %s\n", user.Goal)

		_, err = db.Exec("INSERT INTO users (name, email, goal) VALUES (?, ?, ?)", user.Name, user.Email, user.Goal)
		if err != nil {
			fmt.Println("❌ Greška pri unosu u bazu:")
			fmt.Println(err)
			fmt.Println("\n🔧 Mogući uzroci:")
			fmt.Println("   - Email već postoji u bazi")
			fmt.Println("   - Greška u SQL upitu")
			return
		}
		fmt.Println("✅ Registracija uspešna!\n")
	} else if err != nil {
		fmt.Println("❌ Greška pri čitanju baze:")
		fmt.Println(err)
		return
	} else {
		fmt.Printf("✅ Dobrodošao nazad, %s! Tvoj cilj je: %s\n\n", user.Name, user.Goal)
	}

	// 3️⃣ Generiši plan ishrane
	fmt.Println("🍽️  Korak 4: Generisanje plana ishrane")
	fmt.Println("-------------------------------------")
	mealPlan := generateMealPlan(user.Goal)

	// 4️⃣ Ispiši plan
	fmt.Println("\nTvoj plan ishrane:")
	for _, f := range mealPlan {
		fmt.Printf("- %s: %.2f kcal, P: %.2fg, C: %.2fg, F: %.2fg\n",
			f.Name, f.Calories, f.Protein, f.Carbs, f.Fat)
	}
}
