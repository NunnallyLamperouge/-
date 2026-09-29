SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `chat_groups`;
CREATE TABLE `chat_groups` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL COMMENT '群组名称',
  `description` text COMMENT '群组描述',
  `creator_id` bigint unsigned NOT NULL COMMENT '创建者ID',
  `avatar_url` varchar(255) DEFAULT NULL COMMENT '群组头像',
  `member_count` bigint DEFAULT '1' COMMENT '成员数量',
  PRIMARY KEY (`id`),
  KEY `idx_chat_groups_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `fake_templates`;
CREATE TABLE `fake_templates` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` longtext,
  `type` longtext,
  `description` longtext,
  `user_id` bigint unsigned DEFAULT NULL,
  `content` longtext,
  `code` varchar(30) DEFAULT NULL,
  `slug` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fake_templates_deleted_at` (`deleted_at`),
  KEY `fk_fake_templates_user` (`user_id`),
  CONSTRAINT `fk_fake_templates_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `group_members`;
CREATE TABLE `group_members` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `group_id` bigint unsigned NOT NULL COMMENT '群组ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `nickname` varchar(50) DEFAULT NULL COMMENT '群内昵称',
  `join_time` datetime(3) NOT NULL COMMENT '加入时间',
  `is_admin` tinyint(1) DEFAULT '0' COMMENT '是否为管理员',
  `status` bigint DEFAULT '1' COMMENT '状态 1:正常 2:禁言 3:退出',
  PRIMARY KEY (`id`),
  KEY `idx_group_members_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `guides`;
CREATE TABLE `guides` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `student_uer_id` bigint NOT NULL,
  `teacher_uer_id` bigint NOT NULL,
  `status` bigint NOT NULL,
  `message` text COMMENT '申请理由',
  `student_user_id` bigint unsigned NOT NULL,
  `teacher_user_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_guides_deleted_at` (`deleted_at`),
  KEY `fk_guides_student` (`student_user_id`),
  KEY `fk_guides_teacher` (`teacher_user_id`),
  CONSTRAINT `fk_guides_student` FOREIGN KEY (`student_user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_guides_teacher` FOREIGN KEY (`teacher_user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `image_uploads`;
