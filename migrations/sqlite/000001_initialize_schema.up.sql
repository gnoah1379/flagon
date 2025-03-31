CREATE TABLE users (
    id TEXT NOT NULL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL,
    password TEXT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    avatar_url TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_sso (
    id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    provider_id TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE project_groups (
    id TEXT NOT NULL PRIMARY KEY,
    parent_id TEXT REFERENCES project_groups (id) ON DELETE CASCADE,
    owner_id TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE project_groups_users (
    user_id TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    group_id TEXT NOT NULL REFERENCES project_groups (id) ON DELETE CASCADE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, group_id)
);

CREATE TABLE projects (
    id TEXT NOT NULL PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    group_id TEXT NOT NULL REFERENCES project_groups (id) ON DELETE CASCADE,
    owner_id TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE projects_users (
    user_id TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    project_id TEXT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, project_id)
);

CREATE TABLE project_categories (
    id TEXT NOT NULL PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (project_id, name)
);

CREATE TABLE project_environments (
    id TEXT NOT NULL PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (project_id, name)
);

CREATE TABLE project_features (
    id TEXT NOT NULL PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    category_id TEXT NOT NULL REFERENCES project_categories (id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    default_value BOOLEAN NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (project_id, name)
);

CREATE TABLE project_target_groups (
    id TEXT NOT NULL PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    rollout_percentage INTEGER NOT NULL,
    rules TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (project_id, name)
);

CREATE TABLE project_target_group_environment (
    project_id TEXT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    target_group_id TEXT NOT NULL REFERENCES project_target_groups (id) ON DELETE CASCADE,
    environment_id TEXT NOT NULL REFERENCES project_environments (id) ON DELETE CASCADE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (project_id, target_group_id, environment_id)
);

CREATE TABLE project_feature_flags (
    project_id TEXT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    feature_id TEXT NOT NULL REFERENCES project_features (id) ON DELETE CASCADE,
    environment_id TEXT REFERENCES project_environments (id) ON DELETE CASCADE,
    target_group_id TEXT REFERENCES project_target_groups (id) ON DELETE CASCADE,
    enabled BOOLEAN NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (project_id, feature_id, environment_id, target_group_id)
);

CREATE TABLE access_tokens (
    id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    token_type TEXT NOT NULL,
    group_id TEXT REFERENCES project_groups (id) ON DELETE CASCADE,
    project_id TEXT REFERENCES projects (id) ON DELETE CASCADE,
    environment_id TEXT REFERENCES project_environments (id) ON DELETE CASCADE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, name)
);
