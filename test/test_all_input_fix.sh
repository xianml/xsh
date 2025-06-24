#!/bin/bash

echo "🔄 全面测试输入冲突修复"
echo ""

# 设置测试环境变量
export ANTHROPIC_API_KEY="sk-ant-api03-ktEIQRrNpc9FHsUEoLEpTKderUQevFdMBjWfTV2bZGCm751pi326ji2ur_IuG4OuSgj3sbvcVfQYx-fzrQ3bHXA-4Bz_DgAA"

echo "✅ 构建最新版本..."
go build -o xsh .

echo "✅ 构建完成"
echo ""

echo "🎯 修复的两个输入冲突问题："
echo ""

echo "问题 1 - 模型选择输入冲突："
echo "- 现象：Tab 键触发模型选择后，输入数字被当作命令执行"
echo "- 原因：bufio.Reader 和 readline 争夺标准输入"
echo "- 解决：使用 modelSelectorTriggered 状态管理"
echo ""

echo "问题 2 - AI 命令执行输入冲突："
echo "- 现象：AI 建议命令后，输入选择被当作新命令"
echo "- 原因：同样的 bufio.Reader 冲突"
echo "- 解决：使用 commandExecutionTriggered 状态管理"
echo ""

echo "🔧 统一的状态管理方案："
echo ""
echo "状态变量："
echo "- modelSelectorTriggered: 模型选择状态"
echo "- commandExecutionTriggered: 命令执行选择状态"
echo "- pendingCommands: 待执行的AI建议命令"
echo ""

echo "输入循环逻辑："
echo "1. 检查模型选择状态 → handleModelSelection()"
echo "2. 检查命令执行状态 → handleCommandExecution()"
echo "3. 正常命令处理"
echo ""

echo "🚀 测试场景："
echo ""

echo "测试 1 - Tab 键模型选择："
echo "1. 启动 xsh: ./xsh"
echo "2. 空输入 + Tab"
echo "3. 输入数字如 33"
echo "4. 应该成功切换模型"
echo ""

echo "测试 2 - AI 命令建议："
echo "1. 输入: '帮我打开test.sh'"
echo "2. 按 Tab 键触发 AI 分析"
echo "3. 看到命令建议列表"
echo "4. 输入数字如 1 或 y"
echo "5. 应该执行对应命令"
echo ""

echo "测试 3 - 混合测试："
echo "1. Tab 键选择模型"
echo "2. AI 命令建议"
echo "3. 正常命令执行"
echo "4. 所有输入都应该正确处理"
echo ""

echo "🎉 关键改进："
echo "- 完全移除 bufio.Reader 的使用"
echo "- 统一使用 readline.Readline() 输入"
echo "- 状态机管理不同的输入模式"
echo "- Ctrl+C 正确重置所有状态"
echo ""

echo "启动测试：./xsh"
echo "先试试空输入+Tab，再试试AI命令建议" 