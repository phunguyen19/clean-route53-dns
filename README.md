# Clean Route53 DNS CLI

## About

Motivation: The AWS Rout53 Web UI is limited on listing and filtering data, so this tool is mainly built to use on development, help developer cleaning the DNS records on Route53.

> Note: This tools is under developing base on author's real work demand. Any new feature that you need to use can be created as an issue or make a pull request. Feel free to do it.

## Usage

Set up environment variable by copying and modifying `.env.example` or you can follow the AWS documentation to do it.

Run the command `go run main.go`

The programing will running and clean all `CNAME` records for you.
