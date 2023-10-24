package otp

import (
	"github.com/coretrix/hitrix/service/component/sms"
	"math/rand"
	"strconv"
)

const SMSOTPProviderMobica = "Mobica"

type Mobica struct {
	mobicaProvider sms.IProvider
}

func NewMobicaSMSOTPProvider(mobicaProvider sms.IProvider) *Mobica {
	return &Mobica{
		mobicaProvider: mobicaProvider,
	}
}

func (m *Mobica) GetName() string {
	return SMSOTPProviderMobica
}

func (m *Mobica) GetCode() string {
	rangeMin := 10000
	rangeMax := 100000

	return strconv.Itoa(rand.Intn(rangeMax-rangeMin+1) + rangeMin)
}

func (m *Mobica) GetPhonePrefixes() []string {
	return nil
}

func (m *Mobica) SendOTP(phone *Phone, code string) (string, string, error) {
	_, err := m.mobicaProvider.SendSMSMessage(&sms.Message{
		Text:   code,
		Number: phone.Number,
	})

	return "", "", err
}

func (m *Mobica) Call(_ *Phone, _ string, _ string) (string, string, error) {
	// not implemented
	return "", "", nil
}

func (m *Mobica) VerifyOTP(_ *Phone, code, generatedCode string) (string, string, bool, bool, error) {
	return "", "", true, code == generatedCode, nil
}
