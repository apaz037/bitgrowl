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
	api "github.com/miguelmota/go-coinmarketcap/v2"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "the start command receives a ticker, a retrieval interval, an upperbound, and a lowerbound",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ticker := args[0]
		timerString := args[1]
		upperBound := args[2]
		lowerBound := args[3]

		timer, err := time.ParseDuration(timerString)
		if err != nil {
			fmt.Println(err)
		}

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
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func start(ticker string, timer time.Duration, upperBound, lowerBound float64) error {
	// TODO: need to cast bounds to appropriate datatype before comparisons
	price, err := getPrice(ticker)
	if err != nil {
		return err
	}

	if price < lowerBound {
		fmt.Println("price has fallen beneath lowerBound")
		roar()
		time.Sleep(time.Minute * 20)
	} else if price > upperBound {
		fmt.Println("price has risen above upperbound")
		roar()
		time.Sleep(time.Minute * 20)
	}

	time.Sleep(timer)
	return nil
}

func getPrice(ticker string) (price float64, err error) {
	// TODO: refactor variadic input or make command elligible for variadic input
	priceOptions := &api.PriceOptions{
		Symbol:  ticker,
		Convert: "USD",
	}

	quote, err := api.Price(priceOptions)
	if err != nil {
		return 0, err
	}

	return quote, err
}

func roar() {

}
