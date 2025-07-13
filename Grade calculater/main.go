package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your name: ")
	studentName, _ := reader.ReadString('\n')
	studentName = strings.TrimSpace(studentName)

	fmt.Print("Enter the number of subjects: ")
	numOfSubject, _ := reader.ReadString('\n')
	numOfSubject = strings.TrimSpace(numOfSubject)
	numOfSubjectInt, _ := strconv.Atoi(numOfSubject)

	grades := make(map[string]float64)
	for i := 0; i < numOfSubjectInt; i++ {
		fmt.Print("Enter the subject name: ")
		subject, _ := reader.ReadString('\n')
		subject = strings.TrimSpace(subject)

		fmt.Print("Enter grade for this subject: ")
		gradeStr, _ := reader.ReadString('\n')
		gradeStr = strings.TrimSpace(gradeStr)
		gradeInt, _ := strconv.ParseFloat(gradeStr, 64)

		if gradeInt < 0 || gradeInt > 100 {
			fmt.Println("Invalid grade entered!")
			i--
			continue
		}
		grades[subject] = gradeInt

	}

	averageGrade := gradeCalculater(grades)
	fmt.Printf("Grade report of: %s\n", studentName)

	for subject, grade := range grades {
		fmt.Printf("%s: %.2f\n", subject, grade)
	}
	fmt.Printf("Average Grade: %.2f\n", averageGrade)
}

func gradeCalculater(grades map[string]float64) float64 {
	var n = len(grades)
	var sum float64 = 0.0

	if n == 0 {
		return 0
	}

	for _, grade := range grades {
		sum = sum + grade
	}
	return sum / float64(n)
}
