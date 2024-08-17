//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

// func producer(stream Stream) (tweets []*Tweet) {
// 	for {
// 		tweet, err := stream.Next()
// 		if err == ErrEOF {
// 			return tweets
// 		}

// 		tweets = append(tweets, tweet)
// 	}
// }

// func consumer(tweets []*Tweet) {
// 	for _, t := range tweets {
// 		if t.IsTalkingAboutGo() {
// 			fmt.Println(t.Username, "\ttweets about golang")
// 		} else {
// 			fmt.Println(t.Username, "\tdoes not tweet about golang")
// 		}
// 	}
// }

//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

// func main() {
// 	start := time.Now()
// 	stream := GetMockStream()

// 	// Producer
// 	tweets := producer(stream)

// 	// Consumer
// 	consumer(tweets)

// 	fmt.Printf("Process took %s\n", time.Since(start))
// }

// MY SOLUTION:

func producer(stream Stream, tweetChan chan *Tweet, done chan bool) {
	defer close(tweetChan)
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			done <- true
		}

		tweetChan <- tweet
	}
}

func consumer(t *Tweet) {
	if t.IsTalkingAboutGo() {
		fmt.Println(t.Username, "\ttweets about golang")
	} else {
		fmt.Println(t.Username, "\tdoes not tweet about golang")
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	var wg sync.WaitGroup

	done := make(chan bool)
	tweetChan := make(chan *Tweet)

	// Producer
	go func() {
		producer(stream, tweetChan, done)
	}()

	// Consumer
	go func() {
		for t := range tweetChan {
			wg.Add(1)
			go func(tweet *Tweet) {
				defer wg.Done()
				consumer(tweet)
			}(t)
		}
	}()

	wg.Wait()
	<-done

	fmt.Printf("Process took %s\n", time.Since(start))
}

// Before:
// ❯ make run
// davecheney      tweets about golang
// beertocode      does not tweet about golang
// ironzeb         tweets about golang
// beertocode      tweets about golang
// vampirewalk666  tweets about golang
// Process took 3.578687394s

// After
// ❯  make run
// davecheney      tweets about golang
// beertocode      does not tweet about golang
// ironzeb         tweets about golang
// beertocode      tweets about golang
// TWEETS []
// Process took 1.923645175s
