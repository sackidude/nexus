CREATE TABLE IF NOT EXISTS trials (
    trial_num SMALLINT PRIMARY KEY,
    start_time TIMESTAMP WITHOUT TIME ZONE,
    zero_height FLOAT,
    ml_per_pixel FLOAT,
    yeast_amount FLOAT,
    sugar_amount FLOAT,
    stirring BOOLEAN
);
CREATE TABLE IF NOT EXISTS images (
    id SERIAL PRIMARY KEY,
    file_num SMALLINT NOT NULL,
    trial_num SMALLINT REFERENCES trials (trial_num),
    time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    volume FLOAT DEFAULT 0,
    times_analyzed SMALLINT NOT NULL DEFAULT 0,
    last_requested TIMESTAMP WITH TIME ZONE NULL,
    state char(1) NOT NULL DEFAULT 'U' -- U=unlabeled I=In progress/requested D=Done
);