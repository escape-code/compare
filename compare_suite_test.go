package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCompare(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Compare Suite")
}
