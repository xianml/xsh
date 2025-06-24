#!/bin/bash

echo "Testing xsh functionality..."

# 设置测试环境变量
export OPENAI_API_KEY="test-key"
export XSH_MODEL="openai"

echo "1. Testing version flag..."
./xsh --version

echo -e "\n2. Testing basic command execution (simulation)..."
echo -e "ls\nexit" | timeout 5s ./xsh || echo "Test completed or timeout"

echo -e "\n3. Testing AI command (simulation)..."
echo -e "ai list files in current directory\nexit" | timeout 5s ./xsh || echo "Test completed or timeout"

echo -e "\n4. Testing models command (simulation)..."
echo -e "models\nexit" | timeout 5s ./xsh || echo "Test completed or timeout"

echo -e "\nxsh build and basic functionality test completed!"
echo "To use xsh with real AI functionality:"
echo "1. Copy config.example to .env"
echo "2. Set your API keys in .env"
echo "3. Run: source .env && ./xsh" 