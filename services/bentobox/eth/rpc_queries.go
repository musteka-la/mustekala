package eth

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// getNetworkHeight will send an eth_blockNumber request and
// parse its results
func (e *EthManager) getNetworkHeight() (response int64, err error) {
	body := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":42}`
	target := ethBlockNumber{}
	if err = requestAndParseJSON(e.ethJsonRPC, body, &target); err != nil {
		log.Printf("Query Error: %v", err)
		return -1, err
	}

	response, err = strconv.ParseInt(target.Result[2:], 16, 0)
	if err != nil {
		log.Printf("Parse Error: %v", err)
		return -1, err
	}

	// won't log non error responses as this is a very frequent query
	return
}

// requestAndParseJSON is a helper to send RPC Queries
func requestAndParseJSON(url, body string, target interface{}) error {
	client := &http.Client{
		Timeout: RPC_TIMEOUT,
	}
	request, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return err
	}
	defer request.Body.Close()
	request.Header.Add("Content-Type", "application/json")

	// won't log request as this is a very frequent query
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	return json.NewDecoder(response.Body).Decode(target)
}
