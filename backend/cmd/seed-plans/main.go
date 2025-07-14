package main

import (
	"fmt"
	"log"

	"github.com/lecritique/api/internal/shared/config"
	"github.com/lecritique/api/internal/shared/database"
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

	fmt.Println("Creating subscription plans...")

	// Starter Plan
	err = db.Exec(`
		INSERT INTO subscription_plans (code, name, description, price, currency, features, is_active)
		VALUES (?, ?, ?, ?, ?, ?, true)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price = EXCLUDED.price,
			currency = EXCLUDED.currency,
			features = EXCLUDED.features,
			is_active = EXCLUDED.is_active
	`, "starter", "Starter", "Perfect for small restaurants just getting started", 29.99, "USD", 
	`{"max_restaurants": 1, "max_qr_codes": 15, "max_feedbacks_per_month": 50, "max_team_members": 2, "basic_analytics": true, "advanced_analytics": false, "feedback_explorer": true, "custom_branding": false, "priority_support": false}`).Error

	if err != nil {
		log.Printf("⚠️  Warning: Failed to create Starter plan: %v\n", err)
	} else {
		fmt.Println("✅ Created Starter plan")
	}

	// Professional Plan  
	err = db.Exec(`
		INSERT INTO subscription_plans (code, name, description, price, currency, features, is_active)
		VALUES (?, ?, ?, ?, ?, ?, true)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price = EXCLUDED.price,
			currency = EXCLUDED.currency,
			features = EXCLUDED.features,
			is_active = EXCLUDED.is_active
	`, "professional", "Professional", "For growing restaurant chains and franchises", 79.99, "USD",
	`{"max_restaurants": 5, "max_qr_codes": 125, "max_feedbacks_per_month": 250, "max_team_members": 10, "basic_analytics": true, "advanced_analytics": true, "feedback_explorer": true, "custom_branding": true, "priority_support": false}`).Error

	if err != nil {
		log.Printf("⚠️  Warning: Failed to create Professional plan: %v\n", err)
	} else {
		fmt.Println("✅ Created Professional plan")
	}

	// Premium Plan
	err = db.Exec(`
		INSERT INTO subscription_plans (code, name, description, price, currency, features, is_active)
		VALUES (?, ?, ?, ?, ?, ?, true)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price = EXCLUDED.price,
			currency = EXCLUDED.currency,
			features = EXCLUDED.features,
			is_active = EXCLUDED.is_active
	`, "premium", "Premium", "Unlimited scale with premium support and features", 199.99, "USD",
	`{"max_restaurants": 20, "max_qr_codes": 2000, "max_feedbacks_per_month": 1000, "max_team_members": 50, "basic_analytics": true, "advanced_analytics": true, "feedback_explorer": true, "custom_branding": true, "priority_support": true}`).Error

	if err != nil {
		log.Printf("⚠️  Warning: Failed to create Premium plan: %v\n", err)
	} else {
		fmt.Println("✅ Created Premium plan")
	}

	fmt.Println("\n🎉 Subscription plans created successfully!")
	fmt.Println("📊 Plans available:")
	fmt.Println("   • Starter: $29.99/month - 1 restaurant, 15 QR codes")
	fmt.Println("   • Professional: $79.99/month - 5 restaurants, 125 QR codes") 
	fmt.Println("   • Premium: $199.99/month - 20 restaurants, 2000 QR codes")
}