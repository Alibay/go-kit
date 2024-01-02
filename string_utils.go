package kit

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"io"
	"regexp"
	"strconv"

	gonanoid "github.com/matoous/go-nanoid/v2"
	uuid "github.com/satori/go.uuid"
)

var (
	encoding     = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769")
	digitsRegExp = regexp.MustCompile(`^\d+$`)
)

const (
	baseAlphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	numericAlphabet = "0123456789"
)

// NanoId generates random string
func NanoId() string {
	r, _ := gonanoid.Generate(baseAlphabet, 22)
	return r
}

// NumCode generates random number code with the given size
func NumCode(size int) string {
	r, _ := gonanoid.Generate(numericAlphabet, size)
	return r
}

// NewRandString generates a unique string
func NewRandString() string {
	var b bytes.Buffer
	encoder := base32.NewEncoder(encoding, &b)
	_, _ = encoder.Write(uuid.NewV4().Bytes())
	_ = encoder.Close()
	b.Truncate(26)
	return b.String()
}

// NewId generates UUID
func NewId() string {
	return uuid.NewV4().String()
}

// UUID generates UUID
func UUID(size int) string {
	u := make([]byte, size)
	_, _ = io.ReadFull(rand.Reader, u)
	return hex.EncodeToString(u)
}

// ValidateUUIDs check UUID format nad return error if at least one UUID isn't in a correct format
func ValidateUUIDs(uuids ...string) error {
	for _, u := range uuids {
		if _, err := uuid.FromString(u); err != nil {
			return err
		}
	}
	return nil
}

// Nil returns nil UUID
func Nil() string {
	return uuid.Nil.String()
}

// Strings represents slice of strings
type Strings []string

// Distinct returns distinct slice
func (s Strings) Distinct() Strings {
	var res []string
	m := make(map[string]struct{}, len(s))
	for _, i := range s {
		if _, ok := m[i]; !ok {
			m[i] = struct{}{}
			res = append(res, i)
		}
	}
	return res
}

// Contains check if a strings is in slice
func (s Strings) Contains(str string) bool {
	for _, i := range s {
		if i == str {
			return true
		}
	}
	return false
}

func (s Strings) Intersect(r Strings) Strings {
	res := Strings{}
	rDistinct := r.Distinct()
	for _, i := range s.Distinct() {
		for _, j := range rDistinct {
			if i == j {
				res = append(res, i)
			}
		}
	}
	return res
}

func (s Strings) Equal(r Strings) bool {
	if len(s) != len(r) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] != r[i] {
			return false
		}
	}
	return true
}

func (s Strings) ToMap() map[string]struct{} {
	r := make(map[string]struct{})
	if len(s) == 0 {
		return r
	}
	for _, i := range s {
		r[i] = struct{}{}
	}
	return r
}

func StrToInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func RemoveNonAlfaDigital(str string) string {
	reg := regexp.MustCompile(`[^0-9a-zA-ZА-Яа-я]|\^|\_`)
	return reg.ReplaceAllString(str, "")
}

func Digits(s string) bool {
	if s == "" {
		return false
	}
	return digitsRegExp.MatchString(s)
}
