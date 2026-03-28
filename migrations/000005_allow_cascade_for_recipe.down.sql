BEGIN;

ALTER TABLE recipe_ingredient
DROP CONSTRAINT IF EXISTS recipe_ingredient_recipe_id_ingredient_id_key;

ALTER TABLE instruction
DROP CONSTRAINT IF EXISTS instruction_recipe_id_step_sequence_key;

ALTER TABLE recipe
    ADD COLUMN nutrition_id INTEGER;

UPDATE recipe r
SET nutrition_id = n.id
    FROM nutrition n
WHERE n.recipe_id = r.id;

ALTER TABLE recipe
    ADD CONSTRAINT recipe_nutrition_id_fkey
        FOREIGN KEY (nutrition_id)
            REFERENCES nutrition(id);

ALTER TABLE nutrition
DROP CONSTRAINT IF EXISTS nutrition_recipe_id_fkey;

ALTER TABLE nutrition
DROP CONSTRAINT IF EXISTS nutrition_recipe_id_key;

ALTER TABLE nutrition
DROP COLUMN IF EXISTS recipe_id;

ALTER TABLE recipe
DROP CONSTRAINT IF EXISTS recipe_author_fkey;

ALTER TABLE recipe
    ADD CONSTRAINT recipe_author_fkey
        FOREIGN KEY (author)
            REFERENCES users(id);

ALTER TABLE instruction
DROP CONSTRAINT IF EXISTS instruction_recipe_id_fkey;

ALTER TABLE instruction
    ADD CONSTRAINT instruction_recipe_id_fkey
        FOREIGN KEY (recipe_id)
            REFERENCES recipe(id);

ALTER TABLE recipe_ingredient
DROP CONSTRAINT IF EXISTS recipe_ingredient_recipe_id_fkey;

ALTER TABLE recipe_ingredient
    ADD CONSTRAINT recipe_ingredient_recipe_id_fkey
        FOREIGN KEY (recipe_id)
            REFERENCES recipe(id);

COMMIT;