package main

import (
	"context"
	"fmt"
	"log"

	link "github.com/kuzkuss/url_service/proto/link"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcConnLink, err := grpc.Dial(
		"0.0.0.0:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer grpcConnLink.Close()
	linkManager := link.NewLinksClient(grpcConnLink)

	pbOriginalLink := &link.OriginalLink {
		OriginalLink: "https://www.golang.org",
	}

	shortLink, err := linkManager.CreateShortLink(context.Background(), pbOriginalLink)
	if err != nil {
		log.Fatal(err)
	}

	pbShortLink := &link.ShortLink {
		ShortLink: shortLink.ShortLink,
	}

	originalLink, err := linkManager.GetOriginalLink(context.Background(), pbShortLink)
	if err != nil {
		log.Fatal(err)
	}

	if originalLink.OriginalLink != pbOriginalLink.OriginalLink {
		log.Fatal("incorrect work")
	}

	fmt.Println("SUCCESS")
}
