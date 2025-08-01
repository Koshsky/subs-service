CREATE TABLE users (
    id SERIAL PRIMARY KEY,                    -- Индекс создается автоматически (PK)
    email VARCHAR(255) UNIQUE NOT NULL,        -- Индекс создается автоматически (UNIQUE)
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE       -- Индекс НЕ нужен на старте (мало записей)
);

CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,                     -- Индекс создается автоматически (PK)
    service_name VARCHAR(255) NOT NULL,        -- Индекс НЕ нужен на старте
    price INTEGER NOT NULL CHECK (price > 0),
    user_id INTEGER NOT NULL,                  -- Исправляем тип (был UUID, должен быть INTEGER)
    start_date DATE NOT NULL,                  -- Индекс НЕ нужен на старте
    end_date DATE,                             -- Индекс НЕ нужен на старте
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,       -- Индекс НЕ нужен на старте
    FOREIGN KEY (user_id) REFERENCES users(id) -- Индекс создается автоматически в PostgreSQL
);
