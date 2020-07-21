module github.com/vmware-tanzu/velero

go 1.14

require (
	cloud.google.com/go v0.46.2 // indirect
	github.com/Azure/azure-sdk-for-go v42.0.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.6
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2
	github.com/Azure/go-autorest/autorest/to v0.3.0
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/aws/aws-sdk-go v1.28.2
	github.com/docker/spdystream v0.0.0-20170912183627-bc6354cbbc29 // indirect
	github.com/ernesto-jimenez/gogen v0.0.0-20180125220232-d7d4131e6607 // indirect
	github.com/evanphx/json-patch v4.5.0+incompatible
	github.com/gobwas/glob v0.2.3
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/hashicorp/go-hclog v0.0.0-20180709165350-ff2cf002a8dd
	github.com/hashicorp/go-plugin v0.0.0-20190610192547-a1bc61569a26
	github.com/inconshreveable/go-vhost v0.0.0-20160627193104-06d84117953b // indirect
	github.com/joho/godotenv v1.3.0
	github.com/kubernetes-csi/external-snapshotter/v2 v2.2.0-rc1
	github.com/neelance/astrewrite v0.0.0-20160511093645-99348263ae86 // indirect
	github.com/neelance/sourcemap v0.0.0-20200213170602-2833bce08e4c // indirect
	github.com/onsi/ginkgo v1.13.0
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.8.1
	github.com/pkg/sftp v1.11.0 // indirect
	github.com/prometheus/client_golang v1.0.0
	github.com/robfig/cron v1.1.0
	github.com/rogpeppe/go-charset v0.0.0-20190617161244-0dc95cdf6f31 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	go.etcd.io/etcd v0.0.0-20191023171146-3cf2f69b5738
	golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7
	google.golang.org/genproto v0.0.0-20200731012542-8145dea6a485 // indirect
	google.golang.org/grpc v1.31.0
	google.golang.org/grpc/examples v0.0.0-20200731180010-8bec2f5d898f // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	k8s.io/api v0.18.4
	k8s.io/apiextensions-apiserver v0.18.4
	k8s.io/apimachinery v0.18.4
	k8s.io/cli-runtime v0.18.4
	k8s.io/client-go v0.18.4
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.6.1
	sigs.k8s.io/structured-merge-diff v1.0.2 // indirect
)

replace k8s.io/api => k8s.io/api v0.18.4

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.4

replace k8s.io/apimachinery => k8s.io/apimachinery v0.18.4

replace k8s.io/apiserver => k8s.io/apiserver v0.18.4

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.4

replace k8s.io/client-go => k8s.io/client-go v0.18.4

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.18.4

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.18.4

replace k8s.io/code-generator => k8s.io/code-generator v0.18.4

replace k8s.io/component-base => k8s.io/component-base v0.18.4

replace k8s.io/cri-api => k8s.io/cri-api v0.18.4

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.4

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.18.4

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.18.4

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.18.4

replace k8s.io/kubectl => k8s.io/kubectl v0.18.4

replace k8s.io/kubelet => k8s.io/kubelet v0.18.4

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.18.4

replace k8s.io/metrics => k8s.io/metrics v0.18.4

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.18.4

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.18.4
