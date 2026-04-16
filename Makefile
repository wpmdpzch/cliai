# CLI-AI Makefile

.PHONY: build clean install test run

# 编译
build:
	cd cliai-go && go build -o cliai ./cmd/cliai

# 清理
clean:
	rm -f cliai-go/cliai

# 安装到系统
install: build
	sudo mv cliai-go/cliai /usr/local/bin/cliai

# 本地运行
run:
	cd cliai-go && go run ./cmd/cliai

# 测试
test:
	cd cliai-go && go test ./...

# 查看命令列表
commands:
	cd cliai-go && go run ./cmd/cliai commands

# 帮助
help:
	@echo "CLI-AI Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build     # 编译"
	@echo "  make install   # 编译并安装到系统"
	@echo "  make run       # 本地运行"
	@echo "  make test      # 运行测试"
	@echo "  make commands  # 查看内置命令"
	@echo "  make clean     # 清理编译文件"
