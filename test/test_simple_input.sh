#!/bin/bash

echo "🔄 简化输入测试"
echo ""

# 设置测试环境变量
export ANTHROPIC_API_KEY="sk-ant-api03-ktEIQRrNpc9FHsUEoLEpTKderUQevFdMBjWfTV2bZGCm751pi326ji2ur_IuG4OuSgj3sbvcVfQYx-fzrQ3bHXA-4Bz_DgAA"

echo "✅ 构建最新版本..."
go build -o xsh .

echo "✅ 构建完成"
echo ""

echo "🎯 简化的模型选择方式："
echo ""
echo "现在有两种方式选择模型："
echo ""
echo "方式 1 - 直接命令："
echo "1. 启动 xsh: ./xsh"
echo "2. 输入: m"
echo "3. 按回车"
echo "4. 应该显示所有可用模型"
echo "5. 输入数字选择模型"
echo ""

echo "方式 2 - 完整命令："
echo "1. 启动 xsh: ./xsh"
echo "2. 输入: models"
echo "3. 按回车"
echo "4. 应该显示所有可用模型"
echo "5. 输入数字选择模型"
echo ""

echo "Tab键功能："
echo "- 空输入 + Tab: 显示提示信息"
echo "- 有输入 + Tab: AI 分析"
echo ""

echo "🚀 这种方式避免了复杂的输入冲突问题"
echo ""

echo "启动测试：./xsh"
echo "然后输入 'm' 并按回车" 