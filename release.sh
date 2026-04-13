#!/bin/bash

# 发布 common 套件的脚本
# 用法: ./release.sh v1.2.1 "Add new features"

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查参数
if [ $# -lt 1 ]; then
    echo -e "${RED}❌ 用法: ./release.sh <版本号> [提交信息]${NC}"
    echo ""
    echo "例如:"
    echo "  ./release.sh v1.2.0"
    echo "  ./release.sh v1.2.0 'Add error handling utils'"
    exit 1
fi

VERSION=$1
COMMIT_MSG="${2:-Update common package}"

# 检查版本格式
if ! [[ $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}❌ 版本号格式错误，应该是 v X.Y.Z 格式${NC}"
    exit 1
fi

# 检查是否在 Git 仓库里
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}❌ 错误: 不在 Git 仓库里${NC}"
    exit 1
fi

# 检查是否有未提交的改动
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}❌ 错误: 有未提交的改动，请先 commit${NC}"
    echo ""
    echo "待提交的改动:"
    git diff --name-only
    exit 1
fi

# 检查该版本是否已存在
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    echo -e "${RED}❌ 错误: 版本 $VERSION 已经存在${NC}"
    exit 1
fi

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}发布 common 套件${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}版本: ${GREEN}$VERSION${NC}"
echo -e "${BLUE}提交信息: ${GREEN}$COMMIT_MSG${NC}"
echo ""

read -p "$(echo -e ${YELLOW}确认发布？${NC} [y/N] )" -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${RED}已取消${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}1. 提交代码...${NC}"
git add -A
git commit -m "$COMMIT_MSG" --allow-empty

echo -e "${BLUE}2. 推送到 main...${NC}"
git push origin main

echo -e "${BLUE}3. 打 tag: $VERSION...${NC}"
git tag -a "$VERSION" -m "Release $VERSION"

echo -e "${BLUE}4. 推送 tag 到远程...${NC}"
git push origin "$VERSION"

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ 发布成功！${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "其他项目现在可以用以下命令更新:"
echo -e "${BLUE}go get github.com/FatSong207/common@$VERSION${NC}"
echo ""
