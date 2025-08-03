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
        devShells.default = let
          lib = pkgs.lib;
          libs = with pkgs; [ kdePackages.full ];
        in pkgs.mkShell {
          GOPATH = "/home/readf0x/.config/go";
          packages = with pkgs; [
            go
            delve
            pkg-config
          ];

          LD_LIBRARY_PATH = lib.makeLibraryPath libs;
          PKG_CONFIG_PATH = lib.makeSearchPath "lib/pkgconfig" libs;
          QML2_IMPORT_PATH = lib.makeSearchPath "lib/qt-6/qml" libs;
        };
        packages = rec {
          qtbooru = pkgs.buildGoModule rec {
            name = "qtbooru";
            pname = name;
            version = "indev_v0";

            src = ./.;

            vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";

            ldflags = [ "-s" "-w" ];

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
