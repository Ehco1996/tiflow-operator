/*
Copyright 2022.

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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	pingcapcomv1alpha1 "github.com/pingcap/tiflow-operator/api/v1alpha1"
	"github.com/pingcap/tiflow-operator/pkg/standalone"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// StandaloneReconciler reconciles a Standalone object
// Notes: For test use only
type StandaloneReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Control standalone.ControlInterface
}

// NewStandaloneReconciler Notes: For test use only
func NewStandaloneReconciler(cli client.Client, scheme *runtime.Scheme) *StandaloneReconciler {
	return &StandaloneReconciler{
		Client:  cli,
		Scheme:  scheme,
		Control: standalone.NewDefaultStandaloneControl(cli, scheme),
	}
}

//+kubebuilder:rbac:groups=pingcap.com,resources=standalones,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pingcap.com,resources=standalones/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pingcap.com,resources=standalones/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Standalone object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
// Notes: For test use only
func (r *StandaloneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("start standalone reconcile logic")

	instance := &pingcapcomv1alpha1.Standalone{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("standalone instance is not found")
			return ctrl.Result{}, nil
		}
		logger.Info("find standalone instance error")
		return ctrl.Result{}, err
	}

	logger.Info("start frame standalone reconcile logic", "reconcile", "init")
	if result, err := r.Control.UpdateStandalone(ctx, instance); err != nil {
		return result, err
	}

	logger.Info("standalone reconcile success")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StandaloneReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pingcapcomv1alpha1.Standalone{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
