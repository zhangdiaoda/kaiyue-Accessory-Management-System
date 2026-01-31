-- MySQL dump 10.13  Distrib 8.0.45, for Linux (x86_64)
--
-- Host: localhost    Database: warehouse
-- ------------------------------------------------------
-- Server version	8.0.45

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
-- Table structure for table `announcement`
--

DROP TABLE IF EXISTS `announcement`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `announcement` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(200) NOT NULL,
  `content` text NOT NULL,
  `type` varchar(20) NOT NULL,
  `status` bigint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `announcement`
--

LOCK TABLES `announcement` WRITE;
/*!40000 ALTER TABLE `announcement` DISABLE KEYS */;
/*!40000 ALTER TABLE `announcement` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `borrow_record`
--

DROP TABLE IF EXISTS `borrow_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `borrow_record` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `record_no` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '记录编号',
  `employee_id` bigint NOT NULL COMMENT '员工ID',
  `part_id` bigint NOT NULL COMMENT '配件ID',
  `borrow_quantity` int NOT NULL COMMENT '领用数量',
  `return_quantity` int DEFAULT '0' COMMENT '已归还数量',
  `damaged_quantity` int DEFAULT '0' COMMENT '损毁数量',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '状态：BORROWED/RETURNED/PARTIAL_RETURNED',
  `borrow_time` datetime NOT NULL COMMENT '领用时间',
  `borrow_admin_id` bigint NOT NULL COMMENT '登记管理员ID',
  `return_time` datetime DEFAULT NULL COMMENT '归还时间',
  `return_admin_id` bigint DEFAULT NULL COMMENT '归还登记管理员ID',
  `remark` text COLLATE utf8mb4_unicode_ci COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `record_no` (`record_no`),
  KEY `idx_employee_id` (`employee_id`),
  KEY `idx_part_id` (`part_id`),
  KEY `idx_status` (`status`),
  KEY `idx_borrow_time` (`borrow_time`),
  KEY `borrow_admin_id` (`borrow_admin_id`),
  KEY `return_admin_id` (`return_admin_id`),
  CONSTRAINT `borrow_record_ibfk_1` FOREIGN KEY (`employee_id`) REFERENCES `employee` (`id`),
  CONSTRAINT `borrow_record_ibfk_2` FOREIGN KEY (`part_id`) REFERENCES `part` (`id`),
  CONSTRAINT `borrow_record_ibfk_3` FOREIGN KEY (`borrow_admin_id`) REFERENCES `sys_user` (`id`),
  CONSTRAINT `borrow_record_ibfk_4` FOREIGN KEY (`return_admin_id`) REFERENCES `sys_user` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='领用记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `borrow_record`
--

LOCK TABLES `borrow_record` WRITE;
/*!40000 ALTER TABLE `borrow_record` DISABLE KEYS */;
INSERT INTO `borrow_record` VALUES (1,'BR20260125304873',1,2,3,0,0,'BORROWED','2026-01-25 09:34:34',1,NULL,NULL,'','2026-01-25 09:34:34','2026-01-25 09:34:34'),(2,'BR20260125304887',1,2,4,0,0,'BORROWED','2026-01-25 09:34:48',1,NULL,NULL,'','2026-01-25 09:34:48','2026-01-25 09:34:48'),(3,'BR20260125306224',1,1,3,0,3,'RETURNED','2026-01-25 09:57:05',1,'2026-01-25 11:44:59',1,'','2026-01-25 09:57:05','2026-01-25 11:44:59'),(4,'BR20260125306783',2,2,1,0,0,'BORROWED','2026-01-25 10:06:23',1,NULL,NULL,'222','2026-01-25 10:06:23','2026-01-25 10:06:23'),(5,'BR20260125306853',2,1,1,0,0,'BORROWED','2026-01-25 10:07:33',1,NULL,NULL,'','2026-01-25 10:07:33','2026-01-25 10:07:33'),(6,'BR20260125312699',1,1,1,0,0,'BORROWED','2026-01-25 11:45:00',1,NULL,NULL,'','2026-01-25 11:45:00','2026-01-25 11:45:00');
/*!40000 ALTER TABLE `borrow_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `employee`
--

DROP TABLE IF EXISTS `employee`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `employee` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `employee_no` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `department` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `position` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` bigint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `employee_no` (`employee_no`),
  UNIQUE KEY `employee_no_2` (`employee_no`),
  UNIQUE KEY `idx_employee_employee_no` (`employee_no`),
  KEY `idx_employee_no` (`employee_no`),
  KEY `idx_department` (`department`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='员工表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `employee`
--

LOCK TABLES `employee` WRITE;
/*!40000 ALTER TABLE `employee` DISABLE KEYS */;
INSERT INTO `employee` VALUES (1,'EMP001','张三','生产部','机械工程师','13900139001',1,'2026-01-25 09:00:42.000','2026-01-25 09:00:42.000'),(2,'EMP002','李四','生产部','技术员','13900139002',1,'2026-01-25 09:00:42.000','2026-01-25 09:00:42.000'),(3,'EMP003','王五','质检部','质检员','13900139003',1,'2026-01-25 09:00:42.000','2026-01-25 09:00:42.000');
/*!40000 ALTER TABLE `employee` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `internal_message`
--

DROP TABLE IF EXISTS `internal_message`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `internal_message` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `sender_id` bigint unsigned DEFAULT NULL,
  `receiver_id` bigint unsigned DEFAULT NULL,
  `title` varchar(200) NOT NULL,
  `content` text NOT NULL,
  `is_read` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `internal_message`
--

LOCK TABLES `internal_message` WRITE;
/*!40000 ALTER TABLE `internal_message` DISABLE KEYS */;
INSERT INTO `internal_message` VALUES (1,1,2,'111','111',0,'2026-01-25 11:21:28.257');
/*!40000 ALTER TABLE `internal_message` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `old_part_inventory`
--

DROP TABLE IF EXISTS `old_part_inventory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `old_part_inventory` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `part_id` bigint unsigned NOT NULL,
  `employee_id` bigint unsigned NOT NULL,
  `quantity` bigint NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `old_part_inventory`
--

LOCK TABLES `old_part_inventory` WRITE;
/*!40000 ALTER TABLE `old_part_inventory` DISABLE KEYS */;
/*!40000 ALTER TABLE `old_part_inventory` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `part`
--

DROP TABLE IF EXISTS `part`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `part` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `part_no` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `category_id` bigint NOT NULL COMMENT '分类ID',
  `specification` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '规格型号',
  `unit` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '件' COMMENT '单位',
  `stock_quantity` int DEFAULT '0' COMMENT '当前库存数量',
  `warning_threshold` int DEFAULT '10' COMMENT '预警阈值',
  `price` decimal(10,2) DEFAULT NULL COMMENT '单价',
  `remark` text COLLATE utf8mb4_unicode_ci COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `part_no` (`part_no`),
  UNIQUE KEY `part_no_2` (`part_no`),
  KEY `idx_part_no` (`part_no`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_stock_warning` (`stock_quantity`,`warning_threshold`),
  CONSTRAINT `part_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `part_category` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='配件信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `part`
--

LOCK TABLES `part` WRITE;
/*!40000 ALTER TABLE `part` DISABLE KEYS */;
INSERT INTO `part` VALUES (1,'PART001','硬质合金刀片',2,'CNMG120408','片',95,20,15.50,NULL,'2026-01-25 09:00:42','2026-01-25 11:45:00'),(2,'PART002','高速钢钻头',3,'Φ10mm','支',42,10,8.00,NULL,'2026-01-25 09:00:42','2026-01-25 10:06:23'),(3,'PART003','立铣刀',4,'Φ12mm','支',30,5,25.00,NULL,'2026-01-25 09:00:42','2026-01-25 09:00:42'),(4,'PART004','游标卡尺',5,'0-150mm','把',15,3,45.00,NULL,'2026-01-25 09:00:42','2026-01-25 09:00:42'),(5,'PART005','千分尺',6,'0-25mm','个',8,2,120.00,NULL,'2026-01-25 09:00:42','2026-01-25 09:00:42');
/*!40000 ALTER TABLE `part` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `part_category`
--

DROP TABLE IF EXISTS `part_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `part_category` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称',
  `parent_id` bigint DEFAULT '0' COMMENT '父分类ID（0为顶级）',
  `sort_order` int DEFAULT '0' COMMENT '排序',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='配件分类表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `part_category`
--

LOCK TABLES `part_category` WRITE;
/*!40000 ALTER TABLE `part_category` DISABLE KEYS */;
INSERT INTO `part_category` VALUES (1,'刀具类',0,1,'2026-01-25 09:00:41','2026-01-25 09:00:41'),(2,'刀片',1,1,'2026-01-25 09:00:41','2026-01-25 09:00:41'),(3,'钻头',1,2,'2026-01-25 09:00:41','2026-01-25 09:00:41'),(4,'铣刀',1,3,'2026-01-25 09:00:41','2026-01-25 09:00:41'),(5,'测量工具',0,2,'2026-01-25 09:00:41','2026-01-25 09:00:41'),(6,'卡尺',4,1,'2026-01-25 09:00:41','2026-01-25 09:00:41'),(7,'千分尺',4,2,'2026-01-25 09:00:41','2026-01-25 09:00:41'),(8,'辅助工具',0,3,'2026-01-25 09:00:41','2026-01-25 09:00:41');
/*!40000 ALTER TABLE `part_category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scrap_inventory`
--

DROP TABLE IF EXISTS `scrap_inventory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `scrap_inventory` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `part_id` bigint unsigned NOT NULL,
  `employee_id` bigint unsigned NOT NULL,
  `quantity` bigint NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scrap_inventory`
--

LOCK TABLES `scrap_inventory` WRITE;
/*!40000 ALTER TABLE `scrap_inventory` DISABLE KEYS */;
INSERT INTO `scrap_inventory` VALUES (1,1,1,3,'2026-01-25 11:44:59.399','2026-01-25 11:44:59.399');
/*!40000 ALTER TABLE `scrap_inventory` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_config`
--

DROP TABLE IF EXISTS `sys_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_config` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `config_key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `config_value` text COLLATE utf8mb4_unicode_ci,
  `description` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `config_key` (`config_key`),
  UNIQUE KEY `config_key_2` (`config_key`),
  UNIQUE KEY `idx_sys_config_config_key` (`config_key`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_config`
--

LOCK TABLES `sys_config` WRITE;
/*!40000 ALTER TABLE `sys_config` DISABLE KEYS */;
INSERT INTO `sys_config` VALUES (1,'dingtalk.webhook.url','','钉钉群机器人Webhook地址','2026-01-25 09:00:41.000','2026-01-25 11:09:38.725'),(2,'dingtalk.secret','','钉钉群机器人加签密钥','2026-01-25 09:00:41.000','2026-01-25 11:09:38.629'),(3,'report.auto_push','true','是否自动推送报表','2026-01-25 09:00:41.000','2026-01-25 11:09:39.341'),(4,'report.push_day','1','每月推送日期（1-28）','2026-01-25 09:00:41.000','2026-01-25 11:09:38.349'),(5,'stock.warning.enabled','true','是否启用库存预警','2026-01-25 09:00:41.000','2026-01-25 11:09:38.825'),(6,'copyright','© 2026 凯越机械 版权所有','','2026-01-25 10:56:03.000','2026-01-25 11:09:39.058'),(7,'login_subtitle','蒙阴县凯越工程机械有限公司','','2026-01-25 10:56:03.000','2026-01-25 11:09:39.153'),(8,'brand_name','蒙阴县凯越工程机械有限公司','','2026-01-25 10:56:03.000','2026-01-25 11:09:39.245'),(9,'brand_logo','','','2026-01-25 10:56:04.000','2026-01-25 11:09:38.537'),(10,'company_name','蒙阴县凯越工程机械有限公司','','2026-01-25 11:06:11.056','2026-01-25 11:09:38.962'),(11,'system_name','配件仓储管理系统','','2026-01-25 11:06:11.519','2026-01-25 11:09:38.445'),(12,'backup_path','桌面','数据库备份路径','2026-01-25 13:42:38.349','2026-01-25 13:42:38.349'),(13,'backup_schedule','','数据库备份定时规则','2026-01-25 13:42:38.523','2026-01-25 13:42:38.523');
/*!40000 ALTER TABLE `sys_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user`
--

DROP TABLE IF EXISTS `sys_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `real_name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `role` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `department` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` bigint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `username_2` (`username`),
  UNIQUE KEY `idx_sys_user_username` (`username`),
  KEY `idx_username` (`username`),
  KEY `idx_role` (`role`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员用户表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user`
--

LOCK TABLES `sys_user` WRITE;
/*!40000 ALTER TABLE `sys_user` DISABLE KEYS */;
INSERT INTO `sys_user` VALUES (1,'admin','$2a$10$vLMtSqe5Hvml024.bBBX9Oqs90aebDlYy45v0NdkIkZZS.RN7R06q','系统管理员','SUPER_ADMIN','技术部','13800138000',1,'2026-01-25 09:00:41.000','2026-01-25 09:15:31.000'),(2,'warehouse','$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH','仓库管理员','WAREHOUSE_ADMIN','仓储部','13800138001',1,'2026-01-25 09:00:41.000','2026-01-25 09:00:41.000');
/*!40000 ALTER TABLE `sys_user` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-01-25 13:42:39
