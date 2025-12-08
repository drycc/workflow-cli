package parser

import (
	"fmt"
	"strings"

	"github.com/drycc/controller-sdk-go/api"
)

func parseHeaders(headers []string) ([]*api.KVPair, error) {
	var parsedHeaders []*api.KVPair
	for _, header := range headers {
		parsedHeader, err := parseHeader(header)
		if err != nil {
			return nil, err
		}
		parsedHeaders = append(parsedHeaders, parsedHeader)
	}
	return parsedHeaders, nil
}

func parseHeader(header string) (*api.KVPair, error) {
	headerParts := strings.SplitN(header, ":", 2)
	if len(headerParts) != 2 {
		return nil, fmt.Errorf("could not find separator in header (%s)", header)
	}
	return &api.KVPair{
		Name:  strings.TrimSpace(headerParts[0]),
		Value: strings.TrimSpace(headerParts[1]),
	}, nil
}
