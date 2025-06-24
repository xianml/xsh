#!/bin/bash

echo "🔄 测试异步 Tab 键模型选择功能"
echo ""

# 设置测试环境变量
export ANTHROPIC_API_KEY="sk-ant-api03-ktEIQRrNpc9FHsUEoLEpTKderUQevFdMBjWfTV2bZGCm751pi326ji2ur_IuG4OuSgj3sbvcVfQYx-fzrQ3bHXA-4Bz_DgAA"

echo "✅ 构建最新版本..."
go build -o xsh .

echo "✅ 构建完成"
echo ""

echo "🎯 异步 Tab 键功能测试："
echo ""
echo "新的实现方式："
echo "1. Tab 键在补全器中设置状态标记"
echo "2. 使用 goroutine 异步处理模型选择"
echo "3. 用独立的 bufio.Reader 避免与 readline 冲突"
echo ""

echo "测试步骤："
echo ""

echo "测试 1 - 空输入 Tab 键："
echo "1. 启动 xsh: ./xsh"
echo "2. 不输入任何内容，直接按 Tab 键"
echo "3. 应该显示模型选择界面"
echo "4. 数字输入应该可以正常显示和响应"
echo ""

echo "测试 2 - 有输入 Tab 键："
echo "1. 输入任何文本，例如: 'list files'"
echo "2. 按 Tab 键"
echo "3. 应该触发 AI 分析"
echo ""

echo "测试 3 - 普通命令："
echo "1. 输入: 'm'"
echo "2. 按回车"
echo "3. 应该显示模型选择界面"
echo ""

echo "🔧 技术细节："
echo "- 使用状态标记 modelSelectorTriggered"
echo "- 100ms 延迟让补全器完成"
echo "- bufio.Reader 独立输入处理"
echo "- 避免 readline 输入冲突"
echo ""

echo "🚀 这种方式应该彻底解决输入冲突问题！"
echo ""

echo "启动测试：./xsh"
echo "然后尝试空输入 + Tab 键" 