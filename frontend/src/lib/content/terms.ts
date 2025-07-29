import { APP_CONFIG } from '$lib/constants/config';

export const termsContent = {
  lastUpdated: APP_CONFIG.legal.termsLastUpdated,
  version: APP_CONFIG.legal.termsVersion,
  sections: [
    {
      title: 'Terms of Service',
      content: `Welcome to Kyooar. These terms of service ("Terms") govern your use of our organization feedback management platform. By using Kyooar, you agree to these Terms.`,
    },
    {
      title: '1. Acceptance of Terms',
      content: `By accessing or using Kyooar, you agree to be bound by these Terms. If you do not agree to these Terms, please do not use our services.`,
    },
    {
      title: '2. Service Description',
      content: `Kyooar is a SaaS platform that enables organizations to collect customer feedback through QR codes. The service includes product-specific questionnaires, feedback analytics, and organization management tools.`,
    },
    {
      title: '3. Account Registration',
      content: `To use Kyooar, you must create an account. You agree to provide accurate, current, and complete information during registration and to update such information to keep it accurate, current, and complete.`,
    },
    {
      title: '4. Subscription and Payment',
      content: `Kyooar offers various subscription tiers (Starter, Professional, Enterprise). Payment terms, refund policies, and subscription limits are detailed in your subscription agreement. All payments are processed through Stripe.`,
    },
    {
      title: '5. User Responsibilities',
      content: `You are responsible for maintaining the confidentiality of your account credentials and for all activities that occur under your account. You agree to notify us immediately of any unauthorized use of your account.`,
    },
    {
      title: '6. Acceptable Use',
      content: `You agree not to use Kyooar for any illegal or unauthorized purpose. You must not violate any laws in your jurisdiction, including but not limited to copyright laws.`,
    },
    {
      title: '7. Data Privacy',
      content: `We take data privacy seriously. Customer feedback data collected through your QR codes belongs to you. We do not sell or share this data with third parties without your consent. Please refer to our Privacy Policy for more details.`,
    },
    {
      title: '8. Intellectual Property',
      content: `Kyooar and its original content, features, and functionality are owned by Kyooar and are protected by international copyright, trademark, patent, trade secret, and other intellectual property laws.`,
    },
    {
      title: '9. Service Modifications',
      content: `We reserve the right to modify or discontinue, temporarily or permanently, the service (or any part thereof) with or without notice. Prices for our services are subject to change upon 30 days notice from us.`,
    },
    {
      title: '10. Limitation of Liability',
      content: `To the maximum extent permitted by law, Kyooar shall not be liable for any indirect, incidental, special, consequential, or punitive damages resulting from your use or inability to use the service.`,
    },
    {
      title: '11. Changes to Terms',
      content: `We reserve the right to modify these Terms at our sole discretion. We will notify you of any changes by posting the new Terms on this page and updating the "Last Updated" date. It is your responsibility to review these Terms periodically.

These terms and conditions may change at the discretion of Kyooar. We encourage you to review these terms periodically for any updates.`,
    },
    {
      title: '12. Termination',
      content: `We may terminate or suspend your account and bar access to the service immediately, without prior notice or liability, under our sole discretion, for any reason whatsoever and without limitation.`,
    },
    {
      title: '13. Governing Law',
      content: `These Terms shall be governed and construed in accordance with the laws of the jurisdiction in which Kyooar operates, without regard to its conflict of law provisions.`,
    },
    {
      title: '14. Contact Information',
      content: `If you have any questions about these Terms, please contact us at ${APP_CONFIG.emails.support}.`,
    },
  ],
};
