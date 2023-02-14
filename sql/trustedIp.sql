CREATE DATABASE IF NOT EXISTS BYR_Iptables;
USE BYR_Iptables;
CREATE TABLE IF NOT EXISTS `trustedip` (
                                       `id` int unsigned NOT NULL  AUTO_INCREMENT,
                                       `eip` varchar(20) NOT NULL DEFAULT 0 COMMENT 'ipv4地址字符串',
                                       PRIMARY KEY (`id`)
);
INSERT INTO trustedip (eip) values ('127.0.0.1');