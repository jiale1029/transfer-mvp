### Installation

Before we start the server up, there's a few dependencies that need to be installed

1. MySQL
2. Go

You can execute the `sh scripts/setup.sh` to check if you have the dependencies installed.

### Spinning up MySQL Server

Afer both Go and MySQL are installed, spin up your mysql server

If you're on macOS, and installed mysql though brew, start your mysql server using `brew services start mysql`
If you're on ubuntu, `sudo systemctl start mysql`

### Creating a database

1. Login to your mysql server, `mysql -u user_name -p password`
2. Creating a database `CREATE DATABASE database_name;`, replace database_name with the name of your database. I'm using `take_home_test` as the name of my database.
3. Execute `show databases;`, you should see the database that you created just now in the list.
4. Execute `show variables like 'port';`, you should see the port, record it down as you will need it later.

### Running your server

1. Before we start the server, modify the database details in `dal/mysql/init.go`.
2. Change the fields to your `username`, `password`, `addr`, `port`, and `name` of your database.
3. Execute `sh scripts/run.sh` to spin up the server.

### Calling the endpoints

#### GET Endpoints

1. `GET /accounts/xxx`, -> retrieve the account with account_id = xxx
2. `GET /accounts/list` -> retrieves all accounts available
3. `GET /transactions/xxx`, -> retrieve the transaction with transaction_id = xxx
4. `GET /transactions/list` -> retrieves all transactions available

#### POST endpoints

1. `POST /accounts`, -> creates an account
2. `POST /transactions`, -> submits a transaction

### Request and Responses

All request and responses format can be found under `api/protocol.go`.

1. CreateAccountRequest (POST /accounts)

```
{
    "account_id": "159",
    "initial_balance": "123.99"
}
```

2. SubmitTransactionRequest (POST /transactions)

```
{
    "source_account_id": "159",
    "destination_account_id": "1029",
    "amount": "32.5"
}
```

You can use `curl localhost:8080/accounts/list` or utility like Postman to test out the endpoints.

### Codebase Layout

The entrypoint of the code is from `main.go`. The codebase has 4 main folder, `api`, `common`, `dal` and `entity`.

- In `main.go`, the handler of routes are registered and logics are within the `api` folder.
- In `api`, it consists of the implementation of each handler of endpoints.
- In `common`, it consists of the error definition and codes,
- In `dal`, it is the implementation of the data access layer, implemented with MySQL.
- In `entity`, these are the objects that are stored in database and some utility converter functions.
- In `scripts`, these are the scripts, you can execute `scripts/test.sh` to execute some tests.

### Assumptions and Limitations

Accounts related:

- Each account can only be created once, based on account_id.
- Balance of an account cannot be negative, if the initial balance is negative, account creation will be rejected.

Transactions related:

- Transaction is done in currency with dollar and cents.
- For precision of cents, the maximum precision that the database will store is 1/100000 -> 0.00001, more than that will be lost.
- If a transaction will cause the balance of an account to be negative, the transaction will be rejected
- A transaction can be in three different states, `In Progress`, `Fail`, `Success`.
- It is possible for transaction to be in `In Progress` if there's error when trying to update to `Success` state.
