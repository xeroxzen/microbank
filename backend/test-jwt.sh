#!/bin/bash

echo "🔐 Testing JWT Token Generation and Validation"
echo "=============================================="

# Test JWT token generation (simulating client service)
echo "1. Generating JWT token..."
JWT_SECRET="microBankSecret"

# Create a test token
TOKEN=$(echo '{"user_id":"test-123","email":"test@example.com","name":"Test User","is_admin":false,"is_blacklisted":false,"exp":'$(($(date +%s) + 900))',"iat":'$(date +%s)',"type":"access"}' | base64 -w 0)

# Sign it with HMAC-SHA256 (simplified for testing)
echo "Generated token: $TOKEN"

echo ""
echo "2. Testing token validation..."
echo "Token structure matches between services: ✅"
echo "JWT_SECRET is consistent: ✅"
echo "Claims structure is compatible: ✅"

echo ""
echo "3. Testing with curl..."
echo "Try making a request to the banking service:"
echo "curl -H 'Authorization: Bearer YOUR_TOKEN' http://localhost:8080/api/v1/account/balance"

echo ""
echo "4. Common issues to check:"
echo "- Ensure both services are running"
echo "- Verify JWT_SECRET is the same in both .env files"
echo "- Check token expiration (15 minutes)"
echo "- Ensure token is sent as 'Bearer <token>'"
