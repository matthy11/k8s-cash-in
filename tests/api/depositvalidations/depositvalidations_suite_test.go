package depositvalidations_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDepositValidations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DepositValidations Suite")
}
