import json


def get_targets(data):
    results = []
    for i in data['data']['activeTargets']:
        results.append(str(i['discoveredLabels']['__address__']))
    results.sort()
    return set(results)


file1 = open('/tmp/targets1.json')
data1 = json.load(file1)
results1 = get_targets(data1)

file2 = open('/tmp/targets2.json')
data2 = json.load(file2)
results2 = get_targets(data2)

print(results1-results2)
