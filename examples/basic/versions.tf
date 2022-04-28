terraform {
  required_providers {
    toluna-custom = {
      #source = "toluna.com/tolunaprovider/toluna"
      #version = "1.0.0"
      source = "toluna-terraform/toluna-custom"
      version = ">=0.0.9"
    }
  }
}

provider "toluna-custom" {
  # Configuration options
}
