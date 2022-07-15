package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/serializers"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/store"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

type OAuthConfig struct {
	appId        string
	clientSecret string
	authURL      string
	tokenURL     string
	redirectURI  string
	responseType string
	scope        string
}

// OAuthConfig initialize OAuth configs
func (p *Plugin) OAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		appId:        p.getConfiguration().AzureDevopsOAuthAppID,
		clientSecret: p.getConfiguration().AzureDevopsOAuthClientSecret,
		authURL:      p.getConfiguration().AzureDevopsOAuthAuthorizationURL,
		tokenURL:     p.getConfiguration().AzureDevopsOAuthTokenURL,
		redirectURI:  p.getConfiguration().AzureDevopsOAuthCallbackURL,
		responseType: constants.ResponseType,
		scope:        constants.Scopes, // these scopes must be set in the OAuth app registered with the Azure portal
	}
}

// GenerateOAuthConnectURL generates URL for Azure OAuth authorization
func (p *Plugin) GenerateOAuthConnectURL(mattermostUserID string) string {
	oAuthConfig := p.OAuthConfig()

	var buf bytes.Buffer
	buf.WriteString(oAuthConfig.authURL)
	v := url.Values{
		"response_type": {oAuthConfig.responseType},
		"client_id":     {oAuthConfig.appId},
		"redirect_uri":  {oAuthConfig.redirectURI},
		"scope":         {oAuthConfig.scope},
		"state":         {fmt.Sprintf("%v/%v", model.NewId()[0:15], mattermostUserID)},
	}

	if strings.Contains(oAuthConfig.authURL, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	buf.WriteString(v.Encode())
	return buf.String()
}

// OAuthConnect redirects to the OAuth authorization URL
func (p *Plugin) OAuthConnect(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get("Mattermost-User-ID")
	if mattermostUserID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	channelID := r.URL.Query().Get("channel_id")
	if channelID == "" {
		http.Error(w, "missing channel ID", http.StatusBadRequest)
		return
	}

	if isConnected := p.UserAlreadyConnected(mattermostUserID, channelID); isConnected {
		p.closeBrowserWindowWithHTTPResponse(w)
		p.DM(mattermostUserID, constants.UserAlreadyConnected)
		return
	}

	redirectURL := p.GenerateOAuthConnectURL(mattermostUserID)

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// OAuthComplete captures the redirection request made by the OAuth authorization
func (p *Plugin) OAuthComplete(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing authorization code", http.StatusBadRequest)
		return
	}

	state := r.URL.Query().Get("state")
	if state == "" {
		http.Error(w, "missing authorization state", http.StatusBadRequest)
		return
	}

	err := p.GenerateOAuthToken(code, state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.closeBrowserWindowWithHTTPResponse(w)
}

// GenerateOAuthToken generates OAuth token after successful authorization
func (p *Plugin) GenerateOAuthToken(code string, state string) error {
	if code == "" || state == "" {
		return errors.New("missing code or state")
	}

	mattermostUserID := strings.Split(state, "/")[1]

	oAuthConfig := p.OAuthConfig()

	form := url.Values{
		"client_assertion_type": {constants.ClientAssertionType},
		"client_assertion":      {p.getConfiguration().AzureDevopsOAuthClientSecret},
		"grant_type":            {constants.GrantType},
		"assertion":             {code},
		"redirect_uri":          {p.getConfiguration().AzureDevopsOAuthCallbackURL},
	}

	// Create a HTTP post request
	r, err := http.NewRequest("POST", oAuthConfig.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return errors.Wrap(err, err.Error())
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}

	if res.StatusCode != http.StatusOK {
		errResp := serializers.OAuthErrorResponse{}

		err = json.Unmarshal(resBody, &errResp)
		if err != nil {
			return errors.Wrap(err, err.Error())
		}

		p.DM(mattermostUserID, constants.GenericErrorMessage)
		return errors.Wrap(errors.New(errResp.ErrorMessage), errResp.ErrorDescription)
	}

	successResp := serializers.OAuthSuccessResponse{}

	err = json.Unmarshal(resBody, &successResp)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}

	encryptedAccessToken, err := p.encrypt([]byte(successResp.AccessToken), []byte(p.getConfiguration().EncryptionSecret))
	if err != nil {
		return err
	}

	encryptedRefreshToken, err := p.encrypt([]byte(successResp.RefreshToken), []byte(p.getConfiguration().EncryptionSecret))
	if err != nil {
		return err
	}

	user := store.User{
		MattermostUserID: mattermostUserID,
		AccessToken:      p.encode(encryptedAccessToken),
		RefreshToken:     p.encode(encryptedRefreshToken),
		ExpiresIn:        successResp.ExpiresIn,
	}

	p.Store.StoreUser(&user)

	fmt.Printf("%+v\n", successResp) // TODO: remove later

	p.DM(mattermostUserID, fmt.Sprintf("%s\n\n%s", constants.UserConnected, constants.HelpText))

	return nil
}

// UserAlreadyConnected checks if a user is already connected
func (p *Plugin) UserAlreadyConnected(mattermostUserID string, channelID string) bool {
	user, err := p.Store.LoadUser(mattermostUserID)
	if err != nil {
		errors.Wrap(err, err.Error())
		return false
	}

	if user.AccessToken != "" {
		abc, _ := p.decode(user.AccessToken)
		aa, _ := p.decrypt([]byte(abc), []byte(p.getConfiguration().EncryptionSecret))

		fmt.Printf("%+s token\n", string(aa))
		return true
	}

	return false
}

// closeBrowserWindowWithHTTPResponse closes the browser window
func (p *Plugin) closeBrowserWindowWithHTTPResponse(w http.ResponseWriter) {
	html := `
	<!DOCTYPE html>
	<html>
		<head>
			<script>
				window.close();
			</script>
		</head>
		<body>
			<p>Completed connecting to Azure Devops. Please close this window.</p>
		</body>
	</html>
	`

	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(html)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
