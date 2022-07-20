package plugin

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Brightscout/mattermost-plugin-azure-devops/server/constants"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

// sendEphemeralPostForCommand sends an ephermal message
func (p *Plugin) sendEphemeralPostForCommand(args *model.CommandArgs, text string) (*model.CommandResponse, *model.AppError) {
	post := &model.Post{
		UserId:    p.botUserID,
		ChannelId: args.ChannelId,
		Message:   text,
	}
	_ = p.API.SendEphemeralPost(args.UserId, post)

	return &model.CommandResponse{}, nil
}

// DM posts a simple Direct Message to the specified user
func (p *Plugin) DM(mattermostUserID, format string, args ...interface{}) (string, error) {
	channel, err := p.API.GetDirectChannel(mattermostUserID, p.botUserID)
	if err != nil {
		p.API.LogError("Couldn't get bot's DM channel", "user_id", mattermostUserID, "error", err.Error())
		return "", err
	}
	post := &model.Post{
		ChannelId: channel.Id,
		UserId:    p.botUserID,
		Message:   fmt.Sprintf(format, args...),
	}
	sentPost, err := p.API.CreatePost(post)
	if err != nil {
		p.API.LogError("error occurred while creating post", "error", err.Error())
		return "", err
	}
	return sentPost.Id, nil
}

// encode encodes bytes into base64 string
func (p *Plugin) encode(encrypted []byte) string {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(encrypted)))
	base64.URLEncoding.Encode(encoded, encrypted)
	return string(encoded)
}

// decode decodes a base64 string into bytes
func (p *Plugin) decode(encoded string) ([]byte, error) {
	decoded := make([]byte, base64.URLEncoding.DecodedLen(len(encoded)))
	n, err := base64.URLEncoding.Decode(decoded, []byte(encoded))
	if err != nil {
		return nil, err
	}
	return decoded[:n], nil
}

// encrypt used for generating encrypted bytes
func (p *Plugin) encrypt(plain, secret []byte) ([]byte, error) {
	if len(secret) == 0 {
		return plain, nil
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	sealed := aesgcm.Seal(nil, nonce, plain, nil)
	return append(nonce, sealed...), nil
}

// decrypt used for generating decrypted bytes
func (p *Plugin) decrypt(encrypted, secret []byte) ([]byte, error) {
	if len(secret) == 0 {
		return encrypted, nil
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(encrypted) < nonceSize {
		return nil, errors.New("token too short")
	}

	nonce, encrypted := encrypted[:nonceSize], encrypted[nonceSize:]
	plain, err := aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, err
	}

	return plain, nil
}

func (p *Plugin) GetSiteURL() string {
	return p.getConfiguration().MattermostSiteURL
}

func (p *Plugin) GetPluginURLPath() string {
	return "/plugins/" + constants.PluginID + "/api/v1"
}

func (p *Plugin) GetPluginURL() string {
	return strings.TrimRight(p.GetSiteURL(), "/") + p.GetPluginURLPath()
}

// AddAuthorization function to add authorization to a request.
func (p *Plugin) AddAuthorization(r *http.Request, mattermostUserID string) error {
	user, err := p.Store.LoadUser(mattermostUserID)
	if err != nil {
		return err
	}
	decodedAccessToken, err := p.decode(user.AccessToken)
	if err != nil {
		return err
	}
	decryptedAccessToken, err := p.decrypt(decodedAccessToken, []byte(p.getConfiguration().EncryptionSecret))
	if err != nil {
		return err
	}
	r.Header.Add(constants.Authorization, fmt.Sprintf(constants.Bearer, string(decryptedAccessToken)))
	return nil
}

// TODO: WIP.
// StringToInt function to convert string to int.
// func StringToInt(str string) int {
// 	if str == "" {
// 		return 0
// 	}
// 	val, err := strconv.ParseInt(str, 10, 64)
// 	if err != nil {
// 		return 0
// 	}
// 	return int(val)
// }
