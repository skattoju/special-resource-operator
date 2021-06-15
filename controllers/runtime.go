package controllers

import (
	"context"
	"strings"
	"time"

	srov1beta1 "github.com/openshift-psap/special-resource-operator/api/v1beta1"
	"github.com/openshift-psap/special-resource-operator/pkg/cache"
	"github.com/openshift-psap/special-resource-operator/pkg/clients"
	"github.com/openshift-psap/special-resource-operator/pkg/cluster"
	"github.com/openshift-psap/special-resource-operator/pkg/exit"
	"github.com/openshift-psap/special-resource-operator/pkg/kernel"
	"github.com/openshift-psap/special-resource-operator/pkg/proxy"
	"github.com/openshift-psap/special-resource-operator/pkg/upgrade"

	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
	//machineV1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
)

type ResourceGroupName struct {
	DriverBuild            string `json:"driverBuild"`
	DriverContainer        string `json:"driverContainer"`
	RuntimeEnablement      string `json:"runtimeEnablement"`
	DevicePlugin           string `json:"devicePlugin"`
	DeviceMonitoring       string `json:"deviceMonitoring"`
	DeviceDashboard        string `json:"deviceDashboard"`
	DeviceFeatureDiscovery string `json:"deviceFeatureDiscovery"`
	CSIDriver              string `json:"csiDriver"`
}

type RuntimeInformation struct {
	OnOCP                     bool                           `json:"onOCP"`
	OperatingSystemMajor      string                         `json:"operatingSystemMajor"`
	OperatingSystemMajorMinor string                         `json:"operatingSystemMajorMinor"`
	OperatingSystemDecimal    string                         `json:"operatingSystemDecimal"`
	KernelFullVersion         string                         `json:"kernelFullVersion"`
	KernelPatchVersion        string                         `json:"kernelPatchVersion"`
	DriverToolkitImage        string                         `json:"driverToolkitImage"`
	ClusterVersion            string                         `json:"clusterVersion"`
	ClusterVersionMajorMinor  string                         `json:"clusterVersionMajorMinor"`
	ClusterUpgradeInfo        map[string]upgrade.NodeVersion `json:"clusterUpgradeInfo"`
	UpdateVendor              string                         `json:"updateVendor"`
	PushSecretName            string                         `json:"pushSecretName"`
	OSImageURL                string                         `json:"osImageURL"`
	Proxy                     proxy.Configuration            `json:"proxy"`
	GroupName                 ResourceGroupName              `json:"groupName"`
	SpecialResource           srov1beta1.SpecialResource     `json:"specialresource"`
}

var RunInfo = RuntimeInformation{
	OnOCP:                     true,
	OperatingSystemMajor:      "",
	OperatingSystemMajorMinor: "",
	OperatingSystemDecimal:    "",
	KernelFullVersion:         "",
	KernelPatchVersion:        "",
	ClusterVersion:            "",
	ClusterVersionMajorMinor:  "",
	ClusterUpgradeInfo:        make(map[string]upgrade.NodeVersion),
	UpdateVendor:              "",
	PushSecretName:            "",
	OSImageURL:                "",
	Proxy:                     proxy.Configuration{},
	GroupName: ResourceGroupName{
		DriverBuild:            "driver-build",
		DriverContainer:        "driver-container",
		RuntimeEnablement:      "runtime-enablement",
		DevicePlugin:           "device-plugin",
		DeviceMonitoring:       "device-monitoring",
		DeviceDashboard:        "device-dashboard",
		DeviceFeatureDiscovery: "device-feature-discovery",
		CSIDriver:              "csi-driver"},
	SpecialResource: srov1beta1.SpecialResource{},
}

