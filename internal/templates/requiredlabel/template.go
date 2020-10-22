package requiredlabel

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/check"
	"golang.stackrox.io/kube-linter/internal/diagnostic"
	"golang.stackrox.io/kube-linter/internal/extract"
	"golang.stackrox.io/kube-linter/internal/lintcontext"
	"golang.stackrox.io/kube-linter/internal/matcher"
	"golang.stackrox.io/kube-linter/internal/objectkinds"
	"golang.stackrox.io/kube-linter/internal/stringutils"
	"golang.stackrox.io/kube-linter/internal/templates"
	"golang.stackrox.io/kube-linter/internal/templates/requiredlabel/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Required Label",
		Key:         "required-label",
		Description: "Flag objects not carrying at least one label matching the provided patterns",
		SupportedObjectKinds: check.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.Any},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			keyMatcher, err := matcher.ForString(p.Key)
			if err != nil {
				return nil, errors.Wrap(err, "invalid key")
			}
			valueMatcher, err := matcher.ForString(p.Value)
			if err != nil {
				return nil, errors.Wrap(err, "invalid value")
			}

			return func(_ *lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				labels := extract.Labels(object.K8sObject)
				for k, v := range labels {
					if keyMatcher(k) && valueMatcher(v) {
						return nil
					}
				}
				return []diagnostic.Diagnostic{{
					Message: fmt.Sprintf("no label matching \"%s=%s\" found", p.Key, stringutils.OrDefault(p.Value, "<any>")),
				}}
			}, nil
		}),
	})
}
