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

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	datadoghqcomv1beta1 "github.com/max-rocket-internet/datadog-controller/api/v1beta1"
	"github.com/max-rocket-internet/datadog-controller/datadog"
	"github.com/max-rocket-internet/datadog-controller/utils"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DatadogMonitorReconciler reconciles a DatadogMonitor object
type DatadogMonitorReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Datadog  datadog.Datadog
}

const (
	deletionFinalizer = "datadogmonitors.finalizers.datadoghq.com"
)

// +kubebuilder:rbac:groups=datadoghq.com,resources=datadogmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=datadoghq.com,resources=datadogmonitors/status,verbs=get;update;patch

func (r *DatadogMonitorReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("monitor", req.NamespacedName)

	instance := &datadoghqcomv1beta1.DatadogMonitor{}

	log.V(1).Info("Getting resource from cluster")
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if instance.Status.Id == 0 {
			log.Info("Creating monitor")
			monitorId, err := r.Datadog.CreateMonitor(instance.Spec)

			if err != nil {
				log.Error(err, "Monitor failed to create")

				r.Recorder.Eventf(instance, "Warning", "FailedCreate", fmt.Sprint(err))
				instance.Status.Status = "FailedCreate"
				instance.Status.ObservedGeneration = instance.ObjectMeta.Generation + 1

				if err = r.Update(ctx, instance); err != nil {
					log.Error(err, "Failed to update status after failed monitor creation")
					return ctrl.Result{}, err
				}

				return ctrl.Result{}, err
			} else {
				log.V(1).Info(fmt.Sprintf("Monitor created with ID %v", monitorId))
				r.Recorder.Eventf(instance, "Normal", "SuccessfulCreate", fmt.Sprintf("Monitor created with ID %v", monitorId))

				instance.Status.Id = monitorId
				instance.Status.Url = fmt.Sprintf("https://app.%v/monitors/%v", r.Datadog.Conf.DatadogHost, monitorId)
				instance.Status.Status = "Created"
				instance.Status.ObservedGeneration = instance.ObjectMeta.Generation + 1

				if err = r.Update(ctx, instance); err != nil {
					log.Error(err, "Failed to update status after monitor creation")
					return ctrl.Result{}, err
				}
			}
		} else if instance.ObjectMeta.Generation != instance.Status.ObservedGeneration {
			log.Info("Updating monitor")

			err := r.Datadog.UpdateMonitor(instance.Status.Id, instance.Spec)

			if err != nil {
				log.Error(err, "Monitor update failed")

				r.Recorder.Eventf(instance, "Warning", "FailedUpdate", fmt.Sprint(err))
				instance.Status.Status = "FailedUpdate"
				instance.Status.ObservedGeneration = instance.ObjectMeta.Generation + 1

				if err = r.Update(ctx, instance); err != nil {
					log.Error(err, "Failed to update status after failed monitor update")
					return ctrl.Result{}, err
				}

				return ctrl.Result{}, err
			} else {
				log.V(1).Info(fmt.Sprintf("Monitor updated with ID %v", instance.Status.Id))
				r.Recorder.Eventf(instance, "Normal", "SuccessfulUpdate", fmt.Sprintf("Monitor updated with ID %v", instance.Status.Id))

				instance.Status.Status = "Updated"
				instance.Status.ObservedGeneration = instance.ObjectMeta.Generation + 1

				if err = r.Update(ctx, instance); err != nil {
					log.Error(err, "Failed to update status after monitor update")
					return ctrl.Result{}, err
				}
			}

		} else {
			log.V(1).Info("Skipping as generation is not new")
		}

		if !utils.ContainsString(instance.ObjectMeta.Finalizers, deletionFinalizer) {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, deletionFinalizer)
			log.V(1).Info("Adding finalizer")
			if err := r.Update(context.Background(), instance); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		log.V(1).Info("Deleting monitor")
		if utils.ContainsString(instance.ObjectMeta.Finalizers, deletionFinalizer) {
			if instance.Status.Id == 0 {
				log.V(1).Info("Skipping deletion as monitor was never created")
				return ctrl.Result{}, nil
			}

			if err := r.Datadog.DeleteMonitor(instance.Status.Id); err != nil {
				log.Error(err, "Failed to delete Monitor from datadog")
				return ctrl.Result{}, err
			}

			log.Info("Deleted monitor")

			instance.ObjectMeta.Finalizers = utils.RemoveString(instance.ObjectMeta.Finalizers, deletionFinalizer)
			log.V(1).Info("Removing finalizer")
			if err := r.Update(context.Background(), instance); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (r *DatadogMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&datadoghqcomv1beta1.DatadogMonitor{}).
		Complete(r)
}
