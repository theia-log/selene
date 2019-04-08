#!/usr/bin/env bash

version_tag="$1"

package="github.com/theia-log/selene"
package_split=(${package//\// })
package_name=${package_split[-1]}

build_dir="_build"

platforms=("linux/amd64" "linux/arm" "linux/arm64" "windows/amd64" "darwin/amd64")

if [ ! -d "$build_dir" ]; then
    mkdir "$build_dir"
fi


for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name="$package_name"
    if [ ! -z "$version_tag" ]; then
        output_name+='-'$version_tag
    fi
    output_name=$output_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi  

    echo "Building $output_name"
    env GOOS=$GOOS GOARCH=$GOARCH go build -o "$build_dir/$output_name" $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done

echo "All executables built successfully."
