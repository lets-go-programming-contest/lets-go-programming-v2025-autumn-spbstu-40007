package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	const (
		lowerBoundary = 15
		upperBoundary = 30
	)

	reader := bufio.NewReader(os.Stdin)
	departmentsTemp, _ := reader.ReadString('\n')
	departments, _ := strconv.Atoi(strings.TrimSpace(departmentsTemp))

	for d := 0; d < departments; d++ {
		employeesTemp, _ := reader.ReadString('\n')
		employees, _ := strconv.Atoi(strings.TrimSpace(employeesTemp))

		lowerBorder, upperBorder := lowerBoundary, upperBoundary
		hasError := false
		processedEmployees := 0

		for e := 0; e < employees; e++ {
			processedEmployees = e

			if hasError {
				fmt.Println(-1)
				reader.ReadString('\n') //nolint:errcheck
				//сделано для того, чтобы просто пропустить ввод в случае чего, но линтер ругается на то, что не ловлю ошибку, поэтому считаю, что тут можно это сделать.

				continue
			}
			settingsTemp, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(-1)
				hasError = true

				continue
			}
			settings := strings.TrimSpace(settingsTemp)
			parts := strings.Fields(settings)
			if len(parts) < 2 {
				fmt.Println(-1)
				hasError = true //nolint:wsl
				//линтер ругается на то, что присваивание не стоит вместе с другими присваиваниями, оно тут не нужно и я, честно говоря, и не знаю как ему угодить, так что пишу нолинт.

				continue
			}
			sign, temperature := func() (string, int) {
				temperature, err := strconv.Atoi(parts[1])
				if err != nil {
					fmt.Println(-1)
					hasError = true
				}
				return parts[0], temperature
			}()
			//сразу поясню, для чего тут это. Линтер ругался на то, что у меня до if statment'a идут два присваивания, и поэтому, долго не ломая голову, решил сделать так, как у меня было изначально в каком-то из первых коммитов.

			if sign == ">=" {
				lowerBorder = max(lowerBorder, temperature)
			} else if sign == "<=" {
				upperBorder = min(upperBorder, temperature)
			}

			if lowerBorder > upperBorder {
				fmt.Println(-1)
				hasError = true
			} else {
				fmt.Println(lowerBorder)
			}
		}
		if hasError {
			remaining := employees - (processedEmployees + 1)
			for i := 0; i < remaining; i++ {
				reader.ReadString('\n') //nolint:errcheck
				//такая же история, как и в предыдущем случае
				fmt.Println(-1)
			}
		}
	}
}
