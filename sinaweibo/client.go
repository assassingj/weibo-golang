package sinaweibo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	API_URL = "https://api.weibo.com/"
)

type WeiboAuth struct {
	ClientId     string
	ClientSecret string
	RedirectUri  string
}

type WeiboClient struct {
	Uid               string
	accessTokenResult AccessTokenResult
}

func NewWeiboClient(accessTokenResult AccessTokenResult) *WeiboClient {
	return &WeiboClient{accessTokenResult.Uid, accessTokenResult}
}

type User struct {
	Id         int64
	IdStr      string
	ScreenName string
	Name       string
	Gender     string
	Location   string
}

type Friendship struct {
	Users []User
}

//could be used as map key
// func (u User) Equal(other User) bool {
// 	return u.Id == other.Id
// }

type AccessTokenResult struct {
	//参考http://open.weibo.com/wiki/%E6%8E%88%E6%9D%83%E6%9C%BA%E5%88%B6%E8%AF%B4%E6%98%8E#.E4.BD.BF.E7.94.A8OAuth2.0.E8.B0.83.E7.94.A8API
	AccessToken string `json:"access_token"`
	RemindIn    string `json:"remind_in"`
	ExpiresIn   int    `json:"expires_in"`
	Uid         string `json:"uid"`
}

func (auth *WeiboAuth) GetAuthorizeUrl() string {
	params := url.Values{}
	params.Add("client_id", auth.ClientId)
	params.Add("response_type", "code")
	params.Add("redirect_uri", auth.RedirectUri)
	return fmt.Sprintf("%soauth2/authorize?%s", API_URL, params.Encode())
}

func (auth *WeiboAuth) GetAccessToken(code string) (AccessTokenResult, error) {
	//https://api.weibo.com/oauth2/access_token?client_id=YOUR_CLIENT_ID&client_secret=YOUR_CLIENT_SECRET&grant_type=authorization_code&redirect_uri=YOUR_REGISTERED_REDIRECT_URI&code=CODE
	params := url.Values{}
	params.Add("client_id", auth.ClientId)
	params.Add("client_secret", auth.ClientSecret)
	params.Add("grant_type", "authorization_code")
	params.Add("redirect_uri", auth.RedirectUri)
	params.Add("code", code)
	accessTokenUrl := fmt.Sprintf("%soauth2/access_token?", API_URL)
	response, err := http.PostForm(accessTokenUrl, params)
	if err != nil {
		log.Println("error while get authorize code")
		panic(err)
	}
	defer response.Body.Close()
	accessTokenResult := AccessTokenResult{}
	str, _ := ioutil.ReadAll(response.Body)
	log.Printf("body:%s", str)
	err = json.Unmarshal(str, &accessTokenResult)
	if err != nil {
		log.Println("error while parsing token json")
		panic(err)
	}
	return accessTokenResult, nil
}

func (weiboClient *WeiboClient) httpGet(urlStr string) []byte {
	log.Printf("get url:%s\n", urlStr)
	// client := &http.Client{}
	// request, _ := http.NewRequest("Get", urlStr, nil)
	// request.Header.Set("Authorization", "OAuth2 "+weiboClient.accessTokenResult.AccessToken)
	// response, err := client.Do(request)
	urlStr = fmt.Sprintf("%s&access_token=%s", urlStr, weiboClient.accessTokenResult.AccessToken)
	response, err := http.Get(urlStr)
	if err != nil {
		log.Printf("request error while getting url:%s", urlStr)
		return nil
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		log.Printf("status code error[%d] while getting url:%s\n return body:%s\n",
			response.StatusCode, urlStr, body)
	}
	return body
}

func parseFriend(content []byte) Friendship {
	var f Friendship
	json.Unmarshal(content, &f)
	return f
}

func (weiboClient *WeiboClient) GetFriendshipsFollowers(uid string) Friendship {
	urlStr := fmt.Sprintf("%s2/friendships/followers.json?uid=%s&count=200", API_URL, uid)
	return parseFriend(weiboClient.httpGet(urlStr))
}

func (weiboClient *WeiboClient) GetFriendshipsBilateral(uid string) Friendship {
	urlStr := fmt.Sprintf("%s2/friendships/friends/bilateral.json?uid=%s&count=200", API_URL, uid)
	return parseFriend(weiboClient.httpGet(urlStr))
}

func (weiboClient *WeiboClient) GetFriendships(uid string) Friendship {
	urlStr := fmt.Sprintf("%s2/friendships/friends.json?uid=%s&count=200", API_URL, uid)
	return parseFriend(weiboClient.httpGet(urlStr))
}
