{ pkgs, ... }:

pkgs.buildGoApplication rec {
  pname = "tele2aria";
  version = "0.1.0";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;

  ldflags = [
    "-X github.com/adamgoose/tele2aria/cmd.Version=${version}"
  ];

  buildInputs = [
    pkgs.zlib
    pkgs.tdlib
    pkgs.openssl
  ];
}
