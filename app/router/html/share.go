package html

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"

	"github.com/mamoroom/echo-mvc/app/lib/logger"
	"github.com/mamoroom/echo-mvc/app/lib/share"
	"github.com/mamoroom/echo-mvc/app/lib/share_data"
	"github.com/mamoroom/echo-mvc/app/lib/util"
	"github.com/mamoroom/echo-mvc/app/view/res_redirect"
)

type ShareTmplParam struct {
	//meta
	Title       string
	Description string

	//fb:og
	OgType        string
	OgTitle       string
	OgDescription string
	OgSiteName    string
	OgUrl         string
	OgImage       string
	OgImageWidth  string
	OgImageHeight string

	//tw card
	TwCard        string
	TwTitle       string
	TwDescription string
	TwImage       string
	TwUrl         string

	// others
	Canonical    string
	ShortcutIcon string
}

var OG_SITE_NAME = ""
var OG_IMAGE_WIDTH = "2400"
var OG_IMAGE_HEIGHT = "1260"
var OG_IMAGE_PATH_PREFIX = "/assets/ogp/"
var OG_IMAGE_PATH_SUFFIX = "/ogp.jpg"
var OG_TYPE = "website"
var TW_CARD = "summary_large_image"
var SHORTCUT_IMAGE_PATH_SUFFIX = "favicon.ico"

func ShareHandler(c echo.Context) error {

	_, err := share.DecryptHash2Param(c.Param("hash"))
	//p, err := share.DecryptHash2Param(c.Param("hash"))
	if err != nil {
		return res_redirect.RedirectToIndexHtml(c)
	}

	// bot以外は適切なURLへと遷移
	if !util.CheckCralwerUaValidation(c.Request().UserAgent()) {
		q := map[string]string{}
		return res_redirect.RedirectToNazoDetailHtml(c, q)
	}

	// ogp
	share_data_json, err := share_data.GetShareData("dynamic")
	if err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"type": "GetShareJsonDataError",
		}).Error(err.Error())
		return res_redirect.RedirectToIndexHtml(c)
	}

	return c.Render(http.StatusOK, "share", share_data_json)

}

func ShareStaticHandler(c echo.Context) error {

	// bot以外は適切なURLへと遷移
	if !util.CheckCralwerUaValidation(c.Request().UserAgent()) {
		return res_redirect.RedirectToIndexHtml(c)
	}

	// ogp
	data, err := share_data.GetShareData("static")
	if err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"type": "GetShareJsonDataError",
		}).Error(err.Error())
		return res_redirect.RedirectToIndexHtml(c)
	}

	return c.Render(http.StatusOK, "share", data)
}
