{
  description = "urlittle - URL shortener service";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        urlittlePackage = (pkgs.buildGoModule.override { go = pkgs.go_1_27; }) {
          pname = "urlittle";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-7kFMyJUupLl5oZomwGFSCkMiq4Y0ptP0qLM/ZY4bqyw=";
          modSha256 = "sha256-y2eqpR6HPu7zcz5Ff5Ckb8mhj5ZEwAUcsn618fxMkkY=";
          buildFlags = [ "-ldflags=-s -w" ];
        };
      in
      {
        packages = {
          urlittle = urlittlePackage;
          default = urlittlePackage;
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go_1_27
            gotools
            gcc
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
        };

        nixosModules.default =
          {
            config,
            lib,
            ...
          }:
          {
            options = {
              services.urlittle = {
                enable = lib.mkEnableOption "urlittle URL shortener service";
                port = lib.mkOption {
                  type = lib.types.portNumber;
                  default = 8080;
                  description = "Port to listen on";
                };
                ttl = lib.mkOption {
                  type = lib.types.int;
                  default = 5;
                  description = "Time to live of uploaded urls in minutes";
                };
                databaseUrl = lib.mkOption {
                  type = lib.types.str;
                  default = "postgres://urlittle:urlittle@localhost:5432/urlittle";
                  description = "PostgreSQL connection URL";
                };
                cacheSize = lib.mkOption {
                  type = lib.types.int;
                  default = 1000;
                  description = "Cache size for URL lookups";
                };
                logLevel = lib.mkOption {
                  type = lib.types.enum [
                    "DEBUG"
                    "INFO"
                    "WARN"
                    "ERROR"
                  ];
                  default = "INFO";
                  description = "Log level";
                };
                rateLimiting = lib.mkOption {
                  type = lib.types.submodule {
                    options = {
                      userListSize = lib.mkOption {
                        type = lib.types.int;
                        default = 10000;
                        description = "Max number of client IPs to track for rate limiting";
                      };
                      limit = lib.mkOption {
                        type = lib.types.int;
                        default = 10;
                        description = "Requests per second per client";
                      };
                      burst = lib.mkOption {
                        type = lib.types.int;
                        default = 20;
                        description = "Burst allowance per client";
                      };
                    };
                  };
                  default = { };
                  description = "Rate limiting configuration";
                };
              };
            };

            config = lib.mkIf config.services.urlittle.enable {
              systemd.services.urlittle = {
                description = "urlittle URL shortener";
                wantedBy = [ "multi-user.target" ];
                after = [
                  "network.target"
                  "postgresql.service"
                ];
                serviceConfig = {
                  ExecStart = "${urlittlePackage}/bin/urlittle";
                  Restart = "on-failure";
                  Environment = [
                    "PORT=${toString config.services.urlittle.port}"
                    "TTL_MINUTES=${toString config.services.urlittle.ttl}"
                    "DATABASE_URL=${config.services.urlittle.databaseUrl}"
                    "CACHE_SIZE=${toString config.services.urlittle.cacheSize}"
                    "LOG_LEVEL=${config.services.urlittle.logLevel}"
                    "RATE_LIMIT_USER_LIST_SIZE=${toString config.services.urlittle.rateLimiting.userListSize}"
                    "RATE_LIMIT_LIMIT=${toString config.services.urlittle.rateLimiting.limit}"
                    "RATE_LIMIT_BURST=${toString config.services.urlittle.rateLimiting.burst}"
                  ];
                };
              };
            };
          };
      }
    );
}
