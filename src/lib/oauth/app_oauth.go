package oauth

/*
import (
	"github.com/labstack/echo"

	_ "github.com/mamoroom/echo-mvc/src/lib/util"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"

	"errors"
	_ "fmt"
)

func GetAppAuthURL(c echo.Context) (string, error) {

	providerName, err := get_provider_name(c)
	if err != nil {
		return "", err
	}
	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	state_str := set_state(c)
	if state_str == "" {
		return "", errors.New("Could not set state string")
	}

	pkce_params, err := get_pkce_params(c)
	if err != nil {
		return "", err
	}

	provider.(*openidConnect.Provider).SetAuthCodeOptions(pkce_params)
	sess, err := provider.BeginAuth(state_str)
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	err = store_in_session(providerName, sess.Marshal(), c)

	if err != nil {
		return "", err
	}

	return url, err
}

func get_pkce_params(c echo.Context) (map[string]string, error) {
	params := map[string]string{}

	code_challenge_method := c.QueryParam("code_challenge_method")
	if code_challenge_method != "plain" && code_challenge_method != "S256" {
		return params, errors.New("No paramater found: code_challenge_method")
	}
	params["code_challenge_method"] = code_challenge_method

	code_challenge := c.QueryParam("code_challenge")
	if code_challenge == "" {
		return params, errors.New("no paramater found: code_challenge")
	}
	params["code_challenge"] = code_challenge

	return params, nil
}
*/
