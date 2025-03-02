> Dockerfile has shared network for the same subnet, both should be on the same localhost

Requirements:
1) Go 1.24
2) Docker
3) Docker-compose
4) Make
5) jq
6) Ubuntu

# Default .env file:

- PORT=8080
- DB_URL=postgres://postgres:postgres@db:5432/wallet?sslmode=disable
- POSTGRES_USER=postgres
- POSTGRES_PASSWORD=postgres
- POSTGRES_DB=wallet
- EXCHANGER_GRPC_ADDR=gw-exchanger:50051
- JWT_SECRET=your_jwt_secret

> Makefile scripts for installations
> - make help - to see all make commands
> - make test - grpc testing, need to install grpcurl ```go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest```

> Full install and test, run the following command: ```wget https://raw.githubusercontent.com/CherepanovAndrey-git/WalletApp/ArchFix/scripts/install-wallet.sh && chmod +x install-wallet.sh && ./install-wallet.sh```

> After installation run ```cd ~/gw-wallet && make test-api```

### API
### 1. Регистрация пользователя

Метод: **POST**

URL: **/v1/register**

Заголовки:

**Content-Type: application/json**

Тело запроса:
```json
{
"username": "string",
"password": "string",
"email": "string"
}
```
Ответ:

> Успех: ```201 Created```

```json
{"message": "User registered successfully"}
```

> Ошибка: ```400 Bad Request```
```json
{"error": "Username or email already exists"}
```

### 2. Авторизация пользователя

Метод: **POST**

URL: **/v1/login**

Заголовки:

**Content-Type: application/json**

Тело запроса:
```json
{
"username": "string",
"password": "string"
}
```
Ответ:
> Успех: ```200 OK```
```json
("token": "JWT_TOKEN"}
```
> Ошибка: ```401 Unauthorized```
```json
{"error": "Invalid username or password"}
```

### 3. Создание кошелька

Метод: **POST**

URL: **/v1/create-wallet**

Заголовки:

**Authorization: Bearer JWT_TOKEN**

Ответ:

> Успех: ```201 Created```
```json
{
"message": "Wallet created successfully",
"balances": {
"USD": 0,
"RUB": 0,
"EUR": 0
}
}
```
> Ошибка: ```401 Unauthorized```
```json
{"error": "Unauthorized"}
```

### 4. Пополнение счета

Метод: **POST**

URL: **/v1/wallet/deposit**

Заголовки:

**Authorization: Bearer JWT_TOKEN**

**Content-Type: application/json**

Тело запроса:
```json
{
"amount": 100.00,
"currency": "USD" // (USD, RUB, EUR)
}
```
Ответ:

>Успех: 200 OK
```json
{
"message": "deposit successful",
"new_balance": {
"USD": 1000.50,
"RUB": 500.00,
"EUR": 300.00
}
}
```
> Ошибка: 400 Bad Request
```json
{"error": "Invalid amount or currency"}
```

### 5. Снятие средств

Метод: **POST**

URL: **/v1/wallet/withdraw**

Заголовки:

**Authorization: Bearer JWT_TOKEN**

**Content-Type: application/json**

Тело запроса:
```json
{
"amount": 100.00,
"currency": "USD" // (USD, RUB, EUR)
}
```
Ответ:

>Успех: 200 OK
```json
{
"message": "withdraw successful",
"new_balance": {
"USD": 900.50,
"RUB": 500.00,
"EUR": 300.00
}
}
```
> Ошибка: 400 Bad Request
```json
{"error": "Insufficient funds"}
```

### 6. Получение баланса

Метод: **GET**

URL: **/v1/balance**

Заголовки:

**Authorization: Bearer JWT_TOKEN**

Ответ:

>Успех: 200 OK
```json
{
"balances": {
"USD": 1000.50,
"RUB": 500.00,
"EUR": 300.00
}
}
```
>Ошибка: 401 Unauthorized
```json
{"error": "Unauthorized"}
```

### 7. Получение курсов валют

Метод: **GET**

URL: **/v1/exchange/rates**

Заголовки:

**Authorization: Bearer JWT_TOKEN**

Ответ:

>Успех: 200 OK
```json
{
"rates": {
"USD_RUB": 93.00,
"USD_EUR": 0.92,
"RUB_USD": 0.0108,
"RUB_EUR": 0.00985,
"EUR_USD": 1.087,
"EUR_RUB": 101.5
}
}
```

### 8. Обмен валюты

Метод: **POST**

URL: **/v1/exchange**

Заголовки:

**Authorization: Bearer JWT_TOKEN**

**Content-Type: application/json**

Тело запроса:
```json
{
"from_currency": "USD", // (USD, RUB, EUR)
"to_currency": "RUB",   // (USD, RUB, EUR)
"amount": 100.00
}
```
Ответ:
>Успех: 200 OK
```json
{
"message": "Exchange successful",
"exchanged_amount": "9300.00",
"new_balance": {
"USD": 900.00,
"RUB": 9300.00,
"EUR": 300.00
}
}
```
>Ошибка: 400 Bad Request
```json
{"error": "Invalid currency pair"}
```
>Ошибка: 400 Bad Request
```json
{"error":"Insufficient funds"}
```

### 9. Общие ошибки

Все эндпоинты возвращают следующие ошибки в случае проблем:

>400 Bad Request:
```json
{"error": "Invalid request payload"}
```
> 401 Unauthorized:
```json
{"error": "Unauthorized"}
```
> 500 Internal Server Error:
```json
{"error": "Internal server error"}
```

### Примеры запросов:


1. Регистрация нового пользователя:
```json
POST http://localhost:8080/v1/register
{
"username": "test_user",
"password": "test_password",
"email": "test@example.com"
}
```

2. Авторизация для получения токена:
```json
POST http://localhost:8080/v1/login
{
"username": "test_user",
"password": "test_password"
}
```
3. Создание кошелька:
```json
POST http://localhost:8080/v1/create-wallet
Headers: { "Authorization": "Bearer JWT_TOKEN" }
```
4. Пополнение счета:
```json
POST http://localhost:8080/v1/wallet/deposit
Headers: { "Authorization": "Bearer JWT_TOKEN" }
{
"amount": 1000,
"currency": "USD"
}
```
5. Проверка баланса:
```json
GET http://localhost:8080/v1/balance
Headers: { "Authorization": "Bearer JWT_TOKEN" }
```
6. Обмен валюты:
```json
POST http://localhost:8080/v1/exchange
Headers: { "Authorization": "Bearer JWT_TOKEN" }
{
"from_currency": "USD",
"to_currency": "RUB",
"amount": 100
}
```
