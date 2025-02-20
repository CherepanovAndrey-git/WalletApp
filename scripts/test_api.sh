#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080/v1"


DEPOSIT_AMOUNT=1000
WITHDRAWAL_AMOUNT=500

# Colored output
print_step() {
    echo -e "${BLUE}=== $1 ===${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}


print_step "Registering new user"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/register" \
    -H "Content-Type: application/json" \
    -d '{
        "username": "testuser1234",
        "email": "test123@example.com",
        "password": "password152352358910463"
    }')
echo $REGISTER_RESPONSE | jq
print_success "Registration completed"


print_step "Logging in"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{
        "username": "testuser1234",
        "password": "password152352358910463"
    }')
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    print_error "Failed to get token"
    echo "Login response: $LOGIN_RESPONSE"
    exit 1
fi

print_success "Successfully logged in"
echo "Token: ${TOKEN:0:20}..."


print_step "Creating wallet"
WALLET_RESPONSE=$(curl -s -X POST "$BASE_URL/create-wallet" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")
echo $WALLET_RESPONSE | jq
print_success "Wallet created with initial balance"


print_step "Depositing $DEPOSIT_AMOUNT USD"
DEPOSIT_RESPONSE=$(curl -s -X POST "$BASE_URL/wallet/deposit" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"amount\": $DEPOSIT_AMOUNT,
        \"currency\": \"USD\"
    }")
echo $DEPOSIT_RESPONSE | jq
print_success "Deposit of $DEPOSIT_AMOUNT USD completed"


print_step "Withdrawing $WITHDRAWAL_AMOUNT USD"
WITHDRAW_RESPONSE=$(curl -s -X POST "$BASE_URL/wallet/withdraw" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"amount\": $WITHDRAWAL_AMOUNT,
        \"currency\": \"USD\"
    }")
echo $WITHDRAW_RESPONSE | jq
print_success "Withdrawal of $WITHDRAWAL_AMOUNT USD completed"


print_step "Final balance"
FINAL_BALANCE_RESPONSE=$(curl -s -X GET "$BASE_URL/balance" \
    -H "Authorization: Bearer $TOKEN")
echo $FINAL_BALANCE_RESPONSE | jq
print_success "Test sequence completed"


print_step "Transaction Summary"
echo -e "Initial Deposit: ${GREEN}+$DEPOSIT_AMOUNT USD${NC}"
echo -e "Withdrawal:     ${RED}-$WITHDRAWAL_AMOUNT USD${NC}"
echo -e "Final Balance:  ${GREEN}$(echo $FINAL_BALANCE_RESPONSE | jq -r '.balance') USD${NC}" 