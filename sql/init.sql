-- @block
CREATE TABLE Trials(
    trial_num SMALLINT PRIMARY KEY,
    directory VARCHAR(10) NOT NULL,
    start DATETIME NOT NULL,
    zero_height SMALLINT NOT NULL,
    ml_per_pixel FLOAT NOT NULL,
    description TEXT NULL
);

-- @block
CREATE TABLE Images(
    id INT PRIMARY KEY AUTO_INCREMENT,
    filename VARCHAR(10) NOT NULL,
    trial SMALLINT NOT NULL,
    time DATETIME NOT NULL,
    volume FLOAT NOT NULL DEFAULT '0',
    analyzed TINYINT NULL DEFAULT '0',
    last_analyzed DATETIME NULL
);
ALTER TABLE
    Images ADD CONSTRAINT image_trial_foreign FOREIGN KEY(trial) REFERENCES Trials(trial_num);

-- @block THESE ARE THE READINGS FROM MY OWN EXAMINATION
UPDATE trials
SET zero_height=49, ml_per_pixel=4.2
WHERE trial_num=3;

UPDATE trials
SET zero_height=13, ml_per_pixel=6.024
WHERE trial_num=4;

UPDATE trials
SET zero_height=5, ml_per_pixel=4.762
WHERE trial_num=5;

UPDATE trials
SET zero_height=18, ml_per_pixel=5.525
WHERE trial_num=6;
