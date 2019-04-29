// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gen2brain/beeep"
	api "github.com/miguelmota/go-coinmarketcap/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	High = "HIGH"
	Low = "LOW"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "the start command receives a ticker, a retrieval interval, an upperbound, and a lowerbound",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: configurable timer via flag?
		timer := time.Minute

		ticker := args[0]
		upperBound := args[2]
		lowerBound := args[3]

		upperBoundNum, err := strconv.ParseInt(upperBound, 10, 32)
		if err != nil {
			fmt.Println(err)
		}

		lowerBoundNum, err := strconv.ParseInt(lowerBound, 10, 32)
		if err != nil {
			fmt.Println(err)
		}

		err = start(ticker, timer, float64(upperBoundNum), float64(lowerBoundNum))
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	// TODO: configurable timer via flag? Configurable ICON? Configurable Noise?
	//startCmd.Flags().Int8P("timer", "t", 1, "timer in minutes")
}

func start(ticker string, timer time.Duration, upperBound, lowerBound float64) error {
	for {
		price, err := getPrice(ticker)
		if err != nil {
			return err
		}

		if price < lowerBound {
			err := roar(ticker, "LOW", price)
			if err != nil {
				return err
			}
			time.Sleep(time.Minute * 15)
		} else if price > upperBound {
			err := roar(ticker, "HIGH", price)
			if err != nil {
				return err
			}
			time.Sleep(time.Minute * 15)
		}

		time.Sleep(timer)
		return nil
	}
}

func getPrice(ticker string) (price float64, err error) {
	if len(ticker) > 5  || len(ticker) < 3{
		err := errors.New("Invalid Ticker")
		return 0, err
	}

	priceOptions := &api.PriceOptions{
		Symbol:  ticker,
		Convert: "USD",
	}

	quote, err := api.Price(priceOptions)
	if err != nil {
		return 0, errors.Wrap(err, "Could not retrieve Price: ")
	}

	return quote, err
}

func roar(ticker, status string, price float64) error {
	var msg string
	var title string

	if status == Low {
		title = fmt.Sprintf("%s: Above Threshold", ticker)
		msg = fmt.Sprintf("the price has risen above the set threshold: %v", price)
	} else if status == High {
		title = fmt.Sprintf("%s: Below Threshold", ticker)
		msg = fmt.Sprintf("the price has fallen below the set threshold: %v", price)
	}

	err := beeep.Alert(title, msg, "../assets/baseline_warning_black_18dp.png")
	if err != nil {
		panic(err)
	}

	filename := "YOU_SUFFER.mp3"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	speaker.Play(streamer)

	return nil
}

