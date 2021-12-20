package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	ExecutionTime string
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getFitbitData(dataEndpoint, token string) ([]byte, error) {
	date := time.Now().AddDate(0, 0, -1).Format("2006-01-02.json")
	apiEndPoint := baseURL + dataEndpoint + date
	req, _ := http.NewRequest("GET", apiEndPoint, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	return byteArray, nil
}

func HandleLambdaEvent(ctx context.Context, event Event) error {
	targetName := os.Getenv("TARGET_NAME")

	token, refreshToken := getToken()
	if err := checkToken(token, refreshToken); err != nil {
		return err
	}

	// activity
	byteArray, err := getFitbitData(activityEndpoint, token)
	if err != nil {
		return err
	}
	var activitySummary ActivitySummary
	if err := json.Unmarshal(byteArray, &activitySummary); err != nil {
		os.Exit(1)
	}
	nameTag := fmt.Sprintf("name:%s", targetName)
	tags := []string{nameTag}
	check(sendFloatDDMetrics("activityCalories", float64(activitySummary.Summary.ActivityCalories), tags))
	check(sendFloatDDMetrics("distance", activitySummary.Summary.Distances[0].Distance, tags))
	check(sendFloatDDMetrics("steps", float64(activitySummary.Summary.Steps), tags))

	fmt.Println("ActivityCalories ", activitySummary.Summary.ActivityCalories)
	fmt.Println("Distance ", activitySummary.Summary.Distances[0].Distance)
	fmt.Println("Steps ", activitySummary.Summary.Steps)

	// sleep
	byteArray, err = getFitbitData(sleepEndpoint, token)
	if err != nil {
		return err
	}
	var sleepSummary SleepSummary
	if err := json.Unmarshal(byteArray, &sleepSummary); err != nil {
		return err
	}

	tags = []string{"stage:deep", nameTag}
	check(sendFloatDDMetrics("Sleep", float64(sleepSummary.Summary.Stages.Deep), tags))
	tags = []string{"stage:light", nameTag}
	check(sendFloatDDMetrics("Sleep", float64(sleepSummary.Summary.Stages.Light), tags))
	tags = []string{"stage:rem", nameTag}
	check(sendFloatDDMetrics("Sleep", float64(sleepSummary.Summary.Stages.Rem), tags))
	tags = []string{"stage:wake", nameTag}
	check(sendFloatDDMetrics("Sleep", float64(sleepSummary.Summary.Stages.Wake), tags))

	tags = []string{nameTag}
	check(sendFloatDDMetrics("TotalMinutesAsleep", float64(sleepSummary.Summary.TotalMinutesAsleep), tags))
	check(sendFloatDDMetrics("TotalTimeInBed", float64(sleepSummary.Summary.TotalTimeInBed), tags))

	fmt.Println("Summary ", sleepSummary.Summary)
	fmt.Println("TotalMinutesAsleep ", sleepSummary.Summary.TotalMinutesAsleep)
	fmt.Println("TotalTimeInBed ", sleepSummary.Summary.TotalTimeInBed)

	// weight
	byteArray, err = getFitbitData(weightEndpoint, token)
	if err != nil {
		log.Fatal(err)
	}
	var weightSummary WeightSummary
	if err := json.Unmarshal(byteArray, &weightSummary); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(weightSummary.Weight) != 0 {
		fmt.Println("Weight ", weightSummary.Weight[0].Weight)
		fmt.Println("BMI ", weightSummary.Weight[0].Bmi)

		check(sendFloatDDMetrics("Weight", weightSummary.Weight[0].Weight, tags))
		check(sendFloatDDMetrics("BMI", weightSummary.Weight[0].Bmi, tags))
	}
	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
