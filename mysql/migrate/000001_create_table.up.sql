CREATE TABLE `launchpad` (
    `id` bigInt PRIMARY KEY AUTO_INCREMENT,
    `hash_value` varchar(255) NOT NULL,
    `first_owner_email` varchar(255) NOT NULL,
    `ca_address` varchar(255) NOT NULL,
    `chain_id` varchar(255) NOT NULL,
    `price` integer NOT NULL,
    `airdrop_address` JSON NOT NULL DEFAULT "[]",
    `whitelist_address` JSON NOT NULL DEFAULT "[]",
    `created_at` varchar(255) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX `launchpad_index_0` ON `launchpad` (`first_owner_email`);
