package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	go func() {
		for {
			err := validateCoupon()
			filterError(err)
			if err != nil {
				log.Printf("error from worker 1: %s", err.Error())
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {
		for {
			err := validateCoupon()
			filterError(err)
			if err != nil {
				log.Printf("error from worker 2: %s", err.Error())
			}
			time.Sleep(time.Millisecond * 400)
		}
	}()

	go func() {
		for {
			err := validateCoupon()
			filterError(err)
			if err != nil {
				log.Printf("error from worker 3: %s", err.Error())
			}
			time.Sleep(time.Millisecond * 700)
		}
	}()
	go func() {
		for {
			err := validateCoupon()
			filterError(err)
			if err != nil {
				log.Printf("error from worker 3: %s", err.Error())
			}
			time.Sleep(time.Millisecond * 1000)
		}
	}()

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Println("Signal terminate detected")
	}
	log.Println()
	log.Printf("NUMBER OF REQUEST: %d", numberOfRequests)
	log.Printf("NUMBER OF SUCCESS: %d", numberOfSuccess)
	log.Printf("NUMBER OF FAILED: %d", numberOfFailed)
}

var numberOfRequests int
var numberOfFailed int
var numberOfSuccess int

func filterError(err error) {
	if err == nil {
		numberOfSuccess++
		numberOfRequests++
		return
	}
	if err.Error() == "quantity have changed before update" {
		numberOfFailed++
	}
	if err.Error() == "maximum quantity reached" {
		return
	}
	numberOfRequests++
}

func validateCoupon() error {
	resp, err := http.Get("http://localhost:9090/coupon/validate?code=SHOP-FRIDAYSALE")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(content) != "validate coupon success" {
		return errors.New(string(content))
	}
	return nil
}
