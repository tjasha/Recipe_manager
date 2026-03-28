BEGIN;

-- reversing connection between recipe and nutrition
-- create new keys, migrating existing data, remove old keys and columns, create constrains
ALTER TABLE nutrition
    ADD COLUMN recipe_id INTEGER;

UPDATE nutrition n
SET recipe_id = r.id
    FROM recipe r
WHERE r.nutrition_id = n.id;

ALTER TABLE nutrition
    ALTER COLUMN recipe_id SET NOT NULL;

ALTER TABLE nutrition
    ADD CONSTRAINT nutrition_recipe_id_key UNIQUE (recipe_id);

ALTER TABLE nutrition
    ADD CONSTRAINT nutrition_recipe_id_fkey
        FOREIGN KEY (recipe_id)
            REFERENCES recipe(id)
            ON DELETE CASCADE;

ALTER TABLE recipe
DROP CONSTRAINT IF EXISTS recipe_nutrition_id_fkey;

ALTER TABLE recipe
DROP COLUMN IF EXISTS nutrition_id;

-- if author deleted, set it to null in recipe
ALTER TABLE recipe
DROP CONSTRAINT IF EXISTS recipe_author_fkey;

ALTER TABLE recipe
    ADD CONSTRAINT recipe_author_fkey
        FOREIGN KEY (author)
            REFERENCES users(id)
            ON DELETE SET NULL;

-- instruction.recipe_id ON DELETE CASCADE
ALTER TABLE instruction
DROP CONSTRAINT IF EXISTS instruction_recipe_id_fkey;

ALTER TABLE instruction
    ADD CONSTRAINT instruction_recipe_id_fkey
        FOREIGN KEY (recipe_id)
            REFERENCES recipe(id)
            ON DELETE CASCADE;

-- recipe_ingredient.recipe_id  ON DELETE CASCADE
ALTER TABLE recipe_ingredient
DROP CONSTRAINT IF EXISTS recipe_ingredient_recipe_id_fkey;

ALTER TABLE recipe_ingredient
    ADD CONSTRAINT recipe_ingredient_recipe_id_fkey
        FOREIGN KEY (recipe_id)
            REFERENCES recipe(id)
            ON DELETE CASCADE;

-- recipe_ingredient should always have recipe_id and ingredient_id
ALTER TABLE recipe_ingredient
    ALTER COLUMN recipe_id SET NOT NULL;

ALTER TABLE recipe_ingredient
    ALTER COLUMN ingredient_id SET NOT NULL;

-- all steps in the recipe are unique
ALTER TABLE instruction
    ADD CONSTRAINT instruction_recipe_id_step_sequence_key
        UNIQUE (recipe_id, step_sequence);

-- all ingredients inside the recipe are unique
ALTER TABLE recipe_ingredient
    ADD CONSTRAINT recipe_ingredient_recipe_id_ingredient_id_key
        UNIQUE (recipe_id, ingredient_id);

COMMIT;