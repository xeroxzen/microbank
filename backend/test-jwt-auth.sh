#!/bin/bash

echo "🔍 Testing JWT Authentication Flow"
echo "=================================="

# Test health checks
echo "1. Testing health checks..."
echo "   Client Service:"
curl -s http://localhost:8081/health | jq '.' 2>/dev/null || echo "Client service health check failed"
echo "   Banking Service:"
curl -s http://localhost:8080/health | jq '.' 2>/dev/null || echo "Banking service health check failed"

echo ""
echo "2. Testing protected endpoint without token..."
curl -s -w "HTTP Status: %{http_code}\n" http://localhost:8080/api/v1/account/balance

echo ""
echo "3. Testing protected endpoint with invalid token..."
curl -s -w "HTTP Status: %{http_code}\n" -H "Authorization: Bearer invalid_token" http://localhost:8080/api/v1/account/balance

echo ""
echo "4. To test with valid token:"
echo "   a) First get a token by calling the login endpoint"
echo "   b) Use: curl -H 'Authorization: Bearer YOUR_TOKEN' http://localhost:8080/api/v1/account/balance"
