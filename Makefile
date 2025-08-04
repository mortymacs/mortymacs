all:
	nix-shell --command "hugo server"

clean:
	rm -rf public/