func logRuntimeInformation() {
	log.Info("Runtime Information", "OnOCP", RunInfo.OnOCP)
	log.Info("Runtime Information", "OperatingSystemMajor", RunInfo.OperatingSystemMajor)
	log.Info("Runtime Information", "OperatingSystemMajorMinor", RunInfo.OperatingSystemMajorMinor)
	log.Info("Runtime Information", "OperatingSystemDecimal", RunInfo.OperatingSystemDecimal)
	log.Info("Runtime Information", "KernelFullVersion", RunInfo.KernelFullVersion)
	log.Info("Runtime Information", "KernelPatchVersion", RunInfo.KernelPatchVersion)
	log.Info("Runtime Information", "ClusterVersion", RunInfo.ClusterVersion)
	log.Info("Runtime Information", "ClusterVersionMajorMinor", RunInfo.ClusterVersionMajorMinor)
	log.Info("Runtime Information", "ClusterUpgradeInfo", RunInfo.ClusterUpgradeInfo)
	log.Info("Runtime Information", "UpdateVendor", RunInfo.UpdateVendor)
	log.Info("Runtime Information", "PushSecretName", RunInfo.PushSecretName)
	log.Info("Runtime Information", "OSImageURL", RunInfo.OSImageURL)
	log.Info("Runtime Information", "Proxy", RunInfo.Proxy)
}

func getRuntimeInformation(r *SpecialResourceReconciler) {

	var err error

	err = cache.Nodes(r.specialresource.Spec.NodeSelector, false)
	exit.OnError(errors.Wrap(err, "Failed to cache nodes"))

	RunInfo.OnOCP = cluster.OnOCP()

	RunInfo.OperatingSystemMajor, RunInfo.OperatingSystemMajorMinor, RunInfo.OperatingSystemDecimal, err = cluster.OperatingSystem()
	exit.OnError(errors.Wrap(err, "Failed to get operating system"))

	RunInfo.KernelFullVersion, err = kernel.FullVersion()
	exit.OnError(errors.Wrap(err, "Failed to get kernel version"))

	RunInfo.KernelPatchVersion, err = kernel.PatchVersion(RunInfo.KernelFullVersion)
	exit.OnError(errors.Wrap(err, "Failed to get kernel patch version"))

	RunInfo.ClusterVersion, RunInfo.ClusterVersionMajorMinor, err = cluster.Version()
	cluster.WarnOnK8sFailOnOCP(err, "Failed to get cluster version")

	RunInfo.ClusterUpgradeInfo, err = upgrade.ClusterInfo()
	exit.OnError(errors.Wrap(err, "Failed to get upgrade info"))

	RunInfo.PushSecretName, err = retryGetPushSecretName(r)
	cluster.WarnOnK8sFailOnOCP(err, "Failed to get push secret name")

	RunInfo.OSImageURL, err = cluster.OSImageURL()
	cluster.WarnOnK8sFailOnOCP(err, "Failed to get OSImageURL")

	RunInfo.Proxy, err = proxy.ClusterConfiguration()
	cluster.WarnOnK8sFailOnOCP(err, "Failed to get Proxy Configuration")

	r.specialresource.DeepCopyInto(&RunInfo.SpecialResource)
}

func retryGetPushSecretName(r *SpecialResourceReconciler) (string, error) {
	for i := 0; i < 3; i++ {
		time.Sleep(2 * time.Second)
		pushSecretName, err := getPushSecretName(r)
		if err != nil {
			log.Info("Cannot find Secret builder-dockercfg " + r.specialresource.Spec.Namespace)
			continue
		} else {
			return pushSecretName, err
		}
	}

	return "", errors.New("Cannot find Secret builder-dockercfg")

}

func getPushSecretName(r *SpecialResourceReconciler) (string, error) {

	secrets := &unstructured.UnstructuredList{}

	secrets.SetAPIVersion("v1")
	secrets.SetKind("SecretList")

	log.Info("Getting SecretList in Namespace: " + r.specialresource.Spec.Namespace)
	opts := []client.ListOption{
		client.InNamespace(r.specialresource.Spec.Namespace),
	}
	err := clients.Interface.List(context.TODO(), secrets, opts...)
	if err != nil {
		return "", errors.Wrap(err, "Client cannot get SecretList")
	}

	log.Info("Searching for builder-dockercfg Secret")
	for _, secret := range secrets.Items {
		secretName := secret.GetName()

		if strings.Contains(secretName, "builder-dockercfg") {
			log.Info("Found", "Secret", secretName)
			return secretName, nil
		}
	}

	return "", errors.New("Cannot find Secret builder-dockercfg")
}
