/*
Copyright 2019 Banzai Cloud.

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

package gateways

import (
	"fmt"

	istiov1beta1 "github.com/banzaicloud/istio-operator/pkg/apis/operator/v1beta1"
	"github.com/banzaicloud/istio-operator/pkg/k8sutil"
	"github.com/banzaicloud/istio-operator/pkg/resources"
	"github.com/banzaicloud/istio-operator/pkg/util"
	"github.com/go-logr/logr"
	"github.com/goph/emperror"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler struct {
	resources.Reconciler
}

func New(client client.Client, istio *istiov1beta1.Config) *Reconciler {
	return &Reconciler{
		Reconciler: resources.Reconciler{
			Client: client,
			Owner:  istio,
		},
	}
}

func (r *Reconciler) Reconcile(log logr.Logger) error {
	var rsv = []resources.ResourceVariation{
		r.serviceAccount,
		r.clusterRole,
		r.clusterRoleBinding,
		r.deployment,
		r.service,
		r.horizontalPodAutoscaler,
	}
	for _, res := range append(resources.ResolveVariations("ingressgateway", rsv), resources.ResolveVariations("egressgateway", rsv)...) {
		o := res(r.Owner)
		err := k8sutil.Reconcile(log, r.Client, o)
		if err != nil {
			return emperror.WrapWith(err, "failed to reconcile resource", "resource", o.GetObjectKind().GroupVersionKind())
		}
	}
	return nil
}

func serviceAccountName(gw string) string {
	return fmt.Sprintf("istio-%s-service-account", gw)
}

func clusterRoleName(gw string) string {
	return fmt.Sprintf("istio-%s-cluster-role", gw)
}

func clusterRoleBindingName(gw string) string {
	return fmt.Sprintf("istio-%s-cluster-role-binding", gw)
}

func gatewayName(gw string) string {
	return fmt.Sprintf("istio-%s", gw)
}

func hpaName(gw string) string {
	return fmt.Sprintf("istio-%s-autoscaler", gw)
}

func gwLabels(gw string) map[string]string {
	return map[string]string{
		"app": fmt.Sprintf("istio-%s", gw),
	}
}

func labelSelector(gw string) map[string]string {
	return util.MergeLabels(gwLabels(gw), map[string]string{
		"istio": gw,
	})
}