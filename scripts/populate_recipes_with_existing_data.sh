#!/usr/bin/env bash
set -euo pipefail

if [ -z "${DB_URL:-}" ]; then
  echo "Error: DB_URL is not set."
  echo 'Example: export DB_URL="postgres://user:password@localhost:5432/mydb?sslmode=disable"'
  exit 1
fi

psql "$DB_URL" <<'SQL'
BEGIN;

DO $$
DECLARE
  v_user1_id INT;
  v_user2_id INT;

  v_recipe1_id INT;
  v_recipe2_id INT;
  v_recipe3_id INT;

  v_salmon_id INT;
  v_lemon_id INT;
  v_garlic_id INT;
  v_olive_oil_id INT;
  v_salt_id INT;
  v_black_pepper_id INT;
  v_potato_id INT;
  v_green_beans_id INT;
  v_butter_id INT;
  v_parsley_id INT;

  v_red_lentils_id INT;
  v_onion_id INT;
  v_carrot_id INT;
  v_celery_id INT;
  v_cumin_id INT;
  v_paprika_id INT;
  v_vegetable_broth_id INT;
  v_coconut_milk_id INT;
  v_spinach_id INT;

  v_greek_yogurt_id INT;
  v_strawberries_id INT;
  v_blueberries_id INT;
  v_honey_id INT;
  v_chia_seeds_id INT;
  v_granola_id INT;
  v_banana_id INT;
