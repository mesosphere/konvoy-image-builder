
venv/:
	# lets expect .tool-versions is taken
	python -m venv $@

dev-setup: venv/
	bash $(<)/bin/activate
	pip install -r requirements.txt


molecule: dev-setup
	molecule test --all --destroy always --parallel
