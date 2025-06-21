-- Insert subscription plans
INSERT INTO subscription_plans (id, name, code, description, price, currency, interval, features, is_active) VALUES
(
    '11111111-1111-1111-1111-111111111111',
    'Starter',
    'starter',
    'Perfect for small restaurants getting started',
    29.00,
    'USD',
    'month',
    '{
        "max_restaurants": 1,
        "max_locations_per_restaurant": 1,
        "max_qr_codes_per_location": 10,
        "max_feedbacks_per_month": 500,
        "max_team_members": 2,
        "advanced_analytics": false,
        "custom_branding": false,
        "api_access": false,
        "priority_support": false
    }'::jsonb,
    true
),
(
    '22222222-2222-2222-2222-222222222222',
    'Professional',
    'professional',
    'For growing restaurants with multiple locations',
    79.00,
    'USD',
    'month',
    '{
        "max_restaurants": 3,
        "max_locations_per_restaurant": 3,
        "max_qr_codes_per_location": 50,
        "max_feedbacks_per_month": 2000,
        "max_team_members": 5,
        "advanced_analytics": true,
        "custom_branding": false,
        "api_access": false,
        "priority_support": false
    }'::jsonb,
    true
),
(
    '33333333-3333-3333-3333-333333333333',
    'Enterprise',
    'enterprise',
    'For restaurant chains and large operations',
    199.00,
    'USD',
    'month',
    '{
        "max_restaurants": -1,
        "max_locations_per_restaurant": -1,
        "max_qr_codes_per_location": -1,
        "max_feedbacks_per_month": -1,
        "max_team_members": -1,
        "advanced_analytics": true,
        "custom_branding": true,
        "api_access": true,
        "priority_support": true
    }'::jsonb,
    true
);

-- Insert question templates
INSERT INTO question_templates (category, name, description, text, type, options, min_value, max_value, tags) VALUES
('General', 'Overall Rating', 'Overall satisfaction rating', 'How would you rate your overall experience?', 'rating', NULL, 1, 5, '{"satisfaction", "general"}'),
('Service', 'Service Speed', 'Speed of service rating', 'How would you rate the speed of service?', 'rating', NULL, 1, 5, '{"service", "speed"}'),
('Service', 'Staff Friendliness', 'Staff friendliness rating', 'How friendly was our staff?', 'scale', NULL, 1, 10, '{"service", "staff"}'),
('Food', 'Food Temperature', 'Temperature of the dish', 'Was your food served at the right temperature?', 'single_choice', '{"Too cold", "Just right", "Too hot"}', NULL, NULL, '{"food", "temperature"}'),
('Food', 'Portion Size', 'Portion size satisfaction', 'How satisfied were you with the portion size?', 'single_choice', '{"Too small", "Just right", "Too large"}', NULL, NULL, '{"food", "portion"}'),
('Food', 'Taste Rating', 'Taste satisfaction', 'How would you rate the taste of your dish?', 'rating', NULL, 1, 5, '{"food", "taste"}'),
('Food', 'Presentation', 'Food presentation rating', 'How would you rate the presentation of your dish?', 'rating', NULL, 1, 5, '{"food", "presentation"}'),
('Ambiance', 'Cleanliness', 'Restaurant cleanliness', 'How clean was the restaurant?', 'scale', NULL, 1, 10, '{"ambiance", "cleanliness"}'),
('Ambiance', 'Noise Level', 'Noise level comfort', 'How comfortable was the noise level?', 'single_choice', '{"Too quiet", "Just right", "Too loud"}', NULL, NULL, '{"ambiance", "noise"}'),
('Value', 'Value for Money', 'Price to quality ratio', 'How would you rate the value for money?', 'rating', NULL, 1, 5, '{"value", "price"}'),
('Recommendation', 'Would Recommend', 'Likelihood to recommend', 'Would you recommend us to a friend?', 'yes_no', NULL, NULL, NULL, '{"recommendation"}'),
('Recommendation', 'Return Visit', 'Likelihood to return', 'How likely are you to visit us again?', 'scale', NULL, 1, 10, '{"recommendation", "loyalty"}'),
('Feedback', 'Additional Comments', 'Open feedback', 'Any additional comments or suggestions?', 'text', NULL, NULL, NULL, '{"feedback", "suggestions"}'),
('Dietary', 'Dietary Accommodations', 'Dietary needs satisfaction', 'If you have dietary restrictions, how well were they accommodated?', 'single_choice', '{"Not applicable", "Poorly", "Adequately", "Well", "Excellently"}', NULL, NULL, '{"dietary", "accommodations"}');
