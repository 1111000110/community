CREATE TABLE user (
                         user_id BIGINT AUTO_INCREMENT,
                         password VARCHAR(255) NOT NULL DEFAULT '' COMMENT '用户密码',  -- Be cautious about using empty string as default
                         phone VARCHAR(255) NOT NULL DEFAULT '' COMMENT '手机号码',
                         gender CHAR(10) DEFAULT 'unknown' COMMENT '性别，male|female|unknown',
                         nickname VARCHAR(255) DEFAULT '默认名称' COMMENT '昵称',
                         avatar VARCHAR(255) DEFAULT '' COMMENT '头像URL',  -- Again, empty string may be better as NULL if it’s not always available
                         birth_date BIGINT DEFAULT 0 COMMENT '出生日期（时间戳）',
                         role VARCHAR(50) DEFAULT '' COMMENT '用户角色',
                         status BIGINT DEFAULT 0 COMMENT '用户状态',
                         email VARCHAR(255)  DEFAULT '' COMMENT '电子邮箱',
                         ct BIGINT DEFAULT 0  COMMENT '创建时间',
                         ut BIGINT DEFAULT 0 COMMENT '更新时间',
                         UNIQUE KEY user_id_index (user_id),  -- Changed to a more standard way of creating unique keys
                         UNIQUE KEY mobile_index (phone),
                         PRIMARY KEY (user_id)
) ENGINE=InnoDB COLLATE=utf8mb4_general_ci COMMENT='账户表';
