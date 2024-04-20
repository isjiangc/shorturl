-- ginblog.shorturl definition

CREATE TABLE `shorturl` (
                            `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                            `short_url` longtext NOT NULL,
                            `dest_url` longtext NOT NULL,
                            `valid` tinyint(1) NOT NULL DEFAULT '1',
                            `memo` longtext,
                            `open_type` bigint unsigned NOT NULL DEFAULT '0',
                            `created_at` datetime(3) DEFAULT NULL,
                            `updated_at` datetime(3) DEFAULT NULL,
                            `deleted_at` datetime(3) DEFAULT NULL,
                            PRIMARY KEY (`id`),
                            KEY `idx_shorturl_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;