package webhook

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	once sync.Once
	ws   *webHookServer
	err  error
)

var (
	ignoredNamespaces = []string{ // 忽略某些ns
		metav1.NamespaceSystem,
		metav1.NamespacePublic,
	}

	//禁用注释: "y", "yes", "true", "on":
	admissionWebhookAnnotationDiableKey = "webhook-demo.com/diable"

	admissionWebhookAnnotationRequired = map[string]string{
		"kubernetes.io/ingress-bandwidth": "300M",
		"kubernetes.io/egress-bandwidth":  "300M",
	}
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
	defaulter     = runtime.ObjectDefaulter(runtimeScheme)
)

func init() {
	_ = corev1.AddToScheme(runtimeScheme)
	_ = admissionregistrationv1beta1.AddToScheme(runtimeScheme)
	_ = v1.AddToScheme(runtimeScheme)
}

func NewWebhookServer(webHook WebHookServerParameters) (WebHookServerInt, error) {
	once.Do(func() {
		ws, err = newWebHookServer(webHook)
	})
	return ws, err
}

func newWebHookServer(webHook WebHookServerParameters) (*webHookServer, error) {
	// load tls cert/key file
	tlsCertKey, err := tls.LoadX509KeyPair(webHook.CertFile, webHook.KeyFile)
	if err != nil {
		return nil, err
	}

	ws := &webHookServer{
		server: &http.Server{
			Addr:      fmt.Sprintf(":%v", webHook.Port),
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{tlsCertKey}},
		},
	}

	// add routes
	mux := http.NewServeMux()
	mux.HandleFunc("/mutating", ws.serve)
	mux.HandleFunc("/validating", ws.serve)
	ws.server.Handler = mux
	return ws, nil
}

func (ws *webHookServer) Start() {
	if err := ws.server.ListenAndServeTLS("", ""); err != nil {
		glog.Errorf("Failed to listen and serve webhook server: %v", err)
	}
}

func (ws *webHookServer) Stop() {
	glog.Infof("Got OS shutdown signal, shutting down wenhook server gracefully...")
	ws.server.Shutdown(context.Background())
}

// validate pods
func (whsvr *webHookServer) validating(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	req := ar.Request
	var (
		availableAnnotations            map[string]string
		objectMeta                      *metav1.ObjectMeta
		resourceNamespace, resourceName string
	)

	glog.Infof("AdmissionReview for Kind=%v, Namespace=%v Name=%v (%v) UID=%v patchOperation=%v UserInfo=%v",
		req.Kind, req.Namespace, req.Name, resourceName, req.UID, req.Operation, req.UserInfo)

	switch req.Kind.Kind {
	case "Pod":
		var pod corev1.Pod
		if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
			glog.Errorf("Could not unmarshal raw object: %v", err)
			return &v1beta1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		resourceName, resourceNamespace, objectMeta = pod.Name, pod.Namespace, &pod.ObjectMeta
		availableAnnotations = pod.Annotations
	//其他不支持的类型
	default:
		msg := fmt.Sprintf("\nNot support for this Kind of resource  %v", req.Kind.Kind)
		glog.Errorf(msg)
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: msg,
			},
		}
	}

	if !admissionRequired(ignoredNamespaces, objectMeta) {
		glog.Infof("Skipping validation for %s/%s due to policy check", resourceNamespace, resourceName)
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}

	allowed := true
	var result *metav1.Status
	glog.Info("available annotations:", availableAnnotations)             //实际的anntations
	glog.Info("required annotations", admissionWebhookAnnotationRequired) //需要的anntations
	for _, rl := range admissionWebhookAnnotationRequired {
		if _, ok := availableAnnotations[rl]; !ok {
			allowed = false
			result = &metav1.Status{
				Reason: "required template annotations are not set",
			}
			break
		}
	}

	return &v1beta1.AdmissionResponse{
		Allowed: allowed,
		Result:  result,
	}
}

// main mutation process
func (whsvr *webHookServer) mutating(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	req := ar.Request
	var (
		availableAnnotations            map[string]string
		objectMeta                      *metav1.ObjectMeta
		resourceNamespace, resourceName string
	)

	glog.Infof("AdmissionReview for Kind=%v, Namespace=%v Name=%v (%v) UID=%v patchOperation=%v UserInfo=%v",
		req.Kind, req.Namespace, req.Name, resourceName, req.UID, req.Operation, req.UserInfo)

	switch req.Kind.Kind {
	case "Pod":
		var pod corev1.Pod
		if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
			glog.Errorf("Could not unmarshal raw object: %v", err)
			return &v1beta1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		resourceName, resourceNamespace, objectMeta = pod.Name, pod.Namespace, &pod.ObjectMeta
		availableAnnotations = pod.Annotations
	//其他不支持的类型
	default:
		msg := fmt.Sprintf("\nNot support for this Kind of resource  %v", req.Kind.Kind)
		glog.Errorf(msg)
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: msg,
			},
		}
	}

	if !admissionRequired(ignoredNamespaces, objectMeta) {
		glog.Infof("Skipping validation for %s/%s due to policy check", resourceNamespace, resourceName)
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}

	patchBytes, err := createPatch(availableAnnotations, admissionWebhookAnnotationRequired)
	if err != nil {
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	glog.Infof("AdmissionResponse: patch=%v\n", string(patchBytes))
	return &v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}

// Serve method for webhook server
func (whsvr *webHookServer) serve(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		glog.Error("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("Content-Type=%s, expect application/json", contentType)
		http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	var admissionResponse *v1beta1.AdmissionResponse
	ar := v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		glog.Errorf("Can't decode body: %v", err)
		admissionResponse = &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	} else {
		fmt.Println(r.URL.Path)
		if r.URL.Path == "/mutating" {
			admissionResponse = whsvr.mutating(&ar)
		} else if r.URL.Path == "/validating" {
			admissionResponse = whsvr.validating(&ar)
		}
	}

	admissionReview := v1beta1.AdmissionReview{}
	if admissionResponse != nil {
		admissionReview.Response = admissionResponse
		if ar.Request != nil {
			admissionReview.Response.UID = ar.Request.UID
		}
	}

	resp, err := json.Marshal(admissionReview)
	if err != nil {
		glog.Errorf("Can't encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
	glog.Infof("Ready to write reponse ...")
	if _, err := w.Write(resp); err != nil {
		glog.Errorf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}

// 验证是否需要跳过：忽略的ns，存在diable的注释
func admissionRequired(ignoredList []string, metadata *metav1.ObjectMeta) bool {
	// skip special kubernetes system namespaces
	for _, namespace := range ignoredList {
		if metadata.Namespace == namespace {
			glog.Infof("Skip validation for %v for it's in special namespace:%v", metadata.Name, metadata.Namespace)
			return false
		}
	}

	annotations := metadata.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	var required bool
	switch strings.ToLower(annotations[admissionWebhookAnnotationDiableKey]) {
	default:
		required = true
	case "y", "yes", "true", "on":
		glog.Infof("Skip validation policy for %v/%v: required:%v", metadata.Namespace, metadata.Name, required)
		required = false
	}
	return required
}

func createPatch(target map[string]string, added map[string]string) ([]byte, error) {
	var patch []patchOperation

	values := target
	for key, value := range added {
		if target == nil || target[key] == "" {
			values[key] = value
		}
	}
	patch = append(patch, patchOperation{
		Op:    "add",
		Path:  "/spec/template/metadata/labels",
		Value: values,
	})

	return json.Marshal(patch)
}
