#!/bin/sh

ps aux | grep hello | grep -v grep | awk '{print $1}' | xargs kill -9
