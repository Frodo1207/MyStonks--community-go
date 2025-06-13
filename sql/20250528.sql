CREATE TABLE `users` (
                         id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                         sol_address VARCHAR(128) NOT NULL COMMENT 'sol钱包地址',
                         username VARCHAR(64) DEFAULT NULL COMMENT '昵称，可选',
                         total_points INT NOT NULL DEFAULT 0 COMMENT '用户总积分',

                         created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                         updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                         is_deleted TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除（0否，1是）',

                         PRIMARY KEY (`id`),
                         UNIQUE KEY idx_sol_address (`sol_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

CREATE TABLE `tasks` (
                         id INT NOT NULL COMMENT '任务ID',
                         step INT COMMENT '步骤顺序',
                         title VARCHAR(255) NOT NULL COMMENT '任务标题',
                         description TEXT COMMENT '任务描述',
                         reward INT NOT NULL DEFAULT 0 COMMENT '奖励积分',
                         category VARCHAR(50) NOT NULL COMMENT '任务分类(newbie/daily/other)',
                         icon VARCHAR(100) COMMENT '图标',
                         special_action VARCHAR(100) COMMENT '特殊动作标识',

                         created_by VARCHAR(64) NOT NULL COMMENT '创建人ID',
                         updated_by VARCHAR(64) NOT NULL COMMENT '更新人ID',
                         created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                         updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                         is_deleted TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除（0否，1是）',

                         PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务表';

CREATE TABLE `user_tasks` (
                              id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                              user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
                              sol_address VARCHAR(128) NOT NULL COMMENT '用户Solana钱包地址',
                              task_id INT NOT NULL COMMENT '任务ID',
                              completed_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '完成时间',
                              verified TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已验证（0否，1是）',

                              created_by VARCHAR(64) NOT NULL COMMENT '创建人ID',
                              updated_by VARCHAR(64) NOT NULL COMMENT '更新人ID',
                              created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                              is_deleted TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除（0否，1是）',

                              PRIMARY KEY (`id`),
                              UNIQUE KEY uk_user_task (`user_id`, `task_id`),
                              INDEX idx_sol_address (`sol_address`),
                              INDEX idx_task_id (`task_id`),
                              FOREIGN KEY (`user_id`) REFERENCES `users`(`id`),
                              FOREIGN KEY (`task_id`) REFERENCES `tasks`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户任务完成记录表';

CREATE TABLE `event_record` (
                                id VARCHAR(64) NOT NULL COMMENT '事件记录ID，唯一',
                                event_type VARCHAR(64) NOT NULL COMMENT '事件类型(TASK_COMPLETE/ON_CHAIN_TRANSFER等)',
                                user_id BIGINT UNSIGNED NOT NULL COMMENT '触发事件的用户ID',
                                points_change INT NOT NULL COMMENT '积分变动值',
                                metadata JSON DEFAULT NULL COMMENT '事件具体数据',

                                created_by VARCHAR(64) NOT NULL COMMENT '创建人ID',
                                updated_by VARCHAR(64) NOT NULL COMMENT '更新人ID',
                                created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                is_deleted TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除（0否，1是）',

                                PRIMARY KEY (`id`),
                                KEY idx_user_id (`user_id`),
                                KEY idx_event_type (`event_type`),
                                FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='事件记录表';

CREATE TABLE `telegram_bindings` (
                                     `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                                     `first_name` VARCHAR(255) COMMENT 'Telegram 用户名字',
                                     `addr` VARCHAR(255) NOT NULL COMMENT '用户地址',
                                     `telegram_id` BIGINT NOT NULL COMMENT 'Telegram 用户ID',
                                     `username` VARCHAR(255) COMMENT 'Telegram 用户名',
                                     `photo_url` TEXT COMMENT 'Telegram 头像链接',
                                     `auth_date` BIGINT COMMENT 'Telegram 授权时间戳',
                                     `hash` VARCHAR(255) COMMENT 'Telegram 数据签名 hash',
                                     `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                     `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

                                     PRIMARY KEY (`id`),
                                     UNIQUE KEY `uq_addr` (`addr`),
                                     UNIQUE KEY `uq_telegram_id` (`telegram_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户 Telegram 绑定信息表';