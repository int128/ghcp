GHCP := ../ghcp
GHCP_FLAGS := --debug -u int128 -r sandbox

CIRCLE_BUILD_NUM ?= 0
BRANCH_PREFIX := ghcp-ci-$(CIRCLE_BUILD_NUM)

.PHONY: test
test:
	# create master branch
	$(GHCP) commit $(GHCP_FLAGS) -b $(BRANCH_PREFIX)-master --no-parent -m "Create master branch" fixture1.txt
	# create feature branch from master branch and add a file
	$(GHCP) commit $(GHCP_FLAGS) -b $(BRANCH_PREFIX)-feature --parent $(BRANCH_PREFIX)-master -m "Add fixture2.txt" fixture2.txt
	# add a file to feature branch
	$(GHCP) commit $(GHCP_FLAGS) -b $(BRANCH_PREFIX)-feature -m "Add fixture3.txt" fixture3.txt
	# add a file to feature branch but do nothing
	$(GHCP) commit $(GHCP_FLAGS) -b $(BRANCH_PREFIX)-feature -m "Add fixture1.txt" fixture1.txt
	# create a pull request from feature to master
	$(GHCP) pull-request $(GHCP_FLAGS) -b $(BRANCH_PREFIX)-feature --base $(BRANCH_PREFIX)-master --title "Test: feature -> master"
	# create a release
	$(GHCP) release $(GHCP_FLAGS) -t $(BRANCH_PREFIX) fixture1.txt
