CREATE TABLE "users" (
    "user_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "username" VARCHAR(128) UNIQUE,
    "password" VARCHAR(128),
    "offset" INTEGER DEFAULT 0,
    "multiplier" INTEGER DEFAULT 1, 
    "color" INTEGER DEFAULT 0,
    "update_at" INTEGER DEFAULT 0
);