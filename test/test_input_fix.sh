#!/bin/bash

echo "🔧 测试输入修复"
echo ""

# 设置测试环境变量
export ANTHROPIC_API_KEY="sk-ant-api03-ktEIQRrNpc9FHsUEoLEpTKderUQevFdMBjWfTV2bZGCm751pi326ji2ur_IuG4OuSgj3sbvcVfQYx-fzrQ3bHXA-4Bz_DgAA"

echo "✅ 构建最新版本..."
go build -o xsh .

echo "✅ 构建完成"
echo ""

echo "🎯 输入修复测试："
echo ""
echo "修复的问题："
echo "1. ❌ 之前：使用 fmt.Scanln() - 与 readline 冲突"
echo "2. ✅ 现在：使用 bufio.Reader - 更可靠的输入"
echo "3. ✅ 支持多位数字选择（不仅限于 1-9）"
echo ""

echo "📝 测试步骤："
echo "1. 启动 xsh: ./xsh"
echo "2. 直接按 Tab 键触发模型选择"
echo "3. 应该显示所有可用的 Claude 模型"
echo "4. 输入数字（如：1, 2, 3...）然后按回车"
echo "5. 应该成功选择并切换模型"
echo ""

echo "🚨 如果仍然无法输入："
echo "   - 确保按回车键确认输入"
echo "   - 尝试输入不同的数字"
echo "   - 检查终端是否正常响应"
echo ""

echo "启动测试：./xsh" 