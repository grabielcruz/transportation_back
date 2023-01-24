CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS closed_bills CASCADE;
DROP TABLE IF EXISTS pending_bills CASCADE;
DROP TYPE IF EXISTS bill_status CASCADE;
DROP TABLE IF EXISTS bill_cross CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS person_accounts CASCADE;
DROP TABLE IF EXISTS persons CASCADE;
DROP TABLE IF EXISTS money_accounts CASCADE;
DROP TABLE IF EXISTS currencies CASCADE;

CREATE TABLE currencies (
  currency VARCHAR (3) PRIMARY KEY
);

-- zero currency is 000
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

CREATE TABLE person_accounts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  person_id uuid NOT NULL,
  name VARCHAR NOT NULL,
  description VARCHAR NOT NULL,
  currency VARCHAR (3) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (person_id) REFERENCES persons(id),
  FOREIGN KEY (currency) REFERENCES currencies(currency)
);

CREATE TABLE transactions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  account_id uuid NOT NULL,
  person_id uuid NOT NULL,
  person_account_id uuid,
  person_account_name VARCHAR,
  person_account_description VARCHAR,
  date DATE DEFAULT NOW(),
  amount NUMERIC(17,2) NOT NULL,
  fee NUMERIC(17,2) DEFAULT 0.00 CHECK (fee >= 0),
  amount_with_fee NUMERIC(17,2) NOT NULL,
  description VARCHAR NOT NULL,
  balance NUMERIC(17,2) NOT NULL CHECK (balance >= 0),
  pending_bill_id uuid DEFAULT uuid_nil(),
  closed_bill_id uuid DEFAULT uuid_nil(),
  revert_bill_id uuid DEFAULT uuid_nil(),
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (account_id) REFERENCES money_accounts(id),
  FOREIGN KEY (person_id) REFERENCES persons(id)
);

-- zero transaction
INSERT INTO transactions (id, account_id, person_id, amount, fee, amount_with_fee, description, balance) VALUES (uuid_nil(), uuid_nil(), uuid_nil(), 0, 0, 0, '',0);

CREATE TABLE bill_cross (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  person_id uuid NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  currency VARCHAR (3) NOT NULL,
  balance NUMERIC(17,2) NOT NULL,
  FOREIGN KEY (person_id) REFERENCES persons(id),
  FOREIGN KEY (currency) REFERENCES currencies(currency)
);

INSERT INTO bill_cross (id, person_id, currency, balance) VALUES (uuid_nil(), uuid_nil(), '000', 0);

CREATE TYPE bill_status AS ENUM ('PENDING', 'SOLVED', 'REVERTED', 'GROUPED');

-- to pay: has amount negative
-- to charge: has amount positive
CREATE TABLE pending_bills (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  person_id uuid NOT NULL,
  date DATE DEFAULT NOW(),
  description VARCHAR NOT NULL,
  status bill_status DEFAULT 'PENDING',
  currency VARCHAR (3) NOT NULL,
  amount NUMERIC(17,2) NOT NULL,
   -- both can be null
  parent_transaction_id uuid DEFAULT uuid_nil(),
  parent_bill_cross_id uuid DEFAULT uuid_nil(),
  --
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (person_id) REFERENCES persons(id),
  FOREIGN KEY (currency) REFERENCES currencies(currency),
  FOREIGN KEY (parent_transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
  FOREIGN KEY (parent_bill_cross_id) REFERENCES bill_cross(id)
);

INSERT INTO pending_bills (id, person_id, description, currency, amount) VALUES (uuid_nil(), uuid_nil(), '', '000', 0);


CREATE TABLE closed_bills (
  id uuid PRIMARY KEY,
  person_id uuid NOT NULL,
  date DATE DEFAULT NOW(),
  description VARCHAR NOT NULL,
  status bill_status DEFAULT 'SOLVED',
  currency VARCHAR (3) NOT NULL,
  amount NUMERIC(17,2) NOT NULL,
   -- both can be null
  parent_transaction_id uuid,
  parent_bill_cross_id uuid,
  --
  -- one of these should be not null
  transaction_id uuid,
  bill_cross_id uuid,
  revert_transaction_id uuid,
  --
  post_notes VARCHAR,
  created_at TIMESTAMPTZ DEFAULT NOW(), 
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (person_id) REFERENCES persons(id),
  FOREIGN KEY (currency) REFERENCES currencies(currency),
  FOREIGN KEY (parent_transaction_id) REFERENCES transactions(id),
  FOREIGN KEY (parent_bill_cross_id) REFERENCES bill_cross(id),
  FOREIGN KEY (transaction_id) REFERENCES transactions(id),
  FOREIGN KEY (bill_cross_id) REFERENCES bill_cross(id),
  FOREIGN KEY (revert_transaction_id) REFERENCES transactions(id)
);

INSERT INTO closed_bills (id, person_id, description, currency, amount, transaction_id, bill_cross_id)
  VALUES (uuid_nil(), uuid_nil(), '', '000', 0, uuid_nil(), uuid_nil());

ALTER TABLE transactions
  ADD CONSTRAINT fk_transactions_pending_bills FOREIGN KEY (pending_bill_id) REFERENCES pending_bills (id);
ALTER TABLE transactions
  ADD CONSTRAINT fk_transactions_closed_bills FOREIGN KEY (closed_bill_id) REFERENCES closed_bills (id);
ALTER TABLE transactions
  ADD CONSTRAINT fk_revert_closed_bills FOREIGN KEY (revert_bill_id) REFERENCES closed_bills (id);
