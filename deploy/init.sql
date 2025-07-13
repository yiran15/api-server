-- 用户表
CREATE TABLE `users` (
  `id` BIGINT UNSIGNED PRIMARY KEY auto_increment,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `deleted_at` DATETIME,
  `name` VARCHAR(50) NOT NULL,
  `nick_name` VARCHAR(50),
  `email` VARCHAR(100) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `avatar` VARCHAR(255),
  `mobile` VARCHAR(20),
  `status` TINYINT(1) DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE INDEX `idx_users_deleted_at` ON `users` (`deleted_at`);

-- 角色表
CREATE TABLE `roles` (
  `id` BIGINT UNSIGNED PRIMARY KEY auto_increment,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `deleted_at` DATETIME,
  `name` VARCHAR(50) NOT NULL,
  `description` VARCHAR(255)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE INDEX `idx_roles_deleted_at` ON `roles` (`deleted_at`);
INSERT INTO `roles` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `description`) VALUES (1, '2025-07-19 16:15:22', '2025-07-19 16:15:22', NULL, 'admin', 'admin');

-- 用户角色多对多关联表
CREATE TABLE `user_roles` (
  `user_id` BIGINT UNSIGNED NOT NULL,
  `role_id` BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (`user_id`, `role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--- 接口信息表
CREATE TABLE `apis` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `deleted_at` DATETIME,
  `name` VARCHAR(255) NOT NULL,
  `path` VARCHAR(255) NOT NULL,
  `method` VARCHAR(10) NOT NULL,
  `description` TEXT,
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX `idx_apis_deleted_at` ON `apis` (`deleted_at`);

--- 角色API列表多对多关联表
CREATE TABLE `role_apis` (
  `role_id` BIGINT UNSIGNED NOT NULL,
  `api_id` BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (`role_id`, `api_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `casbin_rule`
(
    id    bigint unsigned primary key auto_increment,
    ptype varchar(100) null COMMENT "p or g",
    v0    varchar(100) null COMMENT "subject",
    v1    varchar(100) null COMMENT "object",
    v2    varchar(100) null COMMENT "action",
    v3    varchar(100) null COMMENT "domain",
    v4    varchar(100) null COMMENT "resource",
    v5    varchar(100) null COMMENT "effect",
    constraint idx_casbin_rule
        unique (ptype, v0, v1, v2, v3, v4, v5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
INSERT INTO `casbin_rule` (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) VALUES (1, 'p', 'admin', '*', '*', NULL, NULL, NULL);