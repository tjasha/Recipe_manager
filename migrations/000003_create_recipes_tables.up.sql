CREATE TABLE nutrition (
    id            SERIAL PRIMARY KEY,
    energy        INTEGER,
    calories      INTEGER,
    fat           DOUBLE PRECISION,
    saturated_fat DOUBLE PRECISION,
    sodium        DOUBLE PRECISION,
    fiber         DOUBLE PRECISION,
    carbohydrate  DOUBLE PRECISION,
    sugar         DOUBLE PRECISION,
    protein       DOUBLE PRECISION,
    salt          DOUBLE PRECISION
);

CREATE TABLE recipe (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    description TEXT,
    portion INTEGER,
    preparation_time INTEGER,
    cooking_time INTEGER,
    nutrition_id INTEGER REFERENCES nutrition(id),
    author INTEGER REFERENCES users(id),
    published BOOLEAN,
    image_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    modified_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE instruction (
    id               SERIAL PRIMARY KEY,
    recipe_id        INTEGER NOT NULL REFERENCES recipe (id),
    step_sequence    INTEGER NOT NULL,
    step_description TEXT,
    image_url        TEXT
);

CREATE TABLE ingredient (
    id SERIAL PRIMARY KEY,
    ingredient TEXT,
    unit TEXT,
    image_url TEXT
);

CREATE TABLE recipe_ingredient (
    id SERIAL PRIMARY KEY,
    recipe_id INTEGER REFERENCES recipe(id),
    ingredient_id INTEGER REFERENCES ingredient(id),
    quantity DOUBLE PRECISION
)