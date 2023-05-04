CREATE TABLE `numbers_allocator` (
                                     `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                     `apply_app_name` varchar(50) DEFAULT NULL COMMENT '申请独占号段的应用名',
                                     `apply_biz_type` varchar(50) DEFAULT NULL COMMENT '申请独占号段的业务类型',
                                     `current_start_id` bigint(20) DEFAULT NULL COMMENT '当前申请号段起始值',
                                     `increment_step` bigint(20) NOT NULL COMMENT '当前号段步长',
                                     `apply_date` varchar(8) DEFAULT NULL COMMENT '号段应用日期',
                                     `version` int(11) DEFAULT NULL COMMENT '版本号',
                                     PRIMARY KEY (`id`),
                                     UNIQUE KEY `idx_apply` (`apply_app_name`,`apply_biz_type`,`apply_date`)
) ENGINE=InnoDB AUTO_INCREMENT=248 DEFAULT CHARSET=utf8mb4;