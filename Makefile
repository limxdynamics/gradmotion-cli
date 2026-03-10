.PHONY: build build-all release-local clean install test npm-version npm-pack npm-publish release-all help

# 默认：编译当前平台
build:
	go build -o gm ./cmd/gradmotion

# 编译所有平台
build-all:
	./scripts/build-all.sh

# 本地打包（GoReleaser snapshot）
release-local:
	./scripts/release-local.sh

# 清理产物
clean:
	rm -rf dist/ gm gm-*

# 安装到本地（macOS）
install: build
	sudo install -m 0755 gm /usr/local/bin/gm

# 运行测试
test:
	go test -v ./...

# npm: 同步版本号（用法：VERSION=0.1.0 make npm-version）
npm-version:
	cd npm && npm version $(VERSION) --no-git-tag-version --allow-same-version

# npm: 本地打包预览
npm-pack:
	cd npm && npm pack

# npm: 发布到 npm registry
npm-publish:
	cd npm && npm publish --access public

# 一键发布（版本对齐 + GitLab Release + npm）
# 用法：VERSION=0.1.3 make release-all
release-all:
	@set -e; \
	if [ -z "$(VERSION)" ]; then \
		echo "ERROR: VERSION is required. Example: VERSION=0.1.3 make release-all"; \
		exit 1; \
	fi; \
	if ! echo "$(VERSION)" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$$'; then \
		echo "ERROR: VERSION must be semver like 0.1.3"; \
		exit 1; \
	fi; \
	if [ -z "$$GITLAB_TOKEN" ]; then \
		echo "ERROR: GITLAB_TOKEN is not set"; \
		exit 1; \
	fi; \
	if ! npm --prefix $(CURDIR)/npm whoami >/dev/null 2>&1; then \
		echo "ERROR: npm is not logged in. Run: npm login"; \
		exit 1; \
	fi; \
	if [ -n "$$(git status --porcelain)" ]; then \
		echo "ERROR: git workspace is not clean. Please commit/stash first."; \
		git status --short; \
		exit 1; \
	fi; \
	echo ">> Sync npm version to $(VERSION)"; \
	npm --prefix $(CURDIR)/npm version $(VERSION) --no-git-tag-version --allow-same-version; \
	echo ">> Commit + push prod"; \
	git add $(CURDIR)/npm/package.json; \
	git commit -m "Release v$(VERSION)"; \
	git push origin prod; \
	echo ">> Create + push tag v$(VERSION)"; \
	git tag -a v$(VERSION) -m "Release v$(VERSION)"; \
	git push origin v$(VERSION); \
	echo ">> Publish GitLab release"; \
	goreleaser release --clean --skip=validate; \
	echo ">> Publish npm package"; \
	(cd $(CURDIR)/npm && npm publish --access public); \
	echo "DONE: release v$(VERSION) published to GitLab and npm."

# 帮助
help:
	@echo "Gradmotion CLI Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build          - 编译当前平台"
	@echo "  make build-all      - 编译所有平台（无需 git）"
	@echo "  make release-local  - GoReleaser 本地打包（需 git）"
	@echo "  make clean          - 清理产物"
	@echo "  make install        - 安装到 /usr/local/bin（需 sudo）"
	@echo "  make test           - 运行测试"
	@echo "  make npm-version    - 同步 npm 包版本（需 VERSION=x.y.z）"
	@echo "  make npm-pack       - 本地打包预览 npm 包"
	@echo "  make npm-publish    - 发布到 npm registry"
	@echo "  make release-all    - 一键发布（需 VERSION=x.y.z 和 GITLAB_TOKEN）"
	@echo ""
	@echo "环境变量："
	@echo "  VERSION=v0.1.0 make build-all   - 指定版本号"
	@echo "  VERSION=0.1.0 make npm-version  - 同步 npm 版本"
	@echo "  GITLAB_TOKEN=xxx VERSION=0.1.3 make release-all"
