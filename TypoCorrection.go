package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"math"
)

func computeTransformTable(x string, y string, cCost int, rCost int, dCost int, iCost int) ([][]int, [][]string) {
	m := len(x)+1
	n := len(y)+1
	cost := make([][]int, m)
	for row := range cost {
		cost[row] = make([]int, n)
	}

	op := make([][]string, m)
	for row := range op {
		op[row] = make([]string, n)
	}

	cost[0][0] = 0;

	for i := 1; i < m; i++ {
		cost[i][0] = i * dCost
		op[i][0] = "d" + string(([]rune(x))[i-1]) 
	}

	for j := 1; j < n; j++ {
		cost[0][j] = j * iCost
		op[0][j] = "i" + string(([]rune(y))[j-1]) 
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if x[i-1] == y[j-1] /*&& x[i] != " "*/ { // if should copy
				cost[i][j] = cost[i-1][j-1] + cCost
				op[i][j] = "c" + string(([]rune(x))[i-1])
			} else {
				cost[i][j] = cost[i-1][j-1] + int(considerMisspellings(([]rune(x))[i-1], ([]rune(y))[j-1]))
				op[i][j] = "r" + string(([]rune(x))[i-1]) + "->" + string(([]rune(y))[j-1])
			}

			if cost[i-1][j] + dCost < cost[i][j] {
				cost[i][j] = cost[i-1][j] + dCost
				op[i][j] = "d" + string(([]rune(x))[i-1])
			}

			if cost[i][j-1] + iCost < cost[i][j] {
				cost[i][j] = cost[i][j-1] + iCost
				op[i][j] = "i" + string(([]rune(y))[j-1])
			}
		}
	}

	return cost, op
}

/* given the operation table from the compute transform table output
 * performed on strings s1 and s2 and the ending indices of the op table,
 * return a sequence of needed operations to transform s1 to s2
 */
func assembleTransformation(op [][]string, i int, j int) string {
	if i == 0 && j == 0 {
		return ""
	}

	if ([]rune(op[i][j]))[0] == 'c' || ([]rune(op[i][j]))[0] == 'r' {
		return assembleTransformation(op, i-1, j-1) + " " + op[i][j] // ?????
	} else if ([]rune(op[i][j]))[0] == 'd' {
		return assembleTransformation(op, i-1, j) + " " + op[i][j]
	} else {
		return assembleTransformation(op, i, j-1) + " " + op[i][j]
	}
}

/* given a sequence of operations to transform one string to another
 * calculate the cost where copy = 0, replace = dependent on 
 * considerMisspellings, delete/insert = +3; use to compare w/ computed
 * lowest cost from last element in computed cost table
 */
func debugComputeCost(opSequence string) int {
	cost := 0
	seqList := strings.Split(strings.TrimSpace(opSequence), " ")
	for _, seq := range seqList {
		runes := []rune(seq)
		if runes[0] == 'r' {
			// runes[2, 3] is "->"
			cost += int(considerMisspellings(runes[1], runes[4]))
		} else if runes[0] != 'c' {
			cost += 3
		}
	}
	return cost
}

// return the cost of transforming a to b considering typo probability
func considerMisspellings(a rune, b rune) float64 {
	iA := int(a)
	iB := int(b)
	return math.Sqrt(math.Abs(float64(iA)-float64(iB))) + float64(1)
}

func main() {
	// collect args
	stringsToConvert := os.Args[1:] // remove first arg which is program name

    // open file & returns os.File type pointer (or an error)
    file, error := os.Open("words")

    // check for error
    if error != nil {
        fmt.Printf("Error while reading the file %v\n", error)
    }

    scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var words []string

	for scanner.Scan() {
		words = append(words, strings.TrimSpace(scanner.Text()))
	}

    // close the file
    file.Close()

	/* now have list of strings to convert & possible strings to convert to
	 * find lowest cost word
	 */
	for _, s1 := range stringsToConvert {
		lowestCostWord := words[0] // set lowestCostWord to first word
		cost, op := computeTransformTable(s1, lowestCostWord, 0, -99, 3, 3)
		assembleTransformation(op, len(s1), len(lowestCostWord))
		lowestCost := cost[len(s1)][len(lowestCostWord)] // set lowestCost to cost for first word
		
		for i := 1; i < len(words); i++ {
			s2 := words[i]
			m := len(s1)+1
			n := len(s2)+1
			cost, op := computeTransformTable(s1, s2, 0, -99, 3, 3)
			assembleTransformation(op, m-1, n-1)
			computedCost := cost[m-1][n-1]
			if  computedCost < lowestCost {
				lowestCostWord = s2
				lowestCost = computedCost
			}
		}

		fmt.Printf("%s -> %s with cost %d\n", s1, lowestCostWord, lowestCost)
	}
}
