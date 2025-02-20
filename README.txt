make install for full installation
make help for other commands
added load2.js for k6 testing
scripts/test_api.sh for api testing
swagger ui: http://localhost:8080/swagger
.env file for environment variables, default values:
PORT=8080
DB_URL=postgres://postgres:postgres@db:5432/wallet?sslmode=disable
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=wallet
JWT_SECRET=<your_secret>

or create your own .env