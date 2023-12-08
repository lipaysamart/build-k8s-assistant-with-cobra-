package cmd

import (
	"fmt"
	"github.com/lipaysamart/build-k8s-assistant-with-cobra/internal/client"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

const (
	EnvToken = "TOKEN"
	EnvLang  = "LANG"
)

func getEnvOrDefault(key, defVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defVal
}

type Options struct {
	Token string
	Lang  string

	kubeConfigFlags *genericclioptions.ConfigFlags
}

func NewOptions() Options {
	return Options{
		kubeConfigFlags: genericclioptions.NewConfigFlags(true),
	}
}

func (o *Options) Complete() error {
	o.Token = os.Getenv(EnvToken)
	if len(o.Token) == 0 {
		return fmt.Errorf("请检查 Api Token with ENV %s", EnvToken)
	}
	o.Lang = getEnvOrDefault(EnvLang, "Chinese")
	return nil
}

func (o *Options) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(o.kubeConfigFlags.KubeConfig, "kubeconfig", *o.kubeConfigFlags.KubeConfig, "Path to the kubeconfig file to use for CLI requests.")
	flags.StringVarP(o.kubeConfigFlags.Namespace, "namespace", "n", *o.kubeConfigFlags.Namespace, "If present, the namespace scope for this CLI requests.")
}

func (o *Options) NewBuilder() *resource.Builder {
	return resource.NewBuilder(o.kubeConfigFlags)
}

func (o *Options) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return o.kubeConfigFlags.ToRawKubeConfigLoader()
}

func (o *Options) NewChatGPTClient(spinnerSuffix string) client.Client {
	return client.NewChatGPTClient(o.Token, spinnerSuffix)

}

func (o *Options) NewKubeClientSet() (kubernetes.Interface, error) {
	cfg, err := o.kubeConfigFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("create rest config: %w", err)
	}
	return kubernetes.NewForConfig(cfg)
}
