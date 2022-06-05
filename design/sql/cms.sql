CREATE TABLE `admin` (
    id bigint AUTO_INCREMENT PRIMARY KEY,
    username varchar(32) NOT NULL UNIQUE,
    password varchar(32) NOT NULL DEFAULT '',
    authority json,
    status tinyint NOT NULL DEFAULT 1 COMMENT 'off(-1),on(1)',
    create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci AUTO_INCREMENT=100 COMMENT='管理员账号';

INSERT INTO `admin`(id,username,password,authority,status)
VALUES (1,'admin','-M2wRJXe1HYVJY-dxqP0cH_SQFQ0_vw8','{}',1); -- admin 初始密码: 123456
