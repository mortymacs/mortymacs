{ pkgs ? import <nixpkgs> {} }:
with pkgs;
mkShell {
    buildInputs = [
        zola
    ];
    shellHook = ''
        echo "Run:"
        echo 'zola serve'
    '';
}
