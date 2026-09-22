
CREATE TABLE  teams (
  id SERIAL PRIMARY KEY,
  ocb_id TEXT NOT NULL,
  series TEXT NOT NULL,
  name TEXT,
  short_name TEXT,
  full_name TEXT,
  UNIQUE (ocb_id, series)
);

CREATE TABLE drivers (
  id SERIAL PRIMARY KEY,
  ocb_id TEXT NOT NULL,
  series TEXT NOT NULL,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  number INT,
  nationality TEXT,
  team_id INT REFERENCES teams(id),
  UNIQUE (ocb_id, series)
);

CREATE TABLE seasons (
  id SERIAL PRIMARY KEY,
  ocb_id TEXT NOT NULL,
  series TEXT NOT NULL,
  year INT NOT NULL,
  round_count INT,
  status TEXT,
  UNIQUE (ocb_id, series)
);

CREATE TABLE events (
  id SERIAL PRIMARY KEY,
  ocb_id TEXT NOT NULL,
  series TEXT NOT NULL,
  name TEXT NOT NULL,
  date_start DATE,
  date_end DATE,
  status TEXT,
  UNIQUE (ocb_id, series)
);

CREATE TABLE sessions (
  id SERIAL PRIMARY KEY,
  ocb_id TEXT NOT NULL UNIQUE,
  event_id INT NOT NULL REFERENCES events(id),
  type TEXT,
  name TEXT,
  status TEXT,
  start_time TIMESTAMPTZ,
  end_time TIMESTAMPTZ
);

CREATE TABLE results (
  id SERIAL PRIMARY KEY,
  session_id INT NOT NULL REFERENCES sessions(id),
  driver_id INT NOT NULL REFERENCES drivers(id),
  team_id INT REFERENCES teams(id),
  car_number TEXT,
  position TEXT,
  points NUMERIC,
  status TEXT,
  laps INT,
  UNIQUE (session_id, driver_id)
);
  


