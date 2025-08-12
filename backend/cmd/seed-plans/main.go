package main

import (
	"fmt"
	"log"

	"kyooar/internal/shared/config"
	"kyooar/internal/shared/database"
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

	var seedRun struct {
		ID string `gorm:"column:id"`
	}
	result := db.Table("seed_runs").Select("id").Where("seed_name = ?", "subscription-plans").First(&seedRun)
	if result.Error == nil {
		fmt.Println("✅ Subscription plans seed already executed, skipping...")
		return
	}

	fmt.Println("Creating subscription plans...")

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
	`, "starter", "Starter", "Perfect for small businesses just getting started", 29.99, "USD",
	1, 10, 500, 2, true, false, true, false, false, 14).Error

	if err != nil {
		log.Printf("⚠️  Warning: Failed to create Starter plan: %v\n", err)
	} else {
		fmt.Println("✅ Created Starter plan")
	}

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
	`, "professional", "Professional", "For growing businesses and multiple locations", 79.99, "USD",
	3, 50, 2000, 5, true, false, true, false, false, 14).Error

	if err != nil {
		log.Printf("⚠️  Warning: Failed to create Professional plan: %v\n", err)
	} else {
		fmt.Println("✅ Created Professional plan")
	}

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
	`, "premium", "Premium", "Enterprise solution with advanced features and priority support", 199.99, "USD",
	10, 200, 5000, 20, true, true, true, false, true, 30).Error

	if err != nil {
		log.Printf("⚠️  Warning: Failed to create Premium plan: %v\n", err)
	} else {
		fmt.Println("✅ Created Premium plan")
	}

	fmt.Println("\n🎉 Subscription plans created successfully!")
	fmt.Println("📊 Plans available:")
	fmt.Println("   • Starter: $29.99/month - 1 organization, 10 QR codes, 500 feedbacks/month, 2 team members")
	fmt.Println("   • Professional: $79.99/month - 3 organizations, 50 QR codes, 2000 feedbacks/month, 5 team members") 
	fmt.Println("   • Premium: $199.99/month - 10 organizations, 200 QR codes, 5000 feedbacks/month, 20 team members + Advanced Analytics")

	err = db.Exec(`INSERT INTO seed_runs (seed_name, version) VALUES (?, ?)`, "subscription-plans", "1.0").Error
	if err != nil {
		log.Printf("⚠️  Warning: Failed to record seed run: %v\n", err)
	}
}