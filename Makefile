default: build plan

deps:
	go install github.com/hashicorp/terraform
	go get github.com/denverdino/aliyungo

build:
	go build -o terraform-provider-aliyun .

test:
	go test -v ./aliyun

plan:
	@terraform plan
