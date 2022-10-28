.PHONY: molecule_test_%
molecule_test_%:
	cd ansible/; molecule test -s $*

.PHONY: molecule_dev_%
molecule_dev_%:
	cd ansible/; molecule test -s $* --destroy=never

.PHONY: molecule_converge_%
molecule_converge_%:
	cd ansible/; molecule converge -s $*

.PHONY: molecule_verify_%
molecule_verify_%:
	cd ansible/; molecule verify -s $*

.PHONY: molecule_destroy_%
molecule_destroy_%:
	cd ansible/; molecule destroy -s $*

.PHONY: molecule_dev molecule_converge molecule_test molecule_verify molecule_destroy
molecule_dev: molecule_dev_default
molecule_converge: molecule_converge_default
molecule_test: molecule_test_default
molecule_verify: molecule_verify_default
molecule_destroy: molecule_destroy_default

.PHONY: molecule
molecule: molecule_dev
	@echo "Molecule ran through. For test driven development you can now write"
	@echo "your test. Afterwards you can run molecule_verify to see the expected"
	@echo "fail. Once the ansible code is added it could be applied with "
	@echo "molecule_converge. Another molecule_verify should succesfully run now"
	@echo "molecule_destroy will destroy your instances and the secutiry groups"
