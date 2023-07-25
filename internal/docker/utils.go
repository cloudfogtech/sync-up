package docker

import (
	"errors"

	"github.com/cloudfogtech/sync-up/internal/utils"
)

func ReadResultChan(resultChan *ResultChan) (string, error) {
	result := ""
	errs := make([]error, 0)
	f := false
	for !f {
		select {
		case r := <-resultChan.Data:
			result += utils.GetPrintable(string(r))
		case err := <-resultChan.Err:
			errs = append(errs, err)
		case <-resultChan.Exit:
			f = true
		}
	}
	return result, errors.Join(errs...)
}
