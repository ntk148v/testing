import copy
import json
import sys

from ansible.parsing.dataloader import DataLoader

try:
    from ansible.inventory.manager import InventoryManager
    A24 = True
except ImportError:
    from ansible.vars import VariableManager
    from ansible.inventory import Inventory
    A24 = False
import etcd3

def put_val(etcdclient, prefix, input):
    prefix = prefix.rstrip("/")
    if isinstance(input, dict):
        for k, v in input.items():
            tmpprefix = copy.deepcopy(prefix)
            tmpprefix = "/".join([tmpprefix, k])
            put_val(etcdclient, tmpprefix, v)
    elif isinstance(input, list):
        for i, v in enumerate(input):
            tmpprefix = copy.deepcopy(prefix)
            tmpprefix = "/".join([tmpprefix, str(i)])
            put_val(etcdclient, tmpprefix, v)
    else:
        print(prefix)
        etcdclient.put(prefix, str(input))

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

# Init etcd3 client
etcdclient = etcd3.client(host="10.4.4.235", port="8379")
put_val(etcdclient, "/test", out)

etcdout = etcdclient.get_prefix("/test")
for k, v in etcdout:
    print(k)
    print(v)

# print(json.dumps(out, indent=4, sort_keys=True))
