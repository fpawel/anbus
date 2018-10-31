--C:\Users\fpawel\AppData\Roaming\Аналитприбор\anbus\series.sqlite
PRAGMA foreign_keys = ON;
PRAGMA encoding = 'UTF-8';

CREATE TABLE IF NOT EXISTS bucket
(
  bucket_id  INTEGER   NOT NULL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL UNIQUE DEFAULT (datetime('now')),
  updated_at TIMESTAMP NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS series
(
  bucket_id INTEGER NOT NULL,
  addr      INTEGER NOT NULL CHECK (addr > 0),
  var       INTEGER NOT NULL CHECK (var >= 0),
  stored_at REAL    NOT NULL,
  value     REAL    NOT NULL,
  FOREIGN KEY (bucket_id) REFERENCES bucket (bucket_id)
    ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS trigger_bucket_updated_at
  AFTER INSERT
  ON series
  FOR EACH ROW
  BEGIN
    UPDATE bucket
    SET updated_at = datetime('now')
    WHERE bucket.bucket_id = new.bucket_id;
  END;


CREATE VIEW IF NOT EXISTS bucket_time AS
  SELECT *,
         cast(strftime('%Y', created_at) AS INT) AS year,
         cast(strftime('%m', created_at) AS INT) AS month,
         cast(strftime('%d', created_at) AS INT) AS day
  FROM bucket;

CREATE VIEW IF NOT EXISTS last_bucket AS
  SELECT *
  FROM bucket_time
  ORDER BY created_at DESC
  LIMIT 1;


--SELECT datetime((julianday(current_timestamp)));
--SELECT (julianday(current_timestamp));
--SELECT datetime(2458402.786550926);
--SELECT julianday('now') - julianday('1776-07-04');

