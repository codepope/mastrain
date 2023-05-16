package main

import (
	"context"
	"fmt"
	"log"

	"mastrain/raindrop"

	"github.com/mattn/go-mastodon"
)

var config map[string]string

func main() {

	envs, error := GetConfig()

	if error != nil {
		log.Fatalf("Error loading .env or ENV: %v", error)
	}

	//fmt.Printf("%v", envs)

	rainclient := raindrop.NewAPI(envs["RAINDROP_SERVER"], envs["RAINDROP_APP_TOKEN"])
	mastodonClient := mastodon.NewClient(&mastodon.Config{envs["MASTODON_SERVER"], envs["MASTODON_CLIENT_ID"], envs["MASTODON_CLIENT_SECRET"], envs["MASTODON_APP_TOKEN"]})
	bookmarks, err := mastodonClient.GetBookmarks(context.Background(), nil)

	if err != nil {
		log.Fatalf("Error getting bookmarks: %v", error)
	}

	var bkurls []string

	for _, bk := range bookmarks {
		bkurls = append(bkurls, bk.URL)
	}

	duplicates, err := rainclient.GetDuplicates(context.Background(), bkurls)

	if err != nil {
		log.Fatalf("Error getting duplicates: %v", err)
	}

	fmt.Printf("Duplicates: %v\n", duplicates)

	// Now we can remove the duplicates from the bookmarks slice
	if duplicates != nil {
		for _, dup := range *duplicates {
			for i, bk := range bookmarks {
				if bk.URL == dup.Link {
					bookmarks = append(bookmarks[:i], bookmarks[i+1:]...)
				}
			}
		}
	}

	fmt.Printf("Bookmarks: %v\n", bookmarks)

	for _, bk := range bookmarks {

		purl, err := rainclient.ParseURL(context.Background(), bk.URL)
		if err != nil {
			fmt.Printf("Error parsing URL: %v\n", err)
		} else {
			fmt.Printf("Tags: %v\n", purl.Item.Meta.Tags)
			purl.Item.Meta.Tags = append(purl.Item.Meta.Tags, "mastodon")
			newdrop := raindrop.Raindrop{Title: purl.Item.Title, Link: bk.URL, Excerpt: purl.Item.Excerpt, CollectionID: -1, Tags: purl.Item.Meta.Tags}
			err := rainclient.SaveRaindrop(context.Background(), &newdrop)
			if err != nil {
				fmt.Printf("Error saving raindrop: %v\n", err)
			}
		}
	}

}
