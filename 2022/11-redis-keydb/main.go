package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var ctx = context.Background()
var results [60]float64

var (
	avgReqTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redis_request_time_avg",
			Help: "Gauge redis average request time for last 60 sec",
		},
		[]string{"redis"},
	)
	maxReqTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "redis_request_time_max",
		Help: "Gauge redis max request time for last 60 sec",
	},
		[]string{"redis"},
	)
	minReqTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "redis_request_time_min",
		Help: "Gauge redis min request time for last 60 sec",
	},
		[]string{"redis"},
	)
	failedReq = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "redis_request_fail",
		Help: "Counter redis key set fails",
	},
		[]string{"redis"},
	)
	RedisAddress = flag.String("redis", "127.0.0.1:6379", "redis cluster addresses")
	RedisKey     = flag.String("key", "test", "test-key name prefix")
)

func findValues(a [60]float64) (min float64, max float64, avg float64) {
	min = 60
	max = 0
	var sum float64 = 0
	var count int = 0

	for _, value := range a {
		if value >= 0 {
			count++
			sum += value
		}
		if value < min && value > 0 {
			min = value
		}
		if value > max {
			max = value
		}
	}
	if count > 0 {
		return min, max, sum / float64(count)
	} else {
		return 60, 60, 60
	}

}

func collect(redis_address string, key string) {
	var i = 0
	for i := 0; i < 60; i++ {
		results[i] = -1
	}
	for {
		i = i % 60
		key := fmt.Sprintf("%s_%d", key, rand.Int())
		start := time.Now()
		rdb := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        []string{redis_address},
			DialTimeout:  1000 * time.Millisecond,
			ReadTimeout:  1000 * time.Millisecond,
			WriteTimeout: 1000 * time.Millisecond,
		})
		err := rdb.SetNX(ctx, key, "ping", 2*time.Second).Err()
		duration := time.Since(start)

		if err != nil {
			results[i] = -1
			failedReq.With(prometheus.Labels{"redis": string(redis_address)}).Inc()
		} else {
			results[i] = duration.Seconds()
		}
		rdb.Close()
		if duration.Milliseconds() < 1000 {
			time.Sleep(time.Duration(1000-duration.Milliseconds()) * time.Millisecond)
		}
		min, max, avg := findValues(results)
		avgReqTime.With(prometheus.Labels{"redis": string(redis_address)}).Set(avg)
		maxReqTime.With(prometheus.Labels{"redis": string(redis_address)}).Set(max)
		minReqTime.With(prometheus.Labels{"redis": string(redis_address)}).Set(min)
		fmt.Printf("Max: %0.4f Min: %0.4f Avg: %0.4f Lst: %0.4f\n", max, min, avg, results[i])
		i++
	}
}

func main() {
	flag.Parse()
	prometheus.MustRegister(avgReqTime)
	prometheus.MustRegister(minReqTime)
	prometheus.MustRegister(maxReqTime)
	prometheus.MustRegister(failedReq)
	failedReq.With(prometheus.Labels{"redis": string(*RedisAddress)}).Add(0)
	go collect(*RedisAddress, *RedisKey)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9379", nil)
}
