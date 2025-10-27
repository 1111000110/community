CREATE TABLE IF NOT EXISTS message (
   session_id varchar(50) NOT NULL,
   message_id bigint NOT NULL,
   send_id bigint NOT NULL,
   reply_id bigint NOT NULL,
   create_time bigint NOT NULL,
   update_time bigint NOT NULL,
   status bigint NOT NULL,
   text text NOT NULL,
   message_type bigint NOT NULL,
   addition text NOT NULL,
   PRIMARY KEY (message_id),  -- 主键设为 message_id（天然唯一）
   KEY idx_session_id (session_id)  -- session_id 作为普通索引（非唯一）
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息表';