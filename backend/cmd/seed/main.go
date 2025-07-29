package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"kyooar/internal/shared/config"
	"kyooar/internal/shared/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	// Disable GORM verbose logging during seeding
	db = db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Get all available subscription plans
	var subscriptionPlans []struct {
		ID   string `gorm:"column:id"`
		Code string `gorm:"column:code"`
		Name string `gorm:"column:name"`
	}
	err = db.Table("subscription_plans").Select("id, code, name").Where("is_active = true AND code != 'free'").Find(&subscriptionPlans).Error
	if err != nil {
		log.Fatal("Failed to get subscription plans:", err)
	}

	if len(subscriptionPlans) == 0 {
		log.Fatal("No subscription plans found!")
	}

	fmt.Printf("✅ Found %d subscription plans\n", len(subscriptionPlans))

	// Create accounts for each subscription plan
	planAccounts := make(map[string][]string) // planCode -> [ownerAccountID, memberAccountID]
	
	for _, plan := range subscriptionPlans {
		fmt.Printf("\n📋 Creating accounts for %s plan...\n", plan.Name)
		
		// Account emails based on plan
		ownerEmail := fmt.Sprintf("admin_%s@kyooar.com", plan.Code)
		memberEmail := fmt.Sprintf("viewer_%s@kyooar.com", plan.Code)
		password := "Pass123!"

		// Create owner account
		var ownerAccountID string
		{
			// Check if owner account exists
			var existingAccount struct {
				ID string `gorm:"column:id"`
			}
			result := db.Table("accounts").Select("id").Where("email = ?", ownerEmail).First(&existingAccount)
			if result.Error == nil {
				fmt.Printf("❌ Owner user with email %s already exists\n", ownerEmail)
				if len(os.Args) > 1 && os.Args[1] == "--force" {
					fmt.Println("🔄 Deleting existing owner user...")
					if err := db.Exec("DELETE FROM accounts WHERE email = ?", ownerEmail).Error; err != nil {
						log.Printf("Failed to delete existing owner user: %v\n", err)
						continue
					}
				} else {
					continue
				}
			}

			// Create owner account
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Failed to hash owner password: %v\n", err)
				continue
			}

			ownerName := fmt.Sprintf("Kyooar %s Owner", plan.Name)
			err = db.Raw(`
				INSERT INTO accounts (email, password_hash, name, is_active, email_verified, email_verified_at)
				VALUES (?, ?, ?, true, true, NOW())
				RETURNING id
			`, ownerEmail, string(hashedPassword), ownerName).Scan(&ownerAccountID).Error
			
			if err != nil {
				log.Printf("Failed to create owner account: %v\n", err)
				continue
			}
			// Owner account created

			// Create owner team member record (owner of their own account)
			err = db.Exec(`
				INSERT INTO team_members (account_id, member_id, role, invited_by, invited_at, accepted_at, created_at, updated_at)
				VALUES (?, ?, 'OWNER', ?, NOW(), NOW(), NOW(), NOW())
			`, ownerAccountID, ownerAccountID, ownerAccountID).Error
			if err != nil {
				log.Printf("Failed to create owner team member record: %v\n", err)
			}

			// Create subscription for owner
			err = db.Exec(`
				INSERT INTO subscriptions (account_id, plan_id, status, current_period_start, current_period_end)
				VALUES (?, ?, 'active', NOW(), NOW() + INTERVAL '1 month')
			`, ownerAccountID, plan.ID).Error
			if err != nil {
				log.Printf("Failed to create owner subscription: %v\n", err)
			}
		}

		// Create member account
		var memberAccountID string
		{
			// Check if member account exists
			var existingAccount struct {
				ID string `gorm:"column:id"`
			}
			result := db.Table("accounts").Select("id").Where("email = ?", memberEmail).First(&existingAccount)
			if result.Error == nil {
				fmt.Printf("❌ Member user with email %s already exists\n", memberEmail)
				if len(os.Args) > 1 && os.Args[1] == "--force" {
					fmt.Println("🔄 Deleting existing member user...")
					if err := db.Exec("DELETE FROM accounts WHERE email = ?", memberEmail).Error; err != nil {
						log.Printf("Failed to delete existing member user: %v\n", err)
						continue
					}
				} else {
					continue
				}
			}

			// Create member account
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Failed to hash member password: %v\n", err)
				continue
			}

			memberName := fmt.Sprintf("Kyooar %s Viewer", plan.Name)
			err = db.Raw(`
				INSERT INTO accounts (email, password_hash, name, is_active, email_verified, email_verified_at)
				VALUES (?, ?, ?, true, true, NOW())
				RETURNING id
			`, memberEmail, string(hashedPassword), memberName).Scan(&memberAccountID).Error
			
			if err != nil {
				log.Printf("Failed to create member account: %v\n", err)
				continue
			}
			// Member account created

			// Create member team member record (owner of their own account)
			err = db.Exec(`
				INSERT INTO team_members (account_id, member_id, role, invited_by, invited_at, accepted_at, created_at, updated_at)
				VALUES (?, ?, 'OWNER', ?, NOW(), NOW(), NOW(), NOW())
			`, memberAccountID, memberAccountID, memberAccountID).Error
			if err != nil {
				log.Printf("Failed to create member team member record: %v\n", err)
			}
		}

		if ownerAccountID != "" && memberAccountID != "" {
			planAccounts[plan.Code] = []string{ownerAccountID, memberAccountID}
		}
	}

	// Create organizations for each subscription plan
	for _, plan := range subscriptionPlans {
		accounts, exists := planAccounts[plan.Code]
		if !exists || len(accounts) != 2 {
			log.Printf("⚠️  Skipping organizations for %s plan - accounts not found\n", plan.Code)
			continue
		}
		
		ownerAccountID := accounts[0]
		viewerAccountID := accounts[1]
		
		fmt.Printf("\n🏢 Creating organizations for %s plan...\n", plan.Name)
		
		// Define organizations for this plan
		planOrganizations := []struct {
			orgName       string
			orgDesc       string
			orgPhone      string
			orgWebsite    string
			products      []Product
			qrCodes       []QRCode
			isTourCompany bool
			location      string
			city          string
		}{
			{
				orgName:       fmt.Sprintf("%s Tours - %s", plan.Name, "Adventure Guides"),
				orgDesc:       fmt.Sprintf("%s tier adventure and cultural tour services", plan.Name),
				orgPhone:      fmt.Sprintf("+1-555-%s1", strings.ToUpper(plan.Code[:4])),
				orgWebsite:    fmt.Sprintf("https://%s-tours.com", strings.ToLower(plan.Code)),
				isTourCompany: true,
				location:      fmt.Sprintf("%s Tourism Center", plan.Name),
				city:          "San Francisco",
				products: []Product{
					{"City Walking Tour", "2-hour guided walking tour of historic downtown", 45.00, "Tours"},
					{"Food & Culture Tour", "4-hour culinary journey through local neighborhoods", 89.00, "Tours"},
					{"Museum Package", "All-day museum pass with guided explanations", 65.00, "Tours"},
					{"Sunset Boat Tour", "Evening boat tour with champagne service", 120.00, "Tours"},
					{"Private Group Tour", "Customized private tour for groups up to 20", 350.00, "Tours"},
					{"Photography Tour", "Guided tour to the best photo spots with tips", 75.00, "Tours"},
				},
				qrCodes: []QRCode{
					{fmt.Sprintf("%s-TOUR-DESK-01", strings.ToUpper(plan.Code)), "Reception Desk", "feedback_point"},
					{fmt.Sprintf("%s-TOUR-BUS-01", strings.ToUpper(plan.Code)), "Tour Bus #1", "vehicle"},
					{fmt.Sprintf("%s-TOUR-BUS-02", strings.ToUpper(plan.Code)), "Tour Bus #2", "vehicle"},
					{fmt.Sprintf("%s-TOUR-MEETING", strings.ToUpper(plan.Code)), "Meeting Point Plaza", "location"},
					{fmt.Sprintf("%s-TOUR-OFFICE", strings.ToUpper(plan.Code)), "Main Office", "feedback_point"},
				},
			},
			{
				orgName:       fmt.Sprintf("%s Print Solutions", plan.Name),
				orgDesc:       fmt.Sprintf("%s tier printing and business services", plan.Name),
				orgPhone:      fmt.Sprintf("+1-555-%s2", strings.ToUpper(plan.Code[:4])),
				orgWebsite:    fmt.Sprintf("https://%s-print.com", strings.ToLower(plan.Code)),
				isTourCompany: false,
				location:      fmt.Sprintf("%s Business Center", plan.Name),
				city:          "New York",
				products: []Product{
					{"B&W Printing", "Black and white printing per page", 0.10, "Printing"},
					{"Color Printing", "Full color printing per page", 0.50, "Printing"},
					{"Poster Printing", "Large format poster printing (24x36)", 25.00, "Printing"},
					{"Business Cards", "500 premium business cards", 45.00, "Printing"},
					{"Binding Service", "Professional document binding", 5.00, "Services"},
					{"International Calls", "Per minute international calling", 0.25, "Call Services"},
					{"Fax Service", "Send or receive fax per page", 2.00, "Services"},
					{"Scanning Service", "Document scanning per page", 0.15, "Services"},
				},
				qrCodes: []QRCode{
					{fmt.Sprintf("%s-PRINT-DESK-01", strings.ToUpper(plan.Code)), "Service Counter 1", "counter"},
					{fmt.Sprintf("%s-PRINT-DESK-02", strings.ToUpper(plan.Code)), "Service Counter 2", "counter"},
					{fmt.Sprintf("%s-PRINT-SELF-01", strings.ToUpper(plan.Code)), "Self Service Station 1", "kiosk"},
					{fmt.Sprintf("%s-PRINT-SELF-02", strings.ToUpper(plan.Code)), "Self Service Station 2", "kiosk"},
					{fmt.Sprintf("%s-PRINT-CALL-01", strings.ToUpper(plan.Code)), "Call Booth 1", "booth"},
					{fmt.Sprintf("%s-PRINT-CALL-02", strings.ToUpper(plan.Code)), "Call Booth 2", "booth"},
					{fmt.Sprintf("%s-PRINT-PICKUP", strings.ToUpper(plan.Code)), "Order Pickup Area", "location"},
				},
			},
		}

		// Create organizations for this plan
		for _, org := range planOrganizations {
			fmt.Printf("\n📋 Creating organization: %s\n", org.orgName)

			// Create organization owned by plan owner
			var organizationID string
			err = db.Raw(`
				INSERT INTO organizations (account_id, name, description, email, phone, website, is_active)
				VALUES (?, ?, ?, ?, ?, ?, true)
				RETURNING id
			`, ownerAccountID, org.orgName, org.orgDesc, fmt.Sprintf("admin_%s@kyooar.com", plan.Code), org.orgPhone, org.orgWebsite).Scan(&organizationID).Error
			
			if err != nil {
				log.Printf("Failed to create organization: %v\n", err)
				continue
			}
			// Organization created

			// Add plan viewer as team member to this organization (check for existing first)
			var existingTeamMember string
			checkResult := db.Raw("SELECT id FROM team_members WHERE account_id = ? AND member_id = ?", ownerAccountID, viewerAccountID).Scan(&existingTeamMember)
			if checkResult.Error != nil || existingTeamMember == "" {
				err = db.Exec(`
					INSERT INTO team_members (account_id, member_id, role, invited_by, invited_at, accepted_at, created_at, updated_at)
					VALUES (?, ?, 'VIEWER', ?, NOW(), NOW(), NOW(), NOW())
				`, ownerAccountID, viewerAccountID, ownerAccountID).Error
				if err != nil {
					log.Printf("Failed to add %s viewer to organization team: %v\n", plan.Code, err)
				}
			}

			// Create location
			locationName := org.location
			err = db.Exec(`
				INSERT INTO locations (organization_id, name, address, city, state, country, postal_code, is_active)
				VALUES (?, ?, ?, ?, ?, ?, ?, true)
			`, organizationID, locationName, "123 Business Plaza", org.city, "CA", "USA", "94102").Error
			if err != nil {
				log.Printf("Failed to create location: %v\n", err)
				continue
			}

			// Create products
			productIDs := make(map[string]string)
			for _, product := range org.products {
				var productID string
				err = db.Raw(`
					INSERT INTO products (organization_id, name, description, price, currency, category, is_active)
					VALUES (?, ?, ?, ?, 'USD', ?, true)
					RETURNING id
				`, organizationID, product.Name, product.Description, product.Price, product.Category).Scan(&productID).Error
				
				if err != nil {
					log.Printf("Failed to create product %s: %v\n", product.Name, err)
				} else {
					// Validate that we got a proper UUID
					if _, err := uuid.Parse(productID); err != nil {
						log.Printf("⚠️  Product %s returned invalid UUID: %s\n", product.Name, productID)
					} else {
						productIDs[product.Name] = productID
					}
				}
			}

			// Create QR codes
			qrCodeIDs := make(map[string]string)
			for _, qr := range org.qrCodes {
				var qrID string
				err = db.Raw(`
					INSERT INTO qr_codes (organization_id, location, code, label, type, is_active, expires_at)
					VALUES (?, ?, ?, ?, ?, true, NOW() + INTERVAL '2 years')
					RETURNING id
				`, organizationID, locationName, qr.Code, qr.Label, qr.Type).Scan(&qrID).Error
				
				if err != nil {
					log.Printf("Failed to create QR code %s: %v\n", qr.Code, err)
				} else {
					// Validate that we got a proper UUID
					if _, err := uuid.Parse(qrID); err != nil {
						log.Printf("⚠️  QR code %s returned invalid UUID: %s\n", qr.Code, qrID)
					} else {
						qrCodeIDs[qr.Code] = qrID
						
						// Add realistic scan data
						scansCount := 50 + rand.Intn(200) // 50-250 scans
						daysAgo := rand.Intn(7) + 1       // Last scan 1-7 days ago
						lastScannedAt := time.Now().AddDate(0, 0, -daysAgo)
						
						// Update QR code with scan data
						err = db.Exec(`
							UPDATE qr_codes 
							SET scans_count = ?, last_scanned_at = ? 
							WHERE id = ?
						`, scansCount, lastScannedAt, qrID).Error
						if err != nil {
							log.Printf("Failed to update scan data for QR code %s: %v\n", qr.Code, err)
						}
					}
				}
			}

			// Create questionnaires
			if org.isTourCompany {
				createTourQuestionnaires(db, organizationID, productIDs)
			} else {
				createPrintQuestionnaires(db, organizationID, productIDs)
			}

			// Create feedback
			if len(qrCodeIDs) > 0 && len(productIDs) > 0 {
				createFeedback(db, organizationID, qrCodeIDs, productIDs, org.isTourCompany)
			}
		}
	}

	fmt.Println("\n🎉 Seed data created successfully!")
	fmt.Println("\n📧 Login credentials (password: Pass123!):")
	
	for _, plan := range subscriptionPlans {
		fmt.Printf("  %s: admin_%s@kyooar.com / viewer_%s@kyooar.com\n", plan.Name, plan.Code, plan.Code)
	}
	
	fmt.Println("\n✨ Each plan has 2 organizations with products, QR codes, and ~920 feedback records")
}

