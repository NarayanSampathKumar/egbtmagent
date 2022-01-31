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
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	mwc "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	apmv1 "github.com/NarayanSampathKumar/egbtmagent/api/v1"
)

// BtmAgentReconciler reconciles a BtmAgent object
type BtmAgentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apm.egapm,resources=btmagents,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apm.egapm,resources=btmagents/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apm.egapm,resources=btmagents/finalizers,verbs=update
//+kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=mutatingwebhookconfigurations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete;
//+kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete;

//////kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update;patch;delete;
//////kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=get;list;watch;create;update;patch;delete
//////kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=get;list;watch;create;update;patch;delete

func (r *BtmAgentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)
	log.Info("1. Reconcile method is called !!!")
	BtmAgent := &apmv1.BtmAgent{}
	err := r.Get(ctx, req.NamespacedName, BtmAgent)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("BtmAgent resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get BtmAgent")
		return ctrl.Result{}, err
	}
	// Check if the namespace already exists, if not create a new one
	//	foundns := &corev1.Namespace{}
	//	BtmAgentNamespace := r.getNamespaceForBtmAgent(BtmAgent)
	//	err = r.Get(ctx, types.NamespacedName{Name: BtmAgentNamespace.Name, Namespace: ""}, foundns)
	//	if err != nil && errors.IsNotFound(err) {
	//		//Define a new Namespace
	//		log.Info("Creating a new Namespace", "Namespace.Name", BtmAgentNamespace.Name)
	//		err = r.Create(ctx, BtmAgentNamespace)
	//		if err != nil {
	//			log.Error(err, "Failed to Create new Namespace", "Namespace.Name", BtmAgentNamespace.Name)
	//			return ctrl.Result{}, err
	//		}
	//		// Namespace Created Successfully - return and reque
	//		return ctrl.Result{Requeue: true}, nil
	//	} else if err != nil {
	//		log.Error(err, "Failed to get Namespace")
	//		return ctrl.Result{}, err
	//	}
	// Check if the serviceaccount already exists, if not create a new one

	foundsa := &corev1.ServiceAccount{}
	BtmAgentServiceaccount := r.getServiceaccountForBtmAgent(BtmAgent)
	err = r.Get(ctx, types.NamespacedName{Name: BtmAgentServiceaccount.Name, Namespace: BtmAgentServiceaccount.Namespace}, foundsa)
	if err != nil && errors.IsNotFound(err) {
		//Define a new ServiceAccount
		log.Info("In line97-Create new SA", "ServiceAccount.Namespace", BtmAgentServiceaccount.Namespace, "Namespace.Name", BtmAgentServiceaccount.Name)
		err = r.Create(ctx, BtmAgentServiceaccount)
		if err != nil {
			log.Error(err, "Fail to Create SA", "SAcc.Namespace", BtmAgentServiceaccount.Namespace, "Ns.Name", BtmAgentServiceaccount.Name)
			return ctrl.Result{}, err
		}
		// Serviceaccount Created Successfully - return and reque
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Serviceaccount")
		return ctrl.Result{}, err
	}
	// Check if the daemonset already exists, if not create a new one
	found := &appsv1.DaemonSet{}
	BtmAgentDaemonset := r.getDaemonsetForBtmAgent(BtmAgent)
	err = r.Get(ctx, types.NamespacedName{Name: BtmAgentDaemonset.Name, Namespace: BtmAgentDaemonset.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new daemonset
		log.Info("Creating a new Daemonset", "Daemonset.Namespace", BtmAgentDaemonset.Namespace, "Daemonset.Name", BtmAgentDaemonset.Name)
		err = r.Create(ctx, BtmAgentDaemonset)
		if err != nil {
			log.Error(err, "Failed-create Daemonset", "Daemonset.Namespace", BtmAgentDaemonset.Namespace, "Daemonset.Name", BtmAgentDaemonset.Name)
			return ctrl.Result{}, err
		}
		// Daemonset created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Daemonset")
		return ctrl.Result{}, err
	}
	// Check if ClusterRole already exists, if not create a new one
	//	foundClusterRole := &rbacv1.ClusterRole{}
	//	BtmAgentClusterRole := r.getClusterRoleForBtmAgent(BtmAgent)
	//	err = r.Get(ctx, types.NamespacedName{Name: BtmAgentClusterRole.Name, Namespace: ""}, foundClusterRole)
	//	if err != nil && errors.IsNotFound(err) {
	//		// Define a new Cluster Role
	//		log.Info("Creating a new ClusterRole", "BTMAgentClusterRole.Name", BtmAgentClusterRole.Name)
	//		err = r.Create(ctx, BtmAgentClusterRole)
	//		if err != nil {
	//			log.Error(err, "Failed-Create -> BtmAgentClusterRole", "BtmAgentClusterRole.Name", BtmAgentClusterRole.Name)
	//			return ctrl.Result{}, err
	//		}
	//		// Cluster Role Created Successfully - return and requeue
	//		return ctrl.Result{Requeue: true}, nil
	//	} else if err != nil {
	//		log.Error(err, "failed to get ClusterRole")
	//		return ctrl.Result{}, err
	//	}
	//	// Check if ClusterRoleBinding already exists, if not create a new one
	//	foundClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
	//	BtmAgentClusterRoleBinding := r.getClusterRoleBindingForBtmAgent(BtmAgent)
	//	err = r.Get(ctx, types.NamespacedName{Name: BtmAgentClusterRoleBinding.Name, Namespace: ""}, foundClusterRoleBinding)
	//	if err != nil && errors.IsNotFound(err) {
	//		// Define a new Cluster Role Binding
	//		log.Info("Creating a new ClusterRoleBinding", "BTMAgentClusterRoleBinding.Name", BtmAgentClusterRoleBinding.Name)
	//		err = r.Create(ctx, BtmAgentClusterRoleBinding)
	//		if err != nil {
	//			log.Error(err, "Failed-Create -> BtmAgentClusterRoleBinding", "BtmAgentClusterRoleBinding.Name", BtmAgentClusterRoleBinding.Name)
	//			return ctrl.Result{}, err
	//		}
	//		// Cluster Role Binding Created Successfully - return and requeue
	//		return ctrl.Result{Requeue: true}, nil
	//	} else if err != nil {
	//		log.Error(err, "failed to get ClusterRoleBinding")
	//		return ctrl.Result{}, err
	//	}
	// Check if the mutatingWebhookConfiguration already exists, if not create a new one
	//    foundMwc := &mwc.MutatingWebhookConfiguration{}
	foundmwc := &mwc.MutatingWebhookConfigurationList{}
	mutatingwebhookConfig := r.getMutatingWebhookConfigurationForBTM(BtmAgent)
	err = r.Get(ctx, types.NamespacedName{Name: mutatingwebhookConfig.Name, Namespace: ""}, foundmwc)
	log.Info("values of BtmAgent", "BtmAgent.MonitoredNamespaces", BtmAgent.MonitoredNamespaces, "BtmAgent.UnMonitoredNamespaces", BtmAgent.UnMonitoredNamespaces, "BtmAgent.matchingLabels", BtmAgent.MatchingLabels)

	// if err != nil && errors.IsNotFound(err) {
	// 	//Define a new Namespace
	// 	log.Info("Creating a new mwc", "mwc.Name", mutatingwebhookConfig.Name)
	// 	err = r.Create(ctx, mutatingwebhookConfig)
	// 	if err != nil {
	// 		log.Error(err, "Failed to Create new mwc", "mwc.Name", mutatingwebhookConfig.Name)
	// 		return ctrl.Result{}, err
	// 	}
	// 	// Namespace Created Successfully - return and reque
	// 	return ctrl.Result{Requeue: true}, nil
	// } else if err != nil {
	// 	log.Error(err, "Failed to get Namespace")
	// 	return ctrl.Result{}, err
	// }

	return ctrl.Result{}, nil
}
func (r *BtmAgentReconciler) getMutatingWebhookConfigurationForBTM(m *apmv1.BtmAgent) *mwc.MutatingWebhookConfiguration {
	mwcPath := "/mutate-core-v1-deployment"
	var sideeffects mwc.SideEffectClass = "None"
	var failurePolicy mwc.FailurePolicyType = "Fail"
	var createoperation mwc.OperationType = "CREATE"
	var updateoperation mwc.OperationType = "UPDATE"
	BtmAgentMwc := &mwc.MutatingWebhookConfiguration{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admissionregistration.k8s.io/v1",
			Kind:       "MutatingWebhookConfiguration",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "mdeployment.kb.io",
		},
		Webhooks: []mwc.MutatingWebhook{
			{
				Name:          "mdeployment.kb.io",
				FailurePolicy: &failurePolicy,
				SideEffects:   &sideeffects,
				AdmissionReviewVersions: []string{
					"v1", "v1beta1",
				},
				ClientConfig: mwc.WebhookClientConfig{
					Service: &mwc.ServiceReference{
						Name:      "webhook-service",
						Namespace: "system",
						Path:      &mwcPath,
					},
				},
				Rules: []mwc.RuleWithOperations{
					{
						Rule: mwc.Rule{
							APIGroups:   []string{"apps"},
							APIVersions: []string{"v1"},
							Resources:   []string{"deployments"},
						},
						Operations: []mwc.OperationType{
							createoperation, updateoperation,
						},
					},
				},
			},
		},
	}
	// Set BtmAgent instance as the owner and controller
	ctrl.SetControllerReference(m, BtmAgentMwc, r.Scheme)
	return BtmAgentMwc
}
func (r *BtmAgentReconciler) getClusterRoleForBtmAgent(m *apmv1.BtmAgent) *rbacv1.ClusterRole {
	btmAgentClusterRole := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRole",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "egagent-role",
		},
		Rules: []rbacv1.PolicyRule{
			{
				NonResourceURLs: []string{"/version", "/healthz"},
				Verbs:           []string{"get"},
				APIGroups:       []string{},
				Resources:       []string{},
			},
			{
				APIGroups: []string{"batch"},
				Resources: []string{"jobs", "cronjobs"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"extensions"},
				Resources: []string{"deployments", "replicasets", "ingresses"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"deployments", "replicasets", "daemonsets", "statefulsets"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces", "events", "services", "endpoints", "nodes", "pods", "replicationcontrollers", "componentstatuses", "resourcequotas", "persistentvolumes", "persistentvolumeclaims"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"endpoints"},
				Verbs:     []string{"create", "update", "patch"},
			},
			{
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingresses"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"apps.openshift.io"},
				Resources: []string{"deploymentconfigs"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups:     []string{"security.openshift.io"},
				ResourceNames: []string{"privileged"},
				Resources:     []string{"securitycontextconstraints"},
				Verbs:         []string{"use"},
			},
		},
	}
	// Set BtmAgent instance as the owner and controller
	ctrl.SetControllerReference(m, btmAgentClusterRole, r.Scheme)
	return btmAgentClusterRole
}
func (r *BtmAgentReconciler) getClusterRoleBindingForBtmAgent(m *apmv1.BtmAgent) *rbacv1.ClusterRoleBinding {
	BtmAgentClusterRoleBinding := &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRoleBinding",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "egagent-role",
			Namespace: "egagent",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "egagent",
				Namespace: "egagent",
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     "egagent-role",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	// Set BtmAgent instance as the owner and controller
	ctrl.SetControllerReference(m, BtmAgentClusterRoleBinding, r.Scheme)
	return BtmAgentClusterRoleBinding
}
func (r *BtmAgentReconciler) getNamespaceForBtmAgent(m *apmv1.BtmAgent) *corev1.Namespace {
	btmAgentNamespace := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "egagent",
		},
	}
	// Set BtmAgent instance as the owner and controller
	ctrl.SetControllerReference(m, btmAgentNamespace, r.Scheme)
	return btmAgentNamespace
}

