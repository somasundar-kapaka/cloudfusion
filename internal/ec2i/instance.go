package ec2i

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	computealphav1 "github.com/somasundar-kapaka/cloudfusion/api/alphav1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func CreateEC2Instance(ctx context.Context, o *computealphav1.EC2Instance) error {

	log := logf.FromContext(ctx)

	awsEC2Client := NewClient(ctx, o.Spec.Region)

	res, err := awsEC2Client.RunInstances(ctx, CreateEC2InstanceReqType(o))
	if err != nil {
		log.Error(err, fmt.Sprintf("error creating aws instance %s", o.Name), "intance-type", o.Spec.InstanceType)
		return err
	}

	err = pollEC2Instance(ctx, awsEC2Client, res)
	if err != nil {
		return err
	}

	return nil
}

func pollEC2Instance(ctx context.Context, awsEC2Client *ec2.Client, res *ec2.RunInstancesOutput) error {

	awsWaiter := ec2.NewInstanceRunningWaiter(awsEC2Client)
	maxWaitTime := 5 * time.Minute

	insId := res.Instances[0].InstanceId

	des := &ec2.DescribeInstancesInput{
		InstanceIds: []string{*insId},
	}
	err := awsWaiter.Wait(ctx, des, maxWaitTime)
	if err != nil {
		return err
	}

	return nil
}

func CreateEC2InstanceReqType(o *computealphav1.EC2Instance) *ec2.RunInstancesInput {

	return &ec2.RunInstancesInput{

		InstanceType: types.InstanceType(o.Spec.InstanceType),
	}
}

func ValidteNewInstanceRequest(o *computealphav1.EC2Instance) error {

	var v ec2.RunInstancesInput
	instanceValues := v.InstanceType.Values()

	switch {

	// Add more cases below as the spec feild validations
	case o.Spec.Region == "":
		return InvalidRegion

	case o.Spec.InstanceType == "", !slices.Contains(instanceValues, types.InstanceType(o.Spec.InstanceType)):
		return fmt.Errorf(InvalidInstanceType.Error(), o.Spec.InstanceType)

	}

	return nil
}



func CheckInstanceExists() {}


