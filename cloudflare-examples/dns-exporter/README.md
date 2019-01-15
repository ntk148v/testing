# CloudFlare backup DNS script

## How to run a script

* Prepare a Python environment.
* Install a requirement package - [cloudflare](https://github.com/cloudflare/python-cloudflare).

```
$ sudo pip install -r requirements.txt
```

* Take a look at [example-vars.env](./example-vars.env) and export these environment variables.

* Run it!

```
$ BACKUP_FILE_PATH=/path/to/backup_file LOG_CFG=/path/to/logging python backup_dns_records.py
# or just
$ python backup_dns_records.py
```