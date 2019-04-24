package objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidAccountName(t *testing.T) {
	assert.True(t, IsValidAccountName("a"))
	assert.True(t, !IsValidAccountName("A"))
	assert.True(t, !IsValidAccountName("0"))
	assert.True(t, !IsValidAccountName("."))
	assert.True(t, !IsValidAccountName("-"))

	assert.True(t, IsValidAccountName("aa"))
	assert.True(t, !IsValidAccountName("aA"))
	assert.True(t, IsValidAccountName("a0"))
	assert.True(t, !IsValidAccountName("a."))
	assert.True(t, !IsValidAccountName("a-"))

	assert.True(t, IsValidAccountName("aaa"))
	assert.True(t, !IsValidAccountName("aAa"))
	assert.True(t, IsValidAccountName("a0a"))
	assert.True(t, IsValidAccountName("a.a"))
	assert.True(t, IsValidAccountName("a-a"))

	assert.True(t, IsValidAccountName("aa0"))
	assert.True(t, !IsValidAccountName("aA0"))
	assert.True(t, IsValidAccountName("a00"))
	assert.True(t, !IsValidAccountName("a.0"))
	assert.True(t, IsValidAccountName("a-0"))

	assert.True(t, IsValidAccountName("aaa-bbb-ccc"))
	assert.True(t, IsValidAccountName("aaa-bbb.ccc"))

	assert.True(t, !IsValidAccountName("aaa,bbb-ccc"))
	assert.True(t, !IsValidAccountName("aaa_bbb-ccc"))
	assert.True(t, !IsValidAccountName("aaa-BBB-ccc"))

	assert.True(t, !IsValidAccountName("1aaa-bbb"))
	assert.True(t, !IsValidAccountName("-aaa-bbb-ccc"))
	assert.True(t, !IsValidAccountName(".aaa-bbb-ccc"))
	assert.True(t, !IsValidAccountName("/aaa-bbb-ccc"))

	assert.True(t, !IsValidAccountName("aaa-bbb-ccc-"))
	assert.True(t, !IsValidAccountName("aaa-bbb-ccc."))
	assert.True(t, !IsValidAccountName("aaa-bbb-ccc.."))
	assert.True(t, !IsValidAccountName("aaa-bbb-ccc/"))

	assert.True(t, !IsValidAccountName("aaa..bbb-ccc"))
	assert.True(t, IsValidAccountName("aaa.bbb-ccc"))
	assert.True(t, IsValidAccountName("aaa.bbb.ccc"))

	assert.True(t, IsValidAccountName("aaa--bbb--ccc"))
	assert.True(t, IsValidAccountName("xn--sandmnnchen-p8a.de"))
	assert.True(t, IsValidAccountName("xn--sandmnnchen-p8a.dex"))
	assert.True(t, IsValidAccountName("xn-sandmnnchen-p8a.de"))
	assert.True(t, IsValidAccountName("xn-sandmnnchen-p8a.dex"))

	assert.True(t, IsValidAccountName("this-label-has-less-than-64-char.acters-63-to-be-really-precise"))
	assert.True(t, !IsValidAccountName("this-label-has-more-than-63-char.act.ers-64-to-be-really-precise"))
	assert.True(t, !IsValidAccountName("none.of.these.labels.has.more.than-63.chars--but.still.not.valid"))
}

func TestIsCheapAccountName(t *testing.T) {
	assert.True(t, IsCheapAccountName("mzt"))
	assert.True(t, IsCheapAccountName("maz1"))
	assert.True(t, !IsCheapAccountName("aay"))
}
