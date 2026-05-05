// Package main is a test application for Databricks workspace acceptance tests.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type vcapServices struct {
	DatabricksWorkspace []struct {
		Credentials struct {
			Host      string `json:"databricks_host"`
			Token     string `json:"databricks_token"`
			ClusterID string `json:"cluster_id"`
		} `json:"credentials"`
	} `json:"csb-databricks-workspace"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		creds, err := getCredentials()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Databricks Host: %s, Cluster ID: %s", creds.Host, creds.ClusterID)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func getCredentials() (struct {
	Host      string
	Token     string
	ClusterID string
}, error) {
	var result struct {
		Host      string
		Token     string
		ClusterID string
	}

	raw := os.Getenv("VCAP_SERVICES")
	if raw == "" {
		return result, fmt.Errorf("VCAP_SERVICES not set")
	}

	var svc vcapServices
	if err := json.Unmarshal([]byte(raw), &svc); err != nil {
		return result, fmt.Errorf("parsing VCAP_SERVICES: %w", err)
	}

	if len(svc.DatabricksWorkspace) == 0 {
		return result, fmt.Errorf("no csb-databricks-workspace binding found")
	}

	creds := svc.DatabricksWorkspace[0].Credentials
	result.Host = creds.Host
	result.Token = creds.Token
	result.ClusterID = creds.ClusterID
	return result, nil
}
