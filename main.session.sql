-- @block
CREATE TABLE Trials(
    trial_num SMALLINT PRIMARY KEY,
    directory VARCHAR(10) NOT NULL,
    start DATETIME NOT NULL,
    zero_height SMALLINT NOT NULL,
    1000_ml_height SMALLINT NOT NULL,
    description TEXT NULL
);

-- @block
CREATE TABLE Images(
    id INT PRIMARY KEY AUTO_INCREMENT,
    filename VARCHAR(10) NOT NULL,
    trial SMALLINT NOT NULL,
    time DATETIME NOT NULL,
    volume FLOAT NOT NULL DEFAULT '0',
    analyzed TINYINT NULL DEFAULT '0'
);
ALTER TABLE
    Images ADD CONSTRAINT `image_trial_foreign` FOREIGN KEY(`trial`) REFERENCES Trials(`trial-number`);

-- @block
SELECT * FROM trials WHERE TRUE LIMIT 1

-- @block
UPDATE Images
SET `times-analyzed`=`times-analyzed`+1, volume=(volume
