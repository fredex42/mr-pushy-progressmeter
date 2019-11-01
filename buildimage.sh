#!/usr/bin/env bash -e

declare -x BASEPATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )

cd $BASEPATH/frontend && npm run build
cd $BASEPATH/server && rm -f server && GOOS=linux go build

BUILDPATH="$BASEPATH/.buildtemp"

echo Build location is $BUILDPATH

if [ -d "$BUILDPATH" ]; then
    echo Clearing out old build from $BUILDPATH
    rm -rf "$BUILDPATH"
fi

mkdir -p $BUILDPATH/tar/public
cp -av $BASEPATH/server/public/     $BUILDPATH/tar/public/
cp -v $BASEPATH/server/server       $BUILDPATH/tar

cd $BUILDPATH/tar
tar cv * | gzip > $BUILDPATH/servercontent.tar.gz

echo Compression completed
cp $BASEPATH/Dockerfile $BUILDPATH
cd $BUILDPATH
docker build . -t guardianmultimedia/mr-pushy-progressmeter

cd $BASEPATH
rm -rf $BUILDPATH

