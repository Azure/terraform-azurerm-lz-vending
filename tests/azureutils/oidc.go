package azureutils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type oidcCredential struct {
	assertion, oidcTokenRequestUrl, oidcTokenRequestToken string
}

func (o *oidcCredential) getAssertion(c context.Context) (string, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, o.oidcTokenRequestUrl, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("oidc: failed to create new request")
	}

	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return "", fmt.Errorf("oidc: failed to parse query string")
	}

	query.Set("audience", "api://AzureADTokenExchange")
	req.URL.RawQuery = query.Encode()

	req.Header = http.Header{
		"Accept":        {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", o.oidcTokenRequestToken)},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("oidc: failed to request token: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("oidc: failed to parse response: %v", err)
	}

	if c := resp.StatusCode; c < 200 || c > 299 {
		return "", fmt.Errorf("oidc: failure... received HTTP status %d with response: %s", resp.StatusCode, body)
	}

	var tokenRes struct {
		Value *string `json:"value"`
	}
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		return "", fmt.Errorf("OidcCredential: failed to unmarshal response: %v", err)
	}
	o.assertion = *tokenRes.Value

	return o.assertion, nil
}
