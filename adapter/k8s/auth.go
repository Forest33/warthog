package k8s

import (
	"context"
	"encoding/base64"
	"fmt"

	"golang.org/x/oauth2"
	auth "golang.org/x/oauth2/google"
	"google.golang.org/api/container/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/exec"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/forest33/warthog/business/entity"
)

var gcsScopes = []string{"https://www.googleapis.com/auth/cloud-platform"}

func (c *Client) gcsAuth(ctx context.Context, r *entity.GCSAuth) (*rest.Config, error) {
	var token *oauth2.Token

	credentials, err := auth.FindDefaultCredentials(ctx, gcsScopes...)
	if err != nil {
		return nil, err
	}

	token, err = credentials.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	containerService, err := container.NewService(ctx)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", r.Project, r.Location, r.Cluster)
	resp, err := containerService.Projects.Locations.Clusters.Get(name).Do()
	if err != nil {
		return nil, err
	}

	cert, err := base64.StdEncoding.DecodeString(resp.MasterAuth.ClusterCaCertificate)
	if err != nil {
		return nil, err
	}

	apiConfig := &api.Config{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: map[string]*api.Cluster{
			r.Cluster: {
				CertificateAuthorityData: cert,
				Server:                   fmt.Sprintf("https://%s", resp.Endpoint),
			},
		},
		Contexts: map[string]*api.Context{
			r.Cluster: {
				Cluster:  r.Cluster,
				AuthInfo: r.Cluster,
			},
		},
		CurrentContext: r.Cluster,
	}

	restConfig, err := clientcmd.NewNonInteractiveClientConfig(
		*apiConfig,
		r.Cluster,
		&clientcmd.ConfigOverrides{
			CurrentContext: r.Cluster,
		},
		nil,
	).ClientConfig()
	if err != nil {
		return nil, err
	}

	restConfig.BearerToken = token.AccessToken

	return restConfig, nil
}
