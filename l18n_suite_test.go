package l18n_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestL18n(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "L18n Suite")
}