BEGIN
  -- =========================================================
  -- USERS
  -- =========================================================
  INSERT INTO users (
    user_name,
    email,
    password,
    access_level,
    created_at,
    modified_at,
    state,
    google_id
  )
  VALUES (
    'Sophie Turner',
    'sophie.turner@example.com',
    '$2a$10$extraSeedHashForSophie123',
    1,
    NOW(),
    NOW(),
    'active',
    NULL
  )
  RETURNING id INTO v_user1_id;

  INSERT INTO users (
    user_name,
    email,
    password,
    access_level,
    created_at,
    modified_at,
    state,
    google_id
  )
  VALUES (
    'Lucas Miller',
    'lucas.miller@example.com',
    '$2a$10$extraSeedHashForLucas456',
    1,
    NOW(),
    NOW(),
    'active',
    NULL
  )
  RETURNING id INTO v_user2_id;

  -- =========================================================
  -- INGREDIENTS: find existing or create new
  -- =========================================================

  SELECT id INTO v_salmon_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Salmon fillet') LIMIT 1;
  IF v_salmon_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Salmon fillet', 'g', 'https://images.unsplash.com/photo-1467003909585-2f8a72700288')
    RETURNING id INTO v_salmon_id;
  END IF;

  SELECT id INTO v_lemon_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Lemon') LIMIT 1;
  IF v_lemon_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Lemon', 'piece', 'https://images.unsplash.com/photo-1590502593747-42a996133562')
    RETURNING id INTO v_lemon_id;
  END IF;

  SELECT id INTO v_garlic_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Garlic') LIMIT 1;
  IF v_garlic_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Garlic', 'clove', 'https://images.unsplash.com/photo-1615477550927-6ec1f2d4cb3a')
    RETURNING id INTO v_garlic_id;
  END IF;

  SELECT id INTO v_olive_oil_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Olive oil') LIMIT 1;
  IF v_olive_oil_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Olive oil', 'tbsp', 'https://images.unsplash.com/photo-1474979266404-7eaacbcd87c5')
    RETURNING id INTO v_olive_oil_id;
  END IF;

  SELECT id INTO v_salt_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Salt') LIMIT 1;
  IF v_salt_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Salt', 'tsp', 'https://images.unsplash.com/photo-1518110925495-5fe2fda0442f')
    RETURNING id INTO v_salt_id;
  END IF;

  SELECT id INTO v_black_pepper_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Black pepper') LIMIT 1;
  IF v_black_pepper_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Black pepper', 'tsp', 'https://images.unsplash.com/photo-1596040033229-a9821ebd058d')
    RETURNING id INTO v_black_pepper_id;
  END IF;

  SELECT id INTO v_potato_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Potato') LIMIT 1;
  IF v_potato_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Potato', 'g', 'https://images.unsplash.com/photo-1518977676601-b53f82aba655')
    RETURNING id INTO v_potato_id;
  END IF;

  SELECT id INTO v_green_beans_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Green beans') LIMIT 1;
  IF v_green_beans_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Green beans', 'g', 'https://images.unsplash.com/photo-1567375698348-5d9d5ae99de0')
    RETURNING id INTO v_green_beans_id;
  END IF;

  SELECT id INTO v_butter_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Butter') LIMIT 1;
  IF v_butter_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Butter', 'g', 'https://images.unsplash.com/photo-1589985270826-4b7bb135bc9d')
    RETURNING id INTO v_butter_id;
  END IF;

  SELECT id INTO v_parsley_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Parsley') LIMIT 1;
  IF v_parsley_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Parsley', 'g', 'https://images.unsplash.com/photo-1502741338009-cac2772e18bc')
    RETURNING id INTO v_parsley_id;
  END IF;

  SELECT id INTO v_red_lentils_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Red lentils') LIMIT 1;
  IF v_red_lentils_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Red lentils', 'g', 'https://images.unsplash.com/photo-1515543904379-3d757afe72e7')
    RETURNING id INTO v_red_lentils_id;
  END IF;

  SELECT id INTO v_onion_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Onion') LIMIT 1;
  IF v_onion_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Onion', 'piece', 'https://images.unsplash.com/photo-1508747703725-719777637510')
    RETURNING id INTO v_onion_id;
  END IF;

  SELECT id INTO v_carrot_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Carrot') LIMIT 1;
  IF v_carrot_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Carrot', 'piece', 'https://images.unsplash.com/photo-1447175008436-170170753d51')
    RETURNING id INTO v_carrot_id;
  END IF;

  SELECT id INTO v_celery_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Celery') LIMIT 1;
  IF v_celery_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Celery', 'stalk', 'https://images.unsplash.com/photo-1615485925873-6b6da2d6b65e')
    RETURNING id INTO v_celery_id;
  END IF;

  SELECT id INTO v_cumin_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Cumin') LIMIT 1;
  IF v_cumin_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Cumin', 'tsp', 'https://images.unsplash.com/photo-1599909535318-67429f0b4d17')
    RETURNING id INTO v_cumin_id;
  END IF;

  SELECT id INTO v_paprika_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Paprika') LIMIT 1;
  IF v_paprika_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Paprika', 'tsp', 'https://images.unsplash.com/photo-1596040033229-a9821ebd058d')
    RETURNING id INTO v_paprika_id;
  END IF;

  SELECT id INTO v_vegetable_broth_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Vegetable broth') LIMIT 1;
  IF v_vegetable_broth_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Vegetable broth', 'ml', 'https://images.unsplash.com/photo-1547592180-85f173990554')
    RETURNING id INTO v_vegetable_broth_id;
  END IF;

  SELECT id INTO v_coconut_milk_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Coconut milk') LIMIT 1;
  IF v_coconut_milk_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Coconut milk', 'ml', 'https://images.unsplash.com/photo-1600788907416-456578634209')
    RETURNING id INTO v_coconut_milk_id;
  END IF;

  SELECT id INTO v_spinach_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Spinach') LIMIT 1;
  IF v_spinach_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Spinach', 'g', 'https://images.unsplash.com/photo-1576045057995-568f588f82fb')
    RETURNING id INTO v_spinach_id;
  END IF;

  SELECT id INTO v_greek_yogurt_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Greek yogurt') LIMIT 1;
  IF v_greek_yogurt_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Greek yogurt', 'g', 'https://images.unsplash.com/photo-1488477181946-6428a0291777')
    RETURNING id INTO v_greek_yogurt_id;
  END IF;

  SELECT id INTO v_strawberries_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Strawberries') LIMIT 1;
  IF v_strawberries_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Strawberries', 'g', 'https://images.unsplash.com/photo-1464965911861-746a04b4bca6')
    RETURNING id INTO v_strawberries_id;
  END IF;

  SELECT id INTO v_blueberries_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Blueberries') LIMIT 1;
  IF v_blueberries_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Blueberries', 'g', 'https://images.unsplash.com/photo-1498557850523-fd3d118b962e')
    RETURNING id INTO v_blueberries_id;
  END IF;

  SELECT id INTO v_honey_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Honey') LIMIT 1;
  IF v_honey_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Honey', 'tbsp', 'https://images.unsplash.com/photo-1587049352851-8d4e89133924')
    RETURNING id INTO v_honey_id;
  END IF;

  SELECT id INTO v_chia_seeds_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Chia seeds') LIMIT 1;
  IF v_chia_seeds_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Chia seeds', 'tbsp', 'https://images.unsplash.com/photo-1515543904379-3d757afe72e7')
    RETURNING id INTO v_chia_seeds_id;
  END IF;

  SELECT id INTO v_granola_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Granola') LIMIT 1;
  IF v_granola_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Granola', 'g', 'https://images.unsplash.com/photo-1515543904379-3d757afe72e7')
    RETURNING id INTO v_granola_id;
  END IF;

  SELECT id INTO v_banana_id FROM ingredient WHERE LOWER(ingredient) = LOWER('Banana') LIMIT 1;
  IF v_banana_id IS NULL THEN
    INSERT INTO ingredient (ingredient, unit, image_url)
    VALUES ('Banana', 'piece', 'https://images.unsplash.com/photo-1574226516831-e1dff420e37f')
    RETURNING id INTO v_banana_id;
  END IF;

  -- =========================================================
  -- RECIPES
  -- recipe no longer stores nutrition_id
  -- =========================================================
  INSERT INTO recipe (
    title,
    description,
    portion,
    preparation_time,
    cooking_time,
    created_at,
    modified_at,
    author,
    published,
    image_url
  )
  VALUES (
    'Lemon Garlic Salmon with Potatoes',
    'Oven-roasted salmon served with tender potatoes and green beans in a light lemon garlic butter sauce.',
    2,
    15,
    30,
    NOW(),
    NOW(),
    v_user1_id,
    TRUE,
    'https://images.unsplash.com/photo-1467003909585-2f8a72700288'
  )
  RETURNING id INTO v_recipe1_id;

  INSERT INTO recipe (
    title,
    description,
    portion,
    preparation_time,
    cooking_time,
    created_at,
    modified_at,
    author,
    published,
    image_url
  )
  VALUES (
    'Creamy Red Lentil Soup',
    'A warm and hearty lentil soup with vegetables, coconut milk, and mild spices.',
    4,
    15,
    35,
    NOW(),
    NOW(),
    v_user2_id,
    TRUE,
    'https://images.unsplash.com/photo-1547592180-85f173990554'
  )
  RETURNING id INTO v_recipe2_id;

  INSERT INTO recipe (
    title,
    description,
    portion,
    preparation_time,
    cooking_time,
    created_at,
    modified_at,
    author,
    published,
    image_url
  )
  VALUES (
    'Berry Yogurt Breakfast Bowl',
    'A fresh breakfast bowl with Greek yogurt, berries, banana, chia seeds, and crunchy granola.',
    1,
    10,
    0,
    NOW(),
    NOW(),
    v_user2_id,
    TRUE,
    'https://images.unsplash.com/photo-1490474418585-ba9bad8fd0ea'
  )
  RETURNING id INTO v_recipe3_id;

  -- =========================================================
  -- NUTRITION
  -- =========================================================
  INSERT INTO nutrition (
    recipe_id,
    calories,
    fat,
    sodium,
    fiber,
    carbohydrate,
    sugar,
    protein
  )
  VALUES (
    v_recipe1_id,
    475,
    21.0,
    430,
    5.0,
    32.0,
    4.0,
    34.0
  );

  INSERT INTO nutrition (
    recipe_id,
    calories,
    fat,
    sodium,
    fiber,
    carbohydrate,
    sugar,
    protein
  )
  VALUES (
    v_recipe2_id,
    398,
    11.0,
    540,
    14.0,
    49.0,
    8.0,
    19.0
  );

  INSERT INTO nutrition (
    recipe_id,
    calories,
    fat,
    sodium,
    fiber,
    carbohydrate,
    sugar,
    protein
  )
  VALUES (
    v_recipe3_id,
    315,
    9.0,
    170,
    7.0,
    38.0,
    16.0,
    12.0
  );

  -- =========================================================
  -- INSTRUCTIONS
  -- =========================================================

  -- Recipe 1
  INSERT INTO instruction (recipe_id, step_sequence, step_description, image_url) VALUES
    (v_recipe1_id, 1, 'Preheat the oven to 200 degrees Celsius and line a baking tray with parchment paper.', 'https://images.unsplash.com/photo-1517248135467-4c7edcad34c4'),
    (v_recipe1_id, 2, 'Cut the potatoes into small chunks and toss them with olive oil, salt, and black pepper.', 'https://images.unsplash.com/photo-1518977676601-b53f82aba655'),
    (v_recipe1_id, 3, 'Spread the potatoes on the tray and roast for 15 minutes.', 'https://images.unsplash.com/photo-1519864600265-abb23847ef2c'),
    (v_recipe1_id, 4, 'Place the salmon and green beans on the tray, then add minced garlic, lemon slices, and small pieces of butter.', 'https://images.unsplash.com/photo-1467003909585-2f8a72700288'),
    (v_recipe1_id, 5, 'Return the tray to the oven and bake for another 12 to 15 minutes until the salmon is cooked through.', 'https://images.unsplash.com/photo-1519864600265-abb23847ef2c'),
    (v_recipe1_id, 6, 'Sprinkle with chopped parsley and serve immediately.', 'https://images.unsplash.com/photo-1502741338009-cac2772e18bc');

  -- Recipe 2
  INSERT INTO instruction (recipe_id, step_sequence, step_description, image_url) VALUES
    (v_recipe2_id, 1, 'Heat olive oil in a large pot and cook the chopped onion, carrot, and celery until softened.', 'https://images.unsplash.com/photo-1512621776951-a57141f2eefd'),
    (v_recipe2_id, 2, 'Add garlic, cumin, and paprika, then stir for 30 seconds until fragrant.', 'https://images.unsplash.com/photo-1599909535318-67429f0b4d17'),
    (v_recipe2_id, 3, 'Add the red lentils and vegetable broth, then bring everything to a boil.', 'https://images.unsplash.com/photo-1547592180-85f173990554'),
    (v_recipe2_id, 4, 'Reduce the heat and simmer for about 25 minutes until the lentils are soft.', 'https://images.unsplash.com/photo-1547592180-85f173990554'),
    (v_recipe2_id, 5, 'Stir in the coconut milk and spinach, and cook for 2 more minutes.', 'https://images.unsplash.com/photo-1576045057995-568f588f82fb'),
    (v_recipe2_id, 6, 'Season with salt and black pepper, then serve warm.', 'https://images.unsplash.com/photo-1547592180-85f173990554');

  -- Recipe 3
  INSERT INTO instruction (recipe_id, step_sequence, step_description, image_url) VALUES
    (v_recipe3_id, 1, 'Spoon the Greek yogurt into a serving bowl.', 'https://images.unsplash.com/photo-1488477181946-6428a0291777'),
    (v_recipe3_id, 2, 'Top with sliced banana, strawberries, and blueberries.', 'https://images.unsplash.com/photo-1464965911861-746a04b4bca6'),
    (v_recipe3_id, 3, 'Sprinkle the chia seeds and granola over the fruit.', 'https://images.unsplash.com/photo-1515543904379-3d757afe72e7'),
    (v_recipe3_id, 4, 'Drizzle with honey just before serving.', 'https://images.unsplash.com/photo-1587049352851-8d4e89133924');

  -- =========================================================
  -- RECIPE_INGREDIENT
  -- =========================================================

  -- Recipe 1
  INSERT INTO recipe_ingredient (recipe_id, ingredient_id, quantity) VALUES
    (v_recipe1_id, v_salmon_id, 300),
    (v_recipe1_id, v_lemon_id, 1),
    (v_recipe1_id, v_garlic_id, 2),
    (v_recipe1_id, v_olive_oil_id, 1),
    (v_recipe1_id, v_salt_id, 0.5),
    (v_recipe1_id, v_black_pepper_id, 0.25),
    (v_recipe1_id, v_potato_id, 400),
    (v_recipe1_id, v_green_beans_id, 200),
    (v_recipe1_id, v_butter_id, 20),
    (v_recipe1_id, v_parsley_id, 10);

  -- Recipe 2
  INSERT INTO recipe_ingredient (recipe_id, ingredient_id, quantity) VALUES
    (v_recipe2_id, v_red_lentils_id, 200),
    (v_recipe2_id, v_onion_id, 1),
    (v_recipe2_id, v_carrot_id, 2),
    (v_recipe2_id, v_celery_id, 2),
    (v_recipe2_id, v_garlic_id, 2),
    (v_recipe2_id, v_olive_oil_id, 1),
    (v_recipe2_id, v_cumin_id, 1),
    (v_recipe2_id, v_paprika_id, 1),
    (v_recipe2_id, v_vegetable_broth_id, 800),
    (v_recipe2_id, v_coconut_milk_id, 200),
    (v_recipe2_id, v_spinach_id, 60),
    (v_recipe2_id, v_salt_id, 0.5),
    (v_recipe2_id, v_black_pepper_id, 0.25);

  -- Recipe 3
  INSERT INTO recipe_ingredient (recipe_id, ingredient_id, quantity) VALUES
    (v_recipe3_id, v_greek_yogurt_id, 200),
    (v_recipe3_id, v_strawberries_id, 80),
    (v_recipe3_id, v_blueberries_id, 60),
    (v_recipe3_id, v_banana_id, 1),
    (v_recipe3_id, v_chia_seeds_id, 1),
    (v_recipe3_id, v_granola_id, 40),
    (v_recipe3_id, v_honey_id, 1);

END $$;

COMMIT;
SQL

echo "Population completed successfully."