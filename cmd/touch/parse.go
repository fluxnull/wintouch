package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Options struct {
	SelectCreated  bool
	SelectModified bool
	SelectAccessed bool
	SelectorSeen   bool

	NoCreate      bool
	Parents       bool
	IgnoreReparse bool
	FollowReparse bool
	Runas         bool

	SyncFile   string
	DateString string
	TouchTime  string
	EditOffset int64
	UseEdit    bool

	ShowHelp    bool
	ShowVersion bool
}

type FileTimes struct {
	Created  time.Time
	Modified time.Time
	Accessed time.Time
}

func parseArgs(args []string) (Options, []string, error) {
	var opt Options
	var files []string

	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			files = append(files, args[i+1:]...)
			break
		}
		if isHelpArg(a) {
			opt.ShowHelp = true
			continue
		}
		if a == "--version" || a == "-v" {
			opt.ShowVersion = true
			continue
		}

		if strings.HasPrefix(a, "--") {
			switch a {
			case "--accessed":
				opt.SelectAccessed = true
				opt.SelectorSeen = true
			case "--modified":
				opt.SelectModified = true
				opt.SelectorSeen = true
			case "--created":
				opt.SelectCreated = true
				opt.SelectorSeen = true
			case "--no-create":
				opt.NoCreate = true
			case "--parents":
				opt.Parents = true
			case "--ignore-reparse":
				opt.IgnoreReparse = true
			case "--follow-reparse":
				opt.FollowReparse = true
			case "--runas":
				opt.Runas = true
			case "--sync":
				v, ni, err := requireNext(args, i, a)
				if err != nil {
					return opt, nil, err
				}
				i = ni
				opt.SyncFile = v
			case "--date":
				v, ni, err := requireNext(args, i, a)
				if err != nil {
					return opt, nil, err
				}
				i = ni
				opt.DateString = v
			case "--edit":
				v, ni, err := requireNext(args, i, a)
				if err != nil {
					return opt, nil, err
				}
				i = ni
				off, err := parseOffset(v)
				if err != nil {
					return opt, nil, err
				}
				opt.EditOffset = off
				opt.UseEdit = true
			default:
				return opt, nil, fmt.Errorf("unknown option %s", a)
			}
			continue
		}

		if strings.HasPrefix(a, "-") && a != "-" {
			if len(a) == 2 {
				ch := a[1]
				switch ch {
				case 's', 'd', 't', 'e':
					v, ni, err := requireNext(args, i, "-"+string(ch))
					if err != nil {
						return opt, nil, err
					}
					i = ni
					if err := assignValueOption(&opt, ch, v); err != nil {
						return opt, nil, err
					}
					continue
				}
			}

			cluster := a[1:]
			if cluster == "" {
				files = append(files, a)
				continue
			}
			for _, ch := range cluster {
				switch ch {
				case 'a':
					opt.SelectAccessed = true
					opt.SelectorSeen = true
				case 'm':
					opt.SelectModified = true
					opt.SelectorSeen = true
				case 'c':
					opt.SelectCreated = true
					opt.SelectorSeen = true
				case 'n':
					opt.NoCreate = true
				case 'p':
					opt.Parents = true
				case 'i':
					opt.IgnoreReparse = true
				case 'f':
					opt.FollowReparse = true
				case 'r':
					opt.Runas = true
				case 's', 'd', 't', 'e':
					return opt, nil, fmt.Errorf("option -%c requires a separated value token; do not place it in a short-option cluster", ch)
				default:
					return opt, nil, fmt.Errorf("unknown option -%c", ch)
				}
			}
			continue
		}

		files = append(files, a)
	}

	if !opt.SelectorSeen {
		opt.SelectCreated = true
		opt.SelectModified = true
		opt.SelectAccessed = true
	}

	if opt.IgnoreReparse && opt.FollowReparse {
		return opt, nil, errors.New("-i/--ignore-reparse and -f/--follow-reparse cannot be used together")
	}

	baseSources := 0
	if opt.SyncFile != "" {
		baseSources++
	}
	if opt.DateString != "" {
		baseSources++
	}
	if opt.TouchTime != "" {
		baseSources++
	}
	if baseSources > 1 {
		return opt, nil, errors.New("use only one base time source: -s, -d, or -t")
	}

	return opt, files, nil
}

func isHelpArg(arg string) bool {
	switch strings.ToLower(arg) {
	case "/?", "/h", "/help", "-h", "-help", "--help":
		return true
	default:
		return false
	}
}

func requireNext(args []string, i int, name string) (string, int, error) {
	if i+1 >= len(args) {
		return "", i, fmt.Errorf("option %s requires a value", name)
	}
	return args[i+1], i + 1, nil
}

