package html

import (
	"github.com/labstack/echo"

	"github.com/mamoroom/echo-mvc/app/config"
	"github.com/mamoroom/echo-mvc/app/router/handler"

	"html/template"
	"io"
)

var conf = config.Conf

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Init(e *echo.Echo) {
	t := &Template{
		templates: template.Must(template.ParseGlob(config.ROOT_PATH + "/view/res_html/*.html")),
	}
	e.Renderer = t
	h := e.Group("")
	h.GET("/share/:hash", ShareHandler)
	h.GET("/share/static/:lang", ShareStaticHandler)

	d := e.Group("")
	d.Use(handler.DebugHandler())
	{
	}
}
