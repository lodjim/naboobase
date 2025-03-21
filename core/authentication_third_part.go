package core

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lodjim/naboobase/models"
	"github.com/lodjim/naboobase/utils"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
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
	"github.com/markbates/goth/providers/mastodon"
	"github.com/markbates/goth/providers/meetup"
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
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/twitterv2"
	"github.com/markbates/goth/providers/typetalk"
	"github.com/markbates/goth/providers/uber"
	"github.com/markbates/goth/providers/vk"
	"github.com/markbates/goth/providers/wecom"
	"github.com/markbates/goth/providers/wepay"
	"github.com/markbates/goth/providers/xero"
	"github.com/markbates/goth/providers/yahoo"
	"github.com/markbates/goth/providers/yammer"
	"github.com/markbates/goth/providers/yandex"
	"github.com/markbates/goth/providers/zoom"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ThirdPartAuthenticator struct {
}

func (thirdPartAuthenticator *ThirdPartAuthenticator) Init(db MongoDBconnector) []Endpoint {
	goth.UseProviders(
		// Use twitterv2 instead of twitter if you only have access to the Essential API Level
		// the twitter provider uses a v1.1 API that is not available to the Essential Level
		twitterv2.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:1555/oauth/twitterv2/callback"),
		// If you'd like to use authenticate instead of authorize in TwitterV2 provider, use this instead.
		// twitterv2.NewAuthenticate(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:1555/auth/twitterv2/callback"),

		twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:1555/oauth/twitter/callback"),
		// If you'd like to use authenticate instead of authorize in Twitter provider, use this instead.
		// twitter.NewAuthenticate(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:1555/auth/twitter/callback"),

		tiktok.New(os.Getenv("TIKTOK_KEY"), os.Getenv("TIKTOK_SECRET"), "http://localhost:1555/oauth/tiktok/callback"),
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost:1555/oauth/facebook/callback"),
		fitbit.New(os.Getenv("FITBIT_KEY"), os.Getenv("FITBIT_SECRET"), "http://localhost:1555/oauth/fitbit/callback"),
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:1555/oauth/google/callback"),
		gplus.New(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"), "http://localhost:1555/oauth/gplus/callback"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:1555/oauth/github/callback"),
		spotify.New(os.Getenv("SPOTIFY_KEY"), os.Getenv("SPOTIFY_SECRET"), "http://localhost:1555/oauth/spotify/callback"),
		linkedin.New(os.Getenv("LINKEDIN_KEY"), os.Getenv("LINKEDIN_SECRET"), "http://localhost:1555/oauth/linkedin/callback"),
		line.New(os.Getenv("LINE_KEY"), os.Getenv("LINE_SECRET"), "http://localhost:1555/oauth/line/callback", "profile", "openid", "email"),
		lastfm.New(os.Getenv("LASTFM_KEY"), os.Getenv("LASTFM_SECRET"), "http://localhost:1555/oauth/lastfm/callback"),
		twitch.New(os.Getenv("TWITCH_KEY"), os.Getenv("TWITCH_SECRET"), "http://localhost:1555/oauth/twitch/callback"),
		dropbox.New(os.Getenv("DROPBOX_KEY"), os.Getenv("DROPBOX_SECRET"), "http://localhost:1555/oauth/dropbox/callback"),
		digitalocean.New(os.Getenv("DIGITALOCEAN_KEY"), os.Getenv("DIGITALOCEAN_SECRET"), "http://localhost:1555/oauth/digitalocean/callback", "read"),
		bitbucket.New(os.Getenv("BITBUCKET_KEY"), os.Getenv("BITBUCKET_SECRET"), "http://localhost:1555/oauth/bitbucket/callback"),
		instagram.New(os.Getenv("INSTAGRAM_KEY"), os.Getenv("INSTAGRAM_SECRET"), "http://localhost:1555/oauth/instagram/callback"),
		intercom.New(os.Getenv("INTERCOM_KEY"), os.Getenv("INTERCOM_SECRET"), "http://localhost:1555/oauth/intercom/callback"),
		box.New(os.Getenv("BOX_KEY"), os.Getenv("BOX_SECRET"), "http://localhost:1555/oauth/box/callback"),
		salesforce.New(os.Getenv("SALESFORCE_KEY"), os.Getenv("SALESFORCE_SECRET"), "http://localhost:1555/oauth/salesforce/callback"),
		seatalk.New(os.Getenv("SEATALK_KEY"), os.Getenv("SEATALK_SECRET"), "http://localhost:1555/oauth/seatalk/callback"),
		amazon.New(os.Getenv("AMAZON_KEY"), os.Getenv("AMAZON_SECRET"), "http://localhost:1555/oauth/amazon/callback"),
		yammer.New(os.Getenv("YAMMER_KEY"), os.Getenv("YAMMER_SECRET"), "http://localhost:1555/oauth/yammer/callback"),
		onedrive.New(os.Getenv("ONEDRIVE_KEY"), os.Getenv("ONEDRIVE_SECRET"), "http://localhost:1555/oauth/onedrive/callback"),
		azuread.New(os.Getenv("AZUREAD_KEY"), os.Getenv("AZUREAD_SECRET"), "http://localhost:1555/oauth/azuread/callback", nil),
		microsoftonline.New(os.Getenv("MICROSOFTONLINE_KEY"), os.Getenv("MICROSOFTONLINE_SECRET"), "http://localhost:1555/oauth/microsoftonline/callback"),
		battlenet.New(os.Getenv("BATTLENET_KEY"), os.Getenv("BATTLENET_SECRET"), "http://localhost:1555/oauth/battlenet/callback"),
		eveonline.New(os.Getenv("EVEONLINE_KEY"), os.Getenv("EVEONLINE_SECRET"), "http://localhost:1555/oauth/eveonline/callback"),
		kakao.New(os.Getenv("KAKAO_KEY"), os.Getenv("KAKAO_SECRET"), "http://localhost:1555/oauth/kakao/callback"),

		// Pointed https://localhost.com to http://localhost:1555/oauth/yahoo/callback
		// Yahoo only accepts urls that starts with https
		yahoo.New(os.Getenv("YAHOO_KEY"), os.Getenv("YAHOO_SECRET"), "https://localhost.com"),
		typetalk.New(os.Getenv("TYPETALK_KEY"), os.Getenv("TYPETALK_SECRET"), "http://localhost:1555/oauth/typetalk/callback", "my"),
		slack.New(os.Getenv("SLACK_KEY"), os.Getenv("SLACK_SECRET"), "http://localhost:1555/oauth/slack/callback"),
		stripe.New(os.Getenv("STRIPE_KEY"), os.Getenv("STRIPE_SECRET"), "http://localhost:1555/oauth/stripe/callback"),
		wepay.New(os.Getenv("WEPAY_KEY"), os.Getenv("WEPAY_SECRET"), "http://localhost:1555/oauth/wepay/callback", "view_user"),
		// By default paypal production auth urls will be used, please set PAYPAL_ENV=sandbox as environment variable for testing
		// in sandbox environment
		paypal.New(os.Getenv("PAYPAL_KEY"), os.Getenv("PAYPAL_SECRET"), "http://localhost:1555/oauth/paypal/callback"),
		steam.New(os.Getenv("STEAM_KEY"), "http://localhost:1555/oauth/steam/callback"),
		heroku.New(os.Getenv("HEROKU_KEY"), os.Getenv("HEROKU_SECRET"), "http://localhost:1555/oauth/heroku/callback"),
		uber.New(os.Getenv("UBER_KEY"), os.Getenv("UBER_SECRET"), "http://localhost:1555/oauth/uber/callback"),
		soundcloud.New(os.Getenv("SOUNDCLOUD_KEY"), os.Getenv("SOUNDCLOUD_SECRET"), "http://localhost:1555/oauth/soundcloud/callback"),
		gitlab.New(os.Getenv("GITLAB_KEY"), os.Getenv("GITLAB_SECRET"), "http://localhost:1555/oauth/gitlab/callback"),
		dailymotion.New(os.Getenv("DAILYMOTION_KEY"), os.Getenv("DAILYMOTION_SECRET"), "http://localhost:1555/oauth/dailymotion/callback", "email"),
		deezer.New(os.Getenv("DEEZER_KEY"), os.Getenv("DEEZER_SECRET"), "http://localhost:1555/oauth/deezer/callback", "email"),
		discord.New(os.Getenv("DISCORD_KEY"), os.Getenv("DISCORD_SECRET"), "http://localhost:1555/oauth/discord/callback", discord.ScopeIdentify, discord.ScopeEmail),
		meetup.New(os.Getenv("MEETUP_KEY"), os.Getenv("MEETUP_SECRET"), "http://localhost:1555/oauth/meetup/callback"),

		// Auth0 allocates domain per customer, a domain must be provided for auth0 to work
		auth0.New(os.Getenv("AUTH0_KEY"), os.Getenv("AUTH0_SECRET"), "http://localhost:1555/oauth/auth0/callback", os.Getenv("AUTH0_DOMAIN")),
		xero.New(os.Getenv("XERO_KEY"), os.Getenv("XERO_SECRET"), "http://localhost:1555/oauth/xero/callback"),
		vk.New(os.Getenv("VK_KEY"), os.Getenv("VK_SECRET"), "http://localhost:1555/oauth/vk/callback"),
		naver.New(os.Getenv("NAVER_KEY"), os.Getenv("NAVER_SECRET"), "http://localhost:1555/oauth/naver/callback"),
		yandex.New(os.Getenv("YANDEX_KEY"), os.Getenv("YANDEX_SECRET"), "http://localhost:1555/oauth/yandex/callback"),
		nextcloud.NewCustomisedDNS(os.Getenv("NEXTCLOUD_KEY"), os.Getenv("NEXTCLOUD_SECRET"), "http://localhost:1555/oauth/nextcloud/callback", os.Getenv("NEXTCLOUD_URL")),
		gitea.New(os.Getenv("GITEA_KEY"), os.Getenv("GITEA_SECRET"), "http://localhost:1555/oauth/gitea/callback"),
		shopify.New(os.Getenv("SHOPIFY_KEY"), os.Getenv("SHOPIFY_SECRET"), "http://localhost:1555/oauth/shopify/callback", shopify.ScopeReadCustomers, shopify.ScopeReadOrders),
		apple.New(os.Getenv("APPLE_KEY"), os.Getenv("APPLE_SECRET"), "http://localhost:1555/oauth/apple/callback", nil, apple.ScopeName, apple.ScopeEmail),
		strava.New(os.Getenv("STRAVA_KEY"), os.Getenv("STRAVA_SECRET"), "http://localhost:1555/oauth/strava/callback"),
		okta.New(os.Getenv("OKTA_ID"), os.Getenv("OKTA_SECRET"), os.Getenv("OKTA_ORG_URL"), "http://localhost:1555/oauth/okta/callback", "openid", "profile", "email"),
		mastodon.New(os.Getenv("MASTODON_KEY"), os.Getenv("MASTODON_SECRET"), "http://localhost:1555/oauth/mastodon/callback", "read:accounts"),
		wecom.New(os.Getenv("WECOM_CORP_ID"), os.Getenv("WECOM_SECRET"), os.Getenv("WECOM_AGENT_ID"), "http://localhost:1555/oauth/wecom/callback"),
		zoom.New(os.Getenv("ZOOM_KEY"), os.Getenv("ZOOM_SECRET"), "http://localhost:1555/oauth/zoom/callback", "read:user"),
		patreon.New(os.Getenv("PATREON_KEY"), os.Getenv("PATREON_SECRET"), "http://localhost:1555/oauth/patreon/callback"),
	)
	var endpoints []Endpoint = []Endpoint{{
		Method:  "GET",
		Path:    "/oauth/:provider",
		Handler: thirdPartAuthenticator.InitAuth(db),
	},
		{
			Method:  "GET",
			Path:    "/oauth/:provider/callback",
			Handler: thirdPartAuthenticator.CallBackAuth(db),
		},
	}
	return endpoints
	/*
		m := map[string]string{
			"amazon":          "Amazon",
			"apple":           "Apple",
			"auth0":           "Auth0",
			"azuread":         "Azure AD",
			"battlenet":       "Battle.net",
			"bitbucket":       "Bitbucket",
			"box":             "Box",
			"dailymotion":     "Dailymotion",
			"deezer":          "Deezer",
			"digitalocean":    "Digital Ocean",
			"discord":         "Discord",
			"dropbox":         "Dropbox",
			"eveonline":       "Eve Online",
			"facebook":        "Facebook",
			"fitbit":          "Fitbit",
			"gitea":           "Gitea",
			"github":          "Github",
			"gitlab":          "Gitlab",
			"google":          "Google",
			"gplus":           "Google Plus",
			"heroku":          "Heroku",
			"instagram":       "Instagram",
			"intercom":        "Intercom",
			"kakao":           "Kakao",
			"lastfm":          "Last FM",
			"line":            "LINE",
			"linkedin":        "LinkedIn",
			"mastodon":        "Mastodon",
			"meetup":          "Meetup.com",
			"microsoftonline": "Microsoft Online",
			"naver":           "Naver",
			"nextcloud":       "NextCloud",
			"okta":            "Okta",
			"onedrive":        "Onedrive",
			"openid-connect":  "OpenID Connect",
			"patreon":         "Patreon",
			"paypal":          "Paypal",
			"salesforce":      "Salesforce",
			"seatalk":         "SeaTalk",
			"shopify":         "Shopify",
			"slack":           "Slack",
			"soundcloud":      "SoundCloud",
			"spotify":         "Spotify",
			"steam":           "Steam",
			"strava":          "Strava",
			"stripe":          "Stripe",
			"tiktok":          "TikTok",
			"twitch":          "Twitch",
			"twitter":         "Twitter",
			"twitterv2":       "Twitter",
			"typetalk":        "Typetalk",
			"uber":            "Uber",
			"vk":              "VK",
			"wecom":           "WeCom",
			"wepay":           "Wepay",
			"xero":            "Xero",
			"yahoo":           "Yahoo",
			"yammer":          "Yammer",
			"yandex":          "Yandex",
			"zoom":            "Zoom",
		}
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		providerIndex := &ProviderIndex{Providers: keys, ProvidersMap: m}
	*/
}

func (thirdPartAuthenticator *ThirdPartAuthenticator) InitAuth(db MongoDBconnector) gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := c.Param("provider")
		q := c.Request.URL.Query()
		q.Add("provider", provider)
		c.Request.URL.RawQuery = q.Encode()
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

func (tpa *ThirdPartAuthenticator) CallBackAuth(db MongoDBconnector) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		provider := c.Param("provider")
		q := c.Request.URL.Query()
		q.Add("provider", provider)
		c.Request.URL.RawQuery = q.Encode()

		userThirdPart, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "authentication failed"})
			return
		}

		// Validate required fields
		if userThirdPart.Email == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "email not provided by provider"})
			return
		}

		var user models.User
		err = db.GetRecord(ctx, "users", bson.M{"email": userThirdPart.Email}, &user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				user = models.User{
					Id:         primitive.NewObjectID(),
					Email:      userThirdPart.Email,
					IsVerified: true,
					OauthId:    userThirdPart.UserID,
					CreatedAt:  time.Now().Format(time.RFC3339),
				}
				if err := db.CreateRecord(ctx, "users", user); err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user creation failed"})
					return
				}
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "database error"})
				return
			}
		}

		token, err := utils.CreateToken(user.Email, user.Id.Hex(), user.IsVerified, user.IsSuperuser)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "token creation failed"})
			return
		}

		c.JSON(http.StatusOK, models.LoginResponse{
			Token:     token,
			TokenType: "Bearer",
		})
	}
}
