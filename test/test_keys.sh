#!/bin/bash

echo "🔥 Testing xsh keyboard events..."

# 设置测试环境变量
export OPENAI_API_KEY="test-key"

echo ""
echo "📋 Testing keyboard event functionality:"
echo ""

echo "1. Normal command test:"
echo -e "ll\nexit" | ./xsh

echo ""
echo "2. AI command test:"
echo -e "ai list files\nexit" | ./xsh

echo ""
echo "3. Model switching test:"
echo -e "models\nexit" | ./xsh

echo ""
echo "🎯 Key functionality summary:"
echo "- Tab key: Triggers AI analysis of current input (interactive only)"
echo "- 'ai <question>': Alternative way to get AI help"
echo "- 'models': Change AI model"
echo "- 'll' and 'la': Work with built-in aliases"
echo "- Ctrl+C: Exit"
echo ""
echo "To test Tab key functionality, run: ./xsh"
echo "Then type something like 'find large files' and press Tab" 