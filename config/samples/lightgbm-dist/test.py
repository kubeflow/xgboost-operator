import argparse
import subprocess


cmd = ["sh", "./test.sh"]

proc = subprocess.Popen(cmd, stdout=subprocess.PIPE)
while (line := proc.stdout.readline()) :
    print(line.decode("utf-8").strip())
