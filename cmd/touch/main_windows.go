//go:build windows

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	FileReadAttributes        = 0x0080
	FileWriteAttributes       = 0x0100
	FileShareRead             = 0x00000001
	FileShareWrite            = 0x00000002
	FileShareDelete           = 0x00000004
	OpenExisting              = 3
	FileAttributeNormal       = 0x00000080
	FileFlagBackupSemantics   = 0x02000000
	FileAttributeReparsePoint = 0x00000400
	SEE_MASK_NOCLOSEPROCESS   = 0x00000040
	SW_SHOWNORMAL             = 1
	WAIT_INFINITE             = 0xffffffff
)

var (
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	shell32               = syscall.NewLazyDLL("shell32.dll")
	procGetFileAttributes = kernel32.NewProc("GetFileAttributesW")
	procShellExecuteEx    = shell32.NewProc("ShellExecuteExW")
	procIsUserAnAdmin     = shell32.NewProc("IsUserAnAdmin")
)

type shellExecuteInfo struct {
	cbSize       uint32
	fMask        uint32
	hwnd         uintptr
	lpVerb       *uint16
	lpFile       *uint16
	lpParameters *uint16
	lpDirectory  *uint16
	nShow        int32
	hInstApp     uintptr
	lpIDList     uintptr
	lpClass      *uint16
	hkeyClass    uintptr
	dwHotKey     uint32
	hIcon        uintptr
	hProcess     syscall.Handle
}

