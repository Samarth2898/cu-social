CREATE TABLE "follows" (
  "FollowingUserId" integer,
  "FollowedUserId" integer
);

CREATE TABLE "users" (
  "UserID" integer PRIMARY KEY,
  "Username" varchar,
  "Password" varchar,
  "ProfilePicture" varchar,
  "Biography" text,
  "Email" varchar,
  "CreatedAt" timestamp
);

CREATE TABLE "posts" (
  "PostId" integer PRIMARY KEY,
  "Title" varchar,
  "Body" text,
  "UserID" integer,
  "Status" varchar,
  "CreatedAt" timestamp
);

COMMENT ON COLUMN "posts"."Body" IS 'Content of the post';

ALTER TABLE "posts" ADD FOREIGN KEY ("UserID") REFERENCES "users" ("UserID");

ALTER TABLE "follows" ADD FOREIGN KEY ("FollowingUserId") REFERENCES "users" ("UserID");

ALTER TABLE "follows" ADD FOREIGN KEY ("FollowedUserId") REFERENCES "users" ("UserID");
