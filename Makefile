.PHONY: all
all: build/hlscale deploy

build/hlscale: tf-update-app
	GOOS=linux GOARCH=amd64 go build -o build/lambda/hlscale ./lambda/.

ci-role:
	cd terraform/ci && terraform init
	cd terraform/ci && terraform apply -auto-approve

.PHONY: clean
clean: go-clean tf-clean-app

deploy: build/hlscale
	cd terraform/app && terraform init
	cd terraform/app && terraform apply -auto-approve

.PHONY: destroy
destroy: destroy-app destroy-infra

destroy-app:
	cd terraform/app && terraform destroy -auto-approve

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

tf-clean-%:
	rm -rf terraform/$*/.terraform

tf-update-%:
	cd terraform/$* && terraform get -update
	cd terraform/$* && terraform init

tf-validate-%:
	cd terraform/$* && terraform validate