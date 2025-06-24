#!/bin/bash

echo "🚀 测试实时模型获取功能"
echo ""

# 设置测试环境变量
export ANTHROPIC_API_KEY="sk-ant-api03-ktEIQRrNpc9FHsUEoLEpTKderUQevFdMBjWfTV2bZGCm751pi326ji2ur_IuG4OuSgj3sbvcVfQYx-fzrQ3bHXA-4Bz_DgAA"

echo "✅ 环境变量设置："
echo "   ANTHROPIC_API_KEY: ${ANTHROPIC_API_KEY:0:20}..."
echo ""

go build -o xsh .

echo "✅ 构建完成"
echo ""

echo "🎯 新功能测试："
echo ""
echo "1. 实时模型获取："
echo "   现在应该显示所有可用的 Claude 模型："
echo "   - claude-3-5-sonnet-20241022"
echo "   - claude-3-5-sonnet-20240620"
echo "   - claude-3-sonnet-20240229"
echo "   - claude-3-opus-20240229"
echo "   - claude-3-haiku-20240307"
echo "   - claude-2.1"
echo "   - claude-2.0"
echo "   - claude-instant-1.2"
echo ""

echo "2. 模型选择修复："
echo "   现在应该能够输入数字选择模型"
echo ""

echo "启动 xsh，直接按 Tab 键测试："
echo "./xsh"
echo ""

echo "📝 测试步骤："
echo "1. 启动 xsh"
echo "2. 直接按 Tab 键"
echo "3. 查看是否显示所有 Claude 模型"
echo "4. 输入数字 (如 1) 选择模型"
echo "5. 检查是否成功切换" 