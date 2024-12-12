{
  deps,
  lib,
  ...
}: let
  inherit (deps.harmony.result.data.harmony.ci.linters) golangci-lint;
in {
  ci.linters.golangci-lint = lib.mkForce (lib.attrsets.recursiveUpdate golangci-lint {
    issues.exclude-dirs = ["internal/example"];

    linters-settings = {
      depguard.rules.all.deny = [
        {
          pkg = "github.com/pkg/errors";
          desc = "use go1.13 errors";
        }
      ];
      importas.alias = [
        {
          pkg = "github.com/google/go-cmp/cmp";
          alias = "gocmp";
        }
        {
          pkg = "github.com/google/go-cmp/cmp/cmpopts";
          alias = "gocmpopts";
        }
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
  });
}
