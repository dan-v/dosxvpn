#!/bin/bash

./platypus -y --admin-privileges --quit-after-execution --app-icon 'appicon.icns'  --name 'dosxvpn'  --interface-type 'None'  --interpreter '/bin/sh'  --bundled-file '../build/osx/x86-64/dosxvpn' --bundled-file '../static'  'run.sh' ../build/osx/x86-64/dosxvpn
