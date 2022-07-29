package migrate
/*
sqlStr := `
CREATE TABLE `sys_user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '账号id',
  `account` varchar(31) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '登录账号(工号)',
  `loginpw` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '登录密码 要求客户端md5之后传到服务端做二次校验',
  `email` varchar(63) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '电子邮箱',
  `mobile` varchar(15) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '手机号码',
  `regtime` int(11) NOT NULL COMMENT '注册时间',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态 2-离职 1-在职',
  `expire` int(11) DEFAULT '0' COMMENT '账号过期时间 0-永不过期',
  `isdup` tinyint(1) DEFAULT '0' COMMENT '是否允许多终端登录 1-允许 2-不允许',
  `stime` int(11) NOT NULL COMMENT '最后操作时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_account` (`account`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1275 DEFAULT CHARSET=utf8;

INSERT INTO `sys_user` VALUES (1, 'admin', '1f9abbabf9926d579a3c5d1140421be8', 'lchenchun@sina.com', '13514076806', 1, 1727452800, 1, 1646036519);

CREATE TABLE `sys_msrv` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `srv` varchar(127) DEFAULT '' COMMENT '服务地址',
  `name` varchar(63) DEFAULT '' COMMENT '服务名称',
  `version` varchar(15) DEFAULT '' COMMENT '版本号',
  `proto` varchar(15) DEFAULT '' COMMENT '协议',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态 1-正常 0-失效',
  `stime` int(11) DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_srv` (`srv`),
  KEY `idx_name_proto_status` (`name`,`proto`,`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

CREATE TABLE `sys_yaml` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `name` varchar(31) DEFAULT NULL COMMENT '配置名称',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态 1-正常 0-副本 编辑的时候保存副本',
  `yaml` text COMMENT '配置内容',
  `userid` int(11) DEFAULT '0' COMMENT '编辑的用户ID',
  `calls` int(11) DEFAULT '0' COMMENT '调用统计',
  `version` varchar(255) DEFAULT '' COMMENT '版本号数据',
  `stime` int(11) DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_userid` (`userid`),
  KEY `idx_name_status` (`name`,`status`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=325 DEFAULT CHARSET=utf8mb4;
`
*/
