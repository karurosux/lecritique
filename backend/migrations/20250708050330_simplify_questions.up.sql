-- Add dish_id column to questions table
ALTER TABLE questions ADD COLUMN dish_id UUID;

-- Update existing questions to use dish_id from their questionnaire
UPDATE questions 
SET dish_id = q.dish_id 
FROM questionnaires q 
WHERE questions.questionnaire_id = q.id 
AND q.dish_id IS NOT NULL;

-- Make dish_id not null and add foreign key constraint
ALTER TABLE questions ALTER COLUMN dish_id SET NOT NULL;
ALTER TABLE questions ADD CONSTRAINT fk_questions_dish 
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE;

-- Add index for better performance
CREATE INDEX idx_questions_dish_id ON questions(dish_id);
CREATE INDEX idx_questions_dish_display_order ON questions(dish_id, display_order);

-- Drop the old questionnaire_id foreign key constraint and column
ALTER TABLE questions DROP CONSTRAINT IF EXISTS fk_questions_questionnaire_id;
ALTER TABLE questions DROP COLUMN questionnaire_id;