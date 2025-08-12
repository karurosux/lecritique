package cron

import (
	"context"
	"log"

	"kyooar/internal/auth/services"
	"github.com/robfig/cron/v3"
)

func SetupDeactivationCron(authService services.AuthService) *cron.Cron {
	c := cron.New()

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