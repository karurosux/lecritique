package cron

import (
	"context"
	"log"

	"lecritique/internal/auth/services"
	"github.com/robfig/cron/v3"
)

// SetupDeactivationCron sets up a daily cron job to process pending account deactivations
func SetupDeactivationCron(authService services.AuthService) *cron.Cron {
	c := cron.New()

	// Run daily at 2 AM to process pending deactivations
	_, err := c.AddFunc("0 2 * * *", func() {
		ctx := context.Background()
		log.Println("Running account deactivation job...")
		
		if err := authService.ProcessPendingDeactivations(ctx); err != nil {
			log.Printf("Error processing pending deactivations: %v", err)
		} else {
			log.Println("Account deactivation job completed successfully")
		}
	})

	if err != nil {
		log.Printf("Failed to schedule deactivation cron job: %v", err)
	}

	return c
}