CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DROP TABLE IF EXISTS trashed_transactions;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS money_accounts;
DROP TABLE IF EXISTS pending_bills;
DROP TABLE IF EXISTS closed_bills;
DROP TABLE IF EXISTS persons;
DROP TABLE IF EXISTS currencies;

CREATE TABLE currencies (
  currency VARCHAR (3) PRIMARY KEY
);

INSERT INTO currencies (currency) VALUES ('VED'), ('USD');

CREATE TABLE money_accounts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  name VARCHAR NOT NULL,
  balance NUMERIC(17,2) DEFAULT 0.00 CHECK (balance >= 0),
  details VARCHAR,
  currency VARCHAR (3) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (currency) REFERENCES currencies(currency)
);

CREATE TABLE persons (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  name VARCHAR NOT NULL,
  document VARCHAR UNIQUE,
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- used in transactions for records without a person
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


-- to pay: has amount negative
-- to charge: has amount positive
CREATE TABLE pending_bills (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  person_id uuid NOT NULL,
  date DATE DEFAULT NOW(),
  description VARCHAR NOT NULL,
  currency VARCHAR (3) NOT NULL,
  amount NUMERIC(17,2) NOT NULL CHECK (amount <> 0),
  -- when pending is zero, the bill is considered to be paid
  pending NUMERIC(17,2) NOT NULL CHECK (pending <> 0), 
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (person_id) REFERENCES persons(id),
  FOREIGN KEY (currency) REFERENCES currencies(currency)
);

CREATE TABLE closed_bills (
  id uuid PRIMARY KEY,
  person_id uuid NOT NULL,
  date DATE DEFAULT NOW(),
  description VARCHAR NOT NULL,
  currency VARCHAR (3) NOT NULL,
  amount NUMERIC(17,2) NOT NULL CHECK (amount <> 0),
  -- when pending is zero, the bill is considered to be paid
  pending NUMERIC(17,2) NOT NULL CHECK (pending = 0), 
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (person_id) REFERENCES persons(id),
  FOREIGN KEY (currency) REFERENCES currencies(currency)
);


