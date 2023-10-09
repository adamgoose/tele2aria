{
  description = "Description for the project";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    devenv.url = "github:cachix/devenv";
    nix2container.url = "github:nlewo/nix2container";
    nix2container.inputs.nixpkgs.follows = "nixpkgs";
    mk-shell-bin.url = "github:rrbutani/nix-mk-shell-bin";
    gomod2nix.url = "github:nix-community/gomod2nix/v1.5.0";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  };

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.devenv.flakeModule
      ];
      systems = [ "x86_64-linux" "i686-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      perSystem = { config, self', inputs', pkgs, system, ... }: rec {
        _module.args.pkgs = import inputs.nixpkgs {
          inherit system;
          overlays = [
            inputs.gomod2nix.overlays.default
          ];
        };

        packages.default = packages.tele2aria;
        packages.tele2aria = pkgs.callPackage ./default.nix { };
        packages.container = pkgs.callPackage ./container.nix {
          pkgsLinux = import inputs.nixpkgs {
            system = "x86_64-linux";
            overlays = [
              inputs.gomod2nix.overlays.default
            ];
          };
        };

        devenv.shells.default = {
          dotenv.enable = true;
          languages.go.enable = true;

          packages = with pkgs; [
            air
            tdlib
            openssl
            gomod2nix
            (mkGoEnv { pwd = ./.; })
          ];

          scripts."build-container".exec = ''
            nix build --builders "$NIX_BUILDERS" .#container
            podman load < result
          '';
        };

      };
      flake = { };
    };
}
