GHCP := ../ghcp
GHCP_FLAGS := --debug -u int128 -r ghcp

GITHUB_RUN_NUMBER ?= 0
BRANCH_PREFIX := ghcp-acceptance-test-$(GITHUB_RUN_NUMBER)
GITHUB_PR_NUMBER := $(word 3,$(subst /, ,$(GITHUB_REF)))

.PHONY: test
test:
	# create master branch
	$(GHCP) commit $(GHCP_FLAGS) -b $(BRANCH_PREFIX)-master --no-parent -m "Create master branch" fixture1.txt
	sleep 1

	# create feature branch from master branch and add a file
	$(GHCP) commit $(GHCP_FLAGS) -b $(BRANCH_PREFIX)-feature --parent $(BRANCH_PREFIX)-master -m "Add fixture2.txt" fixture2.txt
	sleep 1

	# add a file to feature branch
	$(GHCP) commit $(GHCP_FLAGS) -b $(BRANCH_PREFIX)-feature -m "Add fixture3.txt" --author-name="octocat" --author-email="octocat@github.com" fixture3.txt
	sleep 1

	# add a file to feature branch but do nothing
	$(GHCP) commit $(GHCP_FLAGS) -b $(BRANCH_PREFIX)-feature -m "Add fixture1.txt" fixture1.txt
	sleep 1

	# create a pull request from feature to master
	$(GHCP) pull-request $(GHCP_FLAGS) -b $(BRANCH_PREFIX)-feature --base $(BRANCH_PREFIX)-master \
		--title "ghcp: acceptance-test $(BRANCH_PREFIX)" \
		--body "ref: #$(GITHUB_PR_NUMBER)"
	sleep 1

	$(GHCP) pull-request $(GHCP_FLAGS) -b $(BRANCH_PREFIX)-feature --base $(BRANCH_PREFIX)-master \
		--title "ghcp: acceptance-test: This should not be created" \
		--body "ref: #$(GITHUB_PR_NUMBER)"

.PHONY: clean-up
clean-up:
	-git push origin --delete $(BRANCH_PREFIX)-master
	-git push origin --delete $(BRANCH_PREFIX)-feature
