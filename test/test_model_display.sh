#!/bin/bash

echo "🔧 测试模型显示修复"
echo ""

# 设置测试环境变量
export ANTHROPIC_API_KEY="sk-ant-api03-ktEIQRrNpc9FHsUEoLEpTKderUQevFdMBjWfTV2bZGCm751pi326ji2ur_IuG4OuSgj3sbvcVfQYx-fzrQ3bHXA-4Bz_DgAA"
export OPENAI_API_KEY="test-key"
export GOOGLE_API_KEY="test-key"

echo "✅ 设置了多个 API 密钥用于测试"
echo ""

go build -o xsh .

echo "✅ 构建完成"
echo ""

echo "🎯 现在应该显示具体的模型名称："
echo ""
echo "预期显示类似："
echo "1. claude-3-sonnet-20240229 (anthropic) (current)"
echo "2. gpt-3.5-turbo (openai)"
echo "3. gemini-pro (google)"
echo ""
echo "而不是："
echo "1. claude (current)"
echo "2. openai"
echo "3. gemini"
echo ""

echo "启动 xsh 并按 Tab 键测试："
echo "./xsh" 