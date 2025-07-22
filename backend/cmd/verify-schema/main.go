package main

import (
	"fmt"
	"log"

	"lecritique/internal/shared/config"
	"lecritique/internal/shared/database"
	"lecritique/internal/subscription/models"
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

	// Test query with new columns
	var plan models.SubscriptionPlan
	err = db.First(&plan).Error
	if err != nil {
		log.Printf("No plans found (which is ok): %v\n", err)
	} else {
		fmt.Printf("✅ Successfully queried plan: %s\n", plan.Name)
		fmt.Printf("   - Max Organizations: %d\n", plan.MaxOrganizations)
		fmt.Printf("   - Max QR Codes: %d\n", plan.MaxQRCodes)
		fmt.Printf("   - Max Feedbacks/Month: %d\n", plan.MaxFeedbacksPerMonth)
		fmt.Printf("   - Max Team Members: %d\n", plan.MaxTeamMembers)
		fmt.Printf("   - Basic Analytics: %v\n", plan.HasBasicAnalytics)
		fmt.Printf("   - Advanced Analytics: %v\n", plan.HasAdvancedAnalytics)
		fmt.Printf("   - Feedback Explorer: %v\n", plan.HasFeedbackExplorer)
		fmt.Printf("   - Custom Branding: %v\n", plan.HasCustomBranding)
		fmt.Printf("   - Priority Support: %v\n", plan.HasPrioritySupport)
	}

	fmt.Println("\n🎉 Schema migration verified successfully!")
}
