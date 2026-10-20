go list "$1/..." | while IFS= read -r pkg; do
    printf '\n========== %s ==========\n\n' "
$pkg"
    go doc -all -u "$pkg"
done
