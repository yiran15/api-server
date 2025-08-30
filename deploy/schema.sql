-- 用户表
CREATE TABLE `users` (
  `id` BIGINT UNSIGNED PRIMARY KEY auto_increment,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `deleted_at` DATETIME,
  `name` VARCHAR(50) NOT NULL comment '用户名称',
  `nick_name` VARCHAR(50) comment '用户昵称',
  `department` VARCHAR(50) comment '部门',
  `email` VARCHAR(100) NOT NULL comment '用户邮箱',
  `password` VARCHAR(255) NOT NULL comment '用户密码',
  `avatar` VARCHAR(255) comment '用户头像',
  `mobile` VARCHAR(20) comment '用户手机号',
  `status` TINYINT(1) DEFAULT 1 comment '用户状态,1可用,2禁用,3未激活'
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

-- 用户角色多对多关联表
CREATE TABLE `user_roles` (
  `user_id` BIGINT UNSIGNED NOT NULL,
  `role_id` BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (`user_id`, `role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 接口信息表
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

-- 角色API列表多对多关联表
CREATE TABLE `role_apis` (
  `role_id` BIGINT UNSIGNED NOT NULL,
  `api_id` BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (`role_id`, `api_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- casbin 规则表
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

--- feishu_user
-- auto-generated definition
create table feishu_users
(
    uid              bigint auto_increment comment '关联users表中的用户id'
        primary key,
    created_at       datetime(3)       null,
    updated_at       datetime(3)       null,
    deleted_at       datetime(3)       null,
    avatar_big       longtext          null comment '飞书用户avatar_big',
    avatar_middle    longtext          null comment '飞书用户avatar_middle',
    avatar_thumb     longtext          null comment '飞书用户avatar_thumb',
    avatar_url       longtext          null comment '飞书用户avatar_url',
    email            longtext          null comment '飞书用户email',
    employee_no      longtext          null comment '飞书用户employee_no',
    en_name          longtext          null comment '飞书用户en_name',
    enterprise_email longtext          null comment '飞书用户enterprise_email',
    mobile           longtext          null comment '飞书用户mobile',
    name             longtext          null comment '飞书用户name',
    open_id          longtext          null comment '飞书用户open_id',
    tenant_key       longtext          null comment '飞书用户tenant_key',
    union_id         longtext          null comment '飞书用户union_id',
    user_id          varchar(191)      null comment '飞书用户ID'
);

create index idx_feishu_users_deleted_at
    on feishu_users (deleted_at);