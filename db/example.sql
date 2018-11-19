-- MySQL dump 10.13  Distrib 5.7.22, for Linux (x86_64)
--
-- Host: localhost    Database: gdrift
-- ------------------------------------------------------
-- Server version	5.7.22-0ubuntu0.16.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `projects`
--

DROP TABLE IF EXISTS `projects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `projects` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(40) NOT NULL,
  `users_id` int(11) NOT NULL,
  `n_populations` int(11) NOT NULL,
  `date_created` datetime DEFAULT NULL,
  `date_finished` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `projects_users_fk1` (`users_id`),
  CONSTRAINT `projects_users_fk1` FOREIGN KEY (`users_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `projects`
--

LOCK TABLES `projects` WRITE;
/*!40000 ALTER TABLE `projects` DISABLE KEYS */;
INSERT INTO `projects` VALUES (20,'Test Project 01',1,0,'2018-05-17 13:23:06',NULL),(21,'Test Project 01',1,0,'2018-05-17 13:23:23',NULL),(22,'Test Project 01',1,0,'2018-05-17 13:26:12',NULL),(23,'Test Project 01',1,0,'2018-05-17 13:39:50',NULL),(24,'Test Project 01',1,0,'2018-05-17 13:47:44',NULL),(25,'Test Project 01',1,0,'2018-05-17 13:49:04',NULL),(26,'Test Project 01',1,0,'2018-05-17 13:52:44',NULL),(27,'Test Project 01',1,2,'2018-05-17 14:35:22',NULL),(31,'Test Project 01',1,2,'2018-05-17 16:37:43',NULL),(32,'Test Project 01',1,2,'2018-05-17 16:42:56',NULL),(33,'Test Project 01',1,2,'2018-05-17 16:45:22',NULL);
/*!40000 ALTER TABLE `projects` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `simulations`
--

DROP TABLE IF EXISTS `simulations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `simulations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `projects_id` int(11) NOT NULL,
  `batch_size` int(11) NOT NULL,
  `model` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `simulations_projects_fk1` (`projects_id`),
  CONSTRAINT `simulations_projects_fk1` FOREIGN KEY (`projects_id`) REFERENCES `projects` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `simulations`
--

LOCK TABLES `simulations` WRITE;
/*!40000 ALTER TABLE `simulations` DISABLE KEYS */;
INSERT INTO `simulations` VALUES (1,27,1),(2,27,1),(3,27,1),(4,27,1),(6,27,1);
/*!40000 ALTER TABLE `simulations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(40) NOT NULL,
  `password` varchar(40) NOT NULL,
  `name` varchar(40) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'user','pass','Test User');
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

-- Dump completed on 2018-05-22 17:23:52
