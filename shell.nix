{
    pkgs ? import <nixpkgs> { }
}:
pkgs.mkShell {
    name = "sshresume";
    buildInputs = [
        pkgs.go
        pkgs.cowsay
        pkgs.lolcat
    ];
    shellHook = ''
        echo "sshresume" | cowsay | lolcat
    '';
}
