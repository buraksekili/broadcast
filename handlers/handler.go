package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/buraksekili/broadcast/pkg"
)

type Agent struct {
	Log    *log.Logger
	Client *pkg.BroadcastClient
}

// NewHTTPAgent creates a new HTTP Agent for HTTP requests
func NewHTTPAgent() (*Agent, error) {
	conf, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return nil, err
	}

	return &Agent{
		Log:    log.New(os.Stdout, "tyk-broadcast ", log.LstdFlags),
		Client: pkg.NewBroadcastK8s(clientset, conf),
	}, nil
}

func (a *Agent) Broadcast(w http.ResponseWriter, r *http.Request) {
	req := r.Clone(context.Background())

	if err := a.Client.ListPods(context.Background(), req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("failed to broadcast", err)
		return
	}
}
