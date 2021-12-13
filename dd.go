package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	datadog "github.com/DataDog/datadog-api-client-go/api/v1/datadog"
)

func sendFloatDDMetrics(name string, point float64, tags []string) error {
	body := datadog.MetricsPayload{
		Series: []datadog.Series{
			datadog.Series{
				Metric: name,
				Type:   datadog.PtrString("gauge"),
				Points: [][]*float64{
					{
						datadog.PtrFloat64(float64(time.Now().Unix())),
						datadog.PtrFloat64(point),
					},
				},
				Tags: &tags,
			},
		},
	}
	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	resp, r, err := apiClient.MetricsApi.SubmitMetrics(ctx, body, *datadog.NewSubmitMetricsOptionalParameters())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return err
	}

	responseContent, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintf(os.Stdout, "Response from `MetricsApi.SubmitMetrics`:\n%s\n", responseContent)
	return nil
}

func send() {
	body := datadog.MetricsPayload{
		Series: []datadog.Series{
			datadog.Series{
				Metric: "weight",
				Type:   datadog.PtrString("gauge"),
				Points: [][]*float64{
					{
						datadog.PtrFloat64(float64(time.Now().Unix())),
						datadog.PtrFloat64(81.4),
					},
				},
				Tags: &[]string{
					"test:ExampleSubmitmetricsreturnsPayloadacceptedresponse",
				},
			},
		},
	}
	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	resp, r, err := apiClient.MetricsApi.SubmitMetrics(ctx, body, *datadog.NewSubmitMetricsOptionalParameters())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MetricsApi.SubmitMetrics`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}

	responseContent, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintf(os.Stdout, "Response from `MetricsApi.SubmitMetrics`:\n%s\n", responseContent)
}
