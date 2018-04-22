variable "aws_ssh_key_file" {}

variable "aws_region" {
  default     = "eu-west-1"
  description = "AWS Region"
}

variable "ami" {
  default = "ami-94b236ed"
}

variable "instance_type" {
  default = "t2.micro"
}
