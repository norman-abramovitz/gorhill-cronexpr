/*!
 * Copyright 2026 Norman Abramovitz
 *
 * Project: github.com/gorhill/cronexpr
 * File: cronexpr_hash.go
 * License: pick the one which suits you :
 *   GPL v3 see <https://www.gnu.org/licenses/gpl.html>
 *   APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>
 *
 */

package cronexpr

/******************************************************************************/

import (
	"hash/fnv"
)

// H picks values from a hash of a key, so that jobs given the same
// expression spread out over the range: H (one value of the field), H/n (every
// n-th value, from an offset the hash picks), H(a-b) and H(a-b)/n (the same
// within a to b). The choices follow netresearch/go-cron, so a key picks the
// same values in both.

/******************************************************************************/

// hashOf returns the FNV-1a hash of key.
func hashOf(key string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(key))
	return h.Sum64()
}

// rangeSize returns how many values lie from first to last, wrapping past the
// field's maximum when first > last.
func rangeSize(first, last int, desc fieldDescriptor) int {
	if first > last {
		return (desc.max - first + 1) + (last - desc.min + 1)
	}
	return last - first + 1
}

// wrapValue brings a value counted on past the field's maximum back into the
// field.
func wrapValue(v int, desc fieldDescriptor) int {
	if v > desc.max {
		return (v-desc.min)%(desc.max-desc.min+1) + desc.min
	}
	return v
}

// dowCycle maps a day-of-week range onto the seven distinct days: 7 is
// Sunday (0), and 0-7 is the whole week. Other fields are unchanged.
func dowCycle(first, last int, desc fieldDescriptor) (int, int, fieldDescriptor) {
	if desc.name != dowDescriptor.name {
		return first, last, desc
	}
	if first == 0 && last == 7 {
		last = 6
	}
	desc.max = 6
	return first % 7, last % 7, desc
}

// hashOne returns the one value H(first-last) picks.
func hashOne(first, last int, desc fieldDescriptor, hash uint64) int {
	size := rangeSize(first, last, desc)
	v := wrapValue(first+int(hash%uint64(size)), desc)
	if v == 7 && desc.name == dowDescriptor.name {
		v = hashSunday7(first, last, size, hash)
	}
	return v
}

// hashSunday7 resolves a day of week H that picked 7. When the range also
// holds 0, 7 is a second Sunday: pick again among the range's distinct days
// with the rest of the hash, so each day is picked as often. When 7 is the
// range's only Sunday, it is Sunday.
func hashSunday7(first, last, size int, hash uint64) int {
	if first <= last && first != 0 {
		return 0
	}
	// Position of 7 in the walk first, first+1, ... (mod 8); skip it.
	pos7 := (7 - first) % 8
	j := int(hash / uint64(size) % uint64(size-1))
	if j >= pos7 {
		j++
	}
	return (first + j) % 8
}

// hashStep returns the values H(first-last)/step picks: every step-th value,
// from an offset the hash picks. In the day of week it walks the seven
// distinct days.
func hashStep(first, last, step int, desc fieldDescriptor, hash uint64) []int {
	first, last, desc = dowCycle(first, last, desc)
	size := rangeSize(first, last, desc)
	offset := int(hash % uint64(step))
	if offset >= size {
		// A step above the range size: take the value from the rest of the
		// hash, which keeps the choice even
		offset = int(hash / uint64(step) % uint64(size))
	}
	end := last
	if first > last {
		end = last + desc.max - desc.min + 1
	}
	var values []int
	for v := first + offset; v <= end; v += step {
		values = append(values, wrapValue(v, desc))
	}
	return values
}
