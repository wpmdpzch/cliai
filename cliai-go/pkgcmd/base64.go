package pkgcmd

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// ExecBase64 Go 原生 base64 实现
func ExecBase64(args []string) error {
	// 无参数时从 stdin 读取
	if len(args) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		encoded := base64.StdEncoding.EncodeToString(data)
		fmt.Print(encoded)
		return nil
	}

	// 简单参数解析
	decode := false
	var filename string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-d", "--decode":
			decode = true
		default:
			filename = args[i]
		}
	}

	var data []byte
	var err error

	if filename == "" || filename == "-" {
		// 从 stdin 读取
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(filename)
	}

	if err != nil {
		return err
	}

	if decode {
		// 解码
		decoded := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
		n, err := base64.StdEncoding.Decode(decoded, data)
		if err != nil {
			return err
		}
		fmt.Print(string(decoded[:n]))
	} else {
		// 编码
		encoded := base64.StdEncoding.EncodeToString(data)
		fmt.Println(encoded)
	}

	return nil
}
