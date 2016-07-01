package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"net/url"
	//"reflect"
	"time"
)

type ApiConf struct {
	YourAccountName   string `json:"user_name"`
	ConsumerKey       string `json:"consumer_key"`
	ConsumerSecret    string `json:"consumer_secret"`
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("----------------------------------")
	fmt.Println("|       Tweet Sweeper v1.0       |")
	fmt.Println("----------------------------------")

	var apiConf ApiConf
	{
		apiConfPath := flag.String("conf", "config.json", "API Config File")
		flag.Parse()
		data, err_file := ioutil.ReadFile(*apiConfPath)
		check(err_file)
		err_json := json.Unmarshal(data, &apiConf)
		check(err_json)
	}

	anaconda.SetConsumerKey(apiConf.ConsumerKey)
	anaconda.SetConsumerSecret(apiConf.ConsumerSecret)
	client := anaconda.NewTwitterApi(apiConf.AccessToken, apiConf.AccessTokenSecret)

	// Setting parameter using url.Values
	v := url.Values{}
	s := client.UserStream(v)

	for t := range s.C {
		switch status := t.(type) {
		case anaconda.Tweet:
			if status.User.ScreenName == apiConf.YourAccountName {
				go func() {
					fmt.Printf("「%s」を5分後に消します。\n", status.Text)
					// 5分後に削除
					time.Sleep(5 * time.Minute)
					tw, _ := client.DeleteTweet(status.Id, true)
					fmt.Printf("「%s」を消しました。\n", tw.Text)
				}()
			}
		default:
		}
	}
}
