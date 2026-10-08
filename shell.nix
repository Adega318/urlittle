{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  packages = with pkgs; [
    go_1_27
    gotools
    gcc
    (pkgs.writeShellScriptBin "gor" "exec go run .")
    (writeShellScriptBin "up" ''
      case "$1" in
        dev)
          docker compose \
            -f docker-compose.yml \
            -f docker-compose_local.yml \
            up -d postgres
          ;;
        "")
          docker compose up -d
          ;;
        *)
          echo "Usage: up {dev|<empty>}"
          exit 1
          ;;
      esac
    '')
    (writeShellScriptBin "down" ''
      docker compose down
    '')
  ];

  shellHook = ''
    export PORT=8080
    export DATABASE_URL="postgres://urlittle:urlittle@localhost:5432/urlittle"
    export CACHE_SIZE=1000
    export POSTGRES_DB=urlittle
    export POSTGRES_USER=urlittle
    export POSTGRES_PASSWORD=urlittle
  '';
}
