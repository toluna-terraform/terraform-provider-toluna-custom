


provider "aws" {
  region  = "us-east-1"
  profile = "my-profile"
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
  payload = "{\"source_vpc\": \"bread_qa\",\"target_vpc\": \"butter_qa\"}"
}