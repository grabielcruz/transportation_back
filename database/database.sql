CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS closed_bills;
DROP TABLE IF EXISTS pending_bills;
DROP TABLE IF EXISTS bill_cross;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS persons;
DROP TABLE IF EXISTS money_accounts;
DROP TABLE IF EXISTS currencies;

CREATE TABLE currencies (
  currency VARCHAR (3) PRIMARY KEY
);

INSERT INTO currencies (currency) VALUES ('000'), ('VED'), ('USD');

CREATE TABLE money_accounts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  name VARCHAR NOT NULL,
  balance NUMERIC(17,2) DEFAULT 0.00 CHECK (balance >= 0),
  details VARCHAR NOT NULL,
  currency VARCHAR (3) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (currency) REFERENCES currencies(currency)
);

--zero money account
-- used by zero transaction
INSERT INTO money_accounts (id, name, balance, details, currency) VALUES (uuid_nil(), '', 0, '', '000');

CREATE TABLE persons (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  name VARCHAR NOT NULL,
  document VARCHAR UNIQUE,
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- used in transactions for records without a person
-- zero person
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

-- zero transaction
INSERT INTO transactions (id, account_id, person_id, amount, description, balance) VALUES (uuid_nil(), uuid_nil(), uuid_nil(), 0, '',0);

CREATE TABLE bill_cross (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  currency VARCHAR (3) NOT NULL,
  balance NUMERIC(17,2) NOT NULL,
  FOREIGN KEY (currency) REFERENCES currencies(currency)
);

INSERT INTO bill_cross (id, currency, balance) VALUES (uuid_nil(), '000', 0);

-- to pay: has amount negative
-- to charge: has amount positive
CREATE TABLE pending_bills (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  person_id uuid NOT NULL,
  date DATE DEFAULT NOW(),
  description VARCHAR NOT NULL,
  currency VARCHAR (3) NOT NULL,
  amount NUMERIC(17,2) NOT NULL CHECK (amount <> 0),
   -- both can be null
  parent_transaction_id uuid,
  parent_bill_cross_id uuid,
  --
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (person_id) REFERENCES persons(id),
  FOREIGN KEY (currency) REFERENCES currencies(currency),
  FOREIGN KEY (parent_transaction_id) REFERENCES transactions(id),
  FOREIGN KEY (parent_bill_cross_id) REFERENCES bill_cross(id)
);

CREATE TABLE closed_bills (
  id uuid PRIMARY KEY,
  person_id uuid NOT NULL,
  date DATE DEFAULT NOW(),
  description VARCHAR NOT NULL,
  currency VARCHAR (3) NOT NULL,
  amount NUMERIC(17,2) NOT NULL CHECK (amount <> 0),
   -- both can be null
  parent_transaction_id uuid,
  parent_bill_cross_id uuid,
  --
  -- one of these should be not null
  transaction_id uuid,
  bill_cross_id uuid,
  --
  post_notes VARCHAR,
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (person_id) REFERENCES persons(id),
  FOREIGN KEY (currency) REFERENCES currencies(currency),
  FOREIGN KEY (parent_transaction_id) REFERENCES transactions(id),
  FOREIGN KEY (parent_bill_cross_id) REFERENCES bill_cross(id),
  FOREIGN KEY (transaction_id) REFERENCES transactions(id),
  FOREIGN KEY (bill_cross_id) REFERENCES bill_cross(id)
);
