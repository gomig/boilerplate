# GoMig

For more information see [Boilerplate Repository](github.com/gomig/boilerplate).

## Build For Linux

```bash
$Env:GOOS = "linux"; $Env:GOARCH = "amd64"
go build
```

{{ if eq .web "y"}}
## Create Web Layout

For creating base web layout template use following syntax: (`/views/Layout.tpl`)

```handlebars
{{define "base"}}
<html>
    <body>
        <p>template must defined inside base block!</p>
         {{ template "content" . }}
    </body>
</html>
{{end}}
```

Use layout in template: (`/views/pages/home.tpl`)

```handlebars
{{define "content"}}
    <p>Welcome to my site.</p>
{{end}}
```

Render template:

```go
router.Get("/", func(c *fiber.Ctx) error {
    return app.Render(c, "pages/home", "layout", data, 201)
})
```
{{ end }}
