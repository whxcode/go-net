-- MySQL dump 10.13  Distrib 8.0.46, for Linux (x86_64)
--
-- Host: localhost    Database: go-net
-- ------------------------------------------------------
-- Server version	8.0.46-0ubuntu0.24.04.3

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `go-net`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `go-net` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `go-net`;

--
-- Table structure for table `friends`
--

DROP TABLE IF EXISTS `friends`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `friends` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL COMMENT '用户id',
  `friend_id` int NOT NULL COMMENT '好友id',
  `status` tinyint(1) DEFAULT '0' COMMENT '0-待确认, 1-已确认, 2-已拒绝, 3-已删除',
  `remark` varchar(50) DEFAULT '' COMMENT '备注名',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_friend` (`user_id`,`friend_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_friend_id` (`friend_id`),
  CONSTRAINT `fk_commenets_friend_id` FOREIGN KEY (`friend_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_commenets_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `friends`
--

LOCK TABLES `friends` WRITE;
/*!40000 ALTER TABLE `friends` DISABLE KEYS */;
INSERT INTO `friends` VALUES (8,4,2,1,'','2026-08-23 21:02:40','2026-08-26 21:53:53'),(9,4,5,1,'','2026-08-30 22:19:32','2026-08-30 22:24:55'),(10,2,5,1,'','2026-09-06 17:09:31','2026-09-06 17:09:57');
/*!40000 ALTER TABLE `friends` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `group_chats`
--

DROP TABLE IF EXISTS `group_chats`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
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
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='群组表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `group_chats`
--

LOCK TABLES `group_chats` WRITE;
/*!40000 ALTER TABLE `group_chats` DISABLE KEYS */;
INSERT INTO `group_chats` VALUES (1,4,'红警战队','a0212e9a5da9596268f8db145f76dff96524f9ff6cfafd33f6a91c2eceaa5033','今天晚上打红警！！！','2026-08-26 22:15:09','2026-08-30 17:24:49',0,0),(3,4,'熊猫战队','ac2fa8fd352ef7ec7a2845ee49226c40142c18b97dd539b76ab72ffff9ebe54e','今天打CS','2026-08-27 22:13:46','2026-08-30 18:01:32',0,0);
/*!40000 ALTER TABLE `group_chats` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `group_members`
--

DROP TABLE IF EXISTS `group_members`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
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
  `status` tinyint(1) DEFAULT '0' COMMENT '该群员状态 0-正常 1-退群(自己退、被踢出)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_user` (`group_id`,`user_id`),
  KEY `idx_group_id` (`group_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='群成员表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `group_members`
--

LOCK TABLES `group_members` WRITE;
/*!40000 ALTER TABLE `group_members` DISABLE KEYS */;
INSERT INTO `group_members` VALUES (1,1,4,0,'2026-08-26 22:16:33','2026-08-26 22:16:33','2026-08-30 17:25:09',0,0,0),(3,1,2,0,'2026-08-26 22:18:30','2026-08-26 22:18:30','2026-08-30 17:13:50',0,0,0),(5,3,4,0,'2026-08-27 22:14:51','2026-08-27 22:14:51','2026-08-30 17:59:05',0,0,0),(6,3,2,0,'2026-08-29 23:00:50','2026-08-29 23:00:50','2026-08-29 23:00:50',0,0,0),(15,1,5,0,'2026-08-30 22:25:37','2026-08-30 22:25:37','2026-08-30 22:25:37',0,0,0);
/*!40000 ALTER TABLE `group_members` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `messages`
--

DROP TABLE IF EXISTS `messages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `messages` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `sender_id` int NOT NULL COMMENT '发送者ID',
  `receiver_id` int NOT NULL COMMENT '接收者ID',
  `elements` json NOT NULL COMMENT '消息内容元素列表',
  `status` tinyint(1) DEFAULT '0' COMMENT '0-未读, 1-已读',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `type` tinyint(1) DEFAULT '1' COMMENT '通信渠道: 0 好友通信; 1 群组通信',
  PRIMARY KEY (`id`),
  KEY `idx_sender_id` (`sender_id`),
  KEY `idx_receiver_id` (`receiver_id`)
) ENGINE=InnoDB AUTO_INCREMENT=115 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='消息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `messages`
--

LOCK TABLES `messages` WRITE;
/*!40000 ALTER TABLE `messages` DISABLE KEYS */;
INSERT INTO `messages` VALUES (1,1,2,'[{\"type\": 1, \"content\": \"1234\"}]',0,'2026-08-14 22:06:50',1),(3,1,2,'[{\"type\": 1, \"content\": \"1234\"}]',0,'2026-08-14 22:15:18',1),(4,1,2,'[{\"type\": 1, \"content\": \"1234\"}]',0,'2026-08-14 22:15:28',1),(5,1,2,'[{\"type\": 1, \"content\": \"123424\"}]',0,'2026-08-14 22:15:30',1),(6,1,2,'[{\"type\": 1, \"content\": \"1234214\"}]',0,'2026-08-14 22:15:32',1),(7,1,2,'[{\"type\": 1, \"content\": \"哈哈哈哈哈哈哈；你在干嘛\"}]',0,'2026-08-14 22:15:48',1),(8,1,2,'[{\"type\": 1, \"content\": \"1\"}]',0,'2026-08-14 22:33:11',1),(9,1,2,'[{\"type\": 1, \"content\": \"2\"}]',0,'2026-08-14 22:33:12',1),(10,1,2,'[{\"type\": 1, \"content\": \"3\"}]',0,'2026-08-14 22:33:13',1),(11,2,1,'[{\"type\": 1, \"content\": \"1234\"}]',0,'2026-08-14 22:33:17',1),(12,2,1,'[{\"type\": 1, \"content\": \"1234\"}]',0,'2026-08-14 22:33:20',1),(13,2,1,'[{\"type\": 1, \"content\": \"1234\"}]',0,'2026-08-14 22:33:21',1),(14,2,1,'[{\"type\": 1, \"content\": \"1234\"}]',0,'2026-08-14 22:33:22',1),(15,2,1,'[{\"type\": 1, \"content\": \"13241\"}]',0,'2026-08-14 22:33:22',1),(16,2,1,'[{\"type\": 1, \"content\": \"234\"}]',0,'2026-08-14 22:33:22',1),(17,2,1,'[{\"type\": 1, \"content\": \"123\"}]',0,'2026-08-14 22:33:22',1),(18,2,1,'[{\"type\": 1, \"content\": \"4123\"}]',0,'2026-08-14 22:33:23',1),(19,2,1,'[{\"type\": 1, \"content\": \"4\"}]',0,'2026-08-14 22:33:23',1),(20,2,1,'[{\"type\": 1, \"content\": \"123\"}]',0,'2026-08-14 22:33:23',1),(21,1,2,'[{\"type\": 1, \"content\": \"还会\"}]',0,'2026-08-14 22:33:29',1),(22,2,1,'[{\"type\": 1, \"content\": \"4231123423\"}]',0,'2026-08-14 22:33:31',1),(23,1,2,'[{\"type\": 1, \"content\": \"2134214\"}]',0,'2026-08-20 21:45:31',1),(24,2,1,'[{\"type\": 1, \"content\": \"1234324\"}]',0,'2026-08-20 21:45:32',1),(25,1,2,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-20 22:03:24',1),(26,2,1,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-20 22:03:47',1),(27,1,2,'[{\"type\": 0, \"content\": \"123424\"}]',0,'2026-08-20 22:07:43',1),(28,2,1,'[{\"type\": 0, \"content\": \"1234124\"}]',0,'2026-08-20 22:07:52',1),(29,1,2,'[{\"type\": 0, \"content\": \"2314324\"}]',0,'2026-08-20 22:07:55',1),(30,2,1,'[{\"type\": 0, \"content\": \"？？\"}]',0,'2026-08-20 22:08:03',1),(31,2,1,'[{\"type\": 0, \"content\": \"😄\"}]',0,'2026-08-20 22:08:08',1),(32,2,1,'[{\"type\": 0, \"content\": \"哈哈哈\"}]',0,'2026-08-20 22:08:13',1),(33,2,1,'[{\"type\": 0, \"content\": \"💚\"}]',0,'2026-08-20 22:08:35',1),(34,1,2,'[{\"type\": 0, \"content\": \"😀\"}]',0,'2026-08-20 22:08:50',1),(35,1,2,'[{\"type\": 0, \"content\": \"😄\"}]',0,'2026-08-20 22:08:52',1),(36,1,2,'[{\"type\": 0, \"content\": \"🤩\"}]',0,'2026-08-20 22:08:54',1),(37,1,2,'[{\"type\": 0, \"content\": \"☺️\"}]',0,'2026-08-20 22:09:00',1),(38,1,2,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-20 22:19:50',1),(39,2,1,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-20 22:19:54',1),(40,1,2,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-20 22:21:54',1),(41,2,1,'[{\"type\": 0, \"content\": \"1324\"}]',0,'2026-08-20 22:21:58',1),(42,2,1,'[{\"type\": 0, \"content\": \"🍅\"}]',0,'2026-08-20 22:22:04',1),(43,2,1,'[{\"type\": 0, \"content\": \"🍥🥨🍚🌯\"}]',0,'2026-08-20 22:22:13',1),(44,1,2,'[{\"type\": 0, \"content\": \"🫏🫏\"}]',0,'2026-08-20 22:22:19',1),(45,2,1,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-20 22:22:48',1),(46,2,1,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-20 22:22:51',1),(47,1,2,'[{\"type\": 0, \"content\": \"🦄\"}]',0,'2026-08-20 22:22:54',1),(48,4,2,'[{\"type\": 0, \"content\": \"13241234\"}]',0,'2026-08-23 20:50:48',1),(49,2,2,'[{\"type\": 0, \"content\": \"12343214\"}]',0,'2026-08-23 21:31:26',1),(50,4,2,'[{\"type\": 0, \"content\": \"21343124\"}]',0,'2026-08-23 21:31:32',1),(51,2,4,'[{\"type\": 0, \"content\": \"1234234\"}]',0,'2026-08-23 21:31:39',1),(52,2,4,'[{\"type\": 0, \"content\": \"😁\"}]',0,'2026-08-23 21:31:49',1),(53,2,4,'[{\"type\": 0, \"content\": \"肺癌无法哇\"}]',0,'2026-08-23 21:32:53',1),(54,4,2,'[{\"type\": 0, \"content\": \"😀🤩\"}]',0,'2026-08-23 21:32:59',1),(55,2,4,'[{\"type\": 0, \"content\": \"/api/file/75bc2ecde322ddddd4eb2d5b0f79ceed3c76b3e5eca8dae54a8cfe2ef6cbaa3b\"}]',0,'2026-08-23 21:33:10',1),(56,2,4,'[{\"type\": 0, \"content\": \"1\"}]',0,'2026-08-23 21:33:53',1),(57,2,4,'[{\"type\": 0, \"content\": \"/api/file/75bc2ecde322ddddd4eb2d5b0f79ceed3c76b3e5eca8dae54a8cfe2ef6cbaa3b\"}]',0,'2026-08-23 21:33:59',1),(58,2,4,'[{\"hash\": \"75bc2ecde322ddddd4eb2d5b0f79ceed3c76b3e5eca8dae54a8cfe2ef6cbaa3b\", \"name\": \"Snipaste_2026-08-22_13-28-24.png\", \"size\": 71560, \"type\": 1, \"content\": \"/api/file/75bc2ecde322ddddd4eb2d5b0f79ceed3c76b3e5eca8dae54a8cfe2ef6cbaa3b\"}]',0,'2026-08-23 22:24:52',1),(61,4,2,'[{\"hash\": \"247e16352363b2443bf98770801ad129c61d50e91007f13710e92514877abcd0\", \"name\": \"袁境莲-13559766695(1).pdf\", \"size\": 382508, \"type\": 3, \"content\": \"{\\\"url\\\":\\\"/api/file/247e16352363b2443bf98770801ad129c61d50e91007f13710e92514877abcd0\\\",\\\"fileName\\\":\\\"袁境莲-13559766695(1).pdf\\\",\\\"fileSize\\\":382508,\\\"fileType\\\":\\\"application/pdf\\\"}\"}]',0,'2026-08-30 18:14:42',1),(62,4,1,'[{\"type\": 0, \"content\": \"11\"}]',0,'2026-08-30 22:00:04',1),(63,4,1,'[{\"type\": 0, \"content\": \"123\"}]',0,'2026-08-30 22:03:59',1),(64,4,1,'[{\"type\": 0, \"content\": \"12123\"}]',0,'2026-08-30 22:04:39',1),(65,4,1,'[{\"type\": 0, \"content\": \"12341234\"}]',0,'2026-08-30 22:05:24',1),(66,4,1,'[{\"type\": 0, \"content\": \"1234124\"}]',0,'2026-08-30 22:05:28',1),(67,4,1,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-30 22:05:41',1),(68,4,1,'[{\"type\": 0, \"content\": \"1324\"}]',0,'2026-08-30 22:06:07',1),(69,4,1,'[{\"type\": 0, \"content\": \"1234324\"}]',0,'2026-08-30 22:06:19',1),(70,4,1,'[{\"type\": 0, \"content\": \"12341\"}]',0,'2026-08-30 22:06:26',1),(71,4,1,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-30 22:06:48',1),(72,4,1,'[{\"type\": 0, \"content\": \"12342343\"}]',0,'2026-08-30 22:06:52',1),(73,4,1,'[{\"type\": 0, \"content\": \"21341234\"}]',0,'2026-08-30 22:06:54',1),(74,4,1,'[{\"type\": 0, \"content\": \"2134324\"}]',0,'2026-08-30 22:06:58',1),(75,4,1,'[{\"type\": 0, \"content\": \"123423\"}]',0,'2026-08-30 22:07:06',1),(76,4,1,'[{\"type\": 0, \"content\": \"22222\"}]',0,'2026-08-30 22:13:26',1),(77,2,1,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-30 22:18:14',1),(78,4,1,'[{\"type\": 0, \"content\": \"234234\"}]',0,'2026-08-30 22:18:19',1),(79,2,1,'[{\"type\": 0, \"content\": \"2134234234\"}]',0,'2026-08-30 22:18:22',1),(80,4,1,'[{\"type\": 0, \"content\": \"嗷嗷\"}]',0,'2026-08-30 22:18:26',1),(81,2,1,'[{\"type\": 0, \"content\": \"哈哈哈\"}]',0,'2026-08-30 22:18:30',1),(82,5,4,'[{\"type\": 0, \"content\": \"你好 我是王杰瑞\"}]',0,'2026-08-30 22:25:10',0),(83,2,1,'[{\"type\": 0, \"content\": \"欢迎新朋友到来\"}]',0,'2026-08-30 22:26:00',1),(84,5,1,'[{\"type\": 0, \"content\": \"哈哈 谢谢大家。我刚刚来\"}]',0,'2026-08-30 22:26:14',1),(85,5,1,'[{\"type\": 0, \"content\": \"我爱吃鸡蛋饼\"}]',0,'2026-08-30 22:26:44',1),(86,5,1,'[{\"hash\": \"b521d19a604089c333a8b4496a3f20fd9d83ab38a6b081b2697b29cff68cfdeb\", \"name\": \"1000016672.mp4\", \"size\": 26714263, \"type\": 2, \"content\": \"{\\\"url\\\":\\\"/api/file/b521d19a604089c333a8b4496a3f20fd9d83ab38a6b081b2697b29cff68cfdeb\\\",\\\"fileName\\\":\\\"1000016672.mp4\\\",\\\"fileSize\\\":26714263,\\\"fileType\\\":\\\"video/mp4\\\"}\"}]',0,'2026-08-30 22:26:52',1),(87,5,1,'[{\"hash\": \"be3836114fb69f32098d02c17cfda7c69e775f0d5b6574ca01f5df359161b740\", \"name\": \"1788100026990.jpg\", \"size\": 6358240, \"type\": 1, \"content\": \"/api/file/be3836114fb69f32098d02c17cfda7c69e775f0d5b6574ca01f5df359161b740\"}]',0,'2026-08-30 22:27:08',1),(88,5,1,'[{\"type\": 0, \"content\": \"🤚\"}]',0,'2026-08-30 22:27:49',1),(89,4,1,'[{\"type\": 0, \"content\": \"王杰瑞；你学python没？\"}]',0,'2026-08-30 22:28:46',1),(90,4,1,'[{\"type\": 0, \"content\": \"王杰瑞；你学python没？\"}]',0,'2026-08-30 22:29:07',1),(91,4,1,'[{\"type\": 0, \"content\": \"！！！\"}]',0,'2026-08-31 21:09:16',1),(92,5,1,'[{\"type\": 0, \"content\": \"你咋干嘛\"}]',0,'2026-08-31 21:10:19',1),(93,2,1,'[{\"type\": 0, \"content\": \"哈喽\"}]',0,'2026-08-31 21:10:23',1),(94,5,1,'[{\"type\": 0, \"content\": \"❤️\"}]',0,'2026-08-31 21:10:31',1),(95,2,1,'[{\"type\": 0, \"content\": \"❤️\"}]',0,'2026-08-31 21:10:50',1),(96,2,1,'[{\"hash\": \"bb8192a48fdb60c6d1c70a6e483ac09690505321e5a58dccbb0212e61d8737ea\", \"name\": \"IMG_0188.jpeg\", \"size\": 3158503, \"type\": 1, \"content\": \"/api/file/bb8192a48fdb60c6d1c70a6e483ac09690505321e5a58dccbb0212e61d8737ea\"}]',0,'2026-08-31 21:11:04',1),(97,2,1,'[{\"type\": 0, \"content\": \"不是我发的照片\"}]',0,'2026-08-31 21:11:25',1),(98,2,1,'[{\"type\": 0, \"content\": \"{\\\"__paperphone_message\\\":1,\\\"body\\\":\\\"你有病啊\\\",\\\"reply\\\":{\\\"id\\\":\\\"4-1-1788181806513\\\",\\\"senderId\\\":\\\"4\\\",\\\"senderName\\\":\\\"4\\\",\\\"msgType\\\":\\\"text\\\",\\\"preview\\\":\\\"！！！\\\"}}\"}]',0,'2026-08-31 21:11:45',1),(99,2,4,'[{\"type\": 0, \"content\": \"你有病啊\"}]',0,'2026-08-31 21:12:07',0),(100,2,4,'[{\"type\": 0, \"content\": \"发我照片\"}]',0,'2026-08-31 21:12:11',0),(101,5,1,'[{\"type\": 0, \"content\": \"🍋\"}]',0,'2026-08-31 21:12:22',1),(102,2,4,'[{\"type\": 0, \"content\": \"💖💖\"}]',0,'2026-08-31 21:12:29',0),(103,4,1,'[{\"type\": 0, \"content\": \"🫏\"}]',0,'2026-08-31 21:12:45',1),(104,5,1,'[{\"type\": 0, \"content\": \"最后一条\"}]',0,'2026-08-31 21:22:39',1),(105,2,1,'[{\"type\": 0, \"content\": \"🌺\"}]',0,'2026-08-31 21:23:46',1),(106,5,1,'[{\"type\": 0, \"content\": \"啊啊啊啊\"}]',0,'2026-08-31 21:26:13',1),(107,2,4,'[{\"type\": 0, \"content\": \"刚刚\"}]',0,'2026-09-17 22:13:30',0),(108,4,2,'[{\"type\": 0, \"content\": \"1234123\"}]',0,'2026-09-17 22:18:21',0),(109,4,2,'[{\"type\": 0, \"content\": \"123432\"}]',0,'2026-09-17 22:18:24',0),(110,4,2,'[{\"type\": 0, \"content\": \"1234312\"}]',0,'2026-09-17 22:18:26',0),(111,4,2,'[{\"type\": 0, \"content\": \"2134\\n1234\\n1241234\"}]',0,'2026-09-17 22:35:20',0),(112,4,2,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-09-17 22:35:25',0),(113,4,2,'[{\"type\": 0, \"content\": \"12341324\"}]',0,'2026-09-17 22:35:25',0),(114,4,2,'[{\"type\": 0, \"content\": \"12432\"}]',0,'2026-09-17 22:35:25',0);
/*!40000 ALTER TABLE `messages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `moments`
--

DROP TABLE IF EXISTS `moments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `moments` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '朋友圈id',
  `owner_id` bigint NOT NULL COMMENT '朋友圈发布者id',
  `elements` json NOT NULL COMMENT '朋友圈文字内容',
  `status` tinyint DEFAULT '0' COMMENT '朋友圈状态，0-正常，1-删除',
  `like_count` int DEFAULT '0' COMMENT '点赞数',
  `visible` int DEFAULT '0' COMMENT '可见性，0-公开，1-好友可见，2-仅自己可见 3-部分好友可见',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='朋友圈表;内容主题';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `moments`
--

LOCK TABLES `moments` WRITE;
/*!40000 ALTER TABLE `moments` DISABLE KEYS */;
/*!40000 ALTER TABLE `moments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `moments_comments`
--

DROP TABLE IF EXISTS `moments_comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `moments_comments` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '评论id',
  `moment_id` bigint NOT NULL COMMENT '朋友圈id',
  `user_id` bigint NOT NULL COMMENT '评论者id',
  `status` tinyint DEFAULT '0' COMMENT '评论状态，0-正常，1-删除',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `elements` json NOT NULL COMMENT '评论元素',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='朋友圈评论表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `moments_comments`
--

LOCK TABLES `moments_comments` WRITE;
/*!40000 ALTER TABLE `moments_comments` DISABLE KEYS */;
INSERT INTO `moments_comments` VALUES (3,2,4,0,'2026-09-06 15:59:43','2026-09-06 15:59:43','[{\"url\": \"https://graceful-meatloaf.org/\", \"hash\": \"sed ex incididunt\", \"name\": \"东沐阳\", \"size\": 47, \"type\": 3, \"width\": 72, \"height\": 53, \"content\": \"exercitation magna ut\"}, {\"url\": \"https://graceful-meatloaf.org/\", \"hash\": \"sed ex incididunt\", \"name\": \"东沐阳\", \"size\": 47, \"type\": 3, \"width\": 72, \"height\": 53, \"content\": \"exercitation magna ut\"}]'),(4,6,6,0,'2026-09-14 21:53:50','2026-09-14 21:53:50','[{\"type\": 0, \"content\": \"朋友圈评论\"}]'),(7,10,2,0,'2026-09-17 22:12:40','2026-09-17 22:12:40','[{\"type\": 0, \"content\": \"哈哈哈\"}]'),(8,10,4,0,'2026-09-17 22:41:30','2026-09-17 22:41:30','[{\"type\": 0, \"content\": \"你笑傲什么\"}]');
/*!40000 ALTER TABLE `moments_comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `moments_likes`
--

DROP TABLE IF EXISTS `moments_likes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `moments_likes` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '点赞id',
  `moment_id` bigint NOT NULL COMMENT '朋友圈id',
  `user_id` bigint NOT NULL COMMENT '点赞者id',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_moment_user` (`moment_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='朋友圈点赞表;记录谁给谁的朋友圈点赞';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `moments_likes`
--

LOCK TABLES `moments_likes` WRITE;
/*!40000 ALTER TABLE `moments_likes` DISABLE KEYS */;
/*!40000 ALTER TABLE `moments_likes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `moments_privacy`
--

DROP TABLE IF EXISTS `moments_privacy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `moments_privacy` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '隐私设置id',
  `user_id` bigint NOT NULL COMMENT '用户id',
  `target_id` bigint NOT NULL COMMENT '目标用户id',
  `hide_their` tinyint DEFAULT '0' COMMENT '我不看TA的朋友圈，0-不屏蔽，1-屏蔽',
  `hide_mine` tinyint DEFAULT '0' COMMENT '不让TA看我的朋友圈，0-不屏蔽，1-屏蔽',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_target` (`user_id`,`target_id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='朋友圈隐私表;记录谁屏蔽谁的朋友圈';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `moments_privacy`
--

LOCK TABLES `moments_privacy` WRITE;
/*!40000 ALTER TABLE `moments_privacy` DISABLE KEYS */;
INSERT INTO `moments_privacy` VALUES (1,4,2,0,0),(11,2,4,0,0);
/*!40000 ALTER TABLE `moments_privacy` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `moments_visible`
--

DROP TABLE IF EXISTS `moments_visible`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `moments_visible` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '朋友圈id',
  `moment_id` bigint NOT NULL COMMENT '朋友圈id',
  `user_id` bigint NOT NULL COMMENT '用户 id',
  `visible` int DEFAULT '0' COMMENT '0 该好友可见,1 该好友不可见;需要配合coments.visible = 3 的情况',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_moment_user` (`moment_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='该条记录谁可见；谁不可见表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `moments_visible`
--

LOCK TABLES `moments_visible` WRITE;
/*!40000 ALTER TABLE `moments_visible` DISABLE KEYS */;
INSERT INTO `moments_visible` VALUES (1,5,2,0);
/*!40000 ALTER TABLE `moments_visible` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `time_comments`
--

DROP TABLE IF EXISTS `time_comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `time_comments` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '评论id',
  `time_line_id` bigint NOT NULL COMMENT '时间线id',
  `owner_id` int NOT NULL COMMENT '评论者id',
  `elements` json NOT NULL COMMENT '时间线内容',
  `status` tinyint DEFAULT '0' COMMENT '评论状态，0-正常，1-删除',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_time_line_id` (`time_line_id`),
  KEY `idx_owner_id` (`owner_id`),
  CONSTRAINT `fk_commenets_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_comments_time_line` FOREIGN KEY (`time_line_id`) REFERENCES `time_lines` (`id`) ON DELETE CASCADE,
  CONSTRAINT `chk_comments_status` CHECK ((`status` in (0,1)))
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='时间线表;记录用户的时间线信息';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `time_comments`
--

LOCK TABLES `time_comments` WRITE;
/*!40000 ALTER TABLE `time_comments` DISABLE KEYS */;
INSERT INTO `time_comments` VALUES (1,5,4,'[{\"url\": \"https://primary-farmer.name/\", \"hash\": \"magna laborum do\", \"name\": \"禄伟\", \"size\": 36, \"type\": 3, \"width\": 74, \"height\": 21, \"content\": \"sunt qui sit\"}]',1,'2026-09-13 11:42:23','2026-09-13 11:49:19');
/*!40000 ALTER TABLE `time_comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `time_likes`
--

DROP TABLE IF EXISTS `time_likes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `time_likes` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '点赞id',
  `time_line_id` bigint NOT NULL COMMENT '时间线id',
  `owner_id` int NOT NULL COMMENT '点赞者id',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_time_line_id_owner_id` (`time_line_id`,`owner_id`),
  KEY `idx_time_line_id` (`time_line_id`),
  KEY `idx_owner_id` (`owner_id`),
  CONSTRAINT `fk_likes_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_likes_time_line` FOREIGN KEY (`time_line_id`) REFERENCES `time_lines` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='时间线;点赞表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `time_likes`
--

LOCK TABLES `time_likes` WRITE;
/*!40000 ALTER TABLE `time_likes` DISABLE KEYS */;
INSERT INTO `time_likes` VALUES (1,5,4,'2026-09-13 12:06:28'),(7,14,4,'2026-09-14 22:38:43');
/*!40000 ALTER TABLE `time_likes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `time_lines`
--

DROP TABLE IF EXISTS `time_lines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `time_lines` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '时间线ID',
  `owner_id` int NOT NULL COMMENT '时间线发布者id',
  `elements` json NOT NULL COMMENT '时间线内容',
  `like_count` int DEFAULT '0' COMMENT '点赞数',
  `status` tinyint DEFAULT '0' COMMENT '时间线状态，0-正常，1-删除,2-被举报中,3-举报成功,4-举报失败',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_owner_id` (`owner_id`),
  CONSTRAINT `fk_time_line_owner_id` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `chk_like_count` CHECK ((`like_count` >= 0))
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='时间线表;记录用户的时间线信息';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `time_lines`
--

LOCK TABLES `time_lines` WRITE;
/*!40000 ALTER TABLE `time_lines` DISABLE KEYS */;
INSERT INTO `time_lines` VALUES (5,4,'[{\"url\": \"https://low-translation.biz/\", \"hash\": \"in tempor deserunt ipsum veniam\", \"name\": \"运国秀\", \"size\": 50, \"type\": 3, \"width\": 73, \"height\": 67, \"content\": \"cupidatat tempor in\"}, {\"url\": \"https://impeccable-gown.biz/\", \"hash\": \"in\", \"name\": \"宾开慧\", \"size\": 88, \"type\": 3, \"width\": 32, \"height\": 67, \"content\": \"nostrud proident in minim dolor\"}]',1,1,'2026-09-11 23:01:40','2026-09-13 12:10:36'),(6,4,'[{\"url\": \"https://unselfish-stitcher.org/\", \"hash\": \"exercitation qui cillum cupidatat in\", \"name\": \"沐宇航\", \"size\": 58, \"type\": 0, \"width\": 89, \"height\": 31, \"content\": \"in nisi deserunt\"}, {\"url\": \"https://unsightly-octave.org/\", \"hash\": \"elit aute\", \"name\": \"廖振东\", \"size\": 67, \"type\": 0, \"width\": 80, \"height\": 30, \"content\": \"commodo do\"}, {\"url\": \"https://haunting-outrun.com/\", \"hash\": \"consectetur pariatur id\", \"name\": \"权志国\", \"size\": 17, \"type\": 2, \"width\": 38, \"height\": 98, \"content\": \"velit culpa laborum esse\"}]',0,0,'2026-09-13 12:09:26','2026-09-13 12:09:26'),(10,4,'[{\"url\": \"https://defensive-litter.com/\", \"hash\": \"nulla voluptate cillum eu\", \"name\": \"拱建军\", \"size\": 25, \"type\": 0, \"width\": 70, \"height\": 17, \"content\": \"tempor occaecat\"}]',0,0,'2026-09-14 22:08:12','2026-09-14 22:08:12'),(11,4,'[{\"url\": \"https://defensive-litter.com/\", \"hash\": \"nulla voluptate cillum eu\", \"name\": \"拱建军\", \"size\": 25, \"type\": 0, \"width\": 70, \"height\": 17, \"content\": \"tempor occaecat\"}]',0,0,'2026-09-14 22:08:16','2026-09-14 22:08:16'),(12,4,'[{\"url\": \"https://defensive-litter.com/\", \"hash\": \"nulla voluptate cillum eu\", \"name\": \"拱建军\", \"size\": 25, \"type\": 0, \"width\": 70, \"height\": 17, \"content\": \"tempor occaecat\"}]',0,0,'2026-09-14 22:21:42','2026-09-14 22:21:42'),(13,4,'[{\"url\": \"https://defensive-litter.com/\", \"hash\": \"nulla voluptate cillum eu\", \"name\": \"拱建军\", \"size\": 25, \"type\": 0, \"width\": 70, \"height\": 17, \"content\": \"tempor occaecat\"}]',0,0,'2026-09-14 22:22:04','2026-09-14 22:22:04'),(14,4,'[{\"url\": \"https://defensive-litter.com/\", \"hash\": \"nulla voluptate cillum eu\", \"name\": \"拱建军\", \"size\": 25, \"type\": 0, \"width\": 70, \"height\": 17, \"content\": \"tempor occaecat\"}]',1,0,'2026-09-14 22:22:08','2026-09-14 22:38:43');
/*!40000 ALTER TABLE `time_lines` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `avatar` varchar(255) DEFAULT NULL COMMENT '用户头像hash',
  `nickname` varchar(255) DEFAULT NULL COMMENT '用户中文名称',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (2,'whx_1','$2a$10$YBxMVnuiMKAIWGz2QeWeSOAjaMpppcrJh1p5DouxUIH1Kya8G0dyS','2026-08-18 22:17:09','2026-09-17 22:40:54','93d81e39faf33e74233449cfc801eda0cbf57476f86124b2ad6a6b2e8058ac5c','袁境莲'),(4,'whx','$2a$10$YBxMVnuiMKAIWGz2QeWeSOAjaMpppcrJh1p5DouxUIH1Kya8G0dyS','2026-08-22 11:07:00','2026-09-14 22:25:49','adc562d7850362d7cc86da1fd21146d10d5e6c35114a6f1dd80ad8de5921065e','苦立伟'),(5,'whx_2','$2a$10$jJTlGRxNQpyns2G/FNFrDegz/u7UL9JPnRO9Iz1mPMhoJ7t9kJii.','2026-08-30 22:19:10','2026-08-30 22:22:41','ba7c041d8df09bd31af59cb52b8ba0022e29f0465be7159a8ae86193ad618fc6','王杰瑞');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-09-22 21:49:03
