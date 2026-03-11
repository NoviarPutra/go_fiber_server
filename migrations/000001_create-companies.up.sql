CREATE TABLE IF NOT EXISTS "companies" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    "name" text NOT NULL,
    "code" text NOT NULL,
    "logo_url" text,
    "created_at" timestamp with time zone DEFAULT now() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT now() NOT NULL,
    "deleted_at" timestamp with time zone,
    CONSTRAINT "companies_name_check" CHECK (char_length(name) >= 1)
);

CREATE UNIQUE INDEX "idx_companies_code_active_unique" 
ON "companies" ("code") 
WHERE "deleted_at" IS NULL;

CREATE INDEX "idx_companies_active_lookup" 
ON "companies" ("id") 
WHERE "deleted_at" IS NULL;
CREATE INDEX "idx_companies_deleted_at" 
ON "companies" ("deleted_at") 
WHERE "deleted_at" IS NOT NULL;