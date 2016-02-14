package operations

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"strings"
	"testing"
	"time"
)

func TestEvents(t *testing.T) {
	var e2e *EC2Events = NewEC2Event("ap-northeast-1")
	err := e2e.Events()
	if err != nil {
		t.Errorf("Error on Events AWS request: %v", err.Error())
	}
}

func TestResults(t *testing.T) {
	var e2e *EC2Events = NewEC2Event("ap-northeast-1")

	// Set up aws query result data structs
	e2e.InstanceStatuses = []*ec2.InstanceStatus{
		{
			InstanceId: aws.String("i-0ce73d0e"),
			Events: []*ec2.InstanceStatusEvent{
				{
					Code:        aws.String(ec2.EventCodeInstanceReboot),
					Description: aws.String("Description 1"),
					NotAfter:    aws.Time(time.Date(2016, 2, 14, 0, 0, 0, 0, time.UTC)),
					NotBefore:   aws.Time(time.Date(2016, 2, 13, 0, 0, 0, 0, time.UTC)),
				},
				{
					Code:        aws.String(ec2.EventCodeInstanceRetirement),
					Description: aws.String("Description 2"),
					NotAfter:    aws.Time(time.Date(2016, 2, 12, 0, 0, 0, 0, time.UTC)),
					NotBefore:   aws.Time(time.Date(2016, 2, 11, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		{
			InstanceId: aws.String("i-308d5729"),
			Events: []*ec2.InstanceStatusEvent{
				{
					Code:        aws.String(ec2.EventCodeInstanceStop),
					Description: aws.String("Description 3"),
					NotAfter:    aws.Time(time.Date(2016, 2, 10, 0, 0, 0, 0, time.UTC)),
					NotBefore:   aws.Time(time.Date(2016, 2, 9, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
	}

	e2e.InstanceNames = []*string{
		aws.String("opc-ap1"),
		aws.String("opc-ap1"),
		aws.String("opc-cb1"),
	}

	var results []*string = e2e.Results()

	if len(results) != 3 {
		t.Errorf("Result Cound should be 3. but %d", len(results))
	}

	var expline1 = "i-0ce73d0e\topc-ap1\tinstance-reboot\t2016-02-13 00:00:00 +0000 UTC\t2016-02-14 00:00:00 +0000 UTC"
	var expline2 = "i-0ce73d0e\topc-ap1\tinstance-retirement\t2016-02-11 00:00:00 +0000 UTC\t2016-02-12 00:00:00 +0000 UTC"
	var expline3 = "i-308d5729\topc-cb1\tinstance-stop\t2016-02-09 00:00:00 +0000 UTC\t2016-02-10 00:00:00 +0000 UTC"
	var expects []*string = []*string{&expline1, &expline2, &expline3}

	for i, res := range results {
		if *res != *expects[i] {
			t.Errorf("Line %d is not same as expected.", i)
			t.Errorf("Expected: %s", *expects[i])
			t.Errorf("Actual: %s", *res)
		}
	}
}

func TestInstanceNameRequest(t *testing.T) {
	var e2e *EC2Events = NewEC2Event("ap-northeast-1")

	var err error

	_, err = e2e.instanceNameRequest("i-xxxxxxxx", true)
	if strings.Index(err.Error(), "DryRunOperation") == -1 {
		t.Error("Error on DescribeInstances AWS Request.")
	}

}

func TestInstanceName(t *testing.T) {
	var output *ec2.DescribeInstancesOutput = &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						Tags: []*ec2.Tag{
							{
								Key:   aws.String("Name"),
								Value: aws.String("my-instance"),
							},
						},
					},
				},
			},
		},
	}

	var expect *string = aws.String("my-instance")
	var actual *string = instanceName(output)

	if *expect != *actual {
		t.Errorf("Faled to extract Instance Name. Expected %s, Actual %s", *expect, *actual)
	}
}

func TestInstancNames(t *testing.T) {
	var e2e *EC2Events = NewEC2Event("ap-northeast-1")

	e2e.InstanceStatuses = []*ec2.InstanceStatus{
		{
			InstanceId: aws.String("i-0ce73d0e"),
		},
		{
			InstanceId: aws.String("i-xxxxxxxx"),
		},
		{
			InstanceId: aws.String("i-308d5729"),
		},
	}

	e2e.instanceNames()

	var (
		expect1 string    = "opc-ap1"
		expect2 string    = ""
		expect3 string    = "opc-cb1"
		expects []*string = []*string{&expect1, &expect2, &expect3}
		actual  []*string = e2e.InstanceNames
	)

	for i, a := range actual {
		if *a != *expects[i] {
			t.Errorf("Failed to get Instance Names: expect: %s, actual: %s", *expects[i], *a)
		}
	}
}
