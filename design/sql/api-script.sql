CREATE DATABASE IF NOT EXISTS go_project DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE go_project;

CREATE TABLE `user` (
    id bigint AUTO_INCREMENT PRIMARY KEY,
    openid varchar(50) NOT NULL UNIQUE,
    unionid varchar(50) NOT NULL DEFAULT '',
    phone_number varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
    nickname varchar(10) NOT NULL DEFAULT '' COMMENT '昵称',
    avatar_url varchar(150) NOT NULL DEFAULT '' COMMENT '头像链接',
    create_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY (phone_number)
) COMMENT='用户信息';
