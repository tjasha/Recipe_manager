#!/usr/bin/env bash
set -euo pipefail

if [ -z "${DB_URL:-}" ]; then
  echo "Error: DB_URL is not set."
  echo 'Example: export DB_URL="postgres://user:password@localhost:5432/mydb?sslmode=disable"'
  exit 1
fi

psql "$DB_URL" <<'SQL'
BEGIN;

-- Optional cleanup so script can be run multiple times during development
DELETE FROM recipe_ingredient;
DELETE FROM instruction;
DELETE FROM recipe;
DELETE FROM ingredient;
DELETE FROM nutrition;
DELETE FROM users;

-- Reset sequences
ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE nutrition_id_seq RESTART WITH 1;
ALTER SEQUENCE ingredient_id_seq RESTART WITH 1;
ALTER SEQUENCE recipe_id_seq RESTART WITH 1;
ALTER SEQUENCE instruction_id_seq RESTART WITH 1;
ALTER SEQUENCE recipe_ingredient_id_seq RESTART WITH 1;

-- -------------------------------------------------------------------
-- USERS
-- Two email/password users, no google_id
-- user 1 creates 1 recipe
-- user 2 creates 2 recipes
-- -------------------------------------------------------------------
INSERT INTO users (
  user_name,
  email,
  password,
  access_level,
  created_at,
  modified_at,
  state,
  google_id
) VALUES
  (
    'Emma Stone',
    'emma.stone@example.com',
    '$2a$10$demoHashForEmma123456789',
    1,
    NOW(),
    NOW(),
    'active',
    NULL
  ),
  (
    'Daniel Brooks',
    'daniel.brooks@example.com',
    '$2a$10$demoHashForDaniel987654321',
    1,
    NOW(),
    NOW(),
    'active',
    NULL
  );

-- -------------------------------------------------------------------
-- NUTRITION
-- One nutrition row per recipe
-- -------------------------------------------------------------------
INSERT INTO nutrition (
  calories,
  fat,
  sodium,
  fiber,
  carbohydrate,
  sugar,
  protein,
) VALUES
  -- 1. Creamy Tomato Pasta
  (442, 14.5, 520, 6.5, 58.0, 9.0, 14.0),

  -- 2. Chicken Rice Bowl
  (505, 16.2, 610, 5.2, 48.0, 6.0, 35.0),

  -- 3. Banana Oat Pancakes
  (348, 9.4, 280, 5.8, 49.0, 12.0, 13.0);

-- -------------------------------------------------------------------
-- INGREDIENTS
-- Shared ingredients can repeat across recipes
-- -------------------------------------------------------------------
INSERT INTO ingredient (
  ingredient,
  unit,
  image_url
) VALUES
  ('Pasta', 'g', 'https://images.unsplash.com/photo-1551183053-bf91a1d81141'),
  ('Olive oil', 'tbsp', 'https://images.unsplash.com/photo-1474979266404-7eaacbcd87c5'),
  ('Garlic', 'clove', 'https://images.unsplash.com/photo-1615477550927-6ec1f2d4cb3a'),
  ('Onion', 'piece', 'https://images.unsplash.com/photo-1508747703725-719777637510'),
  ('Canned tomatoes', 'g', 'https://images.unsplash.com/photo-1592924357228-91a4daadcfea'),
  ('Heavy cream', 'ml', 'https://images.unsplash.com/photo-1563636619-e9143da7973b'),
  ('Parmesan', 'g', 'https://images.unsplash.com/photo-1486297678162-eb2a19b0a32d'),
  ('Salt', 'tsp', 'https://images.unsplash.com/photo-1518110925495-5fe2fda0442f'),
  ('Black pepper', 'tsp', 'https://images.unsplash.com/photo-1596040033229-a9821ebd058d'),
  ('Fresh basil', 'g', 'https://images.unsplash.com/photo-1601004890684-d8cbf643f5f2'),

  ('Chicken breast', 'g', 'https://images.unsplash.com/photo-1604503468506-a8da13d82791'),
  ('Rice', 'g', 'https://images.unsplash.com/photo-1586201375761-83865001e31c'),
  ('Broccoli', 'g', 'https://images.unsplash.com/photo-1459411621453-7b03977f4bfc'),
  ('Carrot', 'piece', 'https://images.unsplash.com/photo-1447175008436-170170753d51'),
  ('Soy sauce', 'tbsp', 'https://images.unsplash.com/photo-1515003197210-e0cd71810b5f'),
  ('Honey', 'tbsp', 'https://images.unsplash.com/photo-1587049352851-8d4e89133924'),

  ('Banana', 'piece', 'https://images.unsplash.com/photo-1574226516831-e1dff420e37f'),
  ('Oats', 'g', 'https://images.unsplash.com/photo-1515543904379-3d757afe72e7'),
  ('Egg', 'piece', 'https://images.unsplash.com/photo-1506976785307-8732e854ad03'),
  ('Milk', 'ml', 'https://images.unsplash.com/photo-1550583724-b2692b85b150'),
  ('Baking powder', 'tsp', 'https://images.unsplash.com/photo-1585238342024-78d387f4a707'),
  ('Cinnamon', 'tsp', 'https://images.unsplash.com/photo-1604329760661-e71dc83f8f26'),
  ('Maple syrup', 'tbsp', 'https://images.unsplash.com/photo-1587049352846-4a222e784d38');

