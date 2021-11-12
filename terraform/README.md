# terraform/

In this directory is Terraform code that provisions the bucket to which `journal-scraper` will write
data, and which will be used as the web endpoint. Steps to use it:

1. Get AWS credentials in place so that the [AWS
   provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs) can see them
2. Run `terraform init`
3. Run `terraform apply`
