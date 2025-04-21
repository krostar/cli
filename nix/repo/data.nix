{
  deps,
  lib,
  ...
}: let
  inherit (deps.harmony.result.data.harmony.ci.linters) golangci-lint;
in {
  ci.linters.golangci-lint.linters = {
    exclusions = {
      paths = ["internal/example"];
      rules = [
        {
          path = "cfg/source/env/|cfg/source/flag/";
          linters = ["depguard"];
          text = "import 'reflect' is not allowed";
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
}
