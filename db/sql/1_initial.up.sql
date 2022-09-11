CREATE TABLE `user` (
	`id` VARCHAR(20) NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`password_hash` VARCHAR(255) NOT NULL,
	`admin` TINYINT(1) NOT NULL DEFAULT 0,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `user` VALUES ("cceqj5n6i1e7hgou9lv0", "admin", "$2a$10$4nYuciuXWVGOFGxor1LmPeywjlJycdol0uh73v9cl/xdLMgYWcHg2", 1);

CREATE TABLE `product` (
	`id` VARCHAR(20) NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`serving_type` VARCHAR(15) NOT NULL,
	`serving_size` INTEGER NOT NULL,
	`serving_calories` INTEGER NOT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `recipy` (
	`id` VARCHAR(20) NOT NULL,
	`user_id` VARCHAR(20) NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`description` VARCHAR(1023) NOT NULL,
	PRIMARY KEY (`id`),
	CONSTRAINT `user_fk1` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `recipy_product` (
	`recipy_id` VARCHAR(20) NOT NULL,
	`product_id` VARCHAR(20) NOT NULL,
	`quantity` INTEGER NOT NULL,
	PRIMARY KEY (`recipy_id`, `product_id`),
	CONSTRAINT `recipy_fk1` FOREIGN KEY (`recipy_id`) REFERENCES `recipy` (`id`) ON DELETE CASCADE,
	CONSTRAINT `product_fk1` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE
);

CREATE TABLE `rating` (
	`id` VARCHAR(20) NOT NULL,
	`recipy_id` VARCHAR(20) NOT NULL,
	`user_id` VARCHAR(20) NOT NULL,
	`score` INTEGER(4) NOT NULL DEFAULT 10,
	`description` VARCHAR(511) NOT NULL,
	PRIMARY KEY (`id`),
	CONSTRAINT `recipy_fk2` FOREIGN KEY (`recipy_id`) REFERENCES `recipy` (`id`) ON DELETE CASCADE,
	CONSTRAINT `user_fk2` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
