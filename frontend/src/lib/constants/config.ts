/**
 * Application-wide configuration constants
 * Update these values to change contact information across the entire application
 */

export const APP_CONFIG = {
  // Company Information
  company: {
    name: 'LeCritique',
    website: 'https://lecritique.com'
  },

  // Contact Emails
  emails: {
    support: 'support@lecritique.com',
    privacy: 'privacy@lecritique.com',
    billing: 'support@lecritique.com', // Can be different if needed
    noreply: 'noreply@lecritique.com'
  },

  // Client Storage
  localStorageKeys: {
    authToken: 'auth_token',
    authUser: 'auth_user'
  },

  // Legal
  legal: {
    termsVersion: '1.0',
    termsLastUpdated: '2024-01-15',
    privacyVersion: '1.0',
    privacyLastUpdated: '2024-01-15'
  },

  // External Links
  links: {
    github: 'https://github.com/anthropics/claude-code/issues',
    documentation: 'https://docs.lecritique.com'
  },

  // Locales related config
  locales: {
    language: 'en-US',
    defaultDateFormat: {
      year: "numeric",
      month: "short",
      day: "numeric",
    }
  }

} as const;

// Helper function to create mailto links
export function createMailtoLink(email: keyof typeof APP_CONFIG.emails, subject?: string): string {
  const address = APP_CONFIG.emails[email];
  return subject ? `mailto:${address}?subject=${encodeURIComponent(subject)}` : `mailto:${address}`;
}
