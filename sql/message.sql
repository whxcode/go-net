/*
 * 消息记录表格
 * */
CREATE TABLE `messages` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `msg_id` varchar(36) NOT NULL COMMENT '消息唯一标识',
    `sender_id` int(11) NOT NULL COMMENT '发送者ID',
    `receiver_id` int(11) NOT NULL COMMENT '接收者ID',
    `elements` json NOT NULL COMMENT '消息内容元素列表',
    `status` tinyint(1) DEFAULT '0' COMMENT '0-未读, 1-已读',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_msg_id` (`msg_id`),
    KEY `idx_sender_id` (`sender_id`),
    KEY `idx_receiver_id` (`receiver_id`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息表';
