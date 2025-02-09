make install for full installation
make help for other commands
added load2.js for k6 testing
added TestingForAPI.postman_collection.json for postman testing 

POST http://{{baseUrl}}/v1/users
{
    "name": "test",
}
response:
{
    "id": "8ccc2085-64c1-4545-b93e-2989bc048317",
    "created_at": "2025-01-23T13:38:30.795386Z",
    "updated_at": "2025-01-23T13:38:30.795386Z",
    "name": "koi",
    "api_key": "fe0d0f4ad4462ed2dabfe4db97eda10e5d84c093e53c953e4dd3ac4274d4aab1"
}

POST http://{{baseUrl}}/v1/wallets
Auth: ApiKey fe0d0f4ad4462ed2dabfe4db97eda10e5d84c093e53c953e4dd3ac4274d4aab1
response:
{
    "walletId": "f6433faf-3c25-460b-803a-1fc02b1202a1"
}
GET http://{{baseUrl}}/v1/wallets/f6433faf-3c25-460b-803a-1fc02b1202a1/balance
response:
{
    "balance": 0
}
POST http://{{baseUrl}}/v1/wallet
Auth: ApiKey fe0d0f4ad4462ed2dabfe4db97eda10e5d84c093e53c953e4dd3ac4274d4aab1
JSON body
{
  "walletId": "f6433faf-3c25-460b-803a-1fc02b1202a1",
  "operationType": "DEPOSIT",
  "amount": 1000
}
response:
{
  "balance": "1000",
  "walletId": "f6433faf-3c25-460b-803a-1fc02b1202a1"
}

POST http://{{baseUrl}}/v1/wallet
Auth: ApiKey fe0d0f4ad4462ed2dabfe4db97eda10e5d84c093e53c953e4dd3ac4274d4aab1
JSON body
{
  "walletId": "f6433faf-3c25-460b-803a-1fc02b1202a1",
  "operationType": "WITHDRAW",
  "amount": 500
}
response:
{
  "balance": "500",
  "walletId": "f6433faf-3c25-460b-803a-1fc02b1202a1"
}

