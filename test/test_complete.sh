#!/bin/bash

echo "🔥 xsh 完整功能测试 - 最新版本"
echo ""

# 设置测试环境变量
export ANTHROPIC_API_KEY="sk-ant-api03-ktEIQRrNpc9FHsUEoLEpTKderUQevFdMBjWfTV2bZGCm751pi326ji2ur_IuG4OuSgj3sbvcVfQYx-fzrQ3bHXA-4Bz_DgAA"

echo "✅ 环境变量设置："
echo "   ANTHROPIC_API_KEY: ${ANTHROPIC_API_KEY:0:20}..."
echo ""

# 重新构建
go build -o xsh .

echo "✅ xsh 重新构建完成"
echo ""

echo "🎯 功能测试列表："
echo ""

echo "1. 基本命令测试："
echo "   输入: ll"
echo "   预期: 执行 ls -lh 命令"
echo ""

echo "2. 模型选择测试 (新功能 - 推荐)："
echo "   操作: 直接按 Tab 键 (不输入任何内容)"
echo "   预期: 立即显示模型选择菜单"
echo ""

echo "3. AI 分析测试 (Tab 键)："
echo "   输入: find large files"
echo "   按: Tab 键"
echo "   预期: 立即显示 AI 分析结果"
echo ""

echo "4. Execute command 交互测试："
echo "   输入: help me list files"
echo "   按: Tab 键"
echo "   预期: 显示 AI 建议，然后可以选择执行"
echo "   注意: 修复了卡住的 bug"
echo ""

echo "5. Ctrl+C 测试："
echo "   输入任意文本，然后按 Ctrl+C"
echo "   预期: 取消当前行，不退出程序"
echo ""

echo "6. AI 命令测试："
echo "   输入: ai check disk space"
echo "   按: Enter"
echo "   预期: 显示 AI 建议的命令"
echo ""

echo "现在启动 xsh 进行测试："
echo "./xsh"
echo ""
echo "💡 重要更新:"
echo "   ✅ Tab 键智能行为："
echo "      - 空输入时：选择模型"
echo "      - 有输入时：AI 分析"
echo "   ✅ 修复了 Execute command 卡住的 bug"
echo "   ✅ 移除了 models 命令"
echo "   ✅ 简化了键盘快捷键" 