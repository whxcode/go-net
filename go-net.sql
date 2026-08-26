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
-- Table structure for table `friends`
--

DROP TABLE IF EXISTS `friends`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `friends` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `friend_id` bigint NOT NULL COMMENT '好友ID',
  `status` tinyint(1) DEFAULT '0' COMMENT '0-待确认, 1-已确认, 2-已拒绝, 3-已删除',
  `remark` varchar(50) DEFAULT '' COMMENT '备注名',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_friend` (`user_id`,`friend_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_friend_id` (`friend_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `friends`
--

LOCK TABLES `friends` WRITE;
/*!40000 ALTER TABLE `friends` DISABLE KEYS */;
INSERT INTO `friends` VALUES (8,4,2,1,'','2026-08-23 21:02:40','2026-08-26 21:53:53');
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
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='群组表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `group_chats`
--

LOCK TABLES `group_chats` WRITE;
/*!40000 ALTER TABLE `group_chats` DISABLE KEYS */;
INSERT INTO `group_chats` VALUES (1,4,'红警战队','xx','今天打红警','2026-08-26 22:15:09','2026-08-26 22:40:22',0,0);
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
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_user` (`group_id`,`user_id`),
  KEY `idx_group_id` (`group_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='群成员表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `group_members`
--

LOCK TABLES `group_members` WRITE;
/*!40000 ALTER TABLE `group_members` DISABLE KEYS */;
INSERT INTO `group_members` VALUES (1,1,4,2,'2026-08-26 22:16:33','2026-08-26 22:16:33','2026-08-26 22:16:33',0,0),(3,1,2,0,'2026-08-26 22:18:30','2026-08-26 22:18:30','2026-08-26 22:18:30',0,0);
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
  `msg_id` varchar(64) NOT NULL COMMENT '消息唯一标识（雪花ID）',
  `sender_id` int NOT NULL COMMENT '发送者ID',
  `receiver_id` int NOT NULL COMMENT '接收者ID',
  `elements` json NOT NULL COMMENT '消息内容元素列表',
  `status` tinyint(1) DEFAULT '0' COMMENT '0-未读, 1-已读',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_msg_id` (`msg_id`),
  KEY `idx_sender_id` (`sender_id`),
  KEY `idx_receiver_id` (`receiver_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=59 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='消息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `messages`
--

