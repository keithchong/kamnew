package scm

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/redhat-developer/kam/pkg/pipelines/triggers"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreatePushBindingForGitlab(t *testing.T) {
	repo, err := newGitLab("https://gitlab.com/org/fullname/subgroup/repository/subrepo/test")
	assertNoError(t, err)
	want := triggersv1.TriggerBinding{
		TypeMeta: triggers.TriggerBindingTypeMeta,
		ObjectMeta: v1.ObjectMeta{
			Name:      "gitlab-push-binding",
			Namespace: "testns",
		},
		Spec: triggersv1.TriggerBindingSpec{
			Params: []triggersv1.Param{
				{
					Name:  "gitrepositoryurl",
					Value: "$(body.project.git_http_url)",
				},
				{
					Name:  "fullname",
					Value: "$(body.project.path_with_namespace)",
				},
				{
					Name:  triggers.GitRef,
					Value: "$(extensions.ref)",
				},
				{
					Name:  triggers.GitCommitID,
					Value: "$(body.after)",
				},
				{
					Name:  triggers.GitCommitDate,
					Value: "$(body.commits[-1:].timestamp)",
				},

				{
					Name:  triggers.GitCommitMessage,
					Value: "$(body.commits[-1:].message)",
				},
				{
					Name:  triggers.GitCommitAuthor,
					Value: "$(body.commits[-1:].author.name)",
				},
			},
		},
	}
	got, name := repo.CreatePushBinding("testns")
	if name != "gitlab-push-binding" {
		t.Fatalf("CreatePushBinding() returned a wrong binding: want %v got %v", "gitlab-push-binding", name)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("CreatePushBinding() failed:\n%s", diff)
	}
}

func TestCreateCDTriggersForGitLab(t *testing.T) {
	repo, err := NewRepository("http://gitlab.com/org/test")
	assertNoError(t, err)
	rawSecret, err := secretParam("secret", "webhook-secret-key")
	assertNoError(t, err)
	rawFilter, rawOverlays, err := celParams(gitlabPushEventFilters, "org/test")
	assertNoError(t, err)
	name := "test-template"
	want := triggersv1.EventListenerTrigger{
		Name: "test",
		Bindings: []*triggersv1.EventListenerBinding{
			{Ref: "test-binding"},
		},
		Template: &triggersv1.EventListenerTemplate{Ref: &name},
		Interceptors: []*triggersv1.EventInterceptor{
			&triggersv1.TriggerInterceptor{
				Ref: triggersv1.InterceptorRef{
					Name: "gitlab",
				},
				Params: []triggersv1.InterceptorParams{
					{
						Name: "secretRef",
						Value: apiextensionsv1.JSON{
							Raw: rawSecret,
						},
					},
				},
			},
			&triggersv1.TriggerInterceptor{
				Ref: triggersv1.InterceptorRef{
					Name: "cel",
				},
				Params: []triggersv1.InterceptorParams{
					{
						Name: "filter",
						Value: apiextensionsv1.JSON{
							Raw: rawFilter,
						},
					},
					{
						Name: "overlays",
						Value: apiextensionsv1.JSON{
							Raw: rawOverlays,
						},
					},
				},
			},
		},
	}
	got, err := repo.CreatePushTrigger("test", "secret", "ns", "test-template", []string{"test-binding"})
	assertNoError(t, err)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("CreateCDTrigger() failed:\n%s", diff)
	}
}

func TestNewGitlabRepository(t *testing.T) {
	tests := []struct {
		url      string
		repoPath string
		errMsg   string
	}{
		{
			"http://gitlab.com",
			"",
			"invalid repository URL http://gitlab.com: path is empty",
		},
		{
			"http://gitlab.com/a",
			"",
			"invalid repository path for gitlab: /a",
		},
		{
			"http://gitlab.com/foo/bar",
			"foo/bar",
			"",
		},
		{
			"https://gitlab.com/group/subgroup/subgroup/repo.git",
			"group/subgroup/subgroup/repo",
			"",
		},
		{
			"https://gitlaB.com/foo/bar/test.git",
			"foo/bar/test",
			"",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(rt *testing.T) {
			repo, err := NewRepository(tt.url)
			if err != nil {
				if diff := cmp.Diff(tt.errMsg, err.Error()); diff != "" {
					rt.Fatalf("repo path errMsg mismatch: \n%s", diff)
				}
			}
			if repo != nil {
				if diff := cmp.Diff(tt.repoPath, repo.(*repository).path); diff != "" {
					rt.Fatalf("repo path mismatch: got\n%s", diff)
				}
			}
		})
	}
}
