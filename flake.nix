{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils/v1.0.0";
  };

  description = "";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        build = pkgs.buildGoModule {
          pname = "gitlab-notifications";
          version = "0.2.2";
          modSha256 = pkgs.lib.fakeSha256;
          vendorHash = null;
	  src = ./.;
          subPackages = [ "cmd/gitlab-notifications" "cmd/gitlab-notifications-daemon" ];
        };
      in
      rec {
        packages = {
          gitlab-notifications = build;
          default = build;
        };

        devShells = {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              gopls
              golangci-lint
              gotools
              delve
            ];
          };
        };
      }
    );
}

