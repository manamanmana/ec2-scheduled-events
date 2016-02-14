# ec2-scheduled-events
ec2-scheduled-events is a small CLI tool by which you can get EC2 instance events (instance-reboot, system-reboot, system-maintenance, instance-retirement and instance-stop) by specifing AWS region information.

## Installing

```
$ go get -u github.com/manamanmana/ec2-scheduled-events
```

## Configuring
As this CLI is using `aws-sdk-go`, you need to specify some credentials to access with AWS API.
You can set them as `~/.aws/credentials` or environment variables like `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`.

You can see this in detail on [https://github.com/aws/aws-sdk-go](https://github.com/aws/aws-sdk-go) .

## Usage

```
$ ec2-scheduled-events --region=<AWS region code>

Example:

$ ec2-scheduled-events --region=ap-northeast-1
```

## Output

```
CSV format with delimiter "\t".

<Instance ID>\t<Instance Name == Tag:Name>\t<Event Code>\t<Event Start Time>\t<Event End Time>\n
<Instance ID>\t<Instance Name == Tag:Name>\t<Event Code>\t<Event Start Time>\t<Event End Time>\n
...

Event Code == [instance-reboot|system-reboot|system-maintenance|instance-retirement|instance-stop]
```

