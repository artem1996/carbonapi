package stringutils

import "strconv"

// WriteFloat appends the string form of the floating-point number f,
// as generated by FormatFloat, to dst and returns the extended buffer.
func (sb *Builder) WriteFloat(f float64, fmt byte, prec, bitSize int) {
	s := strconv.FormatFloat(f, fmt, prec, bitSize)
	sb.WriteString(s)
}