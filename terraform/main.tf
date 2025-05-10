

# ===================================
# Network
# ===================================
resource "aws_vpc" "go-lamdba-tutorial-vpc" {
  cidr_block       = "10.0.0.0/16"
  instance_tenancy = "default"
  enable_dns_support = true
  enable_dns_hostnames = true

  tags = {
    Name = "go-lamdba-tutorial-vpc"
    type = "go-lamdba-tutorial"
  }
}

resource "aws_internet_gateway" "go-lamdba-tutorial-gw" {
  vpc_id = aws_vpc.go-lamdba-tutorial-vpc.id
  tags = {
    Name = "go-lamdba-tutorial-gw"
    type = "go-lamdba-tutorial"
  }
}

resource "aws_subnet" "go-lamdba-tutorial-public-subnet" {
  vpc_id     = aws_vpc.go-lamdba-tutorial-vpc.id
  cidr_block = "10.0.1.0/24"
  tags = {
    Name = "go-lamdba-tutorial-public-subnet"
    type = "go-lamdba-tutorial"
  }
}

resource "aws_subnet" "go-lamdba-tutorial-private-subnet" {
  vpc_id     = aws_vpc.go-lamdba-tutorial-vpc.id
  cidr_block = "10.0.2.0/24"
  tags = {
    Name = "go-lamdba-tutorial-private-subnet"
    type = "go-lamdba-tutorial"
  }
}

resource "aws_eip" "go-lamdba-tutorial-eip" {
  domain  = "vpc"
    tags = {
    Name = "go-lamdba-tutorial-eip"
    type = "go-lamdba-tutorial"
  }
}

resource "aws_nat_gateway" "go-lamdba-tutorial-ngw" {
  subnet_id = aws_subnet.go-lamdba-tutorial-public-subnet.id
  allocation_id = aws_eip.go-lamdba-tutorial-eip.id
  tags = {
    Name = "go-lamdba-tutorial-ngw"
    type = "go-lamdba-tutorial"
  }
}

resource "aws_route_table" "go-lamdba-tutorial-public-route" {
  vpc_id = aws_vpc.go-lamdba-tutorial-vpc.id
  route {
    cidr_block = "10.0.0.0/16"
    gateway_id = "local"
  }
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.go-lamdba-tutorial-gw.id
  }
  tags = {
    Name = "go-lamdba-tutorial-public-route"
    type = "go-lamdba-tutorial"
  }
}

resource "aws_route_table" "go-lamdba-tutorial-private-route" {
  vpc_id = aws_vpc.go-lamdba-tutorial-vpc.id
  route {
    cidr_block = "10.0.0.0/16"
    gateway_id = "local"
  }
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_nat_gateway.go-lamdba-tutorial-ngw.id
  }
  tags = {
    Name = "go-lamdba-tutorial-private-route"
    type = "go-lamdba-tutorial"
  }
}


# ===================================
# Security
# ===================================
data "aws_iam_policy_document" "go-lamdba-tutorial-assume-role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "go-lamdba-tutorial-lamdba-role" {
  name               = "go-lamdba-tutorial-lamdba-role"
  assume_role_policy = data.aws_iam_policy_document.go-lamdba-tutorial-assume-role.json
}

# Lambda関数へのAPI Gateway呼び出し許可
# refer to: https://docs.aws.amazon.com/lambda/latest/dg/lambda-permissions.html
resource "aws_lambda_permission" "api_gw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.go-lamdba-tutorial-lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.go-lamdba-tutorial-api-gw.execution_arn}/*"
}


# ===================================
# lambda
# ===================================
resource "aws_lambda_function" "go-lamdba-tutorial-lambda" {
  function_name = "go-lamdba-tutorial"
  role          = aws_iam_role.go-lamdba-tutorial-lamdba-role.arn
  architectures = ["arm64"]
  package_type = "Image"
  image_uri     = "428230071946.dkr.ecr.ap-northeast-1.amazonaws.com/yukin520/go-lamdba-tutorial:latest"
  image_config {
    command           = ["/main"]
    working_directory = "/"
  }
  tags = {
    Name = "go-lamdba-tutorial"
    type = "go-lamdba-tutorial"
  }
}


# ===================================
# api gateway
# ===================================
resource "aws_apigatewayv2_api" "go-lamdba-tutorial-api-gw" {
  name          =  "go-lamdba-tutorial-api-gw"
  protocol_type = "HTTP"

  tags = {
    Name = "go-lamdba-tutorial-api-gw"
    type = "go-lamdba-tutorial"
  }
}

resource "aws_apigatewayv2_integration" "go-lamdba-tutorial-api-integration" {
  api_id           = aws_apigatewayv2_api.go-lamdba-tutorial-api-gw.id
  integration_type = "AWS_PROXY"

  connection_type           = "INTERNET"
  description               = "integrate 'go-lamdba-tutorial-api-stage' lambda with api-gateway"
  integration_method        = "POST"
  integration_uri           = aws_lambda_function.go-lamdba-tutorial-lambda.invoke_arn
  passthrough_behavior      = "WHEN_NO_MATCH"
  payload_format_version    = "2.0"
}

resource "aws_apigatewayv2_route" "go-lamdba-tutorial-api-route" {
  api_id    = aws_apigatewayv2_api.go-lamdba-tutorial-api-gw.id
  route_key = "ANY /"

  target = "integrations/${aws_apigatewayv2_integration.go-lamdba-tutorial-api-integration.id}"
}

resource "aws_apigatewayv2_deployment" "go-lamdba-tutorial-deploy" {
  api_id      = aws_apigatewayv2_api.go-lamdba-tutorial-api-gw.id
  description = "deploy api gateway for go-lamdba-tutorial"

  lifecycle {
    create_before_destroy = true
  }
  depends_on = [
    aws_apigatewayv2_route.go-lamdba-tutorial-api-route
  ]
}

resource "aws_apigatewayv2_stage" "go-lamdba-tutorial-api-stage" {
  api_id = aws_apigatewayv2_api.go-lamdba-tutorial-api-gw.id
  name   = "go-lamdba-tutorial-api-stage"

  auto_deploy = true
}



