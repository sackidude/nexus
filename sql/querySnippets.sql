-- @block
SELECT * FROM images LIMIT 10

-- @block
UPDATE Images
SET volume=(1/(analyzed+1))*(analyzed*volume+100), analyzed=analyzed+1
WHERE id=1;

-- @block
UPDATE Images
-- SET volume=0, analyzed=0 COMMENTED OUT FOR SAFETY xD
