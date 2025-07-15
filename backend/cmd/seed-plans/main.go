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
		INSERT INTO subscription_plans (code, name, description, price, currency, 
			max_restaurants, max_qr_codes, max_feedbacks_per_month, max_team_members,
			has_basic_analytics, has_advanced_analytics, has_feedback_explorer, 
			has_custom_branding, has_priority_support, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, true)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price = EXCLUDED.price,
			currency = EXCLUDED.currency,
			max_restaurants = EXCLUDED.max_restaurants,
			max_qr_codes = EXCLUDED.max_qr_codes,
			max_feedbacks_per_month = EXCLUDED.max_feedbacks_per_month,
			max_team_members = EXCLUDED.max_team_members,
			has_basic_analytics = EXCLUDED.has_basic_analytics,
			has_advanced_analytics = EXCLUDED.has_advanced_analytics,
			has_feedback_explorer = EXCLUDED.has_feedback_explorer,
			has_custom_branding = EXCLUDED.has_custom_branding,
			has_priority_support = EXCLUDED.has_priority_support,
			is_active = EXCLUDED.is_active
	`, "starter", "Starter", "Perfect for small restaurants just getting started", 29.99, "USD",
	1, 15, 50, 2, true, false, true, false, false).Error

	if err != nil {
		log.Printf("⚠️  Warning: Failed to create Starter plan: %v\n", err)
	} else {
		fmt.Println("✅ Created Starter plan")
	}

	// Professional Plan  
	err = db.Exec(`
		INSERT INTO subscription_plans (code, name, description, price, currency, 
			max_restaurants, max_qr_codes, max_feedbacks_per_month, max_team_members,
			has_basic_analytics, has_advanced_analytics, has_feedback_explorer, 
			has_custom_branding, has_priority_support, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, true)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price = EXCLUDED.price,
			currency = EXCLUDED.currency,
			max_restaurants = EXCLUDED.max_restaurants,
			max_qr_codes = EXCLUDED.max_qr_codes,
			max_feedbacks_per_month = EXCLUDED.max_feedbacks_per_month,
			max_team_members = EXCLUDED.max_team_members,
			has_basic_analytics = EXCLUDED.has_basic_analytics,
			has_advanced_analytics = EXCLUDED.has_advanced_analytics,
			has_feedback_explorer = EXCLUDED.has_feedback_explorer,
			has_custom_branding = EXCLUDED.has_custom_branding,
			has_priority_support = EXCLUDED.has_priority_support,
			is_active = EXCLUDED.is_active
	`, "professional", "Professional", "For growing restaurant chains and franchises", 79.99, "USD",
	5, 125, 250, 10, true, true, true, true, false).Error

	if err != nil {
		log.Printf("⚠️  Warning: Failed to create Professional plan: %v\n", err)
	} else {
		fmt.Println("✅ Created Professional plan")
	}

	// Premium Plan
	err = db.Exec(`
		INSERT INTO subscription_plans (code, name, description, price, currency, 
			max_restaurants, max_qr_codes, max_feedbacks_per_month, max_team_members,
			has_basic_analytics, has_advanced_analytics, has_feedback_explorer, 
			has_custom_branding, has_priority_support, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, true)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price = EXCLUDED.price,
			currency = EXCLUDED.currency,
			max_restaurants = EXCLUDED.max_restaurants,
			max_qr_codes = EXCLUDED.max_qr_codes,
			max_feedbacks_per_month = EXCLUDED.max_feedbacks_per_month,
			max_team_members = EXCLUDED.max_team_members,
			has_basic_analytics = EXCLUDED.has_basic_analytics,
			has_advanced_analytics = EXCLUDED.has_advanced_analytics,
			has_feedback_explorer = EXCLUDED.has_feedback_explorer,
			has_custom_branding = EXCLUDED.has_custom_branding,
			has_priority_support = EXCLUDED.has_priority_support,
			is_active = EXCLUDED.is_active
	`, "premium", "Premium", "Unlimited scale with premium support and features", 199.99, "USD",
	20, 2000, 1000, 50, true, true, true, true, true).Error

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