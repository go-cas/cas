package cas

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSessionStore(t *testing.T) {
	ss := NewMemorySessionStore()
	require.NotNil(t, ss)
}

func TestSessionStore_Get(t *testing.T) {
	ss := NewMemorySessionStore()
	require.NotNil(t, ss)

	v, ok := ss.Get("key1")
	require.False(t, ok)
	require.Equal(t, "", v)

	err := ss.Set("key1", "value1")
	require.Nil(t, err)

	v, ok = ss.Get("key1")
	require.True(t, ok)
	require.Equal(t, "value1", v)
}

func TestSessionStore_Set(t *testing.T) {
	ss := NewMemorySessionStore()
	require.NotNil(t, ss)

	err := ss.Set("key1", "value1")
	require.Nil(t, err)

	err = ss.Set("key2", "value2")
	require.Nil(t, err)

	v, ok := ss.Get("key1")
	require.True(t, ok)
	require.Equal(t, "value1", v)

	v, ok = ss.Get("key2")
	require.True(t, ok)
	require.Equal(t, "value2", v)

	err = ss.Set("key2", "value2-new")
	require.Nil(t, err)

	v, ok = ss.Get("key2")
	require.True(t, ok)
	require.Equal(t, "value2-new", v)
}

func TestSessionStore_Delete(t *testing.T) {
	ss := NewMemorySessionStore()
	require.NotNil(t, ss)

	err := ss.Set("key1", "value1")
	require.Nil(t, err)

	err = ss.Set("key2", "value2")
	require.Nil(t, err)

	v, ok := ss.Get("key1")
	require.True(t, ok)
	require.Equal(t, "value1", v)

	v, ok = ss.Get("key2")
	require.True(t, ok)
	require.Equal(t, "value2", v)

	err = ss.Delete("key2")
	require.Nil(t, err)

	v, ok = ss.Get("key1")
	require.True(t, ok)
	require.Equal(t, "value1", v)

	v, ok = ss.Get("key2")
	require.False(t, ok)
	require.Equal(t, "", v)
}
