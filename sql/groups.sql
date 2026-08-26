-- PRIMARY KEY (id)  唯一索引 唯一索引；在表中不重复
-- UNIQUE KEY uk_group_user (user_id, friend_id)  唯一索引；在表中不重复
-- 群组表
   CREATE TABLE `group_chats` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `owner_id` bigint NOT NULL COMMENT '群主ID',
  `name` varchar(255) DEFAULT '' COMMENT '群名称',
  `avatar` varchar(255) DEFAULT '' COMMENT '群头像',
  `notice` varchar(255) DEFAULT '' COMMENT '群公告',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `encrypted` tinyint(1) DEFAULT '0' COMMENT '是否加密',
  `is_muted` tinyint(1) DEFAULT '0' COMMENT '是否全员禁言 0-否 1-是',
  PRIMARY KEY (`id`),
  KEY `idx_owner_id` (`owner_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='群组表'

--群成员表
  CREATE TABLE `group_members` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `group_id` bigint NOT NULL COMMENT '群ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role` tinyint(1) DEFAULT '0' COMMENT '角色 0-成员 1-管理员 2-群主',
  `joined_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `is_muted` tinyint(1) DEFAULT '0' COMMENT '该用户是否被禁言 0-否 1-是',
  `is_notify_disabled` tinyint(1) DEFAULT '0' COMMENT '该用户是否关闭该群通知 0-否 1-是',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_user` (`group_id`,`user_id`),
  KEY `idx_group_id` (`group_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='群成员表'
