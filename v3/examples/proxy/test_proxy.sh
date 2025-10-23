#!/bin/bash

# Test proxy functionality

echo "Testing proxy configuration..."
echo ""

# Test direct access to the proxy
echo "1. Testing health endpoint via proxy:"
curl -X GET \
  -H "Origin: https://app-local.wails-awesome.io:3000" \
  "http://wails.localhost/health" 2>&1 | grep -E "(< HTTP|proxy|status)" || echo "No proxy response found"

echo ""
echo "2. Testing Wails runtime API:"
curl -X GET \
  -H "Origin: https://app-local.wails-awesome.io:3000" \
  -H "X-Wails-Window-ID: 1" \
  -v "http://wails.localhost/wails/runtime?object=9&method=0" 2>&1 | grep -E "(< HTTP|< X-|runtime)" || echo "No runtime response found"

echo ""
echo "Done testing proxy configuration."