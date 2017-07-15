package oauth

import (
	"github.com/echo-contrib/sessions"
	"github.com/labstack/echo"
	"github.com/prtytokyo/goth"

	"github.com/mamoroom/echo-mvc/app/lib/util"

	"errors"
	_ "fmt"
)

/*
GetAuthURL starts the authentication process with the requested provided.
It will return a URL that should be used to send users to.

It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".

I would recommend using the BeginAuthHandler instead of doing all of these steps
yourself, but that's entirely up to you.
*/

func GetAuthURL(c echo.Context) (string, error) {

	providerName, err := GetProviderName(c)
	if err != nil {
		return providerName, err
	}
	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	state_str := set_state(c)
	if state_str == "" {
		return "", errors.New("Could not set state string")
	}
	sess, err := provider.BeginAuth(state_str)
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	err = store_oauth_cookie_in_session(providerName, sess.Marshal(), c)

	if err != nil {
		return "", err
	}

	return url, err
}

/*
CompleteUserAuth does what it says on the tin. It completes the authentication
process and fetches all of the basic information about the user from the provider.

It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".

See https://github.com/prtytokyo/goth/examples/main.go to see this in action.
*/
func CompleteUserAuth(c echo.Context) (goth.User, error) {

	providerName, err := GetProviderName(c)
	if err != nil {
		return goth.User{}, err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}

	value, err := get_from_session(providerName, c)
	if err != nil {
		return goth.User{}, err
	}

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		return goth.User{}, err
	}

	user, err := provider.FetchUser(sess)
	if err == nil {
		// user can be found with existing session data
		return user, err
	}

	// get new token and retry fetch
	_, err = sess.Authorize(provider, c.QueryParams())
	if err != nil {
		return goth.User{}, err
	}

	err = store_oauth_cookie_in_session(providerName, sess.Marshal(), c)

	if err != nil {
		return goth.User{}, err
	}

	return provider.FetchUser(sess)
}

// Logout invalidates a user session.
func Logout(c echo.Context) error {

	providerName, err := GetProviderName(c)
	if err != nil {
		return err
	}
	session := sessions.Default(c)
	session.Delete(providerName)
	session.Save()
	return nil
}

func GetProviderName(c echo.Context) (string, error) {
	provider := c.QueryParam("provider")
	if provider == "" {
		provider = c.Param("provider")
	}
	if provider == "" {
		provider = c.QueryParam("provider")
	}
	if provider == "" {
		return provider, errors.New("You must select a provider")
	}
	return provider, nil
}

func ValidateState(c echo.Context) error {
	state_in_session, err := get_from_session("state", c)
	if err != nil {
		return err
	}
	if state_in_session != get_state(c) {
		return errors.New("Invalid token was used")
	}
	return nil
}

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
func set_state(c echo.Context) string {
	/*state := c.QueryParam("state")
	if len(state) > 0 {
		return state
	}*/
	token := util.GenerateBase32RandomKey(20)
	err := store_oauth_cookie_in_session("state", token, c)
	if err != nil {
		return ""
	}
	return token
}

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12

func get_state(c echo.Context) string {
	return c.QueryParam("state")
}

func store_oauth_cookie_in_session(key string, value string, c echo.Context) error {
	session := sessions.Default(c)
	session.Options(sessions.Options{
		Path:     conf.Oauth.Cookie.Path,
		Domain:   conf.Oauth.Cookie.Domain,
		MaxAge:   conf.Oauth.Cookie.MaxAge,
		Secure:   conf.Oauth.Cookie.Secure,
		HttpOnly: conf.Oauth.Cookie.HttpOnly,
	})
	return store_in_session(key, value, session)
}

func store_in_session(key string, value interface{}, session sessions.Session) error {
	session.Set(key, value)
	return session.Save()
}

func get_from_session(key string, c echo.Context) (string, error) {
	session := sessions.Default(c)
	value := session.Get(key)
	if value == nil {
		return "", errors.New("could not find a matching session for this request")
	}

	return value.(string), nil
}
