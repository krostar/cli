{pkgs}:
pkgs.mkShellNoCC {
  nativeBuildInputs = with pkgs; [
    act
    alejandra
    deadnix
    gci
    git
    go_1_21
    gofumpt
    golangci-lint
    govulncheck
    shellcheck
    statix
    yamllint
  ];
}
