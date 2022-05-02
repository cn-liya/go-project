CREATE TABLE `system_action` (
    id int AUTO_INCREMENT PRIMARY KEY,
    key_name varchar(64) NOT NULL UNIQUE,
    level tinyint NOT NULL DEFAULT 2 COMMENT '1:group;2:route',
    title varchar(16) NOT NULL DEFAULT '' COMMENT '备注标题',
    sort int NOT NULL DEFAULT 0 COMMENT '排序值0~9999',
    create_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) COMMENT='系统接口';

CREATE TABLE `system_config` (
    id int AUTO_INCREMENT PRIMARY KEY,
    key_name varchar(64) NOT NULL UNIQUE,
    content json,
    create_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) COMMENT='系统菜单';
INSERT INTO `system_config` (key_name,content) VALUES
('sys_menu_keys', '[]'),
('sys_menu_trees','[]');

CREATE TABLE `admin_role` (
    id int AUTO_INCREMENT PRIMARY KEY,
    name varchar(32) NOT NULL DEFAULT '',
    actions json COMMENT '数组保存system_action.key_name',
    menus json COMMENT '数组保存system_config[sys_menu_keys]元素',
    create_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) COMMENT='管理员角色';

CREATE TABLE `admin_user` (
    id int AUTO_INCREMENT PRIMARY KEY,
    username varchar(32) NOT NULL UNIQUE,
    password varchar(64) NOT NULL DEFAULT '',
    pwd_exp int NOT NULL DEFAULT (UNIX_TIMESTAMP()+7862400) COMMENT '密码过期时间戳',
    role_id int NOT NULL DEFAULT 0 COMMENT 'admin_role.id;-1:super',
    status tinyint NOT NULL DEFAULT 1 COMMENT '-1:off;1:on',
    create_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY(role_id)
) COMMENT='管理员账号';

INSERT INTO `admin_user` (username,password,role_id) VALUES
('super','84BHqjrPHK-07b8xirALIuU7iwjQgzdov-dGqlfWegFueqaSLarejC1bN9Fw',-1),
('admin','84BHqjrPHK-07b8xirALIuU7iwjQgzdov-dGqlfWegFueqaSLarejC1bN9Fw',-1);
-- 受保护的超管账号，初始密码: 123456
