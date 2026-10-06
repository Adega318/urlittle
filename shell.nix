{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  packages = with pkgs; [
    go_1_27
    gotools
    gcc
    (pkgs.writeShellScriptBin "gor" "exec go run .")
  ];
}
