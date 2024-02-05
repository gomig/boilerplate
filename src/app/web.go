package app

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	httput "github.com/gomig/http"
	"github.com/gomig/utils"
	"github.com/microcosm-cc/bluemonday"
)

// SetupWeb driver
func SetupWeb(onError httput.ErrorCallback) {
	erLogger := ErrorLogger()
	if erLogger == nil {
		panic("failed to find Error Logger")
	}

	server := fiber.New(fiber.Config{
		DisableStartupMessage: IsProd(),
		ErrorHandler:          httput.ErrorLogger(erLogger, DateFormatter(), onError),
		ProxyHeader:           fiber.HeaderXForwardedFor,
	})
	_container.Register("--APP-SERVER", server)
}

// Server get web server driver
// leave name empty to resolve default
func Server(names ...string) *fiber.App {
	name := "--APP-SERVER"
	if len(names) > 0 {
		name = names[0]
	}
	if dep, exists := _container.Resolve(name); exists {
		if res, ok := dep.(*fiber.App); ok {
			return res
		}
	}
	return nil
}

// Url get url of file path
// return empty string if file not found
func Url(path string) string {
	path = NormalizeURI(path)
	path = strings.Replace(path, NormalizeURI(PublicPath()), "", 1)
	return utils.If(path == ".", "", path)
}

// UrlOf find url of file (path without public dir)
// this function search public dir
//
// return empty string if file not found
func UrlOf(base, pattern, ignore, ext string) string {
	base = PublicPath(base)
	return Url(FindFile(base, pattern, ignore, ext))
}

// IsLocalRequest check if request from localhost
func IsLocalRequest(c *fiber.Ctx) bool {
	return c.IP() == "127.0.0.1"
}

// MicroserviceKey get microservice key (X-MC-KEY) from to header
func MicroserviceKey(c *fiber.Ctx) string {
	return c.Get("X-MC-KEY", "")
}

// SecureAll secure string from all html tags
func SecureAll(data string, trim bool) string {
	res := bluemonday.StrictPolicy().Sanitize(data)
	if trim {
		return strings.TrimSpace(res)
	}
	return res
}

// SecureScripts secure string from script-like html tags
func SecureScripts(data string, trim bool) string {
	res := bluemonday.UGCPolicy().Sanitize(data)
	if trim {
		return strings.TrimSpace(res)
	}
	return res
}

// NormalizeHtmlText normalize escaped html text
func NormalizeHtmlText(v string) string {
	var patterns = map[string]string{
		"&#39;":   "'",
		"&#34;":   "\"",
		"&#180;":  "´",
		"&#38;":   "&",
		"&#169;":  "©",
		"&#186;":  "°",
		"&#8363;": "€",
		"&#171;":  "«",
		"&#215;":  "x",
		"&#8220;": "“",
		"&#8221;": "”",
		"&#174;":  "®",
		"&#187;":  "»",
		"&#8482;": "™",
	}
	for code, character := range patterns {
		v = strings.ReplaceAll(v, code, character)
	}
	return v
}
