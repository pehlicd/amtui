/*

 █████╗ ███╗   ███╗████████╗██╗   ██╗██╗
██╔══██╗████╗ ████║╚══██╔══╝██║   ██║██║
███████║██╔████╔██║   ██║   ██║   ██║██║
██╔══██║██║╚██╔╝██║   ██║   ██║   ██║██║
██║  ██║██║ ╚═╝ ██║   ██║   ╚██████╔╝██║
╚═╝  ╚═╝╚═╝     ╚═╝   ╚═╝    ╚═════╝ ╚═╝

MIT License

Copyright (c) 2023 Furkan Pehlivan

*/

package main

import (
	"log"

	"github.com/pehlicd/amtui/pkg"
)

func main() {
	tui := pkg.InitTUI()

	if err := tui.Start(); err != nil {
		log.Fatalf("Error running app: %s", err)
	}
}
