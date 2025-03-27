package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResourceQuotaEnforcerSpec defines the desired state of ResourceQuotaEnforcer
type ResourceQuotaEnforcerSpec struct {
	// TargetNamespace is the namespace to monitor and enforce resource quotas for
	// +kubebuilder:validation:Required
	TargetNamespace string `json:"targetNamespace"`

	// ResourceThresholds defines the thresholds for different resources
	// +kubebuilder:validation:Required
	ResourceThresholds []ResourceThreshold `json:"resourceThresholds"`

	// Actions defines what actions to take when thresholds are reached
	// +kubebuilder:validation:Required
	Actions EnforcementActions `json:"actions"`

	// CheckIntervalSeconds defines how often to check resource usage
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=60
	CheckIntervalSeconds int `json:"checkIntervalSeconds,omitempty"`
}

// ResourceThreshold defines a threshold for a specific resource
type ResourceThreshold struct {
	// ResourceName is the name of the resource to monitor (cpu, memory, pods, etc.)
	// +kubebuilder:validation:Required
	ResourceName corev1.ResourceName `json:"resourceName"`

	// WarningThresholdPercent is the percentage of usage that triggers a warning
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:default=80
	WarningThresholdPercent int `json:"warningThresholdPercent,omitempty"`

	// CriticalThresholdPercent is the percentage of usage that triggers critical actions
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:default=90
	CriticalThresholdPercent int `json:"criticalThresholdPercent,omitempty"`
}

// EnforcementActions defines what actions to take when thresholds are reached
type EnforcementActions struct {
	// SendWarningAlert indicates whether to send warning alerts
	// +kubebuilder:default=true
	SendWarningAlert bool `json:"sendWarningAlert,omitempty"`

	// WarningAlertChannel defines where to send warning alerts (e.g., Slack webhook URL)
	// +optional
	WarningAlertChannel string `json:"warningAlertChannel,omitempty"`

	// SendCriticalAlert indicates whether to send critical alerts
	// +kubebuilder:default=true
	SendCriticalAlert bool `json:"sendCriticalAlert,omitempty"`

	// CriticalAlertChannel defines where to send critical alerts
	// +optional
	CriticalAlertChannel string `json:"criticalAlertChannel,omitempty"`

	// BlockNewDeployments indicates whether to block new deployments when critical threshold is reached
	// +kubebuilder:default=false
	BlockNewDeployments bool `json:"blockNewDeployments,omitempty"`
}

// ResourceQuotaEnforcerStatus defines the observed state of ResourceQuotaEnforcer
type ResourceQuotaEnforcerStatus struct {
	// CurrentResourceUsage shows the current resource usage percentages
	// +optional
	CurrentResourceUsage []ResourceUsage `json:"currentResourceUsage,omitempty"`

	// LastCheckedTime is the last time resource usage was checked
	// +optional
	LastCheckedTime metav1.Time `json:"lastCheckedTime,omitempty"`

	// ActiveAlerts shows the active alerts that were triggered
	// +optional
	ActiveAlerts []Alert `json:"activeAlerts,omitempty"`

	// IsBlocking indicates whether the enforcer is currently blocking new deployments
	// +optional
	IsBlocking bool `json:"isBlocking,omitempty"`
}

// ResourceUsage defines the current usage of a specific resource
type ResourceUsage struct {
	// ResourceName is the name of the resource
	ResourceName corev1.ResourceName `json:"resourceName"`

	// UsedPercentage is the percentage of the resource quota being used
	UsedPercentage int `json:"usedPercentage"`

	// CurrentValue is the current usage value
	CurrentValue string `json:"currentValue"`

	// LimitValue is the quota limit value
	LimitValue string `json:"limitValue"`
}

// Alert defines information about an active alert
type Alert struct {
	// ResourceName is the resource that triggered the alert
	ResourceName corev1.ResourceName `json:"resourceName"`

	// Severity is the severity of the alert (warning or critical)
	Severity string `json:"severity"`

	// Message is the alert message
	Message string `json:"message"`

	// TimeTriggered is when the alert was triggered
	TimeTriggered metav1.Time `json:"timeTriggered"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Target Namespace",type=string,JSONPath=`.spec.targetNamespace`
//+kubebuilder:printcolumn:name="Blocking",type=boolean,JSONPath=`.status.isBlocking`
//+kubebuilder:printcolumn:name="Last Checked",type=date,JSONPath=`.status.lastCheckedTime`

// ResourceQuotaEnforcer is the Schema for the resourcequotaenforcers API
type ResourceQuotaEnforcer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceQuotaEnforcerSpec   `json:"spec,omitempty"`
	Status ResourceQuotaEnforcerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ResourceQuotaEnforcerList contains a list of ResourceQuotaEnforcer
type ResourceQuotaEnforcerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceQuotaEnforcer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceQuotaEnforcer{}, &ResourceQuotaEnforcerList{})
}
