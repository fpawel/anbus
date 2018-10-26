--C:\Users\fpawel\AppData\Roaming\Аналитприбор\anbus\series.sqlite
PRAGMA foreign_keys = ON;
PRAGMA encoding = 'UTF-8';

CREATE TABLE IF NOT EXISTS bucket (
  bucket_id  INTEGER   NOT NULL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS series (
  stored_at REAL    NOT NULL, -- DEFAULT (julianday(current_timestamp)),
  bucket_id INTEGER NOT NULL,
  addr      INTEGER NOT NULL CHECK (addr > 0),
  var       INTEGER NOT NULL CHECK (var >= 0),
  value     REAL    NOT NULL,
  FOREIGN KEY (bucket_id) REFERENCES bucket (bucket_id)
    ON DELETE CASCADE
);

CREATE VIEW IF NOT EXISTS last_bucket AS
  SELECT *
  FROM bucket
  ORDER BY created_at DESC
  LIMIT 1;

CREATE VIEW IF NOT EXISTS last_value AS
  SELECT cast(strftime('%Y', datetime(stored_at)) AS INT) AS year,
         cast(strftime('%m', datetime(stored_at)) AS INT) AS month,
         cast(strftime('%d', datetime(stored_at)) AS INT) AS day,
         cast(strftime('%H', datetime(stored_at)) AS INT) AS hour,
         cast(strftime('%M', datetime(stored_at)) AS INT) AS minute,
         cast(strftime('%S', datetime(stored_at)) AS INT) AS second,
         cast(datetime(stored_at) AS TEXT)                AS stored_at,
         bucket_id
  FROM series
  WHERE bucket_id IN (SELECT bucket_id FROM last_bucket)
  ORDER BY stored_at DESC
  LIMIT 1;


CREATE VIEW IF NOT EXISTS bucket_time AS
  SELECT *, cast(strftime('%Y', created_at) AS INT) AS year,
            cast(strftime('%m', created_at) AS INT) AS month,
            cast(strftime('%d', created_at) AS INT) AS day
  FROM bucket
  ORDER BY created_at;

--SELECT datetime((julianday(current_timestamp)));
--SELECT (julianday(current_timestamp));
--SELECT datetime(2458402.786550926);
--SELECT julianday('now') - julianday('1776-07-04');

