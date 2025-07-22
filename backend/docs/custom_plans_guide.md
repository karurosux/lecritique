# Custom Hidden Plans Guide

The subscription system now supports creating custom, hidden plans for specific customers. This allows you to offer special pricing or features without showing them publicly.

## Creating a Custom Hidden Plan

### 1. Via SQL (Recommended)

```sql
INSERT INTO subscription_plans (
    code,
    name,
    description,
    price,
    currency,
    interval,
    is_active,
    is_visible,  -- Set to false to hide from public listing
    features
) VALUES (
    'custom_bigcorp_2024',
    'Custom Enterprise - BigCorp',
    'Special negotiated plan for BigCorp with volume discount',
    499.00,
    'USD',
    'month',
    true,
    false,  -- Hidden plan
    jsonb_build_object(
        'limits', jsonb_build_object(
            'max_organizations', -1,
            'max_feedbacks_per_month', 100000,
            'max_storage_gb', 1000,
            'max_api_calls_per_hour', 100000
        ),
        'flags', jsonb_build_object(
            'advanced_analytics', true,
            'custom_branding', true,
            'api_access', true,
            'priority_support', true,
            'white_label', true,
            'custom_domain', true,
            'dedicated_account_manager', true
        ),
        'custom', jsonb_build_object(
            'contract_terms', '2 years',
            'discount_percentage', 40,
            'sla_uptime', 99.99,
            'support_response_time_hours', 1
        )
    )
);
```

### 2. Via Admin API (Future)

```go
// Handler for admin to create custom plan
func (h *AdminHandler) CreateCustomPlan(c echo.Context) error {
    // Verify admin permissions
    
    var plan models.SubscriptionPlan
    if err := c.Bind(&plan); err != nil {
        return err
    }
    
    plan.IsVisible = false  // Always hidden
    plan.Code = fmt.Sprintf("custom_%s_%d", customerSlug, time.Now().Unix())
    
    return h.planService.CreatePlan(ctx, &plan)
}
```

## Assigning Custom Plans to Users

### 1. Via Admin API

```go
// POST /api/v1/admin/subscriptions/assign-custom-plan
{
    "account_id": "123e4567-e89b-12d3-a456-426614174000",
    "plan_code": "custom_bigcorp_2024"
}
```

### 2. Direct Database Update

```sql
-- Find the user's current subscription
UPDATE subscriptions 
SET plan_id = (SELECT id FROM subscription_plans WHERE code = 'custom_bigcorp_2024')
WHERE account_id = '123e4567-e89b-12d3-a456-426614174000';
```

## Visibility Rules

1. **Public API** (`/api/v1/plans`)
   - Only returns plans where `is_visible = true`
   - Used by frontend subscription selection page

2. **Admin API** (`/api/v1/admin/plans`)
   - Returns all plans including hidden ones
   - Shows visibility status for each plan

3. **User's Current Plan**
   - Users can always see their own plan details
   - Even if the plan is hidden from public listing

## Use Cases

### 1. Volume Discounts
Create hidden plans with special pricing for large customers:
```sql
-- 50% off for customers with 10+ organizations
INSERT INTO subscription_plans (code, name, price, is_visible) 
VALUES ('enterprise_volume_50', 'Enterprise Volume 50% Off', 99.50, false);
```

### 2. Legacy Plans
Keep old customers on discontinued plans:
```sql
-- Mark old plans as hidden instead of deleting
UPDATE subscription_plans 
SET is_visible = false 
WHERE code IN ('legacy_pro', 'legacy_starter');
```

### 3. Beta Features
Test new features with select customers:
```sql
-- Hidden plan with experimental features
INSERT INTO subscription_plans (code, name, features, is_visible)
VALUES ('beta_ai_features', 'Beta AI Features', 
    '{"flags": {"ai_menu_suggestions": true, "ai_review_analysis": true}}', 
    false);
```

### 4. Partner Plans
Special plans for business partners:
```sql
-- Partner plan with revenue sharing
INSERT INTO subscription_plans (code, name, price, features, is_visible)
VALUES ('partner_reseller', 'Partner Reseller Plan', 0, 
    '{"custom": {"revenue_share_percentage": 30}}', 
    false);
```

## Best Practices

1. **Naming Convention**: Use descriptive codes like `custom_[company]_[year]`
2. **Documentation**: Keep a record of why each custom plan was created
3. **Expiration**: Set end dates in the custom metadata
4. **Migration**: Plan how to migrate customers when custom plans expire
5. **Audit Trail**: Log who assigned custom plans and when

## Security Considerations

1. **Admin Only**: Only admins should assign custom plans
2. **Validation**: Verify the plan exists before assignment
3. **Audit Logs**: Track all custom plan assignments
4. **Rate Limiting**: Prevent abuse of custom plan creation

## Example: Complete Custom Plan Workflow

```bash
# 1. Create custom plan for ACME Corp
psql -c "INSERT INTO subscription_plans (code, name, price, is_visible, features) 
         VALUES ('custom_acme_2024', 'ACME Corp Special', 299, false, 
                '{\"limits\": {\"max_organizations\": 50}}')"

# 2. Assign to ACME's account
curl -X POST http://api.lecritique.com/v1/admin/subscriptions/assign-custom-plan \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "account_id": "acme-account-uuid",
    "plan_code": "custom_acme_2024"
  }'

# 3. Verify assignment
curl http://api.lecritique.com/v1/user/subscription \
  -H "Authorization: Bearer $ACME_USER_TOKEN"
# Returns: { "plan": { "name": "ACME Corp Special", "price": 299 } }
```
