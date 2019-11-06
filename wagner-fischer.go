package smetrics

// BytePair structs hold a pair of bytes
type BytePair struct {
	firstByte  byte
	secondByte byte
}

//DefaultSubstitutionWeights holds weights for common misidentifications of
// letters and numbers taken from https://www.ismp.org/resources/misidentification-alphanumeric-symbols
var DefaultSubstitutionWeights = map[BytePair]int{
	BytePair{'l', '1'}: 1,
	BytePair{'b', '6'}: 1,
	BytePair{'o', '0'}: 1,
	BytePair{'g', '9'}: 1,
	BytePair{'q', '9'}: 1,
	BytePair{'G', '6'}: 1,
	BytePair{'F', '7'}: 1,
	BytePair{'Z', '2'}: 1,
	BytePair{'Z', '7'}: 1,
	BytePair{'Q', '2'}: 1,
	BytePair{'O', '0'}: 1,
	BytePair{'B', '8'}: 1,
	BytePair{'D', '0'}: 1,
	BytePair{'S', '5'}: 1,
	BytePair{'S', '8'}: 1,
	BytePair{'Y', '5'}: 1,
	BytePair{'T', '7'}: 1,
	BytePair{'U', '0'}: 1,
	BytePair{'U', '4'}: 1,
	BytePair{'A', '4'}: 1,
	BytePair{'}', '1'}: 1,
	BytePair{'{', '1'}: 1,

	BytePair{'0', '8'}: 3,
	BytePair{'3', '9'}: 3,
	BytePair{'3', '8'}: 3,
	BytePair{'4', '9'}: 3,
	BytePair{'5', '8'}: 3,
	BytePair{'3', '5'}: 3,
	BytePair{'6', '8'}: 3,
	BytePair{'0', '9'}: 3,
	BytePair{'7', '1'}: 3,

	BytePair{'g', 'q'}: 1,
	BytePair{'p', 'n'}: 1,
	BytePair{'m', 'n'}: 1,
	BytePair{'y', 'z'}: 1,
	BytePair{'u', 'v'}: 1,
	BytePair{'c', 'e'}: 1,
	BytePair{'l', 'I'}: 1,

	BytePair{'T', 'I'}: 1,
	BytePair{'D', 'O'}: 1,
	BytePair{'C', 'G'}: 1,
	BytePair{'L', 'I'}: 1,
	BytePair{'M', 'N'}: 1,
	BytePair{'P', 'B'}: 1,
	BytePair{'F', 'R'}: 1,
	BytePair{'U', 'O'}: 1,
	BytePair{'U', 'V'}: 1,
	BytePair{'E', 'F'}: 1,
	BytePair{'V', 'W'}: 1,
	BytePair{'X', 'Y'}: 1,
}

// WagnerFischer computes the Levenshtein Distance using the Wagner-Fisher algorithm
func WagnerFischer(aStr, bStr string, icost, dcost, scost int) int {

	// Convert to runes to support multibyte characters
	a := []rune(aStr)
	b := []rune(bStr)

	// Allocate both rows.
	row1 := make([]int, len(b)+1)
	row2 := make([]int, len(b)+1)
	var tmp []int

	// Initialize the first row.
	for i := 1; i <= len(b); i++ {
		row1[i] = i * icost
	}

	// For each row...
	for i := 1; i <= len(a); i++ {
		row2[0] = i * dcost

		// For each column...
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				row2[j] = row1[j-1]
			} else {
				ins := row2[j-1] + icost
				del := row1[j] + dcost
				sub := row1[j-1] + scost

				if ins < del && ins < sub {
					row2[j] = ins
				} else if del < sub {
					row2[j] = del
				} else {
					row2[j] = sub
				}
			}
		}

		// Swap the rows at the end of each row.
		tmp = row1
		row1 = row2
		row2 = tmp
	}

	// Because we swapped the rows, the final result is in row1 instead of row2.
	return row1[len(row1)-1]
}

// WagnerFischerWithWeightedSubs computes the Levenshtein Distance with substitution weights given as a map of bytes to maps of bytes to int
func WagnerFischerWithWeightedSubs(a, b string, icost, dcost, scost int, substitutionWeights map[BytePair]int) int {
	// Allocate both rows.
	row1 := make([]int, len(b)+1)
	row2 := make([]int, len(b)+1)
	var tmp []int

	// Initialize the first row.
	for i := 1; i <= len(b); i++ {
		row1[i] = i * icost
	}

	// For each row...
	for i := 1; i <= len(a); i++ {
		row2[0] = i * dcost

		// For each column...
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				row2[j] = row1[j-1]
			} else {
				ins := row2[j-1] + icost
				del := row1[j] + dcost
				finalScost := scost
				if weight, foundPair := substitutionWeights[BytePair{a[i-1], b[j-1]}]; foundPair {
					finalScost = weight
				}
				if weight, foundPair := substitutionWeights[BytePair{b[j-1], a[i-1]}]; foundPair {
					if weight < finalScost {
						finalScost = weight
					}
				}
				sub := row1[j-1] + finalScost

				if ins < del && ins < sub {
					row2[j] = ins
				} else if del < sub {
					row2[j] = del
				} else {
					row2[j] = sub
				}
			}
		}

		// Swap the rows at the end of each row.
		tmp = row1
		row1 = row2
		row2 = tmp
	}

	// Because we swapped the rows, the final result is in row1 instead of row2.
	return row1[len(row1)-1]
}
