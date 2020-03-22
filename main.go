package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	domain := os.Getenv("DNS_HOSTED_ZONE_DOMAIN")

	mySession := session.Must(session.NewSession())
	svc := route53.New(mySession)

	res, err := svc.ListHostedZonesByName(&route53.ListHostedZonesByNameInput{DNSName: &domain})
	if err != nil {
		log.Fatal(err)
	}

	var hostedZone *route53.HostedZone

	for _, v := range res.HostedZones {
		if *v.Name == domain {
			hostedZone = v
		}
	}

	if hostedZone == nil {
		log.Fatalf("Hosted zone %v not found", domain)
	}

	var browsed int
	action := "DELETE"
	svc.ListResourceRecordSetsPages(
		&route53.ListResourceRecordSetsInput{
			HostedZoneId: hostedZone.Id,
		},
		func(list *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
			browsed += len(list.ResourceRecordSets)
			var changes []*route53.Change
			for _, v := range list.ResourceRecordSets {
				if *v.Type == "CNAME" {
					changes = append(changes, &route53.Change{
						Action:            &action,
						ResourceRecordSet: v,
					})
				}
			}

			_, err := svc.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
				HostedZoneId: hostedZone.Id,
				ChangeBatch:  &route53.ChangeBatch{Changes: changes},
			})

			if err != nil {
				log.Fatal(err)
			}

			res, err := svc.GetHostedZone(&route53.GetHostedZoneInput{Id: hostedZone.Id})
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Browsed %v records, deleted %v CNAME records, remains %v records", browsed, len(changes), *res.HostedZone.ResourceRecordSetCount)

			return !lastPage
		})
}
