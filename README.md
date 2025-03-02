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