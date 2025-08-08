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
          qtbooru = let
            miqt = pkgs.fetchFromGitHub {
              owner = "mappu";
              repo = "miqt";
              rev = "v0.11.1";
              hash = "sha256-crKIxCrWJrJTmF2JKPy3JpAyloaDoBxjxgxRZjjXHRc=";
            };
          in pkgs.buildGoModule rec {
            name = "qtbooru";
            pname = name;
            version = "v1.0";

            LD_LIBRARY_PATH = lib.makeLibraryPath libs;
            PKG_CONFIG_PATH = lib.makeSearchPath "lib/pkgconfig" libs;

            src = ./.;

            modBuildPhase = ''
              runHook preBuild

              if (( "''${NIX_DEBUG:-0}" >= 1 )); then
                goModVendorFlags+=(-v)
              fi
              go mod vendor
              cp -r ${miqt}/libmiqt vendor/github.com/mappu/miqt

              runHook postBuild
            '';

            vendorHash = "sha256-PLABFPIAbdxWSDfhX+rr4Xh+IOKKCmB6FYQJk9SCta4=";

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
