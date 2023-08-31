CREATE TABLE "history"(
    "id" SERIAL NOT NULL,
    "id_user" INTEGER NOT NULL,
    "segment_name" TEXT NOT NULL,
    "action_type" VARCHAR(255) CHECK
        ("action_type" IN('added', 'removed')) NOT NULL,
        "action_date" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "history" ADD PRIMARY KEY("id");
CREATE TABLE "users"(
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL
);
ALTER TABLE
    "users" ADD PRIMARY KEY("id");
CREATE TABLE "users_segments"(
    "id" SERIAL NOT NULL,
    "id_user" INTEGER NOT NULL,
    "id_segment" INTEGER NOT NULL
);
ALTER TABLE
    "users_segments" ADD PRIMARY KEY("id");
CREATE TABLE "segments"(
    "id" SERIAL NOT NULL,
    "slug" TEXT NOT NULL
);
ALTER TABLE
    "segments" ADD PRIMARY KEY("id");
ALTER TABLE
    "segments" ADD CONSTRAINT "segments_slug_unique" UNIQUE("slug");
ALTER TABLE
    "users_segments" ADD CONSTRAINT "users_segments_id_segment_foreign" FOREIGN KEY("id_segment") REFERENCES "segments"("id");
ALTER TABLE
    "users_segments" ADD CONSTRAINT "users_segments_id_user_foreign" FOREIGN KEY("id_user") REFERENCES "users"("id");
