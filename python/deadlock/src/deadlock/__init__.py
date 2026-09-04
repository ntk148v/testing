import subprocess
import threading

import pymysql

# Two transactions lock the same two rows in opposite order.
#   seeder: locks id=1 ...... then wants id=2
#   reader: locks id=2 ...... then wants id=1
# -> wait cycle -> InnoDB detects and kills the victim with ERROR 1213.


def run_seeder():
    sql = (
        # Lock id=1, sleep holding it to widen the overlap window for the slow
        # docker-CLI session, then lock id=2. (No advisory gate / GET_LOCK.)
        "START TRANSACTION; "
        "UPDATE test_table SET val='s' WHERE id=1; "
        "SELECT SLEEP(8); "
        "UPDATE test_table SET val='s' WHERE id=2; "
        "COMMIT;"
    )
    proc = subprocess.run(
        [
            "docker",
            "exec",
            "-i",
            "mysql",
            "mysql",
            "-uroot",
            "-prootpassword",
            "my_new_db",
            "-e",
            sql,
        ],
        capture_output=True,
        text=True,
    )
    print(f"[seeder] rc={proc.returncode} err={proc.stderr.strip()!r}")


def run_reader():
    sql = (
        # START a transaction so the FOR UPDATE locks persist (the mysql CLI
        # default autocommit would otherwise release each lock at statement
        # end). Lock id=2, hold it, then ask for id=1 (held by the seeder) --
        # opposite lock order. (No advisory gate / GET_LOCK.)
        "START TRANSACTION; "
        "SELECT * FROM test_table WHERE id=2 FOR UPDATE; "
        "SELECT SLEEP(1); "
        "SELECT * FROM test_table WHERE id=1 FOR UPDATE; "
        "COMMIT;"
    )
    proc = subprocess.run(
        [
            "docker",
            "exec",
            "-i",
            "mysql",
            "mysql",
            "-uroot",
            "-prootpassword",
            "my_new_db",
            "-e",
            sql,
        ],
        capture_output=True,
        text=True,
    )
    print(f"[reader] rc={proc.returncode} err={proc.stderr.strip()!r}")


def main():
    conn = pymysql.connect(
        host="127.0.0.1", user="root", password="rootpassword", autocommit=False
    )
    try:
        cur = conn.cursor()
        # Build the schema beforehand (committed), so both subprocesses race
        # only on row locks -- not on DDL.
        cur.execute("CREATE DATABASE IF NOT EXISTS my_new_db;")
        cur.execute("DROP TABLE IF EXISTS my_new_db.test_table;")
        cur.execute(
            "CREATE TABLE my_new_db.test_table (id INT PRIMARY KEY, val VARCHAR(50))"
        )
        cur.execute(
            "INSERT INTO my_new_db.test_table (id, val) VALUES (1, 'a'), (2, 'b');"
        )
        conn.commit()
        print("[setup] schema + table ready")

        # Run both sessions concurrently via subprocess.run(); the deadlock
        # comes from lock-order inversion, not from async process launches.
        t2 = threading.Thread(target=run_seeder)
        t3 = threading.Thread(target=run_reader)
        t2.start()
        t3.start()
        t2.join()
        t3.join()
    finally:
        conn.close()


if __name__ == "__main__":
    main()
