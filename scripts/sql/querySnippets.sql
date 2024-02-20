-- @block
SELECT * FROM images LIMIT 10
-- @block
SELECT *
FROM images 
ORDER BY last_seen ASC 
LIMIT 10

-- @block
UPDATE Images
SET volume=(1/(analyzed+1))*(analyzed*volume+100), analyzed=analyzed+1
WHERE id=1;

-- @block
UPDATE Images
SET volume=0, analyzed=0, request_date=NULL, state='U'

-- @block
SELECT zero_height, 1000_ml_height FROM Trials WHERE trial_num=3 LIMIT 1

-- @block
ALTER TABLE Trials 
RENAME COLUMN 1000_ml_height TO ml_per_pixel;

-- @block
ALTER TABLE Trials
MODIFY ml_per_pixel FLOAT;

-- @block
ALTER TABLE Images
ADD COLUMN last_seen DATETIME;

-- @block
ALTER TABLE Images
RENAME COLUMN last_analyzed TO last_seen

-- @block
UPDATE Images
SET last_seen="2024-02-17 18:13:01"
WHERE id=1;

-- @block
SELECT trial_num FROM Trials

-- @block
SELECT time, volume FROM Images
WHERE trial=3

-- @block
ALTER TABLE Images
RENAME COLUMN last_seen TO request_date

-- @block
ALTER TABLE Images
ADD state CHAR(1) NOT NULL DEFAULT 'U'

-- @block
ALTER TABLE Images
DROP COLUMN state

-- @block
SELECT trial, time, volume, state FROM Images