package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gomig/jalaali"
	"github.com/gomig/utils"
	"github.com/google/uuid"
	"github.com/inhies/go-bytesize"
)

var _pipes template.FuncMap
var _tpl *template.Template

func init() {
	_pipes = make(template.FuncMap, 0)

	// onProd check if app run on production mode
	_pipes["onProd"] = func() bool {
		return IsProd()
	}

	// onDev check if app run on development mode
	_pipes["onDev"] = func() bool {
		return !IsProd()
	}

	// compile template if defined
	_pipes["templateIf"] = func(name string, data any) template.HTML {
		if _tpl != nil {
			t := _tpl.Lookup(name)
			if t != nil {
				buf := bytes.NewBuffer([]byte{})
				err := t.Execute(buf, data)
				if err != nil {
					return ""
				} else {
					return template.HTML(buf.String())
				}
			}
		}
		return ""
	}

	// check if template is defined
	_pipes["defined"] = func(name string) bool {
		if _tpl != nil {
			return _tpl.Lookup(name) != nil
		}
		return false
	}

	// uuid generate unique id
	_pipes["uuid"] = func() string {
		return uuid.NewString()
	}

	// iif ternary like operator
	_pipes["iif"] = func(cond bool, yes, no any) any {
		if cond {
			return yes
		}
		return no
	}

	// numberF format number
	_pipes["numberF"] = func(format string, v ...any) string {
		return utils.FormatNumber(format, v...)
	}

	// regexF format data using regex
	_pipes["regexF"] = func(data, pattern, repl string) string {
		return utils.FormatRx(data, pattern, repl)
	}

	// size format file size
	_pipes["sizeF"] = func(size float64) string {
		return fmt.Sprint(bytesize.New(size))
	}

	// {{if eq .locale "fa"}}
	// jalaali format jalaali date
	_pipes["jalaali"] = func(format string, t time.Time) string {
		return jalaali.NewTehran(t).Format(format)
	}
	// {{ end }}

	// json get json encoded value
	_pipes["json"] = func(data any, fallback string) string {
		bytes, err := json.Marshal(data)
		if err == nil {
			return string(bytes)
		}
		return fallback
	}

	// jsonFrom get json object from values
	_pipes["jsonFrom"] = func(values ...any) string {
		if len(values) > 0 && len(values)%2 == 0 {
			res := "{"
			for i := 0; i < len(values); i += 2 {
				res += fmt.Sprintf("'%s':", values[i]) + fmt.Sprint(values[i+1])
				if i+2 < len(values) {
					res += ","
				}
			}
			return res + "}"
		}
		return "{}"
	}

	// map generate map from key value pairs
	_pipes["map"] = func(values ...any) map[string]any {
		if len(values) > 0 && len(values)%2 == 0 {
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				dict[fmt.Sprint(values[i])] = values[i+1]
			}
			return dict
		}
		return nil
	}

	// param get param from "params" map item
	// params: "name:john"
	_pipes["param"] = func(params string, param string, fallback string) string {
		for _, _param := range strings.Split(strings.TrimSpace(params), "|") {
			if strings.HasPrefix(_param, param+":") {
				return strings.ReplaceAll(_param, param+":", "")
			}
		}
		return fallback
	}

	// option check if option passed in "options"
	// options: "paramA|paramB|paramC"
	_pipes["option"] = func(options string, opt string) bool {
		for _, options := range strings.Split(strings.TrimSpace(options), "|") {
			if options == opt {
				return true
			}
		}
		return false
	}

	// isset check if map field exists
	_pipes["isset"] = func(data map[string]any, field string) bool {
		_, ok := data[field]
		return ok
	}

	// contains check if value contains field
	// this function using json encoder, to detect field exists you must check json key
	_pipes["contains"] = func(data any, field string) bool {
		bytes, err := json.Marshal(data)
		if err != nil {
			return false
		}
		var res map[string]interface{}
		err = json.Unmarshal(bytes, &res)
		if err != nil {
			return false
		}
		_, ok := res[field]
		return ok
	}

	// alter get map field or return fallback
	_pipes["alter"] = func(data map[string]any, field string, fallback any) any {
		if v, ok := data[field]; ok {
			return v
		}
		return fallback
	}

	// config get config items
	_pipes["config"] = func(path string) any {
		return confOrPanic().Cast(path)
	}

	// {{if eq .web "y"}}
	// linebreak convert new line to br
	_pipes["linebreak"] = func(s string) template.HTML {
		return template.HTML(strings.ReplaceAll(s, "\n", "<br />"))
	}

	// css return renderable raw css (no escape)
	_pipes["css"] = func(v string) template.CSS {
		return template.CSS(v)
	}

	// html return renderable raw html (no escape)
	_pipes["html"] = func(v string) template.HTML {
		return template.HTML(v)
	}

	// attr return renderable raw html attr (no escape)
	_pipes["attr"] = func(v string) template.HTMLAttr {
		return template.HTMLAttr(v)
	}

	// attrs generate html attributes from parameters
	// example: attrs "id:test" "class:first second third"
	_pipes["attrs"] = func(attr ...string) template.HTMLAttr {
		res := ""
		for _, at := range attr {
			parts := strings.Split(at, ":")
			if len(parts) == 2 {
				res += parts[0] + "=\"" + parts[1] + "\""
			}
		}
		return template.HTMLAttr(res)
	}

	// js return renderable raw js (no escape)
	_pipes["js"] = func(v string) template.JS {
		return template.JS(v)
	}

	// jss return renderable raw js string (no escape)
	_pipes["jss"] = func(v string) template.JSStr {
		return template.JSStr(v)
	}

	// url return renderable raw url (no escape)
	_pipes["url"] = func(v string) template.URL {
		return template.URL(v)
	}

	// asset find asset url
	// example: asset "dist/js" "vendor-" "" "js"
	_pipes["asset"] = func(path, pattern, ignore, ext string) string {
		return UrlOf(path, pattern, ignore, ext)
	}

	// assets find assets url in public
	// example (get all js): asset "dist/js" "" "" "js"
	_pipes["assets"] = func(base, pattern, ignore, ext string) []string {
		files := make([]string, 0)
		pattern = utils.If(ext == "", pattern+".*", pattern+`.*\.`+ext)
		ignore = utils.If(ignore != "", ".*"+ignore+".*", "")
		base = PublicPath(base)
		for _, f := range utils.FindFile(base, pattern) {
			res := ""
			if ignore == "" {
				res = NormalizeURI(f)
			} else {
				if ok, err := regexp.MatchString(ignore, f); err == nil && !ok {
					res = NormalizeURI(f)
				}
			}
			if res != "" {
				files = append(files, Url(res))
			}
		}
		return files
	}
	// {{ end }}
}

