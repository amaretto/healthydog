package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	token, refreshToken := getToken()
	date := time.Now().Format("2006-01-02.json")

	if err := checkToken(token, refreshToken); err != nil {
		log.Fatal(err)
	}

	// activity
	apiEndPoint := baseURL + activityEndpoint + date
	req, _ := http.NewRequest("GET", apiEndPoint, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	var activitySummary ActivitySummary
	if err := json.Unmarshal(byteArray, &activitySummary); err != nil {
		os.Exit(1)
	}
	fmt.Println("ActivityCalories ", activitySummary.Summary.ActivityCalories)
	fmt.Println("Distance ", activitySummary.Summary.Distances[0].Distance)
	fmt.Println("Steps ", activitySummary.Summary.Steps)

	// sleep
	apiEndPoint = baseURL + sleepEndpoint + date
	req, _ = http.NewRequest("GET", apiEndPoint, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()

	byteArray, _ = ioutil.ReadAll(resp.Body)

	var sleepSummary SleepSummary
	if err := json.Unmarshal(byteArray, &sleepSummary); err != nil {
		os.Exit(1)
	}
	fmt.Println("Summary ", sleepSummary.Summary)
	fmt.Println("TotalMinutesAsleep ", sleepSummary.Summary.TotalMinutesAsleep)
	fmt.Println("TotalTimeInBed ", sleepSummary.Summary.TotalTimeInBed)

	// weight
	apiEndPoint = baseURL + weightEndpoint + date
	req, _ = http.NewRequest("GET", apiEndPoint, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()

	byteArray, _ = ioutil.ReadAll(resp.Body)
	var weightSummary WeightSummary
	if err := json.Unmarshal(byteArray, &weightSummary); err != nil {
		os.Exit(1)
	}
	fmt.Println("Weight ", weightSummary.Weight[0].Weight)
	fmt.Println("BMI ", weightSummary.Weight[0].Bmi)

}
