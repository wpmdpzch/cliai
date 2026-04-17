#!/bin/bash
# =============================================================================
# cliai 命令全面测试脚本
# 安全原则：
#   1. 所有操作在 cliai/test_workspace/ 目录下进行
#   2. 危险命令使用绝对路径和精确目标
#   3. 不会影响系统或其他文件
#   4. 禁止使用会阻塞的命令（如无 timeout 的 ping）
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

# 测试函数 - 只检查命令是否注册
test_cmd_exists() {
    local cmd="$1"
    local category="$2"
    
    TOTAL=$((TOTAL + 1))
    echo -n "[$category] $cmd 是否注册: "
    
    output=$(eval "$CLI commands 2>&1")
    
    if echo "$output" | grep -qw "$cmd"; then
        echo -e "${GREEN}PASS${NC}"
        echo "[PASS] $cmd 注册成功" >> "$RESULTS/full_test.log"
        PASS=$((PASS + 1))
    else
        echo -e "${RED}FAIL${NC}"
        echo "[FAIL] $cmd 注册失败" >> "$RESULTS/full_test.log"
        FAIL=$((FAIL + 1))
    fi
}

# 功能测试函数
test_func() {
    local desc="$1"
    local expected="$2"
    shift 2
    local cmd="$@"
    
    TOTAL=$((TOTAL + 1))
    echo -n "[功能] $desc: $cmd ... "
    
    output=$(eval "$cmd" 2>&1)
    
    if echo "$output" | grep -q "$expected"; then
        echo -e "${GREEN}PASS${NC}"
        PASS=$((PASS + 1))
    else
        echo -e "${RED}FAIL${NC}"
        echo "[FAIL] $desc: $cmd" >> "$RESULTS/full_test.log"
        echo "  Expected: $expected" >> "$RESULTS/full_test.log"
        echo "  Got: $output" >> "$RESULTS/full_test.log"
        FAIL=$((FAIL + 1))
    fi
}

echo "========================================"
echo "cliai 命令全面测试"
echo "========================================"
echo ""

# =============================================================================
# 第一阶段：命令注册测试
# =============================================================================
echo -e "${YELLOW}[阶段 1] 命令注册测试${NC}"
echo ""

# P0 命令 (11个)
for cmd in curl jq grep cat ls head tail ps df free base64; do
    test_cmd_exists "$cmd" "P0"
done

# P1 命令 (12个)
for cmd in wget ping sed awk cut sort wc find du top md5sum sha256sum; do
    test_cmd_exists "$cmd" "P1"
done

# P2 网络命令 (15个)
for cmd in ip ss netstat ifconfig route arp nslookup dig host traceroute tracepath nc telnet nmap tcpdump; do
    test_cmd_exists "$cmd" "P2"
done

# P3 文件命令 (8个)
for cmd in ln chmod chown tar zip unzip rsync tree; do
    test_cmd_exists "$cmd" "P3"
done

# P4 文本命令 (7个)
for cmd in less uniq diff paste tr xargs tee; do
    test_cmd_exists "$cmd" "P4"
done

# P5 系统命令 (8个)
for cmd in kill pkill htop systemctl mount umount lsof strace; do
    test_cmd_exists "$cmd" "P5"
done

# Shell 内建 (12个)
for cmd in cd pwd echo mkdir touch rm cp mv clear exit which history; do
    test_cmd_exists "$cmd" "shell"
done

echo ""

# =============================================================================
# 第二阶段：功能测试 (使用隔离测试数据)
# =============================================================================
echo -e "${YELLOW}[阶段 2] 功能测试${NC}"
echo ""

# 文本命令
test_func "cat 查看文件" "Hello World" $CLI cat test_workspace/text/sample.txt
test_func "head 查看头部" "Hello World" $CLI head -n 2 test_workspace/text/sample.txt
test_func "tail 查看尾部" "Last line" $CLI tail -n 2 test_workspace/text/sample.txt
test_func "grep 搜索" "line" $CLI grep line test_workspace/text/sample.txt
test_func "wc 统计行数" "6" $CLI wc -l test_workspace/text/sample.txt
test_func "sort 排序" "bird" $CLI sort test_workspace/text/file1.txt
test_func "uniq 去重" "aaa" $CLI uniq test_workspace/text/duplicate.txt
test_func "cut 截取" "apple" $CLI cut -d',' -f2 test_workspace/text/sample.csv
test_func "jq 解析" "Alice" $CLI jq '.name' test_workspace/text/sample.json

# 文件命令
test_func "ls 列出目录" "sample.txt" $CLI ls test_workspace/text/
test_func "find 查找文件" "sample.txt" $CLI find test_workspace -name '*.txt'
test_func "du 磁盘使用" "sample.txt" $CLI du test_workspace/text/

# Shell 内建
test_func "echo 输出" "hello_test" $CLI echo hello_test
test_func "pwd 当前目录" "cliai" $CLI pwd
test_func "which 查找" "/bin/ls" $CLI which ls

# 网络命令 (本地)
test_func "ip addr" "inet" $CLI ip addr show lo
test_func "ss socket" "Local" $CLI ss -tuln
test_func "nslookup" "Server" $CLI nslookup localhost

# 系统命令 (只读)
test_func "ps 进程" "PID" $CLI ps aux
test_func "df 磁盘" "Filesystem" $CLI df -h
test_func "free 内存" "Mem:" $CLI free -h

# 编码命令
test_func "base64 编码" "aGVsbG8=" $CLI base64 test_workspace/encoding/test.txt
test_func "md5sum" "5d41402abc4b2a76b9719d911017c592" $CLI md5sum test_workspace/text/sample.txt
test_func "sha256sum" "sha256" $CLI sha256sum test_workspace/text/sample.txt

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
