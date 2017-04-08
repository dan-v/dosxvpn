#!/bin/bash

platypus -y --quit-after-execution --app-icon 'appicon.icns'  --name 'dosxvpn'  --interface-type 'None'  --interpreter '/bin/bash'  --bundled-file '../dosxvpn' --bundled-file '../static'  'run.sh' $@
