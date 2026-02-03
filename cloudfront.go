package encodingx

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/aura-studio/reflectx"
)

var (
	ErrCloudFrontURLSafeWrongValueType = errors.New("encoding cloudfront url safe converts on wrong type value")
)

// CloudFrontURLSafe implements CloudFront's URL-safe base64 encoding.
// It uses standard base64 encoding but replaces:
// - '+' with '-'
// - '=' with '_'
// - removes padding '='
type CloudFrontURLSafe struct{}

func init() {
	register(NewCloudFrontURLSafe())
}

func NewCloudFrontURLSafe() *CloudFrontURLSafe {
	return new(CloudFrontURLSafe)
}

func (c CloudFrontURLSafe) String() string {
	return reflectx.TypeName(c)
}

func (CloudFrontURLSafe) Style() EncodingStyleType {
	return EncodingStyleBytes
}

func (CloudFrontURLSafe) Marshal(v interface{}) ([]byte, error) {
	var data []byte
	switch v := v.(type) {
	case []byte:
		data = v
	case Bytes:
		data = v.Data
	case *Bytes:
		data = v.Data
	default:
		return nil, ErrCloudFrontURLSafeWrongValueType
	}

	// Standard base64 encode
	s := base64.StdEncoding.EncodeToString(data)
	// CloudFront URL-safe replacements
	s = strings.ReplaceAll(s, "+", "-")
	s = strings.ReplaceAll(s, "=", "_")
	s = strings.ReplaceAll(s, "/", "~")
	return []byte(s), nil
}

func (CloudFrontURLSafe) Unmarshal(data []byte, v interface{}) error {
	switch v := v.(type) {
	case *Bytes:
		s := string(data)
		// Reverse CloudFront URL-safe replacements
		s = strings.ReplaceAll(s, "-", "+")
		s = strings.ReplaceAll(s, "_", "=")
		s = strings.ReplaceAll(s, "~", "/")
		decoded, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return err
		}
		v.Data = decoded
		return nil
	default:
		return ErrCloudFrontURLSafeWrongValueType
	}
}

func (c CloudFrontURLSafe) Reverse() Encoding {
	return c
}
