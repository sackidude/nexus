-- @block
SELECT * FROM images LIMIT 10
-- @block
SELECT *
FROM images 
ORDER BY last_analyzed ASC 
LIMIT 10

-- @block
UPDATE Images
SET volume=(1/(analyzed+1))*(analyzed*volume+100), analyzed=analyzed+1, last_analyzed="2024-02-17 15:49:01"
WHERE id=1;

-- @block
UPDATE Images
-- SET volume=0, analyzed=0 -- COMMENTED OUT FOR SAFETY xD;

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
ADD COLUMN last_analyzed DATETIME;

-- @block
UPDATE Images
SET volume=(1/(analyzed+1))*(analyzed*volume+88.2), analyzed=analyzed+1, last_analyzed='2024-02-17 15:57:54'
WHERE id=2;