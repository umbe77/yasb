PKGS := ./actions ./config ./engine ./handlers ./models ./stack ./store
GOTESTSUM := go tool gotestsum

.PHONY: test test-ci

# Human-friendly test run, never cached.
test:
	$(GOTESTSUM) --format testname -- -count=1 $(PKGS)

# CI run: JUnit XML + coverage profile in reports/.
test-ci:
	mkdir -p reports
	$(GOTESTSUM) --format pkgname --junitfile reports/junit.xml -- -count=1 -coverprofile=reports/coverage.out $(PKGS)
