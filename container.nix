{ pkgs, pkgsLinux, ... }:
let
  tele2aria = pkgsLinux.callPackage ./default.nix { };
in
pkgs.dockerTools.buildImage {
  name = "tele2aria";
  config = {
    Cmd = [ "${tele2aria}/bin/tele2aria" ];
  };
}
