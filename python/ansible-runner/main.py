import shutil

import ansible_runner


def artifacts_hanlder(artifacts_dir):
    # Do something here
    print("# Artifacts:", artifacts_dir)


# Run sync
r = ansible_runner.run(private_data_dir="demo",
                       playbook="test.yml",
                       artifacts_handler=artifacts_hanlder)
print("# Status: {} - RC: {}".format(r.status, r.rc))
# successful: 0
print("# Host event:\n")
for each_host_event in r.events:
    print(each_host_event["event"])
print("# Final status:")
print(r.stats)

# Run async using Threading
# Thread & runner
t, r = ansible_runner.run_async(private_data_dir="demo",
                                playbook="test.yml",
                                artifacts_handler=artifacts_hanlder)
# Wait until thread terminates
t.join()
print("Finish roi nha")
print("# Status: {} - RC: {}".format(r.status, r.rc))
# successful: 0
print("# Host event:\n")
for each_host_event in r.events:
    print(each_host_event["event"])
print("# Final status:")
print(r.stats)
