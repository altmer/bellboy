package tumblr

import (
	"net/url"
	"strconv"

	"github.com/kurrik/oauth1a"
)

const (
	apiBlogUrl      = "http://api.tumblr.com/v2/blog/"            // api url to get blog data or write to a blog
	apiUserUrl      = "http://api.tumblr.com/v2/user/"            // api url to get user data or perform actions
	apiTaggedUrl    = "http://api.tumblr.com/v2/tagged?"          // api url to get tagged posts
	requestTokenUrl = "http://www.tumblr.com/oauth/request_token" // oauth request-token URL
	authorizeUrl    = "https://www.tumblr.com/oauth/authorize"    // oauth authorize URL
	accessTokenUrl  = "http://www.tumblr.com/oauth/access_token"  // oauth access-token URL
)

// cient represents api client structure
type client struct {
	oauthService oauth1a.Service    // oauth service used to sign HTTP requests
	config       oauth1a.UserConfig // used within the oauth HTTP signing
	apiKey       string             // consumer key used for certain API requests
}

// API represents all the available methods for interacting with tumblr
type API interface {
	BlogPosts(string, map[string]string) BlogPosts
	PostDelete(string, int) Meta

	UserLikes(map[string]string) Likes
	UserUnlike(int, string) Meta

	UserFollowing(map[string]string) UserFollowing
	UserFollow(string) Meta
}

// New is the initialization method.
// An easy way to get the credentials is to access the interactive console:
// https://api.tumblr.com/console
func New(params map[string]string) API {
	service := &oauth1a.Service{
		RequestURL:   requestTokenUrl,
		AuthorizeURL: authorizeUrl,
		AccessURL:    accessTokenUrl,
		ClientConfig: &oauth1a.ClientConfig{
			ConsumerKey:    params["consumer_key"],
			ConsumerSecret: params["consumer_secret"],
		},
		Signer: new(oauth1a.HmacSha1Signer),
	}
	config := oauth1a.NewAuthorizedConfig(params["oauth_key"], params["oauth_secret"])
	return &client{
		oauthService: *service,
		config:       *config,
		apiKey:       params["consumer_key"],
	}
}

// BlogPosts method retrieves a list of a blog's published posts
// blogHostname - The standard or custom blog hostname (e.g., example.tumblr.com, example.com)
// params - A map of the params that are included in this request. Possible parameters:
//          * type - The type of post to return. text, quote, link, answer, video, audio, photo, chat
//          * id - A specific post ID
//          * tag - Limits the response to posts with the specified tag
//          * limit - The number of posts to return: 1–20, inclusive
//          * offset - Post number to start at
//          * reblog_info - Indicates whether to return reblog information (specify true or false)
//          * notes_info - Indicates whether to return notes information (specify true or false).
//          * filter - Specifies the post format to return, other than HTML (text or raw)
func (api client) BlogPosts(blogHostname string, params map[string]string) BlogPosts {
	var blogPosts BlogPosts
	requestURL := apiBlogUrl + blogHostname + "/posts?"
	urlParams := url.Values{}
	urlParams.Set("api_key", api.apiKey)
	for key, value := range params {
		urlParams.Set(key, value)
	}
	requestURL = requestURL + urlParams.Encode()
	api.info(requestURL, &blogPosts)
	return blogPosts
}

// PostDelete method is used to delete a blog post from a blog
// blogHostname - The standard or custom blog hostname (e.g., example.tumblr.com, example.com)
// id - The ID of the post to delete
func (api client) PostDelete(blogHostname string, id int) Meta {
	requestURL := apiBlogUrl + blogHostname + "/post/delete"
	urlParams := url.Values{}
	urlParams.Set("id", strconv.Itoa(id))
	response := api.post(requestURL, urlParams.Encode())
	return response.Meta
}

// UserInfo method is used to retrieve the user's account information that matches
// the OAuth credentials submitted with the request.
func (api client) UserInfo() UserInfo {
	var userInfo UserInfo
	requestURL := apiUserUrl + "info"
	api.info(requestURL, &userInfo)
	return userInfo
}

// UserLikes method can be used to retrieve the publicly exposed likes from a blog.
// params - A map of the params that are included in this request. Possible parameters:
//          * limit - The number of results to return.  Default: 20 (1–20, inclusive)
//          * offset - Liked post number to start at.  Default: 0 (First post)
//          * before - Retrieve posts liked before the specified timestamp. Default: None
//          * after - Retrieve posts liked after the specified timestamp. Default: None
func (api client) UserLikes(params map[string]string) Likes {
	var userLikes Likes
	requestURL := apiUserUrl + "likes?"
	urlParams := url.Values{}
	for key, value := range params {
		urlParams.Set(key, value)
	}
	requestURL = requestURL + urlParams.Encode()
	api.info(requestURL, &userLikes)
	return userLikes
}

// UserFollowing method is used to retrieve the blogs followed by the user whose OAuth credentials
// are submitted with the request.
// params - A map of the params that are included in this request. Possible parameters:
//          * limit - The number of results to return.  Default: 20 (1–20, inclusive)
//          * offset - Liked post number to start at.  Default: 0 (First post)
func (api client) UserFollowing(params map[string]string) UserFollowing {
	var userFollowing UserFollowing
	requestURL := apiUserUrl + "following?"
	urlParams := url.Values{}
	for key, value := range params {
		urlParams.Set(key, value)
	}
	requestURL = requestURL + urlParams.Encode()
	api.info(requestURL, &userFollowing)
	return userFollowing
}

// UserFollow method is used to follow a specific URL
// followURL - The url to follow, formatted (blogname.tumblr.com, blogname.com)
func (api client) UserFollow(followURL string) Meta {
	requestURL := apiUserUrl + "follow"
	urlParams := url.Values{}
	urlParams.Set("url", followURL)
	response := api.post(requestURL, urlParams.Encode())
	return response.Meta
}

// UserUnfollow method is used to unfollow a specific URL
// unfollowURL - The url to unfollow, formatted (blogname.tumblr.com, blogname.com)
func (api client) UserUnfollow(unfollowURL string) Meta {
	requestURL := apiUserUrl + "unfollow"
	urlParams := url.Values{}
	urlParams.Set("url", unfollowURL)
	response := api.post(requestURL, urlParams.Encode())
	return response.Meta
}

// UserUnlike method is used to unlike a specific blog post
// id - The ID of the blog post to be unliked
// reblogKey - The reblog key string
func (api client) UserUnlike(id int, reblogKey string) Meta {
	requestURL := apiUserUrl + "unlike"
	urlParams := url.Values{}
	urlParams.Set("id", strconv.Itoa(id))
	urlParams.Set("reblog_key", reblogKey)
	response := api.post(requestURL, urlParams.Encode())
	return response.Meta
}
