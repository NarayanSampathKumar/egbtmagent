/*
Copyright 2021.

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

package v1

import (
	"context"
	"encoding/json"
	"net/http"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//+kubebuilder:webhook:path=/mutate-apps-v1-deployment,mutating=true,failurePolicy=fail,sideEffects=None,groups=apps,resources=deployments,verbs=create;update,versions=v1,name=mdeployment.kb.io,admissionReviewVersions=v1

type deploymentModifier struct {
	Client  client.Client
	decoder *admission.Decoder
}

var deploymentlog = logf.Log.WithName("deployment-Modifier-resource")

func NewDeploymentModifier(c client.Client) admission.Handler {
	return &deploymentModifier{Client: c}
}

func (v *deploymentModifier) Handle(ctx context.Context, req admission.Request) admission.Response {
	deploymentlog.Info("Received Mutation Request")
	deployment := &appsv1.Deployment{}
	var resourceInitContainers []apiv1.Container
	var resourceContainers []apiv1.Container
	var newContainers []apiv1.Container
	var resourceVolumes []apiv1.Volume

	err := v.decoder.Decode(req, deployment)
	if err != nil {
		deploymentlog.Info("Could not decode request!", "err", err)
		return admission.Errored(http.StatusBadRequest, err)
	}
	resourceVolumes = deployment.Spec.Template.Spec.Volumes
	resourceInitContainers = deployment.Spec.Template.Spec.InitContainers
	resourceVolumes = append(resourceVolumes, getegbtmVolume())
	resourceInitContainers = append(resourceInitContainers, getInitContainers())
	resourceContainers = deployment.Spec.Template.Spec.Containers
	//append Environment Variables and volumeounts to every container
	for _, container := range resourceContainers {
		container.Env = addEnv(container.Env)
		container.VolumeMounts = addVolumeMount(container.VolumeMounts)
		newContainers = append(newContainers, container)
	}
	deployment.Spec.Template.Spec.Containers = newContainers
	deployment.Spec.Template.Spec.InitContainers = resourceInitContainers
	deployment.Spec.Template.Spec.Volumes = resourceVolumes
	marshaledDeployment, err := json.Marshal(deployment)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledDeployment)
}
func addEnv(podenvvars []apiv1.EnvVar) []apiv1.EnvVar {
	var evar apiv1.EnvVar

	evar.Name = "JAVA_OPTS"
	evar.Value = "-javaagent:/opt/egbtm/eg_btm.jar -Dcom.sun.management.jmxremote -Dcom.sun.management.jmxremote.ssl=false -Dcom.sun.management.jmxremote.authenticate=false -Dcom.sun.management.jmxremote.port=3333 -Dcom.sun.management.jmxremote.rmi.port=3333"
	evar.ValueFrom = nil
	podenvvars = append(podenvvars, evar)

	evar.Name = "EG_COLLECTOR_IP"
	evar.ValueFrom = &apiv1.EnvVarSource{
		FieldRef: &apiv1.ObjectFieldSelector{
			APIVersion: "",
			FieldPath:  "status.hostIP",
		},
	}
	evar.Value = ""
	podenvvars = append(podenvvars, evar)
	return podenvvars
}

//add volume mount for egbt
func addVolumeMount(podvolmnts []apiv1.VolumeMount) []apiv1.VolumeMount {
	var volmnt apiv1.VolumeMount
	volmnt.Name = "egbtm"
	volmnt.MountPath = "/opt/egbtm"
	podvolmnts = append(podvolmnts, volmnt)
	return podvolmnts
}

//returns the init container that we want to attach
func getInitContainers() apiv1.Container {
	return apiv1.Container{
		Name:            "setup-egbtm",
		Image:           "eginnovations/java-profiler:7.2.0",
		ImagePullPolicy: apiv1.PullAlways,
		Command:         []string{"/bin/sh", "copy.sh"},
		Args:            []string{"/opt/egbtm"},
		VolumeMounts: []apiv1.VolumeMount{
			{Name: "egbtm", MountPath: "/opt/egbtm"},
		},
	}
}
func getegbtmVolume() apiv1.Volume {
	return apiv1.Volume{
		Name: "egbtm",
		VolumeSource: apiv1.VolumeSource{
			EmptyDir: &apiv1.EmptyDirVolumeSource{},
		},
	}
}
func (v *deploymentModifier) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
