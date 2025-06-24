#!/bin/bash

echo "🚀 xsh - AI Powered Shell Demo"
echo "=================================="

# 设置测试环境变量
export OPENAI_API_KEY="demo-key-for-testing"

echo ""
echo "✅ Building xsh..."
make build > /dev/null 2>&1

echo ""
echo "📋 Testing basic functionality:"
echo ""

# 测试基本命令
echo "1. Testing basic commands (ll, la, ls):"
echo -e "ll\nla\nls" | ./xsh | head -20

echo ""
echo "2. Testing pwd and other commands:"
echo -e "pwd\nwhoami\nexit" | ./xsh

echo ""
echo "3. Testing version flag:"
./xsh --version

echo ""
echo "🎯 Demo completed!"
echo ""
echo "To start using xsh:"
echo "1. Configure your API keys in .env file"
echo "2. Run: source .env && ./xsh"
echo "3. Try commands like:"
echo "   - ll                    # List files"
echo "   - ai list large files   # AI assistance"  
echo "   - models               # Switch AI models"
echo "   - exit                 # Quit" 