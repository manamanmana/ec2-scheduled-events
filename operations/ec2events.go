package operations

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2Events struct {
	client           *ec2.EC2
	InstanceStatuses []*ec2.InstanceStatus
	InstanceNames    []*string
}

func NewEC2Event(region string) *EC2Events {
	var e *EC2Events = new(EC2Events)
	e.client = ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
	return e
}

func (e *EC2Events) Events() error {
	// Request Input
	var input *ec2.DescribeInstanceStatusInput = &ec2.DescribeInstanceStatusInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("event.code"),
				Values: []*string{
					aws.String("instance-reboot"),
					aws.String("system-reboot"),
					aws.String("system-maintenance"),
					aws.String("instance-retirement"),
					aws.String("instance-stop"),
				},
			},
		},
		MaxResults: aws.Int64(1000),
	}

	// Request response and error
	var (
		resp *ec2.DescribeInstanceStatusOutput
		err  error
	)

	// Instance Status Request
	resp, err = e.client.DescribeInstanceStatus(input)
	if err != nil {
		return err
	}

	// Assign Instance Status Result List into Object
	e.InstanceStatuses = resp.InstanceStatuses
	// In the case Instance Status Result has NextToken
	for resp.NextToken != nil {
		input.NextToken = resp.NextToken
		resp, err = e.client.DescribeInstanceStatus(input)
		if err != nil {
			return err
		}
		e.InstanceStatuses = append(e.InstanceStatuses, resp.InstanceStatuses...)
	}

	// InstanceNames
	e.instanceNames()

	return nil
}

func (e *EC2Events) Results() []*string {
	var (
		results []*string
		i       int = 0
	)

	for _, is := range e.InstanceStatuses {
		for _, ev := range is.Events {
			var res string
			res = fmt.Sprintf("%s\t%s\t%s\t%s\t%s", *is.InstanceId, *e.InstanceNames[i], *ev.Code, *ev.NotBefore, *ev.NotAfter)
			results = append(results, &res)
			i++
		}
	}
	return results
}

func (e *EC2Events) instanceNameRequest(instanceid string, dryrun bool) (*ec2.DescribeInstancesOutput, error) {
	var input *ec2.DescribeInstancesInput = &ec2.DescribeInstancesInput{
		DryRun:      aws.Bool(dryrun),
		InstanceIds: []*string{aws.String(instanceid)},
	}

	return e.client.DescribeInstances(input)
}

func instanceName(output *ec2.DescribeInstancesOutput) *string {
	var tags []*ec2.Tag = output.Reservations[0].Instances[0].Tags
	var iname *string
	var tag *ec2.Tag
	for _, tag = range tags {
		if *tag.Key == "Name" {
			iname = tag.Value
			return iname
		}
	}
	return iname
}

func (e *EC2Events) instanceNames() {
	var is *ec2.InstanceStatus

	for _, is = range e.InstanceStatuses {
		for _, _ = range is.Events {
			var instanceid *string = is.InstanceId
			var iname *string
			resp, err := e.instanceNameRequest(*instanceid, false)
			if err != nil {
				*iname = ""
				continue
			}
			iname = instanceName(resp)
			e.InstanceNames = append(e.InstanceNames, iname)
		}
	}
}
