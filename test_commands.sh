#!/bin/bash
# =============================================================================
# cliai 命令全面测试脚本
# 安全原则：
#   1. 所有操作在 cliai/test_workspace/ 目录下进行
#   2. 危险命令使用绝对路径和精确目标
#   3. 不会影响系统或其他文件
# =============================================================================

set -e  # 遇到错误立即退出

CLI="./cliai-go/cliai"
WORKSPACE="./test_workspace"
RESULTS="$WORKSPACE/test_results"
PASS=0
FAIL=0
TOTAL=0

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 初始化
mkdir -p "$RESULTS"
echo "========================================" > "$RESULTS/full_test.log"
echo "cliai 命令全面测试 - $(date)" >> "$RESULTS/full_test.log"
echo "========================================" >> "$RESULTS/full_test.log"

# 测试函数
test_cmd() {
    local cmd="$1"
    local args="$2"
    local expected_pattern="$3"
    local category="$4"
    local test_name="$5"
    
    TOTAL=$((TOTAL + 1))
    local full_cmd="$cmd $args"
    local full_path_cmd="$CLI $full_cmd"
    
    echo -n "[$category] $test_name: $cmd $args ... "
    
    # 捕获输出和退出码
    set +e
    output=$(cd "$WORKSPACE" && eval "$full_path_cmd" 2>&1)
    exit_code=$?
    set -e
    
    # 检查输出是否符合预期
    if echo "$output" | grep -q "$expected_pattern"; then
        echo -e "${GREEN}PASS${NC}"
        echo "[PASS] $test_name: $full_cmd" >> "$RESULTS/full_test.log"
        PASS=$((PASS + 1))
    else
        echo -e "${RED}FAIL${NC}"
        echo "[FAIL] $test_name: $full_cmd" >> "$RESULTS/full_test.log"
        echo "  Expected: $expected_pattern" >> "$RESULTS/full_test.log"
        echo "  Got: $output" >> "$RESULTS/full_test.log"
        FAIL=$((FAIL + 1))
    fi
}

test_cmd_exists() {
    local cmd="$1"
    local category="$2"
    
    TOTAL=$((TOTAL + 1))
    echo -n "[$category] $cmd 是否注册: "
    
    output=$(eval "$CLI commands 2>&1")
    
    if echo "$output" | grep -qw "$cmd"; then
        echo -e "${GREEN}PASS${NC}"
        PASS=$((PASS + 1))
    else
        echo -e "${RED}FAIL${NC}"
        FAIL=$((FAIL + 1))
    fi
}

echo "========================================"
echo "cliai 命令全面测试"
echo "========================================"
echo ""

# =============================================================================
# 第一阶段：命令注册测试 (验证所有命令都已正确注册)
# =============================================================================
echo -e "${YELLOW}[阶段 1] 命令注册测试${NC}"
echo ""

categories=("network" "text" "file" "system" "encoding" "shell")

for cat in "${categories[@]}"; do
    echo "--- $cat ---"
done

# P0 命令
for cmd in curl jq grep cat ls head tail ps df free base64; do
    test_cmd_exists "$cmd" "P0"
done

# P1 命令
for cmd in wget ping sed awk cut sort wc find du top md5sum sha256sum; do
    test_cmd_exists "$cmd" "P1"
done

# P2 网络命令
for cmd in ip ss netstat ifconfig route arp nslookup dig host traceroute tracepath nc telnet nmap tcpdump; do
    test_cmd_exists "$cmd" "P2"
done

# P3 文件命令
for cmd in ln chmod chown tar zip unzip rsync tree; do
    test_cmd_exists "$cmd" "P3"
done

# P4 文本命令
for cmd in less uniq diff paste tr xargs tee; do
    test_cmd_exists "$cmd" "P4"
done

# P5 系统命令
for cmd in kill pkill htop systemctl mount umount lsof strace; do
    test_cmd_exists "$cmd" "P5"
done

# Shell 内建
for cmd in cd pwd echo mkdir touch rm cp mv clear exit which history; do
    test_cmd_exists "$cmd" "shell"
done

echo ""

# =============================================================================
# 第二阶段：功能测试 (使用 test_workspace 内的测试数据)
# =============================================================================
echo -e "${YELLOW}[阶段 2] 功能测试${NC}"
echo ""

# --- 编码命令 ---
echo -e "${YELLOW}[编码命令]${NC}"

test_cmd "base64" "-d <<< aGVsbG8gd29ybGQ=" "hello world" "encoding" "base64 解码"

test_cmd "md5sum" "text/sample.txt" "5d..." "encoding" "md5sum 计算"

test_cmd "sha256sum" "text/sample.txt" "sha256" "encoding" "sha256sum 计算"

# --- 文本命令 ---
echo ""
echo -e "${YELLOW}[文本命令]${NC}"

