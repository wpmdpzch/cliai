package pkgcmd

import (
	"fmt"
	"io"
	"net"
	"strings"
)

// ExecHost 执行 host 命令
func ExecHost(args []string) error {
	var buf strings.Builder
	err := ExecHostToWriter(args, &buf)
	if err != nil {
		return err
	}
	fmt.Print(buf.String())
	return nil
}

// ExecHostToWriter 执行 host 并写入指定 writer
func ExecHostToWriter(args []string, w io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("用法: host <域名>")
	}

	domain := args[0]

	// 简单的 -v 或 -t 参数忽略（兼容性问题，这里不实现详细模式）
	ips, err := net.LookupHost(domain)
	if err != nil {
		return fmt.Errorf("host: 无法解析 \"%s\": %s", domain, err.Error())
	}

	for _, ip := range ips {
		fmt.Fprintln(w, ip)
	}

	return nil
}
