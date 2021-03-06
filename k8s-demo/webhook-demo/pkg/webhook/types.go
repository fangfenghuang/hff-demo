package webhook

import (
	"net/http"

	"k8s.io/api/admission/v1beta1"
)

type WebHookServerInt interface {
	mutating(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse
	validating(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse
	Start()
	Stop()
}

// Webhook Server parameters
type WebHookServerParameters struct {
	Port     int    // webhook server port
	CertFile string // path to the x509 certificate for https
	KeyFile  string // path to the x509 private key matching `CertFile`
}

type webHookServer struct {
	server *http.Server
}

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}
