package telegram

import "github.com/gotd/td/bin"

func DecodeTgDto[T bin.Decoder](payload []byte, decodedPayload T) error {
	encodedPayload := &bin.Buffer{Buf: payload}
	err := decodedPayload.Decode(encodedPayload)
	if err != nil {
		return err
	}

	return nil
}

func EncodeTgDto[T bin.Encoder](dto T) ([]byte, error) {
	encodedDialogs := bin.Buffer{Buf: []byte{}}
	err := dto.Encode(&encodedDialogs)
	if err != nil {
		return nil, err
	}
	return encodedDialogs.Raw(), nil
}
