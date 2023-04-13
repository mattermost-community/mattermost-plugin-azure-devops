package plugin

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-azure-devops/server/constants"
	"github.com/mattermost/mattermost-plugin-azure-devops/server/serializers"
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

	oAuthState := fmt.Sprintf("%s_%s", model.NewId()[0:15], mattermostUserID)
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

	if isConnected := p.MattermostUserAlreadyConnected(mattermostUserID); isConnected {
		p.CloseBrowserWindowWithHTTPResponse(w)
		if _, DMErr := p.DM(mattermostUserID, constants.MattermostUserAlreadyConnected, false); DMErr != nil {
			p.handleError(w, r, &serializers.Error{Code: http.StatusInternalServerError, Message: DMErr.Error()})
			return
		}
		p.handleError(w, r, &serializers.Error{Code: http.StatusBadRequest, Message: constants.MattermostUserAlreadyConnected})
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

	mattermostUserID := r.Header.Get(constants.HeaderMattermostUserID)
	if err := p.GenerateOAuthToken(code, state, mattermostUserID); err != nil {
		if strings.Contains(err.Error(), "already connected") {
			p.API.LogError(constants.UnableToCompleteOAuth, "Error", constants.ErrorMessageAzureDevopsAccountAlreadyConnected)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		p.API.LogError(constants.UnableToCompleteOAuth, "Error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.CloseBrowserWindowWithHTTPResponse(w)
}

// GenerateOAuthToken generates OAuth token after successful authorization
func (p *Plugin) GenerateOAuthToken(code, state, authenticatedMattermostUserID string) error {
	mattermostUserID := strings.Split(state, "_")[1]

	if mattermostUserID != authenticatedMattermostUserID {
		return errors.New("failed to complete oAuth, mattermost user is not authenticated")
	}

	if err := p.Store.VerifyOAuthState(mattermostUserID, state); err != nil {
		return errors.Wrap(err, "failed to verify oAuth state")
	}

	oauthTokenFormValues := url.Values{
		"client_assertion_type": {constants.ClientAssertionType},
		"client_assertion":      {p.getConfiguration().AzureDevopsOAuthClientSecret},
		"grant_type":            {constants.GrantType},
		"assertion":             {code},
		"redirect_uri":          {fmt.Sprintf("%s%s%s", p.GetSiteURL(), p.GetPluginURLPath(), constants.PathOAuthCallback)},
	}

	if err := p.GenerateAndStoreOAuthToken(mattermostUserID, oauthTokenFormValues, false); err != nil {
		return err
	}

	p.API.PublishWebSocketEvent(
		constants.WSEventConnect,
		nil,
		&model.WebsocketBroadcast{UserId: mattermostUserID},
	)

	if _, err := p.DM(mattermostUserID, fmt.Sprintf("%s\n\n%s", constants.UserConnected, constants.HelpText), false); err != nil {
		return err
	}

	return nil
}

// RefreshOAuthToken refreshes OAuth token
func (p *Plugin) RefreshOAuthToken(mattermostUserID, refreshToken string) error {
	decodedRefreshToken, err := p.Decode(refreshToken)
	if err != nil {
		if _, DMErr := p.DM(mattermostUserID, constants.GenericErrorMessage, false); DMErr != nil {
			return DMErr
		}
		return err
	}

	decryptedRefreshToken, err := p.Decrypt(decodedRefreshToken, []byte(p.getConfiguration().EncryptionSecret))
	if err != nil {
		if _, DMErr := p.DM(mattermostUserID, constants.GenericErrorMessage, false); DMErr != nil {
			return DMErr
		}
		return err
	}

	oauthTokenFormValues := url.Values{
		"client_assertion_type": {constants.ClientAssertionType},
		"client_assertion":      {p.getConfiguration().AzureDevopsOAuthClientSecret},
		"grant_type":            {constants.GrantTypeRefresh},
		"assertion":             {string(decryptedRefreshToken)},
		"redirect_uri":          {fmt.Sprintf("%s%s%s", p.GetSiteURL(), p.GetPluginURLPath(), constants.PathOAuthCallback)},
	}

	return p.GenerateAndStoreOAuthToken(mattermostUserID, oauthTokenFormValues, true)
}

// GenerateAndStoreOAuthToken generates and stores OAuth token
func (p *Plugin) GenerateAndStoreOAuthToken(mattermostUserID string, oauthTokenFormValues url.Values, isTokenRefreshRequest bool) error {
	successResponse, _, err := p.Client.GenerateOAuthToken(oauthTokenFormValues)
	if err != nil {
		if _, DMErr := p.DM(mattermostUserID, constants.GenericErrorMessage, false); DMErr != nil {
			return DMErr
		}
		return errors.Wrap(err, "failed to generate oAuth token")
	}

	userProfile, _, err := p.Client.GetUserProfile(constants.CurrentAzureDevopsUserProfileID, successResponse.AccessToken)
	if err != nil {
		if _, DMErr := p.DM(mattermostUserID, constants.GenericErrorMessage, false); DMErr != nil {
			return DMErr
		}
		return errors.Wrap(err, "failed to fetch user profile")
	}

	azureDevopsUser, err := p.Store.LoadAzureDevopsUserDetails(userProfile.ID)
	if err != nil {
		if _, DMErr := p.DM(mattermostUserID, constants.GenericErrorMessage, false); DMErr != nil {
			return DMErr
		}
		return errors.Wrap(err, "failed to get the user details")
	}

	if !isTokenRefreshRequest && azureDevopsUser.AccessToken != "" {
		if _, DMErr := p.DM(mattermostUserID, fmt.Sprintf(constants.ErrorMessageAzureDevopsAccountAlreadyConnected, userProfile.Email), false); DMErr != nil {
			return errors.Wrap(err, "failed to DM user")
		}

		return fmt.Errorf(constants.ErrorMessageAzureDevopsAccountAlreadyConnected, userProfile.Email)
	}

	encryptedAccessToken, err := p.Encrypt([]byte(successResponse.AccessToken), []byte(p.getConfiguration().EncryptionSecret))
	if err != nil {
		return err
	}

	encryptedRefreshToken, err := p.Encrypt([]byte(successResponse.RefreshToken), []byte(p.getConfiguration().EncryptionSecret))
	if err != nil {
		return err
	}

	tokenExpiryDurationInSeconds, err := strconv.Atoi(successResponse.ExpiresIn)
	if err != nil {
		if _, DMErr := p.DM(mattermostUserID, constants.GenericErrorMessage, false); DMErr != nil {
			return DMErr
		}
		return err
	}

	user := serializers.User{
		MattermostUserID: mattermostUserID,
		AccessToken:      p.Encode(encryptedAccessToken),
		RefreshToken:     p.Encode(encryptedRefreshToken),
		ExpiresAt:        time.Now().UTC().Add(time.Second * time.Duration(tokenExpiryDurationInSeconds)).Unix(),
		UserProfile:      *userProfile,
	}

	if err := p.Store.StoreAzureDevopsUserDetailsWithMattermostUserID(&user); err != nil {
		return err
	}

	return nil
}

// IsAccessTokenExpired checks if a user's access token is expired
func (p *Plugin) IsAccessTokenExpired(mattermostUserID string) (bool, string) {
	azureDevopsUserID, err := p.Store.LoadAzureDevopsUserIDFromMattermostUser(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorLoadingUserData, "Error", err.Error())
		return false, ""
	}

	user, err := p.Store.LoadAzureDevopsUserDetails(azureDevopsUserID)
	if err != nil {
		p.API.LogError(constants.ErrorLoadingUserData, "Error", err.Error())
		return false, ""
	}

	// Consider some buffer for comparing expiry time
	localExpiryTime := time.Unix(user.ExpiresAt, 0).Local()
	if user.AccessToken != "" && time.Until(localExpiryTime) <= time.Minute*constants.TokenExpiryTimeBufferInMinutes {
		return true, user.RefreshToken
	}

	return false, ""
}

// MattermostUserAlreadyConnected checks if a user is already connected
func (p *Plugin) MattermostUserAlreadyConnected(mattermostUserID string) bool {
	azureDevopsUserID, err := p.Store.LoadAzureDevopsUserIDFromMattermostUser(mattermostUserID)
	if err != nil {
		p.API.LogError(constants.ErrorLoadingUserData, "Error", err.Error())
		return false
	}

	user, err := p.Store.LoadAzureDevopsUserDetails(azureDevopsUserID)
	if err != nil {
		p.API.LogError(constants.UnableToCheckIfAlreadyConnected, "Error", err.Error())
		return false
	}

	if user.AccessToken != "" {
		return true
	}

	return false
}

// CloseBrowserWindowWithHTTPResponse closes the browser window
func (p *Plugin) CloseBrowserWindowWithHTTPResponse(w http.ResponseWriter) {
	html := `
	<!DOCTYPE html>
	<html>
		<head>
			<script>
				window.open('','_parent','');
				window.close();
			</script>
		</head>
		<body>
			<p>Completed connecting to Azure DevOps. Please close this window.</p>
		</body>
	</html>
	`

	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(html)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
