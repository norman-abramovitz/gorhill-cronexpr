# cronexpr — build, test and release.
#
# Branches: work flows only into develop, and only develop feeds main.
# Feature branches merge into develop. main receives develop as one
# squash commit whose tree is develop's, and releases are tagged there.
# Never commit on main, never merge main into develop. Tags are vX.Y.Z
# (the module path carries /v2, so Go resolves v2 tags).
#
# Release notes come from changelog.d/ fragments (make changelog-new on
# each feature branch). To release:
#   git switch main
#   git read-tree -u --reset develop     main's tree becomes develop's
#   git commit                           the squash commit
#   make bump minor                      tag it, the notes in the tag body
#   git push fork main vX.Y.Z
# then sweep on a branch into develop (make changelog-sweep); the next
# release carries the sweep to main.

# Tags sit on main only, so develop cannot reach them: take the newest
# vX.Y.Z tag in the repo, not the nearest reachable one, then the next
# patch, as version.mk's default does.
VERSION_CMD := git tag -l 'v[0-9]*' --sort=-v:refname | head -n 1 | sed -e 's/^v//' | awk -F. 'NF==3 && $$3 ~ /^[0-9]+$$/ {print $$1"."$$2"."$$3+1}'

PROJECT := cronexpr
BIN     := bin/$(PROJECT)

include mk/snippets/version.mk
include mk/snippets/snippets.mk
include mk/snippets/bump.mk
include mk/snippets/changelog.mk
include mk/snippets/security.mk
include mk/snippets/help.mk

# Tag with the assembled fragments as the message; refuse to tag a release
# with no notes. bump.mk calls this with the new version as its argument.
VERSION_WRITE_CMD = sh -c 'notes=$$($(MAKE) --no-print-directory changelog-assemble) || exit 1; \
	[ -n "$$notes" ] || { echo "no changelog.d fragments: refusing to tag $$0" >&2; exit 1; }; \
	printf "%s %s\n\n%s\n" "$(PROJECT)" "$$0" "$$notes" | git tag -a -F - "$$0"'

.DEFAULT_GOAL := help
.PHONY: build test bench vet clean version

##@ Build

build: ## Build the cronexpr command-line tool into bin/
	go build -o $(BIN) ./cronexpr

clean: ## Remove build output
	rm -f $(BIN)

##@ Test

test: ## Run the tests with the race detector
	go test -race -count=1 ./...

bench: ## Run the benchmarks
	go test -run '^$$' -bench . -benchmem ./...

vet: ## Run go vet
	go vet ./...

##@ Introspection

version: ## Print resolved version variables
	$(dump.version)
