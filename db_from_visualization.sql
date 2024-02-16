CREATE TABLE `image`(
    `image-number` INT NOT NULL,
    `trial` SMALLINT NOT NULL,
    `time` DATETIME NOT NULL,
    `volume` DOUBLE NULL,
    `times-analyzed` TINYINT NOT NULL DEFAULT '0',
    `filename` VARCHAR(255) NOT NULL
);
ALTER TABLE
    `image` ADD PRIMARY KEY(`image-number`);
CREATE TABLE `trial`(
    `trial-number` SMALLINT NOT NULL,
    `directory-name` VARCHAR(255) NOT NULL,
    `start-time` DATETIME NOT NULL,
    `description` TEXT NULL
);
ALTER TABLE
    `trial` ADD PRIMARY KEY(`trial-number`);
ALTER TABLE
    `image` ADD CONSTRAINT `image_trial_foreign` FOREIGN KEY(`trial`) REFERENCES `trial`(`trial-number`);