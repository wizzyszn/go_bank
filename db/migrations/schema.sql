--Drop tables if they exists for development
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;


--  Accounts Table

create table accounts (
id serial primary key,
email varchar(255) unique not null,
password_hash varchar(255) unique not null,
first_name varchar(100) not null,
last_name varchar(100) not null,
balance decimal(15,2) default 0.00,
currency varchar(3) default 'USD',
status varchar(20) default 'active',
created_at timestamp default current_timestamp,
updated_at timestamp default current_timestamp,
constrant valid_balance check (balance is not null and balance >= 0)
)

--  Transaction Table
create table transactions (
id serial primary key,
from_account_id int references accounts(id),
to_account_id int references accounts(id),
amount decimal(15,2) not null,
type varchar(20) not null,
description text,
status varchar(20) default 'completed'
created_at timestamp default current_timestamp,
constraint valid_amount check (amount > 0),
constraint valid_accounts check (from_account_id is not null or to_account_id is not null),
constraint no_same_account check (from_account_id != to_account_id),
)

-- Session Table
create table sessions (
id varchar(255) primary key,
account_id int references accounts(id) on delete cascade,
expires_at timestamp not null
created_at timestamp default current_timestamp
)

--Indexes

create index idx_transactions_from_account on transactions(from_account_id)
create index idx_transactions_to_account on transactions(to_account_id)
create index idx_transactions_created_at on transactions(created_at)
create index idx_sessions_account_id on sessions(account_id);
create index idx_sessions_expires_at on sessions(expires_at);