package gotrue

import (
	"net/url"
)

func getParametersFromURI(rawURI string) (url.Values, error) {
	uri, err := url.Parse(rawURI)
	if err != nil {
		return nil, err
	}

	values := uri.Query()
	if values == nil {
		values = make(url.Values)
	}

	fragValues, _ := url.ParseQuery(uri.Fragment)
	if fragValues != nil {
		for k, v := range fragValues {
			values[k] = append(values[k], v...)
		}
	}
	return values, nil
}
