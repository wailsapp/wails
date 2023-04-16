// Generate takes a string and returns a number of hashes for it
export function Generate(input) {
    return wails.Plugin("hashes","Generate",input);
}