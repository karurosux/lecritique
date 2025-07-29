# Subscription Components Usage Guide

## Backend Usage

### Middleware Usage

```go
// In your route definitions
analyticsGroup := v1.Group("/analytics")
analyticsGroup.Use(subscriptionMiddleware.RequireFeature(models.FlagAdvancedAnalytics))

// For resource limits
organizationGroup.POST("/organizations",
    subscriptionMiddleware.CheckResourceLimit("organization"),
    organizationHandler.Create,
)

// Get subscription in handler
subscription, err := middleware.GetSubscriptionFromContext(c)
```

## Frontend Usage

### 1. Load Plan Features on App Start

```typescript
// In your root layout or app initialization
import { subscription } from '$lib/stores/subscription';

onMount(async () => {
  await subscription.fetchPlanFeatures();
});
```

### 2. Feature Gate Component

```svelte
<script>
  import { FeatureGate, FEATURES } from '$lib/components/subscription';
</script>

<FeatureGate feature={FEATURES.ADVANCED_ANALYTICS}>
  <!-- This content only shows if user has advanced analytics -->
  <AdvancedAnalyticsDashboard />
</FeatureGate>

<!-- With custom fallback -->
<FeatureGate feature={FEATURES.CUSTOM_BRANDING} showUpgradePrompt={false}>
  <CustomBrandingSettings />
  {#snippet fallback()}
    <p>Upgrade to Professional plan to customize branding</p>
  {/snippet}
</FeatureGate>
```

### 3. Limit Gate Component

```svelte
<script>
  import { LimitGate, LIMITS } from '$lib/components/subscription';
  let organizationCount = 5; // Get this from your data
</script>

<LimitGate limit={LIMITS.RESTAURANTS} currentCount={organizationCount}>
  <button class="btn btn-primary">Add Organization</button>
</LimitGate>
```

### 4. Plan Badge Component

```svelte
<script>
  import { PlanBadge } from '$lib/components/subscription';
</script>

<div class="flex items-center gap-2">
  <span>Current Plan:</span>
  <PlanBadge />
</div>
```

### 5. Using Store Helpers Directly

```svelte
<script>
  import {
    hasFeature,
    getLimit,
    FEATURES,
    LIMITS,
  } from '$lib/stores/subscription';

  let showAnalytics = $derived($hasFeature(FEATURES.ADVANCED_ANALYTICS));
  let maxOrganizations = $derived($getLimit(LIMITS.RESTAURANTS));
</script>

{#if showAnalytics}
  <AnalyticsSection />
{/if}

<p>You can have up to {maxOrganizations} organizations</p>
```

## Common Patterns

### Protecting Routes

```typescript
// In +page.ts
import { subscription, FEATURES } from '$lib/stores/subscription';
import { redirect } from '@sveltejs/kit';
import { get } from 'svelte/store';

export async function load() {
  const hasFeature = get(hasFeature);

  if (!hasFeature(FEATURES.ADVANCED_ANALYTICS)) {
    throw redirect(303, '/subscription?upgrade=true');
  }
}
```

### Conditional UI Elements

```svelte
<nav>
  <a href="/dashboard">Dashboard</a>

  <FeatureGate feature={FEATURES.ADVANCED_ANALYTICS} showUpgradePrompt={false}>
    <a href="/analytics">Analytics</a>
  </FeatureGate>

  <FeatureGate feature={FEATURES.API_ACCESS} showUpgradePrompt={false}>
    <a href="/api">API Settings</a>
  </FeatureGate>
</nav>
```
