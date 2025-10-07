{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    systems.url = "github:nix-systems/default";
    
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = { flake-parts, systems, ... }@inputs:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import systems;

      perSystem = { pkgs, ... }: {
        devShells.default = pkgs.mkShellNoCC {
          packages = with pkgs; [
            go_latest
            gopls
            golangci-lint
          ];
        };
      };
    };
}
