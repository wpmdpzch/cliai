package pkgcmd

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)
// ExecBase64 base64 编码/解码
func ExecBase64(args []string) error {
	var buf bytes.Buffer
	err := ExecBase64ToWriter(args, &buf)
	if err != nil {
		return err
	}
	fmt.Print(buf.String())
	return nil
}

// ExecBase64ToWriter 执行 base64 并写入指定 writer
func ExecBase64ToWriter(args []string, w io.Writer) error {
	// 无参数时从 stdin 读取
	if len(args) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		encoded := base64.StdEncoding.EncodeToString(data)
		fmt.Fprint(w, encoded)
		return nil
	}

	// 参数解析
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
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(filename)
	}

	if err != nil {
		return err
	}

	if decode {
		decoded := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
		n, err := base64.StdEncoding.Decode(decoded, data)
		if err != nil {
			return err
		}
		fmt.Fprint(w, string(decoded[:n]))
	} else {
		encoded := base64.StdEncoding.EncodeToString(data)
		fmt.Fprintln(w, encoded)
	}

	return nil
}