func assignValueOption(opt *Options, ch byte, value string) error {
	switch ch {
	case 's':
		opt.SyncFile = value
	case 'd':
		opt.DateString = value
	case 't':
		opt.TouchTime = value
	case 'e':
		off, err := parseOffset(value)
		if err != nil {
			return err
		}
		opt.EditOffset = off
		opt.UseEdit = true
	default:
		return fmt.Errorf("internal parser error: unsupported value option -%c", ch)
	}
	return nil
}

func parseOffset(s string) (int64, error) {
	orig := s
	neg := false
	if strings.HasPrefix(s, "+") {
		s = s[1:]
	} else if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}
	if !isDigits(s) || !(len(s) == 2 || len(s) == 4 || len(s) == 6) {
		return 0, fmt.Errorf("invalid -e/--edit offset %q: use [+|-][[hh]mm]SS", orig)
	}
	sec, min, hour := 0, 0, 0
	switch len(s) {
	case 2:
		sec = atoi(s[0:2])
	case 4:
		min = atoi(s[0:2])
		sec = atoi(s[2:4])
	case 6:
		hour = atoi(s[0:2])
		min = atoi(s[2:4])
		sec = atoi(s[4:6])
	}
	if min > 59 || sec > 59 {
		return 0, fmt.Errorf("invalid -e/--edit offset %q: minutes and seconds must be 00-59", orig)
	}
	total := int64(hour*3600 + min*60 + sec)
	if neg {
		total = -total
	}
	return total, nil
}

func parseDate(s string) (time.Time, error) {
	orig := s
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("invalid -d/--date value: empty date")
	}
	utc := false
	if strings.HasSuffix(s, "Z") {
		utc = true
		s = strings.TrimSuffix(s, "Z")
	}
	s = strings.Replace(s, "T", " ", 1)
	loc := time.Local
	if utc {
		loc = time.UTC
	}
	for _, layout := range []string{"2006-01-02 15:04:05.999999999", "2006-01-02 15:04:05"} {
		if t, err := time.ParseInLocation(layout, s, loc); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid -d/--date value %q", orig)
}

func parseTouchTime(s string) (time.Time, error) {
	main, sec, err := splitDotSeconds(s)
	if err != nil {
		return time.Time{}, err
	}
	if !isDigits(main) {
		return time.Time{}, fmt.Errorf("invalid -t timestamp %q: use [[CC]YY]MMDDhhmm[.SS]", s)
	}
	now := time.Now()
	year := now.Year()
	var mon, day, hour, min int
	switch len(main) {
	case 8:
		mon = atoi(main[0:2])
		day = atoi(main[2:4])
		hour = atoi(main[4:6])
		min = atoi(main[6:8])
	case 10:
		yy := atoi(main[0:2])
		if yy >= 69 {
			year = 1900 + yy
		} else {
			year = 2000 + yy
		}
		mon = atoi(main[2:4])
		day = atoi(main[4:6])
		hour = atoi(main[6:8])
		min = atoi(main[8:10])
	case 12:
		year = atoi(main[0:4])
		mon = atoi(main[4:6])
		day = atoi(main[6:8])
		hour = atoi(main[8:10])
		min = atoi(main[10:12])
	default:
		return time.Time{}, fmt.Errorf("invalid -t timestamp %q: use [[CC]YY]MMDDhhmm[.SS]", s)
	}
	return checkedDate(year, mon, day, hour, min, sec)
}

func splitDotSeconds(s string) (string, int, error) {
	main := s
	sec := 0
	if dot := strings.LastIndexByte(s, '.'); dot >= 0 {
		main = s[:dot]
		ss := s[dot+1:]
		if len(ss) != 2 || !isDigits(ss) {
			return "", 0, fmt.Errorf("invalid seconds in timestamp %q", s)
		}
		sec = atoi(ss)
		if sec < 0 || sec > 59 {
			return "", 0, fmt.Errorf("seconds out of range in timestamp %q", s)
		}
	}
	return main, sec, nil
}

func checkedDate(year, mon, day, hour, min, sec int) (time.Time, error) {
	if mon < 1 || mon > 12 || day < 1 || day > 31 || hour < 0 || hour > 23 || min < 0 || min > 59 || sec < 0 || sec > 59 {
		return time.Time{}, fmt.Errorf("calendar time out of range")
	}
	t := time.Date(year, time.Month(mon), day, hour, min, sec, 0, time.Local)
	if t.Year() != year || int(t.Month()) != mon || t.Day() != day || t.Hour() != hour || t.Minute() != min || t.Second() != sec {
		return time.Time{}, fmt.Errorf("invalid calendar time")
	}
	return t, nil
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func atoi(s string) int { v, _ := strconv.Atoi(s); return v }
