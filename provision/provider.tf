provider "aws" {
  region = "${var.aws_region}"
}

resource "aws_key_pair" "auth" {
  key_name   = "ce_key"
  public_key = "${file("${var.aws_ssh_key_file}.pub")}"
}

resource "aws_security_group" "instance_group" {
  name        = "casa-esperta-group"
  description = "SSH access and TCP access on 8888"

  # TCP access for Android client requests
  ingress {
    from_port   = 8888
    to_port     = 8888
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "api_gateway" {
  ami                         = "${var.ami}"
  instance_type               = "${var.instance_type}"
  key_name                    = "${aws_key_pair.auth.key_name}"
  security_groups             = ["${aws_security_group.instance_group.name}"]
  associate_public_ip_address = true

  tags {
    Name = "casa-esperta-gateway"
  }
}

output "public_ip" {
  description = "Public IP address assigned to the instance"
  value       = "${aws_instance.api_gateway.public_ip}"
}

output "public_dns" {
  description = "public DNS name assigned to the instance"
  value       = "${aws_instance.api_gateway.public_dns}"
}
