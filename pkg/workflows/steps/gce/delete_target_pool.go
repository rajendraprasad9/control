package gce

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/workflows/steps"
	"google.golang.org/api/compute/v1"
)

const DeleteTargetPoolStepName = "gce_delete_target_pool"

type DeleteTargetPoolStep struct {
	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewDeleteTargetPoolStep() (*DeleteTargetPoolStep, error) {
	return &DeleteTargetPoolStep{
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				deleteTargetPool: func(ctx context.Context, config steps.GCEConfig, targetPoolName string) (*compute.Operation, error) {
					config.AvailabilityZone = "us-central1-a"
					return client.InstanceGroups.Delete(config.ServiceAccount.ProjectID, config.AvailabilityZone, targetPoolName).Do()
				},
			}, nil
		},
	}, nil
}

func (s *DeleteTargetPoolStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {

	logrus.Debugf("Step %s", DeleteTargetPoolStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", DeleteTargetPoolStepName)
	}

	_, err = svc.deleteInstanceGroup(ctx, config.GCEConfig, config.GCEConfig.InstanceGroupName)

	if err != nil {
		logrus.Errorf("Error deleting target pool %v", err)
		return errors.Wrapf(err, "%s deleting target pool caused", DeleteTargetPoolStepName)
	}

	return nil
}

func (s *DeleteTargetPoolStep) Name() string {
	return DeleteTargetPoolStepName
}

func (s *DeleteTargetPoolStep) Depends() []string {
	return nil
}

func (s *DeleteTargetPoolStep) Description() string {
	return "Delete target pool master nodes"
}

func (s *DeleteTargetPoolStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
