#!/bin/sh

tmux new-session -y 40 -d -s 'main';
tmux send-keys "go run ./src/client_rec.go $1 $2 $3 $4" C-m;
tmux split-window -v -p 5;
tmux send-keys "go run ./src/client_in.go $1 $2 $3" C-m;
tmux attach-session -d -t main;
