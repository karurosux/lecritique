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
	`{"max_restaurants": 1, "max_locations_per_restaurant": 3, "max_qr_codes_per_location": 10, "max_feedbacks_per_month": 100, "max_team_members": 2, "advanced_analytics": false, "custom_branding": false, "api_access": false, "priority_support": false}`).Error

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
	`{"max_restaurants": 5, "max_locations_per_restaurant": 10, "max_qr_codes_per_location": 50, "max_feedbacks_per_month": 1000, "max_team_members": 10, "advanced_analytics": true, "custom_branding": true, "api_access": true, "priority_support": false}`).Error

	if err != nil {
		log.Printf("⚠️  Warning: Failed to create Professional plan: %v\n", err)
	} else {
		fmt.Println("✅ Created Professional plan")
	}

	// Enterprise Plan
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
	`, "enterprise", "Enterprise", "Unlimited scale with premium support and features", 199.99, "USD",
	`{"max_restaurants": -1, "max_locations_per_restaurant": -1, "max_qr_codes_per_location": -1, "max_feedbacks_per_month": -1, "max_team_members": -1, "advanced_analytics": true, "custom_branding": true, "api_access": true, "priority_support": true}`).Error

	if err != nil {
		log.Printf("⚠️  Warning: Failed to create Enterprise plan: %v\n", err)
	} else {
		fmt.Println("✅ Created Enterprise plan")
	}

	fmt.Println("\n🎉 Subscription plans created successfully!")
	fmt.Println("📊 Plans available:")
	fmt.Println("   • Starter: $29.99/month - 1 restaurant, 3 locations")
	fmt.Println("   • Professional: $79.99/month - 5 restaurants, 20 locations") 
	fmt.Println("   • Enterprise: $199.99/month - Unlimited everything")
}