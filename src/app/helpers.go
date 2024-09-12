package app

import (
	"io"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gomig/crypto"
	"github.com/gomig/jalaali"
	"github.com/gomig/utils"
	"github.com/google/uuid"
	ptime "github.com/yaa110/go-persian-calendar"
)

// ConfigPath get configs path
func ConfigPath(sub ...string) string {
	return NormalizeURI(path.Join(append([]string{".", "config"}, sub...)...))
}

// LogPath get logs path
func LogPath(sub ...string) string {
	return NormalizeURI(path.Join(append([]string{".", ".logs"}, sub...)...))
}

// StoragePath get storage path
func StoragePath(sub ...string) string {
	return NormalizeURI(path.Join(append([]string{".", ".storage"}, sub...)...))
}

// DatabasePath get storage path
func DatabasePath(sub ...string) string {
	return NormalizeURI(path.Join(append([]string{".", "database"}, sub...)...))
}

// <%if eq .web "y"%>
// PublicPath get public path
func PublicPath(sub ...string) string {
	return NormalizeURI(path.Join(append([]string{".", confOrPanic().Cast("web.public").StringSafe("./public")}, sub...)...))
} // <% end %>

// UniqueFile generate unique hashed filename
func UniqueFile(name string, params ...string) string {
	return utils.VarOrPanic(Crypto().HashFilename(uuid.NewString()+name, crypto.MD5))
}

// NormalizeURI normalize uri path using slashes
func NormalizeURI(url string) string {
	rx := regexp.MustCompile(`\/+`)
	url = rx.ReplaceAllString(filepath.ToSlash(url), "/")
	return path.Join(".", url)
}

// NormalizeURL normalize url path using slashes
func NormalizeURL(url string) string {
	rx := regexp.MustCompile(`\/+`)
	url = rx.ReplaceAllString(filepath.ToSlash(url), "/")
	return "/" + strings.TrimLeft(strings.TrimLeft(url, "."), "/")
}

// FindFile find hashed file path
//
// return empty string if file not found
func FindFile(base, pattern, ignore, ext string) string {
	pattern = utils.If(ext == "", pattern+".*", pattern+`.*\.`+ext)
	ignore = utils.If(ignore != "", ".*"+ignore+".*", "")
	base = "./" + filepath.ToSlash(base)
	for _, f := range utils.FindFile(base, pattern) {
		if ignore == "" {
			return NormalizeURI(f)
		} else {
			if ok, err := regexp.MatchString(ignore, f); err == nil && !ok {
				return NormalizeURI(f)
			}
		}
	}
	return ""
}

// NMinute return duration for n minute
func NMinute(n int) time.Duration {
	return time.Duration(n) * time.Minute
}

// Now get current time in utc
func Now() time.Time {
	return time.Now().UTC()
}

// NowPtr get current time pointer in utc
func NowPtr() *time.Time {
	t := time.Now().UTC()
	return &t
}

// <%if eq .locale "fa"%>
// JTime set jalaali time to date
//
// pass -1 to ignore params
func JTime(time time.Time, hour int, min int, sec int) time.Time {
	jTime := jalaali.NewTehran(time).JTime()
	jTime.SetHour(utils.If(hour == -1, jTime.Hour(), hour))
	jTime.SetMinute(utils.If(min == -1, jTime.Minute(), min))
	jTime.SetSecond(utils.If(sec == -1, jTime.Second(), sec))
	return jTime.Time()
}

// ParseJalaali parse persian date string
func ParseJalaali(str string, h, m, s int) *time.Time {
	isNumber := func(v string) bool { return regexp.MustCompile(`^[0-9]+$`).MatchString(v) }

	// Extract digits
	pattern := regexp.MustCompile(`[^\d]+`)
	digits := strings.Split(pattern.ReplaceAllString(str, "-"), "-")

	// Normalize time
	now := time.Now()
	h = utils.If(h < 0, now.Hour(), h)
	m = utils.If(m < 0, now.Minute(), m)
	s = utils.If(s < 0, now.Second(), s)

	// Generate dates
	if len(digits) == 3 {
		var Y, M, D int
		for i, digit := range digits {
			if isNumber(digit) {
				switch i {
				case 0:
					Y, _ = strconv.Atoi(digit)
				case 1:
					M, _ = strconv.Atoi(digit)
				case 2:
					D, _ = strconv.Atoi(digit)
				}
			}
		}
		if Y > 0 &&
			M > 0 && M < 13 &&
			D > 0 && D < 32 {
			t := ptime.Date(Y, ptime.Month(M), D, h, m, s, 0, jalaali.TehranTz())
			if t.Year() == Y && int(t.Month()) == M && t.Day() == D {
				res := t.Time().UTC()
				return &res
			}
		}
	}
	return nil
}

// FormatJalaali format time as jalaali
func FormatJalaali(time time.Time, format string) string {
	t := jalaali.NewTehran(time).JTime()
	return t.Format(format)
} // <% end %>

// IsProd check if project run on production mode
func IsProd() bool {
	if _container != nil && _container.Exists("--prod") {
		return _container.Cast("--prod").BoolSafe(true)
	}
	if Config() != nil {
		return Config().Cast("prod").BoolSafe(true)
	}
	return true
}

// MCRequest create http microservice request with X-MC-KEY header (read from mc_key config)
func MCRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-MC-KEY", confOrPanic().Cast("mc_key").StringSafe("-"))
	return req, nil
}
