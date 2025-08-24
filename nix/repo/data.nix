{lib, ...}: {
  dev.git-cliff.enable = true;
  ci.testers.go.enable = true;
  ci.linters = {
    commitlint.enable = true;
    editorconfig-checker.settings.Exclude = ["./double/internal/*_generated*.go"];
    golangci-lint = {
      enable = true;
      linters = {
        exclusions = {
          paths = ["internal/example"];
          rules = [
            {
              path = "cfg/source/env/|cfg/source/flag/";
              linters = ["depguard"];
              text = "import 'reflect' is not allowed";
            }
            {
              path = "cfg/source/flag/flag.go";
              linters = ["gosec"];
              text = "Use of unsafe calls should be audited";
            }
            {
              path = "double/internal/generator/";
              linters = ["gosec" "goconst"];
              text = "Potential file inclusion via variable|make it a constant";
            }
          ];
        };

        settings = {
          iface.enable = lib.mkForce ["identical"];
          importas.alias = [
            {
              pkg = "github.com/krostar/cli/mapper/internal";
              alias = "mapper";
            }
            {
              pkg = "github.com/krostar/cli/cfg";
              alias = "clicfg";
            }
          ];
        };
      };
    };
  };
}
