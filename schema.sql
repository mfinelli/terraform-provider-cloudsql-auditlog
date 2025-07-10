-- Copyright (c) Mario Finelli
-- SPDX-License-Identifier: MPL-2.0

CREATE TABLE `audit_log_rules` (
	  `id` bigint NOT NULL AUTO_INCREMENT,
	  `username` varchar(2048) COLLATE utf8mb4_bin NOT NULL,
	  `dbname` varchar(2048) COLLATE utf8mb4_bin NOT NULL,
	  `object` varchar(2048) COLLATE utf8mb4_bin NOT NULL,
	  `operation` varchar(2048) COLLATE utf8mb4_bin NOT NULL,
	  `op_result` char(1) COLLATE utf8mb4_bin NOT NULL,
	  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
