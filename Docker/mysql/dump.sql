-- Adminer 4.8.1 MySQL 8.0.29 dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

USE `mysql1`;

DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `content` varchar(1000) NOT NULL,
  `c_time` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

INSERT INTO `messages` (`id`, `user_id`, `content`, `c_time`) VALUES
(1,	791987055,	'Привет',	1661606243),
(2,	791987055,	'123456',	1661606373),
(3,	791987055,	'kz kzkz',	1661606755),
(4,	791987055,	'hfjgfjgk',	1661606798),
(5,	791987055,	'ghbdtn',	1661606921);

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint NOT NULL,
  `username` varchar(50) DEFAULT NULL,
  `first_name` varchar(50) DEFAULT NULL,
  `last_name` varchar(50) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;


-- 2022-08-27 13:54:54
