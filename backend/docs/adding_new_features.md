# Adding New Features to Subscription Plans

With the new generic feature system, you can add new features without changing any code!

## Example 1: Adding a New Limit (e.g., Email Notifications per Month)

### 1. Update the Database

```sql
-- Add email notification limit to all plans
UPDATE subscription_plans 
SET features = jsonb_set(
    features,
    '{limits,max_emails_per_month}',
    CASE 
        WHEN code = 'starter' THEN '100'
        WHEN code = 'professional' THEN '1000'
        WHEN code = 'enterprise' THEN '-1'
        ELSE '100'
    END::jsonb
);
```

### 2. (Optional) Add to Feature Registry for Better UX

If you want nice formatting and descriptions, add to `feature_registry.go`:

```go
LimitEmailsPerMonth: {
    Key:           "max_emails_per_month",
    Type:          FeatureTypeLimit,
    DisplayName:   "Email Notifications",
    Description:   "Monthly email notifications",
    Unit:          "emails/month",
    UnlimitedText: "Unlimited email notifications",
    Format:        "{value} emails/month",
    Icon:          "mail",
    Category:      "communication",
    SortOrder:     6,
}
```

### 3. Enforce the Limit

```go
// In your email service
func (s *EmailService) SendNotification(accountID uuid.UUID) error {
    // Check limit
    canSend, reason, err := s.usageService.CanAddResource(
        ctx, 
        subscriptionID, 
        "email_notification"
    )
    if !canSend {
        return errors.New(reason)
    }
    
    // Send email
    // ...
    
    // Track usage
    s.usageService.TrackUsage(ctx, subscriptionID, "email_notification", 1)
}
```

## Example 2: Adding a New Feature Flag (e.g., SSO Support)

### 1. Update the Database

```sql
-- Add SSO feature to enterprise plan
UPDATE subscription_plans 
SET features = jsonb_set(
    features,
    '{flags,sso_enabled}',
    'true'::jsonb
)
WHERE code = 'enterprise';
```

### 2. Check the Flag in Code

```go
// In your auth handler
if subscription.Plan.Features.GetFlag("sso_enabled") {
    // Show SSO login option
}
```

## Example 3: Adding Custom Data (e.g., Webhook URLs per Plan)

```sql
-- Add custom webhook limits
UPDATE subscription_plans 
SET features = jsonb_set(
    features,
    '{custom,webhook_urls}',
    CASE 
        WHEN code = 'starter' THEN '[]'::jsonb
        WHEN code = 'professional' THEN '["https://api.lecritique.com/webhook"]'::jsonb
        WHEN code = 'enterprise' THEN '["custom"]'::jsonb
    END
);
```

## Benefits

1. **No Code Deployment** - Marketing can add/change features
2. **Instant Updates** - Changes take effect immediately
3. **A/B Testing** - Test different feature sets easily
4. **Regional Variations** - Different features per region
5. **Custom Plans** - Create one-off plans for specific customers

## Best Practices

1. **Use Descriptive Keys** - `max_emails_per_month` not `mem`
2. **Document Features** - Keep a list of all feature keys
3. **Set Defaults** - Always handle missing features gracefully
4. **Version Plans** - Keep track of plan changes over time
5. **Test Limits** - Ensure enforcement works before enabling