CREATE TABLE `image_uploads` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint DEFAULT NULL,
  `path` longtext,
  `size` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_image_uploads_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `type` bigint NOT NULL COMMENT '消息类型',
  `content` text NOT NULL COMMENT '消息内容',
  `sender_id` bigint unsigned NOT NULL COMMENT '发送者ID',
  `receiver_id` bigint unsigned DEFAULT NULL COMMENT '接收者ID',
  `group_id` bigint unsigned DEFAULT NULL COMMENT '群组ID',
  `is_read` tinyint(1) DEFAULT '0' COMMENT '是否已读',
  `read_time` datetime(3) DEFAULT NULL COMMENT '阅读时间',
  `media_type` varchar(20) DEFAULT NULL COMMENT '媒体类型',
  `media_url` varchar(255) DEFAULT NULL COMMENT '媒体URL',
  PRIMARY KEY (`id`),
  KEY `idx_messages_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `missions`;
CREATE TABLE `missions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` longtext NOT NULL,
  `description` longtext NOT NULL,
  `student_user_id` bigint unsigned NOT NULL,
  `teacher_user_id` bigint unsigned NOT NULL,
  `status` bigint NOT NULL,
  `end_time` datetime(3) NOT NULL,
  `start_time` datetime(3) NOT NULL,
  `priority` tinyint DEFAULT '0',
  `template_id` bigint unsigned NOT NULL,
  `content` longtext,
  `feedback` text,
  `last_reminder_time` datetime(3) DEFAULT NULL COMMENT '上次提醒时间',
  `reminder_sent` tinyint(1) DEFAULT '0',
  `submission_time` datetime(3) DEFAULT NULL COMMENT '提交时间',
  `review_time` datetime(3) DEFAULT NULL COMMENT '审核时间',
  `rejection_reason` text,
  `revision_count` bigint DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_missions_deleted_at` (`deleted_at`),
  KEY `fk_missions_template` (`template_id`),
  KEY `fk_missions_student` (`student_user_id`),
  KEY `fk_missions_teacher` (`teacher_user_id`),
  CONSTRAINT `fk_missions_student` FOREIGN KEY (`student_user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_missions_teacher` FOREIGN KEY (`teacher_user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_missions_template` FOREIGN KEY (`template_id`) REFERENCES `fake_templates` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `notification_settings`;
CREATE TABLE `notification_settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `system_notification` tinyint(1) DEFAULT '1' COMMENT '是否接收系统通知',
  `like_notification` tinyint(1) DEFAULT '1' COMMENT '是否接收点赞通知',
  `comment_notification` tinyint(1) DEFAULT '1' COMMENT '是否接收评论通知',
  `follow_notification` tinyint(1) DEFAULT '1' COMMENT '是否接收关注通知',
  `message_notification` tinyint(1) DEFAULT '1' COMMENT '是否接收私信通知',
  `email_notification` tinyint(1) DEFAULT '0' COMMENT '是否接收邮件通知',
  `push_notification` tinyint(1) DEFAULT '1' COMMENT '是否接收推送通知',
  `thesis_progress_notification` tinyint(1) DEFAULT '1' COMMENT '是否接收论文进度通知',
  `review_feedback_notification` tinyint(1) DEFAULT '1' COMMENT '是否接收导师评审反馈通知',
  `milestone_notification` tinyint(1) DEFAULT '1' COMMENT '是否接收里程碑达成通知',
  `coauthoring_invitation` tinyint(1) DEFAULT '1' COMMENT '是否接收论文合著邀请通知',
  `format_template_update` tinyint(1) DEFAULT '1' COMMENT '是否接收格式模板更新通知',
  `defense_schedule_notification` tinyint(1) DEFAULT '1' COMMENT '是否接收答辩安排通知',
  `thesis_submission_notification` tinyint(1) DEFAULT '1' COMMENT '是否接收论文提交通知',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_notification_settings_user_id` (`user_id`),
  KEY `idx_notification_settings_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `notifications`;
CREATE TABLE `notifications` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `type` bigint NOT NULL COMMENT '通知类型',
  `title` varchar(100) DEFAULT NULL COMMENT '通知标题',
  `content` text COMMENT '通知内容',
  `sender_id` bigint unsigned DEFAULT NULL COMMENT '发送者ID',
  `receiver_id` bigint unsigned NOT NULL COMMENT '接收者ID',
  `is_read` tinyint(1) DEFAULT '0' COMMENT '是否已读',
  `read_time` datetime(3) DEFAULT NULL COMMENT '阅读时间',
  `related_id` bigint unsigned DEFAULT NULL COMMENT '相关内容ID',
  `related_type` varchar(50) DEFAULT NULL COMMENT '相关内容类型',
  `extra_data` text COMMENT '额外数据JSON',
  PRIMARY KEY (`id`),
  KEY `idx_notifications_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `template_histories`;
CREATE TABLE `template_histories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `template_id` bigint unsigned DEFAULT NULL,
  `content` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_template_histories_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(50) DEFAULT NULL,
  `password` varchar(1500) DEFAULT NULL,
  `email` varchar(50) DEFAULT NULL,
  `avatar` varchar(80) DEFAULT NULL,
  `tokens` mediumint DEFAULT NULL,
  `permission` tinyint DEFAULT NULL,
  `group` tinyint DEFAULT NULL,
  `type` tinyint DEFAULT NULL,
  `apartment` varchar(50) DEFAULT NULL,
  `profession` varchar(50) DEFAULT NULL,
  `workspace_id` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `virtual_questions`;
CREATE TABLE `virtual_questions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `virtual_id` bigint DEFAULT NULL,
  `question` longtext,
  `answer` longtext,
  `question_audio_path` longtext,
  `answer_audio_path` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_virtual_questions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `virtuals`;
CREATE TABLE `virtuals` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `template_id` bigint DEFAULT NULL,
  `student_id` bigint unsigned DEFAULT NULL,
  `difficulty` longtext,
  `seconds` bigint DEFAULT NULL,
  `score` bigint DEFAULT NULL,
  `description` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_virtuals_deleted_at` (`deleted_at`),
  KEY `fk_virtuals_user` (`student_id`),
  CONSTRAINT `fk_virtuals_user` FOREIGN KEY (`student_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;
