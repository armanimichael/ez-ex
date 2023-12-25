# ez-ex

Easy expenses' tracker.

Supports Linux system, not guaranteed to work with any other system.

## User-data

User-data is saved into `~/.ez-ex/user-data.db` in a SQLite3 DB.

## CLI App

`make build-cli` will compile the CLI application (then found inside `./out/ez-ex`).

Run `ez-ex -help` for commands.

### Features

- Manage account
    - [x] CLI
        - Create/Soft-Delete accounts
        - View transactions
        - Create/Soft-Delete transactions
        - Upsert Categories / Payees during transaction creation
    - [ ] Web
    - [ ] Mobile App

### Future ideas

#### Any App type

- [ ] Scheduled operations (transactions)
- [ ] Language selection
- [ ] Currency selection
- [ ] Visualize soft-deleted records and hard-delete them if necessary
- [ ] Create backups
- Data Visualization (per time period or absolute)
    - [ ] Earnings vs Expenses
    - [ ] Trends
    - [ ] Hot categories / payees
- [ ] Include / Skip accounts in overviews

#### CLI

For now the CLI app is quite minimal, I may add more functionalities in the future.

- [ ] Manage Categories and Payees
- [ ] Update accounts
- [ ] Update transactions