
resource "aws_iam_role" "oteleport" {
  name = "oteleport-lambda"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}


resource "aws_iam_policy" "oteleport" {
  name   = "oteleport"
  path   = "/"
  policy = data.aws_iam_policy_document.oteleport.json
}

resource "aws_cloudwatch_log_group" "oteleport" {
  name              = "/aws/lambda/oteleport"
  retention_in_days = 7
}

resource "aws_iam_role_policy_attachment" "oteleport" {
  role       = aws_iam_role.oteleport.name
  policy_arn = aws_iam_policy.oteleport.arn
}

data "aws_iam_policy_document" "oteleport" {
  statement {
    actions = [
      "ssm:GetParameter*",
    ]
    resources = ["*"]
  }
  statement {
    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:ListBucket",
      "s3:ListO"
    ]
    resources = ["*"]
  }
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = ["*"]
  }
}

data "archive_file" "oteleport_dummy" {
  type        = "zip"
  output_path = "${path.module}/oteleport_dummy.zip"
  source {
    content  = timestamp()
    filename = "bootstrap"
  }
  depends_on = [
    terraform_data.oteleport_dummy,
  ]
}

resource "terraform_data" "oteleport_dummy" {}

resource "terraform_data" "oteleport_dummy_cleanup" {
  triggers_replace = [
    data.archive_file.oteleport_dummy.output_md5
  ]

  provisioner "local-exec" {
    command = "rm ${data.archive_file.oteleport_dummy.output_path}"
  }
  depends_on = [
    aws_lambda_function.oteleport,
  ]
}

resource "aws_lambda_function" "oteleport" {
  lifecycle {
    ignore_changes = all
  }

  function_name = "oteleport"
  role          = aws_iam_role.oteleport.arn
  architectures = ["arm64"]
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  filename      = data.archive_file.oteleport_dummy.output_path
}

resource "aws_lambda_alias" "oteleport" {
  lifecycle {
    ignore_changes = all
  }
  name             = "current"
  function_name    = aws_lambda_function.oteleport.arn
  function_version = aws_lambda_function.oteleport.version
}


resource "aws_s3_bucket" "oteleport" {
  bucket = "oteleport-test"
}

data "aws_caller_identity" "current" {}
