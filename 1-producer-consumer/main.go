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

func producer(stream Stream, ch chan<- Tweet) (tweets []*Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			break
		}

		ch <- *tweet

	}

	close(ch)
	return tweets
}

func consumer(wg *sync.WaitGroup, ch <-chan Tweet) {
	defer wg.Done()
	for t := range ch {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	tweetChan := make(chan Tweet)

	// Producer
	go producer(stream, tweetChan)

	var wg sync.WaitGroup

	wg.Add(1)
	// Consumer
	go consumer(&wg, tweetChan)

	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
