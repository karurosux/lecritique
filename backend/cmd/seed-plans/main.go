package main

import (
	"fmt"
	"log"

	"kyooar/internal/shared/config"
	"kyooar/internal/shared/database"
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

	// Free Plan
	err = db.Exec(`
		INSERT INTO subscription_plans (code, name, description, price, currency, 
			max_organizations, max_qr_codes, max_feedbacks_per_month, max_team_members,
			has_basic_analytics, has_advanced_analytics, has_feedback_explorer, 
			has_custom_branding, has_priority_support, is_active, is_visible, trial_days)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, true, true, ?)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price = EXCLUDED.price,
			currency = EXCLUDED.currency,
			max_organizations = EXCLUDED.max_organizations,
			max_qr_codes = EXCLUDED.max_qr_codes,
			max_feedbacks_per_month = EXCLUDED.max_feedbacks_per_month,
			max_team_members = EXCLUDED.max_team_members,
			has_basic_analytics = EXCLUDED.has_basic_analytics,
			has_advanced_analytics = EXCLUDED.has_advanced_analytics,
			has_feedback_explorer = EXCLUDED.has_feedback_explorer,
			has_custom_branding = EXCLUDED.has_custom_branding,
			has_priority_support = EXCLUDED.has_priority_support,
			is_active = EXCLUDED.is_active,
			is_visible = EXCLUDED.is_visible,
			trial_days = EXCLUDED.trial_days
	`, "free", "Free", "Perfect for trying out Kyooar", 0.00, "USD",
	1, 3, 25, 1, true, false, false, false, false, 0).Error

	if err != nil {
		log.Printf("⚠️  Warning: Failed to create Free plan: %v\n", err)
	} else {
		fmt.Println("✅ Created Free plan")
	}

	// Starter Plan
	err = db.Exec(`
		INSERT INTO subscription_plans (code, name, description, price, currency, 
			max_organizations, max_qr_codes, max_feedbacks_per_month, max_team_members,
			has_basic_analytics, has_advanced_analytics, has_feedback_explorer, 
			has_custom_branding, has_priority_support, is_active, is_visible, is_popular, trial_days)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, true, true, true, ?)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price = EXCLUDED.price,
			currency = EXCLUDED.currency,
			max_organizations = EXCLUDED.max_organizations,
			max_qr_codes = EXCLUDED.max_qr_codes,
			max_feedbacks_per_month = EXCLUDED.max_feedbacks_per_month,
			max_team_members = EXCLUDED.max_team_members,
			has_basic_analytics = EXCLUDED.has_basic_analytics,
			has_advanced_analytics = EXCLUDED.has_advanced_analytics,
			has_feedback_explorer = EXCLUDED.has_feedback_explorer,
			has_custom_branding = EXCLUDED.has_custom_branding,
			has_priority_support = EXCLUDED.has_priority_support,
			is_active = EXCLUDED.is_active,
			is_visible = EXCLUDED.is_visible,
			is_popular = EXCLUDED.is_popular,
			trial_days = EXCLUDED.trial_days
	`, "starter", "Starter", "Perfect for small organizations just getting started", 29.99, "USD",
	3, 15, 200, 3, true, true, true, false, false, 14).Error

	if err != nil {
		log.Printf("⚠️  Warning: Failed to create Starter plan: %v\n", err)
	} else {
		fmt.Println("✅ Created Starter plan")
	}

	// Professional Plan  
	err = db.Exec(`
		INSERT INTO subscription_plans (code, name, description, price, currency, 
			max_organizations, max_qr_codes, max_feedbacks_per_month, max_team_members,
			has_basic_analytics, has_advanced_analytics, has_feedback_explorer, 
			has_custom_branding, has_priority_support, is_active, is_visible, is_popular, trial_days)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, true, true, false, ?)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price = EXCLUDED.price,
			currency = EXCLUDED.currency,
			max_organizations = EXCLUDED.max_organizations,
			max_qr_codes = EXCLUDED.max_qr_codes,
			max_feedbacks_per_month = EXCLUDED.max_feedbacks_per_month,
			max_team_members = EXCLUDED.max_team_members,
			has_basic_analytics = EXCLUDED.has_basic_analytics,
			has_advanced_analytics = EXCLUDED.has_advanced_analytics,
			has_feedback_explorer = EXCLUDED.has_feedback_explorer,
			has_custom_branding = EXCLUDED.has_custom_branding,
			has_priority_support = EXCLUDED.has_priority_support,
			is_active = EXCLUDED.is_active,
			is_visible = EXCLUDED.is_visible,
			is_popular = EXCLUDED.is_popular,
			trial_days = EXCLUDED.trial_days
	`, "professional", "Professional", "For growing organization chains and franchises", 79.99, "USD",
	10, 50, 1000, 10, true, true, true, true, true, 14).Error

	if err != nil {
		log.Printf("⚠️  Warning: Failed to create Professional plan: %v\n", err)
	} else {
		fmt.Println("✅ Created Professional plan")
	}

	// Premium Plan
	err = db.Exec(`
		INSERT INTO subscription_plans (code, name, description, price, currency, 
			max_organizations, max_qr_codes, max_feedbacks_per_month, max_team_members,
			has_basic_analytics, has_advanced_analytics, has_feedback_explorer, 
			has_custom_branding, has_priority_support, is_active, is_visible, is_popular, trial_days)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, true, true, false, ?)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price = EXCLUDED.price,
			currency = EXCLUDED.currency,
			max_organizations = EXCLUDED.max_organizations,
			max_qr_codes = EXCLUDED.max_qr_codes,
			max_feedbacks_per_month = EXCLUDED.max_feedbacks_per_month,
			max_team_members = EXCLUDED.max_team_members,
			has_basic_analytics = EXCLUDED.has_basic_analytics,
			has_advanced_analytics = EXCLUDED.has_advanced_analytics,
			has_feedback_explorer = EXCLUDED.has_feedback_explorer,
			has_custom_branding = EXCLUDED.has_custom_branding,
			has_priority_support = EXCLUDED.has_priority_support,
			is_active = EXCLUDED.is_active,
			is_visible = EXCLUDED.is_visible,
			is_popular = EXCLUDED.is_popular,
			trial_days = EXCLUDED.trial_days
	`, "premium", "Premium", "Enterprise scale with premium support and features", 199.99, "USD",
	50, 500, 5000, 50, true, true, true, true, true, 30).Error

	if err != nil {
		log.Printf("⚠️  Warning: Failed to create Premium plan: %v\n", err)
	} else {
		fmt.Println("✅ Created Premium plan")
	}

	fmt.Println("\n🎉 Subscription plans created successfully!")
	fmt.Println("📊 Plans available:")
	fmt.Println("   • Free: $0.00/month - 1 organization, 3 QR codes, 25 feedbacks/month")
	fmt.Println("   • Starter: $29.99/month - 3 organizations, 15 QR codes, 200 feedbacks/month")
	fmt.Println("   • Professional: $79.99/month - 10 organizations, 50 QR codes, 1000 feedbacks/month") 
	fmt.Println("   • Premium: $199.99/month - 50 organizations, 500 QR codes, 5000 feedbacks/month")
}
