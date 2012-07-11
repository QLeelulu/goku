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
INSERT INTO `todo` VALUES ('4', '增加摘要功能，除标题外还可以写描述', '0', '2011-06-03 06:08:10');
INSERT INTO `todo` VALUES ('5', '增加分类，不同的事项划分到不同的分类中去', '0', '2011-06-03 06:08:37');
INSERT INTO `todo` VALUES ('6', '增加协作功能，增加用户功能，可以指定转给某个人，且用邮件通知他', '0', '2011-06-03 06:09:18');
INSERT INTO `todo` VALUES ('7', '不直接删除，改为完成，标示为该条事项已完成显示在最下方', '0', '2011-06-03 06:09:41');
INSERT INTO `todo` VALUES ('8', '这是一条测试', '1', '2011-06-04 23:00:47');
INSERT INTO `todo` VALUES ('9', '这是一条测试2', '1', '2011-06-04 23:01:31');


CREATE TABLE `status` (
  `id` int(11) NOT NULL auto_increment,
  `text` varchar(300) default NULL,
  `created_at` datetime default NULL,
  `user_name` varchar(200),
  `json` longtext,
  PRIMARY KEY  (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;