/*
Copyright 2026.

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

package controller

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	computealphav1 "github.com/somasundar-kapaka/cloudfusion/api/alphav1"
	"github.com/somasundar-kapaka/cloudfusion/internal/ec2i"
	"github.com/somasundar-kapaka/cloudfusion/utils"
)

// EC2InstanceReconciler reconciles a EC2Instance object
type EC2InstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=compute.cloudfusion.com,resources=ec2instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=compute.cloudfusion.com,resources=ec2instances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=compute.cloudfusion.com,resources=ec2instances/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EC2Instance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.4/pkg/reconcile
func (r *EC2InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Reconciling ec2 instance request")

	ec2Inc := &computealphav1.EC2Instance{}
	err := r.Client.Get(ctx, req.NamespacedName, ec2Inc)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("ec2 Instance not found", "name", req.Name, "namespace", req.Namespace)
			return ctrl.Result{}, err
		}
		log.Error(err, "error fetching ec2 Instance", "name", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, err
	}
	

	// Add delete functionality
	if ec2Inc.ObjectMeta.DeletionTimestamp != nil {

		// TODO: Delete func  EC2 in AWS
		// TODO: updae the object ec2i
		

		// Write delete functionality
		utils.DeleteFinalizers(&ec2Inc.Finalizers, utils.FinalizerKey)
		err = r.Client.Delete(ctx, ec2Inc)
		if err != nil {
			log.Error(err, "error deleting ec2 instance", "name", ec2Inc.Name, "namespace", ec2Inc.Namespace)
			return ctrl.Result{}, err
		}
		log.Info("ec2 instance deleted successfully", "name", ec2Inc.Name, "namespace", ec2Inc.Namespace)
		// do not reque
		return ctrl.Result{}, nil
	}

	err = ec2i.ValidteNewInstanceRequest(ec2Inc)
	if err != nil {
		log.Error(err, "Invalid ec2 instance spec", "name", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, err
	}

	// TODO: Check finalizers exits if not add one Before creating EC2: utils.ContainsFinalizers
	utils.AddFinalizers(&ec2Inc.Finalizers, utils.FinalizerKey)


	err = ec2i.CreateEC2Instance(ctx, ec2Inc)
	if err != nil {
		log.Error(err, "error creating ec2 instance", "name", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EC2InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// TODO: Explore what can a NewControllerManagedBy do ?
	return ctrl.NewControllerManagedBy(mgr).
		For(&computealphav1.EC2Instance{}).
		Named("ec2instance").
		Complete(r)
}
