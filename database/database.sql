CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DROP TABLE IF EXISTS trashed_transactions;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS money_accounts;
DROP TABLE IF EXISTS persons;
DROP TABLE IF EXISTS currency;

CREATE TABLE money_accounts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  name VARCHAR NOT NULL,
  balance NUMERIC(17,2) DEFAULT 0.00 CHECK (balance >= 0),

  details VARCHAR,
  currency VARCHAR NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE persons (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  name VARCHAR NOT NULL,
  document VARCHAR,
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO persons (id, name, document) VALUES (uuid_nil(), '', '');

CREATE TABLE transactions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  account_id uuid NOT NULL,
  person_id uuid,
  date DATE DEFAULT NOW(),
  amount NUMERIC(17,2) NOT NULL,
  description VARCHAR NOT NULL,
  balance NUMERIC(17,2) NOT NULL CHECK (balance >= 0),
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (account_id) REFERENCES money_accounts(id),
  FOREIGN KEY (person_id) REFERENCES persons(id)
);

CREATE TABLE trashed_transactions (
  id uuid PRIMARY KEY,
  account_id uuid NOT NULL,
  person_id uuid,
  date DATE DEFAULT NOW(),
  amount NUMERIC(17,2) NOT NULL,
  description VARCHAR NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (account_id) REFERENCES money_accounts(id),
  FOREIGN KEY (person_id) REFERENCES persons(id)
);

CREATE TABLE currencies (
  currency VARCHAR (3) PRIMARY KEY
);
