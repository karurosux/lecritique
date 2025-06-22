package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lecritique/api/internal/shared/config"
	"github.com/lecritique/api/internal/shared/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Default test user data
	email := "admin@lecritique.com"
	password := "admin123"
	companyName := "LeCritique Demo Restaurant"

	fmt.Printf("Creating default user: %s\n", email)
	fmt.Printf("Company: %s\n", companyName)
	fmt.Printf("Password: %s\n", password)

	// Check if user already exists
	var existingAccount struct {
		ID string `gorm:"column:id"`
	}
	result := db.Table("accounts").Select("id").Where("email = ?", email).First(&existingAccount)
	if result.Error == nil {
		fmt.Printf("❌ User with email %s already exists (ID: %s)\n", email, existingAccount.ID)
		fmt.Println("Use --force flag to recreate the user")
		if len(os.Args) > 1 && os.Args[1] == "--force" {
			fmt.Println("🔄 Deleting existing user...")
			if err := db.Exec("DELETE FROM accounts WHERE email = ?", email).Error; err != nil {
				log.Fatal("Failed to delete existing user:", err)
			}
		} else {
			os.Exit(1)
		}
	} else if result.Error != gorm.ErrRecordNotFound {
		log.Fatal("Failed to check existing user:", result.Error)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Get starter subscription plan ID
	var subscriptionPlan struct {
		ID string `gorm:"column:id"`
	}
	result = db.Table("subscription_plans").Select("id").Where("code = ? AND is_active = true", "starter").First(&subscriptionPlan)
	subscriptionPlanID := ""
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Printf("⚠️  Warning: Could not find starter subscription plan: %v\n", result.Error)
	} else if result.Error == nil {
		subscriptionPlanID = subscriptionPlan.ID
	}

	// Insert the account
	var newAccount struct {
		ID string `gorm:"column:id"`
	}
	
	err = db.Raw(`
		INSERT INTO accounts (email, password_hash, company_name, is_active, email_verified, email_verified_at)
		VALUES (?, ?, ?, true, true, NOW())
		RETURNING id
	`, email, string(hashedPassword), companyName).Scan(&newAccount).Error
	
	if err != nil {
		log.Fatal("Failed to create account:", err)
	}
	
	accountID := newAccount.ID

	// Create a subscription if we have a plan
	if subscriptionPlanID != "" {
		err = db.Exec(`
			INSERT INTO subscriptions (account_id, plan_id, status, current_period_start, current_period_end)
			VALUES (?, ?, 'active', NOW(), NOW() + INTERVAL '1 month')
		`, accountID, subscriptionPlanID).Error
		
		if err != nil {
			log.Printf("⚠️  Warning: Failed to create subscription: %v\n", err)
		} else {
			fmt.Println("✅ Created subscription with starter plan")
		}
	}

	// Create a sample restaurant
	var newRestaurant struct {
		ID string `gorm:"column:id"`
	}
	err = db.Raw(`
		INSERT INTO restaurants (account_id, name, description, email, phone, website, is_active)
		VALUES (?, ?, ?, ?, ?, ?, true)
		RETURNING id
	`, accountID, "Demo Restaurant", "A sample restaurant for testing LeCritique", email, "+1-555-0123", "https://demo.lecritique.com").Scan(&newRestaurant).Error
	
	if err != nil {
		log.Printf("⚠️  Warning: Failed to create sample restaurant: %v\n", err)
	} else {
		restaurantID := newRestaurant.ID
		fmt.Println("✅ Created sample restaurant")

		// Create a sample location
		var newLocation struct {
			ID string `gorm:"column:id"`
		}
		err = db.Raw(`
			INSERT INTO locations (restaurant_id, name, address, city, state, country, postal_code, is_active)
			VALUES (?, ?, ?, ?, ?, ?, ?, true)
			RETURNING id
		`, restaurantID, "Main Location", "123 Restaurant St", "Food City", "CA", "USA", "12345").Scan(&newLocation).Error
		
		if err != nil {
			log.Printf("⚠️  Warning: Failed to create sample location: %v\n", err)
		} else {
			locationID := newLocation.ID
			fmt.Println("✅ Created sample location")

			// Create a sample QR code
			err = db.Exec(`
				INSERT INTO qr_codes (restaurant_id, location_id, code, label, type, is_active, expires_at)
				VALUES (?, ?, ?, ?, ?, true, NOW() + INTERVAL '1 year')
			`, restaurantID, locationID, "DEMO001", "Table 1", "table").Error
			
			if err != nil {
				log.Printf("⚠️  Warning: Failed to create sample QR code: %v\n", err)
			} else {
				fmt.Println("✅ Created sample QR code: DEMO001")
			}
		}

		// Create sample dishes
		dishes := []struct {
			name        string
			description string
			price       float64
			category    string
		}{
			{"Classic Burger", "Beef patty with lettuce, tomato, and house sauce", 12.99, "Mains"},
			{"Caesar Salad", "Fresh romaine with parmesan and croutons", 8.99, "Salads"},
			{"Chocolate Cake", "Rich chocolate cake with vanilla ice cream", 6.99, "Desserts"},
		}

		for _, dish := range dishes {
			err = db.Exec(`
				INSERT INTO dishes (restaurant_id, name, description, price, currency, category, is_active)
				VALUES (?, ?, ?, ?, 'USD', ?, true)
			`, restaurantID, dish.name, dish.description, dish.price, dish.category).Error
			
			if err != nil {
				log.Printf("⚠️  Warning: Failed to create dish %s: %v\n", dish.name, err)
			} else {
				fmt.Printf("✅ Created dish: %s\n", dish.name)
			}
		}
	}

	fmt.Println("\n🎉 Default user created successfully!")
	fmt.Println("📧 Email:", email)
	fmt.Println("🔑 Password:", password)
	fmt.Println("🏢 Company:", companyName)
	fmt.Printf("🆔 Account ID: %s\n", accountID)
	fmt.Println("\nYou can now login to the frontend with these credentials!")
}