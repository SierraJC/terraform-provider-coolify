package api_test

import (
	"context"
	"log"
	"net/http"
	"testing"

	"terraform-provider-coolify/internal/api"
)

func TestClient_canCall(t *testing.T) {
	// custom HTTP client
	hc := http.Client{}

	// with a raw http.Response
	{
		c, err := api.NewClient("http://192.168.0.4:8001/api/v1", api.WithHTTPClient(&hc))

		if err != nil {
			log.Fatal(err)
		}

		resp, err := c.Version(context.TODO())

		// resp, err := c.GetClient(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Expected HTTP 200 but received %d", resp.StatusCode)
		}
	}

	// or to get a struct with the parsed response body
	{
		// c, err := api.NewClientWithResponses("http://localhost:1234", api.WithHTTPClient(&hc))
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// resp, err := c.GetClientWithResponse(context.TODO())
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// if resp.StatusCode() != http.StatusOK {
		// 	log.Fatalf("Expected HTTP 200 but received %d", resp.StatusCode())
		// }

		// fmt.Printf("resp.JSON200: %v\n", resp.JSON200)
	}

}
