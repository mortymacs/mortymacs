{ pkgs ? import <nixpkgs> {} }:
with pkgs;
mkShell {
    buildInputs = [
        hugo
    ];
    shellHook = ''
        echo "Run:"
        echo 'make'
    '';
}
