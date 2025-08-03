package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Read input from standard input
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()

		// Call the ReverseString function
		output := ReverseString(input)

		// Print the result
		fmt.Println(output)
	}
}

// ReverseString returns the reversed string of s.
func ReverseString(s string) string {
	// TODO: Implement the function
	r := []rune(s)
	i, j := 0, len(r)-1
	for i < j {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// ReverseString returns the reversed string of s.
// func ReverseString(s string) string {
// 	TODO: Implement the function
// 	r := []rune(s)
// 	r1, r2 := r[:len(s)/2], r[len(s)/2:]
// 	l1, l2 := len(r1), len(r2)
// 	for i := 0; i < l1 && i < l2; i++ {
// 		r1[i], r2[l2-i-1] = r2[l2-i-1], r[i]
// 	}
// 	return string(r)
// }
