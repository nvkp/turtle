// Package scanner implements parsing Turtle data provided as a byte slice
// and reading it triple by triple. It handles the compact version of Turtle
// just as the N-triples version where each row corresponds to a single triple.
// It handles @base and @forms. It ignores comments and labels and data types
// assigned to object literals.
package scanner
