package humane_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHumane(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Humane Suite")
}
