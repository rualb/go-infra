import os
import shutil
import subprocess
import sys

"""
git init
git add .
git commit -m ""
git tag "$(cat VERSION)"

#BINARY_NAME = "app.exe" if os.name == "nt" else "app"

local release 
goreleaser build --snapshot

"""

def test():
    print("Testing...")
    env = os.environ.copy()
    #env['CGO_ENABLED'] = '1' #for -race flag #cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in %PATH%
    command = ['go', 'test'
               #, '-race'
               , '-timeout=60s', '-count=1', './...']
    subprocess.run(command, env=env) #, "-v"

def help():
    print("Usage:")
    print("  python build.py test     - Run test")
    print("  python build.py help     - Display this help message")
    
def build():
    print("Building the binary...")
    subprocess.run(["go", "build", "-C", "cmd/go-infra", "-o",f"./../../dist/" ])
 
def run():
    print("Building the binary...")
    subprocess.run(["dist/go-infra", "-config", f"./configs" ]) #go-proxy -config ./configs/go-proxy
def lint():
    print("Linter...")
    subprocess.run(["golangci-lint ", "run"])
 
if len(sys.argv) > 1:
    command = sys.argv[1]
    if command == "test":
        test() 
    elif command == "help":
        help() 
    elif command == "build":
        build() 
    elif command == "run":
        run() 
    elif command == "lint":
        lint() 
    else:
        help()
        exit(1)
else:
    help()









