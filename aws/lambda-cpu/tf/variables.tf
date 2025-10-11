variable "project_name" {
  type    = string
  default = "lambda-cpu"
}

variable "dynamodb_table_name" {
  type    = string
  default = "lambda-cpu-results"
}

variable "architectures" {
  type    = list(string)
  default = ["x86_64", "arm64"]
}

variable "lambda_zip_x86" {
  type    = string
  default = "../lambdasrc/dist/function-x86_64.zip"
}

variable "lambda_zip_arm" {
  type    = string
  default = "../lambdasrc/dist/function-arm64.zip"
}
