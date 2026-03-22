package cmd

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	targetURL   string
	requests    int
	concurrency int
	method string
	bodypayload string
)

var attackCmd = &cobra.Command{
	Use:   "attack",
	Short: "Build an Attack ",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runAttack(method, targetURL, bodypayload, requests, concurrency)
	},
}

func init() {
	rootCmd.AddCommand(attackCmd)

	attackCmd.Flags().StringVarP(&targetURL, "url", "u", "", "Target URL to attack")
	attackCmd.Flags().IntVarP(&requests, "requests", "r", 100, "Total number of requests")
	attackCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 10, "Number of concurrent workers")

	attackCmd.Flags().StringVarP(&method, "method", "m", "GET", "Type Of Request (GET, POST, PUT, DELETE)")
	attackCmd.Flags().StringVarP(&bodypayload, "body", "b", "", "Body For POST/PUT Requests")
	
	attackCmd.MarkFlagRequired("url")
}

type Result struct {
	Duration   time.Duration
	StatusCode int
	Error      error
}

func worker(method string, url string, body string, jobs <-chan int, result chan <- Result, wg *sync.WaitGroup){
	defer wg.Done()

	for range jobs{
		var reqBody io.Reader

		if body != ""{
			reqBody = strings.NewReader(body)
		}

		req, err := http.NewRequest(method, url, reqBody)
		if err == nil && body != "" {
			req.Header.Set("Content-Type", "application/json")
		}

		start := time.Now();

		resp, err := http.DefaultClient.Do(req)

		duration := time.Since(start)

		if err != nil{
			result <- Result{Duration: duration, StatusCode: 0, Error: err}
			continue;
		}

		result <- Result{Duration: duration, StatusCode: resp.StatusCode, Error: nil}

		resp.Body.Close()
	}
}

func runAttack(method, targetUrl, body string, totalRequests int, concurrency int){
	fmt.Printf("Attacking %s %s with %d requests (%d concurrent)...\n", method, targetUrl, totalRequests, concurrency)

	attackStart := time.Now()

	jobs := make(chan int, totalRequests)
	results := make(chan Result, totalRequests)

	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++{
		wg.Add(1)
		go worker(method, targetUrl,body, jobs, results, &wg)
	}

	for i := 0; i < totalRequests; i++ {
		jobs <- i
	}
	close(jobs)

	go func ()  {
		wg.Wait()
		close(results)
	}()
	totalDuration := time.Since(attackStart)
	processResults(results, totalRequests, totalDuration)
}

func processResults(results <- chan Result, total int, totalDuration time.Duration){
	var durations []time.Duration
	successCount := 0
	errorCount := 0

	for res := range results {
		if res.Error != nil || res.StatusCode >= 500 {
			errorCount++
		} else{
			successCount++
		}
		durations = append(durations, res.Duration)
	}

	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	p50Index := total * 50 / 100
	p99Index := total * 99 / 100

	if p99Index >= len(durations){
		p99Index = len(durations)-1
	}

	rps := float64(total) / totalDuration.Seconds()

	fmt.Println("\n--- Attack Complete ---")
	fmt.Printf("Total Time: \t\t%v\n", totalDuration)
	fmt.Printf("Requests/sec (RPS): \t%.2f\n", rps)
	fmt.Println("-----------------------")
	fmt.Printf("Total Requests: \t%d\n", total)
	fmt.Printf("Successful (2xx-4xx): \t%d\n", successCount)
	fmt.Printf("Failed (5xx/Err): \t%d\n", errorCount)
	fmt.Println("-----------------------")
	if len(durations) > 0 {
		fmt.Printf("Fastest Request: \t%v\n", durations[0])
		fmt.Printf("p50 Latency (Median): \t%v\n", durations[p50Index])
		fmt.Printf("p99 Latency: \t\t%v\n", durations[p99Index])
		fmt.Printf("Slowest Request: \t%v\n", durations[len(durations)-1])
	}
}