{
  description = "Description for the project";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [ "x86_64-linux" ];
      perSystem = { system, pkgs, ... }: {
        devShells.default = pkgs.mkShell {
          GOPATH = "/home/readf0x/.config/go";
          packages = [
            pkgs.go
            pkgs.delve
            pkgs.oapi-codegen
          ];
        };
        packages = rec {
          qtbooru = pkgs.buildGoModule rec {
            name = "qtbooru";
            pname = name;
            version = "v0.3.4.1";

            src = ./.;

            vendorHash = "sha256-NHTKwUSIbNCUco88JbHOo3gt6S37ggee+LWNbHaRGEs=";

            # ldflags = [ "-X 'api.E621_URL=https://e621.net/posts.json'" "-X 'api.E926_URL=https://e926.net/posts.json'" ];

            meta = {
              description = "qt booru frontend";
              mainProgram = pname;
            };
          };
          default = qtbooru;
        };
      };
    };
}
