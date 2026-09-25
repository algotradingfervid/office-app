#!/usr/bin/env bash
# Import-boundary check (part of `make check`).
#   1. Shared core (internal/core/...) never imports a form (internal/forms/...).
#   2. A form never imports another form.
#   3. Only cmd/ and internal/testapp import internal/modules (the one module list).
# Go's own internal/ rule already stops anyone importing a module's internal/ packages.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
fail=0
while read -r pkg imports; do
  for imp in $imports; do
    case "$pkg" in
      officeapp/internal/core/*)
        case "$imp" in officeapp/internal/forms/*) echo "BOUNDARY: $pkg imports form $imp"; fail=1 ;; esac ;;
      officeapp/internal/forms/*)
        form="$(echo "$pkg" | cut -d/ -f4 | sed -E 's/(\.test|_test)$//')" # a form's test packages are the form
        case "$imp" in officeapp/internal/forms/*)
          other="$(echo "$imp" | cut -d/ -f4 | sed -E 's/(\.test|_test)$//')"
          if [ "$form" != "$other" ]; then echo "BOUNDARY: form $form imports form $other ($imp)"; fail=1; fi ;;
        esac ;;
    esac
    if [ "$imp" = "officeapp/internal/modules" ]; then
      case "$pkg" in officeapp/cmd/*|officeapp/internal/testapp) ;; *) echo "BOUNDARY: $pkg imports internal/modules"; fail=1 ;; esac
    fi
  done
done < <(go list -test -f '{{.ImportPath}} {{join .Imports " "}}' ./... | sed 's/ \[[^]]*\]//')
if [ "$fail" = 0 ]; then echo "imports: ok"; fi
exit "$fail"
