CREATE UNLOGGED TABLE accounts (
  id SERIAL PRIMARY KEY,
  name VARCHAR(60) NOT NULL,
  account_limit INT,
  balance INT DEFAULT 0,
  CONSTRAINT non_negative_balance CHECK((account_limit + balance) >= 0)
);

CREATE TYPE operation_enum AS ENUM ('c', 'd');

CREATE UNLOGGED TABLE transactions (
  id SERIAL PRIMARY KEY,
  account_id INT,
  amount INT NOT NULL,
  operation operation_enum,
  description VARCHAR(10) NOT NULL,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_account
    FOREIGN KEY(account_id) 
      REFERENCES accounts(id)
);
