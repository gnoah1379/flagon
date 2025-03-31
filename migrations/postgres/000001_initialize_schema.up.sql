CREATE TABLE "users" (
    "id" UUID NOT NULL PRIMARY KEY,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "username" VARCHAR(255) NOT NULL,
    "password" VARCHAR(255) NULL,
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255) NOT NULL,
    "avatar_url" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE "user_sso" (
    "id" UUID NOT NULL PRIMARY KEY,
    "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "provider" VARCHAR(50) NOT NULL,
    "provider_id" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE "project_groups" (
    "id" UUID NOT NULL PRIMARY KEY,
    "parent_id" UUID NULL REFERENCES "project_groups" ("id") ON DELETE CASCADE,
    "owner_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "name" VARCHAR(255) NOT NULL,
    "slug" VARCHAR(255) NOT NULL UNIQUE,
    "description" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE "project_groups_users" (
    "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "group_id" UUID NOT NULL REFERENCES "project_groups" ("id") ON DELETE CASCADE,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("user_id", "group_id")
);
CREATE TABLE "projects" (
    "id" UUID NOT NULL PRIMARY KEY,
    "slug" VARCHAR(255) NOT NULL UNIQUE,
    "group_id" UUID NOT NULL REFERENCES "project_groups" ("id") ON DELETE CASCADE,
    "owner_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE "projects_users" (
    "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "project_id" UUID NOT NULL REFERENCES "projects" ("id") ON DELETE CASCADE,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE "project_categories" (
    "id" UUID NOT NULL PRIMARY KEY,
    "project_id" UUID NOT NULL REFERENCES "projects" ("id") ON DELETE CASCADE,
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE ("project_id", "name")
);
CREATE TABLE "project_environments" (
    "id" UUID NOT NULL PRIMARY KEY,
    "project_id" UUID NOT NULL REFERENCES "projects" ("id") ON DELETE CASCADE,
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE ("project_id", "name")
);
CREATE TABLE "project_features" (
    "id" UUID NOT NULL PRIMARY KEY,
    "project_id" UUID NOT NULL REFERENCES "projects" ("id") ON DELETE CASCADE,
    "name" VARCHAR(255) NOT NULL,
    "category_id" UUID NOT NULL REFERENCES "project_categories" ("id") ON DELETE CASCADE,
    "description" TEXT NOT NULL,
    "default_value" BOOLEAN NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE ("project_id", "name")
);
CREATE TABLE "project_target_groups" (
    "id" UUID NOT NULL PRIMARY KEY,
    "project_id" UUID NOT NULL REFERENCES "projects" ("id") ON DELETE CASCADE,
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL,
    "rollout_percentage" INT NOT NULL,
    "rules" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE ("project_id", "name")
);
CREATE TABLE "project_target_group_environment" (
    "project_id" UUID NOT NULL REFERENCES "projects" ("id") ON DELETE CASCADE,
    "target_group_id" UUID NOT NULL REFERENCES "project_target_groups" ("id") ON DELETE CASCADE,
    "environment_id" UUID NOT NULL REFERENCES "project_environments" ("id") ON DELETE CASCADE,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (
        "project_id",
        "target_group_id",
        "environment_id"
    )
);
CREATE TABLE "project_feature_flags" (
    "project_id" UUID NOT NULL REFERENCES "projects" ("id") ON DELETE CASCADE,
    "feature_id" UUID NOT NULL REFERENCES "project_features" ("id") ON DELETE CASCADE,
    "environment_id" UUID NULL REFERENCES "project_environments" ("id") ON DELETE CASCADE,
    "target_group_id" UUID NULL REFERENCES "project_target_groups" ("id") ON DELETE CASCADE,
    "enabled" BOOLEAN NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (
        "project_id",
        "feature_id",
        "environment_id",
        "target_group_id"
    )
);

CREATE TABLE "access_tokens" (
    "id" UUID NOT NULL PRIMARY KEY,
    "user_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "name" VARCHAR(255) NOT NULL,
    "token" VARCHAR(255) NOT NULL UNIQUE,
    "token_type" VARCHAR(255) NOT NULL,
    "group_id" UUID REFERENCES "project_groups" ("id") ON DELETE CASCADE,
    "project_id" UUID REFERENCES "projects" ("id") ON DELETE CASCADE,
    "environment_id" UUID REFERENCES "project_environments" ("id") ON DELETE CASCADE,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE ("user_id", "name")
);