workflow "lint" {
  resolves = ["golangci-lint"]
  on = "push"
}

action "golangci-lint" {
  uses = "actions-contrib/golangci-lint@master"
  args = "run"
}
