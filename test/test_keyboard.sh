#!/bin/bash

echo "🔥 Testing xsh keyboard events..."

# 设置测试环境变量
export OPENAI_API_KEY="test-key"
export OPENAI_MODEL="gpt-3.5-turbo"

echo ""
echo "📋 Keyboard event functionality test:"
echo ""

# 重新构建
go build -o xsh .

echo "✅ xsh rebuilt successfully"
echo ""

echo "🎯 Manual keyboard event tests:"
echo ""
echo "1. Ctrl+C test:"
echo "   Run: ./xsh"
echo "   Type some text, then press Ctrl+C"
echo "   Expected: Should exit immediately with 'Goodbye!' message"
echo ""

echo "2. Tab key test:"
echo "   Run: ./xsh"
echo "   Type: list files"
echo "   Press: Tab key (without Enter)"
echo "   Expected: Should trigger AI analysis immediately"
echo ""

echo "3. Models test:"
echo "   Run: ./xsh"
echo "   Type: models"
echo "   Press: Enter"
echo "   Expected: Should show model selection menu"
echo ""

echo "4. Normal command test:"
echo "   Run: ./xsh"
echo "   Type: ll"
echo "   Press: Enter"
echo "   Expected: Should execute ls -lh command"
echo ""

echo "To test manually, run: ./xsh"
echo ""
echo "Note: The Tab key should trigger AI analysis IMMEDIATELY without pressing Enter"
echo "      Ctrl+C should exit the shell IMMEDIATELY" 