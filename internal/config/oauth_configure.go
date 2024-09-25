package config

import (
	"github.com/cobaltbase/cobaltbase/internal/ct"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/azuread"
	"github.com/markbates/goth/providers/battlenet"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/box"
	"github.com/markbates/goth/providers/dailymotion"
	"github.com/markbates/goth/providers/deezer"
	"github.com/markbates/goth/providers/digitalocean"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/eveonline"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/fitbit"
	"github.com/markbates/goth/providers/gitea"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/heroku"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/intercom"
	"github.com/markbates/goth/providers/kakao"
	"github.com/markbates/goth/providers/lastfm"
	"github.com/markbates/goth/providers/line"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/microsoftonline"
	"github.com/markbates/goth/providers/naver"
	"github.com/markbates/goth/providers/nextcloud"
	"github.com/markbates/goth/providers/okta"
	"github.com/markbates/goth/providers/onedrive"
	"github.com/markbates/goth/providers/patreon"
	"github.com/markbates/goth/providers/paypal"
	"github.com/markbates/goth/providers/salesforce"
	"github.com/markbates/goth/providers/seatalk"
	"github.com/markbates/goth/providers/shopify"
	"github.com/markbates/goth/providers/slack"
	"github.com/markbates/goth/providers/soundcloud"
	"github.com/markbates/goth/providers/spotify"
	"github.com/markbates/goth/providers/steam"
	"github.com/markbates/goth/providers/strava"
	"github.com/markbates/goth/providers/stripe"
	"github.com/markbates/goth/providers/tiktok"
	"github.com/markbates/goth/providers/tumblr"
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/typetalk"
	"github.com/markbates/goth/providers/uber"
	"github.com/markbates/goth/providers/vk"
	"github.com/markbates/goth/providers/wepay"
	"github.com/markbates/goth/providers/xero"
	"github.com/markbates/goth/providers/yahoo"
	"github.com/markbates/goth/providers/yammer"
	"github.com/markbates/goth/providers/yandex"
	"log"
)

func FetchAllOauthConfigs() {
	var allConfigs []ct.OauthConfig

	err := DB.Find(&allConfigs).Error
	if err == nil {
		log.Println("Fetched All Oauth Configs")
	}
	var providers []goth.Provider
	callbackURL := "http://localhost:3000/api/auth/oauth/callback"
	for _, config := range allConfigs {
		var provider goth.Provider

		switch config.Provider {
		case "amazon":
			provider = amazon.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=amazon")
		case "apple":
			provider = apple.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=apple", nil, apple.ScopeName, apple.ScopeEmail)
		case "auth0":
			provider = auth0.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=auth0", "your-auth0-domain.auth0.com")
		case "azuread":
			provider = azuread.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=azuread", nil)
		case "battlenet":
			provider = battlenet.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=battlenet")
		case "bitbucket":
			provider = bitbucket.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=bitbucket")
		case "box":
			provider = box.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=box")
		case "dailymotion":
			provider = dailymotion.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=dailymotion")
		case "deezer":
			provider = deezer.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=deezer")
		case "digitalocean":
			provider = digitalocean.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=digitalocean")
		case "discord":
			provider = discord.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=discord")
		case "dropbox":
			provider = dropbox.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=dropbox")
		case "eveonline":
			provider = eveonline.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=eveonline")
		case "facebook":
			provider = facebook.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=facebook")
		case "fitbit":
			provider = fitbit.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=fitbit")
		case "gitea":
			provider = gitea.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=gitea")
		case "github":
			provider = github.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=github")
		case "gitlab":
			provider = gitlab.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=gitlab")
		case "google":
			provider = google.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=google")
		case "gplus":
			provider = gplus.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=gplus")
		case "heroku":
			provider = heroku.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=heroku")
		case "instagram":
			provider = instagram.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=instagram")
		case "intercom":
			provider = intercom.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=intercom")
		case "kakao":
			provider = kakao.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=kakao")
		case "lastfm":
			provider = lastfm.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=lastfm")
		case "line":
			provider = line.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=line")
		case "linkedin":
			provider = linkedin.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=linkedin")
		case "microsoftonline":
			provider = microsoftonline.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=microsoftonline")
		case "naver":
			provider = naver.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=naver")
		case "nextcloud":
			provider = nextcloud.NewCustomisedDNS(config.ClientID, config.ClientSecret, callbackURL+"?provider=nextcloud", "your-nextcloud-instance.com")
		case "okta":
			provider = okta.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=okta", "your-okta-instance.okta.com")
		case "onedrive":
			provider = onedrive.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=onedrive")
		//case "openid-connect":
		//	provider = openidConnect.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=openid-connect", "https://your-openid-connect-provider.com")
		case "patreon":
			provider = patreon.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=patreon")
		case "paypal":
			provider = paypal.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=paypal")
		case "salesforce":
			provider = salesforce.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=salesforce")
		case "seatalk":
			provider = seatalk.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=seatalk")
		case "shopify":
			provider = shopify.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=shopify")
		case "slack":
			provider = slack.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=slack")
		case "soundcloud":
			provider = soundcloud.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=soundcloud")
		case "spotify":
			provider = spotify.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=spotify")
		case "steam":
			provider = steam.New(config.ClientID, callbackURL+"?provider=steam")
		case "strava":
			provider = strava.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=strava")
		case "stripe":
			provider = stripe.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=stripe")
		case "tiktok":
			provider = tiktok.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=tiktok")
		case "tumblr":
			provider = tumblr.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=tumblr")
		case "twitch":
			provider = twitch.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=twitch")
		case "twitter":
			provider = twitter.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=twitter")
		case "typetalk":
			provider = typetalk.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=typetalk")
		case "uber":
			provider = uber.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=uber")
		case "vk":
			provider = vk.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=vk")
		case "wepay":
			provider = wepay.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=wepay")
		case "xero":
			provider = xero.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=xero")
		case "yahoo":
			provider = yahoo.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=yahoo")
		case "yammer":
			provider = yammer.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=yammer")
		case "yandex":
			provider = yandex.New(config.ClientID, config.ClientSecret, callbackURL+"?provider=yandex")
		default:
			log.Printf("Warning: Unknown provider %s", config.Provider)
			continue
		}

		providers = append(providers, provider)
	}
	goth.UseProviders(providers...)
}