// serviceAccountForEgAgent returns the serviceAccount for BtmAgent
func (r *BtmAgentReconciler) getServiceaccountForBtmAgent(m *apmv1.BtmAgent) *corev1.ServiceAccount {
	BtmAgentServiceAccount := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "egagent",
			Namespace: "egagent",
		},
	}
	// Set BtmAgent instance as the owner and controller
	ctrl.SetControllerReference(m, BtmAgentServiceAccount, r.Scheme)
	return BtmAgentServiceAccount
}

// daemonsetForEgAgent returns a EgAgent Daemonset object
func (r *BtmAgentReconciler) getDaemonsetForBtmAgent(m *apmv1.BtmAgent) *appsv1.DaemonSet {
	var is_privileged bool = true
	ls := map[string]string{"app": "egagent", "egagent_cr": m.Name}
	BtmAgentDaemonSet := &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "DaemonSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "egagent-narayan",
			Namespace: "egagent",
			Labels: map[string]string{
				"app": "egagent",
			},
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					HostIPC:            true,
					HostNetwork:        true,
					HostPID:            true,
					ServiceAccountName: "egagent",
					Volumes: []corev1.Volume{
						{
							Name: "var-run",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/run",
								},
							},
						},
						{
							Name: "process",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/proc",
								},
							},
						},
						{
							Name: "host-root",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/",
								},
							},
						},
						{
							Name: "btmmount",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/opt/egbtm",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:            "egagent",
							Image:           "eginnovations/agent:718",
							ImagePullPolicy: corev1.PullAlways,
							SecurityContext: &corev1.SecurityContext{
								Privileged: &is_privileged,
							},
							Env: []corev1.EnvVar{
								{
									Name:  "EG_MANAGER",
									Value: "3.141.114.182",
								},
								{
									Name:  "EG_MANAGER_PORT",
									Value: "443",
								},
								{
									Name:  "EG_MANAGER_SSL",
									Value: fmt.Sprintf("%t", true),
								},
								{
									Name:  "JVM_MEMORY",
									Value: "512",
								},
								{
									Name:  "EG_AGENT_IDENTIFIER_ID",
									Value: "ecozxmw5q4yh5typqkioh3ww6pafrf62",
								},
								{
									Name:  "DEBUG",
									Value: fmt.Sprintf("%t", true),
								},
								{
									Name:  "EG_CONTAINERZ_PORT_DISCOVERY",
									Value: fmt.Sprintf("%t", false),
								},
								{
									Name: "EG_HOST_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.hostIP",
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "var-run",
									MountPath: "/var/run",
								},
								{
									Name:      "process",
									MountPath: "/media/proc",
								},
								{
									Name:      "host-root",
									MountPath: "/mnt/",
								},
								{
									Name:      "btmmount",
									MountPath: "/opt/egbtm",
								},
							},
						},
					},
				},
			},
		},
	}
	// Set BtmAgent instance as the owner and controller
	ctrl.SetControllerReference(m, BtmAgentDaemonSet, r.Scheme)
	return BtmAgentDaemonSet
}

// SetupWithManager sets up the controller with the Manager.
func (r *BtmAgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&apmv1.BtmAgent{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&corev1.Namespace{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&mwc.MutatingWebhookConfiguration{}).Complete(r)
}
