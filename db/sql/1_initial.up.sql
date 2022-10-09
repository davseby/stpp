CREATE TABLE `users` (
	`id` VARCHAR(20) NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`password_hash` VARCHAR(255) NOT NULL,
	`admin` TINYINT(1) NOT NULL DEFAULT 0,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `users` (`id`, `name`, `password_hash`, `admin`, `created_at`) VALUES 
("cceqj5n6i1e7hgou9lv0", "admin", "$2a$10$4nYuciuXWVGOFGxor1LmPeywjlJycdol0uh73v9cl/xdLMgYWcHg2", 1, "2022-09-14 12:19:57"),
("cciuk0n6i1e0q9d6pnf0", "user", "$2a$10$xLMHEDwS8M9I9P8dinUPH.qtitZ8tkDjxDq5ZTavX7pViImFzBL7W", 0, "2022-09-12 16:09:23"),
("cciuk2v6i1e0qjrh0hu0", "david", "$2a$10$1jxLIAXDsow2u5ebbgd.k.LqA8QDb9MYbfEkYmXxHGIgU1CpuqtS6", 0, "2022-09-15 04:52:10");

CREATE TABLE `products` (
	`id` VARCHAR(20) NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`serving_type` VARCHAR(15) NOT NULL,
	`serving_size` DECIMAL(18, 4) NOT NULL,
	`serving_calories` INTEGER NOT NULL,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `products` (`id`, `name`, `serving_type`, `serving_size`, `serving_calories`, `created_at`) VALUES 
("cciudtv6i1e0cuscb5jg", "Probio Oats", "grams", "100", "307", "2022-09-14 04:52:10"),
("cciueof6i1e0du7sgp6g", "Banana", "units", "1", "105", "2022-09-15 23:52:10"),
("cciuf5f6i1e0e49j5750", "Cashew Milk", "milliliters", "100", "23", "2022-09-15 11:11:11"),
("cciufd76i1e0ea44drqg", "Stevia", "grams", "100", "2", "2022-09-14 22:22:22"),
("cciufqn6i1e0faaqor4g", "Chocolate Protein Powder", "grams", "25", "110", "2022-09-15 14:32:20"),
("cciug9v6i1e0h9q4e570", "Brown Rice", "grams", "100", "375", "2022-09-14 20:12:15"),
("cciughv6i1e0i77h6vcg", "Chicken Breast", "grams", "100", "165", "2022-09-15 19:53:42"),
("cciugvv6i1e0ip637jpg", "Olive Oil", "milliliters", "100", "900", "2022-09-16 15:32:14"),
("cciuhfv6i1e0jaucp6k0", "Sugar", "grams", "100", "387", "2022-09-15 05:42:16"),
("cciuhj76i1e0jhb8uasg", "Salt", "grams", "100", "0", "2022-09-15 03:52:10"),
("cciuibn6i1e0kff2evfg", "Apple", "units", "1", "95", "2022-09-15 00:30:12"),
("cciulu76i1e0vjpd3i2g", "Nestle Fitness Cereal", "grams", "100", "379", "2022-09-15 18:52:19"),
("cciumsn6i1e11lv61jm0", "Cocumber Spread", "grams", "100", "50", "2022-09-15 09:52:10"),
("cciunif6i1e12b5egfp0", "Rice cake", "units", "1", "20", "2022-09-15 06:59:10"),
("cciuo7n6i1e13guo332g", "Turkey hot smoked fillet", "grams", "100", "95", "2022-09-15 05:52:10");


CREATE TABLE `recipes` (
	`id` VARCHAR(20) NOT NULL,
	`user_id` VARCHAR(20),
	`name` VARCHAR(255) NOT NULL,
	`description` VARCHAR(1023) NOT NULL,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`),
	CONSTRAINT `recipes_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `recipes` (`id`, `user_id`, `name`, `description`, `created_at`) VALUES 
("cciuk9v6i1e0rha6m580", "cciuk2v6i1e0qjrh0hu0", "Overnight Oats", "Healthy morning breakfast.", "2022-09-17 09:10:15"),
("cciuomn6i1e14du2lbe0", "cciuk0n6i1e0q9d6pnf0", "Cocumber spread sandwitches", "Healthy, high protein snack.", "2022-09-18 23:29:00"),
("ccjlpaf6i1e7tbfsih3g", NULL, "Baked Apples", "Healthy desert.", "2022-09-18 23:59:59");

CREATE TABLE `recipe_products` (
	`recipe_id` VARCHAR(20) NOT NULL,
	`product_id` VARCHAR(20) NOT NULL,
	`quantity` DECIMAL(18, 4) NOT NULL,
	PRIMARY KEY (`recipe_id`, `product_id`),
	CONSTRAINT `recipe_products_recipe_fk` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`id`) ON DELETE CASCADE,
	CONSTRAINT `recipe_products_product_fk` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `recipe_products` (`recipe_id`, `product_id`, `quantity`) VALUES 
("cciuk9v6i1e0rha6m580", "cciudtv6i1e0cuscb5jg", "0.33"),
("cciuk9v6i1e0rha6m580", "cciueof6i1e0du7sgp6g", "0.5"),
("cciuk9v6i1e0rha6m580", "cciuf5f6i1e0e49j5750", "2.8"),
("cciuk9v6i1e0rha6m580", "cciufd76i1e0ea44drqg", "0.4"),
("cciuk9v6i1e0rha6m580", "cciufqn6i1e0faaqor4g", "1"),
("cciuk9v6i1e0rha6m580", "cciulu76i1e0vjpd3i2g", "0.33"),
("cciuomn6i1e14du2lbe0", "cciumsn6i1e11lv61jm0", "0.5"),
("cciuomn6i1e14du2lbe0", "cciunif6i1e12b5egfp0", "4"),
("cciuomn6i1e14du2lbe0", "cciuo7n6i1e13guo332g", "1.3"),
("ccjlpaf6i1e7tbfsih3g", "cciuibn6i1e0kff2evfg", "5"),
("ccjlpaf6i1e7tbfsih3g", "cciugvv6i1e0ip637jpg", "0.2"),
("ccjlpaf6i1e7tbfsih3g", "cciuhfv6i1e0jaucp6k0", "2");

CREATE TABLE `plans` (
	`id` VARCHAR(20) NOT NULL,
	`user_id` VARCHAR(20),
	`name` VARCHAR(255) NOT NULL,
	`description` VARCHAR(1023) NOT NULL,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`),
	CONSTRAINT `plans_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `plans` (`id`, `user_id`, `name`, `description`, `created_at`) VALUES 
("cciupgf6i1e15grnp640", "cciuk2v6i1e0qjrh0hu0", "Healthy lifestyle plan.", "Under 2000 calories per day.", "2022-09-18 23:40:00"),
("ccjlq5f6i1e82gbco1jg", "cceqj5n6i1e7hgou9lv0", "A life made to enjoy.", "Sweet.", "2022-09-19 00:44:24");

CREATE TABLE `plan_recipes` (
	`plan_id` VARCHAR(20) NOT NULL,
	`recipe_id` VARCHAR(20) NOT NULL,
	`quantity` INTEGER UNSIGNED NOT NULL,
	PRIMARY KEY (`plan_id`, `recipe_id`),
	CONSTRAINT `plan_recipes_plan_id_fk` FOREIGN KEY (`plan_id`) REFERENCES `plans` (`id`) ON DELETE CASCADE,
	CONSTRAINT `plan_recipes_recipe_id_fk` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `plan_recipes` (`plan_id`, `recipe_id`, `quantity`) VALUES 
("cciupgf6i1e15grnp640", "cciuk9v6i1e0rha6m580", "1"),
("cciupgf6i1e15grnp640", "cciuomn6i1e14du2lbe0", "3"),
("ccjlq5f6i1e82gbco1jg", "cciuomn6i1e14du2lbe0", "4"),
("ccjlq5f6i1e82gbco1jg", "ccjlpaf6i1e7tbfsih3g", "2");
