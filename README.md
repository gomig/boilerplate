# Installation

You need install [mig cli](github.com/gomig/mig) tools to create project.

## Usage

For library usage see [gomig](https://github.com/gomig) libraries docs.

### Access App Dependencies

For accessing app dependencies (config driver, cache driver, etc.) you must use `app` namespace functions.

```go
// github.com/myapp is your app namespace
import "github.com/myapp/src/app"
app.Config()
```

### Predefined Template Pipes

#### onProd

Check if app run on production mode.

```handlebars
{{if onProd}}
<p>application run on production mode</p>
{{end}}
```

#### onDev

Check if app run on development mode.

```handlebars
{{if onDev}}
<p>application run on development mode</p>
{{end}}
```

#### templateIf

Compile template if defined.

**Parameters:**

- `name (string)`: template name.
- `data (any)`: data for template to render.

```handlebars
{{templateIf "header" .}}
```

#### defined

Check if template is defined.

**Parameters:**

-`name (string)`: template name.

```handlebars
{{if defined "header"}}
<p>header template is defined!</p>
{{end}}
```

#### uuid

Generate new uuid.

```handlebars
{{$id :=uuid}}
<div id="{{$id}}">...</div>
```

#### iif

Ternary operator.

**Parameters:**

- `cond (bool)`: boolean condition value.
- `yes (any)`: value to return on true state.
- `no (any)`: value to return on false state.

```handlebars
{{$res := iif true "yes" "no" }}
```

#### numberF

Format number use `message.NewPrinter(language.English)` formatter.

**Parameters:**

- `format (any)`: format pattern or number.
- `v (...any)`: pass to formatter.

```handlebars
<p>Total balance: <strong>{{numberF "%s123456" "$"}}</strong></p>
<p>Simpler: <strong>{{numberF 123456}}</strong></p>
```

#### regexF

Format string using regular expression.

**Parameters:**

- `data (any)`: value to format.
- `pattern (string)`: regex pattern.
- `replacement (string)`: replacement pattern.

```handlebars
<p>Phone:
    <strong>{{regexF
            23418901
            "^(\d{4})(\d{3})(\d{1})$"
            "($1) $2-$3"
        }}</strong></p>
```

#### sizeF

Format file size to friendly string.

**Parameters:**

- `size (numeric)`: size.

```handlebars
<p>File size: <strong>{{sizeF 123451123}}</strong></p>
```

#### jalaali

Format date to jalaali string (use `gomig/jalaali` package).

**Parameters:**

- `format (string)`: date format.
- `time (time.Time)`: date value.

```handlebars
<p>امروز: <strong>{{jalaali "2006-01-02" $date}}</strong></p>
```

#### json

Encode value to json or return fallback on failed.

**Parameters:**

- `data (any)`: value to marshal.
- `fallback (string)`: fallback value if marshalling failed.

```handlebars
<script>
    var data = '{{json $data "{}"}}';
</script>
```

#### map

Generate go map from key value pairs.

**Parameters:**

- `values (...any)`: key value pairs.

```handlebars
{{ $person := map "name" "John" "age" 21 "is_male" true }}
{{ templateIf "person-card" $person }}
```

#### params

Parse parameter value from string or return fallback value (`|` key:value separated string).

**Parameters:**

- `parameters (string)`: parameters string.
- `param (string)`: parameter name.
- `fallback (string)`: fallback value.

```handlebars
{{ $params := "name:John|address:NC, Street21, No 13|age:21" }}
{{ $name := param $params "name" "Unknown"}}
{{ $address := param $params "address" "Unknown"}}
```

#### option

Check if option exists in options str (`|` separated string).

**Parameters:**

- `options (string)`: options string.
- `option (string)`: option to check.

```handlebars
{{ $brands := "Apple|Samsung|Microsoft" }}
{{ $containsApple := option $brands "Apple" }}
{{ $notGoogle := option $brands "Google" }}
```

#### isset

Check if map field exists.

**Parameters:**

- `data (map[string]any)`: map data.
- `field (string)`: field name.

```handlebars
{{ $exists := isset $myMap "user_id" }}
```

#### contains

Check if value contains key. This pipe use json encoder for converting data to map before check.

**Parameters:**

- `data (any)`: data object.
- `field (string)`: field name.

```handlebars
{{ $exists := isset .User "last_activity" }}
```

#### alter

Get map field if exists or return fallback.

**Parameters:**

- `data (map[string]any)`: data map.
- `field (string)`: field name.
- `fallback (any)`: fallback value.

```handlebars
{{ $params := alter . "name" "anonymous" }}
```

#### config

Get app config value.

**Parameters:**

- `path (string)`: config path.

```handlebars
{{ $appTitle := config "app.title" }}
```

#### linebreak

Convert new line to `<br/>` tag.

**Parameters:**

- `data (string)`: data.

```handlebars
<p>{{linebreak $content}}</p>
```

#### css

Return renderable raw css (no escape).

**Parameters:**

- `data (string)`: css raw data.

```handlebars
<html>
    <head>
        <style>
            {{css $rawCssData}}
        </style>
    </head>
</html>
```

#### html

Return renderable raw html (no escape).

**Parameters:**

- `data (string)`: html raw data.

```handlebars
<div>
    {{html $rawData}}
</div>
```

#### attr

Return renderable raw attr (no escape).

**Parameters:**

- `data (string)`: attr raw data.

```handlebars
{{ $rawAttr := attr $extra_attr_str }}
<div {{ $rawAttr }}>...</div>
```

#### attrs

Generate raw html attribute from `key:value` pair.

**Parameters:**

- `attr (...string)`: `key:value` paired attributes.

```handlebars
{{ $attrs := attrs "id:my-div" "data-value:25" "title:my div" }}
<div {{ $attrs }}>...</div>
```

#### js

Return renderable raw js (no escape).

**Parameters:**

- `data (string)`: js raw data.

```handlebars
<html>
    <head>
        <script>
            {{js $rawScript}}
        </script>
    </head>
</html>
```

#### jss

Return renderable raw js string (no escape).

**Parameters:**

- `data (string)`: js raw data.

#### url

Return renderable raw url string (no escape).

**Parameters:**

- `data (string)`: js raw data.

#### asset (only web projects)

Find asset url from file system (public directory).

**Parameters:**

- `path (string)`: base path to search in public directory.
- `pattern (string)`: file pattern to search.
- `ignore (string)`: file pattern to ignore.
- `extension (string)`: file extension.

```handlebars
<!-- files: public/dist/my-style-123541.css -->
<link rel="stylesheet" href="{{asset '/dist' 'my-style' '.map' 'css'}}" />
```

#### assets (only web projects)

Find all asset url from file system (public directory).

**Parameters:**

- `path (string)`: base path to search in public directory.
- `pattern (string)`: file pattern to search.
- `ignore (string)`: file pattern to ignore.
- `extension (string)`: file extension.

```handlebars
{{ range (assets '/dist' '' '.map' 'css')}}
<link rel="stylesheet" href="{{ . }}" />
{{ end}}
```