type Product struct {
	Name        string
	Description string
	Price       float64
	Category    string
}

type QRCode struct {
	Code  string
	Label string
	Type  string
}

func createTourQuestionnaires(db *gorm.DB, orgID string, productIDs map[string]string) []string {
	var questionnaireIDs []string

	// Create questionnaires and questions for each product
	for productName, productID := range productIDs {
		// Create product-specific questionnaire
		var qID string
		err := db.Raw(`
			INSERT INTO questionnaires (organization_id, product_id, name, description, is_active)
			VALUES (?, ?, ?, ?, true)
			RETURNING id
		`, orgID, productID, fmt.Sprintf("%s Feedback", productName), "Help us improve your tour experience").Scan(&qID).Error
		
		if err == nil {
			questionnaireIDs = append(questionnaireIDs, qID)
			
			// Add questions for this product
			questions := []struct {
				text     string
				qtype    string
				required bool
				options  []string
			}{
				{"How would you rate your overall tour experience?", "rating", true, nil},
				{"How knowledgeable was your tour guide?", "scale", true, nil},
				{"Would you recommend this tour to friends?", "yes_no", true, nil},
				{"What type of tour experience do you prefer?", "single_choice", true, []string{"Historical sites", "Cultural experiences", "Nature/Scenic", "Food & Dining", "Adventure activities"}},
				{"What did you enjoy most about the tour?", "text", false, nil},
				{"What aspects could we improve?", "multi_choice", false, []string{"Guide knowledge", "Tour duration", "Group size", "Meeting point", "Price value", "Route selection"}},
				{"Overall satisfaction level", "scale", true, nil},
			}

			for i, q := range questions {
				var optionsParam interface{}
				if q.options != nil {
					// Convert to PostgreSQL array format
					optionsParam = fmt.Sprintf("{%s}", joinStringSlice(q.options, ","))
				}
				
				err = db.Exec(`
					INSERT INTO questions (product_id, text, type, is_required, options, display_order)
					VALUES (?, ?, ?, ?, ?, ?)
				`, productID, q.text, q.qtype, q.required, optionsParam, i+1).Error
				
				if err != nil {
					fmt.Printf("Failed to create question for %s: %v\n", productName, err)
				}
			}
		}
	}

	return questionnaireIDs
}