func main() {
	opt, files, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "touch: %v\n", err)
		os.Exit(2)
	}
	if opt.ShowHelp {
		fmt.Print(HelpText)
		return
	}
	if opt.ShowVersion {
		fmt.Printf("touch.exe %s\n", Version)
		return
	}
	if opt.Runas {
		if err := maybeRunElevated(os.Args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "touch: runas failed: %v\n", err)
			os.Exit(1)
		}
	}
	if len(files) == 0 {
		fmt.Fprint(os.Stderr, HelpText)
		os.Exit(2)
	}

	expanded, err := expandTargets(files, opt.NoCreate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "touch: %v\n", err)
		os.Exit(1)
	}
	if len(expanded) == 0 {
		return
	}

	base, hasBase, err := resolveBaseTimes(opt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "touch: %v\n", err)
		os.Exit(1)
	}
	if hasBase && opt.UseEdit {
		base = addOffsetToSelected(base, opt)
	}

	exitCode := 0
	for _, path := range expanded {
		if err := touchOne(path, opt, base, hasBase); err != nil {
			fmt.Fprintf(os.Stderr, "touch: %s: %v\n", path, err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func resolveBaseTimes(opt Options) (FileTimes, bool, error) {
	if opt.SyncFile != "" {
		if opt.IgnoreReparse {
			isRep, err := isReparsePoint(opt.SyncFile)
			if err != nil {
				return FileTimes{}, false, err
			}
			if isRep {
				return FileTimes{}, false, fmt.Errorf("sync source is a reparse point and -i/--ignore-reparse was used: %s", opt.SyncFile)
			}
		}
		ft, err := getFileTimes(opt.SyncFile)
		return ft, true, err
	}
	if opt.DateString != "" {
		t, err := parseDate(opt.DateString)
		if err != nil {
			return FileTimes{}, false, err
		}
		return FileTimes{Created: t, Modified: t, Accessed: t}, true, nil
	}
	if opt.TouchTime != "" {
		t, err := parseTouchTime(opt.TouchTime)
		if err != nil {
			return FileTimes{}, false, err
		}
		return FileTimes{Created: t, Modified: t, Accessed: t}, true, nil
	}
	now := time.Now()
	return FileTimes{Created: now, Modified: now, Accessed: now}, false, nil
}

func touchOne(path string, opt Options, base FileTimes, hasBase bool) error {
	if opt.IgnoreReparse {
		isRep, err := isReparsePoint(path)
		if err == nil && isRep {
			return nil
		}
		if err != nil && !isNotFound(err) {
			return err
		}
	}

	exists, err := exists(path)
	if err != nil {
		return err
	}
	if !exists {
		if opt.NoCreate {
			return nil
		}
		if opt.Parents {
			dir := filepath.Dir(path)
			if dir != "." && dir != "" {
				if err := os.MkdirAll(windowsPath(dir), 0755); err != nil {
					return err
				}
			}
		}
		f, err := os.OpenFile(windowsPath(path), os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
		exists = true
	}

	ft := base
	if opt.UseEdit && !hasBase {
		cur, err := getFileTimes(path)
		if err != nil {
			return err
		}
		ft = addOffsetToSelected(cur, opt)
	} else if !hasBase {
		// Default base is the time captured once before iteration. This makes multi-file
		// default touch behavior deterministic inside one invocation.
		ft = base
	}

	return setFileTimes(path, ft, opt)
}

func addOffsetToSelected(ft FileTimes, opt Options) FileTimes {
	d := time.Duration(opt.EditOffset) * time.Second
	if opt.SelectCreated {
		ft.Created = ft.Created.Add(d)
	}
	if opt.SelectModified {
		ft.Modified = ft.Modified.Add(d)
	}
	if opt.SelectAccessed {
		ft.Accessed = ft.Accessed.Add(d)
	}
	return ft
}

func exists(path string) (bool, error) {
	_, err := getFileTimes(path)
	if err == nil {
		return true, nil
	}
	if isNotFound(err) {
		return false, nil
	}
	return false, err
}

func isNotFound(err error) bool {
	return os.IsNotExist(err) || errors.Is(err, syscall.ERROR_FILE_NOT_FOUND) || errors.Is(err, syscall.ERROR_PATH_NOT_FOUND)
}

func getFileTimes(path string) (FileTimes, error) {
	h, err := openAttributes(path, false)
	if err != nil {
		return FileTimes{}, err
	}
	defer syscall.CloseHandle(h)
	var info syscall.ByHandleFileInformation
	if err := syscall.GetFileInformationByHandle(h, &info); err != nil {
		return FileTimes{}, err
	}
	return FileTimes{
		Created:  time.Unix(0, info.CreationTime.Nanoseconds()).Local(),
		Accessed: time.Unix(0, info.LastAccessTime.Nanoseconds()).Local(),
		Modified: time.Unix(0, info.LastWriteTime.Nanoseconds()).Local(),
	}, nil
}

func setFileTimes(path string, ft FileTimes, opt Options) error {
	h, err := openAttributes(path, true)
	if err != nil {
		return err
	}
	defer syscall.CloseHandle(h)

	var cPtr, aPtr, mPtr *syscall.Filetime
	if opt.SelectCreated {
		c := toFiletime(ft.Created)
		cPtr = &c
	}
	if opt.SelectAccessed {
		a := toFiletime(ft.Accessed)
		aPtr = &a
	}
	if opt.SelectModified {
		m := toFiletime(ft.Modified)
		mPtr = &m
	}
	return syscall.SetFileTime(h, cPtr, aPtr, mPtr)
}

func toFiletime(t time.Time) syscall.Filetime {
	return syscall.NsecToFiletime(t.UnixNano())
}

func openAttributes(path string, write bool) (syscall.Handle, error) {
	access := uint32(FileReadAttributes)
	if write {
		access = FileWriteAttributes
	}
	flags := uint32(FileAttributeNormal | FileFlagBackupSemantics)
	p, err := syscall.UTF16PtrFromString(windowsPath(path))
	if err != nil {
		return 0, err
	}
	return syscall.CreateFile(p, access, FileShareRead|FileShareWrite|FileShareDelete, nil, OpenExisting, flags, 0)
}

func isReparsePoint(path string) (bool, error) {
	p, err := syscall.UTF16PtrFromString(windowsPath(path))
	if err != nil {
		return false, err
	}
	r1, _, e1 := procGetFileAttributes.Call(uintptr(unsafe.Pointer(p)))
	if r1 == ^uintptr(0) {
		if e1 != syscall.Errno(0) {
			return false, e1
		}
		return false, syscall.GetLastError()
	}
	return uint32(r1)&FileAttributeReparsePoint != 0, nil
}

func windowsPath(path string) string {
	if strings.HasPrefix(path, `\\?\`) || strings.HasPrefix(path, `\\.\`) {
		return path
	}
	abs, err := filepath.Abs(path)
	if err == nil {
		path = abs
	}
	path = filepath.Clean(path)
	if strings.HasPrefix(path, `\\`) {
		return `\\?\UNC\` + strings.TrimPrefix(path, `\\`)
	}
	return `\\?\` + path
}

func expandTargets(args []string, noCreate bool) ([]string, error) {
	out := make([]string, 0, len(args))
	for _, a := range args {
		if strings.ContainsAny(a, "*?") {
			m, err := filepath.Glob(a)
			if err != nil {
				return nil, err
			}
			if len(m) == 0 {
				if noCreate {
					continue
				}
				return nil, fmt.Errorf("wildcard matched no files: %s", a)
			}
			out = append(out, m...)
			continue
		}
		out = append(out, a)
	}
	return out, nil
}

func maybeRunElevated(args []string) error {
	if isElevated() {
		return nil
	}
	cleanArgs := removeRunasArgs(args)
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	verb, _ := syscall.UTF16PtrFromString("runas")
	file, _ := syscall.UTF16PtrFromString(exe)
	params, _ := syscall.UTF16PtrFromString(joinWindowsArgs(cleanArgs))
	dir, _ := syscall.UTF16PtrFromString(wd)

	sei := shellExecuteInfo{
		cbSize:       uint32(unsafe.Sizeof(shellExecuteInfo{})),
		fMask:        SEE_MASK_NOCLOSEPROCESS,
		lpVerb:       verb,
		lpFile:       file,
		lpParameters: params,
		lpDirectory:  dir,
		nShow:        SW_SHOWNORMAL,
	}
	r1, _, e1 := procShellExecuteEx.Call(uintptr(unsafe.Pointer(&sei)))
	if r1 == 0 {
		if e1 != syscall.Errno(0) {
			return e1
		}
		return syscall.GetLastError()
	}
	if sei.hProcess != 0 {
		syscall.WaitForSingleObject(sei.hProcess, WAIT_INFINITE)
		var code uint32
		syscall.GetExitCodeProcess(sei.hProcess, &code)
		syscall.CloseHandle(sei.hProcess)
		os.Exit(int(code))
	}
	os.Exit(0)
	return nil
}

func isElevated() bool {
	r1, _, _ := procIsUserAnAdmin.Call()
	return r1 != 0
}

func removeRunasArgs(args []string) []string {
	out := make([]string, 0, len(args))
	for _, a := range args {
		if a == "--runas" || a == "-r" {
			continue
		}
		if strings.HasPrefix(a, "-") && !strings.HasPrefix(a, "--") && len(a) > 2 {
			cluster := a[1:]
			if isBooleanCluster(cluster) && strings.Contains(cluster, "r") {
				clean := strings.ReplaceAll(cluster, "r", "")
				if clean != "" {
					out = append(out, "-"+clean)
				}
				continue
			}
		}
		out = append(out, a)
	}
	return out
}

func isBooleanCluster(s string) bool {
	if s == "" {
		return false
	}
	for _, ch := range s {
		switch ch {
		case 'a', 'm', 'c', 'n', 'p', 'i', 'f', 'r':
		default:
			return false
		}
	}
	return true
}

func joinWindowsArgs(args []string) string {
	parts := make([]string, len(args))
	for i, a := range args {
		parts[i] = quoteWindowsArg(a)
	}
	return strings.Join(parts, " ")
}

func quoteWindowsArg(s string) string {
	if s == "" {
		return `""`
	}
	if !strings.ContainsAny(s, " \t\n\v\"") {
		return s
	}
	var b strings.Builder
	b.WriteByte('"')
	bs := 0
	for _, r := range s {
		if r == '\\' {
			bs++
			continue
		}
		if r == '"' {
			b.WriteString(strings.Repeat("\\", bs*2+1))
			b.WriteRune(r)
			bs = 0
			continue
		}
		if bs > 0 {
			b.WriteString(strings.Repeat("\\", bs))
			bs = 0
		}
		b.WriteRune(r)
	}
	if bs > 0 {
		b.WriteString(strings.Repeat("\\", bs*2))
	}
	b.WriteByte('"')
	return b.String()
}
