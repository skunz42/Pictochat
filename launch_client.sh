#!/bin/sh

tmux new-session -y 40 -d -s 'main';
tmux send-keys 'go run ./src/client_rec.go' C-m;
tmux split-window -v -p 5;
tmux send-keys 'go run ./src/client_in.go' C-m;
tmux attach-session -d -t main
