package sinaweibo

import (
	"fmt"
	"testing"
)

// var weiboAuth = WeiboAuth{"1373386909", "7f7c91620ef8d737de41e0538bd82d71", "www.lovelin.info"}

// func TestGetAuthUrl(t *testing.T) {
// 	url := weiboAuth.GetAuthorizeUrl()
// 	t.Log(url)
// }

// func TestGetFriendShip(t *testing.T) {
// 	code := "0a4142482337cb6efbfac60a53d19fae"
// 	token, err := weiboAuth.GetAccessToken(code)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	weiboClient := NewWeiboClient(token)
// 	fmt.Println(weiboClient.GetFriendshipsFollowers(token.Uid))
// }

func TestUser(t *testing.T) {
	u1 := User{Id: 111, Name: "1xxx"}
	u2 := User{Id: 112}
	u3 := User{Id: 111, Name: "fdsf"}

	m := make(map[User]bool)
	m[u1] = true
	m[u2] = true
	m[u3] = true
	fmt.Println(m)
}
