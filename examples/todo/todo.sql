/*
MySQL Data Transfer
Source Host: localhost
Source Database: todo
Target Host: localhost
Target Database: todo
Date: 2011/6/25 21:57:10
*/

SET FOREIGN_KEY_CHECKS=0;
-- ----------------------------
-- Table structure for todo
-- ----------------------------
DROP TABLE IF EXISTS `todo`;
CREATE TABLE `todo` (
  `id` int(11) NOT NULL auto_increment,
  `title` varchar(300) default NULL,
  `finished` int(11) default '0',
  `post_date` datetime default NULL,
  PRIMARY KEY  (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records 
-- ----------------------------
INSERT INTO `todo` VALUES ('5', 'add category', '0', '2011-06-03 06:08:37');
INSERT INTO `todo` VALUES ('6', 'add share with other people', '0', '2011-06-03 06:09:18');
INSERT INTO `todo` VALUES ('8', 'this is thest', '1', '2011-06-04 23:00:47');
INSERT INTO `todo` VALUES ('9', 'find my love', '1', '2013-01-01 00:00:00');


CREATE TABLE `status` (
  `id` int(11) NOT NULL auto_increment,
  `text` varchar(300) default NULL,
  `created_at` datetime default NULL,
  `user_name` varchar(200),
  `json` longtext,
  PRIMARY KEY  (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;