-- MySQL dump 10.13  Distrib 8.0.33, for Win64 (x86_64)
--
-- Host: localhost    Database: mydatabase
-- ------------------------------------------------------
-- Server version	8.0.33

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
-- Table structure for table `capitais_temperatura`
--

DROP TABLE IF EXISTS `capitais_temperatura`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `capitais_temperatura` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `UF` varchar(2) DEFAULT NULL,
  `Cidade` varchar(50) DEFAULT NULL,
  `Temperatura` float DEFAULT NULL,
  `Data` date DEFAULT NULL,
  `Hora_Local` time DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `capitais_temperatura`
--

LOCK TABLES `capitais_temperatura` WRITE;
/*!40000 ALTER TABLE `capitais_temperatura` DISABLE KEYS */;
INSERT INTO `capitais_temperatura` VALUES (1,'AC','Rio Branco',28.34,'2023-05-14','18:58:54'),(2,'AL','Maceió',25.69,'2023-05-14','18:58:54'),(3,'AP','Macapá',28.99,'2023-05-14','18:58:54'),(4,'AM','Manaus',31.27,'2023-05-14','18:58:54'),(5,'BA','Salvador',31.14,'2023-05-14','18:58:54'),(6,'CE','Fortaleza',28.07,'2023-05-14','18:58:54'),(7,'DF','Brasília',24.51,'2023-05-14','18:58:54'),(8,'ES','Vitória',23.25,'2023-05-14','18:58:54'),(9,'GO','Goiânia',24.83,'2023-05-14','18:58:54'),(10,'MA','São Luís',29.11,'2023-05-14','18:58:54'),(11,'MT','Cuiabá',25.96,'2023-05-14','18:58:54'),(12,'MS','Campo Grande',22.75,'2023-05-14','18:58:54'),(13,'MG','Belo Horizonte',21.5,'2023-05-14','18:58:54'),(14,'PA','Belém',27.02,'2023-05-14','18:58:54'),(15,'PB','João Pessoa',28.13,'2023-05-14','18:58:54'),(16,'PR','Curitiba',16.02,'2023-05-14','18:58:54'),(17,'PE','Recife',29.02,'2023-05-14','18:58:54'),(18,'PI','Teresina',29.84,'2023-05-14','18:58:54'),(19,'RJ','Rio de Janeiro',22.76,'2023-05-14','18:58:54'),(20,'RN','Natal',28.12,'2023-05-14','18:58:54'),(21,'RS','Porto Alegre',20.31,'2023-05-14','18:58:54'),(22,'RO','Porto Velho',32.03,'2023-05-14','18:58:54'),(23,'RR','Boa Vista',29.99,'2023-05-14','18:58:54'),(24,'SC','Florianópolis',18.57,'2023-05-14','18:58:54'),(25,'SP','São Paulo',17.22,'2023-05-14','18:58:54'),(26,'SE','Aracaju',27.97,'2023-05-14','18:58:54'),(27,'TO','Palmas',30.93,'2023-05-14','18:58:54');
/*!40000 ALTER TABLE `capitais_temperatura` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `follow`
--

DROP TABLE IF EXISTS `follow`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `follow` (
  `followID` int NOT NULL AUTO_INCREMENT,
  `followBy` int NOT NULL,
  `followTo` int NOT NULL,
  `followTime` varchar(255) NOT NULL,
  PRIMARY KEY (`followID`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `follow`
--

LOCK TABLES `follow` WRITE;
/*!40000 ALTER TABLE `follow` DISABLE KEYS */;
INSERT INTO `follow` VALUES (7,6,5,'2017-09-25 10:15:53.957231'),(9,5,6,'2017-09-25 11:13:03.499015'),(10,7,6,'2017-09-25 15:34:00.040075'),(11,8,6,'2017-09-25 15:46:39.580505'),(12,6,7,'2017-09-25 19:42:57.891789');
/*!40000 ALTER TABLE `follow` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `likes`
--

DROP TABLE IF EXISTS `likes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `likes` (
  `likeID` int NOT NULL AUTO_INCREMENT,
  `postID` int NOT NULL,
  `likeBy` int NOT NULL,
  `likeTime` varchar(255) NOT NULL,
  PRIMARY KEY (`likeID`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `likes`
--

LOCK TABLES `likes` WRITE;
/*!40000 ALTER TABLE `likes` DISABLE KEYS */;
INSERT INTO `likes` VALUES (6,9,5,'2017-09-25 16:35:50.840949'),(7,9,6,'2017-09-25 19:10:00.3598');
/*!40000 ALTER TABLE `likes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `posts`
--

DROP TABLE IF EXISTS `posts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `posts` (
  `postID` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `content` text NOT NULL,
  `createdBy` int NOT NULL,
  `createdAt` varchar(255) NOT NULL,
  PRIMARY KEY (`postID`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `posts`
--

LOCK TABLES `posts` WRITE;
/*!40000 ALTER TABLE `posts` DISABLE KEYS */;
INSERT INTO `posts` VALUES (2,'second','second_content',6,'2017-09-23 10:43:5.941602'),(5,'third','third content..',6,'2017-09-23 11:32:45.941602'),(9,'Hello,','World!!',6,'2017-09-24 19:09:37.024131'),(10,'my title..','my content...',5,'2017-09-25 14:20:35.959114'),(11,'ghalib\'s first title..','and this is content!!!',7,'2017-09-25 15:38:51.705595'),(12,'jkj','kj',8,'2017-09-25 15:43:24.782827');
/*!40000 ALTER TABLE `posts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `profile_views`
--

DROP TABLE IF EXISTS `profile_views`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `profile_views` (
  `viewID` int NOT NULL AUTO_INCREMENT,
  `viewBy` int NOT NULL,
  `viewTo` int NOT NULL,
  `viewTime` varchar(255) NOT NULL,
  PRIMARY KEY (`viewID`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `profile_views`
--

LOCK TABLES `profile_views` WRITE;
/*!40000 ALTER TABLE `profile_views` DISABLE KEYS */;
INSERT INTO `profile_views` VALUES (14,6,5,'2017-09-25 10:49:21.373119'),(15,5,6,'2017-09-25 11:13:01.246732'),(16,5,6,'2017-09-25 11:13:03.912177'),(17,5,6,'2017-09-25 11:13:11.465508'),(18,6,5,'2017-09-25 11:30:52.987842'),(19,7,6,'2017-09-25 15:33:58.973987'),(20,7,6,'2017-09-25 15:34:00.460281'),(21,8,6,'2017-09-25 15:46:37.855424'),(22,8,6,'2017-09-25 15:46:40.10399'),(23,6,7,'2017-09-25 19:42:12.359235'),(24,6,7,'2017-09-25 19:42:56.505978'),(25,6,7,'2017-09-25 19:42:58.333437');
/*!40000 ALTER TABLE `profile_views` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(32) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `bio` text NOT NULL,
  `joined` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (5,'takkar','takkar@gmail.com','$2a$10$ttnsVDOPgMlA5vvDE33eneqVO3BHE/zif/axxI5AwNpOuRetkxFk6','','2017-09-23 07:02:14.941602'),(6,'faiyaz','faiyaz@gmail.com','$2a$10$.Wx2jBjYPiMFgWGCW.USze.qFMwrgN1TWOf50CQgqHDBzpcYV2uSa','','2017-09-23 01:22:41.941602'),(7,'ghalib','ghalib@gmail.com','$2a$10$ziw6cqTgpSBIvASZOjTheey8sQYf1iW3HW4N.Xjq4GX6faKqzIrE.','','2017-09-23 01:22:41.941602'),(8,'nature','nature@gmail.com','$2a$10$nBi64BlbJMlzuSJfOhPlXevwdCgHOXKLZQUbJQ1q2Y7Ltbpaf1Woa','','2017-09-25 15:43:08.675349');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users1`
--

DROP TABLE IF EXISTS `users1`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users1` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `username` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users1`
--

LOCK TABLES `users1` WRITE;
/*!40000 ALTER TABLE `users1` DISABLE KEYS */;
INSERT INTO `users1` VALUES (1,NULL,'mat666','rock666@email.com','$2a$10$n6Nl4QNvhBbmpAkVwEXWLeFkxuyUnVBNqahGSTaTbmQIgpf4ZwnRK'),(2,NULL,'mateus','mateuscavalcant7@gmail.com','$2a$10$fxrGTT9vp0kqlwyh8vzX8ODNh1M/9ua4oWhF1reb6dyxNE3ckZ5Wu'),(4,NULL,'mat','mateusc@gmail.com','$2a$10$RWAra.bZxOV5Iz6Lv2aE5.40yP/sxsBywwlHCFw6GjEzMebTblJ0K'),(5,NULL,'mat2222','mat2222@email.com','$2a$10$yimSWKFJKzPtMhJof//eV.rB1.pjuTTRJwzI0x4d5qBA1h06aKUHO'),(8,NULL,'mat222222','mat222222@email.com','$2a$10$vhXiT7UlZIrDjYoiG4uLRekqQN0eRBHwRLONzGA5Ds4.uRb4nyGNG'),(9,NULL,'123444','12212243242@email.com','$2a$10$MWSjDtyHM0N.ule07ZoGfueAZTWjCDUCVN/ME7Ua8/p9jcTeYTwDm'),(10,NULL,'mateus2321432525','mateus78@email.com','$2a$10$3t3nAYkIkbCiLOXOub1FQOtayPzZXpSgtsAcH82GxIzOkNTDTcvNO'),(11,NULL,'mateus7777','mateus77777@email.com','$2a$10$GJZJhK5RNFIVXSYuBcfADeha/6HvptT7TNLTsJHleHUnt8e9.iW7G');
/*!40000 ALTER TABLE `users1` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-01-11 13:31:38
