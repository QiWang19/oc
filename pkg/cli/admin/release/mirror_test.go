package release

import (
	"reflect"
	"testing"

	apicfgv1 "github.com/openshift/api/config/v1"

	"k8s.io/apimachinery/pkg/util/diff"
)

func Test_dedupeSortSources(t *testing.T) {
	tests := []struct {
		name            string
		sources         []apicfgv1.ImageDigestMirrors
		expectedSources []apicfgv1.ImageDigestMirrors
	}{
		{
			name: "single Source single Mirror",
			sources: []apicfgv1.ImageDigestMirrors{
				{
					Source:  "quay/ocp/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/ocp/test"},
				},
			},
			expectedSources: []apicfgv1.ImageDigestMirrors{
				{
					Source:  "quay/ocp/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/ocp/test"},
				},
			},
		},
		{
			name: "single Source multiple Mirrors",
			sources: []apicfgv1.ImageDigestMirrors{
				{
					Source:  "quay/ocp/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/ocp/test"},
				},
				{
					Source:  "quay/ocp/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/another/test"},
				},
				{
					Source:  "quay/ocp/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/ocp/test"},
				},
			},
			expectedSources: []apicfgv1.ImageDigestMirrors{
				{
					Source:  "quay/ocp/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/ocp/test", "registry/another/test"},
				},
			},
		},
		{
			name: "multiple Source single Mirrors",
			sources: []apicfgv1.ImageDigestMirrors{
				{
					Source:  "quay/ocp/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/ocp/test"},
				},
				{
					Source:  "quay/another/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/ocp/test"},
				},
			},
			expectedSources: []apicfgv1.ImageDigestMirrors{
				{
					Source:  "quay/another/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/ocp/test"},
				},
				{
					Source:  "quay/ocp/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/ocp/test"},
				},
			},
		},
		{
			name: "multiple Source multiple Mirrors",
			sources: []apicfgv1.ImageDigestMirrors{
				{
					Source:  "quay/ocp/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/ocp/test"},
				},
				{
					Source:  "quay/ocp/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/another/test"},
				},
				{
					Source:  "quay/another/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/ocp/test"},
				},
				{
					Source:  "quay/another/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/another/test"},
				},
			},
			expectedSources: []apicfgv1.ImageDigestMirrors{
				{
					Source:  "quay/another/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/ocp/test", "registry/another/test"},
				},
				{
					Source:  "quay/ocp/test",
					Mirrors: []apicfgv1.ImageMirror{"registry/ocp/test", "registry/another/test"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uniqueSources := dedupeSortSources(tt.sources)
			if !reflect.DeepEqual(uniqueSources, tt.expectedSources) {
				t.Errorf("%s", diff.ObjectReflectDiff(uniqueSources, tt.expectedSources))
			}
		})
	}
}
