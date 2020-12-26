.PHONY: all
all: build/hlscale deploy

build/hlscale: tf-update
	GOOS=linux GOARCH=amd64 go build -o build/lambda/hlscale ./lambda/.

ci-role:
	cd terraform/ci && terraform init
	cd terraform/ci && terraform apply -auto-approve

.PHONY: clean
clean: go-clean tf-clean

deploy: build/hlscale
	cd terraform && terraform init
	cd terraform && terraform apply -auto-approve

.PHONY: destroy
destroy: destroy-app destroy-infra

destroy-app:
	cd terraform/app && terraform destroy

destroy-infra:
	aws s3 rm s3://hlscale-web --recursive
	cd terraform/infra && terraform destroy -auto-approve

go-clean:
	rm -rf build/*

infra:
	cd terraform/infra && terraform init
	cd terraform/infra && terraform apply -auto-approve
	aws s3 sync web/ s3://hlscale-web

test:
	go test ./..

tf-clean:
	rm -rf terraform/app/.terraform

tf-update:
	cd terraform/app && terraform get -update
	cd terraform/app && terraform init

tf-validate:
	cd terraform/app && terraform validate