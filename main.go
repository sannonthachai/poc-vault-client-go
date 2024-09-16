package main

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/vault-client-go"
)

func initVaultClient() *vault.Client {
	// prepare a client with the given base address
	client, err := vault.New(
		vault.WithAddress("http://192.168.12.105:8200"),
		vault.WithRequestTimeout(10*time.Second),
	)

	if err != nil {
		log.Fatal(err)
	}

	return client
}

func main() {
	cl := initVaultClient()

	usingToken(cl)

	listPath := []string{
		"/kv/data/test-vault-client-go",
		"/kv/data/test-vault-client-go-2",
	}

	for _, s := range listPath {
		data := read(cl, s)

		data = transform(data)

		write(cl, s, data)
	}
}

func usingToken(cl *vault.Client) {
	err := cl.SetToken("") // THIS SHOULD NOT BE HARDCODED IN REAL APP, use os.Getenv() as example of more secured way of storing the token
	if err != nil {
		log.Fatal(err)
	}
}

func write(cl *vault.Client, path string, data map[string]interface{}) {
	_, err := cl.Write(context.Background(), path, map[string]interface{}{
		"data": data,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("secret written successfully")
}

func read(cl *vault.Client, path string) map[string]interface{} {
	resp, err := cl.Read(context.Background(), path)
	if err != nil {
		log.Fatal(err)
	}

	data, ok := resp.Data["data"].(map[string]interface{})
	if !ok {
		log.Fatal("not map interface")
	}

	return data
}

func transform(data map[string]interface{}) map[string]interface{} {
	patch := map[string]interface{}{
		"OpenTelemetry__UseConsoleExporter": "false",
		"OpenTelemetry__EndpointLogs":       "http://192.168.16.81:4317",
		"OpenTelemetry__EndpointMetrics":    "http://192.168.16.81:4317",
		"OpenTelemetry__EndpointTraces":     "http://192.168.16.81:4317",
		"update_test":                       "abcde",
	}

	for k, v := range data {
		if _, ok := patch[k]; !ok {
			patch[k] = v
		}
	}

	return patch
}
