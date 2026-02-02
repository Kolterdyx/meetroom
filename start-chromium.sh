#!/bin/bash
DISPLAY=:0 chromium \
--kiosk \
--noerrdialogs \
--disable-infobars \
--disable-session-crashed-bubble \
--autoplay-policy=no-user-gesture-required \
--usenable-features=VaapiVideoDecoder \
--remote-debugging-port=9222 \
http://localhost:5000/idle