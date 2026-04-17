package pkgcmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExecCd 执行 cd 命令（Go 原生实现）
func ExecCd(args []string) error {
	if len(args) == 0 {
		// cd 到 HOME 目录
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("无法获取 HOME 目录: %v", err)
		}
		return os.Chdir(home)
	}

	dir := args[0]
	// 处理 ~
	if strings.HasPrefix(dir, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("无法获取 HOME 目录: %v", err)
		}
		dir = filepath.Join(home, dir[2:])
	} else if dir == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("无法获取 HOME 目录: %v", err)
		}
		dir = home
	}

	return os.Chdir(dir)
}

// ExecEcho 执行 echo 命令（Go 原生实现）
func ExecEcho(args []string, w io.Writer) error {
	// 处理 -n 参数（不换行）
	noNewline := false
	filtered := args
	if len(args) > 0 && args[0] == "-n" {
		noNewline = true
		filtered = args[1:]
	}

	result := strings.Join(filtered, " ")
	if !noNewline {
		result += "\n"
	}
	fmt.Fprint(w, result)
	return nil
}

// ExecPwd 执行 pwd 命令（Go 原生实现）
func ExecPwd(w io.Writer) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}
	fmt.Fprintln(w, dir)
	return nil
}

// ExecMkdir 执行 mkdir 命令（Go 原生实现）
func ExecMkdir(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("用法: mkdir <目录>")
	}

	// 默认参数
	parents := false
	var dirs []string

	for _, arg := range args {
		if arg == "-p" {
			parents = true
		} else if !strings.HasPrefix(arg, "-") {
			dirs = append(dirs, arg)
		}
	}

	for _, dir := range dirs {
		if parents {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return fmt.Errorf("创建目录 %s 失败: %v", dir, err)
			}
		} else {
			err := os.Mkdir(dir, 0755)
			if err != nil {
				return fmt.Errorf("创建目录 %s 失败: %v", dir, err)
			}
		}
	}
	return nil
}

// ExecTouch 执行 touch 命令（Go 原生实现）
func ExecTouch(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("用法: touch <文件>")
	}

	for _, file := range args {
		if strings.HasPrefix(file, "-") {
			continue
		}
		f, err := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("创建文件 %s 失败: %v", file, err)
		}
		f.Close()
	}
	return nil
}

// ExecRm 执行 rm 命令（Go 原生实现，危险命令标记）
func ExecRm(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("用法: rm <文件>")
	}

	recursive := false
	force := false
	var files []string

	for _, arg := range args {
		switch arg {
		case "-r", "-R":
			recursive = true
		case "-f", "--force":
			force = true
		case "-rf", "-fr", "-rRf", "-Rrf":
			recursive = true
			force = true
		default:
			if !strings.HasPrefix(arg, "-") {
				files = append(files, arg)
			}
		}
	}

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			if !force {
				return fmt.Errorf("文件 %s 不存在", file)
			}
			continue
		}

		if info.IsDir() && !recursive {
			return fmt.Errorf("不能删除目录 %s，请使用 -r", file)
		}

		if info.IsDir() {
			err = os.RemoveAll(file)
		} else {
			err = os.Remove(file)
		}
		if err != nil {
			return fmt.Errorf("删除 %s 失败: %v", file, err)
		}
	}
	return nil
}

// ExecCp 执行 cp 命令（Go 原生实现）
func ExecCp(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("用法: cp <源文件> <目标文件>")
	}

	recursive := false
	var src, dst string

	for _, arg := range args {
		if arg == "-r" || arg == "-R" {
			recursive = true
		} else if !strings.HasPrefix(arg, "-") {
			if src == "" {
				src = arg
			} else {
				dst = arg
			}
		}
	}

	if src == "" || dst == "" {
		return fmt.Errorf("用法: cp <源文件> <目标文件>")
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("源文件 %s 不存在: %v", src, err)
	}

	if srcInfo.IsDir() && !recursive {
		return fmt.Errorf("cp: 目录 %s，请使用 -r", src)
	}

	if srcInfo.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func copyDir(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	os.MkdirAll(dst, info.Mode())

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// ExecMv 执行 mv 命令（Go 原生实现）
func ExecMv(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("用法: mv <源文件> <目标文件>")
	}

	force := false
	var src, dst string

	for _, arg := range args {
		if arg == "-f" || arg == "--force" {
			force = true
		} else if !strings.HasPrefix(arg, "-") {
			if src == "" {
				src = arg
			} else {
				dst = arg
			}
		}
	}

	if src == "" || dst == "" {
		return fmt.Errorf("用法: mv <源文件> <目标文件>")
	}

	// 检查目标是否存在
	if _, err := os.Stat(dst); err == nil && !force {
		return fmt.Errorf("目标文件 %s 已存在，请使用 -f 覆盖", dst)
	}

	return os.Rename(src, dst)
}

// ExecClear 执行 clear 命令（Go 原生实现）
func ExecClear(w io.Writer) error {
	fmt.Fprint(w, "\033[2J\033[H")
	return nil
}

// ExecExit 执行 exit 命令（Go 原生实现）
func ExecExit(args []string) error {
	code := 0
	if len(args) > 0 {
		fmt.Sscanf(args[0], "%d", &code)
	}
	os.Exit(code)
	return nil
}

// ExecWhich 执行 which 命令（Go 原生实现）
func ExecWhich(args []string, w io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("用法: which <命令>")
	}

	cmdName := args[0]

	// 检查是否是内置命令
	if cmd, ok := commands[cmdName]; ok {
		if cmd.Implemented == "builtin" {
			fmt.Fprintf(w, "%s: builtin command\n", cmdName)
			return nil
		}
		if cmd.Implemented == "go" {
			fmt.Fprintf(w, "%s: cliai native command\n", cmdName)
			return nil
		}
	}

	// 检查系统命令路径
	pathDirs := strings.Split(os.Getenv("PATH"), ":")
	for _, dir := range pathDirs {
		path := filepath.Join(dir, cmdName)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			// 检查是否有执行权限
			if info.Mode()&0111 != 0 {
				fmt.Fprintln(w, path)
				return nil
			}
		}
	}

	return fmt.Errorf("%s not found", cmdName)
}

// ExecHistory 执行 history 命令（简单实现）
func ExecHistory(w io.Writer) error {
	// TODO: 实现真正的历史记录
	fmt.Fprintln(w, "1  cd /tmp")
	fmt.Fprintln(w, "2  ls -la")
	fmt.Fprintln(w, "3  pwd")
	return nil
}

// ExecGoCmdInternal 执行 Go 原生内建命令
func ExecGoCmdInternal(name string, args []string, stdout io.Writer) *ExecResult {
	var buf bytes.Buffer
	if stdout == nil {
		stdout = &buf
	}

	var err error
	switch name {
	case "cd":
		err = ExecCd(args)
	case "echo":
		err = ExecEcho(args, stdout)
	case "pwd":
		err = ExecPwd(stdout)
	case "mkdir":
		err = ExecMkdir(args)
	case "touch":
		err = ExecTouch(args)
	case "rm":
		err = ExecRm(args)
	case "cp":
		err = ExecCp(args)
	case "mv":
		err = ExecMv(args)
	case "clear":
		err = ExecClear(stdout)
	case "exit":
		err = ExecExit(args)
	case "which":
		err = ExecWhich(args, stdout)
	case "history":
		err = ExecHistory(stdout)
	default:
		return &ExecResult{Output: "", Error: fmt.Errorf("未知内建命令: %s", name)}
	}

	return &ExecResult{Output: buf.String(), Error: err}
}
