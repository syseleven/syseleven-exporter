{
  description = "Syseleven-exporter, fetch syseleven specific metrics";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      supportedSystems = [
        "x86_64-linux"
        "x86_64-darwin"
        "aarch64-linux"
        "aarch64-darwin"
      ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
      version = builtins.substring 0 8 lastModifiedDate;

      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};

      bin =
        {
          version ? "dev",
          ...
        }:
        pkgs.buildGoModule {
          inherit version;
          pname = "syseleven-exporter";
          src = ./.;
          env = {
            CGO_ENABLED = 0;
          };

          ldflags =
            let
              repo = "github.com/syseleven/syseleven-exporter";
            in
            [
              "-s"
              "-w"
              "-X ${repo}/pkg/version.Version=${version}"
              "-X ${repo}/pkg/version.Revision=${self.rev or self.dirtyRev or "dirty"}"
              "-X ${repo}/pkg/version.Branch=HEAD"
              "-X ${repo}/pkg/version.BuildUser=root"
              "-X ${repo}/pkg/version.BuildDate=${toString (self.lastModified or 0)}"
            ];

          vendorHash = "sha256-loGLZ+6gQd9UNuXoE0LEuqGRy7BDFifUgiuSq7aa5vc=";
          doCheck = false;
        };

      oci =
        {
          name ? "syseleven-exporter",
          tag ? "latest",
          ...
        }:
        pkgs.dockerTools.buildLayeredImage {
          inherit name tag;
          contents = [
            pkgs.cacert
          ];
          config.Entrypoint = [
            "${
              pkgs.callPackage bin {
                version = if tag == "latest" then "dev" else pkgs.lib.removePrefix "v" tag;
              }
            }/bin/syselevenexporter"
          ];
        };
    in
    {
      packages.${system} = {
        default = pkgs.callPackage bin { };
        bin = pkgs.callPackage bin { };
        oci = pkgs.callPackage oci { };
      };

      devShells.${system} = {
        default = pkgs.mkShell {
          shellHook = ''
            echo "nixpkgs revision: ${nixpkgs.rev}"
            echo "    git revision: ${self.rev or self.dirtyRev or "dirty"}"
          '';
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            golangci-lint
            goreleaser
            go-junit-report
            gocover-cobertura
            curl
            kind
            kubectl
            docker
            trivy
            nixfmt-rfc-style
          ];
        };
      };
      formatter.${system} = nixpkgs.legacyPackages.${system}.nixfmt-tree;
    };
}
