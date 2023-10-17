# Go templates

![alt text](fandom.jpeg)

# Tour of Go

## Build and run

 ```
 # Build Cli
 $ go build -o ./cli cmd/main.go
 
 $ ./cli
 
 # Build server
 $ go build -o server main.go
 
 $ ./server
 ```   

### Run CLI

```
$ ./cli create-competition 
-> do something
$ ./cli create-team
-> do something
```

### Run server

> Required to run a project

Config environment variables ->

Edit `config/config.json`

```json
{
  "env": "development",
  "port": 8080,
  "db": "{user}:{password}@tcp({host}:{port})/{database_name}?charset=utf8&parseTime=True&charset=utf8mb4&collation=utf8mb4_unicode_ci"
}
```

Example:
```config/config.json```

```json
{
  "env": "development",
  "port": 8080,
  "db": "root:123456@tcp(localhost:3306)/db_testing?charset=utf8&parseTime=True&charset=utf8mb4&collation=utf8mb4_unicode_ci"
}
```

```
$ ./server
``` 

#### View Swagger Docs

```
http://localhost:8080/docs/index.html
```

#### Do something else

1. Login

```azure
curl
--location 'localhost:8080/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
"Email": "user1@example.com",
"Password": "password1"
}'
```

2. Get user info

```azure
curl
--location 'localhost:8080/me' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXIxQGV4YW1wbGUuY29tIiwiZXhwIjoxNjk4MTI5NzY4LCJpZCI6MSwibmFtZSI6Ikp1c3R1cyIsIm9yaWdfaWF0IjoxNjk3NTI0OTY4fQ.0Rw-Lk2NcTq217X42ppGVIdd-yHLAy0SCff-TdsmYm0'
```

#### Run test

```azure
$ go
test./...
```

#### Worker

Project build a worker management to run background jobs (also call database queue)

1. Create a event

```sql
INSERT INTO `custodial_payment_address_actions` (`created_at`, `updated_at`, `deleted_at`, `entity_id`, `from`, `to`,
                                                 `amount`, `aasm_state`, `stage_status`, `error`, `err_count`,
                                                 `completed_at`)
VALUES (NULL, '2023-10-17 08:53:52.029', NULL, 1, 'userA', 'userB', '1000', 'submitted', 0, NULL, 0,
        '2023-10-17 08:53:52.029');
```

2. Result after run worker

```sql
INSERT INTO `custodial_payment_address_actions` (created_at`, `updated_at`, `deleted_at`, `entity_id`, `from`, `to`, `amount`, `aasm_state`, `stage_status`, `error`, `err_count`, `completed_at`)
VALUES
	(NULL, '2023-10-17 08:53:52.029', NULL, 1, 'userA', 'userB', '1000', 'succeeded', 0, NULL, 0, '2023-10-17 08:53:52.029');
```

Change here is `aasm_state` from `submitted` to `succeeded`

And with idea [AASM - State machines](https://github.com/aasm/aasm), we can add more state to handle more case






  