{
  description = "Description for the project";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [ "x86_64-linux" ];
      perSystem = { system, pkgs, ... }: let
          lib = pkgs.lib;
          libs = with pkgs; [ kdePackages.full ];
        in {
        devShells.default = pkgs.mkShell {
          GOPATH = "/home/readf0x/.config/go";
          packages = with pkgs; [
            go
            delve
            pkg-config
            libsForQt5.qt5.qtbase
          ];

          LD_LIBRARY_PATH = lib.makeLibraryPath libs;
          PKG_CONFIG_PATH = lib.makeSearchPath "lib/pkgconfig" libs;
        };
        packages = rec {
          qtbooru = pkgs.buildGoModule rec {
            name = "qtbooru";
            pname = name;
            version = "v1.0";

            LD_LIBRARY_PATH = lib.makeLibraryPath libs;
            PKG_CONFIG_PATH = lib.makeSearchPath "lib/pkgconfig" libs;

            src = ./.;

            vendorHash = "sha256-WrxgRYSXeXLJwsmNiRCBUCy7YdHIOO2eThhRX+qzI5g=";

            ldflags = [ "-s" "-w" ];
            nativeBuildInputs = [ pkgs.pkg-config ];

            meta = {
              description = "Qt Booru Client";
              mainProgram = pname;
            };
          };
          default = qtbooru;
        };
      };
    };
}
