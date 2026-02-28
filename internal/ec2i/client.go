package ec2i // EC2Instance -> ec2i

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func NewClient(ctx context.Context, region string) *ec2.Client {

	log := logf.FromContext(ctx)

	log.Info("Creating AWS client")

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Error(err, "error creating AWS client")
		return nil
	}

	return ec2.NewFromConfig(cfg)

}
