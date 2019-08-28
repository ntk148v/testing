import sys
import json

from ansible.parsing.dataloader import DataLoader

try:
    from ansible.inventory.manager import InventoryManager
    A24 = True
except ImportError:
    from ansible.vars import VariableManager
    from ansible.inventory import Inventory
    A24 = False

loader = DataLoader()
if A24:
    inventory = InventoryManager(loader, [sys.argv[1]])
    inventory.parse_sources()
else:
    variable_manager = VariableManager()
    inventory = Inventory(loader, variable_manager, sys.argv[1])
    inventory.parse_inventory(inventory.host_list)

out = {'_meta': {'hostvars': {}}}
for group in inventory.groups.values():
    out[group.name] = {
        'hosts': [h.name for h in group.hosts],
        'vars': group.vars,
        'children': [c.name for c in group.child_groups]
    }
for group in inventory.groups.values():
    for c in group.child_groups:
        out[group.name]['hosts'] += out[c.name]['hosts']
    set(out[group.name]['hosts'])
for host in inventory.get_hosts():
    out['_meta']['hostvars'][host.name] = host.vars

print(json.dumps(out, indent=4, sort_keys=True))
