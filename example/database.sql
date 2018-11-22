DROP TABLE IF EXISTS `biz_activity`;
CREATE TABLE `biz_activity` (
  `id` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `uuid` varchar(50) DEFAULT NULL,
  `pc_link` varchar(255) DEFAULT NULL,
  `h5_link` varchar(255) DEFAULT NULL,
  `remark` varchar(255) DEFAULT NULL,
  `create_time` datetime NOT NULL,
  `delete_flag` int(1) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='test';

-- ----------------------------
-- Records of biz_activity
-- ----------------------------
INSERT INTO `biz_activity` VALUES ('165', '安利一波大表哥', null, 'http://www.baidu.com', 'http://www.taobao.com', 'ceshi', '2018-05-23 15:21:22', '0');
INSERT INTO `biz_activity` VALUES ('166', '注册送好礼', null, '', 'www.baidu.com', '测试', '2018-05-24 10:36:31', '0');
INSERT INTO `biz_activity` VALUES ('167', 'hello', null, 'www.baidu.com', 'www.baidu.com', 'ceshi', '2018-05-24 10:41:17', '0');
INSERT INTO `biz_activity` VALUES ('168', 'rs168', null, null, null, null, '0000-00-00 00:00:00', '0');
INSERT INTO `biz_activity` VALUES ('169', 'rs168', null, null, null, null, '0000-00-00 00:00:00', '0');
INSERT INTO `biz_activity` VALUES ('170', 'rs168-10', null, null, null, null, '0000-00-00 00:00:00', '1');