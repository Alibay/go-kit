//go:build integration

package google

import (
	"testing"

	"github.com/Alibay/go-kit"
	"github.com/stretchr/testify/suite"
)

type captchaTestSuite struct {
	kit.Suite
}

func (s *captchaTestSuite) SetupSuite() {
	s.Suite.Init(nil)
}

func TestCaptchaSuite(t *testing.T) {
	suite.Run(t, new(captchaTestSuite))
}

const (
	dummyV2Captcha = "6LeIxAcTAAAAAJcZVRqyHh71UMIEGNQ_MXjiZKhI"
	dummyV2Key     = "6LeIxAcTAAAAAGG-vFI1TnRWxMZNFuojJ4WifJWe"
)

func (s *captchaTestSuite) Test_WhenPassedWithDummyKey() {
	cpt := NewCaptcha(&Config{ReCaptchaSecretV2: dummyV2Key}, s.L())
	r, err := cpt.Verify(s.Ctx, dummyV2Captcha, "0.0.0.0", "v2")
	s.NoError(err)
	s.True(r)
}

func (s *captchaTestSuite) Test_WhenInvalid() {

	test := func(key, cap, ver string) {
		cpt := NewCaptcha(&Config{ReCaptchaSecretV2: key}, s.L())
		r, _ := cpt.Verify(s.Ctx, cap, "0.0.0.0", ver)
		s.False(r)
	}

	test("invalid", dummyV2Captcha, "v2")
	test(dummyV2Key, dummyV2Captcha, "invalid")
	test("", "", "")
}