-- -------------------------------------------------------------------
-- RECIPES
-- -------------------------------------------------------------------
INSERT INTO recipe (
  title,
  description,
  portion,
  preparation_time,
  cooking_time,
  nutrition_id,
  created_at,
  author,
  published,
  image_url
) VALUES
  (
    'Creamy Tomato Pasta',
    'A quick and comforting pasta with a creamy tomato sauce, garlic, and fresh basil.',
    2,
    10,
    20,
    1,
    NOW(),
    1,
    TRUE,
    'https://images.unsplash.com/photo-1621996346565-e3dbc646d9a9'
  ),
  (
    'Honey Soy Chicken Rice Bowl',
    'A balanced rice bowl with pan-seared chicken, vegetables, and a sweet savory sauce.',
    2,
    15,
    20,
    2,
    NOW(),
    2,
    TRUE,
    'https://images.unsplash.com/photo-1547592180-85f173990554'
  ),
  (
    'Banana Oat Pancakes',
    'Soft pancakes made with banana and oats, perfect for a simple breakfast.',
    2,
    10,
    15,
    3,
    NOW(),
    2,
    TRUE,
    'https://images.unsplash.com/photo-1528207776546-365bb710ee93'
  );

-- -------------------------------------------------------------------
-- INSTRUCTIONS
-- Each step stored separately
-- -------------------------------------------------------------------

-- Recipe 1: Creamy Tomato Pasta
INSERT INTO instruction (recipe_id, step_sequence, step_description, image_url) VALUES
  (1, 1, 'Bring a pot of salted water to a boil and cook the pasta until al dente.', 'https://images.unsplash.com/photo-1556761223-4c4282c73f77'),
  (1, 2, 'Heat olive oil in a pan and sauté the chopped onion for 3 minutes until soft.', 'https://images.unsplash.com/photo-1511920170033-f8396924c348'),
  (1, 3, 'Add minced garlic and cook for 30 seconds until fragrant.', 'https://images.unsplash.com/photo-1625944525533-473f1c35fd3b'),
  (1, 4, 'Pour in the canned tomatoes, season with salt and black pepper, and simmer for 8 minutes.', 'https://images.unsplash.com/photo-1512621776951-a57141f2eefd'),
  (1, 5, 'Stir in the heavy cream and parmesan, then mix until the sauce is smooth.', 'https://images.unsplash.com/photo-1473093295043-cdd812d0e601'),
  (1, 6, 'Add the cooked pasta to the sauce and toss well. Finish with fresh basil before serving.', 'https://images.unsplash.com/photo-1498837167922-ddd27525d352');