func createPrintQuestionnaires(db *gorm.DB, orgID string, productIDs map[string]string) []string {
	var questionnaireIDs []string

	// Create questionnaires and questions for each product
	for productName, productID := range productIDs {
		// Create product-specific questionnaire
		var qID string
		err := db.Raw(`
			INSERT INTO questionnaires (organization_id, product_id, name, description, is_active)
			VALUES (?, ?, ?, ?, true)
			RETURNING id
		`, orgID, productID, fmt.Sprintf("%s Feedback", productName), "Help us improve our services").Scan(&qID).Error
		
		if err == nil {
			questionnaireIDs = append(questionnaireIDs, qID)
			
			// Add questions for this product
			questions := []struct {
				text     string
				qtype    string
				required bool
				options  []string
			}{
				{"How satisfied are you with our service?", "rating", true, nil},
				{"How was the waiting time?", "scale", true, nil},
				{"Would you use our services again?", "yes_no", true, nil},
				{"What time of day do you typically visit us?", "single_choice", false, []string{"Morning (8-11 AM)", "Midday (11 AM-2 PM)", "Afternoon (2-5 PM)", "Evening (5-8 PM)"}},
				{"Quality of the service/materials", "rating", true, nil},
				{"Staff helpfulness", "rating", true, nil},
				{"Any suggestions for improvement?", "text", false, nil},
			}

			// Add service-specific questions
			if productName == "International Calls" || productName == "Fax Service" {
				questions = append(questions, struct {
					text     string
					qtype    string
					required bool
					options  []string
				}{"Call/connection quality", "scale", true, nil})
			}

			for i, q := range questions {
				var optionsParam interface{}
				if q.options != nil {
					// Convert to PostgreSQL array format
					optionsParam = fmt.Sprintf("{%s}", joinStringSlice(q.options, ","))
				}
				
				err = db.Exec(`
					INSERT INTO questions (product_id, text, type, is_required, options, display_order)
					VALUES (?, ?, ?, ?, ?, ?)
				`, productID, q.text, q.qtype, q.required, optionsParam, i+1).Error
				
				if err != nil {
					fmt.Printf("Failed to create question for %s: %v\n", productName, err)
				}
			}
		}
	}

	return questionnaireIDs
}

