// Feature types and definitions matching the backend
export type FeatureType = 'limit' | 'flag' | 'custom';

export interface FeatureDefinition {
  key: string;
  type: FeatureType;
  displayName: string;
  description?: string;
  unit?: string;
  unlimitedText?: string;
  format?: string;
  icon?: string;
  category?: string;
  sortOrder: number;
}

// Feature keys matching backend constants
export const LIMITS = {
  RESTAURANTS: 'max_restaurants',
  QR_CODES: 'max_qr_codes',
  FEEDBACKS_PER_MONTH: 'max_feedbacks_per_month',
  TEAM_MEMBERS: 'max_team_members'
} as const;

export const FLAGS = {
  BASIC_ANALYTICS: 'basic_analytics',
  ADVANCED_ANALYTICS: 'advanced_analytics',
  FEEDBACK_EXPLORER: 'feedback_explorer',
  CUSTOM_BRANDING: 'custom_branding',
  PRIORITY_SUPPORT: 'priority_support'
} as const;

// Feature registry with all feature definitions
export const featureRegistry: Record<string, FeatureDefinition> = {
  // Limits
  [LIMITS.RESTAURANTS]: {
    key: LIMITS.RESTAURANTS,
    type: 'limit',
    displayName: 'Restaurants',
    description: 'Maximum number of restaurants',
    unit: 'restaurants',
    unlimitedText: 'Unlimited restaurants',
    format: '{value} restaurant(s)',
    icon: 'store',
    category: 'core',
    sortOrder: 1
  },
  [LIMITS.QR_CODES]: {
    key: LIMITS.QR_CODES,
    type: 'limit',
    displayName: 'QR Codes',
    description: 'Total QR codes across all restaurants',
    unit: 'QR codes',
    unlimitedText: 'Unlimited QR codes',
    format: '{value} QR codes',
    icon: 'qr-code',
    category: 'core',
    sortOrder: 2
  },
  [LIMITS.FEEDBACKS_PER_MONTH]: {
    key: LIMITS.FEEDBACKS_PER_MONTH,
    type: 'limit',
    displayName: 'Monthly Feedbacks',
    description: 'Maximum feedbacks per month',
    unit: 'feedbacks/month',
    unlimitedText: 'Unlimited feedbacks',
    format: '{value} feedbacks/month',
    icon: 'message-square',
    category: 'core',
    sortOrder: 3
  },
  [LIMITS.TEAM_MEMBERS]: {
    key: LIMITS.TEAM_MEMBERS,
    type: 'limit',
    displayName: 'Team Members',
    description: 'Maximum team members',
    unit: 'members',
    unlimitedText: 'Unlimited team members',
    format: '{value} team member(s)',
    icon: 'users',
    category: 'collaboration',
    sortOrder: 4
  },

  // Flags
  [FLAGS.BASIC_ANALYTICS]: {
    key: FLAGS.BASIC_ANALYTICS,
    type: 'flag',
    displayName: 'Basic Analytics',
    description: 'View feedback analytics and insights',
    icon: 'bar-chart-2',
    category: 'analytics',
    sortOrder: 20
  },
  [FLAGS.ADVANCED_ANALYTICS]: {
    key: FLAGS.ADVANCED_ANALYTICS,
    type: 'flag',
    displayName: 'Advanced Analytics',
    description: 'Detailed insights and reporting',
    icon: 'bar-chart',
    category: 'analytics',
    sortOrder: 21
  },
  [FLAGS.FEEDBACK_EXPLORER]: {
    key: FLAGS.FEEDBACK_EXPLORER,
    type: 'flag',
    displayName: 'Feedback Explorer',
    description: 'Browse and search all feedback',
    icon: 'search',
    category: 'analytics',
    sortOrder: 22
  },
  [FLAGS.CUSTOM_BRANDING]: {
    key: FLAGS.CUSTOM_BRANDING,
    type: 'flag',
    displayName: 'Custom Branding',
    description: 'Customize with your brand',
    icon: 'palette',
    category: 'customization',
    sortOrder: 23
  },
  [FLAGS.PRIORITY_SUPPORT]: {
    key: FLAGS.PRIORITY_SUPPORT,
    type: 'flag',
    displayName: 'Priority Support',
    description: '24/7 priority customer support',
    icon: 'headphones',
    category: 'support',
    sortOrder: 24
  }
};

// Helper functions
export function getFeatureDefinition(key: string): FeatureDefinition | undefined {
  return featureRegistry[key];
}

export function getFeaturesByCategory(category: string): FeatureDefinition[] {
  return Object.values(featureRegistry)
    .filter(def => def.category === category)
    .sort((a, b) => a.sortOrder - b.sortOrder);
}

export function formatFeatureValue(key: string, value: number | boolean): string {
  const def = featureRegistry[key];
  if (!def) return '';

  if (def.type === 'limit') {
    const limitValue = value as number;
    if (limitValue === -1) {
      return def.unlimitedText || 'Unlimited';
    }
    
    let result = def.format || '{value} {unit}';
    result = result.replace('{value}', limitValue.toLocaleString());
    result = result.replace('{unit}', def.unit || '');
    return result;
  }

  if (def.type === 'flag' && value === true) {
    return def.displayName;
  }

  return '';
}

// Get all features from a plan in display order
export function getPlanFeatures(plan: any): string[] {
  const features: string[] = [];
  
  // Add base features available for all plans
  features.push('Analytics Dashboard');
  features.push('Feedback Explorer');
  
  // Process limits
  if (plan.features?.limits) {
    Object.entries(plan.features.limits)
      .map(([key, value]) => ({
        key,
        value,
        def: featureRegistry[key]
      }))
      .filter(item => item.def)
      .sort((a, b) => a.def.sortOrder - b.def.sortOrder)
      .forEach(item => {
        const formatted = formatFeatureValue(item.key, item.value as number);
        if (formatted) features.push(formatted);
      });
  }
  
  // Process flags
  if (plan.features?.flags) {
    Object.entries(plan.features.flags)
      .map(([key, value]) => ({
        key,
        value,
        def: featureRegistry[key]
      }))
      .filter(item => item.def && item.value === true)
      .sort((a, b) => a.def.sortOrder - b.def.sortOrder)
      .forEach(item => {
        features.push(item.def.displayName);
      });
  }
  
  return features;
}