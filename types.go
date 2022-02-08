package gotrue

type Provider string

const (
	ProviderApple     Provider = "apple"
	ProviderAzure     Provider = "azure"
	ProviderBitBucket Provider = "bitbucket"
	ProviderDiscord   Provider = "discord"
	ProviderFacebook  Provider = "facebook"
	ProviderGithub    Provider = "github"
	ProviderGitlab    Provider = "gitlab"
	ProviderGoogle    Provider = "google"
	ProviderSlack     Provider = "slack"
	ProviderSpotify   Provider = "spotify"
	ProviderTwitch    Provider = "twitch"
	ProviderTwitter   Provider = "twitter"
)

type AuthChangeEvent string

const (
	PasswordRecoveryEvent AuthChangeEvent = "PASSWORD_RECOVERY"
	SignedInEvent         AuthChangeEvent = "SIGNED_IN"
	SignedOutEvent        AuthChangeEvent = "SIGNED_OUT"
	TokenRefreshedEvent   AuthChangeEvent = "TOKEN_REFRESHED"
	UserDeletedEvent      AuthChangeEvent = "USER_DELETED"
	UserUpdatedEvent      AuthChangeEvent = "USER_UPDATED"
)
