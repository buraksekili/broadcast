package pkg

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"k8s.io/apimachinery/pkg/labels"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	GWPortEnvKey   = "TYK_BROADCAST_GW_LISTENPORT"
	GWNamespaceKey = "TYK_BROADCAST_NS"
)

type BroadcastClient struct {
	Client *kubernetes.Clientset
	Config *rest.Config
}

func NewBroadcastK8s(cl *kubernetes.Clientset, conf *rest.Config) *BroadcastClient {
	return &BroadcastClient{Client: cl, Config: conf}
}

func (b *BroadcastClient) ListPods(ctx context.Context, r *http.Request) error {
	ns := os.Getenv(GWNamespaceKey)

	ls := metav1.LabelSelector{MatchLabels: map[string]string{"app": "gateway-tyk-headless"}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(ls.MatchLabels).String(),
	}

	pl, err := b.Client.CoreV1().Pods(ns).List(ctx, listOptions)
	if err != nil {
		return err
	}

	bodyByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("failed to read body", err)
		return err
	}

	for _, p := range pl.Items {
		gp := os.Getenv(GWPortEnvKey)
		if gp == "" {
			gp = "8080"
		}

		req := r.Clone(ctx)
		req.URL.Host = fmt.Sprintf("%s:%v", p.Status.PodIP, gp)
		req.URL.Scheme = "http"
		req.RequestURI = ""
		req.Body = ioutil.NopCloser(bytes.NewReader(bodyByte))
		req.ContentLength = int64(len(string(bodyByte))) // ???

		c := &http.Client{}

		res, err := c.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		_, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
	}

	return nil
}
