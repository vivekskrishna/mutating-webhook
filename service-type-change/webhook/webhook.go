package hook

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/common/log"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io

// SidecarInjector annotates Pods
type SidecarInjector struct {
	Name          string
	Client        client.Client
	decoder       *admission.Decoder
	SidecarConfig *Config
}

type Config struct {
	Containers []corev1.Container `yaml:"containers"`
}

func shoudInject(pod *corev1.Pod) bool {
	shouldInjectSidecar, err := strconv.ParseBool(pod.Annotations["inject-logging-sidecar"])

	if err != nil {
		shouldInjectSidecar = false
	}

	if shouldInjectSidecar {
		alreadyUpdated, err := strconv.ParseBool(pod.Annotations["logging-sidecar-added"])

		if err == nil && alreadyUpdated {
			shouldInjectSidecar = false
		}
	}

	log.Info("Should Inject: ", shouldInjectSidecar)

	return shouldInjectSidecar
}

// SidecarInjector adds an annotation to every incoming pods.
func (si *SidecarInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
	svc := &corev1.Service{}

	err := si.decoder.Decode(req, svc)
	if err != nil {
		log.Info("Sdecar-Injector: cannot decode")
		return admission.Errored(http.StatusBadRequest, err)
	}

	if svc.Labels != nil {
		for key,_ := range svc.Labels {
			if strings.Contains(key, "gateway-name") && svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
				log.Info("Changing service type from ",corev1.ServiceTypeLoadBalancer," to ",corev1.ServiceTypeClusterIP," for service ", svc.Name)
				svc.Spec.Type = corev1.ServiceTypeClusterIP
				svc.Spec.ExternalTrafficPolicy = ""
			} else {
				log.Info("Skipping service ", svc.Name," as type is ",svc.Spec.Type)
			}
		}
	}



	marshaledSvc, err := json.Marshal(svc)



	if err != nil {
		log.Info("Service: cannot marshal")
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledSvc)
}

// SidecarInjector implements admission.DecoderInjector.
// A decoder will be automatically inj1ected.

// InjectDecoder injects the decoder.
func (si *SidecarInjector) InjectDecoder(d *admission.Decoder) error {
	si.decoder = d
	return nil
}
