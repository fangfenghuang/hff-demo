package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"webhook-demo/pkg/webhook"

	"github.com/golang/glog"
)

var webHook webhook.WebHookServerParameters

func main() {
	// parse parameters
	flag.Parse()

	// init webhook api
	ws, err := webhook.NewWebhookServer(webHook)
	if err != nil {
		panic(err)
	}

	// start webhook server in new routine
	go ws.Start()
	glog.Info("webhook Server started...")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	ws.Stop()
}

func init() {
	flag.IntVar(&webHook.Port, "port", 443, "The port of webhook server to listen.")
	flag.StringVar(&webHook.CertFile, "tlsCertPath", "/etc/webhook/certs/cert.pem", "The path of tls cert")
	flag.StringVar(&webHook.KeyFile, "tlsKeyPath", "/etc/webhook/certs/key.pem", "The path of tls key")
}
