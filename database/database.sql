DROP TABLE IF EXISTS money_accounts;
DROP TABLE IF EXISTS persons;

CREATE TABLE money_accounts (
  id uuid DEFAULT gen_random_uuid (),
  name VARCHAR NOT NULL,
  balance NUMERIC(17,2) DEFAULT 0.00,
  details VARCHAR,
  currency VARCHAR NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE persons (
  id uuid DEFAULT gen_random_uuid (),
  name VARCHAR NOT NULL,
  document VARCHAR,
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW()
);