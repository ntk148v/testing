import requests
import urllib3

urllib3.disable_warnings()


# My repo
repo = "ntk148v/ntk148v"
FLUSH_PKT = "0000"
CAPABILITY_AGENT = "agent"
CAPABILITY_LS_REFS = "ls-refs"
CAPABILITY_FETCH = "fetch"

session = requests.Session()
# Init client request
resp = session.get(f"https://github.com/{repo}/info/refs",
                   params={"service": "git-upload-pack"},
                   headers={"Git-Protocol": "version=2"},
                   verify=False)
# Response in pkt-line format
pkt_lines = resp.content.decode('utf-8').split('\n')
for index, line in enumerate(pkt_lines):
    if line == FLUSH_PKT:
        del pkt_lines[index]
        continue
    pkt_lines[index] = line[4:]
print(pkt_lines)
# WIP :<
