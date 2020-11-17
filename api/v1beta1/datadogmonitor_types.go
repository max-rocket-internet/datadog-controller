/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type DatadogMonitorOptions struct {
	EscalationMessage string `json:"escalation_message,omitempty"`
	// Time (in seconds) to delay evaluation, as a non-negative integer. For example, if the value is set to `300` (5min), the timeframe is set to `last_5m` and the time is 7:00, the monitor evaluates data from 6:50 to 6:55. This is useful for AWS CloudWatch and other backfilled metrics to ensure the monitor always has data during evaluation.
	EvaluationDelay int64 `json:"evaluation_delay,omitempty"`
	// A Boolean indicating whether notifications from this monitor automatically inserts its triggering tags into the title.  **Examples** - If `True`, `[Triggered on {host:h1}] Monitor Title` - If `False`, `[Triggered] Monitor Title`
	IncludeTags bool `json:"include_tags,omitempty"`
	// Whether or not the monitor is locked (only editable by creator and admins).
	Locked bool `json:"locked,omitempty"`
	// How long the test should be in failure before alerting (integer, number of seconds, max 7200).
	MinFailureDuration int64 `json:"min_failure_duration,omitempty"`
	// The minimum number of locations in failure at the same time during at least one moment in the `min_failure_duration` period (`min_location_failed` and `min_failure_duration` are part of the advanced alerting rules - integer, >= 1).
	MinLocationFailed int64 `json:"min_location_failed,omitempty"`
	// Time (in seconds) to allow a host to boot and applications to fully start before starting the evaluation of monitor results. Should be a non negative integer.
	NewHostDelay int64 `json:"new_host_delay,omitempty"`
	// The number of minutes before a monitor notifies after data stops reporting. Datadog recommends at least 2x the monitor timeframe for metric alerts or 2 minutes for service checks. If omitted, 2x the evaluation timeframe is used for metric alerts, and 24 hours is used for service checks.
	NoDataTimeframe int64 `json:"no_data_timeframe,omitempty"`
	// A Boolean indicating whether tagged users is notified on changes to this monitor.
	NotifyAudit bool `json:"notify_audit,omitempty"`
	// A Boolean indicating whether this monitor notifies when data stops reporting.
	NotifyNoData bool `json:"notify_no_data,omitempty"`
	// The number of minutes after the last notification before a monitor re-notifies on the current status. It only re-notifies if it’s not resolved.
	RenotifyInterval int64 `json:"renotify_interval,omitempty"`
	// A Boolean indicating whether this monitor needs a full window of data before it’s evaluated. We highly recommend you set this to `false` for sparse metrics, otherwise some evaluations are skipped. Default is false.
	RequireFullWindow bool                     `json:"require_full_window,omitempty"`
	Thresholds        DatadogMonitorThresholds `json:"thresholds,omitempty"`
	// The number of hours of the monitor not reporting data before it automatically resolves from a triggered state.
	TimeoutH int64 `json:"timeout_h,omitempty"`
}

type DatadogMonitorThresholds struct {
	// The monitor `CRITICAL` threshold.
	// +kubebuilder:validation:Type=number
	Critical float64 `json:"critical,omitempty"`
	// The monitor `CRITICAL` recovery threshold.
	// +kubebuilder:validation:Type=number
	CriticalRecovery float64 `json:"critical_recovery,omitempty"`
	// The monitor `OK` threshold.
	// +kubebuilder:validation:Type=number
	Ok float64 `json:"ok,omitempty"`
	// The monitor UNKNOWN threshold.
	// +kubebuilder:validation:Type=number
	Unknown float64 `json:"unknown,omitempty"`
	// The monitor `WARNING` threshold.
	// +kubebuilder:validation:Type=number
	Warning float64 `json:"warning,omitempty"`
	// The monitor `WARNING` recovery threshold.
	// +kubebuilder:validation:Type=number
	WarningRecovery float64 `json:"warning_recovery,omitempty"`
}

type DatadogMonitorSpec struct {
	// ID of this monitor.
	Id int64 `json:"id,omitempty"`
	// A message to include with notifications for this monitor.
	Message string `json:"message"`
	// Whether or not the monitor is broken down on different groups.
	Multi bool `json:"multi,omitempty"`
	// The monitor name.
	Name    string                `json:"name"`
	Options DatadogMonitorOptions `json:"options,omitempty"`
	// Integer from 1 (high) to 5 (low) indicating alert severity.
	Priority int64 `json:"priority,omitempty"`
	// The monitor query.
	Query string `json:"query"`
	// Tags associated to your monitor.
	Tags []string `json:"tags,omitempty"`
	// The Type of monitor it is. Must be one of: "composite", "event alert", "log alert", "metric alert", "process alert", "query alert", "rum alert", "service check", "synthetics alert", "trace-analytics alert", "slo alert"
	Type string `json:"type,omitempty"`
}

// DatadogMonitorStatus defines the observed state of DatadogMonitor
// Why don't these printcolumns work?
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`,description="The status of the monitor"
// +kubebuilder:printcolumn:name="Id",type=string,JSONPath=`.status.id`,description="The monitor ID in Datadog"
// +kubebuilder:printcolumn:name="Url",type=string,JSONPath=`.status.url`,description="The monitor URL in Datadog"

type DatadogMonitorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// True if BaseRole and Database present in the host

	// Is the monitor created in Datadog
	Status string `json:"status,omitempty"`
	// The monitor ID in Datadog
	Id int64 `json:"id,omitempty"`
	// The monitor URL in Datadog
	Url string `json:"url,omitempty"`
	// The last applied generation. Used to distinguish
	ObservedGeneration int64 `json:"observed_generation,omitempty"`
}

// +kubebuilder:object:root=true

// DatadogMonitor is the Schema for the datadogmonitors API
type DatadogMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatadogMonitorSpec   `json:"spec,omitempty"`
	Status DatadogMonitorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DatadogMonitorList contains a list of DatadogMonitor
type DatadogMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DatadogMonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DatadogMonitor{}, &DatadogMonitorList{})
}
