package serializers

type User struct {
	MattermostUserID string `json:"mattermostUserID"`
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	ExpiresAt        int64  `json:"expiresAt"`
	UserProfile
}
