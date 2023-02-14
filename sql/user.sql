CREATE DATABASE IF NOT EXISTS BYR_Iptables;
USE BYR_Iptables;
CREATE TABLE IF NOT EXISTS `users` (
                         `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
                         `username` varchar(30) NOT NULL COMMENT '用户名',
                         `password` varchar(100) NOT NULL COMMENT '存加密过的密码',
                         `privilege` int(8) NOT NULL DEFAULT 0 COMMENT '用户特权级别',
                         PRIMARY KEY (`id`)
)