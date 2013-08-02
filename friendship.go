package main

import (
	"./sinaweibo"
	"fmt"
	"log"
	"sort"
	"time"
)

type sortedMap struct {
	m map[int64]int
	s []int64
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func sortedKeys(m map[int64]int) []int64 {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]int64, len(m))
	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}

func buildExistenceMap(weiboClient *sinaweibo.WeiboClient, uid string) map[int64]bool {
	existenceMap := make(map[int64]bool)
	f := weiboClient.GetFriendships(uid)
	for _, friend := range f.Users {
		existenceMap[friend.Id] = true
	}
	return existenceMap

}

func main() {
	weiboAuth := sinaweibo.WeiboAuth{"your-app-key", "your-app-secret", "redirect-url"}
	fmt.Printf("please explore url %s\n", weiboAuth.GetAuthorizeUrl())
	var code string
	fmt.Println("please input your code:")
	fmt.Scanf("%s", &code)
	token, err := weiboAuth.GetAccessToken(code)
	if err != nil {
		log.Fatal(err)
	}

	weiboClient := sinaweibo.NewWeiboClient(token)
	bf := weiboClient.GetFriendshipsBilateral(token.Uid)

	existenceMap := buildExistenceMap(weiboClient, token.Uid)
	relationShip := make(map[int64]sinaweibo.User)
	counter := make(map[int64]int)

	for _, friend := range bf.Users {
		secondoryFriends := weiboClient.GetFriendshipsBilateral(friend.IdStr)
		time.Sleep(20 * time.Microsecond)
		for _, secondoryFriend := range secondoryFriends.Users {
			relationShip[secondoryFriend.Id] = secondoryFriend
			counter[secondoryFriend.Id]++
		}
	}

	for _, id := range sortedKeys(counter) {
		if _, ok := existenceMap[id]; ok || fmt.Sprint(id) == token.Uid {
			continue
		}
		u := relationShip[id]
		fmt.Printf("%d,%s,%d\n", u.Id, u.Name, counter[id])
	}

}