func tplBase() string {
	return Config().Cast("view.base").StringSafe("./views")
}

func tplShared() []string {
	return Config().Cast("view.shared").StringSlice([]string{})
}

func tplExt() string {
	return Config().Cast("view.extension").StringSafe("tpl")
}

func tplDLeft() string {
	return Config().Cast("view.delim_left").StringSafe("{{")
}

func tplDRight() string {
	return Config().Cast("view.delim_right").StringSafe("}}")
}

func tplPath(parts ...string) string {
	return filepath.ToSlash(path.Join(parts...))
}

func tplFile(parts ...string) string {
	return tplPath(parts...) + "." + tplExt()
}

// AddPipe add new pipe to template engine
func AddPipe(name string, f any) {
	_pipes[name] = f
}

// ViewExists check if view exists
func ViewExists(view string) bool {
	exists, _ := utils.FileExists(tplFile(tplBase(), view))
	return exists
}

// CompileView compile view
func CompileView(view, layout string, data map[string]any) (string, error) {
	// init
	_tpl = template.New(view).Delims(tplDLeft(), tplDRight())

	// make files array
	files := make([]string, 0)
	files = append(files, tplFile(tplBase(), view))
	if layout != "" {
		files = append(files, tplFile(tplBase(), layout))
	}

	// add partials
	parts := strings.Split(tplPath(view), "/")
	base := utils.If(len(parts) < 2, "", parts[0])
	for _, partial := range utils.FindFile(tplPath(tplBase(), base), `.*\.partial\.`+tplExt()+"$") {
		files = append(files, tplPath("./", partial))
	}

	// add shared
	for _, shared := range tplShared() {
		for _, file := range utils.FindFile(tplPath(shared), `.*\.`+tplExt()) {
			files = append(files, tplPath("./", file))
		}
	}

	// compile
	var writer bytes.Buffer
	if tpl, err := _tpl.Funcs(_pipes).ParseFiles(files...); err != nil {
		return "", err
	} else if err := tpl.ExecuteTemplate(&writer, "base", data); err != nil {
		return "", err
	}
	return writer.String(), nil
}

// {{if eq .web "y"}}
// Render render template to response
func Render(c *fiber.Ctx, view, layout string, data map[string]any, status int) error {
	compiled, err := CompileView(view, layout, data)
	if err != nil {
		return err
	}
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.Status(status).SendString(compiled)
}

// {{end}}
