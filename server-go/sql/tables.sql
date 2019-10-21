/*
SQLyog Ultimate v12.4.1 (64 bit)
MySQL - 5.7.26 : Database - rzas
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`rzas` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `rzas`;

/*Table structure for table `auth_has_permission` */

DROP TABLE IF EXISTS `auth_has_permission`;

CREATE TABLE `auth_has_permission` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `content` varchar(50) NOT NULL,
  `model_id` int(11) NOT NULL COMMENT '模型主键',
  `model_type` varchar(50) NOT NULL COMMENT '模型命名空间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=35 DEFAULT CHARSET=utf8;

/*Data for the table `auth_has_permission` */

insert  into `auth_has_permission`(`id`,`content`,`model_id`,`model_type`) values 
(26,'rule-view',2,'xiaodi\\Permission\\Models\\Role'),
(27,'rule-add',2,'xiaodi\\Permission\\Models\\Role'),
(28,'role-view',2,'xiaodi\\Permission\\Models\\Role'),
(29,'role-add',2,'xiaodi\\Permission\\Models\\Role'),
(30,'account-view',2,'xiaodi\\Permission\\Models\\Role'),
(31,'account-add',2,'xiaodi\\Permission\\Models\\Role'),
(32,'rule-view',3,'xiaodi\\Permission\\Models\\Role'),
(33,'role-view',3,'xiaodi\\Permission\\Models\\Role'),
(34,'account-view',3,'xiaodi\\Permission\\Models\\Role');

/*Table structure for table `auth_permission` */

DROP TABLE IF EXISTS `auth_permission`;

CREATE TABLE `auth_permission` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `action` varchar(50) NOT NULL,
  `name` varchar(100) NOT NULL COMMENT '规则名称',
  `title` varchar(100) NOT NULL,
  `pid` int(11) NOT NULL,
  `status` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=MyISAM AUTO_INCREMENT=17 DEFAULT CHARSET=utf8;

/*Data for the table `auth_permission` */

insert  into `auth_permission`(`id`,`action`,`name`,`title`,`pid`,`status`) values 
(1,'rule','rule','规则管理',0,1),
(2,'view','rule-view','查看',1,1),
(3,'add','rule-add','添加',1,1),
(4,'update','rule-update','更新',1,1),
(5,'delete','rule-delete','删除',1,1),
(6,'role','role','角色管理',0,1),
(7,'view','role-view','查看',6,1),
(8,'add','role-add','添加',6,1),
(9,'update','role-update','更新',6,1),
(10,'delete','role-delete','删除',6,1),
(11,'account','account','管理员管理',0,1),
(12,'view','account-view','查看',11,1),
(13,'add','account-add','添加',11,1),
(14,'update','account-update','更新',11,1),
(15,'delete','account-delete','删除',11,1);

/*Table structure for table `auth_role` */

DROP TABLE IF EXISTS `auth_role`;

CREATE TABLE `auth_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '角色名称',
  `title` varchar(50) NOT NULL,
  `status` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=MyISAM AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

/*Data for the table `auth_role` */

insert  into `auth_role`(`id`,`name`,`title`,`status`) values 
(2,'test','测试',1),
(3,'test2','测试2',1);

/*Table structure for table `auth_role_permission_access` */

DROP TABLE IF EXISTS `auth_role_permission_access`;

CREATE TABLE `auth_role_permission_access` (
  `role_id` int(11) NOT NULL COMMENT '角色主键',
  `permission_id` int(11) NOT NULL COMMENT '规则主键',
  UNIQUE KEY `permission_id` (`permission_id`,`role_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

/*Data for the table `auth_role_permission_access` */

insert  into `auth_role_permission_access`(`role_id`,`permission_id`) values 
(2,2),
(3,2),
(2,3),
(2,4),
(2,7),
(3,7),
(2,8),
(3,8),
(2,9),
(2,12),
(3,12),
(2,13),
(3,13),
(2,14);

/*Table structure for table `auth_user` */

DROP TABLE IF EXISTS `auth_user`;

CREATE TABLE `auth_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `nickname` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(60) COLLATE utf8mb4_unicode_ci NOT NULL,
  `remark` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  `headimg` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '头像',
  `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '邮箱',
  `status` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/*Data for the table `auth_user` */

insert  into `auth_user`(`id`,`name`,`nickname`,`password`,`remark`,`headimg`,`email`,`status`) values 
(1,'admin','超级管理员','$2a$10$HsVrWexilee9vj8BtnjvV.N5PE1Wont8wtalz2KTkdamsdcCssjC6','','static/upload/15716369254046e1b52af2f6e12b937a95f7ce710d.jpg','349119104@qq.com',1),
(2,'test','fanpei','$2a$10$214b0YutP5lbvVN0JJW56O2GCOlmZ5Onjmzmi4FC84vG0s59GpKs2','121313131','static/upload/15716471864046e1b52af2f6e12b937a95f7ce710d.jpg','349119104@qq.com',1);

/*Table structure for table `auth_user_role_access` */

DROP TABLE IF EXISTS `auth_user_role_access`;

CREATE TABLE `auth_user_role_access` (
  `user_id` int(11) NOT NULL COMMENT '用户id',
  `role_id` int(11) NOT NULL COMMENT '角色id',
  UNIQUE KEY `user_id` (`user_id`,`role_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

/*Data for the table `auth_user_role_access` */

insert  into `auth_user_role_access`(`user_id`,`role_id`) values 
(1,1),
(2,2);

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
