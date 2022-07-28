package plugin

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/Brightscout/mattermost-plugin-azure-devops/server/store"
)

type OAuthConfig struct {
	appID        string
	clientSecret string
	authURL      string
	redirectURI  string
	responseType string
	scope        string
}

// OAuthConfig initialize OAuth configs
func (p *Plugin) OAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		appID:        p.getConfiguration().AzureDevopsOAuthAppID,
		clientSecret: p.getConfiguration().AzureDevopsOAuthClientSecret,
		authURL:      fmt.Sprintf("%s%s", constants.BaseOauthURL, constants.PathAuth),
		redirectURI:  fmt.Sprintf("%s%s%s", p.GetSiteURL(), p.GetPluginURLPath(), constants.PathOAuthCallback),
		responseType: constants.ResponseType,
		scope:        constants.Scopes, // these scopes must be set in the OAuth app registered with the Azure portal
	}
}

// GenerateOAuthConnectURL generates URL for Azure OAuth authorization
func (p *Plugin) GenerateOAuthConnectURL(mattermostUserID string) string {
	oAuthConfig := p.OAuthConfig()

	oAuthState := fmt.Sprintf("%v_%v", model.NewId()[0:15], mattermostUserID)
	if err := p.Store.StoreOAuthState(mattermostUserID, oAuthState); err != nil {
		p.API.LogError(fmt.Sprintf(constants.UnableToStoreOauthState, mattermostUserID), "Error", err.Error())
	}

	var stringBuilder strings.Builder
	stringBuilder.WriteString(oAuthConfig.authURL)
	parameterisedURL := url.Values{
		"response_type": {oAuthConfig.responseType},
		"client_id":     {oAuthConfig.appID},
		"redirect_uri":  {oAuthConfig.redirectURI},
		"scope":         {oAuthConfig.scope},
		"state":         {oAuthState},
	}

	if strings.Contains(oAuthConfig.authURL, "?") {
		stringBuilder.WriteByte('&')
	} else {
		stringBuilder.WriteByte('?')
	}
	stringBuilder.WriteString(parameterisedURL.Encode())
	return stringBuilder.String()
}

// OAuthConnect redirects to the OAuth authorization URL
func (p *Plugin) OAuthConnect(w http.ResponseWriter, r *http.Request) {
	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
	// TODO: use checkAuth middleware for this
	if mattermostUserID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	if isConnected := p.UserAlreadyConnected(mattermostUserID); isConnected {
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

	if len(strings.Split(state, "_")) != 2 || strings.Split(state, "_")[1] == "" {
		http.Error(w, constants.InvalidAuthState, http.StatusBadRequest)
		return
	}

	if err := p.GenerateOAuthToken(code, state); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.closeBrowserWindowWithHTTPResponse(w)
}

// GenerateOAuthToken generates OAuth token after successful authorization
func (p *Plugin) GenerateOAuthToken(code, state string) error {
	mattermostUserID := strings.Split(state, "_")[1]

	if err := p.Store.VerifyOAuthState(mattermostUserID, state); err != nil {
		p.DM(mattermostUserID, constants.GenericErrorMessage)
		return errors.Wrap(err, err.Error())
	}

	oauthTokenFormValues := url.Values{
		"client_assertion_type": {constants.ClientAssertionType},
		"client_assertion":      {p.getConfiguration().AzureDevopsOAuthClientSecret},
		"grant_type":            {constants.GrantType},
		"assertion":             {code},
		"redirect_uri":          {fmt.Sprintf("%s%s%s", p.GetSiteURL(), p.GetPluginURLPath(), constants.PathOAuthCallback)},
	}

	successResponse, err := p.Client.GenerateOAuthToken(oauthTokenFormValues.Encode())
	if err != nil {
		p.DM(mattermostUserID, constants.GenericErrorMessage)
		return errors.Wrap(err, err.Error())
	}

	encryptedAccessToken, err := p.encrypt([]byte(successResponse.AccessToken), []byte(p.getConfiguration().EncryptionSecret))
	if err != nil {
		return err
	}

	encryptedRefreshToken, err := p.encrypt([]byte(successResponse.RefreshToken), []byte(p.getConfiguration().EncryptionSecret))
	if err != nil {
		return err
	}

	user := store.User{
		MattermostUserID: mattermostUserID,
		AccessToken:      p.encode(encryptedAccessToken),
		RefreshToken:     p.encode(encryptedRefreshToken),
		ExpiresIn:        successResponse.ExpiresIn,
	}

	p.Store.StoreUser(&user)

	fmt.Printf("%+v\n", successResponse) // TODO: remove later

	p.DM(mattermostUserID, fmt.Sprintf("%s\n\n%s", constants.UserConnected, constants.HelpText))

	return nil
}

// UserAlreadyConnected checks if a user is already connected
func (p *Plugin) UserAlreadyConnected(mattermostUserID string) bool {
	user, err := p.Store.LoadUser(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.UnableToCheckIfAlreadyConnected, "Error", err.Error())
		return false
	}

	if user.AccessToken != "" {
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
