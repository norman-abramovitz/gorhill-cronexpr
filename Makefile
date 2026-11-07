# cronexpr — build, test and release.
#
# Branches: feature branches merge into develop; develop is squash-merged
# into main, and releases are tagged on main only. Tags are vX.Y.Z (the
# module path carries /v2, so Go resolves v2 tags).
#
# Release notes come from changelog.d/ fragments (make changelog-new on
# each feature branch). On main:
#   make bump minor          tag the release, notes in the tag body
#   git push fork vX.Y.Z
#   make changelog-sweep     then commit
# Then merge main back into develop before anything else lands there, and
# run make changelog-sweep again, amending the merge: the squash hides
# main's history, so the merge keeps develop's copies of the fragments.
# git diff main develop should then be empty.

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
