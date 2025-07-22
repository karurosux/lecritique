-- Add product_id column to questions table
ALTER TABLE questions ADD COLUMN product_id UUID;

-- Update existing questions to use product_id from their questionnaire
UPDATE questions 
SET product_id = q.product_id 
FROM questionnaires q 
WHERE questions.questionnaire_id = q.id 
AND q.product_id IS NOT NULL;

-- Make product_id not null and add foreign key constraint
ALTER TABLE questions ALTER COLUMN product_id SET NOT NULL;
ALTER TABLE questions ADD CONSTRAINT fk_questions_product 
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;

-- Add index for better performance
CREATE INDEX idx_questions_product_id ON questions(product_id);
CREATE INDEX idx_questions_product_display_order ON questions(product_id, display_order);

-- Drop the old questionnaire_id foreign key constraint and column
ALTER TABLE questions DROP CONSTRAINT IF EXISTS fk_questions_questionnaire_id;
ALTER TABLE questions DROP COLUMN questionnaire_id;
