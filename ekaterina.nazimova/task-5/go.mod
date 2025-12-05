module github.com/UwUshkin/task-5

go 1.22.7


require (
    github.com/stretchr/testify v1.11.1
    golang.org/x/sync v0.11.0
)

replace github.com/UwUshkin/task-5/pkg/conveyer => ./pkg/conveyer
replace github.com/UwUshkin/task-5/pkg/handlers => ./pkg/handlers