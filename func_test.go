package kit

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetDefault(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "", GetDefault[string]())
	})
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, 0, GetDefault[int]())
	})
	t.Run("float64", func(t *testing.T) {
		assert.Equal(t, 0., GetDefault[float64]())
	})
	t.Run("image.Point", func(t *testing.T) {
		assert.Equal(t, image.Point{}, GetDefault[image.Point]())
	})
	t.Run("*float64", func(t *testing.T) {
		assert.True(t, GetDefault[*float64]() == nil)
	})
}

func Test_First(t *testing.T) {
	t.Run("return value", func(t *testing.T) {
		slice := []*int{IntPtr(1), IntPtr(2), IntPtr(3)}
		assert.Equal(t, slice[1], First(slice, func(i *int) bool {
			return *i == 2
		}))
	})
	t.Run("nil slice", func(t *testing.T) {
		assert.True(t, First(nil, func(i *int) bool {
			return *i == 2
		}) == nil)
	})
	t.Run("empty slice", func(t *testing.T) {
		assert.True(t, First([]*int{}, func(i *int) bool {
			return *i == 2
		}) == nil)
	})
	t.Run("nil if not found", func(t *testing.T) {
		assert.True(t, First([]*int{IntPtr(1), IntPtr(1), IntPtr(3)}, func(i *int) bool {
			return *i == 2
		}) == nil)
	})
}
