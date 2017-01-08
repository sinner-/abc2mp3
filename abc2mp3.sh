#!/usr/bin/env bash

URL="http://abcradiomodhls.abc-cdn.net.au/i/triplej/audio"
SHOW="hip"
NUM="1"
DATE="2016-05-19"

for segment in `seq 1 1081`;
do
    echo -n "Downloading segment $segment..."
    wget -q $URL/$SHOW-$NUM-$DATE.m4a/segment"$segment"_0_a.ts
    if [ $? == 0 ]; then
        echo "successful."
    else
        echo "couldn't find any more segments, breaking."
        break
    fi
    echo "file 'segment"$segment"_0_a.ts'" >> segmentlist
done

if [ -f segmentlist ]; then
    ffmpeg -f concat -i segmentlist -acodec copy output.ts
    ffmpeg -i output.ts -f mp3 -acodec mp3 $SHOW-$NUM-$DATE.mp3
    rm *.ts
    rm segmentlist
fi