LOCK TABLES `messages` WRITE;
/*!40000 ALTER TABLE `messages` DISABLE KEYS */;
INSERT INTO `messages` VALUES (1,'',1,2,'[{\"type\": 1, \"content\": \"1234\"}]',0,'2026-08-14 22:06:50'),(3,'2088268225983811584',1,2,'[{\"type\": 1, \"content\": \"1234\"}]',0,'2026-08-14 22:15:18'),(4,'2088268268291756032',1,2,'[{\"type\": 1, \"content\": \"1234\"}]',0,'2026-08-14 22:15:28'),(5,'2088268276713918464',1,2,'[{\"type\": 1, \"content\": \"123424\"}]',0,'2026-08-14 22:15:30'),(6,'2088268285085749248',1,2,'[{\"type\": 1, \"content\": \"1234214\"}]',0,'2026-08-14 22:15:32'),(7,'2088268350474948608',1,2,'[{\"type\": 1, \"content\": \"哈哈哈哈哈哈哈；你在干嘛\"}]',0,'2026-08-14 22:15:48'),(8,'2088272728799842304',1,2,'[{\"type\": 1, \"content\": \"1\"}]',0,'2026-08-14 22:33:11'),(9,'2088272732130119680',1,2,'[{\"type\": 1, \"content\": \"2\"}]',0,'2026-08-14 22:33:12'),(10,'2088272735670112256',1,2,'[{\"type\": 1, \"content\": \"3\"}]',0,'2026-08-14 22:33:13'),(11,'2088272753621733376',2,1,'[{\"type\": 1, \"content\": \"1234\"}]',0,'2026-08-14 22:33:17'),(12,'2088272766313697280',2,1,'[{\"type\": 1, \"content\": \"1234\"}]',0,'2026-08-14 22:33:20'),(13,'2088272770092765184',2,1,'[{\"type\": 1, \"content\": \"1234\"}]',0,'2026-08-14 22:33:21'),(14,'2088272772357689344',2,1,'[{\"type\": 1, \"content\": \"1234\"}]',0,'2026-08-14 22:33:22'),(15,'2088272773783752704',2,1,'[{\"type\": 1, \"content\": \"13241\"}]',0,'2026-08-14 22:33:22'),(16,'2088272774526144512',2,1,'[{\"type\": 1, \"content\": \"234\"}]',0,'2026-08-14 22:33:22'),(17,'2088272775092375552',2,1,'[{\"type\": 1, \"content\": \"123\"}]',0,'2026-08-14 22:33:22'),(18,'2088272775843155968',2,1,'[{\"type\": 1, \"content\": \"4123\"}]',0,'2026-08-14 22:33:23'),(19,'2088272776367443968',2,1,'[{\"type\": 1, \"content\": \"4\"}]',0,'2026-08-14 22:33:23'),(20,'2088272776933675008',2,1,'[{\"type\": 1, \"content\": \"123\"}]',0,'2026-08-14 22:33:23'),(21,'2088272804280537088',1,2,'[{\"type\": 1, \"content\": \"还会\"}]',0,'2026-08-14 22:33:29'),(22,'2088272809695383552',2,1,'[{\"type\": 1, \"content\": \"4231123423\"}]',0,'2026-08-14 22:33:31'),(23,'2090435056978890752',1,2,'[{\"type\": 1, \"content\": \"2134214\"}]',0,'2026-08-20 21:45:31'),(24,'2090435062171439104',2,1,'[{\"type\": 1, \"content\": \"1234324\"}]',0,'2026-08-20 21:45:32'),(25,'2090439558691819520',1,2,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-20 22:03:24'),(26,'2090439654548443136',2,1,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-20 22:03:47'),(27,'2090440646648139776',1,2,'[{\"type\": 0, \"content\": \"123424\"}]',0,'2026-08-20 22:07:43'),(28,'2090440683054698496',2,1,'[{\"type\": 0, \"content\": \"1234124\"}]',0,'2026-08-20 22:07:52'),(29,'2090440697155948544',1,2,'[{\"type\": 0, \"content\": \"2314324\"}]',0,'2026-08-20 22:07:55'),(30,'2090440727954722816',2,1,'[{\"type\": 0, \"content\": \"？？\"}]',0,'2026-08-20 22:08:03'),(31,'2090440749031100416',2,1,'[{\"type\": 0, \"content\": \"😄\"}]',0,'2026-08-20 22:08:08'),(32,'2090440769058902016',2,1,'[{\"type\": 0, \"content\": \"哈哈哈\"}]',0,'2026-08-20 22:08:13'),(33,'2090440864554815488',2,1,'[{\"type\": 0, \"content\": \"💚\"}]',0,'2026-08-20 22:08:35'),(34,'2090440924998930432',1,2,'[{\"type\": 0, \"content\": \"😀\"}]',0,'2026-08-20 22:08:50'),(35,'2090440932783558656',1,2,'[{\"type\": 0, \"content\": \"😄\"}]',0,'2026-08-20 22:08:52'),(36,'2090440944271757312',1,2,'[{\"type\": 0, \"content\": \"🤩\"}]',0,'2026-08-20 22:08:54'),(37,'2090440970234499072',1,2,'[{\"type\": 0, \"content\": \"☺️\"}]',0,'2026-08-20 22:09:00'),(38,'2090443694476890112',1,2,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-20 22:19:50'),(39,'2090443711417683968',2,1,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-20 22:19:54'),(40,'2090444214004355072',1,2,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-20 22:21:54'),(41,'2090444231628820480',2,1,'[{\"type\": 0, \"content\": \"1324\"}]',0,'2026-08-20 22:21:58'),(42,'2090444255280500736',2,1,'[{\"type\": 0, \"content\": \"🍅\"}]',0,'2026-08-20 22:22:04'),(43,'2090444294702764032',2,1,'[{\"type\": 0, \"content\": \"🍥🥨🍚🌯\"}]',0,'2026-08-20 22:22:13'),(44,'2090444319562403840',1,2,'[{\"type\": 0, \"content\": \"🫏🫏\"}]',0,'2026-08-20 22:22:19'),(45,'2090444441016864768',2,1,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-20 22:22:48'),(46,'2090444455025840128',2,1,'[{\"type\": 0, \"content\": \"1234\"}]',0,'2026-08-20 22:22:51'),(47,'2090444466555981824',1,2,'[{\"type\": 0, \"content\": \"🦄\"}]',0,'2026-08-20 22:22:54'),(48,'2091508450738573312',4,2,'[{\"type\": 0, \"content\": \"13241234\"}]',0,'2026-08-23 20:50:48'),(49,'2091518677626130432',2,2,'[{\"type\": 0, \"content\": \"12343214\"}]',0,'2026-08-23 21:31:26'),(50,'2091518702586433536',4,2,'[{\"type\": 0, \"content\": \"21343124\"}]',0,'2026-08-23 21:31:32'),(51,'2091518733129355264',2,4,'[{\"type\": 0, \"content\": \"1234234\"}]',0,'2026-08-23 21:31:39'),(52,'2091518774501969920',2,4,'[{\"type\": 0, \"content\": \"😁\"}]',0,'2026-08-23 21:31:49'),(53,'2091519043428159488',2,4,'[{\"type\": 0, \"content\": \"肺癌无法哇\"}]',0,'2026-08-23 21:32:53'),(54,'2091519069919383552',4,2,'[{\"type\": 0, \"content\": \"😀🤩\"}]',0,'2026-08-23 21:32:59'),(55,'2091519112055361536',2,4,'[{\"type\": 0, \"content\": \"/api/file/75bc2ecde322ddddd4eb2d5b0f79ceed3c76b3e5eca8dae54a8cfe2ef6cbaa3b\"}]',0,'2026-08-23 21:33:10'),(56,'2091519293580644352',2,4,'[{\"type\": 0, \"content\": \"1\"}]',0,'2026-08-23 21:33:53'),(57,'2091519319497248768',2,4,'[{\"type\": 0, \"content\": \"/api/file/75bc2ecde322ddddd4eb2d5b0f79ceed3c76b3e5eca8dae54a8cfe2ef6cbaa3b\"}]',0,'2026-08-23 21:33:59'),(58,'2091532123260325888',2,4,'[{\"hash\": \"75bc2ecde322ddddd4eb2d5b0f79ceed3c76b3e5eca8dae54a8cfe2ef6cbaa3b\", \"name\": \"Snipaste_2026-08-22_13-28-24.png\", \"size\": 71560, \"type\": 1, \"content\": \"/api/file/75bc2ecde322ddddd4eb2d5b0f79ceed3c76b3e5eca8dae54a8cfe2ef6cbaa3b\"}]',0,'2026-08-23 22:24:52');
/*!40000 ALTER TABLE `messages` ENABLE KEYS */;
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
  `avatar` varchar(255) DEFAULT NULL COMMENT '用户头像；创建时为空',
  `nickname` varchar(255) DEFAULT NULL COMMENT '用户中文名称；可创建之后修改',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (2,'whx_1','$2a$10$YBxMVnuiMKAIWGz2QeWeSOAjaMpppcrJh1p5DouxUIH1Kya8G0dyS','2026-08-18 22:17:09','2026-08-23 21:05:35','xx','袁境莲'),(4,'whx','$2a$10$YBxMVnuiMKAIWGz2QeWeSOAjaMpppcrJh1p5DouxUIH1Kya8G0dyS','2026-08-22 11:07:00','2026-08-23 21:05:46','75bc2ecde322ddddd4eb2d5b0f79ceed3c76b3e5eca8dae54a8cfe2ef6cbaa3b','王恒兴');
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

-- Dump completed on 2026-08-26 22:54:43
