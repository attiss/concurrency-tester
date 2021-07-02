package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	config := edgegrid.Config{
		Host:         os.Getenv("AKAMAI_HOST"),
		ClientToken:  os.Getenv("AKAMAI_CLIENT_TOKEN"),
		ClientSecret: os.Getenv("AKAMAI_CLIENT_SECRET"),
		AccessToken:  os.Getenv("AKAMAI_ACCESS_TOKEN"),
		MaxBody:      4096,
	}

	dnsv2.Init(config)

	zone := os.Getenv("ZONE")
	if zone == "" {
		panic("missing ZONE")
	}

	amountPerGoRoutine := 15
	if os.Getenv("CREATES_PER_GOROUTINE") != "" {
		var convErr error
		amountPerGoRoutine, convErr = strconv.Atoi(os.Getenv("CREATES_PER_GOROUTINE"))
		if convErr != nil {
			panic(fmt.Errorf("failed to convert CREATES_PER_GOROUTINE: %+v", convErr))
		}
	}

	goRoutines := 5
	if os.Getenv("GOROUTINES") != "" {
		var convErr error
		goRoutines, convErr = strconv.Atoi(os.Getenv("GOROUTINES"))
		if convErr != nil {
			panic(fmt.Errorf("failed to convert GOROUTINES: %+v", convErr))
		}
	}

	logger.Info("starting record creation", zap.Int("amountPerGoRoutine", amountPerGoRoutine), zap.Int("goRoutines", goRoutines))

	var wg sync.WaitGroup
	for i := 0; i < goRoutines; i++ {
		wg.Add(1)
		go createRecords(zone, amountPerGoRoutine, logger.With(zap.Int("goRoutineID", i)), &wg)
	}
	wg.Wait()

	logger.Info("finished record creation")
}

func createRecords(zone string, amount int, logger *zap.Logger, wg *sync.WaitGroup) {
	rand.Seed(time.Now().UnixNano())

	numberOfRequests := 0

	for i := 0; i < amount; i++ {
		record := fmt.Sprintf("test-%d.%s", rand.Intn(1000000000), zone)

		recordSet := dnsv2.RecordBody{
			Name:       record,
			RecordType: "CNAME",
			Target:     []string{"this.is.a.test."},
			TTL:        60,
		}

		for {
			numberOfRequests++
			err := recordSet.Save(zone)
			if err == nil {
				logger.Info("record successfully created", zap.String("record", record))
				break
			}
			logger.Error("failed to create record", zap.String("record", record), zap.Error(err))
			<-time.After(1 * time.Second)
		}
	}

	logger.Info("successfully created all records", zap.Int("amountOfRecords", amount), zap.Int("numberOfRequests", numberOfRequests))

	wg.Done()
}
