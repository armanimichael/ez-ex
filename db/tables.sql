CREATE TABLE IF NOT EXISTS accounts
(
    id                          INTEGER PRIMARY KEY,
    name                        TEXT NOT NULL UNIQUE,
    description                 TEXT,
    initial_balance_in_cents    INTEGER NOT NULL DEFAULT 0,
    balance_in_cents            INTEGER NOT NULL DEFAULT 0,
    delete_date_unix            INTEGER
);

CREATE TABLE IF NOT EXISTS payees
(
    id          INTEGER PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE IF NOT EXISTS categories
(
    id          INTEGER PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    description TEXT
);
INSERT OR IGNORE INTO categories (id, name)
VALUES (0, 'no category');

CREATE TABLE IF NOT EXISTS transactions
(
    id                    INTEGER PRIMARY KEY,
    category_id           INTEGER NOT NULL DEFAULT 0,
    payee_id              INTEGER NOT NULL,
    account_id            INTEGER NOT NULL,
    amount_in_cents       INTEGER NOT NULL,
    transaction_date_unix INTEGER NOT NULL,
    update_date_unix      INTEGER,
    delete_date_unix      INTEGER,
    notes                 TEXT,

    FOREIGN KEY (category_id) REFERENCES categories ON DELETE SET DEFAULT,
    FOREIGN KEY (payee_id) REFERENCES payees ON DELETE RESTRICT,
    FOREIGN KEY (account_id) REFERENCES accounts ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS ix_transactions_by_account_id_transaction_date_unix ON transactions (account_id, transaction_date_unix);
