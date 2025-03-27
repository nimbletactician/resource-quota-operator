// controllers/webhook.go
package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	quotaenforcerv1 "github.com/example/resource-quota-enforcer/api/v1"
)

// ResourceQuotaEnforcerWebhook handles admission requests
type ResourceQuotaEnforcerWebhook struct {
	Client  client.Client
	Log     logr.Logger
	decoder *admission.Decoder
}

// Handle handles admission requests and blocks new resources if necessary
func (a *ResourceQuotaEnforcerWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	a.Log.Info("Processing admission request", "namespace", req.Namespace, "name", req.Name, "kind", req.Kind.Kind)

	// Only process creation requests for pods and deployments
	if req.Operation != "CREATE" {
		return admission.Allowed("")
	}

	// Look for any ResourceQuotaEnforcer that targets this namespace and is blocking
	enforcerList := &quotaenforcerv1.ResourceQuotaEnforcerList{}
	if err := a.Client.List(ctx, enforcerList); err != nil {
		a.Log.Error(err, "Failed to list ResourceQuotaEnforcers")
		return admission.Errored(http.StatusInternalServerError, err)
	}

	for _, enforcer := range enforcerList.Items {
		if enforcer.Spec.TargetNamespace == req.Namespace && enforcer.Status.IsBlocking {
			// Get the resources that triggered the critical threshold
			criticalResources := []string{}
			for _, alert := range enforcer.Status.ActiveAlerts {
				if alert.Severity == "critical" {
					criticalResources = append(criticalResources, string(alert.ResourceName))
				}
			}

			msg := fmt.Sprintf("Resource creation blocked: namespace %s has reached resource quota critical threshold for: %v", 
				req.Namespace, criticalResources)
			a.Log.Info(msg)
			return admission.Denied(msg)
		}
	}

	return admission.Allowed("")
}

// InjectDecoder injects the decoder into the webhook
func (a *ResourceQuotaEnforcerWebhook) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
