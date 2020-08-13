// +build generate

package main

import (
	"fmt"
	"log"
	"os"
)

//go:generate go run manifest_gen.go
func main() {
	genIconManifest("manifest.txt")
}

func genIconManifest(file string) {
	w, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()
	fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
<assemblyIdentity
    version="1.0.0.0"
    processorArchitecture="x86"
    name="controls"
    type="win32"
></assemblyIdentity>
<dependency>
    <dependentAssembly>
        <assemblyIdentity
            type="win32"
            name="Microsoft.Windows.Common-Controls"
            version="6.0.0.0"
            processorArchitecture="*"
            publicKeyToken="6595b64144ccf1df"
            language="*"
        ></assemblyIdentity>
    </dependentAssembly>
</dependency>
</assembly>`)
}
