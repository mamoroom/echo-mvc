package res_redirect

import (
	_ "fmt"
	"github.com/labstack/echo"
	"net/http"
)

//303
func RedirectTo(c echo.Context, url string) error {
	return c.Redirect(http.StatusSeeOther, url)
}

func RedirectToIndexHtml(c echo.Context, opts ...map[string]string) error {
	return _redirect_to_html_page(c, "/", opts...)
}

func RedirectToTopHtml(c echo.Context, opts ...map[string]string) error {
	return _redirect_to_html_page(c, "/top.html", opts...)
}

func RedirectToOpeningHtml(c echo.Context, opts ...map[string]string) error {
	return _redirect_to_html_page(c, "/opening.html", opts...)
}

func RedirectToNazoDetailHtml(c echo.Context, opts ...map[string]string) error {
	return _redirect_to_html_page(c, "/mystery.html", opts...)
}

func RedirectToIncentiveHtml(c echo.Context, opts ...map[string]string) error {
	return _redirect_to_html_page(c, "/incentive_inherit.html", opts...)
}

func RedirectToErrorHtml(c echo.Context, opts ...map[string]string) error {
	return _redirect_to_html_page(c, "/systemError.html", opts...)
}

func RedirectTo404Html(c echo.Context, opts ...map[string]string) error {
	return _redirect_to_html_page(c, "/404.html", opts...)
}

func _redirect_to_html_page(c echo.Context, uri string, opts ...map[string]string) error {
	req, _ := http.NewRequest("GET", uri, nil)
	q := req.URL.Query()

	var params = map[string]string{}
	if len(opts) > 0 {
		params = opts[0]
	}
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	return RedirectTo(c, req.URL.String())
}
