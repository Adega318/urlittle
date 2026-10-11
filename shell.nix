{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  packages = with pkgs; [
    go_1_27
    gotools
    gcc
    (pkgs.writeShellScriptBin "gor" "exec go run .")
    (pkgs.writeShellScriptBin "pod" "exec podman-compose \"$@\"")
  ];

  shellHook = ''
    # App 
    export LOG_LEVEL=INFO
    export PORT=8080
    # Store
    export TTL_MINUTES=5
    export DATABASE_URL="postgres://urlittle:urlittle@localhost:5432/urlittle"
    export CACHE_SIZE=1000
    # Rate Limiting
    export RATE_LIMIT_USER_LIST_SIZE=5000
    export RATE_LIMIT_LIMIT=10 
    export RATE_LIMIT_BURST=20
    # DB 
    export POSTGRES_DB=urlittle
    export POSTGRES_USER=urlittle
    export POSTGRES_PASSWORD=urlittle
  '';
}
