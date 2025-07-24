package main

import (
	"fmt"
	"log"
	"os"

	"kyooar/internal/shared/config"
	"kyooar/internal/shared/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	email := "admin@kyooar.com"
	password := "admin123"
	companyName := "Kyooar Demo Organization"

	fmt.Printf("Creating default user: %s\n", email)
	fmt.Printf("Company: %s\n", companyName)
	fmt.Printf("Password: %s\n", password)

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

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

	var newAccount struct {
		ID string `gorm:"column:id"`
	}
	
	err = db.Raw(`
		INSERT INTO accounts (email, password_hash, name, is_active, email_verified, email_verified_at)
		VALUES (?, ?, ?, true, true, NOW())
		RETURNING id
	`, email, string(hashedPassword), companyName).Scan(&newAccount).Error
	
	if err != nil {
		log.Fatal("Failed to create account:", err)
	}
	
	accountID := newAccount.ID

	// Create owner team member record
	err = db.Exec(`
		INSERT INTO team_members (account_id, member_id, role, invited_by, invited_at, accepted_at, created_at, updated_at)
		VALUES (?, ?, 'OWNER', ?, NOW(), NOW(), NOW(), NOW())
	`, accountID, accountID, accountID).Error
	
	if err != nil {
		log.Fatal("Failed to create owner team member record:", err)
	} else {
		fmt.Println("✅ Created owner team member record")
	}

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

	var newOrganization struct {
		ID string `gorm:"column:id"`
	}
	err = db.Raw(`
		INSERT INTO organizations (account_id, name, description, email, phone, website, is_active)
		VALUES (?, ?, ?, ?, ?, ?, true)
		RETURNING id
	`, accountID, "Demo Organization", "A sample organization for testing Kyooar", email, "+1-555-0123", "https://demo.kyooar.com").Scan(&newOrganization).Error
	
	if err != nil {
		log.Printf("⚠️  Warning: Failed to create sample organization: %v\n", err)
	} else {
		organizationID := newOrganization.ID
		fmt.Println("✅ Created sample organization")

		err = db.Exec(`
			INSERT INTO locations (organization_id, name, address, city, state, country, postal_code, is_active)
			VALUES (?, ?, ?, ?, ?, ?, ?, true)
		`, organizationID, "Main Location", "123 Organization St", "Food City", "CA", "USA", "12345").Error
		
		if err != nil {
			log.Printf("⚠️  Warning: Failed to create sample location: %v\n", err)
		} else {
			fmt.Println("✅ Created sample location")

			err = db.Exec(`
				INSERT INTO qr_codes (organization_id, location, code, label, type, is_active, expires_at)
				VALUES (?, ?, ?, ?, ?, true, NOW() + INTERVAL '1 year')
			`, organizationID, "Main Location", "DEMO001", "Table 1", "table").Error
			
			if err != nil {
				log.Printf("⚠️  Warning: Failed to create sample QR code: %v\n", err)
			} else {
				fmt.Println("✅ Created sample QR code: DEMO001")
			}
		}

		products := []struct {
			name        string
			description string
			price       float64
			category    string
		}{
			{"Classic Burger", "Beef patty with lettuce, tomato, and house sauce", 12.99, "Mains"},
			{"Caesar Salad", "Fresh romaine with parmesan and croutons", 8.99, "Salads"},
			{"Chocolate Cake", "Rich chocolate cake with vanilla ice cream", 6.99, "Desserts"},
		}

		for _, product := range products {
			err = db.Exec(`
				INSERT INTO products (organization_id, name, description, price, currency, category, is_active)
				VALUES (?, ?, ?, ?, 'USD', ?, true)
			`, organizationID, product.name, product.description, product.price, product.category).Error
			
			if err != nil {
				log.Printf("⚠️  Warning: Failed to create product %s: %v\n", product.name, err)
			} else {
				fmt.Printf("✅ Created product: %s\n", product.name)
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
