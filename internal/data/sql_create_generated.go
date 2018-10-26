package data

const SQLCreate = `
--C:\Users\fpawel\AppData\Roaming\Аналитприбор\panalib\series.sqlite
PRAGMA foreign_keys = ON;
PRAGMA encoding = 'UTF-8';

CREATE TABLE IF NOT EXISTS bucket (
  bucket_id  INTEGER   NOT NULL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL UNIQUE
);


CREATE VIEW IF NOT EXISTS bucket_time AS
  SELECT *, cast(strftime('%Y', created_at) AS INT) AS year,
            cast(strftime('%m', created_at) AS INT) AS month,
            cast(strftime('%d', created_at) AS INT) AS day
  FROM bucket
  ORDER BY created_at;

CREATE TABLE IF NOT EXISTS series (
  bucket_id      INTEGER NOT NULL,
  place INTEGER NOT NULL CHECK (place >= 0),
  var            INTEGER NOT NULL CHECK (var >= 0),
  seconds_offset REAL    NOT NULL,
  value          REAL    NOT NULL,
  UNIQUE (bucket_id, place, var, seconds_offset),
  FOREIGN KEY (bucket_id) REFERENCES bucket (bucket_id)
    ON DELETE CASCADE
);

CREATE VIEW IF NOT EXISTS series_info AS
  SELECT strftime('%d.%m.%Y %H:%M:%f', bucket.created_at, '+' || series.seconds_offset || ' seconds') AS created_at,
         series.bucket_id,
         series.var,
         series.place,
         series.value
  FROM series
         INNER JOIN bucket ON bucket.bucket_id = series.bucket_id;`
