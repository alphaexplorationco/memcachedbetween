package elasticache

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"strconv"
	"strings"
)

func ClusterNodes(_ *zap.Logger, endpoint string) ([]string, error) {
	if !strings.Contains(endpoint, ":") {
		endpoint = endpoint + ":11211"
	}
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	command := "config get cluster\r\n"
	_, err = fmt.Fprint(conn, command)
	if err != nil {
		return nil, err
	}

	response, err := parseNodes(conn)
	if err != nil {
		return nil, err
	}

	urls, err := parseURLs(response)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func parseNodes(conn io.Reader) (string, error) {
	var response string

	count := 0
	location := 3 // AWS docs suggest that nodes will always be listed on line 3

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		count++
		if count == location {
			response = scanner.Text()
		}
		if scanner.Text() == "END" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return response, nil
}

func parseURLs(response string) ([]string, error) {
	var urls []string

	items := strings.Split(response, " ")

	for _, v := range items {
		fields := strings.Split(v, "|")

		port, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, err
		}

		url := fmt.Sprintf("%s:%d", fields[0], port)
		urls = append(urls, url)
	}

	return urls, nil
}
