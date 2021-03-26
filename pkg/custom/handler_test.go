// +build integration

package custom

import (
	"context"
	"errors"
	"testing"

	xrv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/pivotal-cf/brokerapi/v7/middlewares"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vshn/crossplane-service-broker/pkg/crossplane"
	"github.com/vshn/crossplane-service-broker/pkg/integration"
	"github.com/vshn/crossplane-service-broker/pkg/reqcontext"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestAPIHandler_Endpoints(t *testing.T) {
	type args struct {
		ctx        context.Context
		instanceID string
	}
	ctx := context.WithValue(context.TODO(), middlewares.CorrelationIDKey, "corrid")

	tests := []struct {
		name      string
		args      args
		want      []Endpoint
		wantErr   error
		resources func() (func(c client.Client) error, []client.Object)
	}{
		{
			name: "requires instance to be ready before getting endpoints",
			args: args{
				ctx:        ctx,
				instanceID: "1-1-1",
			},
			resources: func() (func(c client.Client) error, []client.Object) {
				servicePlan := integration.NewTestServicePlan("1", "1-1", crossplane.RedisService)
				objs := []client.Object{
					integration.NewTestService("1", crossplane.RedisService),
					servicePlan.Composition,
					integration.NewTestInstance("1-1-1", servicePlan, crossplane.RedisService, "", ""),
				}
				return nil, objs
			},
			want:    nil,
			wantErr: errors.New(`unable to get secret: secrets "1-1-1" not found`),
		},
		{
			name: "creates a redis instance and gets the endpoints",
			args: args{
				ctx:        ctx,
				instanceID: "1-1-1",
			},
			resources: func() (func(c client.Client) error, []client.Object) {
				servicePlan := integration.NewTestServicePlan("1", "1-1", crossplane.RedisService)
				instance := integration.NewTestInstance("1-1-1", servicePlan, crossplane.RedisService, "", "")
				objs := []client.Object{
					integration.NewTestService("1", crossplane.RedisService),
					integration.NewTestServicePlan("1", "1-2", crossplane.RedisService).Composition,
					servicePlan.Composition,
					instance,
					integration.NewTestSecret(integration.TestNamespace, "1-1-1", map[string]string{
						xrv1.ResourceCredentialsSecretPortKey:     "1234",
						xrv1.ResourceCredentialsSecretEndpointKey: "localhost",
						xrv1.ResourceCredentialsSecretPasswordKey: "supersecret",
						"sentinelPort": "21234",
					}),
				}
				return func(c client.Client) error {
					return integration.UpdateInstanceConditions(ctx, c, servicePlan, instance, xrv1.TypeReady, corev1.ConditionTrue, xrv1.ReasonAvailable)
				}, objs
			},
			want: []Endpoint{
				{
					Destination: "localhost",
					Ports:       "1234",
					Protocol:    "tcp",
				},
				{
					Destination: "localhost",
					Ports:       "21234",
					Protocol:    "tcp",
				},
			},
			wantErr: nil,
		},
	}

	m, logger, cp, err := integration.SetupManager(t)
	require.NoError(t, err, "unable to setup integration test manager")
	defer m.Cleanup()

	handler := NewAPIHandler(cp, logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, objs := tt.resources()
			require.NoError(t, integration.CreateObjects(tt.args.ctx, objs)(m.GetClient()))
			defer func() {
				require.NoError(t, integration.RemoveObjects(tt.args.ctx, objs)(m.GetClient()))
			}()
			if fn != nil {
				require.NoError(t, fn(m.GetClient()))
			}

			got, err := handler.Endpoints(reqcontext.NewReqContext(ctx, logger, nil), tt.args.instanceID)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
