package registry

import (
	"github.com/mpppk/sutaba-server/pkg/application/output"
	"github.com/mpppk/sutaba-server/pkg/interface/converter"
)

type Converter interface {
	NewMessageConverter() output.MessageConverter
}

type ConverterConfig struct {
}

type converterImpl struct {
	messageConverter output.MessageConverter
}

func (p *converterImpl) NewMessageConverter() output.MessageConverter {
	return p.messageConverter
}

func NewConverter(config *ConverterConfig) Converter {
	messageConverter := &converter.SutabaPoliceMessageConverter{}
	return &converterImpl{
		messageConverter: messageConverter,
	}
}
