prereqs="pocketbase go node"

echo 'checking for prereqs'

for cmd in $prereqs; do
    if which -v $cmd > /dev/null; then
        echo "$cmd is installed"
    else
        echo "please install $cmd and try running this script again"
    fi
done