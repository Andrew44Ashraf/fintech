CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    balance NUMERIC NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id),
    amount NUMERIC NOT NULL,
    type VARCHAR(10) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_account ON transactions(account_id);
CREATE INDEX idx_transactions_timestamp ON transactions(timestamp);
