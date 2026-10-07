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

  shellHook = ''
    export PORT=9808
    export DATABASE_URL="postgres://urlittle:urlittle@localhost:5432/urlittle"
    export CACHE_SIZE=1000
  '';
}
