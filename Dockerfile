# docker build -t esperta . --build-arg AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID --build-arg AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY && docker run -it esperta
FROM golang:1.10.1-stretch

# Terraform
ENV TERRAFORM_VERSION=0.11.7
ENV TERRAFORM_SHA256SUM=6b8ce67647a59b2a3f70199c304abca0ddec0e49fd060944c26f666298e23418
RUN apt-get update && apt-get install -y --no-install-recommends unzip vim && \
	curl https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip > terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
	echo "${TERRAFORM_SHA256SUM}  terraform_${TERRAFORM_VERSION}_linux_amd64.zip" > terraform_${TERRAFORM_VERSION}_SHA256SUMS && \
	sha256sum -c terraform_${TERRAFORM_VERSION}_SHA256SUMS && \
	unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /bin && \
	rm -f terraform_${TERRAFORM_VERSION}_linux_amd64.zip

# AWS CLI
ARG AWS_ACCESS_KEY_ID
ENV AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID

ARG AWS_SECRET_ACCESS_KEY
ENV AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY

RUN apt-get install -y \
	python \
	python-pip \
	&& \
	pip install awscli --upgrade --user && \
	rm -rf /var/lib/apt/lists/*

# API Server
RUN go get github.com/gorilla/mux
RUN mkdir -p	/usr/local/go/src/github.com/freddygv/SmartHouse-Server
COPY provision 	/usr/local/go/src/github.com/freddygv/SmartHouse-Server/provision
COPY Makefile 	/usr/local/go/src/github.com/freddygv/SmartHouse-Server/provision
COPY api 		/usr/local/go/src/github.com/freddygv/SmartHouse-Server/api
COPY go 		/usr/local/go/src/github.com/freddygv/SmartHouse-Server/go
COPY pub 		/usr/local/go/src/github.com/freddygv/SmartHouse-Server/pub
COPY main.go 	/usr/local/go/src/github.com/freddygv/SmartHouse-Server

WORKDIR /usr/local/go/src/github.com/freddygv/SmartHouse-Server/provision
RUN terraform init

CMD ["bash"]