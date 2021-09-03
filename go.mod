module github.com/dskatz/generators

go 1.16

// helm
replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
)

// kops
replace (
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.11.0
	k8s.io/api => k8s.io/api v0.19.5
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.5
	k8s.io/apiserver => k8s.io/apiserver v0.19.5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.5
	k8s.io/client-go => k8s.io/client-go v0.19.5
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.19.5
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.19.5
	k8s.io/code-generator => k8s.io/code-generator v0.19.5
	k8s.io/component-base => k8s.io/component-base v0.19.5
	k8s.io/cri-api => k8s.io/cri-api v0.19.5
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.19.5
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.19.5
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.19.5
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.19.5
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.19.5
	k8s.io/kubectl => k8s.io/kubectl v0.19.5
	k8s.io/kubelet => k8s.io/kubelet v0.19.5
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.19.5
	k8s.io/metrics => k8s.io/metrics v0.19.5
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.19.5
)

require (
	github.com/10gen/kanopy v0.1.8
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.42.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	helm.sh/helm/v3 v3.5.3
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/gengo v0.0.0-20200710205751-c0d492a0f3ca
	k8s.io/klog/v2 v2.4.0
	k8s.io/kube-aggregator v0.19.0
	k8s.io/utils v0.0.0-20210305010621-2afb4311ab10
	sigs.k8s.io/yaml v1.2.0
)
