CREATE TABLE "follows" (
  "following_user_id" integer,
  "followed_user_id" integer,
  PRIMARY KEY ("following_user_id", "followed_user_id")
);

CREATE TABLE "users" (
  "user_id" SERIAL PRIMARY KEY,
  "username" varchar,
  "password" varchar,
  "profile_picture" varchar,
  "biography" text,
  "email" varchar,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "posts" (
  "post_id" SERIAL PRIMARY KEY,
  "title" varchar,
  "description" varchar,
  "video_url" text,
  "user_id" integer,
  "status" varchar,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "follows" ADD FOREIGN KEY ("following_user_id") REFERENCES "users" ("user_id");

ALTER TABLE "follows" ADD FOREIGN KEY ("followed_user_id") REFERENCES "users" ("user_id");

INSERT INTO Users (username, email, password, profile_picture, biography)
SELECT
    'User' || generate_series, -- Username
    'user' || generate_series || '@example.com', -- Email
    '$2a$10$jLO/HQ5ECcpmW2HxVlfiK.u0S6Oe58hsNeOIn5XInyXFg61c8kfxG', -- Password
    'https://storage.googleapis.com/cusocial/download.png', -- ProfilePicture
    'Biography for user ' || generate_series -- Biography
FROM generate_series(1, 20);


INSERT INTO Posts (user_id, title, description, status, video_url)
SELECT 
    user_id,
    'Title for post ' || generate_series,
    'Description for post ' || generate_series,
    'Uploaded',
    'https://storage.googleapis.com/cusocial/ForBiggerBlazes.mp4'
FROM Users
CROSS JOIN generate_series(1, 20);

INSERT INTO follows (following_user_id, followed_user_id)
SELECT 
    U1.user_id,
    U2.user_id
FROM Users U1
CROSS JOIN Users U2
WHERE U1.user_id <> U2.user_id -- To avoid users following themselves
LIMIT 20; -- Limit to 20 relationships
