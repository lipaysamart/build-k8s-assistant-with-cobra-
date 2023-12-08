package audit

import (
	"bytes"
	"fmt"
	"github.com/lipaysamart/build-k8s-assistant-with-cobra/cmd"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/meta"
)

func New(opt cmd.Options) *cobra.Command {
	c := &cobra.Command{
		Use:          "audit TYPE NAME",
		Short:        "Audit a resource",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ns, _, _ := opt.ToRawKubeConfigLoader().Namespace()
			obj, err := opt.NewBuilder().
				NamespaceParam(ns).
				Unstructured().
				ResourceNames(args[0], args[1]).
				Do().
				Object()
			if err != nil {
				return fmt.Errorf("get object: %w", err)
			}
			metaObj, err := meta.Accessor(obj)
			if err == nil {
				metaObj.SetManagedFields(nil)
			}
			data, err := yaml.Marshal(obj)
			if err != nil {
				return fmt.Errorf("marshal object: %w", err)
			}
			var buf bytes.Buffer
			_ = promptAudit.Execute(&buf, templateData{
				Data: string(data),
				Lang: opt.Lang,
			})

			return opt.NewChatGPTClient("Auditing...").
				CreateCompletion(cmd.Context(), buf.String(), cmd.OutOrStdout())
		},
	}
	opt.AddFlags(c.Flags())
	return c
}
