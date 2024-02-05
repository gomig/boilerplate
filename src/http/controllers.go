package http

import (
	"__ns__/src/app"
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/gomig/caster"
	"github.com/gomig/utils"
	"golang.org/x/image/webp"
)

// Thumbnail get thumbnail of file
// #GET /thumb/:path
// [200] thumbnail image
func Thumbnail(c *fiber.Ctx) error {
	// parse request params
	_path := c.Params("*", "-")
	_height := caster.NewCaster(c.Query("h")).IntSafe(0)
	_width := caster.NewCaster(c.Query("w")).IntSafe(0)

	// read file
	if ok, _ := utils.FileExists(app.PublicPath(_path)); !ok {
		return c.SendStatus(404)
	}
	content, err := os.ReadFile(app.PublicPath(_path))
	utils.PanicOnError(err)

	// detect and validate mime
	mime := utils.DetectMime(content)
	// video mime
	if mime.Is("video/mp4") {
		_path := strings.Replace(_path, ".mp4", ".jpg", 1)
		if ok, _ := utils.FileExists(app.PublicPath(_path)); !ok {
			return c.SendStatus(404)
		}
		content, err = os.ReadFile(app.PublicPath(_path))
		utils.PanicOnError(err)
		mime = utils.DetectMime(content)
	} else { // images
		mimes := []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
		if !utils.Contains(mimes, mime.String()) {
			c.Set(fiber.HeaderContentType, mime.String())
			return c.Send(content)
		}
	}

	// decode and
	var source image.Image
	if mime.Is("image/webp") {
		source, err = webp.Decode(bytes.NewReader(content))
		utils.PanicOnError(err)
	} else {
		source, _, err = image.Decode(bytes.NewReader(content))
		utils.PanicOnError(err)
	}

	// calc deminsion
	if _height == 0 && _width == 0 {
		_height = 150
	}
	dx := _width
	dy := _height
	if dx == 0 {
		dx = (source.Bounds().Dx() * dy) / source.Bounds().Dy()
	}
	if dy == 0 {
		dy = (source.Bounds().Dy() * dx) / source.Bounds().Dx()
	}

	// generate thumb
	_thumb := imaging.Thumbnail(source, dx, dy, imaging.Lanczos)
	var buf bytes.Buffer
	if mime.Is("image/jpeg") {
		utils.PanicOnError(jpeg.Encode(&buf, _thumb, &jpeg.Options{Quality: jpeg.DefaultQuality}))
	} else {
		utils.PanicOnError(png.Encode(&buf, _thumb))
	}
	c.Set(fiber.HeaderContentType, utils.If(mime.Is("image/jpeg"), "image/jpeg", "image/png"))
	return c.Send(buf.Bytes())
}
