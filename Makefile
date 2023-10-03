.PHONY: install

install-indexer:
	@echo "Installing indexer dependencies..."
	cd ./indexer && go mod download -x && cd ..
	@echo ""

install-server:
	@echo "Installing server dependencies..."
	cd ./server && go mod download -x && cd ..
	@echo ""

install-emails-search-app:
	@echo "Installing emails-search-app dependencies..."
	cd ./emails-search-app && npm install && cd ..
	@echo ""

install: # install dependencies for the project apps
	@make install-indexer
	@make install-server
	@make install-emails-search-app
