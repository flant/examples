provider "aws" {
  profile = "default"
  region  = "us-east-1"
}

resource "aws_iam_user" "cloudwatch_exporter" {
  name = "cloudwatch_exporter"
  path = "/system/"
}

resource "aws_iam_access_key" "cloudwatch_exporter" {
  user = aws_iam_user.cloudwatch_exporter.name
}

resource "aws_iam_policy" "cloudwatch_policy" {
  name = "cloudwatch_policy"
  path = "/"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "cloudwatch:GetMetricStatistics",
        "cloudwatch:ListMetrics",
        "cloudwatch:PutMetricData",
        "ec2:DescribeVolumes",
        "ec2:DescribeTags",
        "logs:PutLogEvents",
        "logs:DescribeLogStreams",
        "logs:DescribeLogGroups",
        "logs:CreateLogStream",
        "logs:CreateLogGroup",
        "ce:GetCostAndUsage"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "ssm:GetParameter"
      ],
      "Resource": [
        "arn:aws:ssm:*:*:parameter/AmazonCloudWatch-*",
        "arn:aws:ce:*:*:/GetCostAndUsage"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_user_policy_attachment" "cloudwatch_exporter" {
  user       = aws_iam_user.cloudwatch_exporter.name
  policy_arn = aws_iam_policy.cloudwatch_policy.arn
}