-- Recipe 2: Honey Soy Chicken Rice Bowl
INSERT INTO instruction (recipe_id, step_sequence, step_description, image_url) VALUES
  (2, 1, 'Cook the rice according to package instructions and set it aside.', 'https://images.unsplash.com/photo-1512058564366-18510be2db19'),
  (2, 2, 'Cut the chicken breast into bite-sized pieces and season lightly with salt and pepper.', 'https://images.unsplash.com/photo-1604503468506-a8da13d82791'),
  (2, 3, 'Heat olive oil in a skillet and cook the chicken until golden and cooked through.', 'https://images.unsplash.com/photo-1515003197210-e0cd71810b5f'),
  (2, 4, 'Add sliced carrot and broccoli to the pan and cook for 4 to 5 minutes until tender-crisp.', 'https://images.unsplash.com/photo-1512621776951-a57141f2eefd'),
  (2, 5, 'Mix soy sauce and honey, pour over the chicken and vegetables, and stir for 1 minute.', 'https://images.unsplash.com/photo-1514517097182-48a729fb8d68'),
  (2, 6, 'Serve the chicken and vegetables over warm rice.', 'https://images.unsplash.com/photo-1547592180-85f173990554');

-- Recipe 3: Banana Oat Pancakes
INSERT INTO instruction (recipe_id, step_sequence, step_description, image_url) VALUES
  (3, 1, 'Mash the banana in a bowl until smooth.', 'https://images.unsplash.com/photo-1574226516831-e1dff420e37f'),
  (3, 2, 'Add oats, eggs, milk, baking powder, and cinnamon, then mix into a batter.', 'https://images.unsplash.com/photo-1488477181946-6428a0291777'),
  (3, 3, 'Heat a lightly oiled pan over medium heat.', 'https://images.unsplash.com/photo-1498837167922-ddd27525d352'),
  (3, 4, 'Pour small portions of batter into the pan and cook until bubbles form on the surface.', 'https://images.unsplash.com/photo-1504754524776-8f4f37790ca0'),
  (3, 5, 'Flip the pancakes and cook the other side until golden.', 'https://images.unsplash.com/photo-1528207776546-365bb710ee93'),
  (3, 6, 'Serve warm with maple syrup.', 'https://images.unsplash.com/photo-1528207776546-365bb710ee93');

-- -------------------------------------------------------------------
-- RECIPE_INGREDIENT
-- -------------------------------------------------------------------

-- Recipe 1: Creamy Tomato Pasta
INSERT INTO recipe_ingredient (recipe_id, ingredient_id, quantity) VALUES
  (1, 1, 200),   -- Pasta
  (1, 2, 1),     -- Olive oil
  (1, 3, 2),     -- Garlic
  (1, 4, 1),     -- Onion
  (1, 5, 400),   -- Canned tomatoes
  (1, 6, 100),   -- Heavy cream
  (1, 7, 40),    -- Parmesan
  (1, 8, 0.5),   -- Salt
  (1, 9, 0.25),  -- Black pepper
  (1, 10, 10);   -- Fresh basil

-- Recipe 2: Honey Soy Chicken Rice Bowl
INSERT INTO recipe_ingredient (recipe_id, ingredient_id, quantity) VALUES
  (2, 11, 300),  -- Chicken breast
  (2, 12, 150),  -- Rice
  (2, 13, 150),  -- Broccoli
  (2, 14, 1),    -- Carrot
  (2, 2, 1),     -- Olive oil
  (2, 15, 2),    -- Soy sauce
  (2, 16, 1),    -- Honey
  (2, 8, 0.25),  -- Salt
  (2, 9, 0.25);  -- Black pepper

-- Recipe 3: Banana Oat Pancakes
INSERT INTO recipe_ingredient (recipe_id, ingredient_id, quantity) VALUES
  (3, 17, 2),    -- Banana
  (3, 18, 100),  -- Oats
  (3, 19, 2),    -- Egg
  (3, 20, 120),  -- Milk
  (3, 21, 1),    -- Baking powder
  (3, 22, 0.5),  -- Cinnamon
  (3, 23, 2);    -- Maple syrup

COMMIT;
SQL

echo "Seed completed successfully."