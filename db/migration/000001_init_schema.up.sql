
CREATE TABLE "organisation" (
  "id" bigserial PRIMARY KEY,
  "country_code" int,
  "merchant_name" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "full_name" varchar,
  "email" varchar(255),
  "encrypted_password" varchar(255),
  "mobile" char(20),
  "organisation_id" bigint,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "location" (
  "id" bigserial PRIMARY KEY,
  "organisation_id" bigint,
  "user_id" bigint,
  "full_name" varchar(255),
  "line1" varchar(255),
  "line2" varchar(255),
  "city" varchar(255),
  "county" varchar(255),
  "country_code" char(3),
  "geo" geometry,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "bookable" (
  "id" bigserial PRIMARY KEY,
  "location_id" bigint NOT NULL,
  "capacity" bigint NOT NULL DEFAULT 1,
  "description" text,
  "default_booking_length" interval,
  "custom_booking_length" boolean NOT NULL DEFAULT false,
  "clear_time" interval,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "booking_lengths" (
  "id" bigserial PRIMARY KEY,
  "bookable_id" bigint,
  "default_booking_length" interval,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "booking" (
  "id" bigserial PRIMARY KEY,
  "bookable_id" bigint,
  "booked_by" bigint NOT NULL,
  "start_time" timestamptz NOT NULL,
  "end_time" timestamptz NOT NULL,
  "capacity" bigint NOT NULL DEFAULT 1,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "gisdata" (
  "id" bigserial PRIMARY KEY,
  "geo" geometry
);


ALTER TABLE "users" ADD FOREIGN KEY ("organisation_id") REFERENCES "organisation" ("id");

ALTER TABLE "location" ADD FOREIGN KEY ("organisation_id") REFERENCES "organisation" ("id");

ALTER TABLE "location" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "bookable" ADD FOREIGN KEY ("location_id") REFERENCES "location" ("id");

ALTER TABLE "booking_lengths" ADD FOREIGN KEY ("bookable_id") REFERENCES "bookable" ("id");

ALTER TABLE "booking" ADD FOREIGN KEY ("bookable_id") REFERENCES "bookable" ("id");

ALTER TABLE "booking" ADD FOREIGN KEY ("booked_by") REFERENCES "users" ("id");