func createFeedback(db *gorm.DB, orgID string, qrCodeIDs map[string]string, productIDs map[string]string, isTourCompany bool) {
	// Generate feedback over the past 90 days for better testing
	now := time.Now()
	
	// Create feedback counts based on actual QR codes created
	feedbackCounts := make(map[string]int)
	
	// Different traffic patterns for different QR code types
	for qrCode := range qrCodeIDs {
		if isTourCompany {
			// Tourist company has more weekend traffic - increased for testing
			if strings.Contains(qrCode, "TOUR-DESK") {
				feedbackCounts[qrCode] = 85  // Reception desk - high traffic
			} else if strings.Contains(qrCode, "TOUR-BUS-01") {
				feedbackCounts[qrCode] = 120 // Main tour bus - highest traffic
			} else if strings.Contains(qrCode, "TOUR-BUS-02") {
				feedbackCounts[qrCode] = 95  // Secondary bus
			} else if strings.Contains(qrCode, "TOUR-MEETING") {
				feedbackCounts[qrCode] = 65  // Meeting point
			} else if strings.Contains(qrCode, "TOUR-OFFICE") {
				feedbackCounts[qrCode] = 75  // Main office
			}
		} else {
			// Print shop has more weekday traffic - increased for testing
			if strings.Contains(qrCode, "PRINT-DESK-01") {
				feedbackCounts[qrCode] = 110 // Main service counter
			} else if strings.Contains(qrCode, "PRINT-DESK-02") {
				feedbackCounts[qrCode] = 98  // Secondary counter
			} else if strings.Contains(qrCode, "PRINT-SELF-01") {
				feedbackCounts[qrCode] = 70  // Self-service stations
			} else if strings.Contains(qrCode, "PRINT-SELF-02") {
				feedbackCounts[qrCode] = 65
			} else if strings.Contains(qrCode, "PRINT-CALL-01") {
				feedbackCounts[qrCode] = 45  // Call booths
			} else if strings.Contains(qrCode, "PRINT-CALL-02") {
				feedbackCounts[qrCode] = 38
			} else if strings.Contains(qrCode, "PRINT-PICKUP") {
				feedbackCounts[qrCode] = 55  // Pickup area
			}
		}
	}

	customerNames := []string{
		"John Smith", "Emma Johnson", "Michael Brown", "Sarah Davis", "William Wilson",
		"Jennifer Garcia", "David Martinez", "Lisa Anderson", "Robert Taylor", "Maria Rodriguez",
		"James Thomas", "Patricia Lee", "Charles White", "Linda Harris", "Daniel Martin",
		"Barbara Thompson", "Joseph Clark", "Elizabeth Lewis", "Thomas Walker", "Susan Hall",
		"Christopher Young", "Jessica Allen", "Matthew King", "Nancy Wright", "Anthony Lopez",
		"Amanda Chen", "Kevin O'Connor", "Sophia Patel", "Ryan Murphy", "Grace Kim",
		"Tyler Jackson", "Olivia Rodriguez", "Nathan Brown", "Isabella Martinez", "Ethan Davis",
		"Mia Thompson", "Alexander Lee", "Charlotte Wilson", "Benjamin Garcia", "Amelia Anderson",
		"Lucas Miller", "Harper Taylor", "Mason Moore", "Evelyn White", "Logan Harris",
		"Abigail Clark", "Jacob Lewis", "Emily Walker", "Michael Hall", "Madison Young",
		"Aiden Allen", "Elizabeth King", "Jackson Wright", "Avery Lopez", "Sebastian Hill",
	}

	customerEmails := []string{
		"john.smith@gmail.com", "emma.j@outlook.com", "m.brown@yahoo.com", "sarah.davis@email.com",
		"w.wilson@company.com", "jen.garcia@gmail.com", "david.m@work.com", "lisa.anderson@email.com",
		"robert.t@business.org", "maria.r@service.com", "james.thomas@email.com", "patricia.lee@gmail.com",
		"c.white@company.com", "linda.h@business.com", "daniel.martin@email.com", "barbara.t@work.org",
		"joseph.clark@email.com", "elizabeth.lewis@gmail.com", "thomas.walker@company.com", "susan.hall@email.com",
		"", "", "", "", "", "", // Some customers don't provide email (20% anonymous)
	}

	// Convert productIDs map to slice for random selection
	var productIDsList []string
	for _, id := range productIDs {
		productIDsList = append(productIDsList, id)
	}

	fmt.Printf("\n📝 Creating feedback for organization...\n")

	for qrCode, count := range feedbackCounts {
		qrID, exists := qrCodeIDs[qrCode]
		if !exists {
			fmt.Printf("⚠️  Skipping feedback for %s - QR code not found\n", qrCode)
			continue
		}
		
		// Validate QR code UUID before using it
		if _, err := uuid.Parse(qrID); err != nil {
			fmt.Printf("⚠️  Skipping feedback for %s - invalid QR UUID: %s\n", qrCode, qrID)
			continue
		}

		for i := 0; i < count; i++ {
			// Spread feedback over past 90 days with realistic patterns
			daysAgo := rand.Intn(90)
			feedbackDate := now.AddDate(0, 0, -daysAgo)
			
			// More feedback during business hours
			hour := 9 + rand.Intn(10) // 9 AM to 7 PM
			if isTourCompany && (feedbackDate.Weekday() == time.Saturday || feedbackDate.Weekday() == time.Sunday) {
				hour = 8 + rand.Intn(12) // 8 AM to 8 PM on weekends
			}
			feedbackDate = time.Date(feedbackDate.Year(), feedbackDate.Month(), feedbackDate.Day(), hour, rand.Intn(60), 0, 0, feedbackDate.Location())

			// Select random product
			productID := productIDsList[rand.Intn(len(productIDsList))]
			
			// Validate product UUID
			if _, err := uuid.Parse(productID); err != nil {
				fmt.Printf("⚠️  Skipping feedback - invalid product UUID: %s\n", productID)
				continue
			}

			// Random customer info
			customerName := customerNames[rand.Intn(len(customerNames))]
			customerEmail := customerEmails[rand.Intn(len(customerEmails))]
			if customerEmail == "" {
				customerEmail = fmt.Sprintf("%s.%d@example.com", qrCode, i)
			}

			// Get questions for this product
			var questions []struct {
				ID      string  `gorm:"column:id"`
				Text    string  `gorm:"column:text"`
				Type    string  `gorm:"column:type"`
				Options *string `gorm:"column:options"`
			}
			db.Table("questions").Select("id, text, type, options").Where("product_id = ?", productID).Find(&questions)

			if len(questions) == 0 {
				continue // Skip if no questions for this product
			}

			// Create responses JSONB array
			responses := make([]map[string]interface{}, 0)
			overallRating := 3 + rand.Intn(3) // 3-5 stars for overall rating

			for _, q := range questions {
				var answer interface{}
				switch q.Type {
				case "rating":
					// Mostly positive ratings with some variation
					answer = 3 + rand.Intn(3) // 3-5 stars
				case "scale":
					// Scale 1-10
					answer = 6 + rand.Intn(5) // 6-10
				case "yes_no":
					// 80% positive
					if rand.Float32() < 0.8 {
						answer = "yes"
					} else {
						answer = "no"
					}
				case "single_choice":
					// Select one random option
					if q.Options != nil && *q.Options != "" {
						// Parse PostgreSQL array format: {option1,option2,option3}
						optionsArray := parsePostgreSQLArray(*q.Options)
						if len(optionsArray) > 0 {
							answer = optionsArray[rand.Intn(len(optionsArray))]
						}
					}
				case "multi_choice":
					// Select 1-3 random options
					if q.Options != nil && *q.Options != "" {
						// Parse PostgreSQL array format: {option1,option2,option3}
						optionsArray := parsePostgreSQLArray(*q.Options)
						if len(optionsArray) > 0 {
							selected := rand.Intn(min(3, len(optionsArray))) + 1
							selectedOptions := make([]string, 0)
							used := make(map[int]bool)
							
							for j := 0; j < selected; j++ {
								for {
									idx := rand.Intn(len(optionsArray))
									if !used[idx] {
										selectedOptions = append(selectedOptions, optionsArray[idx])
										used[idx] = true
										break
									}
								}
							}
							answer = selectedOptions
						}
					}
				case "text":
					// Random feedback comments
					if isTourCompany {
						comments := []string{
							"Great tour! The guide was very knowledgeable.",
							"Enjoyed the experience, would definitely recommend.",
							"Good tour but a bit rushed at some locations.",
							"Amazing views and interesting historical facts.",
							"The guide was friendly and answered all our questions.",
							"Would have preferred a smaller group size.",
							"Excellent value for money!",
							"The meeting point was a bit hard to find.",
						}
						answer = comments[rand.Intn(len(comments))]
					} else {
						comments := []string{
							"Fast service, good quality prints.",
							"Staff was helpful with my printing needs.",
							"Prices are reasonable for the quality.",
							"Quick turnaround time, very satisfied.",
							"The self-service stations are easy to use.",
							"Would appreciate extended hours.",
							"Great for last-minute printing needs.",
							"International calling rates are competitive.",
						}
						answer = comments[rand.Intn(len(comments))]
					}
				}

				if answer != nil {
					responses = append(responses, map[string]interface{}{
						"question_id": q.ID,
						"question":    q.Text,
						"answer":      answer,
					})
				}
			}

			// Create feedback record with properly serialized JSON
			responsesJSON, err := json.Marshal(responses)
			if err != nil {
				log.Printf("Failed to marshal responses: %v\n", err)
				continue
			}

			err = db.Exec(`
				INSERT INTO feedbacks (organization_id, product_id, qr_code_id, customer_name, customer_email, overall_rating, responses, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, orgID, productID, qrID, customerName, customerEmail, overallRating, string(responsesJSON), feedbackDate, feedbackDate).Error

			if err != nil {
				log.Printf("Failed to create feedback: %v\n", err)
				continue
			}
		}
		// Feedback created successfully (silent)
	}
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper function to join string slice
func joinStringSlice(slice []string, separator string) string {
	if len(slice) == 0 {
		return ""
	}
	result := slice[0]
	for i := 1; i < len(slice); i++ {
		result += separator + slice[i]
	}
	return result
}

// Helper function to parse PostgreSQL array format: {option1,option2,option3}
func parsePostgreSQLArray(pgArray string) []string {
	if pgArray == "" || pgArray == "{}" {
		return []string{}
	}
	
	// Remove surrounding braces
	if strings.HasPrefix(pgArray, "{") && strings.HasSuffix(pgArray, "}") {
		pgArray = pgArray[1 : len(pgArray)-1]
	}
	
	// Split by commas and clean up each element
	parts := strings.Split(pgArray, ",")
	result := make([]string, 0, len(parts))
	
	for _, part := range parts {
		cleaned := strings.TrimSpace(part)
		if cleaned != "" {
			result = append(result, cleaned)
		}
	}
	
	return result
}