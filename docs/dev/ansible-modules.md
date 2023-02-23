# Managing Ansible Modules Dependencies

KIB uses Ansible to provision the machine image. Most modules are included in ansible-core, but some are available in collections that we install using ansible-galaxy CLI.

I do not know of a way to reliably enumerate all the necessary modules, but the following script generates a good candidate list. Note that it requires `curl`, `fd`, `gojq`, which you may need to install.

```shell
#!/usr/bin/env bash
# List all ansible keywords
curl -q https://raw.githubusercontent.com/ansible/ansible/devel/lib/ansible/keyword_desc.yml | gojq --yaml-input --raw-output '.. | keys? | flatten[] | select(type=="string")' > ansible_keywords

# List all keys in yml/yaml files in the ansible directory (and subdirectories).
fd . ansible -e 'yml' -e 'yaml' -x gojq --yaml-input --raw-output '.. | keys? | flatten[] | select(type=="string")' > keys

# List every possible module by calculating the set difference (all keys) - (ansible keywords).
diff --new-line-format="" --unchanged-line-format="" keys ansible_keywords | sort | uniq > possible_modules
```

Some non-builtin modules are namespaced, and it's easy to tell the collection to which they belong. For example, `ansible.posix.mount` belongs to `ansible.posix`, and `ansible.utils.cli_parse` belongs to `ansible.utils`.

Other non-builtin modules are NOT namespaced. For example, `redhat_subscription` belongs to the `community.general` module.

Finally, modules can be referenced dynamically. That is, not as a YAML key. For example, the `ansible.netcommon.native` module (which belongs to `ansible.netcommon`) is referenced in a YAML value.

I will look into more reliable methods of enumerating modules in the future. It would be nice to 'dry-run' all tasks, have Ansible throw missing module errors, but continue until it tries all tasks.
