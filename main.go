package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
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

type RefreshResponce struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	UserID       string `json:"user_id"`
}

type Secret struct {
	AccessToken  string `json:"fitbit_token"`
	RefreshToken string `json:"fitbit_refresh_token"`
}

func main() {

	token, _ := getToken("fitbit", "ap-northeast-1")
	date := "2021-11-21.json"

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

	fmt.Println(resp.StatusCode)

	//// sleep
	//apiEndPoint = baseURL + sleepEndpoint + date
	//req, _ = http.NewRequest("GET", apiEndPoint, nil)
	//req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	//resp, err = http.DefaultClient.Do(req)
	//if err != nil {
	//	os.Exit(1)
	//}
	//defer resp.Body.Close()

	//byteArray, _ = ioutil.ReadAll(resp.Body)

	//var sleepSummary SleepSummary
	//if err := json.Unmarshal(byteArray, &sleepSummary); err != nil {
	//	os.Exit(1)
	//}
	//fmt.Println("Summary ", sleepSummary.Summary)
	//fmt.Println("TotalMinutesAsleep ", sleepSummary.Summary.TotalMinutesAsleep)
	//fmt.Println("TotalTimeInBed ", sleepSummary.Summary.TotalTimeInBed)

	//// weight
	//apiEndPoint = baseURL + weightEndpoint + date
	//req, _ = http.NewRequest("GET", apiEndPoint, nil)
	//req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	//resp, err = http.DefaultClient.Do(req)
	//if err != nil {
	//	os.Exit(1)
	//}
	//defer resp.Body.Close()

	//byteArray, _ = ioutil.ReadAll(resp.Body)
	//var weightSummary WeightSummary
	//if err := json.Unmarshal(byteArray, &weightSummary); err != nil {
	//	os.Exit(1)
	//}
	//fmt.Println("Weight ", weightSummary.Weight[0].Weight)
	//fmt.Println("BMI ", weightSummary.Weight[0].Bmi)

}

func getToken(secretName, region string) (token, refreshToken string) {

	sess, err := session.NewSession()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	svc := secretsmanager.New(sess,
		aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	var secretString string
	secretString = *result.SecretString
	var secret Secret
	if err := json.Unmarshal([]byte(secretString), &secret); err != nil {
		os.Exit(1)
	}

	return secret.AccessToken, secret.RefreshToken
}