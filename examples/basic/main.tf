
terraform {
  required_providers {
    toluna-custom = {
      source = "toluna-terraform/toluna-custom"
      version = ">=0.0.9"
    }
  }
}

provider "aws" {
  region  = "us-east-1"
  profile = "my-profile"
}

provider "toluna-custom" {
  
}


resource "toluna-custom_invoke_lambda" "example" {
  region = "us-east-1"
  aws_profile = "my-profile"
  function_name = "my_lambda"
  payload = jsonencode({"name": "example pay load"})
}

resource "toluna-custom_start_codebuild" "example" {
  region = "us-east-1"
  aws_profile = "my-profile"
  project_name = "my_project"
    environment_variables  {
    name = "my-variable"
    value = "FOO"
    type = "PLAINTEXT"
  }
  environment_variables  {
    name = "my-secret-variable"
    value = "BAR"
    type = "PARAMETER_STORE"
  }
  environment_variables  {
    name = "my-other-secret-variable"
    value = "BAR"
    type = "SECRETS_MANAGER"
  }
}