-- @block
CREATE TABLE Trials(
    `trial-number` SMALLINT PRIMARY KEY,
    `directory-name` VARCHAR(10) NOT NULL,
    `start-time` DATETIME NOT NULL,
    `zero-height` SMALLINT NOT NULL,
    `thousand-milliliter-height` SMALLINT NOT NULL,
    `description` TEXT NULL
);

-- @block
CREATE TABLE Images(
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `filename` VARCHAR(10) NOT NULL,
    `trial` SMALLINT NOT NULL,
    `time` DATETIME NOT NULL,
    `volume` FLOAT NULL,
    `times-analyzed` TINYINT NULL DEFAULT '0'
);
ALTER TABLE
    Images ADD CONSTRAINT `image_trial_foreign` FOREIGN KEY(`trial`) REFERENCES Trials(`trial-number`);

-- @block
SELECT id, trial, time FROM images WHERE filename = "1.jpg";