test_cmd "cat" "text/sample.txt" "Hello World" "text" "cat 查看文件"

test_cmd "head" "-n 2 text/sample.txt" "Hello World" "text" "head 查看头部"

test_cmd "tail" "-n 2 text/sample.txt" "Last line" "text" "tail 查看尾部"

test_cmd "grep" "line" "text/sample.txt" "line1" "text" "grep 搜索"

test_cmd "wc" "-l text/sample.txt" "6" "text" "wc 统计行数"

test_cmd "wc" "-w text/sample.txt" "18" "text" "wc 统计词数"

test_cmd "sort" "text/file1.txt" "cat" "text" "sort 排序"

test_cmd "uniq" "text/duplicate.txt" "aaa" "text" "uniq 去重"

test_cmd "cut" "-d',' -f2 text/sample.csv" "apple" "text" "cut 截取"

test_cmd "jq" "." "text/sample.json" "name" "text" "jq 解析 JSON"

# --- 文件命令 ---
echo ""
echo -e "${YELLOW}[文件命令]${NC}"

test_cmd "ls" "$WORKSPACE/text/" "sample.txt" "file" "ls 列出目录"

test_cmd "find" "$WORKSPACE -name '*.txt'" "sample.txt" "file" "find 查找文件"

test_cmd "du" "$WORKSPACE/text/" "sample.txt" "file" "du 磁盘使用"

# --- Shell 内建 ---
echo ""
echo -e "${YELLOW}[Shell 内建命令]${NC}"

test_cmd "echo" "hello_test_123" "hello_test_123" "shell" "echo 输出"

test_cmd "pwd" "" "$WORKSPACE" "shell" "pwd 当前目录"

test_cmd "which" "ls" "/" "shell" "which 查找命令"

# --- 网络命令 (只测试本地/安全目标) ---
echo ""
echo -e "${YELLOW}[网络命令 - 本地测试${NC}"

# 测试 ss, netstat, ip 等能返回本地信息
test_cmd "ip" "addr show lo" "inet" "network" "ip addr"

test_cmd "ss" "-tuln" "Local" "network" "ss socket"

test_cmd "ping" "-c 1 127.0.0.1" "1 packets transmitted" "network" "ping localhost"

test_cmd "nslookup" "localhost" "Server:" "network" "nslookup localhost"

# --- 系统命令 (只测试信息展示) ---
echo ""
echo -e "${YELLOW}[系统命令 - 只读测试${NC}"

test_cmd "ps" "aux" "PID" "system" "ps 进程列表"

test_cmd "df" "-h" "Filesystem" "system" "df 磁盘"

test_cmd "free" "-h" "Mem:" "system" "free 内存"

test_cmd "top" "-b -n 1" "PID" "system" "top 一次"

# =============================================================================
# 第三阶段：组合测试
# =============================================================================
echo ""
echo -e "${YELLOW}[阶段 3] 组合命令测试${NC}"
echo ""

echo "[组合] cat + grep + wc"
output=$(cd "$WORKSPACE" && $CLI "cat text/sample.txt | grep line | wc -l" 2>&1 || true)
if echo "$output" | grep -qE "^[0-9]+"; then
    echo -e "[组合] ${GREEN}PASS${NC}"
    PASS=$((PASS + 1))
else
    echo -e "[组合] ${RED}FAIL${NC}"
    FAIL=$((FAIL + 1))
fi
TOTAL=$((TOTAL + 1))

echo "[组合] echo + tee"
output=$(cd "$WORKSPACE" && $CLI "echo test_tee | tee text/tee_output.txt" 2>&1 || true)
if echo "$output" | grep -q "test_tee"; then
    echo -e "[组合] ${GREEN}PASS${NC}"
    PASS=$((PASS + 1))
else
    echo -e "[组合] ${RED}FAIL${NC}"
    FAIL=$((FAIL + 1))
fi
TOTAL=$((TOTAL + 1))

# =============================================================================
# 第四阶段：危险命令安全测试
# =============================================================================
echo ""
echo -e "${YELLOW}[阶段 4] 危险命令安全测试${NC}"
echo ""

# rm 只测试 help 或者 dry-run，不实际删除任何东西
echo "[安全] rm 命令存在性检查"
test_cmd_exists "rm" "dangerous"

# chmod 只测试自己的测试文件
echo "[安全] chmod/chown 只作用于测试目录"
test_cmd "chmod" "644 text/sample.txt" "text" "shell" "chmod 权限"

echo ""
echo "========================================"
echo "测试完成"
echo "========================================"
echo "总计: $TOTAL"
echo -e "通过: ${GREEN}$PASS${NC}"
echo -e "失败: ${RED}$FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}有测试失败，详见 $RESULTS/full_test.log${NC}"
    exit 1
fi
