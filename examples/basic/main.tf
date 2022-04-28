
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
  payload = "{\"example\": \"my_example\"}"
}

resource "toluna-custom_start_codebuild" "example" {
  region = "us-east-1"
  aws_profile = "my-profile"
  project_name = "my_project"
  payload = "{\"my_property\": \"example\"}"
}