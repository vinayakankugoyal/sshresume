{
    pkgs ? import <nixpkgs> { }
}:
let
  # Define your custom gcloud with the plugin here
  gcloud = pkgs.google-cloud-sdk.withExtraComponents [
    pkgs.google-cloud-sdk.components.gke-gcloud-auth-plugin
  ];
in
pkgs.mkShell {
  # Add the custom variable to your buildInputs
  buildInputs = [
    gcloud
    pkgs.go
    pkgs.cowsay
    pkgs.lolcat
  ];
  shellHook = ''
      echo "sshresume" | cowsay | lolcat
  '';
}
