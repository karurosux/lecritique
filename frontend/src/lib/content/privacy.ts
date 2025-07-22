import { APP_CONFIG } from '$lib/constants/config';

export const privacyContent = {
  lastUpdated: APP_CONFIG.legal.privacyLastUpdated,
  version: APP_CONFIG.legal.privacyVersion,
  sections: [
    {
      title: 'Privacy Policy',
      content: `At Kyooar, we are committed to protecting your privacy and ensuring the security of your personal information. This Privacy Policy explains how we collect, use, and safeguard your data when you use our organization feedback management platform.`
    },
    {
      title: '1. Information We Collect',
      content: `We collect information you provide directly to us, such as:
• Account information (name, email, company name)
• Organization and menu details
• Payment information (processed securely through Stripe)
• Feedback and questionnaire responses from your customers
• Usage data and analytics about how you interact with our platform`
    },
    {
      title: '2. How We Use Your Information',
      content: `We use the information we collect to:
• Provide, maintain, and improve our services
• Process transactions and send billing information
• Send you technical notices and support messages
• Communicate with you about products, services, and promotional offers
• Monitor and analyze trends, usage, and activities
• Detect, investigate, and prevent fraudulent or illegal activities`
    },
    {
      title: '3. Customer Feedback Data',
      content: `Customer feedback collected through your QR codes is your property. We:
• Store this data securely on your behalf
• Never sell or share customer feedback with third parties
• Only access feedback data for technical support when authorized by you
• Provide you tools to export and manage all collected feedback
• Ensure feedback remains anonymous unless customers voluntarily provide contact information`
    },
    {
      title: '4. Data Security',
      content: `We implement appropriate technical and organizational measures to protect your personal information, including:
• Encryption of data in transit and at rest
• Regular security assessments and updates
• Access controls and authentication mechanisms
• Secure data centers with physical security measures
• Regular backups and disaster recovery procedures`
    },
    {
      title: '5. Third-Party Services',
      content: `We work with trusted third-party services to provide our platform:
• Stripe for payment processing
• Cloud infrastructure providers for hosting
• Analytics services to improve our platform
• AI services (Anthropic, OpenAI, Google) for generating questionnaires

These services have their own privacy policies and we recommend reviewing them.`
    },
    {
      title: '6. Data Retention',
      content: `We retain your personal information for as long as necessary to:
• Provide our services to you
• Comply with legal obligations
• Resolve disputes and enforce agreements
• Improve our services

You may request deletion of your account and associated data at any time.`
    },
    {
      title: '7. Your Rights',
      content: `You have the right to:
• Access your personal information
• Correct inaccurate data
• Request deletion of your data
• Export your data in a portable format
• Opt-out of marketing communications
• Update your communication preferences

Contact us at ${APP_CONFIG.emails.privacy} to exercise these rights.`
    },
    {
      title: '8. Cookies and Tracking',
      content: `We use cookies and similar technologies to:
• Keep you logged in
• Remember your preferences
• Analyze platform usage
• Improve user experience

You can control cookie settings through your browser preferences.`
    },
    {
      title: '9. International Data Transfers',
      content: `Your information may be transferred to and processed in countries other than your own. We ensure appropriate safeguards are in place to protect your information in accordance with this Privacy Policy.`
    },
    {
      title: '10. Children\'s Privacy',
      content: `Kyooar is not intended for children under 13 years of age. We do not knowingly collect personal information from children under 13. If we learn we have collected information from a child under 13, we will delete it.`
    },
    {
      title: '11. Changes to This Policy',
      content: `We may update this Privacy Policy from time to time. We will notify you of any changes by posting the new Privacy Policy on this page and updating the "Last Updated" date.

This privacy policy may change at the discretion of Kyooar. We encourage you to review this policy periodically for any updates.`
    },
    {
      title: '12. Contact Us',
      content: `If you have any questions about this Privacy Policy or our privacy practices, please contact us at:

Email: ${APP_CONFIG.emails.privacy}
Address: [Your Company Address]

For general support: ${APP_CONFIG.emails.support}`
    }
  ]
};
