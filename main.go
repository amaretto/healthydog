package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	baseURL          = "https://api.fitbit.com"
	activityEndpoint = "/1/user/-/activities/date/"
	sleepEndpoint    = "/1.2/user/-/sleep/date/"
	weightEndpoint   = "/1/user/-/body/log/weight/date/"
)

type ActivitySummary struct {
	Summary struct {
		ActivityCalories int `json:"activityCalories"`
		Distances        []struct {
			Activity string  `json:"activity"`
			Distance float64 `json:"distance"`
		} `json:"distances"`
		Steps int `json:"steps"`
	} `json:"summary"`
}

type SleepSummary struct {
	Summary struct {
		Stages struct {
			Deep  int `json:"deep"`
			Light int `json:"light"`
			Rem   int `json:"rem"`
			Wake  int `json:"wake"`
		} `json:"stages"`
		TotalMinutesAsleep int `json:"totalMinutesAsleep"`
		TotalTimeInBed     int `json:"totalTimeInBed"`
	} `json:"summary"`
}

type WeightSummary struct {
	Weight []struct {
		Bmi    float64 `json:"bmi"`
		Weight int     `json:"weight"`
	} `json:"weight"`
}

func main() {
	date := "2021-11-21.json"
	token := os.Getenv("FITBIT_TOKEN")
	//refreshToken := os.Getenv("FITBIT_REFRESH_TOKEN")

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
