# AGENTS.md

1. after each changes add all changes using git add . only and gemerate a single line, single sentence concise but understandable commit message of the cjange made. also push
 don't um don't addum don't add any details about the commit author or the commits don't add new lines just the commit message a single sentence
 ensure that in the process of adding a feature you do not break any part of the core code base 
 ensure that everything every feature is been wired up and everything is a closed source code so that every feature is interactive in a closed system of the entire source code so that all the sub commands are wired up not um duplicating functionality everything is wired up in a centralized um system of um so commands command which everything work in together in a very very close loop of um interchanging um interactions between the call functionalities and the um authentication methods egola um everything is wired up replaceably and more lally can be in that changed can be used and still power the coffee with us every single of the um entire tool

---

# Additional Agent Guidelines

1. **Build:** Use `go build ./...` to build all packages. For multi-platform builds, use `./build.sh`.
2. **Test:** Run all tests with `go test ./...`. To run a single test: `go test -run TestFunctionName ./...`.
3. **Lint:** Use `go vet ./...` for static analysis. `golint` is not required. No custom lint config found.
4. **Format:** Format code with `gofmt -s -w .` before committing.
5. **Imports:** Group standard, third-party, and local imports. Use Go's default import order.
6. **Types:** Prefer explicit types. Use Go idioms for interfaces and structs.
7. **Naming:** Use CamelCase for types, functions, and exported variables. Use lowerCamelCase for unexported variables.
8. **Error Handling:** Always check errors. Return early on error. Use `fmt.Errorf` for wrapping.
9. **Comments:** Use full sentences. Exported functions/types must have doc comments.
10. **Tests:** Place tests in `*_test.go` files. Use table-driven tests where possible.
11. **Modules:** Use Go modules (`go.mod`). Keep dependencies up to date.
12. **CLI:** All integrations self-register and are auto-discovered by the CLI.
13. **Centralization:** Wire up all features via the central registry; avoid duplicate functionality.
14. **Commit:** After changes, run `git add .` and commit with a single, concise sentence (no author or extra lines).
15. **No Breaking Changes:** Ensure new features do not break core functionality.
16. **Closed System:** All subcommands and integrations must interact in a tightly-coupled, replaceable, and extensible manner.
17. **No Unused Code:** Remove dead code and unused imports.
18. **Documentation:** Update `README.md` and `USAGE.md` for new features or commands.
19. **Directory Structure:** Follow the modular structure as described in `README.md`.
20. **Contribution:** Follow these guidelines for all PRs and code reviews.
21. NEVER edit .gitignore
22. go through opencode/ project directory. it is an existent tool similar to the current project (codeforgeai.go) and we will be taking cues and architectural logics from it especially in the auth management system, and agent tooling (but not to affect our current codeforgeai.go architecture, but to enhance it). do NOT under any circumstances, modify any of the opencode/ source code. we are simply using it as a reference for building codeforgeai auth/agent system.