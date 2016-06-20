CREATE TABLE `logs` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`firstname` varchar(255) DEFAULT NULL,
	`lastname` varchar(255) DEFAULT NULL,
	`date` datetime DEFAULT NULL,
	`type` varchar(255) DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
