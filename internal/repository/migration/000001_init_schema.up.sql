CREATE TABLE "users" (
                         "id" BIGSERIAL PRIMARY KEY,
                         "username" varchar UNIQUE NOT NULL,
                         "email" varchar UNIQUE NOT NULL,
                         "hashed_password" varchar NOT NULL
);

CREATE TABLE "papers" (
                          "arxiv_id" varchar PRIMARY KEY,
                          "title" varchar NOT NULL,
                          "abstract" varchar NOT NULL,
                          "authors" varchar NOT NULL,
                          "date" date NOT NULL
);

CREATE TABLE "saved_papers" (
                                "id" BIGSERIAL PRIMARY KEY,
                                "user_id" bigserial NOT NULL,
                                "paper_id" varchar NOT NULL,
                                "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "saved_papers" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

-- ALTER TABLE "saved_papers" ADD FOREIGN KEY ("paper_id") REFERENCES "papers" ("arxiv_id");

ALTER TABLE saved_papers ADD CONSTRAINT unique_user_paper_pair UNIQUE (user_id, paper_id);
